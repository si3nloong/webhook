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
	"github.com/si3nloong/webhook/cmd"
	"github.com/si3nloong/webhook/model"
	"github.com/tidwall/gjson"
)

type db struct {
	indexName string
	es        *elasticsearch.Client
}

func New(cfg cmd.Config) (*db, error) {
	esConfig := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		// Username: "user",
		// Password: "pass",
	}

	// es, err := elasticsearch.NewDefaultClient()
	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	log.Println(string(b))

	v := new(db)
	v.indexName = "webhook_index"
	v.es = es
	return v, nil
}

func (c *db) GetLogs(ctx context.Context) (datas []model.Log, err error) {
	var buf bytes.Buffer
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.indexName),
		c.es.Search.WithBody(&buf),
		c.es.Search.WithTrackTotalHits(true),
		c.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf.Reset()

	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, err
	}

	result := gjson.GetBytes(buf.Bytes(), "hits.hits").Array()

	for _, r := range result {
		data := model.Log{}
		if err := json.Unmarshal([]byte(r.Get("_source").Raw), &data); err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	return
}

func (c *db) FindLog(ctx context.Context, id string) (data *model.Log, err error) {
	data = new(model.Log)
	return
}

func (c *db) InsertLog(ctx context.Context, data *model.Log) error {
	blr := new(bytes.Buffer)

	data.ID = ksuid.New()
	data.CreatedAt = time.Now().UTC()
	if err := json.NewEncoder(blr).Encode(data); err != nil {
		return err
	}

	// Instantiate a request object
	req := esapi.IndexRequest{
		Index: c.indexName,
		// DocumentID: strconv.Itoa(i + 1),
		Body:    blr,
		Refresh: "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
