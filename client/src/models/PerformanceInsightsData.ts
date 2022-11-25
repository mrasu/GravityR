import { Type } from "class-transformer";
import { TimeDbLoad } from "@/models/TimeDbLoad";

export class PerformanceInsightsData {
  @Type(() => TimeDbLoad)
  sqlDbLoads: TimeDbLoad[];

  @Type(() => TimeDbLoad)
  tokenizedSqlDbLoads: TimeDbLoad[];
}
