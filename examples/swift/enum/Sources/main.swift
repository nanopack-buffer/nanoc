print("A simple program demonstrating enums in NanoPack.")
print("In NanoPack, enums are serialized into their backing types.")
print("The backing type can be specified in the schema. When none is specified, nanoc will determine the most appropriate type.")
print("")

let date = NpDate(day: 19, week: .monday, month: .june, year: 2000)
let dateData = date.data()!
print("The Date message uses the Week and Month enum, both of which are backed by an int8.")
print("Raw bytes of Date: ", terminator: "")
for b in dateData {
    print("\(b)", terminator: " ")
}
print("")
print("Total bytes:", dateData.count)
print("")

let column = Column(alignment: .center)
let columnData = column.data()!
print("The Column message uses the Alignment enum which is backed by a string.")
print("Raw bytes of Date: ", terminator: "")
for b in columnData {
    print("\(b)", terminator: " ")
}
print("")
print("Total bytes:", columnData.count)
