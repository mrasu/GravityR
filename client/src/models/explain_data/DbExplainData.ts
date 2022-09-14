export abstract class DbExplainData {
  text: string;
  title: string;
  tableName?: string;

  abstract children?: this[];

  abstract calculateEndTime(startTime: number): number;
}
