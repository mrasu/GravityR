import type { TimeDbLoad } from "../../../../models/TimeDbLoad";
import type { DbLoadOfSql } from "../../../../models/DbLoadOfSql";

export type RawLoadOfSql = {
  dbName: string;
  load: DbLoadOfSql;
  date: Date;
};

export const createRawLoads = (
  tokenizedLoads: TimeDbLoad[]
): Record<string, RawLoadOfSql[]> => {
  const rawLoads = {};

  for (const load of tokenizedLoads) {
    for (const dbLoad of load.databases) {
      for (const loadOfSql of dbLoad.sqls) {
        if (!rawLoads[loadOfSql.tokenizedId]) {
          rawLoads[loadOfSql.tokenizedId] = [];
        }
        rawLoads[loadOfSql.tokenizedId].push({
          dbName: dbLoad.name,
          load: loadOfSql,
          date: load.date,
        });
      }
    }
  }

  return rawLoads;
};

export const pickRawLoads = (
  rawLoads: Record<string, RawLoadOfSql[]>,
  size: number,
  tokenizedId: string
): RawLoadOfSql[] => {
  const loads = rawLoads[tokenizedId];
  const sortedLoads = [...loads].sort(
    (l1, l2) => l2.load.loadMax - l1.load.loadMax
  );

  return sortedLoads.slice(0, size);
};
