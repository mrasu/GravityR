import type { IAnalyzeData } from "../../types/gr_param";
import { SeriesData } from "./SeriesData";

class XNumberTree {
  nodes: XNumberNode[];

  constructor() {
    this.nodes = [];
  }
}

class XNumberNode {
  children: XNumberNode[];

  constructor(public xNum: number, public parent: XNumberNode) {
    this.children = [];
  }
}

type XToSeriesRecords = Record<number, SeriesPointer>;
type SeriesPointer = { seriesData: SeriesData; treeNode: XNumberNode };

const createXTree = (
  IAnalyzeDatas: IAnalyzeData[]
): [XNumberTree, XToSeriesRecords] => {
  const xTree = new XNumberTree();
  const xToSeriesRecords: XToSeriesRecords = {};

  let nextXNumber = 0;
  for (let i = 0; i < IAnalyzeDatas.length; i++) {
    const xNode = new XNumberNode(nextXNumber, null);
    xTree.nodes.push(xNode);
    [nextXNumber] = visitAnalyze(
      xNode,
      xToSeriesRecords,
      nextXNumber,
      IAnalyzeDatas[i]
    );
  }

  return [xTree, xToSeriesRecords];
};

const STROKE_COLOR = "#CD2F2A";

const visitAnalyze = (
  parentNode: XNumberNode,
  xToSeriesRecords: XToSeriesRecords,
  nextXNumber: number,
  IAnalyzeData: IAnalyzeData
): [number, number] => {
  let maxFinished = 0;
  const seriesData = new SeriesData({
    x: nextXNumber.toString(),
    y: [0, 0],
    title: IAnalyzeData.title,
    text: IAnalyzeData.text,
    goals: [
      {
        value: 0,
        strokeColor: STROKE_COLOR,
      },
    ],
    IAnalyzeData: IAnalyzeData,
  });
  xToSeriesRecords[nextXNumber] = {
    seriesData: seriesData,
    treeNode: parentNode,
  };

  nextXNumber++;
  let startTime = 0;
  if (IAnalyzeData.children) {
    let childrenFinished = [];
    for (let data of IAnalyzeData.children) {
      let childFinished = 0;
      const xNode = new XNumberNode(nextXNumber, parentNode);
      parentNode.children.push(xNode);
      [nextXNumber, childFinished] = visitAnalyze(
        xNode,
        xToSeriesRecords,
        nextXNumber,
        data
      );
      childrenFinished.push(childFinished);
    }
    startTime = Math.max(...childrenFinished);
  }

  let endTimeFromData: number;
  if (IAnalyzeData.actualLoopCount > 1) {
    if (
      IAnalyzeData.actualLoopCount > 1000 &&
      IAnalyzeData.actualTimeAvg === 0
    ) {
      // ActualTimeAvg doesn't show the value less than 0.000, however,
      // when the number of loop is much bigger, time can be meaningful even the result of multiplication is zero.
      // To handle the problem, assume ActualTimeAvg is some less than 0.000
      endTimeFromData = startTime + 0.0001 * IAnalyzeData.actualLoopCount;
    } else {
      endTimeFromData =
        startTime + IAnalyzeData.actualTimeAvg * IAnalyzeData.actualLoopCount;
    }
  } else {
    endTimeFromData = IAnalyzeData.actualTimeAvg;
  }

  const endTime = Math.max(endTimeFromData, startTime);
  maxFinished = Math.max(maxFinished, endTime);
  seriesData.y = [Math.min(startTime, endTime - 0.0001), endTime];
  seriesData.goals[0].value = endTime;

  return [nextXNumber, maxFinished];
};

export class ExplainAnalyzeTree {
  xTree: XNumberTree;
  xToSeriesRecords: XToSeriesRecords;

  constructor(IAnalyzeDatas: IAnalyzeData[]) {
    [this.xTree, this.xToSeriesRecords] = createXTree(IAnalyzeDatas);
  }

  getFromX(xNum: number): SeriesData {
    return this.xToSeriesRecords[xNum].seriesData;
  }

  getNextXs(xNum: number): number[] {
    const node = this.findNode(xNum, this.xTree.nodes);
    if (!node) return null;
    if (node.children.length === 0) return null;

    return node.children.map((v) => v.xNum);
  }

  private findNode(xNum: number, nodes: XNumberNode[]): XNumberNode {
    if (nodes.length === 0) return null;

    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i];
      if (node.xNum === xNum) return node;

      if (nodes.length > i + 1) {
        const next = nodes[i + 1];
        if (xNum < next.xNum) {
          return this.findNode(xNum, node.children);
        }
      } else {
        return this.findNode(xNum, node.children);
      }
    }
  }

  getSlowestX = () => {
    if (this.xTree.nodes.length == 0) {
      return null;
    }

    const [slowestX] = this.getSlowestXNumber(this.xTree.nodes);
    return slowestX;
  };

  private getSlowestXNumber(aNodes: XNumberNode[]): [number, number] {
    let slowestX = Number.MAX_SAFE_INTEGER;
    let slowestTime = null;

    for (let aNode of aNodes) {
      let x = aNode.xNum;
      let time =
        this.getFromX(aNode.xNum).y[1] - this.getFromX(aNode.xNum).y[0];

      const [childX, childTime] = this.getSlowestXNumber(aNode.children);
      if (time < childTime) {
        x = childX;
        time = childTime;
      }

      if (slowestTime < time) {
        slowestX = x;
        slowestTime = time;
      }
    }

    return [slowestX, slowestTime];
  }

  getFocusedSeriesData(focusNumber: number): SeriesData[] {
    const seriesDatas: SeriesData[] = [];
    this.getSeriesDataRecursive(focusNumber, seriesDatas, this.xTree.nodes);

    seriesDatas.sort((a, b) => a.xNum - b.xNum);
    const skippedData = this.createSkippedSeriesDatas(seriesDatas);

    return seriesDatas.concat(skippedData);
  }

  private getSeriesDataRecursive(
    focusNumber: number,
    seriesArray: SeriesData[],
    nodes: XNumberNode[]
  ) {
    for (let i = 0; i < nodes.length; i++) {
      const child = nodes[i];
      seriesArray.push(this.getFromX(child.xNum));

      if (this.hasFocusNumberInsideNode(focusNumber, nodes, i)) {
        this.getSeriesDataRecursive(focusNumber, seriesArray, child.children);
      }
    }
  }

  private hasFocusNumberInsideNode(
    focusNumber: number,
    nodes: XNumberNode[],
    nodeIndex: number
  ) {
    const node = nodes[nodeIndex];
    if (nodeIndex == nodes.length - 1) {
      return node.xNum < focusNumber;
    }

    const next = nodes[nodeIndex + 1];
    return node.xNum < focusNumber && focusNumber < next.xNum;
  }

  private createSkippedSeriesDatas(
    usingSeriesData: SeriesData[]
  ): SeriesData[] {
    const usingXs = new Set(usingSeriesData.map((v) => v.xNum));
    const missingXs = Object.keys(this.xToSeriesRecords)
      .map((v) => Number(v))
      .filter((v) => !usingXs.has(v));

    const skippedSeriesData: SeriesData[] = [];
    const addedKeys = new Set<number>();
    for (let i = 0; i < missingXs.length; i++) {
      const x = missingXs[i];
      if (addedKeys.has(x)) continue;
      addedKeys.add(x);

      const seriesPointer = this.xToSeriesRecords[x];
      const parentNode = seriesPointer.treeNode.parent;
      if (addedKeys.has(parentNode.xNum)) continue;
      addedKeys.add(parentNode.xNum);
      const parentSeriesData = this.getFromX(parentNode.xNum);

      const [nodeStartAt, _, visitedXs] = this.collectNodeTimeSpan(
        seriesPointer.treeNode
      );
      for (const visitedX of visitedXs) {
        addedKeys.add(visitedX);
      }
      const startAt = Math.min(parentSeriesData.y[0], nodeStartAt);
      const endAt = parentSeriesData.y[1];
      skippedSeriesData.push(parentSeriesData.copyWith([startAt, endAt]));
    }

    return skippedSeriesData;
  }

  private collectNodeTimeSpan(node: XNumberNode): [number, number, number[]] {
    const seriesData = this.getFromX(node.xNum);
    const startAts = [seriesData.y[0]];
    const endAts = [seriesData.y[1]];
    let visitedXs = [node.xNum];

    for (const child of node.children) {
      const [startAt, endAt, childVisitedXs] = this.collectNodeTimeSpan(child);
      startAts.push(startAt);
      endAts.push(endAt);
      visitedXs = visitedXs.concat(childVisitedXs);
    }

    const startAt = Math.min(...startAts);
    const endAt = Math.max(...endAts);

    return [startAt, endAt, visitedXs];
  }
}
