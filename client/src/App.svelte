<script lang="ts">
  import "chota";
  import "reflect-metadata";
  import SuggestPage from "./pages/suggest/SuggestPage.svelte";
  import DigPage from "./pages/dig/DigPage.svelte";
  import { grParam } from "./stores/gr_param.js";
  import { Nav } from "svelte-chota";

  const getCurrentTab = (path: string): string => {
    if ($grParam.dev) {
      if (path === "/suggest") {
        return "suggest";
      }
      if (path === "/dig") {
        return "dig";
      }
    }
    return !!$grParam.suggestData ? "suggest" : "dig";
  };

  let page = getCurrentTab(window.location.pathname);
</script>

<main>
  {#if $grParam.dev}
    <Nav>
      <svelte:fragment slot="left">
        <a href="/">Root</a>
        <a href="/suggest">Suggest</a>
        <a href="/dig">Dig</a>
      </svelte:fragment>
    </Nav>
  {/if}

  <h1 class="text-center">GravityR</h1>

  {#if page === "suggest"}
    <SuggestPage suggestData={$grParam.suggestData} />
  {:else if page === "dig"}
    <DigPage digData={$grParam.digData} />
  {/if}
</main>

<style>
  :global(:root) {
    --color-primary: #01082a;
  }

  h1 {
    margin-bottom: 35px;
  }
</style>
