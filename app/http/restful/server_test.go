package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/si3nloong/webhook/app/entity"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	webhooks []*entity.WebhookRequest
}

func (m *mockRepository) CreateWebhook(ctx context.Context, data *entity.WebhookRequest) error {
	m.webhooks = append(m.webhooks, data)
	return nil
}

func (m *mockRepository) FindWebhook(ctx context.Context, id string) (*entity.WebhookRequest, error) {
	for _, wh := range m.webhooks {
		if wh.ID.String() == id {
			return wh, nil
		}
	}
	return nil, errors.New("data not found")
}

func (m *mockRepository) UpdateWebhook(ctx context.Context, id string, attempt *entity.Attempt) error {
	for _, wh := range m.webhooks {
		if wh.ID.String() == id {
			wh.Attempts = append(wh.Attempts, *attempt)
			break
		}
	}
	return nil
}

type mockWebhookServer struct {
}

func TestServer(t *testing.T) {
	svr := new(Server)
	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// fmt.Fprintf(w, expected)
	// }))
	// defer ts.Close()

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(svr.health)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
