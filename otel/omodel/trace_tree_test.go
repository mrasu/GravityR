package omodel_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraceTree_Compact(t *testing.T) {
	cases := map[string]struct {
		tree *omodel.TraceTree
		want *omodel.TraceTree
	}{
		"removes the same resource": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{2},
						ServiceName: "service",
						Children: []*omodel.Span{
							{
								TraceId:     []byte{111},
								SpanId:      []byte{3},
								ServiceName: "service",
							},
						},
					},
				},
			}}),
			want: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
			}},
		},
		"keeps the last db": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId: []byte{111},
						SpanId:  []byte{2},
						SpanAttributes: map[string]*omodel.AnyValue{
							dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
						},
						ServiceName: "service",
					},
				},
			}}),
			want: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{2},
						ServiceName: "some-db (service)",
						SpanAttributes: map[string]*omodel.AnyValue{
							dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
						},
					},
				},
			}},
		},
		"remove the db.system span when it is not the last": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{2},
						ServiceName: "service",
						SpanAttributes: map[string]*omodel.AnyValue{
							dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
						},
						Children: []*omodel.Span{
							{
								TraceId:     []byte{111},
								SpanId:      []byte{3},
								ServiceName: "some-db",
							},
						},
					},
				},
			}}),
			want: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{3},
						ServiceName: "some-db",
					},
				},
			}},
		},
		"keeps traces when traces visit the same resource": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{2},
						ServiceName: "service",
						Children: []*omodel.Span{
							{
								TraceId:     []byte{111},
								SpanId:      []byte{3},
								ServiceName: "other-service",
							},
							{
								TraceId:     []byte{111},
								SpanId:      []byte{4},
								ServiceName: "other-service",
							},
						},
					},
				},
			}}),
			want: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						TraceId:     []byte{111},
						SpanId:      []byte{3},
						ServiceName: "other-service",
					},
					{
						TraceId:     []byte{111},
						SpanId:      []byte{4},
						ServiceName: "other-service",
					},
				},
			}},
		},
		"keeps when traces visited only root": {
			tree: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
			}},
			want: &omodel.TraceTree{Root: &omodel.Span{
				TraceId:     []byte{111},
				SpanId:      []byte{1},
				ServiceName: "service",
			}},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.tree.Compact()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestTraceTree_MaxSameServiceAccessCount(t *testing.T) {
	cases := map[string]struct {
		tree *omodel.TraceTree
		want int
	}{
		"the number of visited service": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						ServiceName: "service",
						Children: []*omodel.Span{
							{
								ServiceName: "service",
							},
						},
					},
				},
			}}),
			want: 3,
		},
		"the max number of visited the same service": {
			tree: fillParent(&omodel.TraceTree{Root: &omodel.Span{
				ServiceName: "service",
				Children: []*omodel.Span{
					{
						ServiceName: "other-service",
						Children: []*omodel.Span{
							{
								ServiceName: "service",
							},
						},
					},
				},
			}}),
			want: 2,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tree.MaxSameServiceAccessCount())
		})
	}
}

func fillParent(tree *omodel.TraceTree) *omodel.TraceTree {
	queue := []*omodel.Span{tree.Root}

	for len(queue) > 0 {
		span := queue[0]
		queue = queue[1:]

		for _, c := range span.Children {
			c.Parent = span
			queue = append(queue, c)
		}
	}

	return tree
}
