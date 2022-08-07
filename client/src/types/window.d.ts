import type { IGrParam } from "./gr_param";

declare global {
  interface Window {
    grParam: IGrParam;
  }
}
