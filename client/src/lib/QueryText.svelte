<script lang="ts">
  export let query: string;
</script>

<div>{query.trim()}</div>

<style>
  div {
    white-space: pre;
    overflow-x: auto;
    padding-top: 5px;
    padding-bottom: 5px;
  }
</style>
