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
- [CI/CD](#cicd)
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

## 🔄 CI/CD

GitHub Actionsを使用したCI/CDパイプラインを構成しています。

### CI（継続的インテグレーション）

📁 **ワークフロー定義**: [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

`main`ブランチへのPull Request作成時に自動実行されます。

| ツール | 目的 |
|--------|------|
| hadolint | Dockerfileの静的解析・ベストプラクティスチェック |
| trivy | セキュリティ脆弱性スキャン（HIGH/CRITICAL） |

**スキップ対象**: 以下のファイルのみの変更時はCIを実行しません（Dockerビルドに影響しないため）

- `*.md`（README等のドキュメント）
- `docs/**`
- `LICENSE`
- `terraform/**`（インフラ定義）

### CD（継続的デリバリー）

📁 **ワークフロー定義**: [`.github/workflows/cd.yml`](.github/workflows/cd.yml)

タグ（`v*`）をpushすると、AWS ECRにDockerイメージが自動デプロイされます。

```bash
# 例：v1.0.0 タグを作成してpush
git tag v1.0.0
git push origin v1.0.0
```

### AWS設定（OIDC認証）

CDを利用するには、AWS側で以下の設定が必要です。
**Terraform（推奨）** または **手動設定** のいずれかを選択してください。

---

#### 方法A: Terraformで構築（推奨）

`terraform/` ディレクトリにIaCコードが用意されています。

```bash
cd terraform

# 1. セットアップスクリプトを実行（OIDC Provider の有無を自動判定）
./setup.sh

# 2. バックエンド設定ファイルを作成
cp backend.tf.example backend.tf
# backend.tf を編集してS3バケット名を設定

# 3. 初期化 & 適用
terraform init
terraform plan
terraform apply

# 4. 出力されるIAMロールARNをGitHubシークレットに設定
# リポジトリの Settings > Secrets and variables > Actions で以下を設定：
# github_actions_role_arn の値を AWS_ROLE_ARN として設定
```

> **Note**: `setup.sh` は AWS CLI で OIDC Provider の存在を確認し、`terraform.tfvars` を自動生成します。

**変数一覧：**
| 変数名 | 必須 | デフォルト | 説明 |
|--------|------|-----------|------|
| `github_repository` | ✅ | - | GitHubリポジトリ（`owner/repo`形式） |
| `create_oidc_provider` | - | `true` | OIDC Providerを新規作成するか（既存がある場合は`false`） |
| `ecr_repository_name` | - | `go-web-server-docker-example` | ECRリポジトリ名 |
| `aws_region` | - | `ap-northeast-1` | AWSリージョン |

**Terraformで作成されるリソース：**
- ECRリポジトリ（ライフサイクルポリシー付き）
- GitHub OIDC Provider（`create_oidc_provider=true`の場合のみ）
- IAMロール（ECR push権限付き）

---

#### 方法B: 手動で構築

##### 1. ECRリポジトリの作成

```bash
aws ecr create-repository \
  --repository-name go-web-server-docker-example \
  --region ap-northeast-1
```

##### 2. GitHub OIDC Providerの設定

**AWSコンソール:**

1. IAM > ID プロバイダ > プロバイダを追加
2. プロバイダのタイプ: OpenID Connect
3. プロバイダの URL: `https://token.actions.githubusercontent.com`
4. 対象者: `sts.amazonaws.com`

**AWS CLI:**

```bash
# OIDC プロバイダーの作成
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com
```

- OIDCの標準仕様で、JWTトークンには aud (audience) クレームが含まれます。
- これは「このトークンは誰（どのサービス）向けに発行されたのか」を示します。
- セキュリティのため、AWS STS向けトークンだけを受け付けるようにする

GitHub Actionsが発行するJWTトークンの例:
```json
{
  "iss": "https://token.actions.githubusercontent.com",
  "sub": "repo:my-org/my-repo:ref:refs/heads/main",
  "aud": "sts.amazonaws.com",  // ← これがaudience
  "exp": 1234567890
}
```

##### 3. IAMロールの作成

以下の信頼ポリシーでIAMロールを作成します（`OWNER/REPO`を実際の値に置き換えてください）：

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:OWNER/REPO:*"
        }
      }
    }
  ]
}
```

IAMロールに以下のポリシーをアタッチ：

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "arn:aws:ecr:ap-northeast-1:ACCOUNT_ID:repository/go-web-server-docker-example"
    }
  ]
}
```

##### 4. GitHubシークレットの設定

リポジトリの Settings > Secrets and variables > Actions で以下を設定：

| シークレット名 | 値 |
|---------------|-----|
| `AWS_ROLE_ARN` | 作成したIAMロールのARN |

## 👨‍💻 開発者向け情報

### プロジェクト構成

```
.
├── .github/
│   └── workflows/
│       ├── ci.yml            # CI ワークフロー（hadolint, trivy）
│       └── cd.yml            # CD ワークフロー（ECR push）
├── terraform/
│   ├── setup.sh              # セットアップスクリプト（OIDC自動判定）
│   ├── backend.tf.example    # バックエンド設定テンプレート
│   ├── terraform.tfvars.example  # 変数値テンプレート
│   ├── versions.tf           # Terraform/プロバイダーバージョン
│   ├── variables.tf          # 変数定義
│   ├── main.tf               # メインリソース（ECR, OIDC, IAM）
│   └── outputs.tf            # 出力定義
├── main.go                   # メインアプリケーション
├── go.mod                    # Go モジュール定義
├── Dockerfile                # Dockerイメージ定義
├── .dockerignore             # Docker除外ファイル設定
├── .hadolint.yaml            # hadolint設定（Dockerfile Linter）
└── README.md                 # このファイル
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
