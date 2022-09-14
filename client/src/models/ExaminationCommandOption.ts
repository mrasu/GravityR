export class ExaminationCommandOption {
  name: string;
  value: string;
  isShort: boolean;

  ToCommandString(): string {
    const dash = this.isShort ? "-" : "--";
    const escapedValue = `${this.value.replace(/(["$`\\])/g, "\\$1")}`;

    if (escapedValue.match(/^\w+$/)) {
      return `${dash}${this.name} ${escapedValue}`;
    } else {
      return `${dash}${this.name} "${escapedValue}"`;
    }
  }
}
