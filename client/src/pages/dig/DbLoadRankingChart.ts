import type { ApexOptions } from "apexcharts";
import type { TimeDbLoad } from "@/models/TimeDbLoad";
import { GrChart } from "@/lib/GrChart";
import type { GrChartOptions } from "@/lib/GrChart";
import type { GrChartSeries } from "@/lib/GrChart";

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

export type DbLoadRankingChartProp = {
  dbName: string;
  sql: string;
  loadSumPercent: number;
};

type Series = GrChartSeries<DbLoadRankingChartProp>;

export class DbLoadRankingChart extends GrChart<DbLoadRankingChartProp> {
  constructor(
    elm: HTMLElement,
    private readonly timeDbLoads: TimeDbLoad[],
    private readonly dataCount: number
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

  protected assignFunctions(option: GrChartOptions<DbLoadRankingChartProp>) {
    if ("labels" in option.yaxis) {
      option.yaxis.labels.formatter = this.yaxisLabelsFormatter;
    }
  }

  private createDbLoadChartSeries = (targetDatabases: string[]): Series => {
    let props: Record<string, Record<string, DbLoadRankingChartProp>> = {};
    for (const load of this.timeDbLoads) {
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
