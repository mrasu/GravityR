import { ExplainTree } from "@/components/ExplainTree/ExplainTree";
import { MysqlAnalyzeData } from "@/models/explain_data/MysqlAnalyzeData";

describe("ExplainAnalyzeTree", () => {
  const buildMysqlExplainData = (
    time: number,
    children: MysqlAnalyzeData[] = null
  ) => {
    const data = new MysqlAnalyzeData();
    data.text = time.toString();
    data.title = "";
    data.actualTimeAvg = time;
    data.actualLoopCount = 1;
    data.children = children;

    return data;
  };

  describe("getSeriesData", () => {
    const convertXToText = (tree: ExplainTree<MysqlAnalyzeData>, x: number) => {
      return tree.getFromX(x).prop.IAnalyzeData.text;
    };

    it("no data", () => {
      const tree = new ExplainTree([]);
      expect(tree.getSeriesData()).toEqual([]);
    });

    describe("single child only tree", () => {
      const tree = new ExplainTree([
        buildMysqlExplainData(0, [
          buildMysqlExplainData(1, [
            buildMysqlExplainData(2, [
              buildMysqlExplainData(3, [buildMysqlExplainData(4)]),
            ]),
          ]),
        ]),
      ]);

      it("focus at the bottom", () => {
        const resultTexts = tree
          .getSeriesData()
          .map((v) => convertXToText(tree, v.prop.xNum));

        expect(resultTexts).toEqual(["0", "1", "2", "3", "4"]);
      });
    });

    describe("multi child only tree", () => {
      const tree = new ExplainTree([
        buildMysqlExplainData(0, [
          buildMysqlExplainData(1, [
            buildMysqlExplainData(2, [
              buildMysqlExplainData(3, [
                buildMysqlExplainData(4),
                buildMysqlExplainData(5),
              ]),
              buildMysqlExplainData(6, [
                buildMysqlExplainData(7),
                buildMysqlExplainData(8),
              ]),
            ]),
          ]),
          buildMysqlExplainData(9, [
            buildMysqlExplainData(10, [
              buildMysqlExplainData(11, [
                buildMysqlExplainData(12),
                buildMysqlExplainData(13),
              ]),
              buildMysqlExplainData(14, [
                buildMysqlExplainData(15),
                buildMysqlExplainData(16),
              ]),
            ]),
          ]),
        ]),
        buildMysqlExplainData(17, [
          buildMysqlExplainData(18, [
            buildMysqlExplainData(19, [
              buildMysqlExplainData(20, [
                buildMysqlExplainData(21),
                buildMysqlExplainData(22),
              ]),
              buildMysqlExplainData(23, [
                buildMysqlExplainData(24),
                buildMysqlExplainData(25),
              ]),
            ]),
          ]),
          buildMysqlExplainData(26, [
            buildMysqlExplainData(27, [
              buildMysqlExplainData(28, [
                buildMysqlExplainData(29),
                buildMysqlExplainData(30),
              ]),
              buildMysqlExplainData(31, [
                buildMysqlExplainData(32),
                buildMysqlExplainData(33),
              ]),
            ]),
          ]),
        ]),
      ]);

      it("focus at the top-level", () => {
        const resultTexts = tree
          .getSeriesData()
          .map((v) => convertXToText(tree, v.prop.xNum));

        const nums = [];
        for (let i = 0; i < 34; i++) {
          nums.push(i.toString());
        }
        expect(resultTexts).toEqual(nums);
      });
    });
  });
});
