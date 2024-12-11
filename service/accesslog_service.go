package service

import (
	"fmt"
	"time"

	"github.com/dkpcb/finatext_kadai_2/entity"
	"github.com/dkpcb/finatext_kadai_2/infra"
)

type AccessLogService struct {
	LogRepo *infra.AccessLogRepository // ポインタ型
}

// 新しい AccessLogService を作成
func NewAccessLogService(logRepo *infra.AccessLogRepository) *AccessLogService {
	return &AccessLogService{LogRepo: logRepo}
}

// 郵便番号を含むアクセスログを保存
func (s *AccessLogService) SaveAccessLog(postalCode string) error {
	// 現在時刻を取得してログを保存
	now := time.Now()
	return s.LogRepo.InsertAccessLog(postalCode, now)
}

// リクエスト回数を集計し、降順で返す
func (s *AccessLogService) GetAccessLogs() ([]entity.AccessLog, error) {
	// アクセスログを集計
	logs, err := s.LogRepo.GetAccessLogs()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch access logs: %w", err)
	}

	// ログを返す
	return logs, nil
}
