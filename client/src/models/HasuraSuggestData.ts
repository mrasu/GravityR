import type {
  IIndexTarget,
  IExaminationResult,
  IHasuraSuggestData,
} from "@/types/gr_param";
import { BaseSuggestData } from "./BaseSuggestData";
import { plainToInstance } from "class-transformer";
import type { IndexTarget } from "@/models/IndexTarget";
import { PostgresIndexTarget } from "@/models/PostgresIndexTarget";
import { ExaminationResult } from "@/models/ExaminationResult";
import { format } from "sql-formatter";
import { PostgresExplainData } from "@/models/explain_data/PostgresExplainData";

export class HasuraSuggestData extends BaseSuggestData {
  gql: string;
  gqlVariables: Record<string, any>;
  analyzeNodes?: PostgresExplainData[];
  summaryText: string;

  constructor(hasuraData: IHasuraSuggestData) {
    const suggestData = hasuraData.postgres;
    super(suggestData);

    this.query = format(suggestData.query, {
      language: "postgresql",
    });

    this.summaryText = suggestData.summaryText;

    this.gql = suggestData.gql;
    this.gqlVariables = suggestData.gqlVariables;

    this.analyzeNodes = plainToInstance(
      PostgresExplainData,
      suggestData.analyzeNodes
    );
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
