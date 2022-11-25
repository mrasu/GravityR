import "./app.css";
import Mermaid from "./Mermaid.svelte";
import mermaid from "mermaid";

mermaid.mermaidAPI.initialize({
  startOnLoad: false,
  sequence: { messageAlign: "right" },
});

const app = new Mermaid({
  target: document.getElementById("app"),
});

export default app;
