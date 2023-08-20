# Lightest
Lightest is a tester application used to verify that a customer environment's prerequisites for using Lightrun are met.

Currently, the test application supports the following tests:
- **Long polling test** - which tests that the test application can register as a Lightrun agent and perform `getBreakpoints` long polling requests for a period of time.
- **Websocket connectivity test** - which tests that the test application can authenticate as a Lightrun client, establish a websocket connection and maintain its connectivity for a period of time.

## Usage

Since Go supports cross-compiling and comes with a built-in compiler support for multiple OSs and architectures, building a OS and architecture specific executables is simple.
The Github release page should include the latest distributions for a combination of OSs and architectures we would like to support.
Running the tester simply requires running the relevant distributed executable.

### Configuration

Before running the executable, the accompanied `config.json` file should be configured according to your needs.
Configuring it is straightforward, as the necessary configuration file parameters written in the config should be relatively self-explanatory.



