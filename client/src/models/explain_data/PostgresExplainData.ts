import { DbExplainData } from "@/models/explain_data/DbExplainData";
import { Type } from "class-transformer";

export class PostgresExplainData extends DbExplainData {
  @Type(() => PostgresExplainData)
  children?: this[];

  estimatedInitCost: number;
  estimatedCost: number;
  estimatedReturnedRows: number;
  estimatedWidth: number;

  calculateEndTime(_startTime: number): number {
    return this.estimatedCost;
  }
}
