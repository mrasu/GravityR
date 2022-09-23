import SuggestHasuraPage from "@/pages/suggest/hasura/SuggestHasuraPage.svelte";
import { render, screen } from "@testing-library/svelte";
import { HasuraSuggestData } from "@/models/HasuraSuggestData";
import ApexCharts from "apexcharts";

describe("SuggestHasuraPage", () => {
  beforeEach(() => {
    jest.spyOn(ApexCharts.prototype, "render").mockReturnValue(null);
    jest.spyOn(ApexCharts.prototype, "destroy").mockReturnValue(null);
  });

  describe("smoke test", () => {
    const suggestData = new HasuraSuggestData({
      postgres: {
        gql: "query {}",
        gqlVariables: { a: 1 },
        query: "SELECT * FROM users",
        indexTargets: [],
        examinationCommandOptions: [],
        examinationResult: {
          originalTimeMillis: 11,
          indexResults: [],
        },
        analyzeNodes: [],
        planningText: "a",
      },
    });

    it("renders query", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("query");
      expect(data).toBeInTheDocument();
    });

    it("renders sql", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("sql");
      expect(data).toBeInTheDocument();
    });

    it("renders explain", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explain");
      expect(data).toBeInTheDocument();
    });

    it("renders explainChart", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("explainChart");
      expect(data).toBeInTheDocument();
    });

    it("renders suggest", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("suggest");
      expect(data).toBeInTheDocument();
    });

    it("renders examination", () => {
      render(SuggestHasuraPage, {
        suggestData: suggestData,
      });
      const data = screen.getByTestId("examination");
      expect(data).toBeInTheDocument();
    });
  });
});
