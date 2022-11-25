<script lang="ts">
  import { onMount } from "svelte";
  import { Checkbox } from "svelte-chota";
  import type { TimeDbLoad } from "@/models/TimeDbLoad";
  import { DbLoadRankingChart } from "./DbLoadRankingChart";
  import type { DbLoadRankingChartProp } from "./DbLoadRankingChart";
  import CommandCode from "@/lib/components/CommandCode.svelte";
  import { createLoadSumTooltip } from "./util/LoadTooltip";
  import { getAllDatabase, getFirstDatabaseAsArray } from "./util/Databases";
  import { buildCommandText } from "./util/CommandText";

  export let timeDbLoads: TimeDbLoad[];
  let chart: DbLoadRankingChart;
  let rankingChartDiv: HTMLElement;
  let selectedProp: DbLoadRankingChartProp;
  let checked: string[] = [];

  $: allDatabases = getAllDatabase(timeDbLoads);

  onMount(() => {
    checked = getFirstDatabaseAsArray(allDatabases);

    chart = new DbLoadRankingChart(rankingChartDiv, timeDbLoads, 20);
    chart.createTooltipFn = createTooltip;
    chart.onClickedFn = onBarClick;
    chart.render(checked);
  });

  const createTooltip = (prop: DbLoadRankingChartProp): string => {
    return createLoadSumTooltip(prop, 300);
  };

  const onBarClick = (prop: DbLoadRankingChartProp) => {
    selectedProp = prop;
  };

  const onDatabaseChanged = () => {
    chart?.changeDatabase(checked);
  };

  $: commandText = buildCommandText(selectedProp?.dbName, selectedProp?.sql);
</script>

<div class="wrapper">
  {#if allDatabases.length > 1}
    <div>
      Database:
      {#each allDatabases as database}
        <Checkbox
          value={database}
          on:change={onDatabaseChanged}
          bind:group={checked}>{database}</Checkbox
        >
      {/each}
    </div>
  {/if}

  <div id="sql-timeline-chart" bind:this={rankingChartDiv} />

  <div class="margin-top">
    With below command, you can get what index can accelerate your SQL
  </div>

  <div class="margin-top">
    <CommandCode {commandText} height={200} />
  </div>
</div>

<style>
  .wrapper {
    margin-top: 10px;
    margin-bottom: 10px;
  }

  .margin-top {
    margin-top: 10px;
  }

  #sql-timeline-chart {
    padding-bottom: 10px;
  }
</style>
