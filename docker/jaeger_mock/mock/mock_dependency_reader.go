package main

import (
	"context"
	"errors"
	"github.com/jaegertracing/jaeger/model"
	"time"
)

type mockDependencyReader struct{}

func (m *mockDependencyReader) GetDependencies(context.Context, time.Time, time.Duration) ([]model.DependencyLink, error) {
	return nil, errors.New("unexpected call: GetDependencies")
}
