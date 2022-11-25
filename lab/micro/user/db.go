package main

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"
)

var globalDB *gorm.DB

func initDB(dsn string) error {
	var err error
	globalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := globalDB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return err
	}

	return nil
}

func dbWithCtx(ctx context.Context) *gorm.DB {
	return globalDB.WithContext(ctx)
}

type User struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PaymentHistory struct {
	ID              int32
	UserID          int32
	PaymentMethodID int32
	Amount          int32
	PayedAt         time.Time
}

type PaymentMethod struct {
	ID          int32
	UserID      int32
	LastNumbers string
}
