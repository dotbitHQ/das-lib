package http_api

import (
	"context"
	"encoding/json"
	"fmt"
	mylog "github.com/dotbitHQ/das-lib/http_api/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewGormDB(addr, user, password, dbName string, maxOpenConn, maxIdleConn int) (*gorm.DB, error) {
	return NewGormDBWithLog(addr, user, password, dbName, maxOpenConn, maxIdleConn, nil)
}

func NewGormDBWithLog(addr, user, password, dbName string, maxOpenConn, maxIdleConn int, log *mylog.Logger) (*gorm.DB, error) {
	conn := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(conn, user, password, addr, dbName)

	if log == nil {
		log = mylog.NewLogger("gorm", mylog.LevelDebug)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &GormLogger{
			Log: log,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open :%v", err)
	}
	db = db.Debug()
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm db :%v", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)
	return db, nil
}

type GormLogger struct {
	Log *mylog.Logger
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.Log.Infof(msg, data...)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.Log.Warnf(msg, data...)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.Log.Errorf(msg, data...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin).Milliseconds()
	sql, rows := fc()
	sqlInfo := struct {
		Elapsed interface{}
		Rows    interface{}
		Err     error
		Sql     string
	}{
		Elapsed: elapsed,
		Rows:    rows,
		Sql:     sql,
	}
	if err != nil && err.Error() != "" {
		sqlInfo.Err = err
		sqlInfoByte, _ := json.Marshal(sqlInfo)
		g.Log.Error(string(sqlInfoByte))
	} else {
		sqlInfoByte, _ := json.Marshal(sqlInfo)
		g.Log.Info(string(sqlInfoByte))
	}
}
