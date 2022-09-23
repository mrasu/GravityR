<script lang="ts">
  import type { DigData } from "@/models/DigData";
  import { Details, Tab, Tabs } from "svelte-chota";
  import HeavySqlTimeline from "./HeavySqlTimeline.svelte";
  import HeaviestSqlRanking from "./HeaviestSqlRanking.svelte";
  import HeavyTokenizedSqlTimeline from "./HeavyTokenizedSqlTimeline.svelte";
  import HeaviestTokenizedSqlRanking from "./HeaviestTokenizedSqlRanking.svelte";

  export let digData: DigData;
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
      <HeaviestSqlRanking timeDbLoads={digData.sqlDbLoads} />
    </Details>
  </div>

  <div class="details-card">
    <Details class="card" open data-testid="timeline">
      <span slot="summary">Heavy SQL Timeline</span>
      <HeavySqlTimeline timeDbLoads={digData.sqlDbLoads} />
    </Details>
  </div>
{:else}
  <div class="details-card">
    <Details class="card" open>
      <span slot="summary">Heaviest Tokenized SQLs</span>
      <HeaviestTokenizedSqlRanking
        tokenizedTimeDbLoads={digData.tokenizedSqlDbLoads}
        rawTimeDbLoads={digData.sqlDbLoads}
      />
    </Details>
  </div>

  <div class="details-card">
    <Details class="card" open>
      <span slot="summary">Heavy Tokenized SQL Timeline</span>
      <HeavyTokenizedSqlTimeline
        tokenizedTimeDbLoads={digData.tokenizedSqlDbLoads}
        rawTimeDbLoads={digData.sqlDbLoads}
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
