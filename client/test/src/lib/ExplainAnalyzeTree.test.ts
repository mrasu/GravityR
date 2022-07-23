import { ExplainAnalyzeTree } from "../../../src/lib/ExplainAnalyzeChart/ExplainAnalyzeTree";
import type { IAnalyzeData } from "../../../src/types/gr_param";

describe("ExplainAnalyzeTree", () => {
  const buildIAnalyzeData = (time: number, children: IAnalyzeData[] = null) => {
    return {
      text: time.toString(),
      title: "",
      actualTimeAvg: time,
      actualLoopCount: 1,
      children: children,
    };
  };

  describe("getSlowestXNumber", () => {
    it("no data", () => {
      const tree = new ExplainAnalyzeTree([]);
      expect(tree.getSlowestX()).toBe(null);
    });

    describe("single child only tree", () => {
      it("the slowest at the bottom", () => {
        const tree = new ExplainAnalyzeTree([
          buildIAnalyzeData(14, [
            buildIAnalyzeData(13, [
              buildIAnalyzeData(12, [
                buildIAnalyzeData(11, [buildIAnalyzeData(10)]),
              ]),
            ]),
          ]),
        ]);

        const xNum = tree.getSlowestX();
        expect(xNum).toBe(4);
        expect(tree.getFromX(xNum).text).toBe("10");
      });

      it("the slowest in the middle", () => {
        const tree = new ExplainAnalyzeTree([
          buildIAnalyzeData(14, [
            buildIAnalyzeData(13, [
              buildIAnalyzeData(12, [
                buildIAnalyzeData(2, [buildIAnalyzeData(1)]),
              ]),
            ]),
          ]),
        ]);

        const xNum = tree.getSlowestX();
        expect(xNum).toBe(2);
        expect(tree.getFromX(xNum).text).toBe("12");
      });
    });

    it("multi child only tree", () => {
      const tree = new ExplainAnalyzeTree([
        buildIAnalyzeData(14, [
          buildIAnalyzeData(13, [
            buildIAnalyzeData(12, [
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
            ]),
          ]),
          buildIAnalyzeData(1013, [
            buildIAnalyzeData(1012, [
              buildIAnalyzeData(13, [
                buildIAnalyzeData(12),
                buildIAnalyzeData(11),
              ]),
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
            ]),
          ]),
        ]),
        buildIAnalyzeData(114, [
          buildIAnalyzeData(113, [
            buildIAnalyzeData(112, [
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
            ]),
          ]),
          buildIAnalyzeData(113, [
            buildIAnalyzeData(112, [
              buildIAnalyzeData(13, [
                buildIAnalyzeData(12),
                buildIAnalyzeData(11),
              ]),
              buildIAnalyzeData(3, [
                buildIAnalyzeData(2),
                buildIAnalyzeData(1),
              ]),
            ]),
          ]),
        ]),
      ]);

      const xNum = tree.getSlowestX();
      expect(xNum).toBe(10);
      expect(tree.getFromX(xNum).text).toBe("1012");
    });
  });

  describe("getSeriesData", () => {
    const convertXToText = (tree: ExplainAnalyzeTree, x: number) => {
      return tree.getFromX(x).text;
    };

    it("no data", () => {
      const tree = new ExplainAnalyzeTree([]);
      expect(tree.getFocusedSeriesData(0)).toEqual([]);
    });

    describe("single child only tree", () => {
      const tree = new ExplainAnalyzeTree([
        buildIAnalyzeData(0, [
          buildIAnalyzeData(1, [
            buildIAnalyzeData(2, [
              buildIAnalyzeData(3, [buildIAnalyzeData(4)]),
            ]),
          ]),
        ]),
      ]);

      it("focus at the bottom", () => {
        const resultTexts = tree
          .getFocusedSeriesData(4)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual(["0", "1", "2", "3", "4"]);
      });

      it("focus in the middle", () => {
        const resultTexts = tree
          .getFocusedSeriesData(2)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual(["0", "1", "2", "2"]);
      });
    });

    describe("multi child only tree", () => {
      const tree = new ExplainAnalyzeTree([
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
      ]);

      it("focus at the top", () => {
        const resultTexts = tree
          .getFocusedSeriesData(0)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual(["0", "17", "0", "17"]);
      });

      it("focus in the middle", () => {
        const resultTexts = tree
          .getFocusedSeriesData(11)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual([
          "0",
          "1",
          "9",
          "10",
          "11",
          "14",
          "17",
          "1",
          "11",
          "14",
          "17",
        ]);
      });

      it("focus in the middle without siblings", () => {
        const resultTexts = tree
          .getFocusedSeriesData(10)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual([
          "0",
          "1",
          "9",
          "10",
          "17",
          "1",
          "10",
          "17",
        ]);
      });

      it("focus at the top-level", () => {
        const resultTexts = tree
          .getFocusedSeriesData(17)
          .map((v) => convertXToText(tree, v.xNum));

        expect(resultTexts).toEqual(["0", "17", "0", "17"]);
      });
    });
  });
});
