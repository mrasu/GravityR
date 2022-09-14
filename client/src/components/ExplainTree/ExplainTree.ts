import type { ExplainTreeSeriesData } from "./ExplainTreeChart";
import type { DbExplainData } from "@/models/explain_data/DbExplainData";

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

type XToSeriesRecords<D extends DbExplainData> = Record<
  number,
  SeriesPointer<D>
>;
type SeriesPointer<D extends DbExplainData> = {
  seriesData: ExplainTreeSeriesData<D>;
  treeNode: XNumberNode;
};

const createXTree = <D extends DbExplainData>(
  IAnalyzeDatas: DbExplainData[]
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
      IAnalyzeDatas[i]
    );
  }

  return [xTree, xToSeriesRecords];
};

const STROKE_COLOR = "#CD2F2A";

const visitAnalyze = <D extends DbExplainData>(
  parentNode: XNumberNode,
  xToSeriesRecords: XToSeriesRecords<D>,
  nextXNumber: number,
  IAnalyzeData: D
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
        data
      );
      childrenFinished.push(childFinished);
    }
    startTime = Math.max(...childrenFinished);
  }

  const endTimeFromData = IAnalyzeData.calculateEndTime(startTime);
  const endTime = Math.max(endTimeFromData, startTime);
  maxFinished = Math.max(maxFinished, endTime);
  seriesData.y = [Math.min(startTime, endTime - 0.0001), endTime];
  seriesData.goals[0].value = endTime;

  return [nextXNumber, maxFinished];
};

export class ExplainTree<D extends DbExplainData> {
  xTree: XNumberTree;
  xToSeriesRecords: XToSeriesRecords<D>;

  constructor(IAnalyzeDatas: DbExplainData[]) {
    [this.xTree, this.xToSeriesRecords] = createXTree(IAnalyzeDatas);
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
