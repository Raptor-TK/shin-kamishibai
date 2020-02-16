# Shin-Kamishibai

> pronounce shin kami-shi-bye

A basic manga reader in browser, server hosted in remote pc. This comes with browser client and server.

Written in Go language, with zero dependency and portability and support widest variety of devices in mind (so no fancy HTML5 and no javascript either).

This is a continuation of my previous project kamishibai.

Please send in pictures when you got the client running on your cool hardware.

## Aim

To make manga reader (client side) work with old/outdated/retro/obscure software, hardware. Netscape/IE/Opera, early Firefox, Windows 9x/XP, Kindle 1/2/DX, Android 4, etc.

## Requirements

Go version >= 1.2 [https://golang.org](https://golang.org)
Browser support cookie and css

## Build

`go build -o shin-kamishibai`

## Build (using Make)

For Linux and Mac and Windows
`go build *.go`

For Windows 2000
`go build`

## Run

1. Copy sample-config.json, and put to \$HOME/etc/shin-kamishibai/config.json
2. Edit and save config
3. Copy `web` to the same place as config
4. Start by running `./shin-kamishibai` in terminal
5. Open web browser and browse `http://<host>:<port>`

## License:

---

BSD 3-clause
Mac Ma gitmac at runbox.com (C) 2013-2020 Copyright
