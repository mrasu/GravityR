<script lang="ts">
  import { onMount } from "svelte";
  import { ExplainAnalyzeTree } from "./ExplainAnalyzeTree";
  import type { SeriesData } from "./SeriesData";
  import type { IAnalyzeData } from "../../types/gr_param";
  import { getHighlightIndex } from "../../contexts/HighlightIndexContext";
  import { ExplainTree } from "./ExplainTree";

  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  export let analyzeNodes: IAnalyzeData[];
  $: explainAnalyzeTree = new ExplainAnalyzeTree(analyzeNodes);

  let chartDiv: HTMLElement;
  onMount(() => {
    const tree = new ExplainTree(chartDiv, explainAnalyzeTree);
    tree.onMouseEntered = (x: number) => {
      $highlightIndex = x;
    };
    tree.createTooltipFn = createTooltip;
    tree.render();
  });

  const createTooltip = (data: SeriesData) => {
    const trimmedText = data.text.trim();
    const text =
      trimmedText.length > 50
        ? trimmedText.substring(0, 50) + "..."
        : trimmedText.substring(0, 50);

    const initCost = data.IAnalyzeData.estimatedInitCost;
    const InitCostHtml = initCost
      ? `<div>Initial cost (estimated) : ${initCost}</div>`
      : "";

    const avgText = data.IAnalyzeData.actualTimeAvg?.toString() ?? "-";
    const avgHtml =
      data.IAnalyzeData.actualLoopCount > 1
        ? `<div>Time per loop (actual) : ${avgText}</div>`
        : `<div>Time (actual) : ${avgText}</div>`;

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
        <div>Cost (estimated) : ${data.IAnalyzeData.estimatedCost ?? "-"}</div>
        <div>Rows returned (estimated) : ${
          data.IAnalyzeData.estimatedReturnedRows ?? "-"
        }</div>
        <div>Time first row (actual) : ${
          data.IAnalyzeData.actualTimeFirstRow ?? "-"
        }</div>
        ${avgHtml}
        <div>Rows returned (actual) : ${
          data.IAnalyzeData.actualReturnedRows ?? "-"
        }</div>
        <div>Loop count (actual) : ${
          data.IAnalyzeData.actualLoopCount ?? "-"
        }</div>
      </div>`;
  };
</script>

<div id="tree" bind:this={chartDiv} />

<style>
</style>
