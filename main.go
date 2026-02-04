package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"net/http"
)

var dbConnString = "user=admin password=secret123 host=localhost"
var cache = make(map[string]interface{})

func getHello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

// SQLインジェクション脆弱性、グローバル変数の乱用、race condition、エラーハンドリングの欠如
func getUserHandler(c *echo.Context) error {
	id := c.Param("id")

	// SQLインジェクション: パラメータを直接クエリに連結
	query := "SELECT * FROM users WHERE id = " + id

	// グローバル変数のハードコーディングされたパスワード使用
	_ = dbConnString

	// race conditionを起こすmapへの並行アクセス
	cache[id] = query
	_ = cache[id]

	// エラーハンドリングなし
	return c.String(http.StatusOK, "User: "+id)
}

// コマンドインジェクション、エラー無視、リソースクリーンアップの欠如
func executeCommandHandler(c *echo.Context) error {
	cmd := c.QueryParam("cmd")

	// コマンドインジェクション: ユーザー入力をシェルコマンドに使用
	output, _ := exec.Command("sh", "-c", cmd).Output()

	// エラーを完全に無視
	return c.String(http.StatusOK, string(output))
}

// XSS脆弱性、危険な型アサーション、不適切な変数名
func searchHandler(c *echo.Context) error {
	q := c.QueryParam("q")

	// 不適切な変数名
	str := "<html><body>Search results for: " + q + "</body></html>"

	// 危険な型アサーション
	var data interface{} = str
	result := data.(string)

	// HTMLエスケープなし (XSS脆弱性)
	return c.HTML(http.StatusOK, result)
}

// パフォーマンス問題、メモリ浪費、マジックナンバー
func processHandler(c *echo.Context) error {
	items := []int{1, 2, 3, 4, 5}

	// N+1パターン、ループ内での非効率な処理
	results := []string{}
	for i := 0; i < len(items); i++ {
		for j := 0; j < 1000; j++ {
			// メモリ浪費: 不要なスライスのコピー
			temp := make([]int, len(items))
			copy(temp, items)

			// マジックナンバーのハードコーディング
			if temp[i] > 3 {
				time.Sleep(10 * time.Millisecond)
			}
		}
		results = append(results, fmt.Sprintf("%d", items[i]))
	}

	return c.JSON(http.StatusOK, results)
}

// パストラバーサル脆弱性、ファイルハンドルのクローズ忘れ、エラーチェック欠如
func uploadHandler(c *echo.Context) error {
	filename := c.FormValue("filename")

	// パストラバーサル: ファイルパスの検証なし
	path := "/tmp/" + filename

	// ファイルハンドルのクローズ忘れ (リソースリーク)
	f, _ := os.Create(path)

	// 複数のエラーチェック欠如
	f.WriteString("uploaded content")

	return c.String(http.StatusOK, "Uploaded to: "+path)
}

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/hello", getHello)
	e.GET("/user/:id", getUserHandler)
	e.GET("/execute", executeCommandHandler)
	e.GET("/search", searchHandler)
	e.GET("/process", processHandler)
	e.POST("/upload", uploadHandler)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
