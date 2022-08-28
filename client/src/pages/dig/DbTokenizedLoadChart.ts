import type { ApexOptions } from "apexcharts";
import type { TimeDbLoad } from "@/models/TimeDbLoad";
import { GrChart } from "@/lib/GrChart";
import type { GrChartOptions } from "@/lib/GrChart";
import type { GrChartSeries } from "@/lib/GrChart";
import { createLoadPropChartSeries } from "./util/chart/LoadPropChartSereies";
import { createRawLoads, pickRawLoads } from "./util/chart/TokenizedLoad";
import type { RawLoadOfSql } from "./util/chart/TokenizedLoad";

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

export type DbTokenizedLoadChartProp = {
  dbName: string;
  sql: string;
  loadMaxPercent: number;
  tokenizedId: string;
  date: Date;
};

type Series = GrChartSeries<DbTokenizedLoadChartProp>;

export class DbTokenizedLoadChart extends GrChart<DbTokenizedLoadChartProp> {
  private readonly rawLoads: Record<string, RawLoadOfSql[]>;

  constructor(
    elm: HTMLElement,
    private readonly tokenizedTimeDbLoads: TimeDbLoad[],
    rawTimeDbLoads: TimeDbLoad[],
    private readonly threshold: number
  ) {
    super(elm);

    this.rawLoads = createRawLoads(rawTimeDbLoads);
  }

  render(databases: string[]) {
    const option = JSON.parse(JSON.stringify(DEFAULT_CHART_OPTION));
    option.series = this.createDbLoadChartSeries(databases);

    this.renderChart(option);
  }

  changeDatabase(databases: string[]) {
    this.render(databases);
  }

  getRawLoads(size: number, tokenizedId: string): RawLoadOfSql[] {
    return pickRawLoads(this.rawLoads, size, tokenizedId);
  }

  protected assignFunctions(option: GrChartOptions<DbTokenizedLoadChartProp>) {
    if ("labels" in option.yaxis) {
      option.yaxis.labels.formatter = this.yaxisLabelsFormatter;
    }
  }

  private createDbLoadChartSeries = (targetDatabases: string[]): Series => {
    return createLoadPropChartSeries<DbTokenizedLoadChartProp>(
      this.tokenizedTimeDbLoads,
      targetDatabases,
      this.threshold,
      (timeDbLoad, dbLoad, loadOfSql) => ({
        x: timeDbLoad.date,
        y: loadOfSql.loadMax * 100,
        prop: {
          date: timeDbLoad.date,
          loadMaxPercent: loadOfSql.loadMax * 100,
          dbName: dbLoad.name,
          tokenizedId: loadOfSql.tokenizedId,
          sql: loadOfSql.sql,
        },
      })
    );
  };

  private yaxisLabelsFormatter = (value: number): string => {
    return value.toFixed(0);
  };
}
