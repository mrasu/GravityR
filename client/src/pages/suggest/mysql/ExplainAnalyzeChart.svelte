<script lang="ts">
  import { onMount } from "svelte";
  import { ExplainTree } from "@/components/ExplainTree/ExplainTree";
  import { getHighlightIndex } from "@/contexts/HighlightIndexContext";
  import { ExplainTreeChart } from "@/components/ExplainTree/ExplainTreeChart";
  import type { ExplainTreeChartProp } from "@/components/ExplainTree/ExplainTreeChart";
  import type { MysqlAnalyzeData } from "@/models/explain_data/MysqlAnalyzeData";

  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  export let chartDescription: string;
  export let analyzeNodes: MysqlAnalyzeData[];
  $: explainTree = new ExplainTree(analyzeNodes);

  let chartDiv: HTMLElement;
  onMount(() => {
    const chart = new ExplainTreeChart(chartDiv, explainTree, chartDescription);
    chart.onDataPointMouseEnter = (
      prop: ExplainTreeChartProp<MysqlAnalyzeData>
    ) => {
      $highlightIndex = prop.xNum;
    };
    chart.createTooltipFn = createTooltip;
    chart.render();

    return () => {
      chart.destroy();
    };
  });

  const createTooltip = (prop: ExplainTreeChartProp<MysqlAnalyzeData>) => {
    const trimmedText = prop.IAnalyzeData.text.trim();
    const text =
      trimmedText.length > 50
        ? trimmedText.substring(0, 50) + "..."
        : trimmedText.substring(0, 50);

    const initCost = prop.IAnalyzeData.estimatedInitCost;
    const InitCostHtml = initCost
      ? `<div>Initial cost (estimated) : ${initCost}</div>`
      : "";

    const avgText = prop.IAnalyzeData.actualTimeAvg?.toString() ?? "-";
    const avgHtml =
      prop.IAnalyzeData.actualLoopCount > 1
        ? `<div>Time per loop (actual) : ${avgText}</div>`
        : `<div>Time (actual) : ${avgText}</div>`;

    const {
      tableName,
      estimatedCost,
      estimatedReturnedRows,
      actualTimeFirstRow,
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
        ${InitCostHtml}
        <div>Cost (estimated) : ${estimatedCost ?? "-"}</div>
        <div>Rows returned (estimated) : ${estimatedReturnedRows ?? "-"}</div>
        <div>Time first row (actual) : ${actualTimeFirstRow ?? "-"}</div>
        ${avgHtml}
        <div>Rows returned (actual) : ${actualReturnedRows ?? "-"}</div>
        <div>Loop count (actual) : ${actualLoopCount ?? "-"}</div>
      </div>`;
  };
</script>

<div id="tree" bind:this={chartDiv} />

<style>
</style>
