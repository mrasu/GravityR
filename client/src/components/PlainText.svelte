<script lang="ts">
  export let text: string;
</script>

<div>{text.trim()}</div>

<style>
  div {
    white-space: pre;
    overflow-x: auto;
    padding-top: 5px;
    padding-bottom: 5px;
  }
</style>
