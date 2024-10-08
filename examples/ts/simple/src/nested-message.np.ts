// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, NanoBufWriter, type NanoPackMessage } from "nanopack";

class NestedMessage implements NanoPackMessage {
  public static TYPE_ID = 2309634176;

  public readonly typeId: number = 2309634176;

  public readonly headerSize: number = 8;

  constructor(public stringField: string) {}

  public static fromBytes(
    bytes: Uint8Array,
  ): { bytesRead: number; result: NestedMessage } | null {
    const reader = new NanoBufReader(bytes);
    return NestedMessage.fromReader(reader);
  }

  public static fromReader(
    reader: NanoBufReader,
    offset = 0,
  ): { bytesRead: number; result: NestedMessage } | null {
    let ptr = offset + 8;

    const stringFieldByteLength = reader.readFieldSize(0, offset);
    const stringField = reader.readString(ptr, stringFieldByteLength);
    ptr += stringFieldByteLength;

    return { bytesRead: ptr - offset, result: new NestedMessage(stringField) };
  }

  public writeTo(writer: NanoBufWriter, offset = 0): number {
    let bytesWritten = 8;

    writer.writeTypeId(2309634176, offset);

    const stringFieldByteLength = writer.appendString(this.stringField);
    writer.writeFieldSize(0, stringFieldByteLength, offset);
    bytesWritten += stringFieldByteLength;

    return bytesWritten;
  }

  public bytes(): Uint8Array {
    const writer = new NanoBufWriter(8);
    this.writeTo(writer);
    return writer.bytes;
  }
}

export { NestedMessage };
