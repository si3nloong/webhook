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

func (c *db) Incr(ctx context.Context, id string) error {
	var buf bytes.Buffer
	buf.WriteString(`{
		"script" : {
			"source": "ctx._source.counter += params.count",
			"lang": "painless",
			"params" : {
				"count" : 1
			}
		}
	}`)
	res, err := c.client.Update(
		c.indexName,
		id,
		&buf,
		c.client.Update.WithContext(ctx),
		c.client.Update.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (c *db) GetWebhooks(ctx context.Context, curCursor string, limit uint) (datas []*entity.WebhookRequest, nextCursor string, err error) {
	var buf bytes.Buffer
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.indexName),
		c.client.Search.WithBody(&buf),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	buf.Reset()
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, "", err
	}

	result := gjson.GetBytes(buf.Bytes(), "hits.hits").Array()

	for _, r := range result {
		data := entity.WebhookRequest{}
		err = json.Unmarshal([]byte(r.Get("_source").Raw), &data)
		if err != nil {
			return
		}

		datas = append(datas, &data)
	}
	return
}

func (c *db) FindWebhook(ctx context.Context, id string) (data *entity.WebhookRequest, err error) {
	// Instantiate a request object
	req := esapi.GetRequest{
		Index:      c.indexName,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var o struct {
		Source interface{} `json:"_source"`
	}

	data = new(entity.WebhookRequest)
	o.Source = data
	err = json.NewDecoder(res.Body).Decode(&o)
	return
}

func (c *db) CreateWebhook(ctx context.Context, data *entity.WebhookRequest) error {
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

	if _, err := req.Do(ctx, c.client); err != nil {
		return err
	}

	return nil
}
