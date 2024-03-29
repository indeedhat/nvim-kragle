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
let s:buffer_clean = v:false

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
    call s:select("Remote Client:", l:server_list, v:true, { 
    \    server_path -> "" != server_path && KragleRemoteFocus(server_path)
    \ })
endfunction

function kragle#OpenOnRemote(file)
    let l:server_list = KragleListServers()
    call s:select("Remote Client:", l:server_list, v:true, {
    \    server_path -> "" != server_path && KragleRemoteOpen(a:file, server_path)
    \ })
endfunction

function kragle#Quit(save, force)
    let l:write = v:true == a:save ? 'w' : ''
    let l:quit = 'qa'
    let l:force = v:true == a:force ? '!' : ''

    call KragleCommandAll(l:write . l:quit . l:force)
endfunction

function kragle#SwitchToBuffer()
    let l:file_list = KragleListAllFiles()
    call s:select("Switch to file", l:file_list, v:false, {
    \    file_path -> '' != file_path && timer_start(0, { ->  execute( "e " . fnameescape(file_path)) })
    \ })
endfunction

function kragle#AdoptBuffer()
    let l:file_list = KragleListRemoteFiles()
    call s:select("Adopt file", l:file_list, v:false, { 
    \    file_path -> "" != file_path && KragleAdoptBuffer(file_path) 
    \ })
endfunction

function kragle#OrphanBuffer()
    let l:file_path = expand("%:p")
    if '' == l:file_path 
        return
    endif

    let l:server_list = KragleListServers()
    call s:select("Orphan file:", l:server_list, v:true, {
    \    server_path -> '' != server_path && KragleOrphanBuffer(l:file_path, server_path) 
    \ })
endfunction


" Private 
" """""""
" of course some of these can all publicly be called in vim private is mearly intention
function! s:select(message, options, auto_pick, cb)
    if empty(a:options) || 0 == len(copy(a:options)) 
        return 
    elseif 1 == len(copy(a:options)) && v:true == a:auto_pick
        " if there is only one option might as well auto pick it
        call a:cb(a:options[0])
        return
    endif

    if get(g:, 'loaded_fzf', 0)
        call s:fzf_select(a:message, a:options, a:cb)
    else
        call s:input_select(a:message, a:options, a:cb))
    endif
endfunction

function! s:fzf_select(message, choices, cb)
    call fzf#run({
    \    'source': s:add_number_to_choices(a:choices),
    \    'down': '35%',
    \    'options': [
         \    '--tiebreak=index',
         \    '--layout=reverse-list',
    \    ],  
    \    'sink': {key -> a:cb(a:choices[key - 1])}
    \ })
endfunction

function! s:input_select(message, options, cb)
    let l:choice = inputlist([a:message] + map(copy(a:options), '(v:key + 1).". ".v:val'))
    if 1 > l:choice || len(copy(a:options)) < l:choice
        return
    endif

    call a:cb(a:options[l:choice -1])
endfunction

" This has been lifted directly from the phpactor vim plugin src
function! s:add_number_to_choices(choices)
    return map(copy(a:choices), {key, value -> key + 1 .') '. value})
endfunction

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

        let s:buffer_clean = v:false
    endif
endfunction

function! kragle#getConfig()
    return {
    \    "client_root": getcwd(),
    \    "server_name": v:servername,
    \    "log_path": s:log_path,
    \    "same_root": s:same_root,
    \    "use_tabs": s:use_tabs,
    \ }
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
\    {'type': 'function', 'name': 'KragleAdoptBuffer', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleInit', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleListAllFiles', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleListRemoteFiles', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleRemoteOpen', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleRemoteFocus', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleRemoteFocusBuffer', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleOrphanBuffer', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleListServers', 'sync': 1, 'opts': {}},
\    {'type': 'function', 'name': 'KragleCommandAll', 'sync': 1, 'opts': {}},
\ ])

call s:init()


" Auto Command bindings
" """""""""""""""""""""
augroup Kragle
    autocmd!
    autocmd BufEnter * call kragle#bufEnter()
    autocmd SwapExists * call kragle#swapExists()
augroup END
