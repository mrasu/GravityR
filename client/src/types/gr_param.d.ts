export interface IGrParam {
  dev?: boolean;
  suggestData?: ISuggestData;
  digData?: IDigData;
}

export interface ISuggestData {
  mysql?: IMysqlSuggestData;
  postgres?: IPostgresSuggestData;
  hasura?: IHasuraSuggestData;
}

export abstract interface IDbSuggestData {
  query: string;
  indexTargets?: IIndexTarget[];
  examinationCommandOptions: IExaminationCommandOption[];
  examinationResult?: IExaminationResult;
}

export interface IMysqlSuggestData extends IDbSuggestData {
  analyzeNodes?: IMysqlAnalyzeData[];
}

export interface IPostgresSuggestData extends IDbSuggestData {
  analyzeNodes?: IPostgresAnalyzeData[];
  summaryText: string;
}

export interface IHasuraSuggestData {
  postgres: IHasuraPostgresSuggestData;
}

export interface IHasuraPostgresSuggestData extends IDbSuggestData {
  gql: string;
  gqlVariables: Record<string, any>;
  analyzeNodes?: IPostgresExplainData[];
  summaryText: string;
}

export interface IDigData {
  sqlDbLoads: ITimeDbLoad[];
  tokenizedSqlDbLoads: ITimeDbLoad[];
}

export abstract interface IDbAnalyzeData {
  text: string;
  title: string;
  tableName?: string;
  actualTimeFirstRow?: number;
  actualTimeAvg?: number;
  actualReturnedRows?: number;
  actualLoopCount?: number;

  children?: IDbAnalyzeData[];
}

export interface IMysqlAnalyzeData extends IDbAnalyzeData {
  estimatedInitCost?: number;
  estimatedCost?: number;
  estimatedReturnedRows?: number;

  children?: IMysqlAnalyzeData[];
}

export interface IPostgresAnalyzeData extends IDbAnalyzeData {
  estimatedInitCost: number;
  estimatedCost: number;
  estimatedReturnedRows: number;
  estimatedWidth: number;

  children?: IPostgresAnalyzeData[];
}

export abstract interface IDbExplainData {
  text: string;
  title: string;
  tableName?: string;

  children?: IDbAnalyzeData[];
}

export interface IPostgresExplainData extends IDbExplainData {
  estimatedInitCost: number;
  estimatedCost: number;
  estimatedReturnedRows: number;
  estimatedWidth: number;

  children?: IPostgresExplainData[];
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
