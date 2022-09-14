import { DbAnalyzeData } from "@/models/explain_data/DbAnalyzeData";
import { Type } from "class-transformer";
import "reflect-metadata";

export class MysqlAnalyzeData extends DbAnalyzeData {
  @Type(() => MysqlAnalyzeData)
  children?: this[];

  estimatedInitCost?: number;
  estimatedCost?: number;
  estimatedReturnedRows?: number;

  calculateEndTime(startTime: number): number {
    if (this.actualLoopCount > 1) {
      if (this.actualLoopCount > 1000 && this.actualTimeAvg === 0) {
        // ActualTimeAvg doesn't show the value less than 0.000, however,
        // when the number of loop is much bigger, time can be meaningful even the result of multiplication is zero.
        // To handle the problem, assume ActualTimeAvg is some less than 0.000
        return startTime + 0.0001 * this.actualLoopCount;
      } else {
        return startTime + this.actualTimeAvg * this.actualLoopCount;
      }
    } else {
      return this.actualTimeAvg;
    }
  }
}
