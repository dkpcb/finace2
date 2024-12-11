#  API ドキュメント

## 機能概要

1. **ルートエンドポイント**  
   エンドポイント: `GET http://localhost:8080/`  
   メッセージを返す。

2. **データの取得**  
   エンドポイント: `GET http://localhost:8080/address?postal_code=[7桁の郵便番号]`  
   レスポンス例:
   ```json
   {
       "postal_code": "5016121",
       "hit_count": 6,
       "address": "岐阜県岐阜市柳津町",
       "tokyo_sta_distance": 278.3
   }
   ```

3. **アクセスログの取得**  
   エンドポイント: `GET http://localhost:8080/address/access_logs`  
   各郵便番号のリクエスト回数を集計して返す。  
   レスポンス例:
   ```json
   {
       "access_logs": [
           {
               "postal_code": "1020073",
               "request_count": 7
           },
           {
               "postal_code": "1000001",
               "request_count": 5
           }
       ]
   }
   ```

4. **データベースへのログ記録**  
   `/address` へのリクエストを MySQL データベースに記録する。

---

## 前提条件

- [Go 1.22+](https://golang.org/doc/install) がインストールされていること
- [Docker](https://www.docker.com/) がインストールされていること
- [Docker Compose](https://docs.docker.com/compose/) がインストールされていること

---

## セットアップと実行手順

### ステップ 1: Docker コンテナのビルドと実行

1. Docker イメージをビルドする:
   ```bash
   make build
   ```

2. サービスを起動する:
   ```bash
   make up
   ```

3. サービスを停止してクリーンアップする:
   ```bash
   make down
   ```

#### リクエスト例:

1. 住所データの取得:
   ```bash
   curl "http://localhost:8080/address?postal_code=5016121"
   ```

2. アクセスログの取得:
   ```bash
   curl "http://localhost:8080/address/access_logs"
   ```

---
