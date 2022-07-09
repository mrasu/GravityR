export interface IGrParam {
  analyzeNodes?: IAnalyzeData[];
  query: string;
  indexTargets?: IIndexTarget[];
  examinationCommandOptions: IExaminationCommandOption[];
  examinationResult?: IExaminationResult;
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
