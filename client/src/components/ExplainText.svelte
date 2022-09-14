<script lang="ts">
  import type { IDbAnalyzeData } from "@/types/gr_param";
  import { getHighlightIndex } from "@/contexts/HighlightIndexContext";

  export let analyzeNodes: IDbAnalyzeData[];
  export let trailingText = "";
  export let highlightIndexKey: symbol;
  let highlightIndex = getHighlightIndex(highlightIndexKey);

  $: texts = (() => {
    const res = [];
    const visited = [];
    for (const data of analyzeNodes) {
      let current = data;
      const stack = [data];
      while (current) {
        if (!visited.includes(current)) {
          if (current.children) {
            for (let i = current.children.length - 1; i >= 0; i--) {
              stack.push(current.children[i]);
            }
          }
          res.push(current);
        }
        visited.push(current);
        current = stack.pop();
      }
    }
    return res.map((v) => v.text);
  })();
</script>

<div>
  {#each texts as text, i}
    <span class:bold={$highlightIndex === i}>{`${text}`}</span><br />
  {/each}
  {#if trailingText.length > 0}
    <span>-----------</span><br />
    <span>{trailingText}</span>
  {/if}
</div>

<style>
  div {
    white-space: pre;
    overflow-x: auto;
    padding-top: 5px;
    padding-bottom: 5px;
  }

  .bold {
    font-weight: bold;
  }
</style>
