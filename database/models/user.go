package models

type User struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Date    string `json:"date"`
	City    string `json:"city"`
	Country string `json:"country"`
}
