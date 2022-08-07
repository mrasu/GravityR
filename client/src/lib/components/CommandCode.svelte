<script lang="ts">
  export let commandText: string;
  export let height: number;
</script>

<div class="command-code" style:height={`${height}px`}>
  <code>{commandText}</code>
</div>

<style>
  .command-code {
    overflow: auto;
    background-color: var(--bg-secondary-color);
  }
</style>
