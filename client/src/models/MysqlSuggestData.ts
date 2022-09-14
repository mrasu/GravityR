import type {
  IMysqlSuggestData,
  IIndexTarget,
  IExaminationResult,
} from "@/types/gr_param";
import { BaseSuggestData } from "./BaseSuggestData";
import { plainToInstance } from "class-transformer";
import { IndexTarget } from "@/models/IndexTarget";
import { ExaminationResult } from "@/models/ExaminationResult";
import { MysqlAnalyzeData } from "@/models/explain_data/MysqlAnalyzeData";

export class MysqlSuggestData extends BaseSuggestData {
  analyzeNodes?: MysqlAnalyzeData[];

  constructor(suggestData: IMysqlSuggestData) {
    super(suggestData);

    this.analyzeNodes = plainToInstance(
      MysqlAnalyzeData,
      suggestData.analyzeNodes
    );
  }

  protected createIndexTargets(targets?: IIndexTarget[]): IndexTarget[] {
    return targets?.map((v) => plainToInstance(IndexTarget, v));
  }

  protected createExaminationResult(
    result?: IExaminationResult
  ): ExaminationResult {
    return result ? plainToInstance(ExaminationResult, result) : undefined;
  }
}
