package common

type Response struct {
	Message string `json:"message"`
}

type Result struct {
	Data interface{} `json:"data"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

type Error struct {
	Message string `json:"errorMessage"`
}

type DeleteStatus struct {
	Deleted bool `json:"deleted"`
}
