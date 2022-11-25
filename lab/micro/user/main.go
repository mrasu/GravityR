package main

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	dbDSN = "root:@tcp(user_db_mysql:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	err := initDB(dbDSN)
	if err != nil {
		panic(err)
	}

	err = initTraceProvider("http://jaeger:14268/api/traces")
	if err != nil {
		panic(err)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(usersHandler), "users")
	http.Handle("/users/", otelHandler)
	fmt.Println("Running a user service at http://localhost:9002")
	http.ListenAndServe(":9002", nil)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	db := dbWithCtx(r.Context())
	fmt.Println("=========")
	fmt.Println(r.Header)
	fmt.Println(trace.SpanContextFromContext(r.Context()).TraceID().String())

	sub := strings.TrimPrefix(r.URL.Path, "/users")
	_, file := filepath.Split(sub)
	id, err := strconv.Atoi(file)
	if err != nil {
		panic(err)
	}
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.Key("user_id").Int(id))
	fmt.Printf("getting user data for id=%d", id)

	var u *User
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		panic(err)
	}
	var histories []*PaymentHistory
	if err := db.Where("user_id = ?", u.ID).Find(&histories).Error; err != nil {
		panic(err)
	}

	var resHistories []map[string]any
	for _, h := range histories {
		var m *PaymentMethod
		if err := db.Where("id = ?", h.PaymentMethodID).Find(&m).Error; err != nil {
			panic(err)
		}
		resHistories = append(resHistories, map[string]any{
			"amount":       h.Amount,
			"last_numbers": m.LastNumbers,
		})
	}

	res := map[string]any{
		"id":        u.ID,
		"name":      u.Name,
		"email":     u.Email,
		"histories": resHistories,
	}
	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	w.Write(out)
}
