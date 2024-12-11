# ビルド用のステージ
FROM golang:1.21-alpine as server-build

WORKDIR /app

# 必要なツールのインストール
RUN apk upgrade --update && \
    apk --no-cache add git && \
    apk --no-cache add ca-certificates

# go.mod と go.sum をコピーして依存関係を解決
COPY go.mod go.sum ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main .

# 実行用の軽量イメージ
FROM alpine:3.18

WORKDIR /app

# 必要なツールのインストール
RUN apk --no-cache add ca-certificates

# ビルド済みのファイルをコピー
COPY --from=server-build /app/main .

# アプリケーションを実行
CMD ["./main"]
