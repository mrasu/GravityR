import { getContext, setContext } from "svelte";
import { writable } from "svelte/store";
import type { Writable } from "svelte/types/runtime/store";

type numberOrUndef = number | undefined;

export function createHighlightIndexContext() {
  const symbol = Symbol();
  const x = writable<numberOrUndef>(undefined);
  setContext<Writable<numberOrUndef>>(symbol, x);

  return symbol;
}

export function getHighlightIndex(symbol: symbol): Writable<numberOrUndef> {
  return getContext<Writable<numberOrUndef>>(symbol);
}
