# LPFR-O-MATIC

Simple watchdog for lpfr server process

## Usage:

```lpfr-o-matic [flags]```

Flags:

```
--checkurl string           url of the LPFR to check (default "http://localhost:7555")
--config string             config file (default is ./.lpfr-o-matic.yaml)
--exepath string            LPFR executable name / shortcut (default "lpfr.lnk")
-h, --help                      help for lpfr-o-matic
--interval int              interval (in seconds) to perform checks (default 10)
--middleware string         middleware application to start on successful status
--nopin                     skip automatic pin setup
--pin string                smart card pin code
--telegram                  send telegram status messages
--telegram-api-key string   telegram bot api key
--telegram-chat-id string   telegram chat id
--telegram-sender string    sender identification
-v, --version                   version for lpfr-o-matic
```

Settings may be written to configuration file named ```.lpfr-o-matic.yaml``` (in the same folder as lpfr-o-matic 
executable). Example:

```yaml
exepath: D:\project\source\repositories\fake-lpfr\fake-lpfr.exe
checkurl: http://localhost:29000
pin: 1234
telegram: true
telegram-channel-id: ......
telegram-api-key: ....
telegram-sender: sender string here
middleware: some-middleware.lnk
interval: 20
```