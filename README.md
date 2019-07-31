# Kragle
multiple instances of neovim held together with krazy glue

### Note
This has only been tested on Fedora 29 using neovim 0.3.8\
I have only tested it with two instances as thats what i use for development but i see no readson why
it wouldnt work with more.

## Installation
Vim Plug\
`Plug 'indeedhat/kragle'`\

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

## TODO
- [x] make plugin load prebuilt binary from plugin dir
- [x] optionally limit the server to interacting with instances using the same working directory
- [x] switch to currently open file (all windows)
- [x] Adopt a file from a remote into the local (close it there and open it here)
- [x] Move current file to remote instance
- [ ] make work with quick fix lists (:cn)
- [ ] focus remote client when opening a file in it
- [ ] Open a file on a specific window
- [ ] fuzzy find all open files
- [ ] Test the crap out of it and see if it breaks

## Credit
This was inspired by [codeape2/vim-multiple-monitors](https://github.com/codeape2/vim-multiple-monitors) and some
of the vimscript is pulled from that plugn
