package models

type Card struct {
	Id       string `json:"id"`
	Last4    string `json:"last_4"`
	ExpMonth uint8  `json:"exp_month"`
	ExpYear  uint16 `json:"exp_year"`
	Name     string `json:"name"`
	Default  bool   `json:"default"`
}
