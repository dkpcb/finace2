package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dkpcb/finatext_kadai_2/config"
	"github.com/dkpcb/finatext_kadai_2/server"
	"github.com/stretchr/testify/assert"
)

// MockConfigLoader は config.New をモックする関数
type MockConfigLoader func() (*config.Config, error)

// TestRun は Run 関数のテスト
func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// サーバー起動を非同期で実行
		go func() {
			time.Sleep(2 * time.Second) // サーバーが起動するのを待機
			cancel()                    // シャットダウンシグナル送信
		}()

		err := server.RunWithMockConfig(ctx, mockConfigSuccess)
		assert.NoError(t, err)
	})

	t.Run("failure_initializeConfig", func(t *testing.T) {
		ctx := context.Background()

		// Mock Config Initialization Error
		err := server.RunWithMockConfig(ctx, mockConfigFailure)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock config error")
	})
}

// mockConfigSuccess は成功する MockConfigLoader
func mockConfigSuccess() (*config.Config, error) {
	return &config.Config{
		Port: ":8080",
		DSN:  "root:password@tcp(localhost:3306)/",
	}, nil
}

// mockConfigFailure は失敗する MockConfigLoader
func mockConfigFailure() (*config.Config, error) {
	return nil, errors.New("mock config error")
}
