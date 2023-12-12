package http

import (
	"encoding/json"
	"time"

	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/utils"
)

// JSON response format in random data api
// comment useless field to only decode JSON to required field
type randomDataAPIResponse struct {
	// Id int
	// Uid string
	// password string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// UserName string
	// Email string
	// Avatar string
	// Gender string
	// PhoneNumber string `json:"phone_number"`
	// SocialInsuranceNumber string `json:"social_insurance_number"`
	DateOfBirth string `json:"date_of_birth"`
	// Employment randomDateAPIResponseEmployment
	Address randomDataAPIResponseAddress
	// CreditCard   randomDataAPIResponseCreditCard `json:"credit_card"`
	// Subscription randomDataAPIResponseSubscription
}

// type randomDateAPIResponseEmployment struct {
// 	Title    string
// 	KeySkill string `json:"key_skill"`
// }

type randomDataAPIResponseAddress struct {
	City          string
	StreetName    string `json:"street_name"`
	StreetAddress string `json:"street_address"`
	ZipCode       string `json:"zip_code"`
	State         string
	Country       string
	Coordinates   randomDataAPIResponseCoordinates
}

type randomDataAPIResponseCoordinates struct {
	Lat float64
	Lng float64
}

// type randomDataAPIResponseCreditCard struct {
// 	CCNumber string `json:"cc_number"`
// }

// type randomDataAPIResponseSubscription struct {
// 	Plan          string
// 	Status        string
// 	PaymentMethod string `json:"payment_method"`
// 	Term          string
// }

type RandomDataAPITransformer struct {
	logger utils.Logger
}

func NewRandomDataAPITransformer(logger utils.Logger) *RandomDataAPITransformer {
	return &RandomDataAPITransformer{
		logger: logger,
	}
}

// parse JSON from random data api
func (r *RandomDataAPITransformer) Transform(rawData []byte) (transformation.TransformedData, error) {
	r.logger.Debugf("random data API rawData", rawData)

	var structedData *randomDataAPIResponse

	err := json.Unmarshal(rawData, &structedData)
	if err != nil {
		return transformation.TransformedData{}, err
	}

	// parse data of birth from string to time
	date, err := time.Parse(time.DateOnly, structedData.DateOfBirth)
	if err != nil {
		return transformation.TransformedData{}, err
	}

	return transformation.TransformedData{
		FirstName:   structedData.FirstName,
		LastName:    structedData.LastName,
		DateOfBirth: date,
		Address: transformation.StructuredAddress{
			City:          structedData.Address.City,
			StreetName:    structedData.Address.StreetName,
			StreetAddress: structedData.Address.StreetAddress,
			ZipCode:       structedData.Address.ZipCode,
			State:         structedData.Address.State,
			Country:       structedData.Address.Country,
			Latitude:      structedData.Address.Coordinates.Lat,
			Longitude:     structedData.Address.Coordinates.Lng,
		},
	}, nil
}
