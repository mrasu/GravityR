import { IndexColumn } from "./IndexColumn";
import { Type } from "class-transformer";

export class IndexTarget {
  tableName: string;
  @Type(() => IndexColumn)
  columns: IndexColumn[];

  toString(): string {
    const columns = this.columns.map((v) => v.name).join(", ");
    return `Table=${this.tableName} / Column=${columns}`;
  }

  toGrIndexOption(): string {
    return `"${this.tableName}:${this.toGrColumnOption()}"`;
  }

  toAlterAddSQL(): string {
    const columns = this.columns.map((v) => v.name).join(", ");

    return `ALTER TABLE ${this.tableName} ADD INDEX (${columns});`;
  }

  toExecutableText(): string {
    return this.toAlterAddSQL();
  }

  private toGrColumnOption(): string {
    return this.columns.map((v) => v.name).join("+");
  }
}
