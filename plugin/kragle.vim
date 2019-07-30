if exists('g:loaded_kragle')
    finish
endif
let s:kragle_job = expand("<sfile>:p:h:h") . '/kragle'
let g:loaded_kragle = 1


" Config variables
" """"""""""""""""
let s:log_path = ""
let s:same_root = v:true

if exists("g:kragle_log_path")
    let s:log_path = g:kragle_log_path
endif
if exists("g:kragle_same_root")
    let s:same_root = g:kragle_same_root
endif



" Public API
" """"""""""
function kragle#SwitchToBuffer()
    let a:file_list = GetFileList()

    let a:file_path = s:select("Switch to file", a:file_list)

    if "" != a:file_path
        edit a:file_path
    endif
endfunction


" Private 
" """""""
" of course some of these can all publicly be called in vim private is mearly intention
function! s:select(message, options)
    let s:choice = inputlist([a:message] + map(copy(a:options), '(v:key+1).". ".v:val'))
    if 1 > s:choice || len(copy(a:options)) < s:choice
        return ""
    endif

    return a:options[s:choice -1]
endfunction


let s:buffer_clean = v:false
function! kragle#swapExists()
    echom "Swap file found for " . expand("<afile>") . ", attempting open on other server."

    let opened = RemoteOpen(expand("<afile>:p"))
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
        else
            echom expand("<afile>")
        endif

        echom "Cleaning up"
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
        \}
endfunction


" Setup kragle plugin
" """""""""""""""""""
function! s:RegisterKragle(host) abort
    return jobstart([s:kragle_job], {'rpc': v:true})
endfunction

call remote#host#Register(s:kragle_job, 'x', function('s:RegisterKragle'))

call remote#host#RegisterPlugin(s:kragle_job, '0', [
\ {'type': 'function', 'name': 'KragleInit', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'RemoteOpen', 'sync': 1, 'opts': {}},
\ ])

call KragleInit(v:servername)


" Auto Command bindings
" """""""""""""""""""""
augroup Kragle
    autocmd!
    autocmd BufEnter * call kragle#bufEnter()
    autocmd SwapExists * call kragle#swapExists()
augroup END
