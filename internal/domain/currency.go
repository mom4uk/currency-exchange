package domain

type Currency struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}
