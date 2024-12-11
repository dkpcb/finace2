package handler

import (
	"net/http"

	"github.com/dkpcb/finatext_kadai_2/config"
	"github.com/dkpcb/finatext_kadai_2/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	AddressService   *service.AddressService
	AccessLogService *service.AccessLogService
	Cfg              *config.Config
}

// 新しい Handler を作成
func NewHandler(addressService *service.AddressService, accessLogService *service.AccessLogService, cfg *config.Config) *Handler {
	return &Handler{
		AddressService:   addressService,
		AccessLogService: accessLogService,
		Cfg:              cfg,
	}
}

// ルーティングを登録
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.HandleRoot)
	e.GET("/address", h.HandleAddress)
	e.GET("/address/access_logs", h.HandleAccessLogs)
}

// ルートエンドポイントを処理
func (h *Handler) HandleRoot(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

// HandleAddress は住所検索のエンドポイントを処理
func (h *Handler) HandleAddress(c echo.Context) error {
	// クエリパラメータから郵便番号を取得
	postalCode := c.QueryParam("postal_code")
	if postalCode == "" {
		// 郵便番号がない場合は400エラーを返す
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "postal_code is required"})
	}

	// アクセスログを保存
	if err := h.AccessLogService.SaveAccessLog(postalCode); err != nil {
		// ログ保存でエラーが発生した場合は500エラーを返す
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save access log"})
	}

	// サービス層で住所データを取得
	address, err := h.AddressService.GetAddress(postalCode)
	if err != nil {
		// サービス層でエラーが発生した場合は500エラーを返す
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if address == nil {
		// 該当する住所がない場合は404エラーを返す
		return c.JSON(http.StatusNotFound, map[string]string{"error": "address not found"})
	}

	// 正常時は200 OKと住所データを返す
	return c.JSON(http.StatusOK, address)
}

// アクセスログを集計して返す
func (h *Handler) HandleAccessLogs(c echo.Context) error {
	// アクセスログの集計結果を取得
	logs, err := h.AccessLogService.GetAccessLogs()
	if err != nil {
		// エラーが発生した場合は500エラーを返す
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 正常時は200 OKとアクセスログを返す
	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_logs": logs,
	})
}
