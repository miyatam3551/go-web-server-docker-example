# Simple Go Web Server

Go言語で作成したシンプルなWebサーバです。Dockerのベストプラクティスに準拠した構成になっています。

## 📋 目次

- [機能](#機能)
- [前提条件](#前提条件)
- [クイックスタート](#クイックスタート)
  - [ローカル実行](#ローカル実行)
  - [Docker実行](#docker実行)
- [APIエンドポイント](#apiエンドポイント)
- [Dockerのベストプラクティス](#dockerのベストプラクティス)
- [開発者向け情報](#開発者向け情報)

## 🚀 機能

- **シンプルなREST API**: JSON形式でレスポンスを返すAPIエンドポイント
- **ヘルスチェック**: サーバの状態を確認できるエンドポイント
- **ログ機能**: リクエストのログを出力
- **柔軟なポート設定**: コマンドライン引数や環境変数でポート番号を変更可能

## 📦 前提条件

### ローカル実行の場合
- Go 1.24以上

### Docker実行の場合
- Docker

## 🏃 クイックスタート

### ローカル実行

1. **プロジェクトのクローン**
```bash
git clone <repository-url>
cd create-cluster-at-aws
```

2. **依存関係のダウンロード**
```bash
go mod download
```

3. **サーバの起動**
```bash
# デフォルトポート（8080）で起動
go run main.go

# カスタムポートで起動
go run main.go --port=3000
```

4. **動作確認**
```bash
curl http://localhost:8080/health
```

### Docker実行

1. **イメージのビルド**
```bash
docker build -t simple-go-webserver .
```

2. **コンテナの起動**
```bash
# デフォルトポート（8080）で起動
docker run -p 8080:8080 simple-go-webserver

# カスタムポート（3000）で起動
docker run -p 3000:3000 simple-go-webserver --port=3000
```

3. **動作確認**
```bash
curl http://localhost:8080/health
```

## 🔌 APIエンドポイント

### 1. ルート `/`
サーバのウェルカムメッセージを返します。

**リクエスト例:**
```bash
curl http://localhost:8080/
```

**レスポンス例:**
```json
{
  "message": "Welcome to Simple Go Web Server!"
}
```

### 2. ヘルスチェック `/health`
サーバの状態を確認します。

**リクエスト例:**
```bash
curl http://localhost:8080/health
```

**レスポンス例:**
```json
{
  "status": "ok",
  "message": "Server is running",
  "time": "2025-12-02T10:00:00Z"
}
```

### 3. 挨拶 `/hello`
名前を指定して挨拶メッセージを返します。

**リクエスト例:**
```bash
# 名前を指定
curl http://localhost:8080/hello?name=Alice

# 名前を指定しない場合
curl http://localhost:8080/hello
```

**レスポンス例:**
```json
{
  "message": "Hello, Alice!",
  "name": "Alice"
}
```

## 🐳 Dockerのベストプラクティス

このプロジェクトでは、以下のDockerベストプラクティスを実装しています：

### 1. ✅ マルチステージビルド
- **ビルド環境**と**実行環境**を分離
- イメージサイズを大幅に削減（1GB → 数十MB）
- セキュリティリスクを低減

```dockerfile
FROM golang:1.24 AS builder
# ... ビルド処理 ...

FROM gcr.io/distroless/base-debian12
# ... 実行環境 ...
```

### 2. ✅ キャッシュの最大活用
- `go.mod`と`go.sum`を先にコピー
- 依存関係のダウンロードをキャッシュ化
- ビルド時間を大幅短縮

```dockerfile
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
```

### 3. ✅ 軽量なベースイメージ（distroless）
- シェルや不要なツールを含まない
- 攻撃対象領域を最小化
- イメージサイズが小さい

### 4. ✅ .dockerignoreの活用
- 不要なファイルをビルドコンテキストから除外
- ビルド速度を向上
- セキュリティリスクを低減

### 5. ✅ ENTRYPOINTとCMDの適切な使い分け
- `ENTRYPOINT`: メインコマンドを固定
- `CMD`: デフォルト引数（実行時に上書き可能）

```dockerfile
ENTRYPOINT ["./server"]
CMD ["--port=8080"]
```

### 6. ✅ 非rootユーザーでの実行
- セキュリティベストプラクティス
- 権限昇格のリスクを低減

```dockerfile
USER nonroot
```

### 7. 📊 イメージサイズの比較

| 構成 | イメージサイズ |
|------|---------------|
| golang:1.24のみ | ~900MB |
| マルチステージ + distroless | ~30MB |

**削減率: 約97%！**

## 👨‍💻 開発者向け情報

### プロジェクト構成

```
.
├── main.go              # メインアプリケーション
├── go.mod               # Go モジュール定義
├── Dockerfile           # Dockerイメージ定義
├── .dockerignore        # Docker除外ファイル設定
└── README.md            # このファイル
```

### ビルドとテスト

```bash
# ローカルビルド
go build -o server

# 実行
./server --port=8080

# Dockerイメージのビルド
docker build -t simple-go-webserver .

# Dockerコンテナの実行
docker run -p 8080:8080 simple-go-webserver
```

### セキュリティスキャン（推奨）

```bash
# Trivyでの脆弱性チェック（Docker経由で実行）
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy:latest image --severity CRITICAL,HIGH simple-go-webserver

# Trivyをローカルにインストールする場合（macOS）
brew install trivy
trivy image --severity CRITICAL,HIGH simple-go-webserver
```

### Linter（推奨）

```bash
# Hadolintのインストール（macOS）
brew install hadolint

# Dockerfileの静的解析
hadolint Dockerfile
```

## 📚 参考資料

- [Go公式ドキュメント](https://go.dev/doc/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)

## 📝 ライセンス

MIT License

## 🤝 コントリビューション

プルリクエストを歓迎します！
