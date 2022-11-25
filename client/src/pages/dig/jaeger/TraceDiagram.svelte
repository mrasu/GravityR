<script lang="ts">
  import mermaid from "mermaid";
  import { onMount } from "svelte";
  import type { TraceTree } from "@/models/trace_diagram/TraceTree";

  export let trace: TraceTree;

  const rerender = () => {
    const diagram = trace.toMermaidDiagram();
    mermaid.render("mermaidSVG", diagram, (code) => {
      mermaidDiv.innerHTML = code;
    });
  };

  let mermaidDiv: HTMLElement;

  onMount(() => {
    rerender();
  });
</script>

<div class="container">
  <h2>TraceId: {trace.traceId}</h2>
  <h2>Open detail</h2>
  <div bind:this={mermaidDiv} />
</div>

<style>
</style>
