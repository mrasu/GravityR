import { ITraceNode } from "@/models/trace_diagram/ITraceNode";

export class TraceNode extends ITraceNode {
  constructor(serviceName: string, children: ITraceNode[] = []) {
    super(serviceName, children);
  }

  toMermaidStartArrow(fromServiceName: string): string {
    return `${fromServiceName}->>${this.serviceName}: ${this.serviceName}`;
  }

  toMermaidEndArrow(fromServiceName: string): string {
    return `${this.serviceName}->>${fromServiceName}: ${fromServiceName}`;
  }
}
