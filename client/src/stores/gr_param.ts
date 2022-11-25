import { readable } from "svelte/store";
import { SuggestData } from "@/models/SuggestData";
import { DigData } from "@/models/DigData";

export const grParam = readable({
  dev: window.grParam.dev || false,
  suggestData: window.grParam.suggestData
    ? new SuggestData(window.grParam.suggestData)
    : null,

  digData: window.grParam.digData ? new DigData(window.grParam.digData) : null,
});
