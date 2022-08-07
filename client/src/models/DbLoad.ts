import { Type } from "class-transformer";
import { DbLoadOfSql } from "./DbLoadOfSql";

export class DbLoad {
  name: string;

  @Type(() => DbLoadOfSql)
  sqls: DbLoadOfSql[];
}
