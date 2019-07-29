# Kragle
multi monitor support for neovim... sort of\
No more swap file madness when using multiple neovim instances on a single project

### Note
This has only been tested on Fedora 29 using neovim 0.3.8\
I have only tested it with two instances as thats what i use for development but i see no readson why
it wouldnt work with more.

## Installation
Vim Plug\
`Plug 'indeedhat/kragle'`\
At the moment you will need to manually move the kragle binary to somewhere in your path for this plugin to work
i will fix this at some point but right now its late and i need to be up early so im going to bed

## TODO
- [x] make plugin load prebuilt binary from plugin dir
- [ ] optionally limit the server to interacting with instances using the same working directory
- [ ] Open a file on a specific window
- [ ] Test the crap out of it and see if it breaks

## Credit
This was inspired by [codeape2/vim-multiple-monitors](https://github.com/codeape2/vim-multiple-monitors) and some
of the vimscript is pulled from that plugn
