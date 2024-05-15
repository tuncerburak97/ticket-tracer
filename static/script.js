document.addEventListener('DOMContentLoaded', function () {
    const fromStationSelect = document.getElementById('from-station');
    const toStationSelect = document.getElementById('to-station');
    const departureDateInput = document.getElementById('departure-date');
    const searchBtn = document.getElementById('search-btn');
    const modal = document.getElementById('modal');
    const modalContent = document.querySelector('.modal-content');
    const closeBtn = document.querySelector('.close-btn');
    const trainDetailsContainer = document.getElementById('train-details');
    const addBtn = document.getElementById('add-btn');
    const loading = document.getElementById('loading');

    // Success modal elements
    const successModal = document.getElementById('success-modal');
    const successCloseBtn = document.querySelector('.success-close-btn');

    // Eklenen email input alanı ve hata mesajı
    const emailInputDiv = document.createElement('div');
    emailInputDiv.className = 'form-group';
    emailInputDiv.innerHTML = `<label for="email">Email</label> <input type="email" id="email" required> <small id="email-error" style="color: red; display: none;">Geçerli bir email adresi girin.</small>`;
    modalContent.insertBefore(emailInputDiv, addBtn);

    fetch('/tcdd/load')
        .then(response => response.json())
        .then(data => {
            data.response.forEach(station => {
                const option = document.createElement('option');
                option.value = station.stationID;
                option.text = station.stationName;
                option.dataset.toStations = JSON.stringify(station.toStationList);
                fromStationSelect.add(option);
            });
        })
        .catch(error => {
            console.error('Error:', error);
            alert('İstasyonlar yüklenirken hata oluştu.');
        });

    fromStationSelect.addEventListener('change', function () {
        const selectedOption = fromStationSelect.options[fromStationSelect.selectedIndex];
        const toStations = JSON.parse(selectedOption.dataset.toStations || '[]');
        toStationSelect.innerHTML = '<option value="">İstasyon Seç</option>';
        toStationSelect.disabled = toStations.length === 0;
        toStations.forEach(toStation => {
            const option = document.createElement('option');
            option.value = toStation.toStationId;
            option.text = toStation.toStationName;
            toStationSelect.add(option);
        });
    });

    searchBtn.addEventListener('click', function () {
        const stationID = fromStationSelect.value;
        const toStationID = toStationSelect.value;
        const departureDate = departureDateInput.value;
        const stationName = fromStationSelect.options[fromStationSelect.selectedIndex].text;
        const toStationName = toStationSelect.options[toStationSelect.selectedIndex].text;
        if (!stationID || !toStationID || !departureDate) {
            alert('Lütfen tüm alanları doldurun.');
            return;
        }
        const formattedDate = formatDateForServer(departureDate);
        const queryRequest = {
            gidisTarih: formattedDate,
            binisIstasyonId: parseInt(stationID),
            inisIstasyonId: parseInt(toStationID),
            binisIstasyon: stationName,
            inisIstasyonu: toStationName
        };
        loading.style.display = "flex";
        fetch('/tcdd/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(queryRequest),
        })
            .then(response => response.json())
            .then(result => {
                showModal(result.details);
                loading.style.display = "none";
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Sorgulama sırasında bir hata oluştu.');
                loading.style.display = "none";
            });
    });

    closeBtn.addEventListener('click', function () {
        modal.style.display = "none";
    });

    window.onclick = function (event) {
        if (event.target == modal) {
            modal.style.display = "none";
        } else if (event.target == successModal) {
            successModal.style.display = "none";
        }
    }

    let selectedTrains = [];

    trainDetailsContainer.addEventListener('click', function (event) {
        if (event.target.classList.contains('train-detail') && !event.target.classList.contains('disabled')) {
            const trainID = event.target.dataset.trainId;
            const trainIndex = selectedTrains.indexOf(trainID);
            if (trainIndex > -1) {
                selectedTrains.splice(trainIndex, 1);
                event.target.classList.remove('selected');
            } else if (selectedTrains.length < 3) {
                selectedTrains.push(trainID);
                event.target.classList.add('selected');
            } else {
                alert('En fazla 3 sefer seçebilirsiniz.');
            }
        }
    });

    addBtn.addEventListener('click', function () {
        const email = document.getElementById('email').value;
        const emailError = document.getElementById('email-error');

        // Email regex kontrolü
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            emailError.style.display = 'block';
            return;
        } else {
            emailError.style.display = 'none';
        }

        if (selectedTrains.length === 0) {
            alert('Lütfen en az bir sefer seçin.');
            return;
        }

        const requestPayload = {
            request: selectedTrains.map(trainID => {
                const detail = trainDetailsContainer.querySelector(`[data-train-id="${trainID}"]`);
                return {
                    trainID: parseInt(trainID),
                    tourID: parseInt(detail.dataset.tourId),
                    gidisTarih: formatDateForServer(detail.dataset.departureDate),
                    binisIstasyonId: parseInt(detail.dataset.departureStationId),
                    inisIstasyonId: parseInt(detail.dataset.arrivalStationId),
                    binisIstasyon: detail.dataset.departureStationName,
                    inisIstasyonu: detail.dataset.arrivalStationName,
                    email
                };
            })
        };

        fetch('/tcdd/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestPayload),
        })
            .then(response => response.text())
            .then(result => {
                modal.style.display = "none";
                successModal.style.display = "block";
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Sefer eklerken bir hata oluştu.');
            });
    });

    successCloseBtn.addEventListener('click', function () {
        successModal.style.display = "none";
    });

    function showModal(details) {
        trainDetailsContainer.innerHTML = '';
        selectedTrains = []; // Modal açıldığında seçili trenlerin sıfırlanması

        details.forEach(detail => {
            const formattedDepartureDate = formatTimeForDisplay(detail.departureDate);
            const formattedArrivalDate = formatTimeForDisplay(detail.arrivalDate);
            const normalSeatsCount = detail.emptyPlace.normalPeopleEmptyPlaceCount;

            const trainDetailDiv = document.createElement('div');
            trainDetailDiv.classList.add('train-detail');
            if (normalSeatsCount > 0) {
                trainDetailDiv.classList.add('disabled');
            }
            trainDetailDiv.dataset.trainId = detail.trainID;
            trainDetailDiv.dataset.tourId = detail.tourID;
            trainDetailDiv.dataset.departureDate = detail.departureDate;
            trainDetailDiv.dataset.departureStationId = detail.departureStationID;  // Fix the dataset field
            trainDetailDiv.dataset.arrivalStationId = detail.arrivalStationID;      // Fix the dataset field
            trainDetailDiv.dataset.departureStationName = detail.departureStation;
            trainDetailDiv.dataset.arrivalStationName = detail.arrivalStation;
            trainDetailDiv.innerHTML = `
                <h3>${detail.departureStation} - ${detail.arrivalStation}</h3>
                <p>Gidiş: ${formattedDepartureDate}</p>
                <p>Varış: ${formattedArrivalDate}</p>
                <div class="seat-info">
                    <span class="seat-total">Boş Koltuklar: ${detail.emptyPlace.totalEmptyPlaceCount}</span>
                    <span class="seat-disabled">Engelli Koltuklar: ${detail.emptyPlace.disabledPlaceCount}</span>
                    ${normalSeatsCount === 0 ? `<span class="seat-normal">Normal Koltuklar: ${normalSeatsCount}</span>` : `<span class="seat-normal disabled">Normal Koltuklar: ${normalSeatsCount}</span>`}
                </div>
            `;
            trainDetailsContainer.appendChild(trainDetailDiv);
        });
        modal.style.display = "block";
    }

    function formatDateForServer(date) {
        const dateObj = new Date(date);
        const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        const day = ('0' + dateObj.getDate()).slice(-2); // 2 basamaklı gün
        const month = monthNames[dateObj.getMonth()]; // Ay ismi
        const year = dateObj.getFullYear(); // Yıl
        const hours = ('0' + (dateObj.getHours() % 12 || 12)).slice(-2); // 12 saat formatı
        const minutes = ('0' + dateObj.getMinutes()).slice(-2); // 2 basamaklı dakika
        const seconds = '00'; // Saniye
        const period = dateObj.getHours() >= 12 ? 'PM' : 'AM'; // AM or PM
        return `${month} ${day}, ${year} ${hours}:${minutes}:${seconds} ${period}`;
    }

    function formatTimeForDisplay(dateTime) {
        const dateObj = new Date(dateTime);
        const hours = ('0' + dateObj.getHours()).slice(-2);
        const minutes = ('0' + dateObj.getMinutes()).slice(-2);
        return `${hours}:${minutes}`;
    }
});