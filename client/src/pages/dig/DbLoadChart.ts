import type { ApexOptions } from "apexcharts";
import type { TimeDbLoad } from "../../models/TimeDbLoad";
import { GrChart } from "../../lib/GrChart";
import type { GrChartOptions } from "../../lib/GrChart";
import type { GrChartSeries } from "../../lib/GrChart";
import { createLoadPropChartSeries } from "./util/chart/LoadPropChartSereies";

const DEFAULT_CHART_OPTION: ApexOptions = {
  series: [],
  chart: {
    type: "bar",
    stacked: true,
    zoom: { enabled: true },
    animations: {
      enabled: false,
    },
    events: {},
    height: 350,
    width: 3000,
  },
  tooltip: {
    enabled: true,
  },
  dataLabels: {
    enabled: false,
  },
  xaxis: {
    type: "datetime",
    tickPlacement: "on",
  },
  yaxis: {
    labels: {
      show: true,
    },
    title: {
      text: "Load of database (%)",
      style: {
        fontSize: "15px",
      },
    },
  },
  legend: {
    show: true,
  },
};

export type DbLoadChartProp = {
  dbName: string;
  sql: string;
  loadMaxPercent: number;
  date: Date;
};

type Series = GrChartSeries<DbLoadChartProp>;

export class DbLoadChart extends GrChart<DbLoadChartProp> {
  constructor(
    elm: HTMLElement,
    private readonly timeDbLoads: TimeDbLoad[],
    private readonly threshold: number
  ) {
    super(elm);
  }

  render(databases: string[]) {
    const option = JSON.parse(JSON.stringify(DEFAULT_CHART_OPTION));
    option.series = this.createDbLoadChartSeries(databases);

    this.renderChart(option);
  }

  changeDatabase(databases: string[]) {
    this.render(databases);
  }

  protected assignFunctions(option: GrChartOptions<DbLoadChartProp>) {
    if ("labels" in option.yaxis) {
      option.yaxis.labels.formatter = this.yaxisLabelsFormatter;
    }
  }

  private createDbLoadChartSeries = (targetDatabases: string[]): Series => {
    return createLoadPropChartSeries<DbLoadChartProp>(
      this.timeDbLoads,
      targetDatabases,
      this.threshold,
      (timeDbLoad, dbLoad, loadOfSql) => ({
        x: timeDbLoad.date,
        y: loadOfSql.loadMax * 100,
        prop: {
          date: timeDbLoad.date,
          loadMaxPercent: loadOfSql.loadMax * 100,
          dbName: dbLoad.name,
          sql: loadOfSql.sql,
        },
      })
    );
  };

  private yaxisLabelsFormatter = (value: number): string => {
    return value.toFixed(0);
  };
}
