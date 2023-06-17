# MQTT Connectivity Testing Tool

This application is an MQTT Testing Tool for testing MQTT connections and sending/receiving messages from topics.

The tool is created for the purpose of testing connectivity to our MQTT broker with certificate-key pair that we provided.
And it was created since not everyone has any (easy) way of testing connectivity with MQTT protocol.

Tool itself is very low-effort console application. It does unstructured logging to console for debugging purposes.

## Features

- Connect to an MQTT broker using TLS with client certificate authentication
- Subscribe to one or multiple MQTT topics
- Display received messages in the console
- Gracefully reconnects and resubscribes to topics upon connection loss
- Listens for a TCP connections so one can connect to it  and try publishing.

## Dependencies

For MQTT connections the tool uses [eclipse/paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang)
which is under license: Eclipse Public License - v 2.0 (EPL-2.0)

## Installation

You can download the pre-compiled binaries from the GitHub [releases](https://github.com/Kapacity-io/mqtt-test-client/releases) page or compile the tool yourself using the provided Makefile.

After downloading the binary, you may need to give it "execute" permissions.
If you're on a Unix-like system (such as Linux or MacOS), 
you can do this by navigating to the directory containing the downloaded file and running:

```sh
chmod +x <your-binary-file>
```

To compile from source, ensure you have a working Go environment and `make` installed, and then run `make` from the project root:

```bash
make
```

Or by default with Go tools from the project root:
```
go mod download

go build -o mqtt-test-client main.go

```

## Usage

```
./mqtt-test-client --mqtt-broker-addr=<broker URL> \
    --mqtt-broker-client-id=<client ID> \
    --mqtt-broker-device-cert=<path to certificate> \
    --mqtt-broker-device-key=<path to private key> \
    --mqtt-broker-root-ca=<path to root CA> \
    --mqtt-topics-to-sub=<comma-separated list of topics>
```

Flags:
-    -mqtt-broker-addr: URL of the MQTT broker, for example: "tcps://foo-ats.iot.eu-north-1.amazonaws.com:8883/mqtt"
-    -mqtt-broker-client-id: Client ID used for connecting to the broker (default is "mqtt-test-client")
-    -mqtt-broker-device-cert: Filepath to the certificate file
-    -mqtt-broker-device-key: Filepath to the private key for the certificate
-    -mqtt-broker-root-ca: Filepath to the root CA file, for example "AmazonRootCA1.pem"
-    -mqtt-topics-to-sub: Comma-separated list of topics to subscribe to

---

To test publish one can connect to a running applications port 9596 with a tool like netcat or telnet
for example with netcat:

```
$ nc localhost 9596
```

after connection is made, type the message as <topic-to-pub>|<message-to-publish>
for example. Newline (enter) will send the message.

```
my/test/topic|{"foo": "bar", "quz": 22.13}
```

the tool will split the sent line as topic + message and send it to a broker.


## License

This project is licensed under the MIT License - see the LICENSE file for details
