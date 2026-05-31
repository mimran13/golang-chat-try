package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go-chat/internal/config"
)

func connectDB(dbConfig config.DBConfig) *sql.DB {
      dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
          dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name,
      )
      db, err := sql.Open("mysql", dsn)
      if err != nil {
          slog.Error("failed to open db", "error", err)
          os.Exit(1)
      }
      if err := db.Ping(); err != nil {
          slog.Error("failed to ping db", "error", err)
          os.Exit(1)
      }
      slog.Info("database connected")
      db.SetMaxOpenConns(25)
      db.SetMaxIdleConns(10)
      db.SetConnMaxLifetime(5 * time.Minute)
      return db
  }
  
func runMigrations(dbConfig config.DBConfig) {

  dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
 
m, err := migrate.New("file://migrations", "mysql://"+dsn)
      if err != nil {
          slog.Error("failed to init migrations", "error", err)
          os.Exit(1)
      }
      if err := m.Up(); err != nil && err != migrate.ErrNoChange {
          slog.Error("failed to run migrations", "error", err)
          os.Exit(1)
      }
      slog.Info("migrations applied successfully")
}