import type { IDbLoad, ITimeDbLoad } from "../types/gr_param";
import digSqlDbLoads from "./digSqlDbLoads.json";
import digTokenizedSqlDbLoads from "./digTokenizedSqlDbLoads.json";
import jaegerTraces from "./jaegerTraces.json";

type timeDbLoad = { timestamp: string; databases: IDbLoad[] };

const toTimeDbLoads = (loads: timeDbLoad[]) => {
  const timeDbLoads: ITimeDbLoad[] = [];
  for (const sql of loads) {
    timeDbLoads.push({
      ...sql,
      timestamp: new Date(sql.timestamp).getTime(),
    });
  }
  return timeDbLoads;
};

export const dummyDigData = {
  performanceInsights: {
    sqlDbLoads: toTimeDbLoads(digSqlDbLoads),
    tokenizedSqlDbLoads: toTimeDbLoads(digTokenizedSqlDbLoads),
  },
  jaeger: jaegerTraces,
};
