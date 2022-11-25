<script lang="ts">
  import { Details, Tab, Tabs } from "svelte-chota";
  import HeavySqlTimeline from "./HeavySqlTimeline.svelte";
  import HeaviestSqlRanking from "./HeaviestSqlRanking.svelte";
  import HeavyTokenizedSqlTimeline from "./HeavyTokenizedSqlTimeline.svelte";
  import HeaviestTokenizedSqlRanking from "./HeaviestTokenizedSqlRanking.svelte";
  import type { PerformanceInsightsData } from "@/models/PerformanceInsightsData";

  export let performanceInsightsData: PerformanceInsightsData;
  let activeTab = "plain";
</script>

<div class="tabs">
  <Tabs bind:active={activeTab}>
    <Tab tabid="plain">Plain</Tab>
    <Tab tabid="tokenized">Tokenized</Tab>
  </Tabs>
</div>

{#if activeTab === "plain"}
  <div class="details-card">
    <Details class="card" open data-testid="sqls">
      <span slot="summary">Heaviest SQLs</span>
      <HeaviestSqlRanking timeDbLoads={performanceInsightsData.sqlDbLoads} />
    </Details>
  </div>

  <div class="details-card">
    <Details class="card" open data-testid="timeline">
      <span slot="summary">Heavy SQL Timeline</span>
      <HeavySqlTimeline timeDbLoads={performanceInsightsData.sqlDbLoads} />
    </Details>
  </div>
{:else}
  <div class="details-card">
    <Details class="card" open>
      <span slot="summary">Heaviest Tokenized SQLs</span>
      <HeaviestTokenizedSqlRanking
        tokenizedTimeDbLoads={performanceInsightsData.tokenizedSqlDbLoads}
        rawTimeDbLoads={performanceInsightsData.sqlDbLoads}
      />
    </Details>
  </div>

  <div class="details-card">
    <Details class="card" open>
      <span slot="summary">Heavy Tokenized SQL Timeline</span>
      <HeavyTokenizedSqlTimeline
        tokenizedTimeDbLoads={performanceInsightsData.tokenizedSqlDbLoads}
        rawTimeDbLoads={performanceInsightsData.sqlDbLoads}
      />
    </Details>
  </div>
{/if}

<style>
  .tabs {
    margin-bottom: 20px;
  }

  .details-card {
    margin-bottom: 30px;
  }
</style>
