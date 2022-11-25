import { plainToInstance } from "class-transformer";
import { JaegerData } from "@/models/JaegerData";
import { render, screen } from "@testing-library/svelte";
import DigJaegerPage from "@/pages/dig/jaeger/DigJaegerPage.svelte";

describe("DigJaegerPage", () => {
  describe("smoke test", () => {
    const jaegerData = plainToInstance(JaegerData, {
      uiPath: "http://localhost:16686",
      slowThresholdMilli: 100,
      sameServiceThreshold: 5,
      slowTraces: [],
      sameServiceTraces: [],
    });

    it("renders slowTraces", () => {
      render(DigJaegerPage, {
        jaegerData: jaegerData,
      });
      const data = screen.getByTestId("slowTraces");
      expect(data).toBeInTheDocument();
    });

    it("renders sameServiceTraces", () => {
      render(DigJaegerPage, {
        jaegerData: jaegerData,
      });
      const data = screen.getByTestId("sameServiceTraces");
      expect(data).toBeInTheDocument();
    });
  });
});
