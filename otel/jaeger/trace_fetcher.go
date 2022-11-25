package jaeger

import (
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"github.com/mrasu/GravityR/infra/jaeger"
	"github.com/mrasu/GravityR/lib"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/rs/zerolog/log"
	"time"
)

type TraceFetcher struct {
	cli *jaeger.Client
}

func NewTraceFetcher(cli *jaeger.Client) *TraceFetcher {
	return &TraceFetcher{cli: cli}
}

func (tf *TraceFetcher) FetchCompactedTraces(size int32, start, end time.Time, serviceName string, durationMin time.Duration) ([]*omodel.TraceTree, error) {
	interval := 30 * time.Minute

	var res []*omodel.TraceTree
	times := lib.GenerateTimeRanges(start, end, interval)
	for _, start := range times {
		param := &api_v3.TraceQueryParameters{
			ServiceName:  serviceName,
			StartTimeMin: &types.Timestamp{Seconds: start.Unix()},
			StartTimeMax: &types.Timestamp{Seconds: start.Add(interval).Unix()},
			NumTraces:    size,
			DurationMin:  tf.toProtoDuration(durationMin),
		}
		log.Info().Msgf("Getting data for %s", start.Format(time.RFC3339))

		trees, err := tf.findCompactedTraceTrees(param)
		if err != nil {
			return nil, err
		}
		res = append(res, trees...)
	}

	return res, nil
}

func (tf *TraceFetcher) toProtoDuration(duration time.Duration) *types.Duration {
	secs := int64(duration.Truncate(time.Second).Seconds())
	nanos := int32(duration.Nanoseconds() - int64(time.Duration(secs)*time.Second))

	return &types.Duration{
		Seconds: secs,
		Nanos:   nanos,
	}
}

func (tf *TraceFetcher) findCompactedTraceTrees(param *api_v3.TraceQueryParameters) ([]*omodel.TraceTree, error) {
	traces, err := tf.cli.FindTraces(param)
	if err != nil {
		return nil, err
	}

	var compactedTrees []*omodel.TraceTree
	for _, trace := range traces {
		trees := ToTraceTrees(trace.Spans)
		for _, tree := range trees {
			compacted := tree.Compact()
			compactedTrees = append(compactedTrees, compacted)
		}
	}

	return compactedTrees, nil
}
