import { DbExplainData } from "@/models/explain_data/DbExplainData";

export abstract class DbAnalyzeData extends DbExplainData {
  actualTimeFirstRow?: number;
  actualTimeAvg?: number;
  actualReturnedRows?: number;
  actualLoopCount?: number;
}
