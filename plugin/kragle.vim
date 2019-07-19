if exists('g:loaded_kragle')
    finish
endif
let g:loaded_kragle = 1

function! s:RegisterKragle(host) abort
    return jobstart(['kragle'], {'rpc': v:true})
endfunction

call remote#host#Register('kragle', 'x', function('s:RegisterKragle'))

call remote#host#RegisterPlugin('kragle', '0', [
\ {'type': 'function', 'name': 'KragleLog', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'KragleInit', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'RemoteOpen', 'sync': 1, 'opts': {}},
\ ])

call KragleInit(v:servername)

if exists("g:kragle_log_path")
    call KragleLog(g:kragle_log_path)
endif

let s:buffer_to_cleanup = ""

function! Swap_Exists()
    echom "Swap file found for " . expand("<afile>") . ", attempting open on other server."

    let opened = RemoteOpen(expand("<afile>:p"))
    if "opened" != opened 
        echom "Could not find remote file"
        return
    endif

    let s:buffer_to_cleanup = expand("<afile>:p")
    let v:swapchoice = "a"
endfunction

function! Buf_Enter()
    if s:buffer_to_cleanup != ""
        echom "Cleaning up " . s:buffer_to_cleanup
        execute "bdelete " . s:buffer_to_cleanup
        let s:buffer_to_cleanup = ""
    endif
endfunction

augroup Kragle
    autocmd!
    autocmd BufEnter * call Buf_Enter()
    autocmd SwapExists * call Swap_Exists()
augroup END
