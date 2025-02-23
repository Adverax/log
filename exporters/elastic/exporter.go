package elasticExporter

import (
	"context"
	"github.com/adverax/log"
	"github.com/olivere/elastic/v7"
)

type Exporter struct {
	client    *elastic.Client
	formatter log.Formatter
	index     string
}

func New(
	client *elastic.Client,
	formatter log.Formatter,
	index string,
) *Exporter {
	return &Exporter{
		client:    client,
		formatter: formatter,
		index:     index,
	}
}

func (that *Exporter) Export(ctx context.Context, entry *log.Entry) {
	data, err := that.formatter.Format(entry)
	if err != nil {
		return
	}

	_, _ = that.client.Index().
		Index(that.index).
		BodyJson(data).
		Do(ctx)
	// nothing
}
