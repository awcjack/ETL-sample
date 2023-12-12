package transformation

import "time"

// transformed (unified) data format
type TransformedData struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Address     StructuredAddress
}

type StructuredAddress struct {
	City          string
	StreetName    string
	StreetAddress string
	ZipCode       string
	State         string
	Country       string
	Latitude      float64
	Longitude     float64
}

type Transformer interface {
	Transform(data string) TransformedData
}
