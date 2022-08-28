<script lang="ts">
  import { createHighlightIndexContext } from "@/contexts/HighlightIndexContext";
  import QueryText from "@/components/QueryText.svelte";
  import ExplainAnalyzeText from "@/components/ExplainAnalyzeText.svelte";
  import ExplainAnalyzeChart from "./ExplainAnalyzeChart.svelte";
  import IndexSuggestion from "@/components/IndexSuggestion.svelte";
  import ExaminationResultTable from "@/components/ExaminationResultTable.svelte";
  import type { MysqlSuggestData } from "@/models/MysqlSuggestData";
  import DetailsCard from "@/components/DetailsCard.svelte";

  export let suggestData: MysqlSuggestData;
  const highlightIndexKey = createHighlightIndexContext();
</script>

<main>
  {#if suggestData.query}
    <DetailsCard title="SQL">
      <QueryText query={suggestData.query} />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Explain Tree" open>
      <ExplainAnalyzeText
        {highlightIndexKey}
        analyzeNodes={suggestData.analyzeNodes}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Execution Timeline" open>
      <ExplainAnalyzeChart
        {highlightIndexKey}
        analyzeNodes={suggestData.analyzeNodes}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.indexTargets}
    <DetailsCard title="Index suggestion" open={!suggestData.examinationResult}>
      <IndexSuggestion
        subCommand="mysql"
        examinationCommandOptions={suggestData.examinationCommandOptions}
        indexTargets={suggestData.indexTargets}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.examinationResult}
    <DetailsCard title="Examination Result" open>
      <ExaminationResultTable result={suggestData.examinationResult} />
    </DetailsCard>
  {/if}
</main>

<style>
</style>
