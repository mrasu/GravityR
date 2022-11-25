import { Type } from "class-transformer";
import { OtelTraceSpan } from "@/models/OtelTraceSpan";
import { TraceTree } from "@/models/trace_diagram/TraceTree";

export class OtelCompactedTrace {
  traceId: string;
  sameServiceAccessCount: number;
  timeConsumingServiceName: string;

  @Type(() => OtelTraceSpan)
  root: OtelTraceSpan;

  toTraceTree(): TraceTree {
    const duration = this.root.endTimeMillis - this.root.startTimeMillis;
    const root = this.root.toTraceNode();
    return new TraceTree(
      this.traceId,
      duration,
      this.sameServiceAccessCount,
      this.timeConsumingServiceName,
      root
    );
  }
}
