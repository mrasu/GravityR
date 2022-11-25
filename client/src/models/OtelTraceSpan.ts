import { Type } from "class-transformer";
import { TraceNode } from "@/models/trace_diagram/TraceNode";
import { TraceRepeatNode } from "@/models/trace_diagram/TraceRepeatNode";

export class OtelTraceSpan {
  name: string;
  spanId: string;
  startTimeMillis: number;
  endTimeMillis: number;
  serviceName: string;

  @Type(() => OtelTraceSpan)
  children: OtelTraceSpan[];

  toTraceNode(): TraceNode {
    let children: TraceNode[] = [];
    let i = 0;
    for (; i < this.children.length; i++) {
      const child = this.children[i];
      if (child.isLast) {
        const repeatCount = this.countRepeatedChildService(child, i);
        children = children.concat(this.createTraceNodes(repeatCount, child));

        i = i + repeatCount - 1;
      } else {
        children.push(child.toTraceNode());
      }
    }

    return new TraceNode(this.serviceName, children);
  }

  private countRepeatedChildService(
    child: OtelTraceSpan,
    currentIndex: number
  ): number {
    if (currentIndex === this.children.length - 1) return 1;

    let cnt = 1;
    for (let i = currentIndex + 1; i < this.children.length; i++) {
      const otherChild = this.children[i];
      if (otherChild.isLast === false) break;
      if (child.isSameService(otherChild) === false) break;

      cnt++;
    }

    return cnt;
  }

  private createTraceNodes(
    repeatCount: number,
    span: OtelTraceSpan
  ): TraceNode[] {
    let children: TraceNode[] = [];
    const loopCount = Math.min(4, repeatCount);
    for (let j = 0; j < loopCount; j++) {
      if (j < 3) {
        children.push(span.toTraceNode());
        continue;
      }

      if (j + 1 === repeatCount) {
        children.push(span.toTraceNode());
      } else {
        children.push(new TraceRepeatNode(repeatCount - j, span.serviceName));
      }
    }

    return children;
  }

  isSameService(other: OtelTraceSpan): boolean {
    return other.serviceName === this.serviceName;
  }

  get isLast(): boolean {
    return this.children.length === 0;
  }
}
