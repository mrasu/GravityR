import { plainToInstance, Type } from "class-transformer";
import { PerformanceInsightsData } from "@/models/PerformanceInsightsData";
import { JaegerData } from "@/models/JaegerData";
import type { IDigData } from "@/types/gr_param";

export class DigData {
  @Type(() => PerformanceInsightsData)
  performanceInsightsData?: PerformanceInsightsData;

  @Type(() => JaegerData)
  jaegerData?: JaegerData;

  constructor(iDigData: IDigData) {
    if (iDigData.performanceInsights) {
      this.performanceInsightsData = plainToInstance(
        PerformanceInsightsData,
        iDigData.performanceInsights
      );
    }

    if (iDigData.jaeger) {
      this.jaegerData = plainToInstance(JaegerData, iDigData.jaeger);
    }
  }
}
