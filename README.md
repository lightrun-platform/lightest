# Lightest
Lightest is a tester application used to verify that a customer environment's prerequisites for using Lightrun are met.

Currently, the test application supports the following tests:
- **Long polling test** - which tests that the test application can register as a Lightrun agent and perform `getBreakpoints` long polling requests for a period of time.
- **Websocket connectivity test** - which tests that the test application can authenticate as a Lightrun client, establish a websocket connection and maintain its connectivity for a period of time.

## Usage

Since Go supports cross-compiling and comes with a built-in compiler support for multiple OSs and architectures, building a OS and architecture specific executables is simple.
The [Github release page](https://github.com/lightrun-platform/lightest/releases) should include the latest distributions for a combination of OSs and architectures we would like to support.
Running the tester simply requires running the relevant distributed executable.  
e.g., for running the linux, 32 bit executable, simply run:
```
./lightest-linux-x86
```
You will most likely need to give the executable execution permissions:
```
chmod +x lightest-linux-x86
```

### Configuration

Before running the executable, the accompanied `config.json` file should be configured according to your needs and reside in the same folder as the executable's.
Configuring it is straightforward, as the necessary configuration file parameters written in the config should be relatively self-explanatory.  
Nonetheless, note that the following fields must be filled:  
* `agent.apiKey` - agent-polling test prerequisite.
* `userEmail` - websocket-connection test prerequisite.
* `userPassword` - websocket-connection test prerequisite
* `companyId` - agent-polling and websocket-connection tests prerequisites.


## Build

It should be quite easy building the Lightest executables on your own.  
The following is a bash function that can be used to build Lightest for the currently supported set of platforms and architectures:

```bash
function crossBuildGo() {
    GOOS=windows GOARCH=amd64 go build -o lightest-windows-amd64
    GOOS=windows GOARCH=386 go build -o lightest-windows-x86
    GOOS=linux GOARCH=arm64 go build -o lightest-linux-arm64
    GOOS=linux GOARCH=386 go build -o lightest-linux-x86
    GOOS=linux GOARCH=amd64 go build -o lightest-linux-amd64
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o lightest-alpine-arm64
    GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o lightest-alpine-x86
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lightest-alpine-amd64
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o lightest-darwin-amd64
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o lightest-darwin-arm64
}
```
