import { SimpleMessage } from "./simple-message.np.js";

console.log(
	"a simple program demonstrating conversion between NanoPack data & TypeScript class.",
);

const message = new SimpleMessage(
	"hello world",
	123,
	123.456,
	null,
	[1, 2, 3],
	new Map([
		["mai", true],
		["sakurajima", true],
	]),
);

const bytes = message.bytes();
console.log("raw bytes: ", [...bytes]);
console.log("");
console.log("total bytes:", bytes.byteLength);

const { result: decoded } = SimpleMessage.fromBytes(bytes)!;
console.log("string field:", decoded.stringField);
console.log("int field:", decoded.intField);
console.log("double field:", decoded.doubleField);
console.log("array field:", decoded.arrayField);
decoded.mapField.forEach((v, k) => {
	console.log(`map entry: ${k} => ${v}`);
});
