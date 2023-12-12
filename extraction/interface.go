package extraction

import (
	"sync"

	"github.com/awcjack/ETL-sample/transformation"
)

type DataSourceExtration interface {
	Extract(source string, transformer func(data []byte) (transformation.TransformedData, error), dataPipeline chan<- transformation.TransformedData, wg *sync.WaitGroup) error
}
