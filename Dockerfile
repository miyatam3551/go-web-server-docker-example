# ===== ビルドステージ =====
# Goのコンパイル環境を持つイメージを使用し、"builder"という名前をつける
FROM golang:1.24 AS builder

LABEL version="1.0.0" \
      maintainer="miyatam3551" \
      description="simple go webserver"

WORKDIR /app

# 依存関係ファイルを先にコピー（キャッシュ効率化）
COPY go.mod go.sum* ./

# 依存関係のダウンロード（go.mod/go.sumが変更されない限りキャッシュ利用）
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（静的リンクでビルド）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# ===== 実行ステージ =====
# distroless: シェルや不要なツールを含まない軽量・セキュアなイメージ
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

# 非rootユーザーに切り替える（nonrootは distroless に組み込み済み、UID:65532）
USER nonroot

# builderステージからビルド済みバイナリだけをコピー
COPY --from=builder /app/server .

# メインコマンドを固定
ENTRYPOINT ["./server"]

# デフォルトのポート番号（docker run時に上書き可能）
CMD ["--port=8080"]
