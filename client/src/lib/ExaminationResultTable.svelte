<script lang="ts">
  import type { ExaminationResult } from "../models/ExaminationResult";

  export let result: ExaminationResult;
</script>

<table class="striped">
  <thead>
    <tr>
      <th class="time-ms"><div>Execution Time</div></th>
      <th class="efficiency"><div>Efficiency</div></th>
      <th class="index"><div>Index</div></th>
    </tr>
    <tr>
      <td class="time-ms"
        ><div>
          {result.originalTimeMillis.toLocaleString()}ms
        </div></td
      >
      <td class="efficiency"><div>-</div></td>
      <td class="index"><div>Original</div></td>
    </tr>

    {#each result.indexResults as time}
      <tr>
        <td class="time-ms">{time.executionTimeMillis.toLocaleString()}ms</td>
        <td class="efficiency"
          >{time.toEfficiency(result.originalTimeMillis)}%</td
        >
        <td class="index">{time.toIndex()}</td>
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
    width: 200px;
    padding-left: 20px;
  }
  .efficiency {
    width: 200px;
  }
  .index {
    width: 100%;
  }
</style>
