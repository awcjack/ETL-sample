package http_test

import (
	"testing"
	"time"

	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/transformation/http"
	"github.com/sirupsen/logrus"
)

func TestTransform(t *testing.T) {
	dep := newRandomDataAPITransformerDependencies()

	type testcase struct {
		testcase       string
		rawData        []byte
		expectedResult transformation.TransformedData
		expectedError  bool
		isPanic        bool
	}

	dob, _ := time.Parse(time.DateOnly, "1981-08-30")

	testcases := []testcase{
		{
			testcase: "Normal",
			rawData:  []byte("{\"id\":596,\"uid\":\"96bedfef-4de2-4b5f-8cb0-adafc1b8fdce\",\"password\":\"IfuewXQ4P0\",\"first_name\":\"Chasidy\",\"last_name\":\"Kirlin\",\"username\":\"chasidy.kirlin\",\"email\":\"chasidy.kirlin@email.com\",\"avatar\":\"https://robohash.org/autodiodolorem.png?size=300x300\u0026set=set1\",\"gender\":\"Polygender\",\"phone_number\":\"+269 460.093.9024\",\"social_insurance_number\":\"357134402\",\"date_of_birth\":\"1981-08-30\",\"employment\":{\"title\":\"Administration Assistant\",\"key_skill\":\"Problem solving\"},\"address\":{\"city\":\"Marionland\",\"street_name\":\"Domingo Green\",\"street_address\":\"82204 Wisoky Canyon\",\"zip_code\":\"43072-8812\",\"state\":\"Washington\",\"country\":\"United States\",\"coordinates\":{\"lat\":-50.65341353032217,\"lng\":-93.89954802799431}},\"credit_card\":{\"cc_number\":\"4403-8715-0240-9153\"},\"subscription\":{\"plan\":\"Silver\",\"status\":\"Active\",\"payment_method\":\"Visa checkout\",\"term\":\"Annual\"}}"),
			expectedResult: transformation.TransformedData{
				FirstName:   "Chasidy",
				LastName:    "Kirlin",
				DateOfBirth: dob,
				Address: transformation.StructuredAddress{
					City:          "Marionland",
					StreetName:    "Domingo Green",
					StreetAddress: "82204 Wisoky Canyon",
					ZipCode:       "43072-8812",
					State:         "Washington",
					Country:       "United States",
					Latitude:      -50.65341353032217,
					Longitude:     -93.89954802799431,
				},
			},
			expectedError: false,
			isPanic:       false,
		},
		{
			testcase: "Missing field",
			rawData:  []byte("{\"id\":596,\"uid\":\"96bedfef-4de2-4b5f-8cb0-adafc1b8fdce\",\"password\":\"IfuewXQ4P0\",\"last_name\":\"Kirlin\",\"username\":\"chasidy.kirlin\",\"email\":\"chasidy.kirlin@email.com\",\"avatar\":\"https://robohash.org/autodiodolorem.png?size=300x300\u0026set=set1\",\"gender\":\"Polygender\",\"phone_number\":\"+269 460.093.9024\",\"social_insurance_number\":\"357134402\",\"date_of_birth\":\"1981-08-30\",\"employment\":{\"title\":\"Administration Assistant\",\"key_skill\":\"Problem solving\"},\"address\":{\"city\":\"Marionland\",\"street_name\":\"Domingo Green\",\"street_address\":\"82204 Wisoky Canyon\",\"zip_code\":\"43072-8812\",\"state\":\"Washington\",\"country\":\"United States\",\"coordinates\":{\"lat\":-50.65341353032217,\"lng\":-93.89954802799431}},\"credit_card\":{\"cc_number\":\"4403-8715-0240-9153\"},\"subscription\":{\"plan\":\"Silver\",\"status\":\"Active\",\"payment_method\":\"Visa checkout\",\"term\":\"Annual\"}}"),
			expectedResult: transformation.TransformedData{
				FirstName:   "",
				LastName:    "Kirlin",
				DateOfBirth: dob,
				Address: transformation.StructuredAddress{
					City:          "Marionland",
					StreetName:    "Domingo Green",
					StreetAddress: "82204 Wisoky Canyon",
					ZipCode:       "43072-8812",
					State:         "Washington",
					Country:       "United States",
					Latitude:      -50.65341353032217,
					Longitude:     -93.89954802799431,
				},
			},
			expectedError: false,
			isPanic:       false,
		},
		{
			testcase:       "Missing date of birth field",
			rawData:        []byte("{\"id\":596,\"uid\":\"96bedfef-4de2-4b5f-8cb0-adafc1b8fdce\",\"password\":\"IfuewXQ4P0\",\"first_name\":\"Chasidy\",\"last_name\":\"Kirlin\",\"username\":\"chasidy.kirlin\",\"email\":\"chasidy.kirlin@email.com\",\"avatar\":\"https://robohash.org/autodiodolorem.png?size=300x300\u0026set=set1\",\"gender\":\"Polygender\",\"phone_number\":\"+269 460.093.9024\",\"social_insurance_number\":\"357134402\",\"employment\":{\"title\":\"Administration Assistant\",\"key_skill\":\"Problem solving\"},\"address\":{\"city\":\"Marionland\",\"street_name\":\"Domingo Green\",\"street_address\":\"82204 Wisoky Canyon\",\"zip_code\":\"43072-8812\",\"state\":\"Washington\",\"country\":\"United States\",\"coordinates\":{\"lat\":-50.65341353032217,\"lng\":-93.89954802799431}},\"credit_card\":{\"cc_number\":\"4403-8715-0240-9153\"},\"subscription\":{\"plan\":\"Silver\",\"status\":\"Active\",\"payment_method\":\"Visa checkout\",\"term\":\"Annual\"}}"),
			expectedResult: transformation.TransformedData{},
			expectedError:  true,
			isPanic:        false,
		},
		{
			testcase:       "Empty data",
			rawData:        []byte("{}"),
			expectedResult: transformation.TransformedData{},
			expectedError:  true,
			isPanic:        false,
		},
		{
			testcase:       "Empty data",
			rawData:        nil,
			expectedResult: transformation.TransformedData{},
			expectedError:  true,
			isPanic:        false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testcase, func(t *testing.T) {
			v, err := dep.randomDataAPITransformerHandler.Transform(tc.rawData)
			isPanic := false
			defer func() {
				if err := recover(); err != nil {
					isPanic = true
				}
			}()
			if tc.isPanic != isPanic {
				t.Errorf("expected panic %v but panic %v", tc.isPanic, isPanic)
			}
			if tc.expectedError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("not expected error, but got %v", err)
			}

			if v != tc.expectedResult {
				t.Errorf("expected %v, but got %v", tc.expectedResult, v)
			}
		})
	}
}

type randomDataAPITransformerDependencies struct {
	randomDataAPITransformerHandler *http.RandomDataAPITransformer
}

func newRandomDataAPITransformerDependencies() randomDataAPITransformerDependencies {
	logger := logrus.NewEntry(logrus.StandardLogger())
	randomDataAPITransformerHandler := http.NewRandomDataAPITransformer(logger)

	return randomDataAPITransformerDependencies{
		randomDataAPITransformerHandler: randomDataAPITransformerHandler,
	}
}
