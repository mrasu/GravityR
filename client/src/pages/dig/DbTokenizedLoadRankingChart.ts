import type { ApexOptions } from "apexcharts";
import type { TimeDbLoad } from "@/models/TimeDbLoad";
import { GrChart } from "@/lib/GrChart";
import type { GrChartOptions } from "@/lib/GrChart";
import type { GrChartSeries } from "@/lib/GrChart";
import { createRawLoads, pickRawLoads } from "./util/chart/TokenizedLoad";
import type { RawLoadOfSql } from "./util/chart/TokenizedLoad";

const DEFAULT_CHART_OPTION: ApexOptions = {
  chart: {
    height: 350,
    type: "bar",
  },
  plotOptions: {
    bar: {
      columnWidth: "45%",
      distributed: true,
    },
  },
  tooltip: {
    enabled: true,
  },
  dataLabels: {
    enabled: false,
  },
  legend: {
    show: false,
  },
  xaxis: {
    labels: {
      show: false,
    },
  },
  yaxis: {
    labels: {
      show: true,
    },
    title: {
      text: "Total Load of database (%)",
      style: {
        fontSize: "15px",
      },
    },
  },
};

export type DbTokenizedLoadRankingChartProp = {
  dbName: string;
  sql: string;
  loadSumPercent: number;
  tokenizedId: string;
  date: Date;
};

type Series = GrChartSeries<DbTokenizedLoadRankingChartProp>;

export class DbTokenizedLoadRankingChart extends GrChart<DbTokenizedLoadRankingChartProp> {
  private readonly rawLoads: Record<string, RawLoadOfSql[]>;

  constructor(
    elm: HTMLElement,
    private readonly tokenizedTimeDbLoads: TimeDbLoad[],
    rawTimeDbLoads: TimeDbLoad[],
    private readonly dataCount: number
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

  protected assignFunctions(
    option: GrChartOptions<DbTokenizedLoadRankingChartProp>
  ) {
    if ("labels" in option.yaxis) {
      option.yaxis.labels.formatter = this.yaxisLabelsFormatter;
    }
  }

  private createDbLoadChartSeries = (targetDatabases: string[]): Series => {
    let props: Record<
      string,
      Record<string, DbTokenizedLoadRankingChartProp>
    > = {};
    for (const load of this.tokenizedTimeDbLoads) {
      for (const dbLoad of load.databases) {
        if (!targetDatabases.includes(dbLoad.name)) continue;

        if (!props[dbLoad.name]) {
          props[dbLoad.name] = {};
        }
        const dbProps = props[dbLoad.name];
        for (const loadOfSql of dbLoad.sqls) {
          if (dbProps[loadOfSql.sql]) {
            dbProps[loadOfSql.sql].loadSumPercent += loadOfSql.loadSum * 100;
          } else {
            dbProps[loadOfSql.sql] = {
              dbName: dbLoad.name,
              sql: loadOfSql.sql,
              loadSumPercent: loadOfSql.loadSum * 100,
              tokenizedId: loadOfSql.tokenizedId,
              date: load.date,
            };
          }
        }
      }
    }

    const allProps = Object.keys(props)
      .map((dbName) =>
        Object.keys(props[dbName]).map((sql) => props[dbName][sql])
      )
      .flat();

    const elements = allProps
      .sort((p1, p2) => p2.loadSumPercent - p1.loadSumPercent)
      .slice(0, this.dataCount)
      .map((prop) => ({
        x: prop.sql,
        y: prop.loadSumPercent,
        prop: prop,
      }));

    return [
      {
        data: elements,
      },
    ];
  };

  private yaxisLabelsFormatter = (value: number): string => {
    return value.toFixed(0);
  };
}
