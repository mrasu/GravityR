package dig

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iJaeger "github.com/mrasu/GravityR/infra/jaeger"
	"github.com/mrasu/GravityR/lib"
	"github.com/mrasu/GravityR/otel/jaeger"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"time"
)

var JaegerCmd = &cobra.Command{
	Use:   "jaeger",
	Short: "Dig OpenTelemetry data stored in Jaeger",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := jaegerR.prepare()
		if err != nil {
			return err
		}

		return jaegerR.run()
	},
}

func init() {
	flg := JaegerCmd.Flags()
	flg.StringVarP(&jaegerR.serviceName, "service-name", "s", "", "[Required] name of a service to dig")
	flg.IntVarP(&jaegerR.durationMilli, "min-duration-millisecond", "d", 1000, "Minimum duration to search")
	flg.IntVarP(&jaegerR.sameServiceThreshold, "same-service-threshold", "t", 5, "Threshold for same-service-trace")
	flg.IntVarP(&jaegerR.count, "count", "c", 500, "The number of traces get from Jaeger")
	flg.StringVar(&jaegerR.startStr, "start-from", "", "Start time by RFC3339. (Default: 24 hours ago)")

	err := cobra.MarkFlagRequired(flg, "service-name")
	if err != nil {
		panic(err)
	}
}

var jaegerR = jaegerRunner{}

type jaegerRunner struct {
	serviceName          string
	durationMilli        int
	sameServiceThreshold int
	count                int

	startStr string
	start    time.Time
	end      time.Time

	grpcAddress string
	uiURL       string
}

func (jr *jaegerRunner) prepare() error {
	now := time.Now()

	if jr.startStr != "" {
		var err error
		jr.start, err = time.Parse(time.RFC3339, jr.startStr)
		if err != nil {
			return errors.Wrap(err, "invalid format at StartAt")
		}
	} else {
		jr.start = now.Add(-24 * time.Hour)
	}

	jr.end = jr.start.Add(24 * time.Hour)
	if jr.end.After(now) {
		jr.end = now
	}

	uiURL, err := lib.GetEnv("JAEGER_UI_URL")
	if err != nil {
		return err
	}
	jr.uiURL = uiURL

	grpcAddress, err := lib.GetEnv("JAEGER_GRPC_ADDRESS")
	if err != nil {
		return err
	}
	jr.grpcAddress = grpcAddress

	return nil
}

func (jr *jaegerRunner) run() error {
	cli, err := iJaeger.Open(jr.grpcAddress, false)
	if err != nil {
		return err
	}
	defer cli.Close()

	return jr.dig(flag.AppFlag.Output, cli)
}

func (jr *jaegerRunner) dig(outputPath string, cli *iJaeger.Client) error {
	log.Info().Msgf("Dig traces from %s to %s taking more than %d millis)", jr.start.Format(time.RFC3339), jr.end.Format(time.RFC3339), jr.durationMilli)
	tf := jaeger.NewTraceFetcher(cli)
	traces, err := tf.FetchCompactedTraces(int32(jr.count), jr.start, jr.end, jr.serviceName, time.Duration(jr.durationMilli)*time.Millisecond)
	if err != nil {
		return err
	}

	tp := jaeger.NewTracePicker(traces)
	slowTraces := tp.PickSlowTraces(100)
	sameServiceTraces := tp.PickSameServiceAccessTraces(jr.sameServiceThreshold)

	log.Debug().Msg("found traces: ----start")
	for _, tree := range slowTraces {
		s := tree.Root
		log.Printf(
			`[%x] "%s" at %s`,
			s.TraceId,
			s.Name,
			s.ServiceName,
		)
	}
	log.Debug().Msg("found traces: ----end ")

	if outputPath != "" {
		err = jr.createHTML(outputPath, slowTraces, sameServiceTraces)
		if err != nil {
			return err
		}

		util.LogResultOutputPath(outputPath)
	}

	return nil
}

func (jr *jaegerRunner) createHTML(outputPath string, slowTraces, sameServiceTraces []*omodel.TraceTree) error {
	slows := lo.Map(slowTraces, func(trace *omodel.TraceTree, _ int) *viewmodel.VmOtelCompactTrace {
		return trace.ToCompactViewModel()
	})
	sameServices := lo.Map(sameServiceTraces, func(trace *omodel.TraceTree, _ int) *viewmodel.VmOtelCompactTrace {
		return trace.ToCompactViewModel()
	})

	bo := html.NewDigJaegerBuildOption(
		jr.uiURL,
		jr.durationMilli,
		jr.sameServiceThreshold,
		slows,
		sameServices,
	)

	err := html.CreateHtml(html.TypeMermaid, outputPath, bo)
	if err != nil {
		return err
	}
	return nil
}
