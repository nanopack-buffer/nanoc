// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, NanoBufWriter, type NanoPackMessage } from "nanopack";

class Person implements NanoPackMessage {
  public static TYPE_ID = 1;

  constructor(
    public firstName: string,
    public middleName: string | null,
    public lastName: string,
    public age: number,
    public otherFriend: Person | null,
  ) {}

  public static fromBytes(
    bytes: Uint8Array,
  ): { bytesRead: number; result: Person } | null {
    const reader = new NanoBufReader(bytes);
    return Person.fromReader(reader);
  }

  public static fromReader(
    reader: NanoBufReader,
  ): { bytesRead: number; result: Person } | null {
    let ptr = 24;

    const firstNameByteLength = reader.readFieldSize(0);
    const firstName = reader.readString(ptr, firstNameByteLength);
    ptr += firstNameByteLength;

    let middleName: string | null;
    if (reader.readFieldSize(1) >= 0) {
      const middleNameByteLength = reader.readFieldSize(1);
      middleName = reader.readString(ptr, middleNameByteLength);
      ptr += middleNameByteLength;
    } else {
      middleName = null;
    }

    const lastNameByteLength = reader.readFieldSize(2);
    const lastName = reader.readString(ptr, lastNameByteLength);
    ptr += lastNameByteLength;

    const age = reader.readInt8(ptr);
    ptr += 1;

    let otherFriend: Person | null;
    if (reader.readFieldSize(4) >= 0) {
      const maybeOtherFriend = Person.fromReader(reader.newReaderAt(ptr));
      if (!maybeOtherFriend) {
        return null;
      }
      otherFriend = maybeOtherFriend.result;
      ptr += maybeOtherFriend.bytesRead;
    } else {
      otherFriend = null;
    }

    return {
      bytesRead: ptr,
      result: new Person(firstName, middleName, lastName, age, otherFriend),
    };
  }

  public get typeId(): number {
    return 1;
  }

  public bytes(): Uint8Array {
    const writer = new NanoBufWriter(24);
    writer.writeTypeId(1);

    const firstNameByteLength = writer.appendString(this.firstName);
    writer.writeFieldSize(0, firstNameByteLength);

    if (this.middleName) {
      const middleNameByteLength = writer.appendString(this.middleName);
      writer.writeFieldSize(1, middleNameByteLength);
    } else {
      writer.writeFieldSize(1, -1);
    }

    const lastNameByteLength = writer.appendString(this.lastName);
    writer.writeFieldSize(2, lastNameByteLength);

    writer.appendInt8(this.age);
    writer.writeFieldSize(3, 1);

    if (this.otherFriend) {
      const otherFriendData = this.otherFriend.bytes();
      writer.appendBytes(otherFriendData);
      writer.writeFieldSize(4, otherFriendData.byteLength);
    } else {
      writer.writeFieldSize(4, -1);
    }

    return writer.bytes;
  }

  public bytesWithLengthPrefix(): Uint8Array {
    const writer = new NanoBufWriter(24 + 4, true);
    writer.writeTypeId(1);

    const firstNameByteLength = writer.appendString(this.firstName);
    writer.writeFieldSize(0, firstNameByteLength);

    if (this.middleName) {
      const middleNameByteLength = writer.appendString(this.middleName);
      writer.writeFieldSize(1, middleNameByteLength);
    } else {
      writer.writeFieldSize(1, -1);
    }

    const lastNameByteLength = writer.appendString(this.lastName);
    writer.writeFieldSize(2, lastNameByteLength);

    writer.appendInt8(this.age);
    writer.writeFieldSize(3, 1);

    if (this.otherFriend) {
      const otherFriendData = this.otherFriend.bytes();
      writer.appendBytes(otherFriendData);
      writer.writeFieldSize(4, otherFriendData.byteLength);
    } else {
      writer.writeFieldSize(4, -1);
    }

    writer.writeLengthPrefix(writer.currentSize - 4);

    return writer.bytes;
  }
}

export { Person };