---
weight: 999
url: "/ZSH_\\:_un_shell_très_pratique/"
title: "ZSH: A Very Practical Shell"
description: "Configuration and usage guide for ZSH shell, a powerful alternative to bash with advanced features like auto-completion, customizable prompts, and better history management."
categories: ["Linux", "Debian"]
date: "2009-12-11T21:11:00+02:00"
lastmod: "2009-12-11T21:11:00+02:00"
tags: ["Servers", "Shell", "Terminal", "Configuration", "Productivity"]
toc: true
---

## Introduction

I've been thinking about this for a while, but I was waiting to have a sufficiently nice ZSH configuration to publish it. After today's mishap (running `rm -Rf *` in the root directory), I decided to finish it quickly. So here's my configuration which can still be optimized, but is already very good for everyday use :-).

It works on Linux, Mac, BSD, and Windows (Cygwin).

The content of the files below may not be completely up-to-date. If you want my latest ZSH configuration, here's my git repository address: [https://www.deimos.fr/gitweb](https://www.deimos.fr/gitweb)

## Installation

If you don't have ZSH installed, install it in the simplest way possible. On Debian for example:

```bash
apt-get install zsh
```

## Configuration

We'll first configure the current account and you'll see if it suits you or not.

### ~/.zshrc

```bash
#!/bin/zsh

# Version detection
ZSH_VERSION_TYPE=old

if [[ $ZSH_VERSION == 3.1.<6->* || $ZSH_VERSION == 3.2.<->* || $ZSH_VERSION == 4.<->* ]] ; then
    if which zstyle > /dev/null ; then
    ZSH_VERSION_TYPE=new
    fi  
fi

# Uname system
export MYSYSTEM=`uname`

# If you don't want history enter 1
# Do not forget to delete ~/.zshhistory
export NOHIST="0"

# Using environnement files
for envfile in ~/.zsh/* ; do
    test -f $envfile && source $envfile
done

# Default Umask
umask 022 

# Make default color
if [ -z "$nocolor" ] ; then
     c6 && c2; 
else PS1="%n@%m %~ %% "
     export PS1 
fi
```

Then create the .zsh directory:

```bash
mkdir .zsh
```

### ~/.zsh/alias

```bash
#!/bin/zsh

# Usefull alias
alias utar="tar -xvzf"
alias ..='cd ..'
alias ...='cd ../..'
alias uu='source ~/.zshrc &> /dev/null'

# ls
if [ $MYSYSTEM = "Linux" ] ; then
        alias ls='ls --color=auto'
        alias l='ls --color=auto -lg'
        alias ll='ls --color=auto -lag | $PAGER'
fi

if [ $MYSYSTEM = "OpenBSD" ] ; then
        if [ -x /usr/local/bin/gls ] ; then
                alias ls='gls --color=auto'
                alias l='gls --color=auto -lg'
                alias ll='gls --color=auto -lag | $PAGER'
        else
                alias ls='ls'
                alias l='ls -lg'
                alias ll='ls -lag | $PAGER'
        fi  
fi

if [ $MYSYSTEM = "Darwin" ] ; then
        alias ls='ls -G'
        alias l='ls -lGg'
        alias ll='ls -lGag | $PAGER'
fi
```

### ~/.zsh/complet

```bash
#!/bin/zsh

autoload -U compinit
autoload -U colors

## Check if we are using Cygwin ##
if [ `uname` = "CYGWIN*" ] ; then
    compinit -u
else
    compinit
fi

zstyle ':completion:*:*:cd:*' tag-order local-directories path-directories
zstyle ':completion:*' menu select=2
zstyle ':completion:*' select-prompt %SScrolling active: current selection at %p%s
zstyle ':completion:*:rm:*' ignore-line yes 
zstyle ':completion:*:mv:*' ignore-line yes 
zstyle ':completion:*:cp:*' ignore-line yes 

zstyle ':completion:*' verbose yes 
zstyle ':completion:*:descriptions' format '--==[ %B%d%b ]==--'
zstyle ':completion:*:messages' format '--==[ %d ]==--'
zstyle ':completion:*:warnings' format 'No matches for: %d'
zstyle ':completion:*:corrections' format '%B%d (errors: %e)%b'
zstyle ':completion:*' group-name ''

# Color completion
zstyle ':completion:*' list-colors ''
zstyle ':completion:*:*: kill :*:processes' list-colors '=( #b) #([0-9]#)*=0=01;34'
zstyle ':completion:*' list-colors 'di=01;34'

local _myhosts
if [ -d ~/.ssh ]; then
  if [ -f ~/.ssh/known_hosts ];then
    _myhosts=(${=${${(f)"$(<$HOME/.ssh/known_hosts)"}%%[# ]*}//,/ })
   fi  
fi
zstyle ':completion:*' hosts $_myhosts
```

### ~/.zsh/env

```bash
#!/bin/zsh

## Check if we are using Cygwin ##
if [ `uname` = "CYGWIN*" ] ; then
    export TERM=cygwin
else
    export TERM=xterm-color
fi

usernames=( $(cut -d: -f1 /etc/passwd) )
groups=( $(cut -d: -f1 /etc/group) )

case "$TERM" in
    xterm*|rxvt|linux|cygwin)
    ;;  
    *)  
    nocolor=yes
    ;;  
esac

# Set locales
if [ $MYSYSTEM = "SunOS" ] ; then
    export LANGUAGE=fr_FR.ISO8859-15
    export LC_ALL=fr_FR.ISO8859-15
    export LANG=fr_FR.ISO8859-15
else
    export LANGUAGE=fr_FR.UTF8
    export LC_ALL=fr_FR.UTF8
    export LANG=fr_FR.UTF8
fi

export LESSCHARSET=latin9
export MINICOM="-c on"
export LESS="-S -g"
export JAVA_HOME="/usr/lib/jvm/cacao/"
LOGCHECK=60
WATCHFMT="%n has %a %l from %M"


# CVS
export CVS_RSH=/usr/bin/ssh
export CVSROOT=:ext:user@host:/var/lib/cvs

# History
export HISTSIZE=5000
export HISTFILE=$HOME/.zshhistory
if [ $NOHIST = "1" ] ; then
    export SAVEHIST=0
else
    export SAVEHIST=$HISTSIZE
fi

# Defaut editor
export EDITOR=vim

export LISTMAX=0

# Use most if possible
if [[ -x /usr/bin/most || -x /opt/local/bin/most ]] ; then
    export PAGER=most
else
    export PAGER=more
fi

export BLOCK_SIZE=human-readable
export LS_COLORS='no=00:fi=00:di=0;34:ln=01;36:pi=40;33:so=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arj=01;31:*.taz=01;31:*.lzh=01;31:*.zip=01;31:*.z=01;31:*.Z=01;31:*.gz=01;31:*.bz2=01;31:*.rar=01,31:*.par2=01,31:*.deb=01;31:*.rpm=01;31:*.jpg=01;35:*.gif=01;35:*.bmp=01;35:*.pgm=01;35:*.pbm=01;35:*.ppm=01;35:*.tga=01;35:*.png=01;35:*.GIF=01;35:*.JPG=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.mpg=01;37:*.avi=00;35:*.gl=01;37:*.dl=01;37:*.mly=01;37:*.mll=01;37:*.mli=01;37:*.ml=01;37:*.cpp=01;37:*.cc=01;37:*.c=01;37:*.hh=01;37:*.h=01;37:*Makefile=4;32:*.pl=4;32:*.sh=4;32:*.ps=01;34:*.pdf=01;34:*TODO=01;37:*TOGO=01;37:*README=01;37:*LINKS=01;37:*.y=01;37:*.l=01;37:*.algo=01;37'

limit core 0
```

### ~/.zsh/fonctions

```bash
#!/bin/zsh

# Prompt definition
function setprompt { #
    export PROMPT="%{$COLOR2%}[%{$COLOR2%}%D{%H:%M:%S}%{$COLOR2%}] %{$COLOR4%}- %{$COLOR3%}[~%{$COLOR3%}]%{$COLOR4%}
%{$COLOR4%}%n%{$COLOR1%}@%{$COLOR4%}%m%{$COLOR3%}%{$COLOR1%} %{$COLOR1%}>%{$COLOR1%}%{$COLOR1%} "
    export RPROMPT="%{^[[A%}%{$COLOR3%}[%{$COLOR6%}%D{%d-%m-%Y}%{$COLOR3%}]%{$COLOR5%}%{^[[B%}"
    export PROMPT2="%{$COLOR3%}[%{$COLOR1%}%(#.#.$)%{$COLOR3%}]%{$COLOR4%} "
    export SPROMPT="%{$COLOR3%}[%{$COLOR1%}correct%{$COLOR3%}[%{$COLOR1%}'%R'%{$COLOR3%}]%{$COLOR1%}->%{$COLOR3%}[%{$COLOR1%}'%r'%{$COLOR3%}]] [%{$COLOR1%}n%{$COLOR3%}/%{$COLOR1%}y%{$COLOR3%}/%{$COLOR1%}a%{$COLOR3%}/%{$COLOR1%}e%{$COLOR3%}]%{$COLOR4%} "
}

function c1 {
    export COLOR1="^[[0;31m"
    export COLOR2="^[[1;31m"
    export COLOR3="^[[1;30m"
    export COLOR4="^[[0;31m"
    setprompt
}

function c2 {
    # Red
    export COLOR1="^[[1;31m"
    # Green
    export COLOR2="^[[1;32m"
    # White
    export COLOR3="^[[1;33m"
    # Blue
    export COLOR4="^[[1;34m"
    # Grey
    export COLOR5="^[[0;37m"
    # Purple
    export COLOR6="^[[1;35m"
    setprompt
}

function c3 {
    export COLOR1="^[[0;33m"
    export COLOR2="^[[1;33m"
    export COLOR3="^[[1;30m"
    export COLOR4="^[[0;33m"
    setprompt
}
function c4 {
    export COLOR1="^[[0;34m"
    export COLOR2="^[[1;34m"
    export COLOR3="^[[1;30m"
    export COLOR4="^[[0;34m"
    setprompt
}

function c5 {
    export COLOR1="^[[0;35m"
    export COLOR2="^[[1;35m"
    export COLOR3="^[[1;30m"
    export COLOR4="^[[0;35m"
    setprompt
}

function c6 {
    export STATUS_WR="^[[4;37m"
    export STATUS_COLOR="^[[1;33m"
    export LOGIN_COLOR="^[[0;40m"
    export HOST_COLOR="^[[1;37m"
    export COLOR1="^[[0;37m"
    export COLOR2="^[[0;33m"
    export COLOR3="^[[0;31m"
    export COLOR4="^[[0;0m"
    setprompt
}

function c7 {
    export COLOR1="^[[0;37m"
    export COLOR2="^[[1;37m"
    export COLOR3="^[[1;30m"
    export COLOR4="^[[0;0m"
    setprompt
}
```

### ~/.zsh/path

```bash
#!/bin/sh

# OpenBSD Respository
test $MYSYSTEM = "OpenBSD" && export PKG_PATH=ftp://ftp.arcane-networks.fr/pub/OpenBSD/`uname -r`/packages/`machine -a`/

# DarwinPorts Mac OS X
if [[ $MYSYSTEM = "Darwin" && -d /opt/local/bin ]] ; then
    PATH="$PATH:/opt/local/bin:/opt/local/sbin"
fi

# Solaris
if [[ $MYSYSTEM = "SunOS" && -d /usr/openwin/bin/ ]] ; then
    PATH="$PATH:/usr/openwin/bin/:/usr/X11/bin:/opt/csw/bin"
fi

# Define scripts
test -d "$HOME/.scripts" && PATH="$PATH:$HOME/.scripts"

# Define Path
PATH="$PATH:/usr/bin:/sbin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/bin"
export PATH
```

### ~/.zsh/term

```bash
#!/bin/zsh

set convert-meta off # Don't strip high bit when reading or displaying. 
set input-meta on  
set output-meta on
set append history # multiple parallel zsh sessions will all have their history lists added to the history

# No bips
unsetopt beep
unsetopt hist_beep
unsetopt list_beep

unsetopt ignore_eof # Logout Ctrl+D
setopt chase_links # Properly handle symbolic links
setopt hist_verify # Don't execute command when searching history with !
setopt auto_list
setopt auto_cd # If command is invalid but matches a subdirectory name, execute 'cd subdirectory'
setopt auto_remove_slash # When last character of completion is '/' and space is typed after, the '/' is deleted

function common_terms () { 
    bindkey "\e[2~" quoted-insert 
    bindkey "\e[3~" delete-char 
    bindkey "\e[5~" beginning-of-history 
    bindkey "\e[6~" end-of-history 
} # Make the Home, End, and Delete keys work on common terminals. 
if [[ "$TERM" == "linux" ]] ; then 
    common_terms 
    bindkey "\e[1~" beginning-of-line 
    bindkey "\e[4~" end-of-line 
elif [[ "$TERM" == "rxvt" ]] ; then 
    common_terms 
    bindkey "\e[7~" beginning-of-line 
    bindkey "\e[8~" end-of-line 
elif [[ "$TERM" == xterm* ]] ; then 
    common_terms 
    bindkey "\e[1~" beginning-of-line
    bindkey "\e[4~" end-of-line 
fi

bindkey -s '^X^Z' '%-^M'
bindkey '^[e' expand-cmd-path
bindkey -s '^X?' '\eb=\ef\C-x*'
bindkey '^[^I' reverse-menu-complete
bindkey '^[p' history-beginning-search-backward 
bindkey '^[n' history-beginning-search-forward 
bindkey '^W' kill -region bindkey '^I' expand-or-complete-prefix 
bindkey -s '^[[Z' '\t' 
bindkey  '^?' backward-delete-char
if which setterm > /dev/null ; then
    setterm -hbcolor bright white
    setterm -ulcolor cyan
fi
```

## Usage

For usage, here are some basic tricks or features I've implemented:

- `uu`: allows you to refresh your shell (when making modifications or adding new software)
- `utar`: equivalent to tar -xzvf
- `..`: equivalent to cd ..
- `...`: equivalent to cd ../..

For the rest, I'll let you look at the file contents; I try to comment sufficiently.

And for a result, you get something like this:

```
[20:07:45] - [~]                                                   [01-10-2009]
pmavro@pm-laptop >
```
