package main

import (
	"context"
	"errors"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

type mockSpanReader struct {
	getTracesFn func() ([]*model.Trace, error)
}

// GetTrace retrieves the trace with a given id.
//
// If no spans are stored for this trace, it returns ErrTraceNotFound.
func (m *mockSpanReader) GetTrace(context.Context, model.TraceID) (*model.Trace, error) {
	return nil, errors.New("unexpected call: GetTrace")
}

// GetServices returns all service names known to the backend from spans
// within its retention period.
func (m *mockSpanReader) GetServices(context.Context) ([]string, error) {
	return nil, errors.New("unexpected call: GetServices")
}

// GetOperations returns all operation names for a given service
// known to the backend from spans within its retention period.
func (m *mockSpanReader) GetOperations(context.Context, spanstore.OperationQueryParameters) ([]spanstore.Operation, error) {
	return nil, errors.New("unexpected call: GetOperations")
}

// FindTraces returns all traces matching query parameters. There's currently
// an implementation-dependent abiguity whether all query filters (such as
// multiple tags) must apply to the same span within a trace, or can be satisfied
// by different spans.
//
// If no matching traces are found, the function returns (nil, nil).
func (m *mockSpanReader) FindTraces(_ context.Context, param *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	traces, err := m.getTracesFn()
	if err != nil {
		return nil, err
	}
	var res []*model.Trace
	for _, t := range traces {
		for _, s := range t.Spans {
			if len(s.References) == 0 || s.References[0].SpanID == 0 {
				start := s.StartTime
				if param.StartTimeMin.Before(start) && param.StartTimeMax.After(start) {
					res = append(res, t)
					break
				}
			}
		}
	}
	return res, nil
}

// FindTraceIDs does the same search as FindTraces, but returns only the list
// of matching trace IDs.
//
// If no matching traces are found, the function returns (nil, nil).
func (m *mockSpanReader) FindTraceIDs(context.Context, *spanstore.TraceQueryParameters) ([]model.TraceID, error) {
	return nil, errors.New("unexpected call: FindTraceIDs")
}
