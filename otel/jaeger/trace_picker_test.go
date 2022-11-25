package jaeger_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/otel/jaeger"
	"github.com/mrasu/GravityR/otel/omodel"
	"testing"
)

var simpleTrace = &omodel.TraceTree{
	Root: &omodel.Span{
		Name: "root",
		Children: []*omodel.Span{
			{
				Name: "child1",
			},
		},
	},
}

var n1Trace = &omodel.TraceTree{
	Root: &omodel.Span{
		Name:        "root",
		ServiceName: "root",
		Children: []*omodel.Span{
			{Name: "child1", ServiceName: "child"},
			{Name: "child2", ServiceName: "child"},
			{Name: "child3", ServiceName: "child"},
			{Name: "child4", ServiceName: "child"},
			{Name: "child5", ServiceName: "child"},
			{Name: "child6", ServiceName: "child"},
		},
	},
}

func TestTracePicker_PickSameResourceAccessTraces(t *testing.T) {
	cases := map[string]struct {
		threshold int
		traces    []*omodel.TraceTree
		want      []*omodel.TraceTree
	}{
		"pick more than threshold": {
			threshold: 5,
			traces: []*omodel.TraceTree{
				simpleTrace, n1Trace,
			},
			want: []*omodel.TraceTree{
				n1Trace,
			},
		},
		"not pick when threshold is higher": {
			threshold: 100,
			traces: []*omodel.TraceTree{
				simpleTrace, n1Trace,
			},
			want: []*omodel.TraceTree{},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			tp := jaeger.NewTracePicker(tt.traces)
			traces := tp.PickSameServiceAccessTraces(tt.threshold)
			if diff := cmp.Diff(tt.want, traces); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestTracePicker_PickSlowTraces(t *testing.T) {
	traces := []*omodel.TraceTree{
		buildTimeTrace(0),
		buildTimeTrace(10),
		buildTimeTrace(20),
		buildTimeTrace(30),
		buildTimeTrace(40),
		buildTimeTrace(4),
		buildTimeTrace(3),
		buildTimeTrace(2),
		buildTimeTrace(1),
	}
	cases := map[string]struct {
		num    int
		traces []*omodel.TraceTree
		want   []*omodel.TraceTree
	}{
		"pick ordered slow traces": {
			num:    3,
			traces: traces,
			want: []*omodel.TraceTree{
				traces[4], traces[3], traces[2],
			},
		},
		"pick all when num is higher the number of traces": {
			num:    10,
			traces: traces,
			want: []*omodel.TraceTree{
				traces[4], traces[3], traces[2], traces[1], traces[5], traces[6], traces[7], traces[8], traces[0],
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			tp := jaeger.NewTracePicker(tt.traces)
			traces := tp.PickSlowTraces(tt.num)
			if diff := cmp.Diff(tt.want, traces); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func buildTimeTrace(duration int) *omodel.TraceTree {
	return &omodel.TraceTree{
		Root: &omodel.Span{
			Name:              "root",
			StartTimeUnixNano: 0,
			EndTimeUnixNano:   uint64(duration),
		},
	}
}
