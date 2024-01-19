import {NpDate} from "./np-date.np.js";
import {Week} from "./week.np.js";
import {Month} from "./month.np.js";
import {Column} from "./column.np.js";
import {Alignment} from "./alignment.np.js";

console.log("A simple program demonstrating enums in NanoPack.")
console.log("In NanoPack, enums are serialized into their backing types.")
console.log("The backing type can be specified in the schema. When none is specified, nanoc will determine the most appropriate type.")

const date = new NpDate(19, Week.MONDAY, Month.JUNE, 2000)
const dateBytes = date.bytes()
console.log("The Date message uses the Week and Month enum, both of which are backed by an int8.")
console.log("Raw bytes of Date: ", [...dateBytes])
console.log("Total bytes:", dateBytes.byteLength)
console.log("")

const column = new Column(Alignment.CENTER)
const columnBytes = column.bytes()
console.log("The Column message uses the Alignment enum which is backed by a string.")
console.log("Raw bytes of Column: ", [...columnBytes])
console.log("Total bytes:", columnBytes.byteLength)
