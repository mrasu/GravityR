import ApexCharts from "apexcharts";
import type { ApexOptions } from "apexcharts";
import type { ExplainAnalyzeTree } from "./ExplainAnalyzeTree";
import type { SeriesData } from "./SeriesData";

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
  title: {
    text: "Execution time based timeline from EXPLAIN ANALYZE",
  },
};

const BAR_HEIGHT = 30;

type ChartW = { config: { series: { data: SeriesData[] }[] } };
type ChartContext = { w: ChartW };
type ChartPoint = {
  dataPointIndex: number;
  seriesIndex: number;
  w: ChartW;
};
type ChartConfig = { seriesIndex: number; dataPointIndex: number };

export class ExplainTree {
  private chart?: ApexCharts;
  private animating = false;
  private mouseEnteringX?: number;

  public onMouseEntered = (_: number): void => {};
  public createTooltipFn = (_: SeriesData): string => "";

  constructor(
    private elm: HTMLElement,
    private explainAnalyzeTree: ExplainAnalyzeTree
  ) {}

  render() {
    const seriesData2 = this.explainAnalyzeTree.getFocusedSeriesData(
      this.explainAnalyzeTree.getSlowestX()
    );
    const option = this.createChartOption();
    option.series = [{ data: seriesData2 }];
    option.chart.height = this.calculateChartHeightPx(seriesData2);

    this.chart = new ApexCharts(this.elm, option);
    this.chart.render();
  }

  private createChartOption(): ApexOptions {
    const option = JSON.parse(JSON.stringify(DEFAULT_CHART_OPTION));

    option.tooltip.custom = this.tooltipCustom;
    option.dataLabels.formatter = this.dataLabelsFormatter;
    option.chart.events = {
      animationEnd: this.animationEnd,
      dataPointSelection: this.dataPointSelection,
      dataPointMouseEnter: this.dataPointMouseEnter,
      dataPointMouseLeave: this.dataPointMouseLeave,
    };

    return option;
  }

  private tooltipCustom = ({
    seriesIndex,
    dataPointIndex,
    w: {
      config: { series },
    },
  }: ChartPoint): string => {
    const data = series[seriesIndex].data[dataPointIndex];
    return this.createTooltipFn(data);
  };

  private dataLabelsFormatter = (
    _: any,
    {
      seriesIndex,
      dataPointIndex,
      w: {
        config: { series },
      },
    }: ChartPoint
  ) => {
    return series[seriesIndex].data[dataPointIndex].title;
  };

  private animationEnd = () => {
    this.animating = false;
    if (this.mouseEnteringX) {
      this.onMouseEntered(this.mouseEnteringX);
    }
  };

  private dataPointSelection = (
    _: any,
    {
      w: {
        config: { series },
      },
    }: ChartContext,
    { seriesIndex, dataPointIndex }: ChartConfig
  ) => {
    const x = series[seriesIndex].data[dataPointIndex].xNum;
    this.updateSeriesData(series[seriesIndex].data, x);
  };

  private dataPointMouseEnter = (
    event: { path: SVGMPathElement[] },
    {
      w: {
        config: { series },
      },
    }: ChartContext,
    { seriesIndex, dataPointIndex }: ChartConfig
  ) => {
    const x = series[seriesIndex].data[dataPointIndex].xNum;
    this.mouseEnteringX = x;
    if (this.animating) return;

    event.path[0].style.cursor = "pointer";
    this.onMouseEntered(x);
  };

  private dataPointMouseLeave = () => {
    this.mouseEnteringX = undefined;
  };

  private updateSeriesData = async (currentSeries: SeriesData[], x: number) => {
    return new Promise<void>((resolve) => {
      setTimeout(async () => {
        const nexts = this.explainAnalyzeTree.getNextXs(x);
        if (!nexts) return;

        const nextX = nexts[0];
        const isOpen = !!currentSeries.find((d) => d.xNum === nextX);
        let seriesData: SeriesData[];
        if (isOpen) {
          seriesData = this.explainAnalyzeTree.getFocusedSeriesData(x);
        } else {
          seriesData = this.explainAnalyzeTree.getFocusedSeriesData(nextX);
        }

        this.animating = true;
        await this.chart.updateOptions({
          series: [
            {
              data: seriesData,
            },
          ],
          chart: { height: this.calculateChartHeightPx(seriesData) },
        });
        resolve();
      });
    });
  };

  private calculateChartHeightPx = (seriesData: SeriesData[]): string => {
    const barCount = new Set(seriesData.map((data) => data.x)).size;
    const height = barCount * BAR_HEIGHT + 130;
    return `${height}px`;
  };
}