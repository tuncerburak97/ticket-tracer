package train

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"ticker-tracer/client/notification/mail"
	emailModel "ticker-tracer/client/notification/mail/model"
	"ticker-tracer/client/tcdd"
	tcddClientCommonModel "ticker-tracer/client/tcdd/model/common"
	tcddClientRequest "ticker-tracer/client/tcdd/model/request"
	tcddClientResponse "ticker-tracer/client/tcdd/model/response"
	"ticker-tracer/model/entity"
	"ticker-tracer/repository"
	"time"
)

type TrainScheduler struct {
	tcddClient          *tcdd.TcddHttpClient
	mailClient          *mail.MailHttpClient
	stations            *tcddClientResponse.StationLoadResponse
	once                sync.Once
	mu                  sync.Mutex
	isZeroRequestLogged bool
}

var trainSchedulerInstance *TrainScheduler

func GetTrainSchedulerInstance() *TrainScheduler {
	if trainSchedulerInstance == nil {
		trainSchedulerInstance = NewTrainScheduler(tcdd.GetTcddHttpClientInstance(),
			mail.GetMailHttpClientInstance())
		trainSchedulerInstance.isZeroRequestLogged = false
	}
	return trainSchedulerInstance

}

func NewTrainScheduler(tcddClient *tcdd.TcddHttpClient,
	mailClient *mail.MailHttpClient,

) *TrainScheduler {
	return &TrainScheduler{
		tcddClient: tcddClient,
		mailClient: mailClient,
	}
}

func (ts *TrainScheduler) getStations() (*tcddClientResponse.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		stationLoadRequest := tcddClientRequest.StationLoadRequest{
			Language:    0,
			ChannelCode: "3",
			Date:        "Nov 10, 2011 12:00:00 AM",
			SalesQuery:  true,
		}
		ts.stations, err = ts.tcddClient.LoadAllStation(stationLoadRequest)
	})
	return ts.stations, err
}

func (ts *TrainScheduler) Run() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var ticketRequestRepository = repository.GetTicketRequestRepository()
	var pendingRequests, err = ticketRequestRepository.FindByStatus("PENDING")
	if err != nil {
		log.Printf("Error getting pending requests: %v", err)
		return
	}

	if len(pendingRequests) == 0 && !ts.isZeroRequestLogged {
		log.Printf("No pending request found")
		ts.isZeroRequestLogged = true
		return
	}

	if len(pendingRequests) == 0 {
		return
	}

	log.Printf("Running train scheduler with %d requests", len(pendingRequests))

	if _, err := ts.getStations(); err != nil {
		log.Printf("Error getting stations: %v", err)
		return
	}
	var foundedRequestIDList = make([]string, 0)

	for _, searchTrainRequest := range pendingRequests {
		foundedRequestIDS := ts.processRequest(searchTrainRequest)
		if foundedRequestIDS != "" {
			foundedRequestIDList = append(foundedRequestIDList, foundedRequestIDS)
		}
	}

	var filterFoundedRequests = make([]entity.TicketRequest, 0)
	for _, foundedRequestID := range foundedRequestIDList {
		for _, request := range pendingRequests {
			if request.ID == foundedRequestID {
				filterFoundedRequests = append(filterFoundedRequests, request)
			}
		}
	}

	ts.UpdateTicketRequestStatusToFound(filterFoundedRequests)
}

func (ts *TrainScheduler) processRequest(request entity.TicketRequest) (requestID string) {

	criteria := tcddClientRequest.Criteria{
		SalesChannel:       3,
		DepartureStation:   request.DepartureStation,
		IsMapDeparture:     false,
		ArrivalStation:     request.ArrivalStation,
		IsMapArrival:       false,
		DepartureDate:      request.DepartureDate,
		IsRegional:         false,
		OperationType:      0,
		PassengerCount:     1,
		IsTransfer:         true,
		DepartureStationID: request.DepartureStationID,
		ArrivalStationID:   request.ArrivalStationID,
		TravelType:         1,
	}

	search, err := ts.tcddClient.TripSearch(tcddClientRequest.TripSearchRequest{
		ChannelCode: 3,
		Language:    0,
		Criteria:    criteria,
	})
	if err != nil {
		log.Printf("Error searching trip: %v", err)
		return
	}

	b := search.TripSearchResponseInfo.ResponseCode != "000"
	if b {
		log.Printf("Error searching trip: %v", search.TripSearchResponseInfo.ResponseMsg)
		return
	}
	remainingDisabledNumber, found := ts.findTrip(search, request.TourID)
	if found {
		var requestId = ts.handleFoundTrip(request, int(remainingDisabledNumber), search.SearchResult[0].ArrivalDate)
		if requestId != "" {
			return requestId
		}
		return ts.handleNotFoundTrip(request)

	}

	return ts.handleNotFoundTrip(request)

}
func (ts *TrainScheduler) handleFoundTrip(request entity.TicketRequest, remainingDisabledNumber int, arrivalDate string) (requestID string) {

	placeSearch, err := ts.tcddClient.StationEmptyPlaceSearch(tcddClientRequest.StationEmptyPlaceSearchRequest{
		ChannelCode:   "3",
		Language:      0,
		TourTitleID:   request.TourID,
		DepartureStID: request.DepartureStationID,
		ArrivalStID:   int(request.ArrivalStationID),
	})
	if err != nil {
		log.Printf("Error getting empty place: %v", err)
		return
	}

	totalEmptyPlace := calculateTotalEmptyPlace(placeSearch.EmptyPlaceList)
	availablePlace := totalEmptyPlace - remainingDisabledNumber
	if availablePlace > 0 {

		locationSelectionWagonRequestList := getLocationSelectionWagonRequestList(placeSearch.EmptyPlaceList, request)
		reservedSeats := ts.reserveSeat(locationSelectionWagonRequestList, request)

		departureValidation := true
		departureDateFormat, err := time.Parse("Jan 02, 2006 03:04:05 PM", request.DepartureDate)
		if err != nil {
			fmt.Println("Departure date parse edilemedi:", err)
			fmt.Println("Departure Date:", request.DepartureDate)
			departureValidation = false
		}

		arrivalValidation := true
		arrivalDateFormat, err := time.Parse("Jan 02, 2006 03:04:05 PM", request.ArrivalDate)
		if err != nil {
			fmt.Println("Arrival Date parse edilemedi:", err)
			fmt.Println("Arrival Date:", request.ArrivalDate)
			arrivalValidation = false
		}

		// Türkçe tarih formatını oluşturma ve yazdırma
		var departureDateStr string
		if departureValidation {
			departureDateStr = formatTurkishDate(departureDateFormat)
		} else {
			departureDateStr = request.DepartureDate
		}

		var arrivalDateStr string
		if arrivalValidation {
			arrivalDateStr = formatTurkishDate(arrivalDateFormat)
		} else {
			arrivalDateStr = request.ArrivalDate
		}

		log.Printf("Found trip for request: %s and email: %s date: %s from: %s to: %s",
			request.ID,
			request.Email,
			request.DepartureDate,
			request.DepartureStation,
			request.ArrivalStation)

		sendEmail(
			request.Email,
			availablePlace,
			departureDateStr,
			arrivalDateStr,
			request.DepartureStation,
			request.ArrivalStation,
			reservedSeats)
		return request.ID
	}

	return ""
}

func (ts *TrainScheduler) reserveSeat(
	locationSelectionWagonRequestList []tcddClientRequest.LocationSelectionWagonRequest,
	request entity.TicketRequest,
) []tcddClientCommonModel.ReserveSeatDetail {

	reservedSeats := make([]tcddClientCommonModel.ReserveSeatDetail, 0)
	totalReservedSeat := 0

	for _, locationSelectionWagonRequest := range locationSelectionWagonRequestList {
		seats := ts.processWagonRequest(locationSelectionWagonRequest, request, &totalReservedSeat)
		if seats != nil {
			reservedSeats = append(reservedSeats, seats...)
		}
	}

	return reservedSeats
}

func (ts *TrainScheduler) processWagonRequest(
	locationSelectionWagonRequest tcddClientRequest.LocationSelectionWagonRequest,
	request entity.TicketRequest,
	totalReservedSeat *int,
) []tcddClientCommonModel.ReserveSeatDetail {

	reservedSeats := make([]tcddClientCommonModel.ReserveSeatDetail, 0)

	locationSelectionWagonResponse, err := ts.tcddClient.LocationSelectionWagon(locationSelectionWagonRequest)
	if err != nil {
		log.Printf("Error selecting wagon: %v", err)
		return nil
	}
	if locationSelectionWagonResponse.ResponseInfo.ResponseCode != "000" {
		log.Printf("Error selecting wagon: %v", locationSelectionWagonResponse.ResponseInfo.ResponseMsg)
		return nil
	}
	for _, locationSelectionWagon := range locationSelectionWagonResponse.LocationSelectionWagonResponseData.SeatInformationList {
		if locationSelectionWagon.Status == 0 {
			if *totalReservedSeat >= 3 {
				break
			}

			checkSeatRequest := tcddClientRequest.CheckSeatRequest{
				ChannelCode:             "3",
				Language:                0,
				SelectedSeatWagonNumber: locationSelectionWagon.WagonOrderNo,
				SelectedSeatNumber:      locationSelectionWagon.SeatNo,
				TourId:                  strconv.FormatInt(request.TourID, 10),
			}
			checkSeatResponse, err := ts.tcddClient.CheckSeat(checkSeatRequest)
			if err != nil {
				log.Printf("Error reserving seat: %v", err)
				return nil
			}
			if checkSeatResponse.ResponseInfo.ResponseCode != "000" {
				log.Printf("Error reserving seat: %v", checkSeatResponse.ResponseInfo.ResponseMsg)
				return nil
			}

			reserveSeatRequest := tcddClientRequest.ReserveSeatRequest{
				ChannelCode:        "3",
				Language:           0,
				TourID:             int(request.TourID),
				WagonOrder:         locationSelectionWagon.WagonOrderNo,
				SeatNo:             locationSelectionWagon.SeatNo,
				Gender:             "M",
				ArrivalStationID:   int(request.ArrivalStationID),
				DepartureStationID: int(request.DepartureStationID),
				Minute:             10,
				Huawei:             false,
			}

			reserveSeatResponse, err := ts.tcddClient.ReserveSeat(reserveSeatRequest)
			if err != nil {
				log.Printf("Error reserving seat: %v", err)
				return nil
			}
			if reserveSeatResponse.ResponseInfo.ResponseCode != "000" {
				log.Printf("Error reserving seat: %v", reserveSeatResponse.ResponseInfo.ResponseMsg)
				return nil
			}

			log.Printf("Seat reserved for request: %s, Email: %s, Date: %s, From: %s, To: %s",
				request.ID,
				request.Email,
				request.DepartureDate,
				request.DepartureStation,
				request.ArrivalStation)

			reservedSeats = append(reservedSeats, tcddClientCommonModel.ReserveSeatDetail{
				SeatNo:       locationSelectionWagon.SeatNo,
				WagonOrderNo: locationSelectionWagon.WagonOrderNo,
			})
			*totalReservedSeat++
		}
	}

	return reservedSeats
}

func (ts *TrainScheduler) findTrip(search *tcddClientResponse.TripSearchResponse, tourID int64) (int64, bool) {
	for _, trip := range search.SearchResult {
		if trip.TourID == tourID {
			if len(trip.WagonTypesEmptyPlace) > 0 {
				return trip.WagonTypesEmptyPlace[0].RemainingDisabledNumber, true
			}
		}
	}
	return 0, false
}

func calculateTotalEmptyPlace(emptyPlaceList []tcddClientResponse.EmptyPlace) int {
	totalEmptyPlace := 0
	for _, emptyPlace := range emptyPlaceList {
		totalEmptyPlace += emptyPlace.EmptyPlace
	}
	return totalEmptyPlace
}

func sendEmail(recipient string,
	availablePlace int,
	departureDate string,
	arrivalDate string,
	departureStation string,
	arrivalStation string,
	reservedSeats []tcddClientCommonModel.ReserveSeatDetail,
) {

	{
		body := fmt.Sprintf(`
<html>
<head>
<style>
table {
  font-family: Arial, sans-serif;
  border-collapse: collapse;
  width: 100%%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}

.margin-top {
  margin-top: 20px;
}
</style>
</head>
<body>
<p>Merhaba,</p>
<p>Aradığınız trende boş yer bulundu. &#128522;</p>
<table>
  <tr>
    <th>Kalan Boş Yer Sayısı</th>
    <th>Kalkış Zamanı</th>
    <th>Varış Zamanı</th>
    <th>Kalkış İstasyonu</th>
    <th>Varış İstasyonu</th>
  </tr>
  <tr>
    <td>%d</td>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
    <td>%s</td>
  </tr>
</table>

<div class="margin-top">
<p>Sizin için aşağıdaki koltuklar rezerv edilmiştir. 10 dakika boyunca koltuk diğer kullanıcılar için görünür olmayacaktır. Bu maili aldıktan 10 dakika sonra koltuk kilidi kalkmış olacaktır. İlgili koltuğu 10 dakika sonra kontrol edebilirsiniz!</p>
<table>
  <tr>
    <th>Vagon No</th>
    <th>Koltuk No</th>
  </tr>
`, availablePlace, departureDate, arrivalDate, departureStation, arrivalStation)

		for _, seat := range reservedSeats {
			body += fmt.Sprintf(`
  <tr>
    <td>%d</td>
    <td>%s</td>
  </tr>
`, seat.WagonOrderNo, seat.SeatNo)
		}

		body += `
</table>
</div>

<p>Tekrardan bu yolculuga dair bildirimleri takip etmek isterseniz uygulama üzerinden aynı talebi oluşturabilirsiniz</p>
<p>İyi yolculuklar dileriz!</p>
</body>
</html>`

		email := emailModel.Email{
			To:      recipient,
			Subject: "Tren Bilet Uyarısı",
			Body:    body,
		}

		// Send the email
		err := trainSchedulerInstance.mailClient.SendEmail(email)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}
		log.Printf("Email sent to: %s", recipient)
	}
}

func (ts *TrainScheduler) UpdateTicketRequestStatusToFound(foundedRequests []entity.TicketRequest) {

	var ticketRequestRepository = repository.GetTicketRequestRepository()
	now := time.Now()
	for _, request := range foundedRequests {
		request.Status = "FOUND"
		request.UpdatedAt = now
		totalAttempt := request.TotalAttempt
		request.TotalAttempt = totalAttempt + 1
		err := ticketRequestRepository.Update(&request)
		if err != nil {
			log.Printf("Error updating ticket request status to found: %v", err)
		}
	}
}

func getLocationSelectionWagonRequestList(emptyPlaceList []tcddClientResponse.EmptyPlace, request entity.TicketRequest) []tcddClientRequest.LocationSelectionWagonRequest {
	response := make([]tcddClientRequest.LocationSelectionWagonRequest, 0)
	for _, emptyPlace := range emptyPlaceList {
		if emptyPlace.EmptyPlace > 0 {
			response = append(response, tcddClientRequest.LocationSelectionWagonRequest{
				ChannelCode:          "3",
				Language:             0,
				TourTitleID:          strconv.FormatInt(request.TourID, 10),
				WagonOrderNo:         emptyPlace.WagonOrderNo,
				DepartureStationName: request.DepartureStation,
				ArrivalStationName:   request.ArrivalStation,
			})
		}
	}
	return response
}

func formatTurkishDate(t time.Time) string {
	// Ay isimlerini Türkçe karşılıklarıyla değiştirin
	months := map[string]string{
		"January":   "Ocak",
		"February":  "Şubat",
		"March":     "Mart",
		"April":     "Nisan",
		"May":       "Mayıs",
		"June":      "Haziran",
		"July":      "Temmuz",
		"August":    "Ağustos",
		"September": "Eylül",
		"October":   "Ekim",
		"November":  "Kasım",
		"December":  "Aralık",
	}

	// Günü, ayı, yılı, saati ve dakikayı formatlayın
	day := t.Day()
	month := months[t.Month().String()]
	year := t.Year()
	hour := t.Hour()
	minute := t.Minute()

	// Türkçe formatta string oluşturma
	return fmt.Sprintf("%02d-%s-%d %02d:%02d", day, month, year, hour, minute)
}

func (ts *TrainScheduler) handleNotFoundTrip(request entity.TicketRequest) (requestID string) {

	/*
		log.Printf("Trip not found for request: %s and email: %s date: %s from: %s to: %s",
			request.ID,
			request.Email,
			request.DepartureDate,
			request.DepartureStation,
			request.ArrivalStation)


	*/
	totalAttempt := request.TotalAttempt
	request.TotalAttempt = totalAttempt + 1
	now := time.Now()
	request.UpdatedAt = now

	var ticketRequestRepository = repository.GetTicketRequestRepository()
	err := ticketRequestRepository.Update(&request)
	if err != nil {
		log.Printf("Error updating ticket request: %v", err)
	}

	return ""
}
