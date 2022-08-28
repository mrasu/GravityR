<script lang="ts">
  import { Details } from "svelte-chota";

  export let title: string;
  export let open = false;
</script>

<div class="details-card">
  <Details class="card" {open}>
    <span slot="summary">{title}</span>
    <slot />
  </Details>
</div>

<style>
  .details-card {
    margin-bottom: 30px;
  }
</style>
