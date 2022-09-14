<script lang="ts">
  import { onMount } from "svelte";
  import { ExplainTree } from "@/components/ExplainTree/ExplainTree";
  import { getHighlightIndex } from "@/contexts/HighlightIndexContext";
  import { ExplainTreeChart } from "@/components/ExplainTree/ExplainTreeChart";
  import type { ExplainTreeChartProp } from "@/components/ExplainTree/ExplainTreeChart";
  import type { PostgresExplainData } from "@/models/explain_data/PostgresExplainData";

  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  export let chartDescription: string;
  export let analyzeNodes: PostgresExplainData[];

  $: explainTree = new ExplainTree(analyzeNodes);

  let chartDiv: HTMLElement;
  onMount(() => {
    const chart = new ExplainTreeChart(chartDiv, explainTree, chartDescription);
    chart.onDataPointMouseEnter = (
      prop: ExplainTreeChartProp<PostgresExplainData>
    ) => {
      $highlightIndex = prop.xNum;
    };
    chart.createTooltipFn = createTooltip;
    chart.render();

    return () => {
      chart.destroy();
    };
  });

  const createTooltip = (prop: ExplainTreeChartProp<PostgresExplainData>) => {
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
      </div>`;
  };
</script>

<div id="tree" bind:this={chartDiv} />

<style>
</style>
