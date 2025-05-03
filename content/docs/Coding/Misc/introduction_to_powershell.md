---
weight: 999
url: "/Introduction_au_PowerShell/"
title: "Introduction to PowerShell"
description: "An introduction to Microsoft PowerShell, explaining its object-oriented nature and advantages compared to traditional command-line interfaces."
categories: ["Windows"]
date: "2008-09-24T15:02:00+02:00"
lastmod: "2008-09-24T15:02:00+02:00"
tags: ["Windows", "Development", "Divers"]
toc: true
---

## Introduction

### Presentation of PowerShell

You may be wondering how PowerShell will help you, especially since you've learned to do without it until now.

Due to the limited commands provided by the "cmd.exe" environment (inherited from MS-DOS), most system administrators were forced to turn to "pseudo" programming languages such as Perl, KixStart, or VBScript to perform simple tasks.

Will Microsoft succeed in reconciling administrators of its products with a command-line environment (or shell)? This is not an easy bet, considering that Windows users have a more graphical culture than command-line. However, the graphical interface quickly shows its limitations when it comes to performing many repetitive actions. Therefore, to perform these tasks, there is no alternative other than turning to scripting and thus to PowerShell...

I can see some of you smirking behind your screen. Those people think that Unix has been doing this very well for a long time and that Microsoft hasn't invented anything! Well, these people are only half wrong. Indeed, Microsoft could have thought of it sooner, but Microsoft didn't just try to do as well as Unix shells, it tried to do much better.
Let's stop any controversy right away; the purpose of this article is not to compare Unix shells to PowerShell, but rather to try to understand what PowerShell really brings.

### New Features

Purists or "pro-Microsoft" people could make an endless list, but I will give you my perception of things, that of a confirmed system administrator.

Here are, in my opinion, the most important advances:

- PowerShell is object-oriented
- PowerShell provides access to .Net objects
- PowerShell brings an exceptional richness of commands
- PowerShell offers consistency of commands and associated parameters
- PowerShell provides comprehensive help on commands
- PowerShell is easy to learn

#### Object Orientation of PowerShell

In a traditional environment, executing each command returns text. Take for example the "DIR" command. It provides a textual list of files and directories in return. In PowerShell, the equivalent of this command is called "Get-Childitem". Although its execution returns roughly the same thing to the screen as the "DIR" command, it actually returns a list of objects to the screen. These objects are most often of file or directory type, but they can also be registry keys. You will discover that most PowerShell commands are generic.

Additionally, when you use the | "pipe" to pass the result of one command to another command. For example: "dir | more". In PowerShell, you pass complete objects instead of passing text.

You may not understand the subtlety for now, but you will discover that this is one of the great strengths of PowerShell.

So to continue the previous example, you might want to get only the size of a file.
How to do this with "command prompt and the dir command"?

Let's try...

```text
C:\TEMP> Dir monfichier.txt
Le volume dans le lecteur C n'a pas de nom.
 Le numéro de série du volume est 78D5-739D

 Répertoire de C:\TEMP

28/11/2006  13:00            82 944 monfichier.txt
               1 fichier(s)           82 944 octets
               0 Rép(s)   3 140 132 864 octets libres
```

It's not easy now to get just the file size. You'll have to count the number of lines and characters to cut the string and finally hope to get the size by passing this command via the pipe to another command. Tedious!

And in VBScript? It's possible but in several lines of code.

And now in PowerShell?

```
PS C:\TEMP> $b = Get-ChildItem monfichier.txt; $b.length
```

Similarly, we could very easily get the file's creation date or last access date with the following commands:

```
$b = Get-ChildItem monfichier.txt; $b.CreationTime
$b = Get-ChildItem monfichier.txt; $b.LastAccessTime
```

#### The Power of .Net

All PowerShell commands have been written based on .Net object classes, whereas initially one might have thought it was WMI. This allows for two things: first, you have access to the entire .Net framework classes (which offers extended functionality compared to WMI), and second, if you are familiar with .Net development, it will be all the easier for you to learn PowerShell (and vice versa of course!).

Let's not forget that "who can do more, can do less," so to not fail this motto PowerShell is also capable of calling WSH, COM, and WMI objects in addition to .Net objects.

#### The Richness of Commands

PowerShell includes a set of nearly 130 commands, while CMD has only 70.
You will discover this through the next article dedicated to commands.

#### The Consistency of Commands and Their Parameters

All commands are part of the same object class. They therefore naturally inherit the same parameters. Thus, if you master the parameters of one command, you master them virtually for all others.

#### Help on Commands

PowerShell has very comprehensive built-in help. It has 3 levels of detail. The first is obtained by asking for help without specifying parameters, such as:

```
PS C:\> Get-Help Get-ChildItem
```

The second is obtained by adding the "-detailed" parameter and the last with the "-full" parameter. Moreover, the help is in English, so why not take advantage of it!?

You can also get help on a very specific topic, for example on the use of the pipe. To do this, use the following command:

```
PS C:\> Get-Help about_pipeline
```

To get the list of help topics, type the command:

```
PS C:\> Get-Help about*
```
