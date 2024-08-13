// The Swift Programming Language
// https://docs.swift.org/swift-book
import NanoPack

let clientChannel = NPInMemoryRPCClientChannel()
let serverChannel = NPInMemoryRPCServerChannel()

clientChannel.sendTo(serverChannel)
serverChannel.replyTo(clientChannel)

let client = ExampleServiceClient(channel: clientChannel)
let server = ExampleServiceServer(channel: serverChannel, delegate: ExampleService())

client.add(3, 4) {
    print($0)
}

client.subtract(123, 23) {
    print($0)
}
