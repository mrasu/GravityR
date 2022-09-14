<script lang="ts">
  import { createHighlightIndexContext } from "@/contexts/HighlightIndexContext";
  import ExplainText from "@/components/ExplainText.svelte";
  import ExplainChart from "./ExplainChart.svelte";
  import IndexSuggestion from "@/components/IndexSuggestion.svelte";
  import ExaminationResultTable from "@/components/ExaminationResultTable.svelte";
  import DetailsCard from "@/components/DetailsCard.svelte";
  import type { HasuraSuggestData } from "@/models/HasuraSuggestData";
  import PlainText from "@/components/PlainText.svelte";

  export let suggestData: HasuraSuggestData;
  const highlightIndexKey = createHighlightIndexContext();

  $: gqlText = `${suggestData.gql}
---- variables ----
${JSON.stringify(suggestData.gqlVariables)}`;
</script>

<main>
  {#if suggestData.gql}
    <DetailsCard title="GraphQL query" open>
      <PlainText text={gqlText} />
    </DetailsCard>
  {/if}

  {#if suggestData.query}
    <DetailsCard title="SQL">
      <PlainText text={suggestData.query} />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Explain Tree" open>
      <ExplainText
        {highlightIndexKey}
        analyzeNodes={suggestData.analyzeNodes}
        trailingText={suggestData.planningText}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.analyzeNodes}
    <DetailsCard title="Estimated Timeline" open>
      <ExplainChart
        {highlightIndexKey}
        chartDescription="Estimated cost based timeline from EXPLAIN"
        analyzeNodes={suggestData.analyzeNodes}
      />
    </DetailsCard>
  {/if}

  {#if suggestData.indexTargets}
    <DetailsCard title="Index suggestion" open={!suggestData.examinationResult}>
      <IndexSuggestion
        subCommand="hasura"
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
