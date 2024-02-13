// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import Foundation
import NanoPack

let Widget_typeID: TypeID = 1

class Widget: NanoPackMessage {
  var typeID: TypeID { return 1 }

  let id: Int32

  static func from(data: Data) -> Widget? {
    switch data.readTypeID() {
    case 1: return Widget(data: data)
    case 2: return Text(data: data)
    default: return nil
    }
  }

  static func from(data: Data, bytesRead: inout Int) -> Widget? {
    switch data.readTypeID() {
    case 1: return Widget(data: data, bytesRead: &bytesRead)
    case 2: return Text(data: data, bytesRead: &bytesRead)
    default: return nil
    }
  }

  init(id: Int32) {
    self.id = id
  }

  required init?(data: Data) {
    var ptr = data.startIndex + 8

    let id: Int32 = data.read(at: ptr)
    ptr += 4

    self.id = id
  }

  required init?(data: Data, bytesRead: inout Int) {
    var ptr = data.startIndex + 8

    let id: Int32 = data.read(at: ptr)
    ptr += 4

    self.id = id

    bytesRead = ptr - data.startIndex
  }

  func data() -> Data? {
    var data = Data()
    data.reserveCapacity(8)

    withUnsafeBytes(of: Int32(Widget_typeID)) {
      data.append(contentsOf: $0)
    }

    data.append([0], count: 1 * 4)

    data.write(size: 4, ofField: 0)
    data.append(int: id)

    return data
  }
}