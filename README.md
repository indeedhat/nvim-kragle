# Kragle
> Multi terminal window support for neovim

## Platform Requirements
This has only been tested on Fedora 29 using neovim 0.3.8\
I can think of no reason why this would not work on other linux distributions although it is as of
yet untested.

The binary is prebuilt and included in the repo so should not need go to be installed on the machine to work (linux only)

**This is a neovim plugin and has no support for vanilla vim**

There is currently no support for windows as i dont currently run it but i will add support if it is requested

`xdotool` is required for focusing the remote instance when opening/moving buffers to it.\
if xdtool is not installed on the system it will fallback to calling foreground() however i havnt had
any luck with that function actually working on osx or with i3wm. It may work better on a more
traditional window manager in linux but i havnt tested it.

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

`g:kragle_use_tabs` defualt v:true\
when moving files between clients or opening on a remote should tabe be used over e

## Public API
`kragle#SwitchToBuffer()`\
Pick from a list of buffers open in all connected windows and switch to said file\
(file will open in whatever window it belongs to)

`kragle#AdoptBuffer()`\
Pick from a list of buffers open in remote windows and move one of them to the current one

`kragle#OrpahBuffer()`\
Move a buffer from the current instance to a remote one\
If only one remote buffer is open it will auto move it otherwise you will be asked to chose

`kragle#Quit(save, force)` (`:qa`)\
Quit all connected clients (including self)\
save: bool (save all files before quitting `:wqa`)\
force: bool (force the quit, ignore errors etc `:qa!`)

`kragle#FocusRemote()`\
Switch window focus to a remote instance.\
This will auto focus if there is only one remote instance otherwise it will prompt for selection

`kragle#OpenOnRemote(path)`\
Open a buffer by path on a remote client.\
If more than one remote client exists it will ask for a choice\
path: string (the full path of the file to be opened)

## My Bindings
```vim
" adopt a buffer from a remote terminal
noremap <Leader>ka :call kragle#AdoptBuffer()<CR>

"move current buffer to a remote terminal
noremap <Leader>ko :call kragle#OrphanBuffer()<CR>

" Quit all connected terminals
noremap <Leader>kq :call kragle#Quit(v:false, v:false)<CR>

" Save and force quit all connected terminals
noremap <Leader>kQ :call kragle#Quit(v:true, v:true)<CR>

" Focus a remote terminal
noremap <Leader>kf :call kragle#FocusRemote()<CR>

" Select and switch to an open buffer (all connected terminals)
noremap <Leader>kl :call kragle#SwitchToBuffer()<CR>
```

## RoadMap
- [x] make plugin load prebuilt binary from plugin dir
- [x] optionally limit the server to interacting with instances using the same working directory
- [x] switch to currently open file (all windows)
- [x] Adopt a file from a remote into the local (close it there and open it here)
- [x] Move current file to remote instance
- [x] Close all instances (:qa)
- [ ] Try and find a better way to name servers
- [ ] make work with quick fix lists (:cn)
- [x] focus remote client when opening a file in it (requires xdotool)
- [x] Open a file on a specific window
- [ ] fuzzy find all open files (possibly look into intergrating remote open into ctrlP or NERDTree)
- [ ] Write some vader tests (i think for the most part the go is untestable)
- [ ] Improve the documentation

## Credit
This was inspired by [codeape2/vim-multiple-monitors](https://github.com/codeape2/vim-multiple-monitors) and some
of the vimscript is pulled from that plugn
