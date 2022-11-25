import { ITraceNode } from "@/models/trace_diagram/ITraceNode";

const SERVICE_NAME = "repeat";

export class TraceRepeatNode extends ITraceNode {
  constructor(
    public repeat: number,
    public toService: string,
    children: ITraceNode[] = []
  ) {
    super(SERVICE_NAME, children);
  }

  toMermaidStartArrow(fromServiceName: string): string {
    let texts = [];
    const repeatCount = Math.min(this.repeat, 3);
    for (let j = 0; j < repeatCount; j++) {
      if (j === 0) {
        const text = `${fromServiceName}->>${this.toService}: More ${this.repeat} times...`;
        texts.push(text);
      } else {
        texts.push(`${fromServiceName}->>${this.toService}: `);
      }
    }

    return texts.join("\n");
  }

  toMermaidEndArrow(_: string): string {
    return "";
  }
}
