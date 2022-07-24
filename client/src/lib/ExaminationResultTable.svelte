<script lang="ts">
  import type { ExaminationResult } from "../models/ExaminationResult";

  export let result: ExaminationResult;
  $: sortedIndexResults = [...result.indexResults].sort(
    (a, b) => a.executionTimeMillis - b.executionTimeMillis
  );
</script>

<div>Times after adding index</div>

<table class="striped">
  <thead>
    <tr>
      <th class="time-ms"><div>Execution Time</div></th>
      <th class="reduction"><div>Reduction</div></th>
      <th class="index"><div>Columns</div></th>
      <th class="index"><div>SQL</div></th>
    </tr>
    <tr>
      <td class="time-ms"
        ><div>
          {result.originalTimeMillis.toLocaleString()}ms
        </div></td
      >
      <td class="reduction"><div>-</div></td>
      <td class="index"><div>(Original)</div></td>
      <td class="index" />
    </tr>

    {#each sortedIndexResults as indexResult}
      {@const reducedPercent = indexResult.toReductionPercent(
        result.originalTimeMillis
      )}
      <tr>
        <td class="time-ms">
          <div>{indexResult.executionTimeMillis.toLocaleString()}ms</div></td
        >
        <td class="reduction"><div>{reducedPercent}%</div></td>
        <td class="index">{indexResult.toIndex()}</td>
        <td class="index">{indexResult.indexTarget.toAlterAddSQL()}</td>
      </tr>
    {/each}
  </thead>
</table>

<style>
  table {
    table-layout: fixed;
    margin-bottom: 20px;
    padding: 10px;
  }
  .time-ms {
    width: 150px;
    text-align: end;
  }
  .time-ms > div {
    margin-right: 20px;
  }
  .reduction {
    width: 100px;
    text-align: end;
  }
  .reduction > div {
    margin-right: 20px;
  }
  .index {
    width: 100%;
  }
</style>
