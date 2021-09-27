package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/segmentio/ksuid"
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/cmd"
	"github.com/tidwall/gjson"
)

type db struct {
	indexName string
	client    *elasticsearch.Client
	timeout   time.Duration
}

func New(cfg cmd.Config) (*db, error) {
	esConfig := elasticsearch.Config{
		Addresses: []string{
			cfg.Elasticsearch.Host,
		},
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
		APIKey:   cfg.Elasticsearch.ApiKey,
	}

	// es, err := elasticsearch.NewDefaultClient()
	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}

	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	log.Println(string(b))

	v := new(db)
	v.indexName = "webhook_index"
	v.client = client
	v.timeout = 1 * time.Minute
	return v, nil
}

func (c *db) GetLogs(ctx context.Context, curCursor string, limit uint) (datas []*entity.WebhookRequest, nextCursor string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var buf bytes.Buffer
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.indexName),
		c.client.Search.WithBody(&buf),
		c.client.Search.WithTrackTotalHits(true),
		c.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	buf.Reset()

	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, "", err
	}

	log.Println(buf.String())
	result := gjson.GetBytes(buf.Bytes(), "hits.hits").Array()

	for _, r := range result {
		data := entity.WebhookRequest{}
		if err := json.Unmarshal([]byte(r.Get("_source").Raw), &data); err != nil {
			return nil, "", err
		}
		datas = append(datas, &data)
	}

	return
}

func (c *db) FindLog(ctx context.Context, id string) (data *entity.WebhookRequest, err error) {
	// Instantiate a request object
	req := esapi.GetRequest{
		Index:      c.indexName,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	log.Println(buf.String())

	data = new(entity.WebhookRequest)
	return
}

func (c *db) InsertLog(ctx context.Context, data *entity.WebhookRequest) error {
	blr := new(bytes.Buffer)

	nilTime := time.Time{}
	utcNow := time.Now().UTC()
	if data.ID == ksuid.Nil {
		data.ID = ksuid.New()
	}
	if data.CreatedAt == nilTime {
		data.CreatedAt = utcNow
	}
	if data.UpdatedAt == nilTime {
		data.UpdatedAt = utcNow
	}
	if err := json.NewEncoder(blr).Encode(data); err != nil {
		return err
	}

	// Instantiate a request object
	req := esapi.IndexRequest{
		Index:      c.indexName,
		DocumentID: data.ID.String(),
		Body:       blr,
		Refresh:    "true",
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
