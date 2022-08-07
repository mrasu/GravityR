<script lang="ts">
  import { createHighlightIndexContext } from "../../contexts/HighlightIndexContext";
  import { Details } from "svelte-chota";
  import QueryText from "./QueryText.svelte";
  import ExplainAnalyzeText from "./ExplainAnalyzeText.svelte";
  import ExplainAnalyzeChart from "./ExplainAnalyzeChart/ExplainAnalyzeChart.svelte";
  import IndexSuggestion from "./IndexSuggestion.svelte";
  import ExaminationResultTable from "./ExaminationResultTable.svelte";
  import type { SuggestData } from "../../models/SuggestData";

  export let suggestData: SuggestData;
  const highlightIndexKey = createHighlightIndexContext();
</script>

<main>
  {#if suggestData.query}
    <div class="details-card">
      <Details class="card">
        <span slot="summary">SQL</span>
        <QueryText query={suggestData.query} />
      </Details>
    </div>
  {/if}

  {#if suggestData.analyzeNodes}
    <div class="details-card">
      <Details class="card" open>
        <span slot="summary">Explain Tree</span>
        <ExplainAnalyzeText
          {highlightIndexKey}
          analyzeNodes={suggestData.analyzeNodes}
        />
      </Details>
    </div>
  {/if}

  {#if suggestData.analyzeNodes}
    <div class="details-card">
      <Details class="card" open>
        <span slot="summary">Execution Timeline</span>
        <ExplainAnalyzeChart
          {highlightIndexKey}
          analyzeNodes={suggestData.analyzeNodes}
        />
      </Details>
    </div>
  {/if}

  {#if suggestData.indexTargets}
    <div class="details-card">
      <Details class="card" open={!suggestData.examinationResult}>
        <span slot="summary">Index suggestion</span>
        <IndexSuggestion
          examinationCommandOptions={suggestData.examinationCommandOptions}
          indexTargets={suggestData.indexTargets}
        />
      </Details>
    </div>
  {/if}

  {#if suggestData.examinationResult}
    <div class="details-card">
      <Details class="card" open>
        <span slot="summary">Examination Result</span>
        <ExaminationResultTable result={suggestData.examinationResult} />
      </Details>
    </div>
  {/if}
</main>

<style>
  .details-card {
    margin-bottom: 30px;
  }
</style>
