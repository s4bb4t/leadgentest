package models

type Building struct {
	Title  string `json:"title" validate:"required,min=2,max=60,alphanumunicode"`
	City   string `json:"city" validate:"required,min=2,max=60,alpha"`
	Year   int    `json:"year" validate:"required"`
	Floors int    `json:"floors" validate:"required"`
}

type Buildings struct {
	Data []Building    `json:"data"`
	Meta BuildingsMeta `json:"meta"`
}

type BuildingsMeta struct {
	TotalAmount uint64 `json:"total_amount"`
	Query       Query  `json:"query,omitempty"`
}

type Query struct {
	City   string `json:"city,omitempty"`
	Year   int    `json:"year,omitempty"`
	Floors int    `json:"floors,omitempty"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
