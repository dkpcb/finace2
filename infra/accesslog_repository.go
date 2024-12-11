package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dkpcb/finatext_kadai_2/entity"
)

type AccessLogRepository struct {
	DB *sql.DB
}

// 新しい AccessLogRepository を作成
func NewAccessLogRepository(db *sql.DB) *AccessLogRepository {
	return &AccessLogRepository{DB: db}
}

// // 郵便番号と作成日時を含むアクセスログを保存
// func (r *AccessLogRepository) InsertAccessLog(postalCode string, createdAt time.Time) error {
// 	query := "INSERT INTO access_logs (postal_code, created_at) VALUES (?, ?)"
// 	_, err := r.DB.Exec(query, postalCode, createdAt)
// 	return err
// }

func (r *AccessLogRepository) InsertAccessLog(postalCode string, createdAt time.Time) error {
	query := "INSERT INTO access_logs (postal_code, created_at) VALUES (?, ?)"
	_, err := r.DB.Exec(query, postalCode, createdAt)
	if err != nil {
		fmt.Printf("Failed to execute query: %s with error: %v\n", query, err)
	}
	return err
}

// リクエスト回数を集計し、降順で返す
func (r *AccessLogRepository) GetAccessLogs() ([]entity.AccessLog, error) {
	// SQLクエリの更新：リクエスト回数を集計し降順で返す
	query := `
		SELECT postal_code, COUNT(*) as request_count
		FROM access_logs
		GROUP BY postal_code
		ORDER BY request_count DESC
	`
	// クエリを実行
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch access logs: %w", err)
	}
	defer rows.Close()

	// 結果を AccessLog 構造体のスライスに詰める
	var logs []entity.AccessLog
	for rows.Next() {
		var log entity.AccessLog
		// スキャン時に ID と CreatedAt を削除
		if err := rows.Scan(&log.PostalCode, &log.RequestCount); err != nil {
			return nil, fmt.Errorf("failed to scan access log: %w", err)
		}
		logs = append(logs, log)
	}

	// 行の処理中にエラーが発生していないかを確認
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over access logs: %w", err)
	}

	// 最終的なスライスを返す
	return logs, nil
}
