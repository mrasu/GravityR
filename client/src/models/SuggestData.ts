import type { ISuggestData } from "@/types/gr_param";
import { MysqlSuggestData } from "./MysqlSuggestData";
import { PostgresSuggestData } from "./PostgresSuggestData";
import { HasuraSuggestData } from "@/models/HasuraSuggestData";

export class SuggestData {
  mysqlSuggestData?: MysqlSuggestData;
  postgresSuggestData?: PostgresSuggestData;
  hasuraSuggestData?: HasuraSuggestData;

  constructor(iSuggestData: ISuggestData) {
    if (iSuggestData.mysql) {
      this.mysqlSuggestData = new MysqlSuggestData(iSuggestData.mysql);
    }

    if (iSuggestData.postgres) {
      this.postgresSuggestData = new PostgresSuggestData(iSuggestData.postgres);
    }

    if (iSuggestData.hasura) {
      this.hasuraSuggestData = new HasuraSuggestData(iSuggestData.hasura);
    }
  }
}
