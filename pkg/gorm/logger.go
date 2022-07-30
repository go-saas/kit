package gorm

import (
	"context"
	"errors"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DefaultSlowThreshold = 200 * time.Millisecond
)

// Logger GORM log adapter
type Logger struct {
	klog.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

var _ logger.Interface = (*Logger)(nil)

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Log(klog.LevelInfo, klog.DefaultMessageKey, fmt.Sprintf(s, i...))
	}
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Log(klog.LevelWarn, klog.DefaultMessageKey, fmt.Sprintf(s, i...))
	}
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Log(klog.LevelError, klog.DefaultMessageKey, fmt.Sprintf(s, i...))
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	th := l.SlowThreshold
	if th <= 0 {
		th = DefaultSlowThreshold
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.Log(klog.LevelError, klog.DefaultMessageKey, fmt.Sprintf("%s", err), "sql", sql, "rows", rows, "elapsed", float64(elapsed.Nanoseconds())/1e6)
	case elapsed > th && th != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()

		slowLog := fmt.Sprintf("SLOW SQL >= %v", th)
		l.Log(klog.LevelWarn, klog.DefaultMessageKey, slowLog, "sql", sql, "rows", rows, "elapsed", float64(elapsed.Nanoseconds())/1e6)

	case l.LogLevel == logger.Info:
		sql, rows := fc()
		l.Log(klog.LevelInfo, "sql", sql, "rows", rows, "elapsed", float64(elapsed.Nanoseconds())/1e6)

	}
}
