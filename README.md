# Aternos API
**Unofficial** [aternos.org](https://aternos.org/) API, initially inspired by the now deprecated [AternosAPI](https://github.com/Duerocraft/AternosAPI) python library. Also 'bypasses' cloudflare without resorting to browser automation or other stupid hacks.

## Disclaimer
This is NOT an AFK bot and never will be! Please be respectful and don't abuse an awesome free service such as Aternos. I do not take any responsibility for what you do with this solution, but I will gladly take this repository down if it results in abuse of the service.

By using this package you are automatically breaking the [TOS](https://aternos.gmbh/en/aternos/terms). In other words your account *could* be suspended at any time. Use at your own risk.

## Usage
### Library
See [examples](./examples) (easy) or the [CLI source code](./cmd) (advanced).

### CLI
This project also comes with a simple command line application to start and stop your server.

Installation & usage:
```
$ go install github.com/sleeyax/aternos-api
$ aternos-api -h
```

## Projects
Projects that are using this package:
* [sleeyax/aternos-discord-bot](https://github.com/sleeyax/aternos-discord-bot) (WIP)
