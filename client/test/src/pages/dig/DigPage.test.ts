import DigPage from "@/pages/dig/DigPage.svelte";
import { render, screen } from "@testing-library/svelte";
import { plainToInstance } from "class-transformer";
import { DigData } from "@/models/DigData";
import ApexCharts from "apexcharts";

describe("DigPage", () => {
  beforeEach(() => {
    jest.spyOn(ApexCharts.prototype, "render").mockReturnValue(null);
  });

  describe("smoke test", () => {
    const digData = plainToInstance(DigData, {
      sqlDbLoads: [],
      tokenizedSqlDbLoads: [],
    });

    it("renders sqls", () => {
      render(DigPage, {
        digData: digData,
      });
      const data = screen.getByTestId("sqls");
      expect(data).toBeInTheDocument();
    });

    it("renders timeline", () => {
      render(DigPage, {
        digData: digData,
      });
      const data = screen.getByTestId("timeline");
      expect(data).toBeInTheDocument();
    });
  });
});
