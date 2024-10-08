// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import Foundation
import NanoPack

let NestedMessage_typeID: TypeID = 2_309_634_176

class NestedMessage: NanoPackMessage {
  var typeID: TypeID { return 2_309_634_176 }

  var headerSize: Int { return 8 }

  let stringField: String

  init(stringField: String) {
    self.stringField = stringField
  }

  required init?(data: Data) {
    var ptr = data.startIndex + 8

    let stringFieldSize = data.readSize(ofField: 0)
    guard let stringField = data.read(at: ptr, withLength: stringFieldSize) else {
      return nil
    }
    ptr += stringFieldSize

    self.stringField = stringField
  }

  required init?(data: Data, bytesRead: inout Int) {
    var ptr = data.startIndex + 8

    let stringFieldSize = data.readSize(ofField: 0)
    guard let stringField = data.read(at: ptr, withLength: stringFieldSize) else {
      return nil
    }
    ptr += stringFieldSize

    self.stringField = stringField

    bytesRead = ptr - data.startIndex
  }

  func write(to data: inout Data, offset: Int) -> Int {
    let dataCountBefore = data.count

    data.reserveCapacity(offset + 8)

    data.append(typeID: TypeID(NestedMessage_typeID))
    data.append([0], count: 1 * 4)

    data.write(size: stringField.lengthOfBytes(using: .utf8), ofField: 0, offset: offset)
    data.append(string: stringField)

    return data.count - dataCountBefore
  }

  func data() -> Data? {
    var data = Data()
    _ = write(to: &data, offset: 0)
    return data
  }
}
