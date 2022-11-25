type LoadProp = {
  dbName: string;
  sql: string;
  loadMaxPercent: number;
  date: Date;
};

export const createLoadMaxTooltip = (
  { date, dbName, loadMaxPercent, sql }: LoadProp,
  maxHeight: number
): string => {
  return `
    <div style='max-width: 500px; white-space: normal; max-height: ${maxHeight}px'>
      ${date.toUTCString()}<br />
      Database: ${dbName}<br />
      <div style='margin-bottom: 10px;'>
        db.load.avg:&nbsp;
        <span style='font-weight: bold;'>
          ${loadMaxPercent.toFixed(2)}%
        </span>
      </div>
      <span style='white-space: break-spaces'>${sql}</span>
    </div>
    `;
};

type LoadSumProp = {
  dbName: string;
  sql: string;
  loadSumPercent: number;
};

export const createLoadSumTooltip = (
  { dbName, loadSumPercent, sql }: LoadSumProp,
  maxHeight: number
): string => {
  return `
    <div style='max-width: 500px; white-space: normal; max-height: ${maxHeight}px'>
      Database: ${dbName}<br />
      <div style='margin-bottom: 10px;'>
        Sum of db.load.avg:&nbsp;
        <span style='font-weight: bold;'>
          ${loadSumPercent.toFixed(2)}%
        </span>
      </div>
      <span style='white-space: break-spaces'>${sql}</span>
    </div>
    `;
};
