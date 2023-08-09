package db

import (
	"context"
	"time"

	"github.com/capell/capell_scan/lib/logger"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	opt           *Option
	slowThreshold time.Duration
}

func NewGormLogger(opt *Option) *GormLogger {
	st := time.Duration(opt.SlowThreshold * float64(time.Second))
	l := &GormLogger{
		opt:           opt,
		slowThreshold: st,
	}
	return l
}

func (l *GormLogger) LogMode(lvl glog.LogLevel) glog.Interface {
	n := *l
	return &n
}

func (l GormLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	logger.Info("%v "+fmt, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GormLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	logger.Warn("%v "+fmt, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GormLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	logger.Error("%v "+fmt, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	pos := utils.FileWithLineNum()
	if err != nil {
		logger.Error("%v sql:%v rows:%v err:%v", pos, sql, rows, err)
	}
	elapsed := time.Since(begin)
	if err == nil && elapsed >= l.slowThreshold {
		logger.Warn("%v sql:%v rows:%v slow:%v", pos, sql, rows, elapsed)
	}
}
