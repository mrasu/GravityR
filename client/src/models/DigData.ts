import { Type } from "class-transformer";
import { TimeDbLoad } from "./TimeDbLoad";

export class DigData {
  @Type(() => TimeDbLoad)
  sqlDbLoads: TimeDbLoad[];

  @Type(() => TimeDbLoad)
  tokenizedSqlDbLoads: TimeDbLoad[];
}
