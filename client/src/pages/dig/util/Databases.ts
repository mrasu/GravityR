import type { TimeDbLoad } from "../../../models/TimeDbLoad";

export const getAllDatabase = (dbLoads?: TimeDbLoad[]): string[] => {
  if (!dbLoads) return [];

  const dbs = dbLoads.map((c) => c.databases.map((d) => d.name)).flat();
  return Array.from(new Set(dbs)).sort();
};

export const getFirstDatabaseAsArray = (databases: string[]): string[] => {
  if (databases.length === 0) {
    return [];
  }

  return [databases[0]];
};
