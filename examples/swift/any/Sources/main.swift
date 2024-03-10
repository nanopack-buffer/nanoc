// The Swift Programming Language
// https://docs.swift.org/swift-book

import Foundation

print("The any keyword in NanoPack allows you to store any NanoPack data type.")

let clickEvent = ClickEvent(x: 23.4, y: 12.34, timestamp: Int64(Date().timeIntervalSince1970))
let invokeCallback = InvokeCallback(handle: 123, args: clickEvent.data()!)
let invokeCallbackData = invokeCallback.data()!

print("Raw bytes of InvokeCallback: ", terminator: "")
for b in invokeCallbackData {
    print("\(b)", terminator: " ")
}
print("")
print("Total bytes:", invokeCallbackData.count)
print("===================================")

let invokeCallbackParsed = InvokeCallback(data: invokeCallbackData)!
print("callback handle: \(String(describing: invokeCallbackParsed.handle))")
let clickEventParsed = ClickEvent(data: invokeCallbackParsed.args)!
print("click event x, y: \(clickEventParsed.x), \(clickEventParsed.y)")
print("timestamp: \(clickEventParsed.timestamp)")
