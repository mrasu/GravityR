import type { ISuggestData } from "@/types/gr_param";
import { MysqlSuggestData } from "./MysqlSuggestData";
import { PostgresSuggestData } from "./PostgresSuggestData";

export class SuggestData {
  mysqlSuggestData?: MysqlSuggestData;
  postgresSuggestData?: PostgresSuggestData;

  constructor(iSuggestData: ISuggestData) {
    if (iSuggestData.mysql) {
      this.mysqlSuggestData = new MysqlSuggestData(iSuggestData.mysql);
    }
    if (iSuggestData.postgres) {
      this.postgresSuggestData = new PostgresSuggestData(iSuggestData.postgres);
      console.log(this.postgresSuggestData);
    }
  }
}
