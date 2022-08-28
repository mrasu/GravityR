import { IndexTarget } from "@/models/IndexTarget";

export class PostgresIndexTarget extends IndexTarget {
  toAlterAddSQL(): string {
    const columns = this.columns.map((v) => v.name).join(", ");

    return `CREATE INDEX ON ${this.tableName} (${columns});`;
  }
}
