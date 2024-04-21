package models

type Car struct {
	RegNum          string `json:"regNum"`
	Mark            string `json:"mark"`
	Model           string `json:"model"`
	Year            int    `json:"year"`
	OwnerName       string `json:"ownerName"`
	OwnerSurname    string `json:"ownerSurname"`
	OwnerPatronymic string `json:"ownerPatronymic"`
}
