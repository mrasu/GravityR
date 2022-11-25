import type { ITraceNode } from "@/models/trace_diagram/ITraceNode";

export class TraceTree {
  constructor(
    public traceId: string,
    public durationMilli: number,
    public sameServiceAccessCount: number,
    public timeConsumingServiceName: string,
    public root: ITraceNode
  ) {}

  toMermaidDiagram(): string {
    if (this.root.children.length === 0) {
      return this.toRootOnlyMermaidDiagram();
    }

    return "sequenceDiagram\r\n" + this.toMermaidDiagramRecursive(this.root);
  }

  private toRootOnlyMermaidDiagram(): string {
    return `
      sequenceDiagram
      ${this.root.serviceName}->>${this.root.serviceName}: ${" "}
    `;
  }

  private toMermaidDiagramRecursive(node: ITraceNode): string {
    let text = "";
    const serviceName = node.serviceName;

    for (let i = 0; i < node.children.length; i++) {
      const childTrace = node.children[i];
      text += childTrace.toMermaidStartArrow(serviceName) + "\n";
      if (childTrace.hasChild()) {
        text += `activate ${childTrace.serviceName}\n`;
      }

      text += this.toMermaidDiagramRecursive(childTrace);

      if (childTrace.hasChild()) {
        // Not show return arrow when the child trace is at the tail to make diagram smaller.
        text += childTrace.toMermaidEndArrow(serviceName) + "\n";

        text += `deactivate ${childTrace.serviceName}\n`;
      }
    }
    return text;
  }
}
