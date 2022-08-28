import ApexCharts from "apexcharts";
import type { ApexOptions } from "apexcharts";
import type {
  ChartConfig,
  ChartContext,
  ChartPoint,
} from "../types/apexcharts";

export type GrChartPoint<T> = ChartPoint<GrChartSeries<T>>;

export type GrChartOptions<T> = Omit<ApexOptions, "series"> & {
  series: GrChartSeries<T>;
};

export type GrChartSeries<T> = Omit<ApexAxisChartSeries[number], "data"> &
  {
    data: GrChartSeriesData<T>[];
  }[];

export type GrChartSeriesData<T> =
  ApexAxisChartSeries[number]["data"][number] & {
    prop?: T;
  };

export abstract class GrChart<T> {
  private chart?: ApexCharts;

  createTooltipFn = (_: T): string => "";
  onClickedFn = (_: T): void => {};
  onDataPointMouseEnter = (_: T): void => {};

  protected constructor(private readonly elm: HTMLElement) {}

  destroy() {
    this.chart?.destroy();
  }

  protected renderChart(option: GrChartOptions<T>) {
    this.destroy();

    this.assignOptionFunctions(option);
    this.chart = new ApexCharts(this.elm, option);
    this.chart.render();
  }

  private assignOptionFunctions(option: GrChartOptions<T>) {
    option.tooltip.custom = this.tooltipCustom;
    option.legend.formatter = this.legendFormatter;

    option.chart.events = {
      click: this.eventsClick,
      dataPointMouseEnter: this.dataPointMouseEnter,
    };

    this.assignFunctions(option);
  }

  protected abstract assignFunctions(option: GrChartOptions<T>): void;

  private tooltipCustom = ({
    seriesIndex,
    dataPointIndex,
    w: {
      config: { series },
    },
  }: ChartPoint<GrChartSeries<T>>): string => {
    const data = series[seriesIndex].data[dataPointIndex];
    return this.createTooltipFn(data.prop);
  };

  private legendFormatter = (seriesName: string): string => {
    if (seriesName.length <= 20) {
      return seriesName;
    }

    return seriesName.substring(0, 20) + "...";
  };

  private eventsClick = (
    _event: any,
    {
      w: {
        config: { series },
      },
    }: ChartContext<GrChartSeries<T>>,
    { seriesIndex, dataPointIndex }: ChartConfig
  ) => {
    if (seriesIndex < 0) return;

    this.onClickedFn(series[seriesIndex].data[dataPointIndex].prop);
  };

  private dataPointMouseEnter = (
    event: { path: SVGMPathElement[] },
    {
      w: {
        config: { series },
      },
    }: ChartContext<T>,
    { seriesIndex, dataPointIndex }: ChartConfig
  ) => {
    event.path[0].style.cursor = "pointer";

    this.onDataPointMouseEnter(series[seriesIndex].data[dataPointIndex].prop);
  };
}
