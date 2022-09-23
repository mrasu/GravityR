import SuggestMysqlPage from "@/pages/suggest/mysql/SuggestMysqlPage.svelte";
import { render, screen } from "@testing-library/svelte";
import { MysqlSuggestData } from "@/models/MysqlSuggestData";
import ApexCharts from "apexcharts";

describe("SuggestMysqlPage", () => {
  beforeEach(() => {
    jest.spyOn(ApexCharts.prototype, "render").mockReturnValue(null);
    jest.spyOn(ApexCharts.prototype, "destroy").mockReturnValue(null);
  });

  describe("smoke test", () => {
    const suggestData = new MysqlSuggestData({
      query: "SELECT * FROM users",
      indexTargets: [],
      examinationCommandOptions: [],
      examinationResult: {
        originalTimeMillis: 11,
        indexResults: [],
      },
      analyzeNodes: [],
    });

    it("renders sql", () => {
      render(SuggestMysqlPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("sql");
      expect(data).toBeInTheDocument();
    });

    it("renders explain", () => {
      render(SuggestMysqlPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explain");
      expect(data).toBeInTheDocument();
    });

    it("renders explainChart", () => {
      render(SuggestMysqlPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explainChart");
      expect(data).toBeInTheDocument();
    });

    it("renders suggest", () => {
      render(SuggestMysqlPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("suggest");
      expect(data).toBeInTheDocument();
    });

    it("renders examination", () => {
      render(SuggestMysqlPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("examination");
      expect(data).toBeInTheDocument();
    });
  });
});
