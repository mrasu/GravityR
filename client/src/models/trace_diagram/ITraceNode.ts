export abstract class ITraceNode {
  protected constructor(
    public serviceName: string,
    public children: ITraceNode[]
  ) {}

  abstract toMermaidStartArrow(fromServiceName: string): string;
  abstract toMermaidEndArrow(fromServiceName: string): string;

  hasChild(): boolean {
    return this.children.length > 0;
  }
}
