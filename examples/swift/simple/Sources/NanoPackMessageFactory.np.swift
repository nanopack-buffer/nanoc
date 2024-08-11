// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import Foundation
import NanoPack

func makeNanoPackMessage(from data: Data) -> NanoPackMessage? {
    let typeID = data.readTypeID()
    switch typeID {
    case 2309634176: return NestedMessage(data: data)
    case 3338766369: return SimpleMessage(data: data)
    default: return nil
    }
}

func makeNanoPackMessage(from data: Data, bytesRead: inout Int) -> NanoPackMessage? {
    let typeID = data.readTypeID()
    switch typeID {
    case 2309634176: return NestedMessage(data: data, bytesRead: &bytesRead)
    case 3338766369: return SimpleMessage(data: data, bytesRead: &bytesRead)
    default: return nil
    }
}
