import { ExaminationIndexResult } from "./ExaminationIndexResult";
import { Type } from "class-transformer";

export class ExaminationResult {
  originalTimeMillis: number;

  @Type(() => ExaminationIndexResult)
  indexResults: ExaminationIndexResult[];
}
