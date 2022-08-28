import type {
  IIndexTarget,
  IPostgresAnalyzeData,
  IPostgresSuggestData,
  IExaminationResult,
} from "@/types/gr_param";
import { BaseSuggestData } from "./BaseSuggestData";
import { plainToInstance } from "class-transformer";
import type { IndexTarget } from "@/models/IndexTarget";
import { PostgresIndexTarget } from "@/models/PostgresIndexTarget";
import { ExaminationResult } from "@/models/ExaminationResult";

export class PostgresSuggestData extends BaseSuggestData {
  analyzeNodes?: IPostgresAnalyzeData[];
  planningText: string;

  constructor(suggestData: IPostgresSuggestData) {
    super(suggestData);

    this.analyzeNodes = suggestData.analyzeNodes;
    this.planningText = suggestData.planningText;
  }

  protected createIndexTargets(targets?: IIndexTarget[] | null): IndexTarget[] {
    return targets?.map((v) => plainToInstance(PostgresIndexTarget, v));
  }

  protected createExaminationResult(
    result?: IExaminationResult
  ): ExaminationResult {
    const res = result ? plainToInstance(ExaminationResult, result) : undefined;
    if (!res) return res;

    for (const result of res.indexResults) {
      result.indexTarget = plainToInstance(
        PostgresIndexTarget,
        result.indexTarget
      );
    }

    return res;
  }
}
