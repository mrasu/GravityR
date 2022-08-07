import type { IAnalyzeData } from "../../../types/gr_param";
import type { ExplainTreeSeriesData } from "./ExplainTreeChart";

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
type SeriesPointer = {
  seriesData: ExplainTreeSeriesData;
  treeNode: XNumberNode;
};

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
  const seriesData: ExplainTreeSeriesData = {
    x: nextXNumber.toString(),
    y: [0, 0],
    goals: [
      {
        value: 0,
        strokeColor: STROKE_COLOR,
      },
    ],
    prop: {
      xNum: nextXNumber,
      IAnalyzeData: IAnalyzeData,
    },
  };
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

  getFromX(xNum: number): ExplainTreeSeriesData {
    return this.xToSeriesRecords[xNum].seriesData;
  }

  getSeriesData(): ExplainTreeSeriesData[] {
    const seriesDatas: ExplainTreeSeriesData[] = [];
    this.getSeriesDataRecursive(seriesDatas, this.xTree.nodes);

    seriesDatas.sort((a, b) => a.prop.xNum - b.prop.xNum);
    return seriesDatas;
  }

  private getSeriesDataRecursive(
    seriesArray: ExplainTreeSeriesData[],
    nodes: XNumberNode[]
  ) {
    for (let i = 0; i < nodes.length; i++) {
      const child = nodes[i];
      seriesArray.push(this.getFromX(child.xNum));

      this.getSeriesDataRecursive(seriesArray, child.children);
    }
  }
}
