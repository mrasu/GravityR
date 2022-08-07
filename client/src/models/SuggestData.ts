import type { IAnalyzeData, ISuggestData } from "../types/gr_param";
import { IndexTarget } from "./IndexTarget";
import { ExaminationCommandOption } from "./ExaminationCommandOption";
import { ExaminationResult } from "./ExaminationResult";
import { plainToInstance } from "class-transformer";

export class SuggestData {
  analyzeNodes?: IAnalyzeData[];
  query: string;
  indexTargets?: IndexTarget[];
  examinationCommandOptions: ExaminationCommandOption[];
  examinationResult?: ExaminationResult;

  constructor(iSuggestData: ISuggestData) {
    this.analyzeNodes = iSuggestData.analyzeNodes;
    this.query = iSuggestData.query;

    this.indexTargets = iSuggestData.indexTargets?.map((v) =>
      plainToInstance(IndexTarget, v)
    );

    this.examinationCommandOptions =
      iSuggestData.examinationCommandOptions?.map((v) =>
        plainToInstance(ExaminationCommandOption, v)
      );

    this.examinationResult = iSuggestData.examinationResult
      ? plainToInstance(ExaminationResult, iSuggestData.examinationResult)
      : undefined;
  }
}
