import { Type } from "class-transformer";
import { DbLoad } from "./DbLoad";

export class TimeDbLoad {
  timestamp: number;

  @Type(() => DbLoad)
  databases: DbLoad[];

  get date(): Date {
    return new Date(this.timestamp);
  }
}
