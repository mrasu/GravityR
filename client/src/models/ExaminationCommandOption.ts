export class ExaminationCommandOption {
  name: string;
  value: string;
  isShort: boolean;

  ToCommandString(): string {
    const dash = this.isShort ? "-" : "--";

    if (this.value.match(/^\w+$/)) {
      return `${dash}${this.name} ${this.value}`;
    } else {
      return `${dash}${this.name} "${this.value}"`;
    }
  }
}
