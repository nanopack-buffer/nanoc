import NanoPack

print("a simple program demonstrating conversion between NanoPack data & Swift struct.")

let message = SimpleMessage(
    stringField: "hello world",
    intField: 123456,
    doubleField: 123.456,
    optionalField: nil,
    arrayField: [1, 2, 3],
    mapField: ["hello": true],
    anyMessage: NestedMessage(stringField: "nested")
)

let data = message.data()!
print("raw bytes: ", terminator: "")
for b in data {
    print("\(b)", terminator: " ")
}

print("")
print("total bytes:", data.count)

let decoded = SimpleMessage(data: data)
print(decoded?.mapField)
