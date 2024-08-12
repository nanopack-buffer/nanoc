// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import Foundation
import NanoPack

let SimpleMessage_typeID: TypeID = 3_338_766_369

class SimpleMessage: NanoPackMessage {
  var typeID: TypeID { return 3_338_766_369 }

  var headerSize: Int { return 32 }

  let stringField: String
  let intField: Int32
  let doubleField: Double
  let optionalField: String?
  let arrayField: [UInt8]
  let mapField: [String: Bool]
  let anyMessage: NanoPackMessage

  init(
    stringField: String, intField: Int32, doubleField: Double, optionalField: String?,
    arrayField: [UInt8], mapField: [String: Bool], anyMessage: NanoPackMessage
  ) {
    self.stringField = stringField
    self.intField = intField
    self.doubleField = doubleField
    self.optionalField = optionalField
    self.arrayField = arrayField
    self.mapField = mapField
    self.anyMessage = anyMessage
  }

  required init?(data: Data) {
    var ptr = data.startIndex + 32

    let stringFieldSize = data.readSize(ofField: 0)
    guard let stringField = data.read(at: ptr, withLength: stringFieldSize) else {
      return nil
    }
    ptr += stringFieldSize

    let intField: Int32 = data.read(at: ptr)
    ptr += 4

    let doubleField: Double = data.read(at: ptr)
    ptr += 8

    var optionalField: String?
    if data.readSize(ofField: 3) < 0 {
      optionalField = nil
    } else {
      let optionalFieldSize = data.readSize(ofField: 3)
      guard let optionalField_ = data.read(at: ptr, withLength: optionalFieldSize) else {
        return nil
      }
      optionalField = optionalField_
      ptr += optionalFieldSize
    }

    let arrayFieldByteSize = data.readSize(ofField: 4)
    let arrayFieldItemCount = arrayFieldByteSize / 1
    let arrayField = data[ptr..<ptr + arrayFieldByteSize].withUnsafeBytes {
      [UInt8]($0.bindMemory(to: UInt8.self).lazy.map { $0.littleEndian })
    }
    ptr += arrayFieldByteSize

    let mapFieldItemCount = data.readSize(at: ptr)
    ptr += 4
    var mapField: [String: Bool] = [:]
    mapField.reserveCapacity(mapFieldItemCount)
    for i in 0..<mapFieldItemCount {
      let iKeySize = data.readSize(at: ptr)
      ptr += 4
      guard let iKey = data.read(at: ptr, withLength: iKeySize) else {
        return nil
      }
      ptr += iKeySize
      let iValue: Bool = data.read(at: ptr)
      ptr += 1
      mapField[iKey] = iValue
    }

    let anyMessageByteSize = data.readSize(ofField: 6)
    guard let anyMessage = makeNanoPackMessage(from: data[ptr...]) else {
      return nil
    }
    ptr += anyMessageByteSize

    self.stringField = stringField
    self.intField = intField
    self.doubleField = doubleField
    self.optionalField = optionalField
    self.arrayField = arrayField
    self.mapField = mapField
    self.anyMessage = anyMessage
  }

  required init?(data: Data, bytesRead: inout Int) {
    var ptr = data.startIndex + 32

    let stringFieldSize = data.readSize(ofField: 0)
    guard let stringField = data.read(at: ptr, withLength: stringFieldSize) else {
      return nil
    }
    ptr += stringFieldSize

    let intField: Int32 = data.read(at: ptr)
    ptr += 4

    let doubleField: Double = data.read(at: ptr)
    ptr += 8

    var optionalField: String?
    if data.readSize(ofField: 3) < 0 {
      optionalField = nil
    } else {
      let optionalFieldSize = data.readSize(ofField: 3)
      guard let optionalField_ = data.read(at: ptr, withLength: optionalFieldSize) else {
        return nil
      }
      optionalField = optionalField_
      ptr += optionalFieldSize
    }

    let arrayFieldByteSize = data.readSize(ofField: 4)
    let arrayFieldItemCount = arrayFieldByteSize / 1
    let arrayField = data[ptr..<ptr + arrayFieldByteSize].withUnsafeBytes {
      [UInt8]($0.bindMemory(to: UInt8.self).lazy.map { $0.littleEndian })
    }
    ptr += arrayFieldByteSize

    let mapFieldItemCount = data.readSize(at: ptr)
    ptr += 4
    var mapField: [String: Bool] = [:]
    mapField.reserveCapacity(mapFieldItemCount)
    for i in 0..<mapFieldItemCount {
      let iKeySize = data.readSize(at: ptr)
      ptr += 4
      guard let iKey = data.read(at: ptr, withLength: iKeySize) else {
        return nil
      }
      ptr += iKeySize
      let iValue: Bool = data.read(at: ptr)
      ptr += 1
      mapField[iKey] = iValue
    }

    let anyMessageByteSize = data.readSize(ofField: 6)
    guard let anyMessage = makeNanoPackMessage(from: data[ptr...]) else {
      return nil
    }
    ptr += anyMessageByteSize

    self.stringField = stringField
    self.intField = intField
    self.doubleField = doubleField
    self.optionalField = optionalField
    self.arrayField = arrayField
    self.mapField = mapField
    self.anyMessage = anyMessage

    bytesRead = ptr - data.startIndex
  }

  func write(to data: inout Data, offset: Int) -> Int {
    let dataCountBefore = data.count

    data.reserveCapacity(offset + 32)

    data.append(typeID: TypeID(SimpleMessage_typeID))
    data.append([0], count: 7 * 4)

    data.write(size: stringField.lengthOfBytes(using: .utf8), ofField: 0, offset: offset)
    data.append(string: stringField)

    data.write(size: 4, ofField: 1, offset: offset)
    data.append(int: intField)

    data.write(size: 8, ofField: 2, offset: offset)
    data.append(double: doubleField)

    if let optionalField = self.optionalField {
      data.write(size: optionalField.lengthOfBytes(using: .utf8), ofField: 3, offset: offset)
      data.append(string: optionalField)
    } else {
      data.write(size: -1, ofField: 3, offset: offset)
    }

    data.write(size: arrayField.count * 1, ofField: 4, offset: offset)
    for i in arrayField {
      data.append(int: i)
    }

    data.append(size: mapField.count)
    var mapFieldByteSize = 4 + 1
    for (iKey, iValue) in mapField {
      data.append(size: iKey.lengthOfBytes(using: .utf8))
      data.append(string: iKey)
      data.append(bool: iValue)
      mapFieldByteSize += 1
    }
    data.write(size: mapFieldByteSize, ofField: 5, offset: offset)

    let anyMessageByteSize = anyMessage.write(to: &data, offset: data.count)
    data.write(size: anyMessageByteSize, ofField: 6, offset: offset)

    return data.count - dataCountBefore
  }

  func data() -> Data? {
    var data = Data()
    _ = write(to: &data, offset: 0)
    return data
  }
}