// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, type NanoPackMessage } from "nanopack";

import { NestedMessage } from "./nested-message.np.js";
import { SimpleMessage } from "./simple-message.np.js";

function makeNanoPackMessage(
  reader: NanoBufReader,
  offset = 0,
): { bytesRead: number; result: NanoPackMessage } | null {
  switch (reader.readTypeId(offset)) {
    case 2309634176:
      return NestedMessage.fromReader(reader, offset);
    case 3338766369:
      return SimpleMessage.fromReader(reader, offset);
    default:
      return null;
  }
}

export { makeNanoPackMessage };
