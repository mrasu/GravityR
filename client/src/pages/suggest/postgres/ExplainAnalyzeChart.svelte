<script lang="ts">
  import { onMount } from "svelte";
  import { ExplainTree } from "@/components/ExplainTree/ExplainTree";
  import { getHighlightIndex } from "@/contexts/HighlightIndexContext";
  import { ExplainTreeChart } from "@/components/ExplainTree/ExplainTreeChart";
  import type { ExplainTreeChartProp } from "@/components/ExplainTree/ExplainTreeChart";
  import type { PostgresAnalyzeData } from "@/models/explain_data/PostgresAnalyzeData";

  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  export let chartDescription: string;
  export let analyzeNodes: PostgresAnalyzeData[];

  $: explainTree = new ExplainTree(analyzeNodes);

  let chartDiv: HTMLElement;
  onMount(() => {
    const chart = new ExplainTreeChart(chartDiv, explainTree, chartDescription);
    chart.onDataPointMouseEnter = (
      prop: ExplainTreeChartProp<PostgresAnalyzeData>
    ) => {
      $highlightIndex = prop.xNum;
    };
    chart.createTooltipFn = createTooltip;
    chart.render();

    return () => {
      chart.destroy();
    };
  });

  const createTooltip = (prop: ExplainTreeChartProp<PostgresAnalyzeData>) => {
    const trimmedText = prop.IAnalyzeData.text.trim();
    const text =
      trimmedText.length > 50
        ? trimmedText.substring(0, 50) + "..."
        : trimmedText.substring(0, 50);

    const {
      tableName,
      estimatedInitCost,
      estimatedCost,
      estimatedReturnedRows,
      estimatedWidth,
      actualTimeFirstRow,
      actualTimeAvg,
      actualReturnedRows,
      actualLoopCount,
    } = prop.IAnalyzeData;

    return `
      <div style='text-align: left; padding: 5px'>
        <div style='margin-bottom: 5px; font-weight: bold'>${text}</div>
        <div>Table:
          ${
            tableName ? `<span style='font-weight: bold'>${tableName}` : "-"
          }</span>
        </div>
        <div>Initial cost (estimated) : ${estimatedInitCost}</div>
        <div>Cost (estimated) : ${estimatedCost}</div>
        <div>Rows returned (estimated) : ${estimatedReturnedRows}</div>
        <div>Width (estimated) : ${estimatedWidth}</div>
        <!-- actualTimeAvg seems not the average time of loops for PostgreSQL even doc says so... -->
        <div>Time (actual) :  ${
          actualTimeFirstRow ? `${actualTimeFirstRow} ~ ${actualTimeAvg}` : "-"
        }</div>
        <div>Rows returned (actual) : ${actualReturnedRows ?? "-"}</div>
        <div>Loop count (actual) : ${actualLoopCount ?? "-"}</div>
      </div>`;
  };
</script>

<div id="tree" bind:this={chartDiv} />

<style>
</style>
