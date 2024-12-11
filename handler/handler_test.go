package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dkpcb/finatext_kadai_2/config"
	"github.com/dkpcb/finatext_kadai_2/entity"
	"github.com/dkpcb/finatext_kadai_2/handler"
	"github.com/dkpcb/finatext_kadai_2/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockAddressRepository は entity.AddressRepository を模倣
type MockAddressRepository struct{}

func (m *MockAddressRepository) FetchAddressData(postalCode string) ([]entity.AddressLocation, error) {
	if postalCode == "5016121" {
		return []entity.AddressLocation{
			{Prefecture: "岐阜県", City: "岐阜市", Town: "柳津町", Lat: 35.355743, Lon: 136.725408},
		}, nil
	}
	if postalCode == "9999999" {
		return nil, nil
	}
	return nil, errors.New("mock error")
}

// MockAccessLogRepository は infra.AccessLogRepository を模倣
type MockAccessLogRepository struct{}

func (m *MockAccessLogRepository) InsertAccessLog(postalCode string, createdAt time.Time) error {
	if postalCode == "error" {
		return errors.New("mock save error")
	}
	return nil
}

func (m *MockAccessLogRepository) GetAccessLogs() ([]entity.AccessLog, error) {
	return []entity.AccessLog{
		{PostalCode: "1020073", RequestCount: 7},
		{PostalCode: "1000001", RequestCount: 5},
		{PostalCode: "5300001", RequestCount: 2},
	}, nil
}

func TestHandler_HandleAddress(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{Port: ":8080", ExternalAPI: "https://example.com"}

	// モックをサービスにラップ
	addressService := service.NewAddressService(&MockAddressRepository{}, cfg.ExternalAPI)
	accessLogService := service.NewAccessLogService(&MockAccessLogRepository{})

	h := handler.NewHandler(addressService, accessLogService, cfg)

	tests := []struct {
		name           string
		postalCode     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "正常な郵便番号",
			postalCode:     "5016121",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"postal_code":"5016121","hit_count":1,"common_address":"岐阜県岐阜市柳津町","tokyo_sta_distance":277.7}`,
		},
		{
			name:           "郵便番号が空",
			postalCode:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"postal_code is required"}`,
		},
		{
			name:           "存在しない郵便番号",
			postalCode:     "9999999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"address not found"}`,
		},
		{
			name:           "ログ保存エラー",
			postalCode:     "error",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to save access log"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/address?postal_code="+tt.postalCode, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.HandleAddress(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestHandler_HandleAccessLogs(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{Port: ":8080"}

	// モックをサービスにラップ
	addressService := service.NewAddressService(&MockAddressRepository{}, cfg.ExternalAPI)
	accessLogService := service.NewAccessLogService(&MockAccessLogRepository{})

	h := handler.NewHandler(addressService, accessLogService, cfg)

	req := httptest.NewRequest(http.MethodGet, "/address/access_logs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleAccessLogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody := `{
		"access_logs": [
			{"postal_code": "1020073", "request_count": 7},
			{"postal_code": "1000001", "request_count": 5},
			{"postal_code": "5300001", "request_count": 2}
		]
	}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}
