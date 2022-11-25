import DigPerformanceInsightsPage from "@/pages/dig/performance_insights/DigPerformanceInsightsPage.svelte";
import { render, screen } from "@testing-library/svelte";
import { plainToInstance } from "class-transformer";
import { PerformanceInsightsData } from "@/models/PerformanceInsightsData";
import ApexCharts from "apexcharts";

describe("DigPerformanceInsightsPage", () => {
  beforeEach(() => {
    jest.spyOn(ApexCharts.prototype, "render").mockReturnValue(null);
  });

  describe("smoke test", () => {
    const performanceInsightsData = plainToInstance(PerformanceInsightsData, {
      sqlDbLoads: [],
      tokenizedSqlDbLoads: [],
    });

    it("renders sqls", () => {
      render(DigPerformanceInsightsPage, {
        performanceInsightsData: performanceInsightsData,
      });
      const data = screen.getByTestId("sqls");
      expect(data).toBeInTheDocument();
    });

    it("renders timeline", () => {
      render(DigPerformanceInsightsPage, {
        performanceInsightsData: performanceInsightsData,
      });
      const data = screen.getByTestId("timeline");
      expect(data).toBeInTheDocument();
    });
  });
});
