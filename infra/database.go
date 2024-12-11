package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBManager はデータベース操作の責務を持つ構造体
type DBManager struct {
	DB *sql.DB
}

// DBManager を初期化
func NewDBManager(dsn string) (*DBManager, error) {
	db, err := initializeConnection(dsn)
	if err != nil {
		return nil, err
	}

	return &DBManager{DB: db}, nil
}

// iデータベース接続を初期化
func initializeConnection(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 2; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("Database connection successful!")
				return db, nil
			}
		}
		// リトライ中はログを一度だけ表示
		if i == 0 {
			fmt.Println("Retrying database connection...")
		}
		time.Sleep(5 * time.Second)
	}

	// すべてのリトライが失敗した場合のみエラーを返す
	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
}

// データベースとテーブルを作成
func (m *DBManager) InitializeSchema(databaseName string) error {
	// データベース作成
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", databaseName)
	if _, err := m.DB.Exec(createDBQuery); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// 使用するデータベースを指定
	useDBQuery := fmt.Sprintf("USE %s;", databaseName)
	if _, err := m.DB.Exec(useDBQuery); err != nil {
		return fmt.Errorf("failed to use database: %w", err)
	}

	// テーブル作成
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS access_logs (
		id INT AUTO_INCREMENT NOT NULL,
		postal_code VARCHAR(8) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`
	if _, err := m.DB.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Database and table initialized successfully.")
	return nil
}
