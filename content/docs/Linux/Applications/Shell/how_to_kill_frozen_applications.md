---
weight: 999
url: "/Comment_tuer_une_application_qui_a_planté,_ou_comment_faire_face_à_un_écran_«_figé_»_\\?/"
title: "How to Kill a Crashed Application or Handle a Frozen Screen"
description: "Methods to recover from frozen applications or system hangs in Linux without resorting to hard resets"
categories: ["Linux", "Troubleshooting"]
date: "2008-06-25T20:31:00+02:00"
lastmod: "2008-06-25T20:31:00+02:00"
tags: ["crash", "freeze", "recovery", "sysrq", "kill"]
toc: true
---

You might have been told that "Linux is great because it never crashes!" Well, here you are, helplessly staring at your frozen screen for the past 5 minutes, with an unresponsive keyboard and a mouse you're shaking in all directions... Don't panic: there's a way to solve the problem, and in a much less brutal way than using the Power or Reset buttons on your computer!

Several scenarios may occur:

* If the crashed program was launched from the command line, you still have access to the launch terminal, and your keyboard works, simply press [Ctrl]+[C]. You'll regain control, and all windows of the application will disappear.

* If the program was launched graphically (via a menu or double-clicking an icon), open a terminal and type the command: ~$ killall -9 program_name.

* If the previous solution doesn't work, you can replace the program name with its PID (Process IDentifier). To determine the PID of the crashed application, use the command:

```bash
~$ ps -ef 
```

The pattern corresponds to a string contained in the name of the offending program. For example, if the Gnome editor, Gedit, has crashed, enter the command:

```bash
~$ ps -ef 
```

The application closes and you regain control.

* If only the mouse has abandoned you, don't forget that the keyboard shortcut [Alt]+[F2] allows you to call the application launcher. You can then type the name of your command terminal (gnome-terminal under Gnome or konsole under KDE), then enter a command to kill the application.

* If the X graphical server has crashed and your keyboard still works, it will be impossible to launch or use a graphical application. Therefore, you can't use your graphical console to kill the application that caused your graphical server to freeze. In this case, you need to switch to tty text mode (an absolute text mode that allows direct communication with the core of your machine).

To do this, on some distributions, repeatedly pressing the [Windows] key on the keyboard should do the trick. This key is sometimes configured to allow switching from your graphical interface to different tty interfaces. Then, after entering your login and password, you can use the kill command as specified above.

However, the [Windows] key may not be configured this way. In that case, the shortcuts [Ctrl]+[Alt]+[F1], [Ctrl]+[Alt]+[F2], [Ctrl]+[Alt]+[F3], [Ctrl]+[Alt]+[F4], [Ctrl]+[Alt]+[F5], etc., allow you to switch to tty1, tty2, tty3, etc. modes respectively. To return to graphical mode, use [Alt]+[F7].

This actually depends on what is defined by default in the /etc/inittab file (lines like 1:2345:respawn:/sbin/getty 38400 tty1, 2:23:respawn:/sbin/getty 38400 tty2, etc.), but this is what is observed on most distributions.

* If despite stopping the responsible applications, your system is still frozen, then your only option is to enter the following command in tty mode, which allows for a "clean" restart:

```bash
~$ shutdown -r now
```

When the problem is related to the X server, the simplest solution is to restart it. To do this, just type the keyboard combination [Ctrl]+[Alt]+[Backspace]. The X server will then kill all graphical applications and thus the current user sessions, then restart itself. This brings you back to the login screen.

Another solution is to use a series of keyboard combinations, based on the [SysRq] key (same key as [Print screen]), to be typed in a well-defined order:

* [Alt]+[SysRq]+[R]: places the keyboard in "raw mode". Then, try pressing [Ctrl]+[Alt]+[Backspace] to kill the X server. If that doesn't work, continue with what follows.
* [Alt]+[SysRq]+[S]: this allows writing all unsaved data to disk (referred to as disk "synchronization").
* [Alt]+[SysRq]+[E]: to send a termination signal to all processes, except init.
* [Alt]+[SysRq]+[I]: to kill all active processes, except init.
* [Alt]+[SysRq]+[U]: to unmount, then remount all partitions in read-only mode (this will avoid a filesystem check at restart).
* [Alt]+[SysRq]+[B]: to restart the system. You can also press the reset button on your machine.

These key combinations allow sending commands directly to the kernel, commands that will allow recording open files despite the absence of a graphical interface, since the latter is frozen.

Warning, for this to work, your kernel must have been compiled with support for "magic keys" (the CONFIG_MAGIC_SYSRQ option must be set to "y"), and it must be activated in /proc (which is almost always the case in common distributions like Ubuntu or Mandriva). To verify:

```bash
~$ grep CONFIG_MAGIC_SYSRQ /boot/config-2.6.15-27-386
CONFIG_MAGIC_SYSRQ=y
~$ cat /proc/sys/kernel/sysrq
1
```
