# slide-remote-control

Tool to feed Google Slides via http, for Linux Desktop

## Requirements

- xdotool
- ngrok (optional)

## Installation

If you want to use ngrok, install ngrok and authentiacte to it using `ngrok config add-authtoken` command.

Next, install this by `go install` script

```sh
go install github.com/ebiyuu1121/slide-remote-control@latest
```

## Usage

Start http server by command below,

```sh
slide-remote-control
slide-remote-control -ngrok # use ngrok
```

And please access printed url

