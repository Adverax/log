package elasticExporter

import (
	"context"
	"github.com/adverax/log"
	jsonFormatter "github.com/adverax/log/formatters/json"
	"github.com/olivere/elastic/v7"
	"os"
)

func ExampleExporter() {
	formatter, err := jsonFormatter.NewBuilder().Build()
	if err != nil {
		panic(err)
	}

	client, err := elastic.NewClient(
		elastic.SetURL("http://"+os.Getenv("ELASTICSEARCH_HOST")+":9200"),
		elastic.SetSniff(false), // Отключение sniffing, если необходимо
	)
	if err != nil {
		panic(err)
	}

	logger, err := log.NewBuilder().
		WithLevel(log.InfoLevel).
		WithExporter(
			New(client, formatter, "log"),
		).
		Build()
	if err != nil {
		panic(err)
	}

	logger.Info(context.Background(), "Hello, World!")
}
