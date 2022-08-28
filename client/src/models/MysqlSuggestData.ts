import type {
  IMysqlAnalyzeData,
  IMysqlSuggestData,
  IIndexTarget,
  IExaminationResult,
} from "@/types/gr_param";
import { BaseSuggestData } from "./BaseSuggestData";
import { plainToInstance } from "class-transformer";
import { IndexTarget } from "@/models/IndexTarget";
import { ExaminationResult } from "@/models/ExaminationResult";

export class MysqlSuggestData extends BaseSuggestData {
  analyzeNodes?: IMysqlAnalyzeData[];

  constructor(suggestData: IMysqlSuggestData) {
    super(suggestData);

    this.analyzeNodes = suggestData.analyzeNodes;

    this.indexTargets = suggestData.indexTargets?.map((v) =>
      plainToInstance(IndexTarget, v)
    );

    this.examinationResult = suggestData.examinationResult
      ? plainToInstance(ExaminationResult, suggestData.examinationResult)
      : undefined;

    for (const result of this.examinationResult.indexResults) {
      result.indexTarget = plainToInstance(IndexTarget, result.indexTarget);
    }
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
