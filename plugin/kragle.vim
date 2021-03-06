if exists('g:loaded_kragle')
    finish
endif
let s:kragle_job = expand("<sfile>:p:h:h") . '/kragle'
let g:loaded_kragle = 1


" Config variables
" """"""""""""""""
let s:log_path = ""
let s:same_root = v:true
let s:use_tabs = v:true
let s:window_id = ""

if exists("g:kragle_log_path")
    let s:log_path = g:kragle_log_path
endif
if exists("g:kragle_same_root")
    let s:same_root = g:kragle_same_root
endif
if exists("g:kragle_use_tabs")
    let s:use_tabs = g:kragle_use_tabs
endif



" Public API
" """"""""""
function kragle#FocusRemote()
    let l:server_list = KragleListServers()
    let l:server_path = s:select("Remote Client:", l:server_list, v:true)

    if "" != l:server_path
        echo "calling KragleRemoteFocus"
        call KragleRemoteFocus(l:server_path)
    endif
endfunction

function kragle#OpenOnRemote(file)
    let l:server_list = KragleListServers()
    let l:server_path = s:select("Remote Client:", l:server_list, v:true)

    if "" != l:server_path
        call KragleRemoteOpen(a:file, l:server_path)
    endif
endfunction

function kragle#Quit(save, force)
    let l:command = ""
    if v:true == a:save
        let l:command .= "w"
    endif

    let l:command .= "qa"

    if v:true == a:force 
        let l:command .= "!"
    endif

    call KragleCommandAll(l:command)
endfunction

function kragle#SwitchToBuffer()
    let l:file_list = KragleListAllFiles()
    let s:file_path = s:select("Switch to file", l:file_list, v:false)

    if "" == s:file_path
        return 
    endif

    execute "e " . fnameescape(s:file_path)
endfunction

function kragle#AdoptBuffer()
    let l:file_list = KragleListRemoteFiles()
    let s:file_path = s:select("Adopt file", l:file_list, v:false)

    if "" == s:file_path
        return
    endif

    echo 'Adopting ' . fnameescape(s:file_path)
    call KragleAdoptBuffer(s:file_path)
endfunction

function kragle#OrphanBuffer()
    let l:server_list = KragleListServers()
    let l:server_path = s:select("Orphan file:", l:server_list, v:true)

    if "" != l:server_path
        call KragleOrphanBuffer(expand("%:p"), l:server_path)
    endif
endfunction

" Private 
" """""""
" of course some of these can all publicly be called in vim private is mearly intention
function! s:select(message, options, auto_pick)
    if empty(a:options) || 0 == len(copy(a:options)) 
        return ""
    elseif 1 == len(copy(a:options)) && v:true == a:auto_pick
        " if there is only one option might as well auto pick it
        return a:options[0]
    endif

    let l:choice = inputlist([a:message] + map(copy(a:options), '(v:key+1).". ".v:val'))
    if 1 > l:choice || len(copy(a:options)) < l:choice
        return ""
    endif

    return a:options[l:choice -1]
endfunction


let s:buffer_clean = v:false
function! kragle#swapExists()
    echom "Swap file found for " . expand("<afile>") . ", attempting open on other server."

    let opened = KragleRemoteFocusBuffer(expand("<afile>:p"))
    if "opened" != opened 
        echom "Could not find remote file"
        return
    endif

    let s:buffer_clean = v:true 
    let v:swapchoice = "a"
endfunction

function! kragle#bufEnter()
    if s:buffer_clean
        if "" == expand("<afile>")
            return 
        endif

        " execute "bdelete"
        let s:buffer_clean = v:false
    endif
endfunction

function! kragle#getConfig()
    return {
        \"client_root": getcwd(),
        \"server_name": v:servername,
        \"log_path": s:log_path,
        \"same_root": s:same_root,
        \"use_tabs": s:use_tabs,
        \}
endfunction

function! kragle#focus()
    if "" != s:window_id
        call system("xdotool windowfocus " . s:window_id)
    else
        call foreground()
    endif

endfunction

function! kragle#trackWindowId()
    if executable("xdotool")
        let s:window_id = system("xdotool getactivewindow")
    endif
endfunction

function! s:init()
    " initialize the go plugin
    call KragleInit(v:servername)

    " attempt to get the window id for later use with focus
    call kragle#trackWindowId()
endfunction


" Setup kragle plugin
" """""""""""""""""""
function! s:RegisterKragle(host) abort
    return jobstart([s:kragle_job], {'rpc': v:true})
endfunction

call remote#host#Register(s:kragle_job, 'x', function('s:RegisterKragle'))

call remote#host#RegisterPlugin(s:kragle_job, '0', [
\ {'type': 'function', 'name': 'KragleAdoptBuffer', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleInit', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleListAllFiles', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleListRemoteFiles', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleRemoteOpen', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleRemoteFocus', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleRemoteFocusBuffer', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleOrphanBuffer', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleListServers', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'KragleCommandAll', 'sync': 1, 'opts': {}},
\ ])

call s:init()


" Auto Command bindings
" """""""""""""""""""""
augroup Kragle
    autocmd!
    autocmd BufEnter * call kragle#bufEnter()
    autocmd SwapExists * call kragle#swapExists()
    " autocmd FocusGained * call kragle#trackWindowId()
augroup END
