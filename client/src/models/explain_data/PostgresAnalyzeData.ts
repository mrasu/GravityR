import { Type } from "class-transformer";
import { DbAnalyzeData } from "@/models/explain_data/DbAnalyzeData";

export class PostgresAnalyzeData extends DbAnalyzeData {
  @Type(() => PostgresAnalyzeData)
  children?: this[];

  estimatedInitCost: number;
  estimatedCost: number;
  estimatedReturnedRows: number;
  estimatedWidth: number;

  calculateEndTime(_startTime: number): number {
    //  actualTimeAvg seems not the average time of loops for PostgreSQL even doc says so
    return this.actualTimeAvg;
  }
}
