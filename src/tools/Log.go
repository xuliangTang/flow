package tools

import (
	"flow/src/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"sync"
	"time"
)

var logging *LoggingImpl
var logger *zap.Logger
var loggerOnce sync.Once

type LevelEnablerFunc func(lvl *zapcore.Level) bool

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

type LoggingImpl struct {
	*zap.Logger
}

// Logger 获取日志对象
func Logger() *LoggingImpl {
	loggerOnce.Do(func() {
		// 设置多log文件和轮转
		tops := getTops()
		var cores []zapcore.Core
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}

		for _, top := range tops {
			top := top

			lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return top.Lef(&lvl)
			})

			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   top.Filename,
				MaxSize:    top.Ropt.MaxSize,
				MaxBackups: top.Ropt.MaxBackups,
				MaxAge:     top.Ropt.MaxAge,
				Compress:   top.Ropt.Compress,
			})

			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg.EncoderConfig),
				zapcore.AddSync(w),
				lv,
			)
			cores = append(cores, core)
		}

		logger = zap.New(zapcore.NewTee(cores...))
		defer logger.Sync() // flushes buffer, if any
		logging = &LoggingImpl{logger}
	})

	return logging
}

// 重写父类方法
func (this *LoggingImpl) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("stack", this.GetStack()))

	this.Logger.Error(msg, fields...)
}

// GetStack 获取堆栈信息
func (this *LoggingImpl) GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("==> %s\n", string(buf[:n]))
}

func getTops() []TeeOption {
	var tops = []TeeOption{
		{
			Filename: conf.Config.Server.AppPath + "/storage/logs/access.log",
			Ropt: RotateOptions{
				MaxSize:    255,
				MaxAge:     60,
				MaxBackups: 5,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl <= zapcore.InfoLevel
			},
		},
		{
			Filename: conf.Config.Server.AppPath + "/storage/logs/error.log",
			Ropt: RotateOptions{
				MaxSize:    255,
				MaxAge:     90,
				MaxBackups: 5,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl > zapcore.InfoLevel
			},
		},
	}

	return tops
}

// RequestHandler 请求日志
func RequestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		execTime := endTime.Sub(startTime) // 响应时间

		requestMethod := ctx.Request.Method
		requestURI := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		requestIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()
		Logger().Info(fmt.Sprintf("%s - %s %s[%d]", requestIP, requestMethod, requestURI, statusCode),
			zap.String("execTime", execTime.String()),
			zap.String("userAgent", userAgent),
		)
	}
}
