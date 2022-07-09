<script lang="ts">
  import ApexCharts from "apexcharts";
  import { onMount } from "svelte";
  import { ExplainAnalyzeTree } from "./ExplainAnalyzeTree";
  import type { SeriesData } from "./SeriesData";
  import type { IAnalyzeData } from "../../types/gr_param";
  import { getHighlightIndex } from "../../contexts/HighlightIndexContext";

  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  // TODO: wrap Tree and encapsulate `animating` and `mouseEnteringX` and expose current entering x like `onXChanged` method
  let animating = false;
  let mouseEnteringX: number | undefined = undefined;
  const options = {
    series: [
      {
        data: [],
      },
    ],
    dataLabels: {
      enabled: true,
      textAnchor: "middle",
      style: {
        color: "rgb(0 143 251 / 20%)",
      },
      formatter: (
        _: any,
        {
          seriesIndex,
          dataPointIndex,
          w: {
            config: { series },
          },
        }
      ) => {
        return series[seriesIndex].data[dataPointIndex].title;
      },
    },
    grid: {
      xaxis: { lines: { show: true } },
      yaxis: { lines: { show: false } },
    },
    chart: {
      height: "", // change dynamically
      width: "100%",
      type: "rangeBar",
      events: {
        animationEnd: () => {
          animating = false;
          if (mouseEnteringX) {
            $highlightIndex = mouseEnteringX;
          }
        },
        dataPointSelection: (
          _: any,
          {
            w: {
              config: { series },
            },
          },
          { seriesIndex, dataPointIndex }
        ) => {
          const x = series[seriesIndex].data[dataPointIndex].x;
          updateSeriesData(series[seriesIndex].data, Number(x));
        },
        dataPointMouseEnter: (
          event,
          {
            w: {
              config: { series },
            },
          },
          { seriesIndex, dataPointIndex }
        ) => {
          const x = Number(series[seriesIndex].data[dataPointIndex].x);
          mouseEnteringX = x;
          if (animating) return;

          event.path[0].style.cursor = "pointer";
          $highlightIndex = x;
        },
        dataPointMouseLeave: () => {
          mouseEnteringX = undefined;
        },
      },
    },
    plotOptions: {
      bar: {
        horizontal: true,
        barHeight: "80%",
      },
    },
    fill: {
      type: "solid",
      opacity: 0.6,
    },
    xaxis: {
      type: "numeric",
      title: {
        text: "Execution time (ms)",
      },
    },
    yaxis: {
      show: true,
      labels: {
        show: false,
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    stroke: {
      colors: ["transparent"],
      width: 0,
    },
    tooltip: {
      enabled: true,
      custom: ({
        seriesIndex,
        dataPointIndex,
        w: {
          config: { series },
        },
      }) => {
        const data = series[seriesIndex].data[dataPointIndex];
        const trimmedText = data.text.trim();
        const text =
          trimmedText.length > 50
            ? trimmedText.substring(0, 50) + "..."
            : trimmedText.substring(0, 50);

        const initCost = data.IAnalyzeData.estimatedInitCost;
        const InitCostHtml = initCost
          ? `<div>EstimatedInitCost: ${initCost}</div>`
          : "";

        return `
          <div style='text-align: left; padding: 5px'>
            <div style='margin-bottom: 5px; font-weight: bold'>${text}</div>
            <div>Table:
              ${
                data.IAnalyzeData.tableName
                  ? `<span style='font-weight: bold'>${data.IAnalyzeData.tableName}`
                  : "-"
              }</span>
            </div>
            ${InitCostHtml}
            <div>EstimatedCost: ${data.IAnalyzeData.estimatedCost ?? "-"}</div>
            <div>EstimatedReturnedRows: ${
              data.IAnalyzeData.estimatedReturnedRows ?? "-"
            }</div>
            <div>ActualTimeFirstRow: ${
              data.IAnalyzeData.actualTimeFirstRow ?? "-"
            }</div>
            <div>ActualTimeAvg: ${data.IAnalyzeData.actualTimeAvg ?? "-"}</div>
            <div>ActualReturnedRow: ${
              data.IAnalyzeData.actualReturnedRows ?? "-"
            }</div>
            <div>ActualLoopCount: ${
              data.IAnalyzeData.actualLoopCount ?? "-"
            }</div>
          </div>`;
      },
    },
    title: {
      text: "Execution time based timeline from EXPLAIN ANALYZE",
    },
  };

  export let analyzeNodes: IAnalyzeData[];
  $: explainAnalyzeTree = new ExplainAnalyzeTree(analyzeNodes);

  const barHeight = 30;
  const calculateChartHeightPx = (seriesData: SeriesData[]): string => {
    const height =
      new Set(seriesData.map((highlightIndex) => highlightIndex.x)).size *
        barHeight +
      130;
    return `${height}px`;
  };

  const updateSeriesData = async (currentSeries: SeriesData[], x: number) => {
    return new Promise<void>((resolve) => {
      setTimeout(async () => {
        const nexts = explainAnalyzeTree.getNextXs(x);
        if (!nexts) return;

        const nextX = nexts[0];
        const isOpen = !!currentSeries.find((d) => Number(d.x) === nextX);
        let seriesData: SeriesData[];
        if (isOpen) {
          seriesData = explainAnalyzeTree.getFocusedSeriesData(x);
        } else {
          seriesData = explainAnalyzeTree.getFocusedSeriesData(nextX);
        }

        animating = true;
        await chart.updateOptions({
          series: [
            {
              data: seriesData,
            },
          ],
          chart: { height: calculateChartHeightPx(seriesData) },
        });
        resolve();
      });
    });
  };

  let chartDiv: HTMLElement;
  let chart: ApexCharts;
  onMount(() => {
    const seriesData2 = explainAnalyzeTree.getFocusedSeriesData(
      explainAnalyzeTree.getSlowestX()
    );
    options.series = [{ data: seriesData2 }];
    options.chart.height = calculateChartHeightPx(seriesData2);
    chart = new ApexCharts(chartDiv, options);
    chart.render();
  });
</script>

<div id="chart" bind:this={chartDiv} />

<style>
</style>
