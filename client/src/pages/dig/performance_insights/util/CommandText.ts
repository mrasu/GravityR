export const buildCommandText = (dbName: string, sql: string): string => {
  const db = dbName ?? "...";
  const query = sql ?? "...";

  return `DB_DATABASE="${db}" gr db suggest -o "suggest.html" \\\n\t-q "${query}"`;
};
