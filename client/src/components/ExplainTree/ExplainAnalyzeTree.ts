import type { IDbAnalyzeData } from "@/types/gr_param";
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

type XToSeriesRecords<D extends IDbAnalyzeData> = Record<
  number,
  SeriesPointer<D>
>;
type SeriesPointer<D extends IDbAnalyzeData> = {
  seriesData: ExplainTreeSeriesData<D>;
  treeNode: XNumberNode;
};

const createXTree = <D extends IDbAnalyzeData>(
  IAnalyzeDatas: IDbAnalyzeData[],
  considersActualTimeIsAverageIfLooped: boolean
): [XNumberTree, XToSeriesRecords<D>] => {
  const xTree = new XNumberTree();
  const xToSeriesRecords: XToSeriesRecords<D> = {};

  let nextXNumber = 0;
  for (let i = 0; i < IAnalyzeDatas.length; i++) {
    const xNode = new XNumberNode(nextXNumber, null);
    xTree.nodes.push(xNode);
    [nextXNumber] = visitAnalyze(
      xNode,
      xToSeriesRecords,
      nextXNumber,
      IAnalyzeDatas[i],
      considersActualTimeIsAverageIfLooped
    );
  }

  return [xTree, xToSeriesRecords];
};

const STROKE_COLOR = "#CD2F2A";

const visitAnalyze = <D extends IDbAnalyzeData>(
  parentNode: XNumberNode,
  xToSeriesRecords: XToSeriesRecords<D>,
  nextXNumber: number,
  IAnalyzeData: D,
  considersActualTimeIsAverageIfLooped: boolean
): [number, number] => {
  let maxFinished = 0;
  const seriesData: ExplainTreeSeriesData<D> = {
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
        data,
        considersActualTimeIsAverageIfLooped
      );
      childrenFinished.push(childFinished);
    }
    startTime = Math.max(...childrenFinished);
  }

  let endTimeFromData: number;
  if (
    considersActualTimeIsAverageIfLooped &&
    IAnalyzeData.actualLoopCount > 1
  ) {
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

export class ExplainAnalyzeTree<D extends IDbAnalyzeData> {
  xTree: XNumberTree;
  xToSeriesRecords: XToSeriesRecords<D>;

  constructor(
    IAnalyzeDatas: IDbAnalyzeData[],
    considersActualTimeIsAverageIfLooped: boolean
  ) {
    [this.xTree, this.xToSeriesRecords] = createXTree(
      IAnalyzeDatas,
      considersActualTimeIsAverageIfLooped
    );
  }

  getFromX(xNum: number): ExplainTreeSeriesData<D> {
    return this.xToSeriesRecords[xNum].seriesData;
  }

  getSeriesData(): ExplainTreeSeriesData<D>[] {
    const seriesDatas: ExplainTreeSeriesData<D>[] = [];
    this.getSeriesDataRecursive(seriesDatas, this.xTree.nodes);

    seriesDatas.sort((a, b) => a.prop.xNum - b.prop.xNum);
    return seriesDatas;
  }

  private getSeriesDataRecursive(
    seriesArray: ExplainTreeSeriesData<D>[],
    nodes: XNumberNode[]
  ) {
    for (let i = 0; i < nodes.length; i++) {
      const child = nodes[i];
      seriesArray.push(this.getFromX(child.xNum));

      this.getSeriesDataRecursive(seriesArray, child.children);
    }
  }
}
