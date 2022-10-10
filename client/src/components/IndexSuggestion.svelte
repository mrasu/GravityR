<script lang="ts">
  import { Checkbox } from "svelte-chota";
  import type { IndexTarget } from "@/models/IndexTarget";
  import type { ExaminationCommandOption } from "@/models/ExaminationCommandOption";

  const subCommands = {
    mysql: "mysql",
    postgres: "postgres",
    hasuraPostgres: "hasura postgres",
  } as const;
  type SubCommandKey = keyof typeof subCommands;

  export let subCommandKey: SubCommandKey;
  export let examinationCommandOptions: ExaminationCommandOption[];
  export let indexTargets: IndexTarget[];
  let checked = indexTargets.map(() => true);

  // Use function index of svelte because <code> needs to remove leading space but its formatter adds spaces to align lines.
  const buildCommandText = (targets: IndexTarget[]): string => {
    let command = `gr db suggest ${subCommands[subCommandKey]} --with-examine \\\n`;
    let iOption: string;

    if (targets.length === 0) {
      iOption = "\t-i ... \\\n";
    } else {
      iOption = targets
        .map((it) => it.toGrIndexOption())
        .map((option) => `\t-i ${option} \\\n`)
        .join("");
    }

    const option = examinationCommandOptions
      .map((v) => `\t${v.ToCommandString()}`)
      .join(" \\\n");

    return command + iOption + option;
  };

  $: commandText = buildCommandText(indexTargets.filter((_, i) => checked[i]));
</script>

<div>Adding index to below columns might improve performance of the query.</div>

<ul>
  {#each indexTargets as it, index}
    <li>
      <Checkbox class="aaa" value={index} bind:checked={checked[index]}>
        {it.toString()}
      </Checkbox>
    </li>
  {/each}
</ul>

<div>
  With below command, you can examine the efficiency of indexes by temporarily
  adding them in your database.
</div>

<div class="command-code">
  <code>$ {commandText}</code>
</div>

<style>
  ul {
    margin-top: 10px;
    margin-bottom: 20px;
    padding-left: 10px;
    list-style-type: none;
  }

  .command-code {
    margin-top: 5px;
  }
</style>
