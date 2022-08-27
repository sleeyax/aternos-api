# Aternos API
**Unofficial** [aternos.org](https://aternos.org/) API, initially inspired by the now deprecated [AternosAPI](https://github.com/Duerocraft/AternosAPI) python library. Also 'bypasses' cloudflare without resorting to browser automation or other stupid hacks.

## Disclaimer
This is NOT an AFK bot and never will be! Please be respectful and don't abuse an awesome free service such as Aternos. I do not take any responsibility for what you do with this solution, but I will gladly take this repository down if it results in abuse of the service.

By using this package you are automatically breaking the [TOS](https://aternos.gmbh/en/aternos/terms). In other words your account *could* be suspended at any time. Use at your own risk.

## Usage
### Library
See [examples](./examples) (easy) or the [CLI source code](./cmd) (advanced). See also the auto-generated pkg.go.dev reference documentation [here](https://pkg.go.dev/github.com/sleeyax/aternos-api).

### CLI
This project also comes with a simple command line application to start and stop your server.

Download the binary for your operating system from [releases](https://github.com/sleeyax/aternos-api/releases).

OR:

Manual build & usage instructions:
```
$ git clone https://github.com/sleeyax/aternos-api.git
$ cd aternos-api
$ go mod download
$ go run cmd/main.go
```

Unfortunately the command `go install github.com/sleeyax/aternos-api@latest` is not supported due to a limitation in go regarding 'replace directives'.

## Projects
Projects that are using this package:
* [sleeyax/aternos-discord-bot](https://github.com/sleeyax/aternos-discord-bot)

Made something cool? Let me know or create a PR to add your project to this list!

## License
Licensed under `GNU General Public License v3.0`.

[TL;DR](https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)):

> You may copy, distribute and modify the software as long as you track changes/dates in source files. 
> Any modifications to or software including (via compiler) GPL-licensed code must also be made available under the GPL along with build & install instructions.
