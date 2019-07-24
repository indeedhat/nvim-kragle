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

let s:buffer_clean = v:false

function! Swap_Exists()
    echom "Swap file found for " . expand("<afile>") . ", attempting open on other server."

    let opened = RemoteOpen(expand("<afile>:p"))
    if "opened" != opened 
        echom "Could not find remote file"
        return
    endif

    let s:buffer_clean = v:true 
    let v:swapchoice = "a"
endfunction

function! Buf_Enter()
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

augroup Kragle
    autocmd!
    autocmd BufEnter * call Buf_Enter()
    autocmd SwapExists * call Swap_Exists()
augroup END
