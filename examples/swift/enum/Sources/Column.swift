// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import Foundation
import NanoPack

let Column_typeID: TypeID = 2

class Column: NanoPackMessage {
  var typeID: TypeID { return 2 }

  let alignment: Alignment

  init(alignment: Alignment) {
    self.alignment = alignment
  }

  required init?(data: Data) {
    var ptr = 8

    let alignmentSize = data.readSize(ofField: 0)
    guard let alignmentRawValue = data.read(at: ptr, withLength: alignmentSize) else {
      return nil
    }
    guard let alignment = Alignment(rawValue: alignmentRawValue) else {
      return nil
    }

    self.alignment = alignment
  }

  required init?(data: Data, bytesRead: inout Int) {
    var ptr = 8

    let alignmentSize = data.readSize(ofField: 0)
    guard let alignmentRawValue = data.read(at: ptr, withLength: alignmentSize) else {
      return nil
    }
    guard let alignment = Alignment(rawValue: alignmentRawValue) else {
      return nil
    }

    self.alignment = alignment

    bytesRead = ptr
  }

  func data() -> Data? {
    var data = Data()
    data.reserveCapacity(8)

    withUnsafeBytes(of: Int32(Column_typeID)) {
      data.append(contentsOf: $0)
    }

    data.append([0], count: 1 * 4)

    data.write(size: alignment.rawValue.lengthOfBytes(using: .utf8), ofField: 0)
    data.append(string: alignment.rawValue)

    return data
  }
}
