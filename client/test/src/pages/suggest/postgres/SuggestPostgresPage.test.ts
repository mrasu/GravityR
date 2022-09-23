import SuggestPostgresPage from "@/pages/suggest/postgres/SuggestPostgresPage.svelte";
import { render, screen } from "@testing-library/svelte";
import { PostgresSuggestData } from "@/models/PostgresSuggestData";
import ApexCharts from "apexcharts";

describe("SuggestPostgresPage", () => {
  beforeEach(() => {
    jest.spyOn(ApexCharts.prototype, "render").mockReturnValue(null);
    jest.spyOn(ApexCharts.prototype, "destroy").mockReturnValue(null);
  });

  describe("smoke test", () => {
    const suggestData = new PostgresSuggestData({
      query: "SELECT * FROM users",
      indexTargets: [],
      examinationCommandOptions: [],
      examinationResult: {
        originalTimeMillis: 11,
        indexResults: [],
      },
      analyzeNodes: [],
      planningText: "",
    });

    it("renders sql", () => {
      render(SuggestPostgresPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("sql");
      expect(data).toBeInTheDocument();
    });

    it("renders explain", () => {
      render(SuggestPostgresPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explain");
      expect(data).toBeInTheDocument();
    });

    it("renders explainChart", () => {
      render(SuggestPostgresPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explainChart");
      expect(data).toBeInTheDocument();
    });

    it("renders suggest", () => {
      render(SuggestPostgresPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("suggest");
      expect(data).toBeInTheDocument();
    });

    it("renders examination", () => {
      render(SuggestPostgresPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("examination");
      expect(data).toBeInTheDocument();
    });
  });
});
