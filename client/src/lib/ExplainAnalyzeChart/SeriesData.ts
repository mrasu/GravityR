import type { IAnalyzeData } from "../../types/gr_param";

export class SeriesData {
  x: string;
  y: number[];
  title?: string;
  text?: string;
  goals?: {
    value: number;
    strokeColor: string;
  }[];
  IAnalyzeData: IAnalyzeData;

  constructor({
    x,
    y,
    IAnalyzeData,
    title,
    text,
    goals,
  }: {
    x: string;
    y: number[];
    title?: string;
    text?: string;
    goals?: {
      value: number;
      strokeColor: string;
    }[];
    IAnalyzeData: IAnalyzeData;
  }) {
    this.x = x;
    this.y = y;
    this.IAnalyzeData = IAnalyzeData;
    this.title = title;
    this.text = text;
    this.goals = goals;
  }

  copyWith(y: [number, number], title?: string) {
    return new SeriesData({
      ...this,
      y,
      title,
    });
  }

  get xNum(): number {
    return Number(this.x);
  }
}
