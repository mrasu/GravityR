<script lang="ts">
  import "chota";
  import "reflect-metadata";
  import SuggestMysqlPage from "./pages/suggest/mysql/SuggestMysqlPage.svelte";
  import DigPage from "./pages/dig/DigPage.svelte";
  import { grParam } from "./stores/gr_param.js";
  import { Nav } from "svelte-chota";
  import SuggestPostgresPage from "./pages/suggest/postgres/SuggestPostgresPage.svelte";

  const getCurrentTab = (path: string): string => {
    if ($grParam.dev) {
      if (path === "/suggest_mysql") {
        return "suggest_mysql";
      }
      if (path === "/suggest_postgres") {
        return "suggest_postgres";
      }
      if (path === "/dig") {
        return "dig";
      }
    }

    if (!!$grParam.suggestData) {
      if (!!$grParam.suggestData.mysqlSuggestData) {
        return "suggest_mysql";
      } else {
        return "suggest_postgres";
      }
    } else {
      return "dig";
    }
  };

  let page = getCurrentTab(window.location.pathname);
</script>

<main>
  {#if $grParam.dev}
    <Nav>
      <svelte:fragment slot="left">
        <a href="/">Root</a>
        <a href="/suggest_mysql">Suggest (MySQL)</a>
        <a href="/suggest_postgres">Suggest (Postgres)</a>
        <a href="/dig">Dig</a>
      </svelte:fragment>
    </Nav>
  {/if}

  <h1 class="text-center">GravityR</h1>

  {#if page === "suggest_mysql"}
    <SuggestMysqlPage suggestData={$grParam.suggestData.mysqlSuggestData} />
  {:else if page === "suggest_postgres"}
    <SuggestPostgresPage
      suggestData={$grParam.suggestData.postgresSuggestData}
    />
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
