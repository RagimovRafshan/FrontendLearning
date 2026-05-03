package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"messenger/config"
)

// ElasticSearchClient обертка над elasticsearch.Client для работы с ElasticSearch
type ElasticSearchClient struct {
	client      *elasticsearch.Client
	indexPrefix string
}

// NewElasticSearchClient создает новый клиент ElasticSearch
func NewElasticSearchClient(cfg config.ElasticSearchConfig) (*ElasticSearchClient, error) {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Sniff:     cfg.Sniff,
		Transport: &httpTransportWithTimeout{timeout: cfg.Timeout},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ElasticSearch client: %w", err)
	}

	// Проверка подключения при старте
	if cfg.HealthCheck {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		
		res, err := esClient.Info(esapi.Info{}.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to ElasticSearch: %w", err)
		}
		defer res.Body.Close()
		
		if res.IsError() {
			return nil, fmt.Errorf("ElasticSearch returned error: %s", res.String())
		}
	}

	client := &ElasticSearchClient{
		client:      esClient,
		indexPrefix: cfg.IndexPrefix,
	}

	return client, nil
}

// httpTransportWithTimeout кастомный транспорт с таймаутом
type httpTransportWithTimeout struct {
	timeout time.Duration
}

func (t *httpTransportWithTimeout) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(req.Context(), t.timeout)
	defer cancel()
	
	req = req.WithContext(ctx)
	return http.DefaultTransport.RoundTrip(req)
}

// Client возвращает базовый клиент elasticsearch
func (c *ElasticSearchClient) Client() *elasticsearch.Client {
	return c.client
}

// IndexPrefix возвращает префикс индексов
func (c *ElasticSearchClient) IndexPrefix() string {
	return c.indexPrefix
}

// HealthCheck проверяет подключение к ElasticSearch
func (c *ElasticSearchClient) HealthCheck(ctx context.Context) error {
	res, err := c.client.Info(c.client.Info.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("ElasticSearch health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ElasticSearch returned error status: %s", res.String())
	}

	return nil
}

// buildIndexName строит полное имя индекса с префиксом
func (c *ElasticSearchClient) buildIndexName(indexName string) string {
	return fmt.Sprintf("%s_%s", c.indexPrefix, indexName)
}

// CreateIndex создает индекс с настройками по умолчанию
func (c *ElasticSearchClient) CreateIndex(ctx context.Context, indexName string, mappings map[string]interface{}) error {
	fullIndexName := c.buildIndexName(indexName)

	body := make(map[string]interface{})
	if mappings != nil {
		body["mappings"] = mappings
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal index settings: %w", err)
	}

	res, err := c.client.Indices.Create(
		fullIndexName,
		c.client.Indices.Create.WithBody(bytes.NewReader(bodyBytes)),
		c.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}

// IndexDocument индексирует документ
func (c *ElasticSearchClient) IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	fullIndexName := c.buildIndexName(indexName)

	bodyBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	var res *esapi.Response
	var err2 error
	
	if docID != "" {
		res, err2 = c.client.Index(
			fullIndexName,
			bytes.NewReader(bodyBytes),
			c.client.Index.WithDocumentID(docID),
			c.client.Index.WithContext(ctx),
		)
	} else {
		res, err2 = c.client.Index(
			fullIndexName,
			bytes.NewReader(bodyBytes),
			c.client.Index.WithContext(ctx),
		)
	}

	if err2 != nil {
		return fmt.Errorf("failed to index document: %w", err2)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document: %s", res.String())
	}

	return nil
}

// GetDocument получает документ по ID
func (c *ElasticSearchClient) GetDocument(ctx context.Context, indexName, docID string) (map[string]interface{}, error) {
	fullIndexName := c.buildIndexName(indexName)

	res, err := c.client.Get(
		fullIndexName,
		docID,
		c.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, ErrDocumentNotFound
		}
		return nil, fmt.Errorf("failed to get document: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid document format")
	}

	return source, nil
}

// DeleteDocument удаляет документ по ID
func (c *ElasticSearchClient) DeleteDocument(ctx context.Context, indexName, docID string) error {
	fullIndexName := c.buildIndexName(indexName)

	res, err := c.client.Delete(
		fullIndexName,
		docID,
		c.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return ErrDocumentNotFound
		}
		return fmt.Errorf("failed to delete document: %s", res.String())
	}

	return nil
}

// Search выполняет поиск по индексу
func (c *ElasticSearchClient) Search(ctx context.Context, indexName string, query map[string]interface{}, size int) ([]map[string]interface{}, error) {
	fullIndexName := c.buildIndexName(indexName)

	var bodyBytes []byte
	var err error

	if query != nil {
		bodyBytes, err = json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal search query: %w", err)
		}
	} else {
		// Поиск по всем документам
		bodyBytes = []byte(`{"query": {"match_all": {}}}`)
	}

	res, err := c.client.Search(
		c.client.Search.WithIndex(fullIndexName),
		c.client.Search.WithBody(bytes.NewReader(bodyBytes)),
		c.client.Search.WithSize(size),
		c.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}

	hitArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}

	documents := make([]map[string]interface{}, len(hitArray))
	for i, hit := range hitArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		documents[i] = source
	}

	return documents, nil
}

// BulkDocuments выполняет массовую индексацию документов
func (c *ElasticSearchClient) BulkDocuments(ctx context.Context, indexName string, docs []BulkDoc) error {
	fullIndexName := c.buildIndexName(indexName)

	var body bytes.Buffer
	for _, doc := range docs {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": fullIndexName,
			},
		}
		if doc.ID != "" {
			meta["index"].(_map[string]interface{})["_id"] = doc.ID
		}

		metaJSON, _ := json.Marshal(meta)
		body.Write(metaJSON)
		body.WriteByte('\n')

		docJSON, _ := json.Marshal(doc.Document)
		body.Write(docJSON)
		body.WriteByte('\n')
	}

	res, err := c.client.Bulk(
		bytes.NewReader(body.Bytes()),
		c.client.Bulk.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index failed: %s", res.String())
	}

	return nil
}

// BulkDoc структура для массовой индексации
type BulkDoc struct {
	ID       string
	Document interface{}
}

// Ошибки
var (
	ErrDocumentNotFound = fmt.Errorf("document not found")
)
