package main

import (
	"context"
	"sync"

	"github.com/awcjack/ETL-sample/config"
	"github.com/awcjack/ETL-sample/extraction"
	"github.com/awcjack/ETL-sample/loading"
	"github.com/awcjack/ETL-sample/transformation"
	"github.com/awcjack/ETL-sample/transformation/http"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func main() {
	// Init logrus as logger
	logger := logrus.New()

	// Loading config from env
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Not able to load config ", err)
	}

	// set logger level from config
	logLevel, err := logrus.ParseLevel(config.Application.LogLevel)
	if err != nil {
		logger.Errorf("Not able to parse log level %v, changing to default level info level", err)
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Create repository implementation (possible to switch to other datastore implementation)
	var db *sqlx.DB
	logger.Info("Loading datastore Reopsitory")
	if config.Database.Type == "postgresql" {
		// Start PostgreSQL connection
		db, err = loading.NewPostgreSQLConnection(config.Database)
		if err != nil {
			logger.Fatal("Not able to create new postgresql connection", err)
		}
	} else {
		logger.Fatal("Not implemented other datastore repository")
	}
	// Create repostiroy based on PostgreSQL db connection
	repo := loading.NewPostgreSQLRepository(db, logger)

	// create data channel for passing data from transformer to data store
	structedDataChan := make(chan transformation.TransformedData, config.Application.ProcessPipelineSize)
	// dedicate go routine for storing processed data to datastore
	go loading.SaveData(context.Background(), repo, structedDataChan, config.Application.BulkInsert, config.Application.BulkInsertSize, config.Application.BulkInsertInterval)

	// data extraction processor
	httpExtractionProcessor := extraction.NewHttpExtraction(logger)
	// waitgroup to make sure the application won't close before all extraction processor fail
	var wg sync.WaitGroup
	// for loop to create extract go routine based on config
	for _, datasource := range config.Datasource {
		wg.Add(1)
		if datasource.Type == "http" && datasource.Transformer == "random-data-api" {
			logger.Debugf("datasource %s is starting", datasource.Name)
			// create transformer based on type
			randomDataAPITransformer := http.NewRandomDataAPITransformer(logger)
			// dedicate go routine for starting extract data from data source which allow getting data from different data source simultaneously
			go httpExtractionProcessor.Extract(datasource.Source, randomDataAPITransformer.Transform, structedDataChan, &wg)
		} else {
			// close wait group if no extract function started or failed
			wg.Done()
			logger.Errorf("datasource %s type : %s, transformer: %s not yet implemented\n", datasource.Name, datasource.Type, datasource.Transformer)
		}
	}

	wg.Wait()
}
