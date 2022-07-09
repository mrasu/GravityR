import { IndexTarget } from "./IndexTarget";
import { Type } from "class-transformer";

export class ExaminationIndexResult {
  @Type(() => IndexTarget)
  indexTarget: IndexTarget;
  executionTimeMillis: number;

  toEfficiency(originalTimeMillis: number): string {
    return (
      (1 - this.executionTimeMillis / originalTimeMillis) *
      100
    ).toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  }

  toIndex(): string {
    return this.indexTarget.toString();
  }
}
