package extraction

import (
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/utils"
)

type HttpExtraction struct {
	logger utils.Logger
}

func NewHttpExtraction(logger utils.Logger) *HttpExtraction {
	return &HttpExtraction{
		logger: logger,
	}
}

// extract function to run GET request to url
// transform data using transformer function in params
// pass data to data channel
func (h *HttpExtraction) Extract(url string, transformer func(data []byte) (transformation.TransformedData, error), dataPipeline chan<- transformation.TransformedData, wg *sync.WaitGroup) error {
	h.logger.Debugf("HTTP source: ", url)

	defer wg.Done()

	// random seed based on current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// fetch data from data source (url)
		resp, err := http.Get(url)

		if err != nil {
			return err
		}

		// read response body
		body, err := io.ReadAll((resp.Body))

		if err != nil {
			return err
		}

		// close response body to release memory
		resp.Body.Close()

		// transform data based on transformer function from params
		transformedData, err := transformer(body)
		if err != nil {
			return err
		}

		h.logger.Debugf("inserted data to channel", transformedData)
		// push transformed data to channel for storing data to storage
		dataPipeline <- transformedData

		// random delay from 250ms to 750ms
		randomDelay := r.Float32()*500 + 250
		time.Sleep(time.Duration(randomDelay) * time.Millisecond)
	}
}
