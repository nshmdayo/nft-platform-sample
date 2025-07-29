# NFT Platform for Academic Papers - Backend

学術論文のためのNFTプラットフォームのバックエンドAPI

## 概要

このプロジェクトは、Go言語で実装された学術論文のNFTプラットフォームのバックエンドです。研究者が論文を投稿し、査読を受け、NFTとして発行できるシステムを提供します。

## 機能

- **ユーザー認証**: JWT認証によるユーザー登録・ログイン
- **論文管理**: 論文の作成、編集、削除、検索
- **査読システム**: ピアレビューとスコアリング機能
- **NFT統合**: 論文の査読完了後のNFT変換（準備中）

## 技術スタック

- **言語**: Go 1.22+
- **Webフレームワーク**: Gin
- **データベース**: PostgreSQL with GORM
- **認証**: JWT
- **ログ**: slog (JSON format)

## アーキテクチャ

Clean Architectureパターンを採用:
- `models/`: データモデル
- `repository/`: データアクセス層
- `service/`: ビジネスロジック層
- `handlers/`: HTTPハンドラー
- `middleware/`: ミドルウェア
- `router/`: ルーティング設定

## セットアップ

### 環境変数

`.env`ファイルを作成し、以下の環境変数を設定してください:

```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/nft_platform?sslmode=disable

# Server
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# IPFS (将来の実装用)
IPFS_API_URL=http://localhost:5001

# Ethereum (将来の実装用)
ETHEREUM_RPC_URL=http://localhost:8545
ETHEREUM_PRIVATE_KEY=your-private-key
PAPER_CONTRACT_ADDRESS=0x...
REVIEW_CONTRACT_ADDRESS=0x...
```

### データベースセットアップ

PostgreSQLを起動し、データベースを作成:

```bash
createdb nft_platform
```

### アプリケーション実行

```bash
# 依存関係のインストール
go mod tidy

# アプリケーション実行
go run cmd/server/main.go
```

## API エンドポイント

### 認証

- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `GET /api/v1/auth/profile` - プロフィール取得（認証必要）

### 論文

- `POST /api/v1/papers` - 論文作成（認証必要）
- `GET /api/v1/papers` - 論文一覧取得
- `GET /api/v1/papers/my` - 自分の論文一覧（認証必要）
- `GET /api/v1/papers/search?q=query` - 論文検索
- `GET /api/v1/papers/:id` - 論文詳細取得
- `PUT /api/v1/papers/:id` - 論文更新（認証必要）
- `DELETE /api/v1/papers/:id` - 論文削除（認証必要）
- `POST /api/v1/papers/:id/submit` - 査読提出（認証必要）

### 査読

- `POST /api/v1/reviews` - 査読作成（認証必要）
- `GET /api/v1/reviews/my` - 自分の査読一覧（認証必要）
- `GET /api/v1/reviews/pending` - 査読待ち論文一覧（認証必要）
- `GET /api/v1/reviews/:id` - 査読詳細取得（認証必要）
- `PUT /api/v1/reviews/:id` - 査読更新（認証必要）
- `DELETE /api/v1/reviews/:id` - 査読削除（認証必要）
- `GET /api/v1/papers/:paper_id/reviews` - 論文の査読一覧（認証必要）
- `GET /api/v1/papers/:paper_id/score` - 論文スコア取得（認証必要）

## プロジェクト構造

```
cmd/
  server/
    main.go           # アプリケーションエントリーポイント
internal/
  config/
    config.go         # 設定管理
  database/
    connection.go     # データベース接続
  handlers/
    auth_handler.go   # 認証ハンドラー
    paper_handler.go  # 論文ハンドラー
    review_handler.go # 査読ハンドラー
  middleware/
    auth.go          # 認証ミドルウェア
    middleware.go    # その他ミドルウェア
  models/
    user.go          # ユーザーモデル
    paper.go         # 論文モデル
    review.go        # 査読モデル
    nft_metadata.go  # NFTメタデータモデル
  repository/
    user_repository.go    # ユーザーリポジトリ
    paper_repository.go   # 論文リポジトリ
    review_repository.go  # 査読リポジトリ
  router/
    router.go        # ルーティング設定
  service/
    auth_service.go  # 認証サービス
    paper_service.go # 論文サービス
    review_service.go # 査読サービス
  utils/
    jwt.go           # JWT ユーティリティ
    password.go      # パスワードハッシュ
pkg/
  logger/
    logger.go        # ログ設定
```

## 今後の実装予定

- [ ] IPFS統合によるファイルストレージ
- [ ] Ethereumスマートコントラクト統合
- [ ] NFT発行機能
- [ ] ファイルアップロード機能
- [ ] 通知システム
- [ ] 詳細なアクセス制御

## 開発

### テスト実行

```bash
go test ./...
```

### ビルド

```bash
go build -o bin/server cmd/server/main.go
```

## ライセンス

MIT License