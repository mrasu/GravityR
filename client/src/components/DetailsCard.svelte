<script lang="ts">
  import { Details } from "svelte-chota";

  export let title: string;
  export let open = false;
  export let testId = "";
</script>

<div class="details-card">
  <Details class="card" {open} data-testid={testId || null}>
    <span slot="summary">{title}</span>
    <slot />
  </Details>
</div>

<style>
  .details-card {
    margin-bottom: 30px;
  }
</style>
