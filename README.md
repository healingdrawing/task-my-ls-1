# my-ls-1</h3>

## The `my-ls-1` is a project, consisting on creating your own `ls` command.

## Table of Contents

- [About](#about)
- [Usage](#usage)
- [Authors](#authors)

[⬆️](#table-of-contents)

## About

`my-ls-1` is a command-line tool written in Go that allows users to list the files and folders in a directory.  
Designed to behave similarly to the original `ls` command with a few variations, including at least the following flags:

- `-l` to display the files and directories in a long format, including file permissions, ownership, size, and modification date.
- `-R` to recursively list subdirectories.
- `-a` to include hidden files and directories in the list.
- `-r` to reverse the order of the list.
- `-t` to sort the list by modification time.

[⬆️](#table-of-contents)

## Usage

### Installing

To install `my-ls-1`, follow these steps:

- clone the repository
- open terminal inside `my-ls-1` folder
- execute inside terminal: `go build .`

The `my-ls-1` executable will be created

### Running

Inside terminal: `./my-ls-1 [flags] [path]`

- `[flags]` are optional and can be at least any combination of -l, -R, -a, -r, and -t.
- `[path]` is optional. If no path is specified, the tool will default to the current directory.

### Examples

- List the files and directories in the current directory:
  `./my-ls-1`

- List the files and directories in a specific directory:
  `./my-ls-1 /path/to/directory`

- List the files and directories in a specific directory in long format:
  `./my-ls-1 -l /path/to/directory`

- List the files and directories in a specific directory recursively:
  `./my-ls-1 -R /path/to/directory`

- List all files and directories in a specific directory, including hidden files and directories:
  `./my-ls-1 -a /path/to/directory`

- Sort the files and directories in a specific directory by modification time:
  `./my-ls-1 -t /path/to/directory`

- Sort the files and directories in a specific directory in reverse order:
  `./my-ls-1 -r /path/to/directory`

[⬆️](#table-of-contents)

## Authors

- [@Mijan](https://01.gritlab.ax/git/Mijan)
- [@MariaSagulin](https://01.gritlab.ax/git/MariaSagulin)
- [@healingdrawing](https://healingdrawing.github.io "aka @lenivaya10003 on grit:lab Åland")
