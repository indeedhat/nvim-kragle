# Kragle
multiple instances of neovim held together with krazy glue

## Platform Requirements
This has only been tested on Fedora 29 using neovim 0.3.8\
I can think of no reason why this would not work on other linux distributions although it is as of
yet untested.

As for windows it is totally untested as i currently dont run windows at all

The binary is prebuilt and included in the repo so should not need go to be installed on the machine to work (linux only)

** This is a neovim plugin and has no support for vanilla vim **

## Installation
It can of course be installed with your package manager of choice, mine is Plug
Vim Plug\
`Plug 'indeedhat/kragle'`

### Aditional steps for osx
The binary shipped in this repo is built for linux and will not work on osx.\
To get kragle working on osx you will need to have go >= 1.11 installed 
    - cd into the kragle location, on with Plug this is ~/.vim/plugged/nvim-kragle
    - rebuild `go build .`


## Options
`g:kragle_log_path` default ""\
the core go binary will output to this log file if set

`g:kragle_same_root` default v:true\
if set to true then the server will only attempt to open files on vim instances that share the same file root
as the current instance.\
Files open in other instances will show the default swap found message

## Public API
`kragle#SwitchToBuffer()`\
Pick from a list of buffers open in all connected windows and switch to said file\
(file will open in whatever window it belongs to)

`kragle#AdoptBuffer()`\
Pick from a list of buffers open in remote windows and move one of them to the current one

`kragle#OrpahBuffer()`\
Move a buffer from the current instance to a remote one\
If only one remote buffer is open it will auto move it otherwise you will be asked to chose

## TODO
- [x] make plugin load prebuilt binary from plugin dir
- [x] optionally limit the server to interacting with instances using the same working directory
- [x] switch to currently open file (all windows)
- [x] Adopt a file from a remote into the local (close it there and open it here)
- [x] Move current file to remote instance
- [ ] Try and find a better way to name buffers
- [ ] make work with quick fix lists (:cn)
- [ ] focus remote client when opening a file in it
- [ ] Open a file on a specific window
- [ ] fuzzy find all open files
- [ ] Write some vader tests (i think for the most part the go is untestable)
- [ ] Improve the documentation

## Credit
This was inspired by [codeape2/vim-multiple-monitors](https://github.com/codeape2/vim-multiple-monitors) and some
of the vimscript is pulled from that plugn
