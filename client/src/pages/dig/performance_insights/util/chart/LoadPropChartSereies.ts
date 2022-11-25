import type { GrChartSeries, GrChartSeriesData } from "@/lib/GrChart";
import type { TimeDbLoad } from "@/models/TimeDbLoad";
import type { DbLoadOfSql } from "@/models/DbLoadOfSql";
import type { DbLoad } from "@/models/DbLoad";

type SeriesRecord<T> = Record<string, Record<string, GrChartSeries<T>[number]>>;

export const createLoadPropChartSeries = <T>(
  dbLoads: TimeDbLoad[],
  targetDatabases: string[],
  threshold: number,
  buildDataFn: (
    timeDbLoad: TimeDbLoad,
    dbLoad: DbLoad,
    loadOfSql: DbLoadOfSql
  ) => GrChartSeriesData<T>
): GrChartSeries<T> => {
  let series: GrChartSeries<T> = [];
  let knownSeries: SeriesRecord<T> = {};

  for (let i = 0; i < dbLoads.length; i++) {
    const timeDbLoad = dbLoads[i];
    for (const dbLoad of timeDbLoad.databases) {
      if (!targetDatabases.includes(dbLoad.name)) continue;

      for (const loadOfSql of dbLoad.sqls) {
        if (loadOfSql.loadMax < threshold) continue;

        if (!knownSeries[dbLoad.name]) {
          knownSeries[dbLoad.name] = {};
        }
        const dbSeries = knownSeries[dbLoad.name];
        if (!dbSeries[loadOfSql.sql]) {
          const seriesData = {
            database: dbLoad.name,
            name: loadOfSql.sql,
            data: Array(i).fill({
              x: timeDbLoad.date,
              y: 0,
            }),
          };
          dbSeries[loadOfSql.sql] = seriesData;
          series.push(seriesData);
        }

        dbSeries[loadOfSql.sql].data.push(
          buildDataFn(timeDbLoad, dbLoad, loadOfSql)
        );
      }
    }

    for (const s of series) {
      if (s.data.length < i + 1) {
        s.data.push({
          x: timeDbLoad.date,
          y: 0,
        });
      }
    }
  }

  return series;
};
