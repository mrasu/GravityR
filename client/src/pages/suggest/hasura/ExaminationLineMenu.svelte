<script lang="ts">
  import { Button } from "svelte-chota";
  import type { ExaminationIndexResult } from "@/models/ExaminationIndexResult";

  export let indexResult: ExaminationIndexResult;

  const onCopyAsCurlButtonClicked = () => {
    navigator.clipboard.writeText(indexResult.indexTarget.toExecutableText());
  };
</script>

<span class="pull-right">
  <Button autoclose dropdown outline class="is-center examination-result-menu">
    <Button class="menu-button" on:click={onCopyAsCurlButtonClicked}
      >Copy as Curl</Button
    >
  </Button>
</span>

<style>
  span :global(.examination-result-menu) {
    padding-left: 10px;
    padding-right: 10px;
  }

  span :global(.examination-result-menu::after) {
    content: "•••";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) rotate(90deg);
    font-size: 12px;
    letter-spacing: 1px;
  }

  span :global(.examination-result-menu + div) {
    right: 0;
    left: auto;
    padding: 0;
  }

  span :global(.menu-button) {
    background-color: inherit;
  }
</style>
