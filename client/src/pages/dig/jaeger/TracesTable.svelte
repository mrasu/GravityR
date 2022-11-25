<script lang="ts">
  import TraceDiagram from "@/pages/dig/jaeger/TraceDiagram.svelte";
  import type { TraceTree } from "@/models/trace_diagram/TraceTree";
  import { Modal } from "svelte-chota";

  export let uiPath: string;
  export let traceTrees: TraceTree[];
  export let description: string;

  $: sortedTraceTrees = traceTrees.sort(
    (a, b) => b.durationMilli - a.durationMilli
  );

  let showingTrace: TraceTree;
  const showDiagram = (trace: TraceTree) => {
    showingTrace = trace;
    showModal = !!trace;
  };

  let showModal = false;
</script>

<div>{description}</div>

<table class="striped">
  <thead>
    <tr>
      <th class="time-ms"><div>Trace ID</div></th>
      <th class="reduction"><div>Duration</div></th>
      <th class="index"><div># of same service</div></th>
      <th class="index"><div>Time consuming service</div></th>
      <th class="index" />
      <th class="index" />
    </tr>
  </thead>
  {#if sortedTraceTrees.length === 0}
    <tr>
      <td class="no-trace" colspan="6">No trace found</td>
    </tr>
  {:else}
    {#each sortedTraceTrees as traceTree}
      <tr>
        <td>{traceTree.traceId}</td>
        <td class="right-align"
          >{Intl.NumberFormat().format(traceTree.durationMilli)}ms</td
        >
        <td class="center-align"
          >{Intl.NumberFormat().format(traceTree.sameServiceAccessCount)} times</td
        >
        <td>{traceTree.timeConsumingServiceName}</td>
        <td>
          <button
            on:click={() => {
              showDiagram(traceTree);
            }}>Show diagram</button
          >
        </td>
        <td
          ><a
            href={new URL(`/trace/${traceTree.traceId}`, uiPath).href}
            target="_blank"
            rel="noopener noreferrer">Open detail</a
          ></td
        >
      </tr>
    {/each}
  {/if}
</table>

<div class="trace-modal">
  <Modal bind:open={showModal}>
    <TraceDiagram trace={showingTrace} />
  </Modal>
</div>

<style>
  table {
    margin-bottom: 20px;
    padding: 10px;
  }

  .trace-modal :global(.modal) {
    height: 80%;
    width: 80%;
    overflow: auto;
  }
  td.right-align {
    text-align: right;
    padding-right: 10px;
  }
  td.center-align {
    text-align: center;
  }
</style>
