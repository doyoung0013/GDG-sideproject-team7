package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/GDG-KHU-Side/backend-side-project/config"
)

// Open은 MySQL 커넥션 풀을 만들고 연결을 확인합니다.
func Open(ctx context.Context, conf config.DBConfig) (*sql.DB, error) {
	mc := mysql.NewConfig()
	mc.User = conf.User
	mc.Passwd = conf.Password
	mc.Net = "tcp"
	mc.Addr = conf.Host + ":" + conf.Port
	mc.DBName = conf.DBName
	mc.ParseTime = true
	mc.ClientFoundRows = true // UPDATE 시 "변경된 행"이 아니라 "일치한 행" 수를 반환
	mc.Params = map[string]string{"charset": "utf8mb4"}

	conn, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := conn.PingContext(pingCtx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return conn, nil
}
