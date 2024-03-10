// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, NanoBufWriter, type NanoPackMessage } from "nanopack";

import type { TWeek } from "./week.np.js";
import type { TMonth } from "./month.np.js";

class NpDate implements NanoPackMessage {
  public static TYPE_ID = 1732634645;

  public readonly typeId: number = 1732634645;

  public readonly headerSize: number = 20;

  constructor(
    public day: number,
    public week: TWeek,
    public month: TMonth,
    public year: number,
  ) {}

  public static fromBytes(
    bytes: Uint8Array,
  ): { bytesRead: number; result: NpDate } | null {
    const reader = new NanoBufReader(bytes);
    return NpDate.fromReader(reader);
  }

  public static fromReader(
    reader: NanoBufReader,
  ): { bytesRead: number; result: NpDate } | null {
    let ptr = 20;

    const day = reader.readInt8(ptr);
    ptr += 1;

    const week = reader.readInt8(ptr) as TWeek;
    ptr += 1;

    const month = reader.readInt8(ptr) as TMonth;
    ptr += 1;

    const year = reader.readInt32(ptr);
    ptr += 4;

    return { bytesRead: ptr, result: new NpDate(day, week, month, year) };
  }

  public writeTo(writer: NanoBufWriter, offset: number = 0): number {
    let bytesWritten = 20;

    writer.writeTypeId(1732634645, offset);

    writer.appendInt8(this.day);
    writer.writeFieldSize(0, 1, offset);
    bytesWritten += 1;

    writer.appendInt8(this.week);
    writer.writeFieldSize(1, 1, offset);
    bytesWritten += 1;

    writer.appendInt8(this.month);
    writer.writeFieldSize(2, 1, offset);
    bytesWritten += 1;

    writer.appendInt32(this.year);
    writer.writeFieldSize(3, 4, offset);
    bytesWritten += 4;

    return bytesWritten;
  }

  public bytes(): Uint8Array {
    const writer = new NanoBufWriter(20);
    this.writeTo(writer);
    return writer.bytes;
  }
}

export { NpDate };
