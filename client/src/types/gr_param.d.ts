export interface IGrParam {
  dev?: boolean;
  suggestData?: ISuggestData;
  digData?: IDigData;
}

export interface ISuggestData {
  analyzeNodes?: IAnalyzeData[];
  query: string;
  indexTargets?: IIndexTarget[];
  examinationCommandOptions: IExaminationCommandOption[];
  examinationResult?: IExaminationResult;
}

export interface IDigData {
  sqlDbLoads: ITimeDbLoad[];
  tokenizedSqlDbLoads: ITimeDbLoad[];
}

export interface IAnalyzeData {
  text: string;
  title: string;
  tableName?: string;
  estimatedInitCost?: number;
  estimatedCost?: number;
  estimatedReturnedRows?: number;
  actualTimeFirstRow?: number;
  actualTimeAvg?: number;
  actualReturnedRows?: number;
  actualLoopCount?: number;

  children?: IAnalyzeData[];
}

interface IIndexTarget {
  tableName: string;
  columns: IIndexColumn[];
}

interface IIndexColumn {
  name: string;
}

interface IExaminationCommandOption {
  isShort: boolean;
  name: string;
  value: string;
}

interface IExaminationResult {
  originalTimeMillis: number;
  indexResults: IExaminationIndexResult[];
}

interface IExaminationIndexResult {
  indexTarget: IIndexTarget;
  executionTimeMillis: number;
}

export interface ITimeDbLoad {
  timestamp: number;
  databases: IDbLoad[];
}

export interface IDbLoad {
  name: string;
  sqls: IDbLoadOfSql[];
}

export interface IDbLoadOfSql {
  sql: string;
  loadMax: number;
  loadSum: number;
  tokenizedId: string;
}
