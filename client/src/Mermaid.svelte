<script lang="ts">
  import "chota";
  import "reflect-metadata";
  import { Nav } from "svelte-chota";
  import { grParam } from "@/stores/gr_param.js";
  import DigJaegerPage from "@/pages/dig/jaeger/DigJaegerPage.svelte";

  const getCurrentTab = (path: string): string => {
    if ($grParam.dev) {
      if (path === "dig_jaeger") {
        return "dig_jaeger";
      }
    }

    return "dig_jaeger";
  };

  let page = getCurrentTab(
    new URL(window.location.href).searchParams.get("page")
  );
</script>

<main>
  {#if $grParam.dev}
    <Nav>
      <svelte:fragment slot="left">
        <a href="/roots/mermaid.html">Root</a>
        <a href="/roots/mermaid.html?page=dig_jaeger">Dig (Jaeger)</a>
        <a href="/">Apexcharts</a>
      </svelte:fragment>
    </Nav>
  {/if}

  <h1 class="text-center">GravityR</h1>

  {#if page === "dig_jaeger"}
    <DigJaegerPage jaegerData={$grParam.digData.jaegerData} />
  {/if}
</main>
