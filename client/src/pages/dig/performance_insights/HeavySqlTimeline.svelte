<script lang="ts">
  import { onMount } from "svelte";
  import { DbLoadChart } from "./DbLoadChart";
  import type { DbLoadChartProp } from "./DbLoadChart";
  import { Checkbox } from "svelte-chota";
  import type { TimeDbLoad } from "@/models/TimeDbLoad";
  import CommandCode from "@/lib/components/CommandCode.svelte";
  import { createLoadMaxTooltip } from "./util/LoadTooltip";
  import { getAllDatabase, getFirstDatabaseAsArray } from "./util/Databases";
  import { buildCommandText } from "./util/CommandText";

  export let timeDbLoads: TimeDbLoad[];
  let chart: DbLoadChart;
  let sqlTimelineChartDiv: HTMLElement;
  let selectedProp: DbLoadChartProp;
  let checked: string[] = [];

  $: allDatabases = getAllDatabase(timeDbLoads);

  onMount(() => {
    checked = getFirstDatabaseAsArray(allDatabases);

    chart = new DbLoadChart(sqlTimelineChartDiv, timeDbLoads, 0.005);
    chart.createTooltipFn = createTooltip;
    chart.onClickedFn = onBarClick;
    chart.render(checked);
  });

  const createTooltip = (prop: DbLoadChartProp): string => {
    return createLoadMaxTooltip(prop, 300);
  };

  const onBarClick = (prop: DbLoadChartProp) => {
    selectedProp = prop;
  };

  const onDatabaseChanged = () => {
    if (!checked.includes(selectedProp?.dbName)) {
      selectedProp = undefined;
    }
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

  <div id="sql-timeline-chart" bind:this={sqlTimelineChartDiv} />

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
    overflow-y: hidden;
    padding-bottom: 10px;
  }
</style>
