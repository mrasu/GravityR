// Make possible to work autocomplete for the instance of ApexChart().
import "apexcharts";

type ChartW<T> = { config: { series: T } };
type ChartContext<T> = { w: ChartW<T> };
type ChartPoint<T> = {
  dataPointIndex: number;
  seriesIndex: number;
  w: ChartW<T>;
};
type ChartConfig = { seriesIndex: number; dataPointIndex: number };
