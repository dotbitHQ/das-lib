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

var log = mylog.NewLogger("gorm", mylog.LevelDebug)

func NewGormDB(addr, user, password, dbName string, maxOpenConn, maxIdleConn int) (*gorm.DB, error) {
	conn := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(conn, user, password, addr, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &GormLogger{},
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

type GormLogger struct{}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(msg, data...)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Errorf(msg, data...)
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
	if err != nil {

		sqlInfo.Err = err
		sqlInfoByte, _ := json.Marshal(sqlInfo)
		log.Error(string(sqlInfoByte))

	} else {
		sqlInfoByte, _ := json.Marshal(sqlInfo)
		log.Info(string(sqlInfoByte))
	}
}
