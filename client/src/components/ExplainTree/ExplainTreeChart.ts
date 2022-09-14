import type { ApexOptions } from "apexcharts";
import type { ExplainTree } from "./ExplainTree";
import type {
  GrChartOptions,
  GrChartPoint,
  GrChartSeriesData,
} from "@/lib/GrChart";
import { GrChart } from "@/lib/GrChart";
import type { DbExplainData } from "@/models/explain_data/DbExplainData";

const DEFAULT_CHART_OPTION: ApexOptions = {
  series: [
    {
      data: [],
    },
  ],
  dataLabels: {
    enabled: true,
    textAnchor: "middle",
  },
  grid: {
    xaxis: { lines: { show: true } },
    yaxis: { lines: { show: false } },
  },
  chart: {
    width: "100%",
    type: "rangeBar",
  },
  plotOptions: {
    bar: {
      horizontal: true,
      barHeight: "80%",
    },
  },
  fill: {
    type: "solid",
    opacity: 0.6,
  },
  xaxis: {
    type: "numeric",
    title: {
      text: "Execution time (ms)",
    },
  },
  yaxis: {
    show: true,
    labels: {
      show: false,
    },
    axisBorder: {
      show: false,
    },
    axisTicks: {
      show: false,
    },
  },
  stroke: {
    colors: ["transparent"],
    width: 0,
  },
  tooltip: {
    enabled: true,
  },
  legend: {
    show: true,
  },
};

const BAR_HEIGHT = 30;

export type ExplainTreeChartProp<D extends DbExplainData> = {
  xNum: number;
  IAnalyzeData: D;
};

export type ExplainTreeSeriesData<D extends DbExplainData> = GrChartSeriesData<
  ExplainTreeChartProp<D>
>;

export class ExplainTreeChart<D extends DbExplainData> extends GrChart<
  ExplainTreeChartProp<D>
> {
  constructor(
    elm: HTMLElement,
    private explainTree: ExplainTree<D>,
    private description: string
  ) {
    super(elm);
  }

  render() {
    const seriesData = this.explainTree.getSeriesData();
    const option = JSON.parse(JSON.stringify(DEFAULT_CHART_OPTION));
    option.series = [{ data: seriesData }];
    option.chart.height = this.calculateChartHeightPx(seriesData);
    option.title = { text: this.description };

    this.renderChart(option);
  }

  private calculateChartHeightPx = (
    seriesData: ExplainTreeSeriesData<D>[]
  ): string => {
    const barCount = new Set(seriesData.map((data) => data.prop.xNum)).size;
    const height = barCount * BAR_HEIGHT + 130;
    return `${height}px`;
  };

  protected assignFunctions(
    option: GrChartOptions<ExplainTreeChartProp<D>>
  ): void {
    option.dataLabels.formatter = this.dataLabelsFormatter;
  }

  private dataLabelsFormatter = (
    _: any,
    {
      seriesIndex,
      dataPointIndex,
      w: {
        config: { series },
      },
    }: GrChartPoint<ExplainTreeChartProp<D>>
  ) => {
    return series[seriesIndex].data[dataPointIndex].prop.IAnalyzeData.title;
  };
}
