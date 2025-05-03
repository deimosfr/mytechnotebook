---
weight: 999
url: "/CVS_\\:_Utilisation_de_CVS/"
title: "CVS: Using CVS"
description: "A guide to the basics of using CVS (Concurrent Versions System) for version control, including project management, adding and removing files, and identifying changes."
categories: ["Development", "Version Control"]
date: "2007-05-20T20:37:00+02:00"
lastmod: "2007-05-20T20:37:00+02:00"
tags: ["cvs", "version control", "development", "source code management"]
toc: true
---

## Introduction

**CVS**, an acronym for **Concurrent Versions System**, is free software (GPL license) for version management, successor to SCCS. Although it is still widely used in the open source software domain, it is now obsolete with the arrival of its successor Subversion. Since it helps sources converge towards the same destination, we say that CVS manages concurrent versions. It can function both in command line mode and through a graphical interface. It consists of client modules and one or more server modules for exchange areas.

It's worth noting that there are also decentralized software options like Bazaar, Darcs, Git, or Monotone, all under Open Source licenses.

Among the most popular graphical interfaces, notable ones include Cervisia for Linux and TortoiseCVS for Windows.

## Prerequisites

For prerequisites, you need a few things in your environment. I strongly recommend adding them to your shell load file (e.g., ~/.bashrc, ~/.zshrc):

```bash
export CVS_RSH=/usr/bin/ssh
export CVSROOT=:ext:xxx@mycompany.com:/var/lib/cvs
```

In the first line, we need to indicate the transport method for CVS. In this case, it's SSH.
For the second, we indicate the hostname where the CVS server is located, as well as the folder where the repository is located.

Reload your shell and you're good to go.

## Usage

### Projects

* Downloading a project:

```bash
cvs checkout project_name
```

* Updating a project after updates (**this does not upload our updates**):

```bash
cvs update
```

* Creating a new project (make sure you're in the concerned project first!):

```bash
cvs import project_name creator release
```

* Destroying a project:

```bash
cvs release -d project_name
```

### Adding Files

* Adding files:

```bash
cvs add file_name
```

* Updating files:

```bash
cvs commit files_to_update
```

* Updating files with comments at the same time:

```bash
cvs commit -m "My comments" files_to_update
```

### Removing Files

To remove files, you need to:

* Delete the file on your local machine:

```bash
rm file_name
```

* Remove the file from cvs:

```bash
cvs remove file_name
```

* Commit the changes:

```bash
cvs commit file_name
```

### Identification

* View differences between server modifications and your own:

```bash
cvs diff
```
