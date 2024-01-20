// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, NanoBufWriter, type NanoPackMessage } from "nanopack";

import type { TAlignment } from "./alignment.np.js";

class Column implements NanoPackMessage {
  public static TYPE_ID = 2;

  constructor(public alignment: TAlignment) {}

  public static fromBytes(
    bytes: Uint8Array,
  ): { bytesRead: number; result: Column } | null {
    const reader = new NanoBufReader(bytes);
    return Column.fromReader(reader);
  }

  public static fromReader(
    reader: NanoBufReader,
  ): { bytesRead: number; result: Column } | null {
    let ptr = 8;

    const alignmentByteLength = reader.readFieldSize(0);
    const alignment = reader.readString(ptr, alignmentByteLength) as TAlignment;
    ptr += alignmentByteLength;

    return { bytesRead: ptr, result: new Column(alignment) };
  }

  public get typeId(): number {
    return 2;
  }

  public bytes(): Uint8Array {
    const writer = new NanoBufWriter(8);
    writer.writeTypeId(2);

    const alignmentByteLength = writer.appendString(this.alignment);
    writer.writeFieldSize(0, alignmentByteLength);

    return writer.bytes;
  }

  public bytesWithLengthPrefix(): Uint8Array {
    const writer = new NanoBufWriter(8 + 4, true);
    writer.writeTypeId(2);

    const alignmentByteLength = writer.appendString(this.alignment);
    writer.writeFieldSize(0, alignmentByteLength);

    writer.writeLengthPrefix(writer.currentSize - 4);

    return writer.bytes;
  }
}

export { Column };