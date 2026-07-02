package models

type ClientError struct {
	MessageToUser string `json:"messageToUser"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
}
