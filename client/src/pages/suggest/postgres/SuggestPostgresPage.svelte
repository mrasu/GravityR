<script lang="ts">
  import { createHighlightIndexContext } from "@/contexts/HighlightIndexContext";
  import ExplainText from "@/components/ExplainText.svelte";
  import ExplainAnalyzeChart from "./ExplainAnalyzeChart.svelte";
  import IndexSuggestion from "@/components/IndexSuggestion.svelte";
  import ExaminationResultTable from "@/components/ExaminationResultTable.svelte";
  import type { PostgresSuggestData } from "@/models/PostgresSuggestData";
  import DetailsCard from "@/components/DetailsCard.svelte";
  import PlainText from "@/components/PlainText.svelte";

  export let suggestData: PostgresSuggestData;
  const highlightIndexKey = createHighlightIndexContext();
</script>

<main>
  {#if suggestData.query}
    <DetailsCard title="SQL" open testId="sql">
      <PlainText text={suggestData.query} />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Explain Tree" open testId="explain">
      <ExplainText
        {highlightIndexKey}
        analyzeNodes={suggestData.analyzeNodes}
        trailingText={suggestData.summaryText}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Execution Timeline" open testId="explainChart">
      <ExplainAnalyzeChart
        {highlightIndexKey}
        chartDescription="Execution time based timeline from EXPLAIN ANALYZE"
        analyzeNodes={suggestData.analyzeNodes}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.indexTargets}
    <DetailsCard
      title="Index suggestion"
      open={!suggestData.examinationResult}
      testId="suggest"
    >
      <IndexSuggestion
        subCommandKey="postgres"
        examinationCommandOptions={suggestData.examinationCommandOptions}
        indexTargets={suggestData.indexTargets}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.examinationResult}
    <DetailsCard title="Examination Result" open testId="examination">
      <ExaminationResultTable result={suggestData.examinationResult} />
    </DetailsCard>
  {/if}
</main>

<style>
</style>
