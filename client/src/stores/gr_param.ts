import { readable } from "svelte/store";
import { SuggestData } from "../models/SuggestData";
import { DigData } from "../models/DigData";
import { plainToInstance } from "class-transformer";

export const grParam = readable({
  dev: window.grParam.dev || false,
  suggestData: window.grParam.suggestData
    ? new SuggestData(window.grParam.suggestData)
    : null,

  digData: window.grParam.digData
    ? plainToInstance(DigData, window.grParam.digData)
    : null,
});
