package initializer

import (
	"time"

	"log/slog"

	logging "github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/log"
)

// 初期化したロガーをパッケージ全体で使用可能にする
var Logger *slog.Logger

func init() {
	// ロガーの初期化 (log.NewLoggerを使用)
	Logger = logging.NewLogger()

	// タイムゾーンを設定
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60) // 日本時間に設定
}
