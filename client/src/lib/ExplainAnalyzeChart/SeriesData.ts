import type { IAnalyzeData } from "../../types/gr_param";

export interface SeriesData {
  x: string;
  y: number[];
  title?: string;
  text?: string;
  goals?: {
    value: number;
    strokeColor: string;
  }[];
  IAnalyzeData: IAnalyzeData;
}
