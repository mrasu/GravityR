<script lang="ts">
  import "chota";
  import "reflect-metadata";
  import ExplainAnalyzeChart from "./lib/ExplainAnalyzeChart/ExplainAnalyzeChart.svelte";
  import { grParam } from "./stores/gr_param";
  import { Details } from "svelte-chota";
  import ExplainAnalyzeText from "./lib/ExplainAnalyzeText.svelte";
  import QueryText from "./lib/QueryText.svelte";
  import IndexSuggestion from "./lib/IndexSuggestion.svelte";
  import ExaminationResultTable from "./lib/ExaminationResultTable.svelte";
  import { createHighlightIndexContext } from "./contexts/HighlightIndexContext";

  const highlightIndexKey = createHighlightIndexContext();
</script>

<main>
  <h1 class="text-center">GravityR</h1>

  {#if $grParam.query}
    <div class="details-card">
      <Details class="card">
        <span slot="summary">SQL</span>
        <QueryText query={$grParam.query} />
      </Details>
    </div>
  {/if}

  {#if $grParam.analyzeNodes}
    <div class="details-card">
      <Details class="card">
        <span slot="summary">Explain Tree</span>
        <ExplainAnalyzeText
          {highlightIndexKey}
          analyzeNodes={$grParam.analyzeNodes}
        />
      </Details>
    </div>
  {/if}

  {#if $grParam.analyzeNodes}
    <div class="details-card">
      <Details class="card" open>
        <span slot="summary">Execution Timeline</span>
        <ExplainAnalyzeChart
          {highlightIndexKey}
          analyzeNodes={$grParam.analyzeNodes}
        />
      </Details>
    </div>
  {/if}

  {#if $grParam.indexTargets}
    <div class="details-card">
      <!--<Details class="card" open={!$grParam.examinationResult}>-->
      <Details class="card" open>
        <span slot="summary">Index suggestion</span>
        <IndexSuggestion
          examinationCommandOptions={$grParam.examinationCommandOptions}
          indexTargets={$grParam.indexTargets}
        />
      </Details>
    </div>
  {/if}

  {#if $grParam.examinationResult}
    <div class="details-card">
      <Details class="card" open>
        <span slot="summary">Examination Result</span>
        <ExaminationResultTable result={$grParam.examinationResult} />
      </Details>
    </div>
  {/if}
</main>

<style>
  :global(:root) {
    --color-primary: #01082a;
  }

  h1 {
    margin-bottom: 35px;
  }

  .details-card {
    margin-bottom: 30px;
  }
</style>
