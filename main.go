package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// レスポンス用の構造体
type HealthResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type HelloResponse struct {
	Message string `json:"message"`
	Name    string `json:"name,omitempty"`
}

// ヘルスチェックエンドポイント
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status:  "ok",
		Message: "Server is running",
		Time:    time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// ルートエンドポイント
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HelloResponse{
		Message: "Welcome to Simple Go Web Server!",
	}
	json.NewEncoder(w).Encode(response)
}

// 挨拶エンドポイント
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	response := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
		Name:    name,
	}
	json.NewEncoder(w).Encode(response)
}

// ログミドルウェア
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func main() {
	// コマンドライン引数のパース
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// 環境変数からポート番号を取得（コマンドライン引数より優先度は低い）
	if envPort := os.Getenv("PORT"); envPort != "" && *port == "8080" {
		port = &envPort
	}

	// ルーターの設定
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hello", helloHandler)

	// ミドルウェアを適用
	handler := loggingMiddleware(mux)

	// サーバーの起動
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
