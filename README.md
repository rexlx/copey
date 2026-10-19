# copey

A lightweight, zero-dependency CLI tool to quickly copy file contents or piped terminal output straight to your system clipboard without needing to select and drag across large terminal buffers.

## Features

- Piped Input: Read directly from standard input (stdin) via pipelines and shell redirections.
- File Input: Pass one or multiple files as arguments to concatenate and copy their contents.
- Cross-Platform: Uses native clipboard utilities across macOS (pbcopy), Linux (wl-copy, xclip, or xsel), and Windows (clip).
- Zero External Dependencies: Implemented purely using the Go standard library.

## Installation

Ensure you have Go installed, then build the binary:

go build -o copey main.go

Move the compiled binary into your system PATH (e.g., /usr/local/bin or ~/.local/bin):

mv copey ~/.local/bin/

## Usage

1. Piped Input (Standard In)

Pipe command outputs, logs, diffs, or stream results directly into copey:

git diff | copey
cat error.log | grep "FATAL" | copey
curl -s https://api.ipify.org | copey

Redirect files via stdin:

copey < config.yaml

2. File Arguments

Pass a file path directly to copy its contents:

copey main.go

Pass multiple files to concatenate and copy them together:

copey header.txt body.txt footer.txt

## Requirements

- macOS: Built-in (pbcopy).
- Linux: Wayland (wl-clipboard) or X11 (xclip / xsel).
  - Fedora/RHEL: sudo dnf install wl-clipboard xclip
  - Debian/Ubuntu: sudo apt install wl-clipboard xclip
  - Arch Linux: sudo pacman -S wl-clipboard xclip
- Windows: Built-in (clip.exe).