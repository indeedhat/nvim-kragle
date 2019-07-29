# Kragle
multiple instances of neovim held together with krazy glue

### Note
This has only been tested on Fedora 29 using neovim 0.3.8\
I have only tested it with two instances as thats what i use for development but i see no readson why
it wouldnt work with more.

## Installation
Vim Plug\
`Plug 'indeedhat/kragle'`\
At the moment you will need to manually move the kragle binary to somewhere in your path for this plugin to work
i will fix this at some point but right now its late and i need to be up early so im going to bed

## Options
`g:kragle_log_path` string\
the core go binary will output to this log file if set

`g:kragle_same_root` bool\
if set to true then the server will only attempt to open files on vim instances that share the same file root
as the current instance.\
Files open in other instances will show the default swap found message

## TODO
- [x] make plugin load prebuilt binary from plugin dir
- [x] optionally limit the server to interacting with instances using the same working directory
- [ ] Open a file on a specific window
- [ ] Move current file to remote instance
- [ ] Move remote file to self
- [ ] Test the crap out of it and see if it breaks
- [ ] focus remote client when opening a file in it
- [ ] fuzzy find all open files

## Credit
This was inspired by [codeape2/vim-multiple-monitors](https://github.com/codeape2/vim-multiple-monitors) and some
of the vimscript is pulled from that plugn
