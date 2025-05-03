---
weight: 999
url: "/Vim_\\:_Les_indispendables_&_Quick_Reference_Card/"
title: "Vim: Essential Commands and Quick Reference Card"
description: "A guide to essential Vim commands, functions, and configuration options with a quick reference card for Vim users."
categories: ["Development", "Tools"]
date: "2013-05-07T13:17:00+02:00"
lastmod: "2013-05-07T13:17:00+02:00"
tags: ["Vim", "Editor", "Development", "Configuration"]
toc: true
---

## References

Here are some useful commands that I don't always remember:

Replace text in the current line:

```
:s/original/destination/g
```

Replace text throughout the document:

```
:%s/original/destination/g
```

Remove unnecessary spaces and tabs at the end of lines:

```
:%s/\s\+$//
```

Here are also the [Quick Reference Card](/pdf/vimqrc.pdf) and [the same in French](/pdf/vimqrcfr.pdf) :-).

## The Functions

Here is what I consider essential in VIM for good file editing and better navigation:

```
" /etc/vim/vimrc, ~/.vimrc or ~/.exrc
" VIM Configuration File

" Made by Pierre Mavro
" Version 0.3

" Show matching parentheses
set showmatch
" Line numbering
set number
" Try to keep the cursor in the same column when changing lines
set nostartofline
" Auto-completion options
set wildmode=list:full
" Always keep one visible line above the cursor on the screen
set scrolloff=1
" Euro encoding with accents
set enc=UTF-8
" Character font for Gvim that supports the euro symbol
set guifont=-misc-fixed-medium-r-semicondensed-*-*-111-75-75-c-*-iso8859-15
" Enable syntax highlighting
syntax on
" Use the standard color scheme
colorscheme default
" Show cursor position 'line,column'
set ruler
" Make searches case insensitive
set ignorecase
" Highlight search expressions
set hlsearch

" Tabulation section
set autoindent
set expandtab
set shiftwidth=4
set softtabstop=4
set tabstop=4

" Message function
function! s:DisplayStatus(msg)
    echohl Todo
    echo a:msg
    echohl None
endfunction

" Mouse state mode
let s:mouseActivation = 1 
" Activating / Desactivating mouse mode
function! ToogleMouseActivation()
    if (s:mouseActivation)
        let s:mouseActivation = 0 
        set mouse=n
        set paste
        call s:DisplayStatus('Paste Mode Desactivated')
    else
        let s:mouseActivation = 1 
        set mouse=a
        set nopaste
        call s:DisplayStatus('Paste Mode Activated')
    endif
endfunction
" Activating mouse mode by default
" set mouse=a
set paste

" Cleaning function
" Call this function : ':call Clean()'
function! Clean()
    %retab
    %s/^M//g
    %s/\s\+$//
    call s:DisplayStatus('Code cleaned')
endfunction

" Advanced completion
" Use : 'Ctrl+x & Ctrl+o | Ctrl+x & Ctrl+k | Ctrl+x & Ctrl+n'
function! MultipleAutoCompletion()
    if &omnifunc != ''
        return "\<C-x>\<C-o>"
    elseif &dictionary != ''
        return "\<C-x>\<C-k>"
    else
        return "\<C-x>\<C-n>"
    endif
endfunction

" Fold function
function! MyFoldFunction()
    let line = getline(v:foldstart)
    let sub = substitute(line, '/\*\|\*/\|^\s+', '', 'g')
    let lines = v:foldend - v:foldstart + 1
    return '[+] '. v:folddashes.sub . '...' . lines . 'lines...' .
    getline(v:foldend)
endfunction
" Activate Fold function
" Use : za (open or close), zm (close all) or zr (close zr)
set foldenable
set fillchars=fold:=
set foldtext=MyFoldFunction()

" Comment / Uncomment with # or //
" Use : Shift+v, select wished lines and press F5 or F6
map <F5> :s.^.#. <CR> :noh <CR>
map <F6> :s.^#.. <CR> :noh <CR>
" Use : Shift+v, select wished lines and press F7 or F8
map <F7> :s.^.\/\/. <CR> :noh <CR>
map <F8> :s.^\/\/.. <CR> :noh <CR>

" Indent or not from 4 spaces
" Use : Shift+v, select wished lines and press F3 or F4
map <F3> :s.^    .. <CR> :noh <CR>
map <F4> :s.^.    . <CR> :noh <CR>
```

This should be placed in "**~/.exrc**". You will then have auto-indentation, line numbering, etc...

If you want to get the latest version of my vimrc, go to my git: [https://www.deimos.fr/gitweb](https://www.deimos.fr/gitweb)

## Resources
- [https://vimcasts.org/](https://vimcasts.org/)
- [https://vim-adventures.com/](https://vim-adventures.com/)
- [Vim Config Generator](https://yoursachet.com/)
- [Quick Reference Card](/pdf/vimqrc.pdf)
- [Quick Reference Card (French)](/pdf/vimqrcfr.pdf)
- [The new and improved Vim editor](/pdf/au-speakingunix_vim-pdf.pdf)
