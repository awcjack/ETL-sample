package extraction

import (
	"sync"

	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/utils"
)

type FileExtraction struct {
	logger utils.Logger
}

func NewFileExtraction(logger utils.Logger) *HttpExtraction {
	return &HttpExtraction{
		logger: logger,
	}
}

func (h *FileExtraction) Extract(path string, transformer func(data []byte) (transformation.TransformedData, error), dataPipeline chan<- transformation.TransformedData, wg *sync.WaitGroup) error {
	// to be implement
	h.logger.Errorf("Not yet implemented")
	wg.Done()
	return errorNotImplemented
}
