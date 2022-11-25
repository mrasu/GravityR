<script lang="ts">
  import TracesTable from "@/pages/dig/jaeger/TracesTable.svelte";
  import type { JaegerData } from "@/models/JaegerData";
  import type { TraceTree } from "@/models/trace_diagram/TraceTree";
  import type { OtelCompactedTrace } from "@/models/OtelCompactedTrace";
  import DetailsCard from "@/components/DetailsCard.svelte";

  export let jaegerData: JaegerData;

  const toTraceTree = (traces: OtelCompactedTrace[]): TraceTree[] => {
    return traces.map((t) => t.toTraceTree());
  };

  $: slowTraceTrees = toTraceTree(jaegerData.slowTraces);
  $: sameServiceTraceTrees = toTraceTree(jaegerData.sameServiceTraces);
</script>

<DetailsCard title="Slow traces" open testId="slowTraces">
  <TracesTable
    uiPath={jaegerData.uiPath}
    traceTrees={slowTraceTrees}
    description={`traces took more than ${Intl.NumberFormat().format(
      jaegerData.slowThresholdMilli
    )} ms`}
  />
</DetailsCard>

<DetailsCard title="Same service traces" open testId="sameServiceTraces">
  <TracesTable
    uiPath={jaegerData.uiPath}
    traceTrees={sameServiceTraceTrees}
    description={`Traces accessed the same service more than ${
      jaegerData.sameServiceThreshold
    } times and took more than ${Intl.NumberFormat().format(
      jaegerData.slowThresholdMilli
    )} ms`}
  />
</DetailsCard>
