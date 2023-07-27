# A Peer to Peer example with gRPC

In [main.go](main.go) you can find a simple example of how to set up a peer-to-peer connection and how they send messages to each other.

The [Proto file](proto/Hello.proto) defines a simple grp method `Hello` that takes a HelloMessage and returns a HelloMessage. The `HelloMessage` contains a `message` and a `sender` of the message.
