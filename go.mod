module github.com/ganiskowicz/aternos-api

go 1.17

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/dop251/goja v0.0.0-20211217115348-3f9136fa235d
	github.com/gorilla/websocket v1.5.0
	github.com/refraction-networking/utls v1.0.0
	github.com/sleeyax/gotcha v0.1.3
	github.com/sleeyax/gotcha/adapters/fhttp v0.0.0-20220513160314-4b06cd561da9
	github.com/useflyent/fhttp v0.0.0-20211004035111-333f430cfbbf
)

require (
	github.com/Sleeyax/urlValues v1.0.0 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/dlclark/regexp2 v1.4.1-0.20201116162257-a2a8dda75c91 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace github.com/refraction-networking/utls => github.com/sleeyax/utls v1.1.1

replace github.com/gorilla/websocket => github.com/sleeyax/websocket v1.5.1-0.20220512160613-502bd65db8ae
