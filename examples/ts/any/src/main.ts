import {ClickEvent} from "./click-event.np.js";
import {InvokeCallback} from "./invoke-callback.np.js";
import {NanoBufReader} from "nanopack";

console.log("The any keyword in NanoPack allows you to store any NanoPack data type.")

const clickEvent = new ClickEvent(6.19, 10.8, BigInt(new Date().getMilliseconds()))
const invokeCallback = new InvokeCallback(123, new NanoBufReader(clickEvent.bytes()))
const bytes = invokeCallback.bytes()
console.log("Raw bytes of Date: ", [...bytes])
console.log("Total bytes:", bytes.byteLength)
console.log("")
