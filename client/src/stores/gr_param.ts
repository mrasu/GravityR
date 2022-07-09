import { readable } from "svelte/store";
import { ExaminationResult } from "../models/ExaminationResult";
import { plainToInstance } from "class-transformer";
import { IndexTarget } from "../models/IndexTarget";
import { ExaminationCommandOption } from "../models/ExaminationCommandOption";

export const grParam = readable({
  analyzeNodes: window.gr.analyzeNodes,
  query: window.gr.query,

  indexTargets: window.gr.indexTargets?.map((v) =>
    plainToInstance(IndexTarget, v)
  ),

  examinationCommandOptions: window.gr.examinationCommandOptions?.map((v) =>
    plainToInstance(ExaminationCommandOption, v)
  ),

  examinationResult: window.gr.examinationResult
    ? plainToInstance(ExaminationResult, window.gr.examinationResult)
    : undefined,
});
