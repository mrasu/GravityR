import { ExplainAnalyzeTree } from "@/components/ExplainTree/ExplainAnalyzeTree";
import type { IDbAnalyzeData } from "@/types/gr_param";

describe("ExplainAnalyzeTree", () => {
  const buildIAnalyzeData = (
    time: number,
    children: IDbAnalyzeData[] = null
  ) => {
    return {
      text: time.toString(),
      title: "",
      actualTimeAvg: time,
      actualLoopCount: 1,
      children: children,
    };
  };

  describe("getSeriesData", () => {
    const convertXToText = (
      tree: ExplainAnalyzeTree<IDbAnalyzeData>,
      x: number
    ) => {
      return tree.getFromX(x).prop.IAnalyzeData.text;
    };

    it("no data", () => {
      const tree = new ExplainAnalyzeTree([], true);
      expect(tree.getSeriesData()).toEqual([]);
    });

    describe("single child only tree", () => {
      const tree = new ExplainAnalyzeTree(
        [
          buildIAnalyzeData(0, [
            buildIAnalyzeData(1, [
              buildIAnalyzeData(2, [
                buildIAnalyzeData(3, [buildIAnalyzeData(4)]),
              ]),
            ]),
          ]),
        ],
        true
      );

      it("focus at the bottom", () => {
        const resultTexts = tree
          .getSeriesData()
          .map((v) => convertXToText(tree, v.prop.xNum));

        expect(resultTexts).toEqual(["0", "1", "2", "3", "4"]);
      });
    });

    describe("multi child only tree", () => {
      const tree = new ExplainAnalyzeTree(
        [
          buildIAnalyzeData(0, [
            buildIAnalyzeData(1, [
              buildIAnalyzeData(2, [
                buildIAnalyzeData(3, [
                  buildIAnalyzeData(4),
                  buildIAnalyzeData(5),
                ]),
                buildIAnalyzeData(6, [
                  buildIAnalyzeData(7),
                  buildIAnalyzeData(8),
                ]),
              ]),
            ]),
            buildIAnalyzeData(9, [
              buildIAnalyzeData(10, [
                buildIAnalyzeData(11, [
                  buildIAnalyzeData(12),
                  buildIAnalyzeData(13),
                ]),
                buildIAnalyzeData(14, [
                  buildIAnalyzeData(15),
                  buildIAnalyzeData(16),
                ]),
              ]),
            ]),
          ]),
          buildIAnalyzeData(17, [
            buildIAnalyzeData(18, [
              buildIAnalyzeData(19, [
                buildIAnalyzeData(20, [
                  buildIAnalyzeData(21),
                  buildIAnalyzeData(22),
                ]),
                buildIAnalyzeData(23, [
                  buildIAnalyzeData(24),
                  buildIAnalyzeData(25),
                ]),
              ]),
            ]),
            buildIAnalyzeData(26, [
              buildIAnalyzeData(27, [
                buildIAnalyzeData(28, [
                  buildIAnalyzeData(29),
                  buildIAnalyzeData(30),
                ]),
                buildIAnalyzeData(31, [
                  buildIAnalyzeData(32),
                  buildIAnalyzeData(33),
                ]),
              ]),
            ]),
          ]),
        ],
        true
      );

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
