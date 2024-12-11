package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	// テスト用の期待値
	wantPort := ":3333"
	wantEnv := "test"
	wantExternalAPIURL := "https://testapi.example.com/api/json?method=searchByPostal&postal="

	// 環境変数を設定
	t.Setenv("PORT", wantPort)
	t.Setenv("TODO_ENV", wantEnv)
	t.Setenv("EXTERNAL_API", wantExternalAPIURL)

	// Config を作成
	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}

	// ポートの確認
	if got.Port != wantPort {
		t.Errorf("Port mismatch: want %s, got %s", wantPort, got.Port)
	}

	// 環境名の確認
	if got.Env != wantEnv {
		t.Errorf("Environment mismatch: want %s, got %s", wantEnv, got.Env)
	}

	// 外部API URLの確認
	if got.ExternalAPI != wantExternalAPIURL {
		t.Errorf("External API URL mismatch: want %s, got %s", wantExternalAPIURL, got.ExternalAPI)
	}
}
