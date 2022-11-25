const { Resource } = require('@opentelemetry/resources');
const { HttpInstrumentation } = require ('@opentelemetry/instrumentation-http');
const { SimpleSpanProcessor } = require ("@opentelemetry/sdk-trace-base");
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node");
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { ExpressInstrumentation } = require ('@opentelemetry/instrumentation-express');
const { GraphQLInstrumentation } = require("@opentelemetry/instrumentation-graphql");
const { JaegerExporter } = require("@opentelemetry/exporter-jaeger");
const { W3CTraceContextPropagator } = require("@opentelemetry/core");
const opentelemetry = require("@opentelemetry/api");

registerInstrumentations({
  instrumentations: [
    new HttpInstrumentation(),
    new ExpressInstrumentation(),
    new GraphQLInstrumentation()
  ]
})

const provider = new NodeTracerProvider({
  resource: Resource.default().merge(new Resource({
    "service.name": "bff"
  }))
})

const jaegerExporter = new JaegerExporter({
  endpoint: "http://jaeger:14268/api/traces"
})
provider.addSpanProcessor(
  new SimpleSpanProcessor(jaegerExporter)
);

opentelemetry.propagation.setGlobalPropagator( new W3CTraceContextPropagator())
provider.register()
