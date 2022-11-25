<script lang="ts">
  import "chota";
  import "reflect-metadata";
  import SuggestMysqlPage from "./pages/suggest/mysql/SuggestMysqlPage.svelte";
  import DigPerformanceInsightsPage from "./pages/dig/performance_insights/DigPerformanceInsightsPage.svelte";
  import { grParam } from "./stores/gr_param.js";
  import { Nav } from "svelte-chota";
  import SuggestPostgresPage from "./pages/suggest/postgres/SuggestPostgresPage.svelte";
  import SuggestHasuraPage from "@/pages/suggest/hasura/SuggestHasuraPage.svelte";

  const getCurrentTab = (path: string): string => {
    if ($grParam.dev) {
      if (path === "suggest_mysql") {
        return "suggest_mysql";
      }
      if (path === "suggest_postgres") {
        return "suggest_postgres";
      }
      if (path === "suggest_hasura") {
        return "suggest_hasura";
      }
      if (path === "dig") {
        return "dig";
      }
    }

    if (!!$grParam.suggestData) {
      if (!!$grParam.suggestData.mysqlSuggestData) {
        return "suggest_mysql";
      } else if (!!$grParam.suggestData.postgresSuggestData) {
        return "suggest_postgres";
      } else if (!!$grParam.suggestData.hasuraSuggestData) {
        return "suggest_hasura";
      }
    } else {
      return "dig";
    }
  };

  let page = getCurrentTab(
    new URL(window.location.href).searchParams.get("page")
  );
</script>

<main>
  {#if $grParam.dev}
    <Nav>
      <svelte:fragment slot="left">
        <a href="/">Root</a>
        <a href="/?page=suggest_mysql">Suggest (MySQL)</a>
        <a href="/?page=suggest_postgres">Suggest (Postgres)</a>
        <a href="/?page=suggest_hasura">Suggest (Hasura)</a>
        <a href="/?page=dig">Dig</a>
        <a href="/roots/mermaid.html">Mermaid</a>
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
  {:else if page === "suggest_hasura"}
    <SuggestHasuraPage suggestData={$grParam.suggestData.hasuraSuggestData} />
  {:else if page === "dig"}
    <DigPerformanceInsightsPage
      performanceInsightsData={$grParam.digData.performanceInsightsData}
    />
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
