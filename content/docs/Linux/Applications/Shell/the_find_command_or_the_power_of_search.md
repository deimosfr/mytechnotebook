---
weight: 999
url: "/La_commande_find_ou_la_puissance_de_la_recherche/"
title: "The find command or the power of search"
description: "A comprehensive guide on using the find command in Linux for powerful file searching capabilities"
categories: ["Linux"]
date: "2011-02-09T14:06:00+02:00"
lastmod: "2011-02-09T14:06:00+02:00"
tags: ["Linux", "Command line", "File search", "System administration"]
toc: true
---

## Introduction

The find command is one of the most effective advanced commands in system administration—it requires administrator privileges.

The find command is an extremely powerful tool with about fifty options, allowing searches defined according to very fine criteria: files, directories, symbolic links... based on the name (case-sensitive or not), according to the owner, size, date, etc., in a single or multiple filesystem! It's worth adding that its options are defined according to two categories: selection and execution. This gives you an idea of the full extent of this command's power.

The syntax of the find command is unfortunately not the most common. A few tests will convince you, if you are not already convinced, to adopt it as soon as possible.

As an anecdote, it's thanks to this formidable command that I found, on my external drive, a lost file that I had been sorely missing for several years:

```bash
find /media/IOMEGA_HDD -print0 | xargs --null grep chaplin
```

As shown in the line above, the find command works wonderfully with the xargs command. Using the combined resources of these commands, I was searching through my disk's directory structure for all files containing the word "chaplin"... find was able to trace my old file inserted in a backup file dating from 2004! Hat's off to the artist!

## Proof by example

Beyond the undeniable ergonomics of graphical tools, what should make the difference is above all the efficiency and speed of the result. So, let's put the locate, find commands and the Beagle Search graphical tool to the test.

For this search example, let's make the task a bit more complicated for our competitors. The test machine includes several hard drives, therefore several partitions, several file systems, several OSs (Free only... that goes without saying!). Some partitions are dedicated only to pure and simple storage, and it's of course on these partitions that the search will have to prove itself.

### locate and Beagle

With an elementary syntax, the result of the locate command does not correspond at all to what was expected. To get satisfaction, I would have had to create my own index and certainly specify more search criteria (path, file type, etc.). But let's keep it simple.

```bash
locate artemis* /usr/share/app-install/icons/_usr_share_pixmaps_artemis.xpm
/usr/share/app-install/desktop/artemis.desktop
```

In graphical mode, Beagle doesn't do much better: only one result for a plain text file that has absolutely nothing to do with the object being searched!

### The find command... the one and only winner!

Here's the find command in action with a basic syntax, so as not to disadvantage its competitors too much:

```bash
sudo find / -name artemis* -print
[sudo] password for zarer:
/media/disk/documents_pro/2007_2008/lectures/artemis_fowl.odt
/usr/share/app-install/desktop/artemis.desktop
```

Bingo! On the first try! Wow! For me, it doesn't take more to convince me: the object being searched is found quickly, across different file systems, without entering many obscure options. Remarkable.

By adding just one option, the owner, the result would have been unmistakable:

```bash
sudo find / -name artemis* -user zarer -print
[sudo] password for zarer:
/media/disk/documents_pro/2007_2008/lectures/artemis_fowl.odt
```

*gnome-search-tool: find, locate and grep combined!*

The File Search application is the only one, to my knowledge (immediately available to a user freshly arrived on Linux with the Gnome desktop), that can compete with find. The reason:

The File Search uses the UNIX commands find, grep and locate. By default, during an elementary search, the File Search first uses the locate command, and secondly the find command, which is slower but more effective.

Well yeah! Find is behind it all!

## Introduction to the find command

As you'll have understood, this command is among the "great classics" for administration: find is a basic tool, whether in command line mode or in shell scripts.

The find command searches for objects (files, directories, links, etc.) in the directory tree that starts at the directory given as an argument, according to defined criteria, and executes an action on that object. The most commonly used selection criterion is the name of the object (file, directory, ...) and the most frequent action is displaying the access path to the searched object.

Its general syntax is as follows:

```bash
find /search_directory -option1 -option2 ...
```

Here are some simple examples of searching in all or part of the file system(s):

## Example 1

In this first example, the command must search for a specific file named my_first_file from the root directory (see the UNIX directory tree diagram) and display the result on the screen:

```bash
find / -type f -name my_first_file -print
```

## Example 2

Even the most organized users know how difficult it is sometimes to remember the exact name given to a particular file. The find command also allows the use of wildcard or substitution characters:

```bash
find / -type f -name *file -print
```

In a file (or directory) name, wildcards or substitution characters play a special role for the command interpreter: they replace a variable number of characters according to the wildcard used.

### The asterisk (*)

The star (or asterisk) replaces any number of characters (including none). It can be placed anywhere in the name. This character is often used to list files based on a part of the name.

### The question mark (?)

The question mark replaces only one character and one only. It can be placed anywhere in the file (directory) name, as in the examples below for the search for a file named file1 or file10:

```bash
find / -type f -name file? -print
find / -type f -name file?? -print
```

### Brackets ([])

Brackets are identical to the question mark except that the substitution is only made with one of the characters presented between brackets: [abc], [0-9], [!abc] for a character different from a, b and c, [!0-5], ...

## Example 3

In this third example, the command must list all directories contained in my home directory and display the result on the screen:

```bash
find /home/zarer -type d -print
```

Some useful options for advanced search

The options of the find command are divided into two categories: selection and execution.

### Options for selection criteria

The selection criteria are very numerous and apply to the type (file, directory, ...) and attributes (owner, group, permissions, creation or modification date, size, ...) of the search objects. Let's just present a few of them (see the find command manual page):

* name: performs a search based on the name (without the directories of the access path). If the defined name contains wildcards, it is preferable to surround it with quotation marks or apostrophes.
* iname: identical to -name but without differentiating between uppercase and lowercase.

```bash
find / -iname "*file*" -print
```

* type: This option specifies the type of object being searched. It must be followed by a letter that defines it. For example: d = directory, f = regular or "normal" file, l = symbolic link.

```bash
find / -type f -iname "artemis*" -print
```

The example below searches from the root directory for all files with the .odt extension and displays them:

```bash
find / -type f -iname "*.odt" -print
```

* user: File belonging to the specified user. Extremely practical for the administrator of one or more multi-user stations.
* group: File belonging to the specified group.
* mount: The search is limited to a single file system; it must not be carried out on directories located on other file systems. This is an alternative to the -xdev option, ensuring compatibility with older versions of find.

Warning! A certain order is required! The options are called "positional". Their control is carried out in the order of the command line. Otherwise, you risk getting a message like:

find: WARNING: you have specified the -mount option after an argument that is not a -type option but options are positional (-mount affects the tests specified before as well as after)

It's easy to understand that the option limiting the search to the file system alone should be placed before the name of the file (or directory) being searched and, once located, the path to the object should then be displayed:

```bash
find / -mount -type f -iname "*.pdf" -user zarer -print
```

## Options for execution criteria

Once the search criteria are specified, it's possible to ask the find command to perform different operations on the found objects. The main operation is simply to display the result of the search (the access path to the object). Again, there are also many options here. Refer to the find command manual page for details.

The selection criteria apply to the type (file, directory, ...) and attributes of the search objects. Let's just present a few of them:

* print: Displays the search result in the form of a list of access paths to objects (files, directories...) meeting the selection criteria. This is the most common use of the find command.
* exec command ;: Executes the command. All arguments to find are considered as arguments for the command line, until a ";" is encountered. The string "{}" is replaced by the name of the file being processed, in all its occurrences. These two strings may need to be protected from expansion by the shell, using the escape character ("\") or protection with apostrophes. The command is executed from the starting directory.

The search result can thus serve other commands (destruction, backup, page-by-page display, etc.), as in the example syntax of the line below which chains the rm command that will destroy without confirmation request the result of the search:

```bash
find /home -iname my_file_test -exec rm -f {} \;
```

If you have fears or doubts when executing a risky command, it's preferable in this case to use the following option:

* ok command ;: This option is identical to -exec but first asks the user. If the answer doesn't start with "y" or "Y", the command will not be executed:

```bash
find /home -iname my_file_test -ok rm -f {} \;
< rm ... /home/zarer/Desktop/my_file_test >?
```

* There is an even more suitable solution for the case where you have a very large number of files. This solution has the advantage of being more optimized:

```bash
find / -name myfiles -print0 | xargs -0 rm -f
```

The 0 allows you to overcome problems with spaces and carriage returns.

* If you want to search for files older than 7 days:

```bash
find . -type f -mtime +7 -exec ls -l {} \;
```

* And if we want to compress them:

```bash
find /path/to/files -type f -mtime +7 | grep -v \.gz | xargs gzip
```

* Now, if we want to copy files while respecting a directory structure:

```bash
find / -name "lshell*" | cpio -pduvm .
```

## Search with multiple criteria

When several options are used simultaneously, all criteria are checked. Between each option, the find command uses an implicit logical operator: "And".

It's possible to modify this rule by applying explicit operators. These operators are rarely used but it can be interesting to know them for, for example, grouping several names. The syntax of the line below allows you to search for an object whose name contains two expressions:

```bash
find / \( -iname *image* -a -iname *chaplin* \) -print
```

### Logical operators are presented in decreasing order of priority

* \( expression \): Parentheses force priority. As in algebra, what is between parentheses must be evaluated before any other operation.

* ! expression or -not expression: True if the expression is false. Simply consists of taking the expression as false. For example, you can search for a file that does not belong to a particular user.

* expression1 expression2: AND (implicit). Expression2 will not be evaluated if expression1 is false.
* expression1 -a expression2 or expression1 -and expression2: Like expression1 expression2. The -a is therefore optional.

* expression1 -o expression2 or expression1 -or expression2: OR. Expression2 is not evaluated if expression1 is true.

* expression1 , expression2: List. Both expressions will always be evaluated. The value of expression1 is ignored; the value of the list is that of expression2.

Happy searching!
