import { Type } from "class-transformer";
import { OtelCompactedTrace } from "@/models/OtelCompactedTrace";

export class JaegerData {
  uiPath: string;
  slowThresholdMilli: number;
  sameServiceThreshold: number;

  @Type(() => OtelCompactedTrace)
  slowTraces: OtelCompactedTrace[];

  @Type(() => OtelCompactedTrace)
  sameServiceTraces: OtelCompactedTrace[];
}
