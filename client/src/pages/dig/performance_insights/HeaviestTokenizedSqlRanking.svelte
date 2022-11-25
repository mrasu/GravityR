<script lang="ts">
  import { onMount } from "svelte";
  import { Checkbox, Card, Radio } from "svelte-chota";
  import type { TimeDbLoad } from "@/models/TimeDbLoad";
  import CommandCode from "@/lib/components/CommandCode.svelte";
  import { createLoadSumTooltip } from "./util/LoadTooltip";
  import { DbTokenizedLoadRankingChart } from "./DbTokenizedLoadRankingChart";
  import type { DbTokenizedLoadRankingChartProp } from "./DbTokenizedLoadRankingChart";
  import type { RawLoadOfSql } from "./util/chart/TokenizedLoad";
  import { getAllDatabase, getFirstDatabaseAsArray } from "./util/Databases";
  import { buildCommandText } from "./util/CommandText";

  export let tokenizedTimeDbLoads: TimeDbLoad[];
  export let rawTimeDbLoads: TimeDbLoad[];
  let chart: DbTokenizedLoadRankingChart;
  let rankingChartDiv: HTMLElement;
  let selectedProp: DbTokenizedLoadRankingChartProp;
  let selectedRawLoadOfSql: RawLoadOfSql;
  let checked: string[] = [];

  $: allDatabases = getAllDatabase(tokenizedTimeDbLoads);

  onMount(() => {
    checked = getFirstDatabaseAsArray(allDatabases);

    chart = new DbTokenizedLoadRankingChart(
      rankingChartDiv,
      tokenizedTimeDbLoads,
      rawTimeDbLoads,
      20
    );

    chart.createTooltipFn = createTooltip;
    chart.onClickedFn = onBarClick;
    chart.render(checked);
  });

  const createTooltip = (prop: DbTokenizedLoadRankingChartProp): string => {
    return createLoadSumTooltip(prop, 300);
  };

  const onBarClick = (prop: DbTokenizedLoadRankingChartProp) => {
    selectedProp = prop;
    selectedRawLoadOfSql = undefined;
  };

  const onDatabaseChanged = () => {
    if (!checked.includes(selectedProp?.dbName)) {
      selectedProp = undefined;
      selectedRawLoadOfSql = undefined;
    }
    chart?.changeDatabase(checked);
  };

  $: commandText = buildCommandText(
    selectedRawLoadOfSql?.dbName,
    selectedRawLoadOfSql?.load?.sql
  );
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

  <div id="detail-sql-card">
    <Card>
      <h4 slot="header">Plain SQLs</h4>
      <div class="detail-sql-card-content">
        {#if selectedProp}
          {@const rawLoads = chart.getRawLoads(10, selectedProp.tokenizedId)}
          <table class="striped">
            <thead>
              <tr>
                <th class="radio" />
                <th class="db-name"><div>Database</div></th>
                <th class="db-load-avg"><div>db.load.avg</div></th>
                <th class="date"><div>Date</div></th>
                <th class="sql"><div>SQL</div></th>
              </tr>
            </thead>
            {#if rawLoads.length > 0}
              {#each rawLoads as rawLoad}
                <tr>
                  <td class="radio"
                    ><Radio
                      value={rawLoad}
                      bind:group={selectedRawLoadOfSql}
                    /></td
                  >
                  <td class="db-name">{rawLoad.dbName}</td>
                  <td class="db-load-avg"
                    ><div>
                      {(rawLoad.load.loadMax * 100).toFixed(2)}%
                    </div></td
                  >
                  <td class="date"><div>{rawLoad.date.toUTCString()}</div></td>
                  <td class="sql"><div>{rawLoad.load.sql}</div></td>
                </tr>
              {/each}
            {:else}
              <tr>
                <td>No data</td>
              </tr>
            {/if}
          </table>
        {:else}
          Click bar to see SQLs
        {/if}
      </div>
    </Card>
  </div>

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

  #detail-sql-card {
    margin-top: 30px;
  }
  .detail-sql-card-content {
    height: 300px;
    overflow: auto;
  }

  #detail-sql-card table .radio {
    width: 30px;
  }
  #detail-sql-card table .db-name {
    width: 100px;
  }
  #detail-sql-card table .db-load-avg {
    width: 100px;
  }
  #detail-sql-card table .date {
    width: 270px;
  }
</style>
