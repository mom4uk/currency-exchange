package domain

type Currency struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Sign string `json:"sign"`
}
