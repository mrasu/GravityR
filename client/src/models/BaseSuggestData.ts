import type { IndexTarget } from "./IndexTarget";
import { ExaminationCommandOption } from "./ExaminationCommandOption";
import type { ExaminationResult } from "./ExaminationResult";
import type {
  IDbSuggestData,
  IExaminationResult,
  IIndexTarget,
} from "@/types/gr_param";
import { plainToInstance } from "class-transformer";

export abstract class BaseSuggestData {
  query: string;
  indexTargets?: IndexTarget[];
  examinationCommandOptions: ExaminationCommandOption[];
  examinationResult?: ExaminationResult;

  protected constructor(suggestData: IDbSuggestData) {
    this.query = suggestData.query;
    this.indexTargets = this.createIndexTargets(suggestData.indexTargets);

    this.examinationCommandOptions = suggestData.examinationCommandOptions?.map(
      (v) => plainToInstance(ExaminationCommandOption, v)
    );

    this.examinationResult = this.createExaminationResult(
      suggestData.examinationResult
    );
  }

  protected abstract createIndexTargets(
    targets?: IIndexTarget[]
  ): IndexTarget[];

  protected abstract createExaminationResult(
    result?: IExaminationResult
  ): ExaminationResult;
}
