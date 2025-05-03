---
weight: 999
url: "/Tmux_\\:_le_multiplexeur_de_terminal_remplaçant_de_screen/"
title: "Tmux: The Terminal Multiplexer Replacing Screen"
description: "Guide to installing, configuring and using Tmux, a powerful terminal multiplexer which serves as a modern replacement for screen."
categories: ["Linux"]
date: "2013-01-25T07:47:00+02:00"
lastmod: "2013-01-25T07:47:00+02:00"
tags: ["Terminal", "Tmux", "Screen", "Linux", "Command Line", "SSH"]
toc: true
---

## Introduction

[Tmux](https://tmux.sourceforge.net/) has many advantages over screen and serves as a terminal multiplexer. This functionality is extremely practical and even essential once you've started using it.

## Installation

```bash
aptitude install tmux
```

## Usage

To use tmux, simply launch it:

```bash
tmux
```

Like screen, tmux uses a key combination to access its internal functions. By default, it's "Ctrl+b" (which can be modified) that is used.
To start, you can easily get help (don't forget to press "Ctrl+b" before pressing any key):

{{< table "table-hover table-striped" >}}
| Description | Keys |
|------------|------|
| Get help | ? |
{{< /table >}}

### Window Management

You can manage your windows as follows (don't forget to press "Ctrl+b" before pressing any key):

{{< table "table-hover table-striped" >}}
| Description | Keys |
|------------|------|
| Create a new window | c |
| Get a list of open windows | w |
| Move to the next window | n |
| Move to the previous window | p |
| Move to the last used window | l |
| Move to a window by its number | 0 1 2 3 4 5 6 7 8 9 |
| Search in window buffers | f then "window search name" |
| Rename the current window | , |
| Force close a window | & |
| Display time | t |
{{< /table >}}

### Split

You can split the screen in several ways:

{{< table "table-hover table-striped" >}}
| Description | Keys |
|------------|------|
| Horizontally split the screen | " |
| Vertically split the screen | % |
| Move to the previous pane | { |
| Move to the next pane | } or o |
| Move to the pane corresponding to the key | ← → ↑ ↓ |
| Get pane numbers | q |
| Change visual organization of panes | [space] |
| Resize a pane | Alt+(← → ↑ ↓) |
| Convert a pane from a split into a window | ! |
| Convert a window for integration into a split<br>_ -h: horizontally<br>_ -s 0.0: window 0 and pane 0<br>\* -p 75: taking 75% of window | :joinp -h -s 0.0 -p 75 |
{{< /table >}}

### History

By default, Tmux keeps only 2000 lines of history. Here's how to navigate:

{{< table "table-hover table-striped" >}}
| Description | Keys |
|------------|------|
| Scroll up through history | ↑↑ (PageUP) |
| Scroll down through history after scrolling up | ↓↓ (PageDOWN) |
| Select lines from history (after PageUP) | [space] then (↑/↓) |
| Copy selection | [enter] |
| Paste selection | = |
{{< /table >}}

### Sessions

Session management is something very practical. It's always useful to be able to exit an SSH session and leave time-consuming tasks running, or to protect against network disconnections. That's why when you're in tmux, it's possible to detach from your tmux:

{{< table "table-hover table-striped" >}}
| Description | Keys |
|------------|------|
| Detach from tmux session | d |
| List tmux sessions | s |
| Switch to next tmux session | ) |
| Switch to previous tmux session | ( |
{{< /table >}}

Then reattach later:

```bash
tmux a
```

This command also allows multiple participants to see exactly the same thing.

## Cheat Sheet

I've created a cheat sheet for those interested:

- [Tmux Cheat Sheet PDF French](/pdf/tmux_cheat_sheet_fr.pdf) - [LaTeX Fr](/others/tmux_cheat_sheet_fr.tex)
- [Tmux Cheat Sheet PDF English](/pdf/tmux_cheat_sheet_en.pdf) - [LaTeX En](/others/tmux_cheat_sheet_en.tex)

## Customization

You can customize all kinds of things (I'll let you read the man page as it's so complete), and here's my configuration:

```bash
# Tmux configuration
# Pierre Mavro

# Default shell
set -g default-command zsh

# Screen addict (replacing Ctrl+b by Ctrl+a)
#set -g prefix C-a
#unbind C-b
#bind C-a send-prefix

# Enable utf8
set -g status-utf8 on
setw -g utf8 on

# Same hack as screen to scroll terminal (xterm ...)
set -g terminal-overrides 'xterm*:smcup@:rmcup@'

# Scrollback buffer n lines
set -g history-limit 100000

# Lock session after delay (in seconds)
set -g lock-after-time 720
# Use vlock to lock session (aptitude install vlock)
set -g lock-command vlock
# To unlock as a user
#vlock -a
# To unlock as root
#vlock -sn

# Highlight active window
set-window-option -g window-status-current-bg red

# Split easier keys (| for horizontal and - for vertical)
bind | split-window -h
bind - split-window -v

# Set window notifications when somethings new happen
setw -g monitor-activity on
set -g visual-activity on
#set -g visual-bell on

# Start window to number to 1 (default 0)
set -g base-index 1

# Automatically set window title
setw -g automatic-rename

# force resize local window
#setw -g aggressive-resize on

# Enable mouse
#set -g mouse-select-pane on
#setw -g mode-mouse on
```

## Launch a Program with Tmux at Boot

You may want to launch weechat (IRC client) or other program in tmux when your machine starts. To do this, it's easy:

```bash
su - deimos -c "tmux new-session -d 'weechat-curses'"
```

Here we're asking tmux to create a new session via the deimos user.

## Resources

- http://tmux.sourceforge.net/
- http://myhumblecorner.wordpress.com/2011/08/30/screen-to-tmux-a-humble-quick-start-guide/
- http://blog.hawkhost.com/2010/06/28/tmux-the-terminal-multiplexer/
- http://blog.hawkhost.com/2010/07/02/tmux-%E2%80%93-the-terminal-multiplexer-part-2/
- http://www.dayid.org/os/notes/tm.html
- http://tmux.svn.sourceforge.net/viewvc/tmux/trunk/examples/
- http://linux-attitude.fr/post/configuration-de-tmux?utm_source=rss&utm_medium=rss&utm_campaign=configuration-de-tmux
