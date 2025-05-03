---
weight: 999
url: "/Initiation_à_Vi/"
title: "Introduction to Vi"
description: "A comprehensive guide to learning and using the Vi text editor on Linux systems, from basic commands to advanced features."
categories: ["Linux"]
date: "2009-09-19T22:04:00+02:00"
lastmod: "2009-09-19T22:04:00+02:00"
tags: ["Linux", "Text Editor", "Vi", "Terminal", "Commands"]
toc: true
---

## Introduction

All Linux users have Vi by default on their system, but few know how to use it. Let's discover this very useful editor.

## Why would I use Vi?

If you've ever seen someone learning to use Vi, you've probably noticed that it often involves numerous expletives and a lot of frustration... And yet, a few moments of keyboard struggle can be very beneficial in the long run.

### What is this complicated tool?

Vi, which you would pronounce [vi:aj] if you want to be trendy, is the basic editor for all Unix systems. It works in console mode, on all types of terminals, and allows you to read, write, and edit files. Vi is based on the rudimentary Ex, a line-by-line editor: it's actually the result of Ex's "visual" command, which changes the edited file from line-by-line mode to full-screen mode. In simple terms, this allows you to see the content of the file you want to work with.

It is much more powerful than other graphical or semi-graphical editors, such as Gedit or Xedit, because it has many features that don't exist elsewhere. 
The disadvantage is that to master the tools Vi offers, you need to learn a particular syntax, as Vi works by typing specific commands for the program. Plus, there are many of them, and they often don't resemble anything you've seen before.

Using Vi therefore requires tedious learning. Here we propose to give you some keys to get started, and to show you how practical this lightweight little editor is.

### How does it work?

To complicate matters a bit, Vi is a modal editor, meaning that the keys you press will not have the same effect depending on which mode you are in. There are three modes:

* insert mode: this mode allows you to insert text into your file;
* command mode: in command mode, you can move the cursor through the text and type a large number of commands that will allow multiple operations, such as searching for strings, moving them, replacing them...
* Ex or command line mode: these are the commands of the Ex editor, which don't work quite like those of Vi.

When Vi starts up, you're in command mode. You enter insert mode by typing one of the insertion commands, which we'll see below (i I a A o O). To exit, press [Esc], which brings you back to command mode. To enter Ex mode, type : or /, and commands are only executed when you press [Enter].

The key is not to confuse modes! If you're in text mode and think you're in command mode, any command you try to launch will add text; vice versa, when you think you're in text mode, some characters you think you're inserting may have unfortunate consequences if they turn out to launch deletion commands, for example! And it is true that this happens often at the beginning...

### What will it be useful for?

So why would we bother with this complicated editor? Well, simply because one day or another you'll definitely need it, and it's better to be prepared for that eventuality! Have you ever been unable to start X and had to tinker with its configuration files? Have you ever had to remotely access a computer and modify its system operation options? And without looking for complicated situations, Vi is also very practical for programming, whether for a simple shell script or for using complex languages.

## First steps with Vi

To use Vi, it's necessary to proceed step by step, otherwise you'll quickly be overwhelmed by the number of commands to memorize. Here we'll first see some basic approaches, which are sufficient for regular use of Vi, but which don't allow you to discover all its effectiveness.

### Editing, saving a file and quitting Vi

To start Vi, just type its name in a console. Opening a file is not much more complicated:

```bash
vi my_file.txt
```

Your file opens and fills almost the entire terminal, except for the last line which displays some information about it: it is on this last line that you will type your commands, and where error messages will appear.

To quit Vi, you'll need to type the command `:q`.

If my_file.txt didn't exist yet, it will be created when you launch the command `$ vi my_file.txt`. But be careful, you'll need to save it for it to be saved. For that, use the command `:w`.

So get into the habit before quitting Vi of launching a little `:wq`. If you've modified your file but don't save it when quitting, Vi will tell you with an error message:

```
File modified since last complete write; write or use ! to override.
```

To force quit without saving changes, use `:q!`.

### Moving around in a file

Explaining how to move around in a file may seem quite stupid... But keep in mind that Vi was built when keyboard navigation keys didn't exist! While you can use them, it's good to know the commands to use, just in case:

* h moves the cursor one character to the left;
* l moves the cursor one character to the right;
* j moves the cursor one line down;
* k moves the cursor one line up;
* $ moves the cursor to the last character of the line;
* 0 (zero) moves the cursor to the first character of the line.

In addition to these basic movement commands, there are others that are much more powerful. We'll only mention a few here:

* [Ctrl] + [F] allows you to move the cursor forward one screen;
* [Ctrl] + [B] allows you to move the cursor back one screen;
* G will bring your cursor to the last line of the file;
* 3G will bring the cursor to the third line of the file;
* 3w or 3b will move the cursor three words to the right or to the left;
* 3| will bring the cursor to the third column of the file.

### Writing and correcting

Now let's get to serious business and insert text into our file. Vi has a whole series of commands that allow you to switch from command mode to insert mode in different ways:

* after typing a, you will insert your character after the cursor;
* i inserts the character before the cursor;
* A inserts the text after the last character of the current line;
* I inserts the text at the beginning of the current line;
* o inserts the text in a new line below the current line;
* O inserts the text in a new line above the current line.

Correcting your text with Vi also presents many possibilities (remember to exit insert mode by pressing [Esc] to make corrections):

* x deletes the character under the cursor; combined with a repetition factor, the same command will delete the specified number of characters from the cursor (example: 4x);
* dw deletes the characters from the one under the cursor to the beginning of the next word, space included; 5dw will delete the five words following the character under the cursor;
* D deletes the entire line from the character under the cursor;
* dd deletes the current line; 4dd will delete four lines from the current line;
* rt will replace the character under the cursor with the character t; rb will replace it with b, etc.;
* C replaces the entire line from the character under the cursor with what you will type next, the insertion ending with [Esc];
* ~ will change the case of the character under the cursor: a lowercase will become an uppercase, and vice versa.

In case of an error, you also have the possibility to use the u command, which cancels the last modification made to the text.

## Configuring Vi

At this stage, and to avoid quite a few disappointments, we'll give you a trick... It is possible to know at any time which operating mode of Vi you are in... For this, switch to command mode, or check that you are there by pressing [Esc], and type the command:

```bash
:set showmode
```

Then, launch the command with [Enter]. From that moment, and for the entire current session, you will see the mode you are in on the right side of the command line: Insert or Command. Practical, isn't it? You've just modified the default configuration of Vi. And there are many other configuration options.

### Vi configuration options

The command `:set all` will give you an overview of all Vi configuration options. With the command `:set`, you will only see the default options that you have modified.
Most Vi options work like switches. To activate an option, type `:set Option_name` and to deactivate it `:set noOption_name`.
Here's an overview of useful options to start with:

* showmode / noshowmode: to activate or deactivate the display of the mode you are in;
* verbose / noverbose: by default, when Vi cannot do what you ask, the program warns you with an audible beep; the verbose option transforms this beep into a written message, very useful for beginners;
* number / nonumber: this option allows you to number the lines of the file;
* wrapscan / nowrapscan: when you are searching your file, activating this option will allow you to continue the search at the beginning of the text; otherwise the search is stopped at the end of the text;
* showmatch / noshowmatch: this option is very useful for programming: it allows, every time you close a brace or parenthesis, to highlight the corresponding opening brace or parenthesis for a few seconds.

### Making configuration changes permanent

The option modifications you make during a Vi session are not saved. When you exit the program, they are canceled, and the default configuration is reset at the next launch.

To make these modifications permanent, you need to create a `.exrc` file in your home directory:

```bash
cd
vi .exrc
```

Then indicate in the created file the options you want to start Vi with by writing their name preceded by set (without the :):

```bash
set verbose
set showmode
set number
```

Then save and quit with `:wq`.

## Advanced Features

Now, you already know enough to use Vi in a very profitable way.
However, here's a small overview of some advanced functions to show you that Vi has nothing to envy from a graphical editor.

### Copy - paste

Vi has a buffer memory whose content can then be copied to where you want in the text.

The command yy or Y copies the current line into the buffer memory, and 8yy copies eight lines from that same line. The P command inserts the buffer content before the current line, and p inserts it after.

You can also copy a series of characters...

* y$ to copy to the end of the line;
* y8w to copy eight words from the current word;
* ... and then paste it where you want with p and P.

Note also that everything that is deleted is also kept in buffer memory. So you can easily cut and paste. For example, ddp will delete the current line and copy it below the next line.

### Search - replace

Two commands allow you to search for a string of characters in a file:

* /string lets you search for occurrences of string forward from the cursor position;
* ?string lets you search for occurrences of string backward.

These two commands are activated by pressing [Enter]. The cursor then comes to rest on the first character of the found string.

To continue the search in the same direction, use the n command; to continue it in the opposite direction, use N.

If Vi doesn't find the occurrence, the message "Pattern not found" appears. When a search has been successful, the end of the search is indicated by the message "Reached end-of-file without finding the pattern" or "Reached top-fo-file without finding the pattern".

If you had previously activated the wrapscan option, the search will be carried out throughout the document, and the message announcing the end of the search becomes "Search wrapped".

Vi also allows you to replace searched strings with others. This time it's the s command that will be used. It is activated on the command line, so you introduce it with :, and launch it with [Enter]. Its operation assumes several options. The basic syntax is as follows:

```bash
:s/string_to_replace/replacement_string/
```

This command replaces the first occurrence of the specified characters in the current line. It is possible to specify the line or series of lines where you want to perform the replacement:

```bash
:1,3s/string_to_replace/replacement_string/
```

replaces the first occurrence of each line from line 1 to line 3. Note that . will designate the active line, and $ the last line of the file.

It is also possible to specify that all occurrences on the same line should be replaced by adding the g option:

```bash
:s/string_to_replace/replacement_string/g
```

For a confirmation to be requested before each replacement, add the c option:

```bash
:s/string_to_replace/replacement_string/c
```

You will then respond to the "Confirm change?" message with y (yes) or n (no).

And there you go! You can now use Vi in a quite effective way!! Practice, and you can then discover its expert features!

## Resources
- http://www.unixgarden.com/index.php/comprendre/initiation-a-vi
