package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dkpcb/finatext_kadai_2/config"
	"github.com/dkpcb/finatext_kadai_2/handler"
	"github.com/dkpcb/finatext_kadai_2/infra"
	"github.com/dkpcb/finatext_kadai_2/service"
	"github.com/labstack/echo/v4"
)

// Run はアプリケーションのエントリーポイント
func Run(ctx context.Context) error {
	// 設定の初期化
	cfg, err := initializeConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// データベースと依存性の初期化
	dbManager, services, err := initializeDependencies(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}
	defer dbManager.DB.Close()

	// サーバーの起動とシャットダウン管理
	return startServer(ctx, cfg, services)
}

// 設定を初期化
func initializeConfig() (*config.Config, error) {
	return config.New()
}

// データベースと依存性を初期化
func initializeDependencies(cfg *config.Config) (*infra.DBManager, *service.ServiceRegistry, error) {
	// データベース接続を初期化
	dbManager, err := infra.NewDBManager(cfg.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize DB manager: %w", err)
	}

	// スキーマを初期化
	if err := dbManager.InitializeSchema("finatext_db"); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// リポジトリを初期化
	addressRepo := infra.NewAddressRepository(cfg.ExternalAPI)
	accessLogRepo := infra.NewAccessLogRepository(dbManager.DB)

	// サービスを初期化
	services := &service.ServiceRegistry{
		Address:   service.NewAddressService(addressRepo, cfg.ExternalAPI),
		AccessLog: service.NewAccessLogService(accessLogRepo),
	}

	return dbManager, services, nil
}

// サーバーを起動し、シグナルを監視してグレースフルシャットダウンを実行
func startServer(ctx context.Context, cfg *config.Config, services *service.ServiceRegistry) error {
	e := echo.New()

	// ルートを登録
	h := handler.NewHandler(services.Address, services.AccessLog, cfg)
	h.RegisterRoutes(e)

	// シグナルの監視
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// サーバーを非同期で起動
	go func() {
		if err := e.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// シグナルを受信するまで待機
	<-ctx.Done()

	// グレースフルシャットダウン
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shut down server: %w", err)
	}

	e.Logger.Info("Server gracefully stopped")
	return nil
}

// package server

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"github.com/dkpcb/finatext_kadai_2/config"
// 	"github.com/dkpcb/finatext_kadai_2/handler"
// 	"github.com/dkpcb/finatext_kadai_2/infra"
// 	"github.com/dkpcb/finatext_kadai_2/service"
// 	"github.com/labstack/echo/v4"
// )

// // ConfigLoader は設定をロードする関数
// type ConfigLoader func() (*config.Config, error)

// // Run はアプリケーションのエントリーポイント
// func Run(ctx context.Context) error {
// 	return RunWithMockConfig(ctx, config.New)
// }

// // RunWithMockConfig は設定ローダーを差し替え可能なエントリーポイント
// func RunWithMockConfig(ctx context.Context, loadConfig ConfigLoader) error {
// 	// 設定の初期化
// 	cfg, err := loadConfig()
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize config: %w", err)
// 	}

// 	// データベースと依存性の初期化
// 	dbManager, services, err := initializeDependencies(cfg)
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize dependencies: %w", err)
// 	}
// 	defer dbManager.DB.Close()

// 	// サーバーの起動とシャットダウン管理
// 	return startServer(ctx, cfg, services)
// }

// // データベースと依存性を初期化
// func initializeDependencies(cfg *config.Config) (*infra.DBManager, *service.ServiceRegistry, error) {
// 	dbManager, err := infra.NewDBManager(cfg.DSN)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	addressRepo := infra.NewAddressRepository(cfg.ExternalAPI)
// 	accessLogRepo := infra.NewAccessLogRepository(dbManager.DB)

// 	services := &service.ServiceRegistry{
// 		Address:   service.NewAddressService(addressRepo, cfg.ExternalAPI),
// 		AccessLog: service.NewAccessLogService(accessLogRepo),
// 	}

// 	return dbManager, services, nil
// }

// // サーバーを起動し、シグナルを監視してグレースフルシャットダウンを実行
// func startServer(ctx context.Context, cfg *config.Config, services *service.ServiceRegistry) error {
// 	e := echo.New()

// 	// ルートを登録
// 	h := handler.NewHandler(services.Address, services.AccessLog, cfg)
// 	h.RegisterRoutes(e)

// 	// シグナルの監視
// 	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
// 	defer stop()

// 	// サーバーを非同期で起動
// 	go func() {
// 		if err := e.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
// 			e.Logger.Fatal("shutting down the server")
// 		}
// 	}()

// 	// シグナルを受信するまで待機
// 	<-ctx.Done()

// 	// グレースフルシャットダウン
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := e.Shutdown(shutdownCtx); err != nil {
// 		return fmt.Errorf("failed to shut down server: %w", err)
// 	}

// 	e.Logger.Info("Server gracefully stopped")
// 	return nil
// }
