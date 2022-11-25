package omodel

import (
	"bytes"
	"fmt"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/samber/lo"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"sort"
)

type Span struct {
	Parent            *Span
	Name              string
	TraceId           []byte
	SpanId            []byte
	StartTimeUnixNano uint64
	EndTimeUnixNano   uint64
	SpanAttributes    map[string]*AnyValue
	ServiceName       string
	Children          []*Span
}

func (s *Span) SortChildren() {
	sort.Slice(s.Children, func(i, j int) bool {
		return s.Children[i].StartTimeUnixNano < s.Children[j].StartTimeUnixNano
	})

	for _, c := range s.Children {
		c.SortChildren()
	}
}

func (s *Span) GetDBSystem() string {
	v := s.SpanAttributes[string(semconv.DBSystemKey)]
	if v == nil {
		return ""
	}
	return v.GetStringValue()
}

func (s *Span) IsLastDbSpan() bool {
	if len(s.Children) > 0 {
		return false
	}

	_, ok := s.SpanAttributes[string(semconv.DBSystemKey)]
	return ok
}

func (s *Span) ShallowCopyWithoutDependency() *Span {
	return &Span{
		Parent:            nil,
		Name:              s.Name,
		TraceId:           s.TraceId,
		SpanId:            s.SpanId,
		StartTimeUnixNano: s.StartTimeUnixNano,
		EndTimeUnixNano:   s.EndTimeUnixNano,
		SpanAttributes:    s.SpanAttributes,
		ServiceName:       s.ServiceName,
		Children:          nil,
	}
}

func (s *Span) IsSame(other *Span) bool {
	return bytes.Equal(s.TraceId, other.TraceId) && bytes.Equal(s.SpanId, other.SpanId)
}

func (s *Span) CalcConsumedTime() uint64 {
	var childTimes [][]uint64
	for _, c := range s.Children {
		childTimes = append(childTimes, []uint64{c.StartTimeUnixNano, c.EndTimeUnixNano})
	}
	lib.Sort(childTimes, func(t []uint64) uint64 { return t[0] })

	consumedTime := uint64(0)
	lastEndTime := s.StartTimeUnixNano
	for _, t := range childTimes {
		if lastEndTime < t[0] {
			consumedTime += t[0] - lastEndTime
		}
		if lastEndTime < t[1] {
			lastEndTime = t[1]
		}
	}
	consumedTime += s.EndTimeUnixNano - lastEndTime

	return consumedTime
}

func (s *Span) ToViewModel() *viewmodel.VmOtelTraceSpan {
	children := lo.Map(s.Children, func(c *Span, _ int) *viewmodel.VmOtelTraceSpan { return c.ToViewModel() })

	return &viewmodel.VmOtelTraceSpan{
		Name:            s.Name,
		SpanId:          fmt.Sprintf("%x", s.SpanId),
		StartTimeMillis: int(s.StartTimeUnixNano / 1000 / 1000),
		EndTimeMillis:   int(s.EndTimeUnixNano / 1000 / 1000),
		ServiceName:     s.ServiceName,
		Children:        children,
	}
}
