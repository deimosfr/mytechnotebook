---
weight: 999
url: "/Mise_en_place_d'un_serveur_et_client_Git/"
title: "Setting Up a Git Server and Client"
description: "This guide covers installation and configuration of Git server and client, with details on managing repositories, branches, tags, and advanced features."
categories: ["Nginx", "Debian", "Linux"]
date: "2015-02-02T20:33:00+02:00"
lastmod: "2015-02-02T20:33:00+02:00"
tags: ["git", "version control", "repository", "server", "client"]
toc: true
---

![Git Logo](/images/git-logo.avif)

## Introduction

[Git](https://en.wikipedia.org/wiki/Git) is a distributed version control system. It's free software created by Linus Torvalds, the creator of the Linux kernel, and distributed under GNU GPL version 2.

Like BitKeeper, Git doesn't rely on a centralized server. It's a low-level tool, designed to be simple and highly efficient, whose main task is to manage the evolution of a file tree.

Git indexes files based on their checksum calculated with the SHA-1 function. When a file isn't modified, the checksum doesn't change, and the file is stored only once. However, if the file is modified, both versions are stored on the disk.

Git wasn't initially a version control system in the strict sense. Linus Torvalds explained that "in many ways, you can consider git as a filesystem: it allows associative addressing, and has the notion of versioning, but most importantly, I designed it by solving the problem from a filesystem specialist's perspective (my job is kernels!), and I had absolutely no interest in creating a traditional version control system." It has since evolved to incorporate all the features of a version control system.

Git is considered to be highly performant, to the point where some other version control systems (Darcs, Arch) that don't use databases have shown interest in Git's file storage system for their own operations. However, they would continue to offer more advanced features.

## Installation

First, let's install Git on the server:

```bash
aptitude install git git-core
```

If you want to make the server accessible via a `git://` address, you'll also need to install:

```bash
aptitude install git-daemon-run
```

## Configuration

### Server

First, we'll create the repository:

```bash
$ mkdir myproject
$ cd myproject
$ git init
Initialized empty Git repository in .git/
```

Now that our project is created, let's add a small file to make sure everything is working correctly:

```bash
touch test
```

Then we'll add it and commit it:

```bash
 $ git add .
 $ git commit -m "test"
 Created initial commit c491bd6: test
 1 files changed, 1 insertions(+), 0 deletions(-)
 create mode 100644 afile
```

To check the status at any time, run:

```bash
git status
```

Next, we clone this repository:

```bash
$ cd ..
$ git clone --bare myproject myproject.git
Initialized empty Git repository in /var/cache/git/myproject.git/
0 blocks
$ ls myproject.git
branches  config  description  HEAD  hooks  info  objects  refs
```

- clone: creates a copy of the repository locally
- --bare: creates a copy containing only the info from the myproject folder
- myproject.git: folder that git has created for us

#### The git protocol

If you want to make git accessible via its protocol (`git://`), you need to configure the following file like this:

```bash {linenos=table}
#!/bin/sh
exec 2>&1
echo 'git-daemon starting.'
exec chpst -ugitdaemon \
  "$(git --exec-path)"/git-daemon --verbose --base-path=/var/cache /var/cache/gi
```

Then for each project you want to make accessible from outside, you'll need to place a 'git-daemon-export-ok' file:

```bash
touch /var/cache/git/myproject.git/git-daemon-export-ok
```

Restart the daemon and it will now be accessible on port 9418.

### Client

You should first [install Git as done on the server](#installation).

Then with your current user, inform git about your identity (make sure your user exists on the server):

```bash
git config --global user.email "xxx@mycompany.com"
git config --global user.name "deimosfr"
```

This should create the following in your ~/.gitconfig:

```bash
[user]
       email = xxx@mycompany.com
       name = Deimos
```

If you want to have different identities for different repositories, simply remove global (by placing yourself in the repository in question):

```bash
git config user.email "xxx@mycompany.com"
git config user.name "Deimos"
```

Now we will retrieve our project from our git. We're using the SSH method here, but git can also work with rsync, http, https, and git-daemon.

#### With SSH

Here's the SSH method:

```bash
 $ git clone ssh://username@host/var/cache/git/myproject
 Initialized empty Git repository in /var/cache/git/myproject/.git/
 Password:
 remote: Counting objects: 1, done.
 remote: Total 1 (delta 0), reused 0 (delta 0)
 Receiving objects: 100% (1/1), done.
 $ ls
 myproject
 $ ls myproject
 test
```

- For repository updates, you'll need to use the push option. First, changes must be made locally (that's the advantage of git), such as an add and a commit (we'll see this just after):

```bash
git push ssh://username@host/var/cache/git/myproject.git
```

or

```bash
git push ssh://username@host/var/cache/git/myproject.git master
```

'master' here corresponds to the branch name we're interested in (see below for using branches).

- To add a file to an SSH repository:

```bash
git remote add test ssh://username@host/var/cache/git/myproject.git
git commit -a
git push test
```

In the future, to avoid dealing with SSH and its special commands, add this to your ~/.gitconfig file:

```bash
[remote "myproject"]
      url = ssh://username@host/var/cache/git/myproject/
```

#### With Git-daemon

Git will refuse to synchronize if the folder in question doesn't contain a file called '**git-daemon-export-ok**'. Once this file is created, the folder will be accessible to everyone:

```bash
git clone git://serveur.git/git/myproject
```

or

```bash
git clone git://serveur.git/git/myproject.git
```

We can use the following options:

- '**--export-all**': this option no longer requires the 'git-daemon-export-ok' file
- '**--user-path=gitexport**': this option will allow URLs on user home directories. So `git://deimos-laptop/~deimos/myproject.git` will point to /home/deimos/gitexport/myproject.git

#### On HTTP (with Nginx)

I spent quite a bit of time getting Git over http(s) and Gitweb to coexist, but it's working now.

{{< alert context="info" text="Prefer the <a href='./Gitweb_:_Installation_et_configuration_d'une_interface_web_pour_git#Nginx'>Gitweb</a> method alone if you don't need git over http(s)." />}}

Here's the method I used:

```bash {linenos=table}
server {
    listen 80;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name git.deimos.fr;
    root /usr/share/gitweb/;

    access_log /var/log/nginx/git.deimos.fr_access.log;
    error_log /var/log/nginx/git.deimos.fr_error.log;

    index gitweb.cgi;

    # Drop config
    include drop.conf;

    # Git over https
    location /git/ {
        alias /var/cache/git/;
        if ($scheme = http) {
            rewrite ^ https://$host$request_uri permanent;
        }
    }

    # Gitweb
    location ~ gitweb\.cgi {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass  unix:/run/fcgiwrap.socket;
    }
}
```

Here I have my git over https working (http is redirected to https) as well as my gitweb since everything that matches gitweb.cgi is caught. Now, for the git part, we'll need to authorize the repositories we want. For that, we'll need to rename a file in our repository and run a command:

```bash
cd /var/cache/git/myrepo.git
hooks/post-update{.sample,}
su - www-data -c 'cd /var/cache/git/myrepo.git && /usr/lib/git-core/git-update-server-info'
```

Replace www-data with the user who has rights to the repository. Use www-data so that nginx has the rights. Then you have permissions to clone:

```bash
git clone http://www.deimos.fr/git/deimosfr.git deimosfr
```

#### Configuring a proxy

If you need to go through a proxy with Git, you can do it like this (for its global part):

```bash
git config --global http.proxy http://proxy:8080
```

To check the current settings:

```bash
git config --get http.proxy
```

If you need to remove this configuration:

```bash
git config --global --unset http.proxy
```

## Usage

### Updating your local repository

If you want to update your local repository from the server, it's simple:

```bash
git pull
```

or

```bash
git pull ssh://username@host/var/cache/git/myproject/
```

### Adding files

To add files in bulk, use this command:

```bash
git add *
```

### Committing changes

If you want to both log and commit at the same time (otherwise leave out -m...):

```bash
git commit -a -m "My first addition"
```

This will update the repository locally (don't forget to push afterwards if you're working with a remote server).

### Changing the default editor and viewer

You can force vim and less for example by entering these lines:

```bash
git config --global core.editor vim
git config --global core.pager "less -FSRX"
```

### Getting logs

You can check the logs of previous commits at any time with this command:

```bash
git log
```

If you want to get logs for a specific file:

```bash
git log test
```

If your file has been renamed, use the **--follow** option to see its new name as well.
To see logs from the last 2 days:

```bash
git log --since="2 days ago"
```

### Searching

You can search for content in files like this:

```bash
git grep deimos *
```

### Ignoring files

You may want to exclude certain files such as vim temporary files. You'll need to create a .gitignore file at the root of the working copy, which will generally be added to the repository itself, or in the .git/info/exclude file. Here's an example:

```bash
*~
*.o
```

### Making sure there's nothing to commit

You can ensure there's nothing to commit like this:

```bash
> git stash
No local changes to save
```

### Configuring aliases

You can configure aliases for longer commands. For example, if you're used to [SVN](Installation_et_configuration_d'un_repository_SVN.html):

```bash
git config --global alias.co 'checkout'
git config --global alias.ci 'commit -a -m'
```

Here 'git co' corresponds to the 'git checkout' command. This configuration (whether global or not) can be located in 2 places:

- ~/.gitconfig to benefit from it in all your repositories.
- .git/config of a project to restrict its access to that single project.

### Revert: canceling a commit

Git revert allows you to cancel the previous commit by introducing a new one. First, we'll do a git log to get the commit number:

```bash
$ git log
commit : bfbfefa4d9116beff394ad29e2721b1ab6118df0
```

Then we'll do a revert to cancel the old one:

```bash
git revert bfbfefa4d9116beff394ad29e2721b1ab6118df0
```

In case of conflict, simply make manual changes, do another add, and then commit.

Note: you can also use the first few characters of the commit identifier (SHA-1) such as bfbfe. Git will automatically complete it (as long as another commit doesn't start with the same characters).

### Reset: deleting all commits from an old one

Git reset will completely remove all commits made after a certain version:

```bash
git reset bfbfefa4d9116beff394ad29e2721b1ab6118df0
```

Git reset allows you to rewrite history. You can verify this with a 'git log'. This can have serious consequences if you're working with others on a project. Not only because the local repositories of other people will need to be resynchronized, but also because there can be many conflicts. It is therefore advised to use this command with caution.
Then you must redo a '**git add**' of your files, then a '**git commit**'.

If you want the local working copy to reflect the repository, add the --hard option during reset:

```bash
git reset **--hard** bfbfefa4d9116beff394ad29e2721b1ab6118df0
```

In case you're having a terrible day and deleted months of work and want to go back, git has preserved the previous state, so do:

```bash
git reset --hard ORIG_HEAD
```

HEAD means the last commit by default.

If you get a message like:

```
! [rejected]        master -> master (non-fast forward)
```

And you really want to force the commit to look exactly like what you have on your machine:

```bash
git push origin +master
```

### Restoring a file from an old commit

You can choose to restore only a particular file from an old commit via the checkout command:

```bash
git checkout bfbfe test
```

### Deleting a file

There are 2 methods; my preferred one, which I find simple:

```bash
git rm my_file
```

The 2nd method consists of doing a standard rm, then at the next git commit -a, it will detect that this file no longer exists and will remove it from the repository.

### Moving a file

To move a file within your git repository, use this command:

```bash
git mv test1 test2
```

### Listing the repository contents

You can use this command to list the contents of a repository:

```bash
git ls-files
```

Note: empty directories will not be displayed

### Creating a repository archive

To create an archive in tar gzip format, we'll use the archive option:

```bash
git archive --format=tar --prefix=version_1/ HEAD | gzip > ../version_1.tgz
```

- --prefix: allows giving a folder that will be created upon decompression.
- HEAD: allows specifying the commit (here the 1st)

### Viewing a specific file from a commit

It's possible to view a file (README) from a commit:

```bash
git show 291d2ee31feb4a318f77201dea941374aae279a5:README
```

or see all diffs at once if you don't specify the file:

```bash
git show 291d2ee31feb4a318f77201dea941374aae279a5
```

### Rewriting history

If you forgot to do something in a previous commit and want to go back, you need to do an interactive rebase by indicating the earliest commit you want to modify:

```bash
git rebase -i cbde26ad^
```

In the editor, change 'pick' to 'edit' on the line(s) you want to modify. Save and exit. Make all your changes and then validate the commit:

```bash
git commit -a --amend
or
git commit -a --amend --no-edit
```

Then, to move on to the next commit:

```bash
git rebase --continue
```

Until it's finished.

### Making diffs between different versions

It's very easy with git to make diffs between different versions. For example, from HEAD to 2 versions back:

```bash
git diff HEAD~2 HEAD
```

### Modifying the text of the last commit

To see the state of the last log, we can do:

```bash
git log -n 1
```

To modify the text of this last log, we'll use the amend option:

```bash
git commit -a --amend
```

This allows you to change just the text; there won't be a new commit number. Only the identification number will change.

### Modifying the date of the last commit

It's possible to modify the date of the last commit:

```bash
GIT_COMMITTER_DATE="`date`" git commit --amend --date "`date`"
```

Or you can specify the date explicitly:

```bash
GIT_COMMITTER_DATE="Fri Nov 17 12:00:00 CET 2014" git commit --amend --date "Fri Nov 17 12:00:00 CET 2014"
```

### Modifying the author of commits

When using multiple Git accounts on the same machine, it's easy to make a mistake when cloning a repository and not setting the right username and email address. When you realize it, it's too late and you need to rewrite part of the history to correct the information:

```bash {linenos=table,hl_lines=[3,4,5,6]}
> git filter-branch --env-filter '
oldname="old username"
oldemail="old email address"
newname="new username"
newemail="old username address"
[ "$GIT_AUTHOR_EMAIL"="$oldemail" ] && GIT_AUTHOR_EMAIL="$newemail"
[ "$GIT_COMMITTER_EMAIL"="$oldemail" ] && GIT_COMMITTER_EMAIL="$newemail"
[ "$GIT_AUTHOR_NAME"="$oldname" ] && GIT_AUTHOR_NAME="$newname"
[ "$GIT_COMMITTER_NAME"="$oldname" ] && GIT_COMMITTER_NAME="$newname"
' HEAD
```

### Synchronizing only part of a repository

If you want to have only part of a repository (also called sparse), you'll need to first retrieve the entire repository, then specify what interests you. Only what interests you will remain:

```bash
git clone ssh://deimos@git/var/cache/git/git_deimosfr .
git config core.sparsecheckout true
echo configs/puppet/ > .git/info/sparse-checkout
git read-tree -m -u HEAD
```

Here in the git_deimosfr repository, only the "configs/puppet/" folder interests me. You can add multiple folders to the ".git/info/sparse-checkout" file line by line if you want to keep multiple files.

### Using an external git repository within a git

If, for example, you already have a Git in which you want to integrate another Git, but external, you'll need to use submodules. For my part, I have a [puppet](./puppet_:_solution_de_gestion_de_fichier_de_configuration.html) folder with lots of modules inside. Some of these modules weren't created by me and are maintained by other people who have a Git. I want some modules to point to external gits. Here's how to proceed:

- I add the external source to this git:

```bash
> git submodule add git://git.black.co.at/module-common configs/puppet/modules/common
Cloning into configs/puppet/modules/common...
remote: Counting objects: 423, done.
remote: Compressing objects: 100% (298/298), done.
remote: Total 423 (delta 155), reused 203 (delta 70)
Receiving objects: 100% (423/423), 53.96 KiB, done.
Resolving deltas: 100% (155/155), done.
```

Now you just need to commit and push if you want to make this configuration valid on the server.

- Now let's consider the case of a client who does a pull. They will get your entire tree structure, but not the external links; for that, they will need to execute the following commands:

```bash
> git submodule init
Submodule 'configs/puppet/modules/common' (git://git.black.co.at/module-common) registered for path 'configs/puppet/modules/common'
> git submodule update
Cloning into configs/puppet/modules/common...
remote: Counting objects: 423, done.
remote: Compressing objects: 100% (298/298), done.
remote: Total 423 (delta 155), reused 203 (delta 70)
Receiving objects: 100% (423/423), 53.96 KiB, done.
Resolving deltas: 100% (155/155), done.
Submodule path 'configs/puppet/modules/common': checked out 'ba28f3004d402c250ef3099f95a1ae13740b009f'
```

If during a clone you want the submodules to also be taken, you need to add the "--recursive" option. Example:

```bash
git clone --recursive git://git.deimos.fr/git/git_deimosfr
```

### Branches

#### Creating a branch

If, for example, the current version you have on git suits you and you want to make it stable, you need to use branches and create a new branch as a development branch:

```bash
git branch devel
```

#### Sending a branch to the server

To send a branch to the server:

```bash
git push origin <my_branch>
```

#### Listing branches

Now, the devel branch is created; to check:

```bash
$ git branch -a
  devel
* master
```

The '\*' indicates the current branch (in use).

When using remote branches, you can list them with this option:

```bash
$ git branch -r
  devel
* master
```

#### Changing branches

To change branches, simply use the checkout option:

```bash
$ git checkout devel
Switched to branch "devel"
```

Let's check:

```bash
$ git branch
* devel
  master
```

If we want to go back to the previous branch, we can at any time do:

```bash
$ git checkout -f master
$ git branch
* master
  devel
```

The -f option corresponds to the --hard option of reset which allows synchronizing local changes made to the repository.

If you're working on a remote branch, do this for example:

```bash
$ git checkout deimos/myproject -b master
$ git branch
* master
  devel
```

#### Getting a branch from a remote server

Once your "pull" is done, change branches like this (you can name it differently if you wish):

```bash
git checkout -b <new_branch> origin/<new_remote_branch>
```

- new_branch: the name of my new local branch
- new_remote_branch: the name of the new remote branch (server side)

#### Deleting a branch

To delete a branch, nothing simpler:

```bash
git branch -d devel
```

If this works, perfect; if you want to force in case of a problem, use '-D' instead of '-d'. To apply the changes to the remote server:

```bash
git push origin :devel
```

This will delete the devel branch on the remote server.

#### Merging branches

The merge option consists of completely merging 2 branches. Like for example merging a devel branch with a stable one so that devel becomes stable:

```bash
git checkout master
git merge devel
```

#### Rebase: Applying a patch to multiple branches

Imagine we've been working on the devel branch but a bug in master has been discovered. It's normal that if master gets the fix, the devel branch should get it too. We'll place ourselves in the master branch, apply the patch, then merge the patch with the devel branch. **We place ourselves on the devel branch, then execute these commands:**

```bash
git checkout devel
git rebase master
```

It's possible to do it interactively with the -i argument:

```bash
git rebase -i master
```

#### Creating a branch during a restoration

You can create a branch [when restoring a file](#restoring-a-file-from-an-old-commit) by first placing yourself in the branch and using the -b option:

```bash
$ git checkout bfbfefa4d9116beff394ad29e2721b1ab6118df0 -b version_2
Switched to a new branch "version_2"
```

### Tags

Tags are very useful for marking a specific version. For example, to release version 1.0, you can create a tag and reuse it later to download this particular version. You can tag anything; it allows for easy referencing.

#### Creating a Tag

To create a tag, it's easy:

```bash
git tag '1.0'
```

{{< alert context="info" text="It's possible to sign tags via GPG" />}}

#### Listing tags

To list available tags:

```bash
git tag
```

#### Deleting a tag

To delete a tag:

```bash
git tag -d v0.1
git push origin :refs/tags/v0.1
```

#### Sending tags to the server

Once the tags are created locally, you can push them to the server:

```bash
git push --tags
```

### Managing conflicts

You can manage conflicts interactively:

```bash
git mergetool
```

### Automatic resolution of already seen conflicts

It's possible to ask git to automatically resolve problems it has already encountered:

```bash
git config --global rerere.enabled 1
```

### Using git to find a bug

Imagine a sneaky bug suddenly appears and disrupts your project. After investigations worthy of an English pipe smoker, you discover that the nasty malfunction has only appeared recently. It's clearly a regression!

Sure, you could waste time fixing this regression the classic way, by diving into dusty code full of virtual cobwebs, but it would be much simpler to discover exactly when it appeared, which commit corresponds to it, to have a very precise idea of its cause.

This is exactly what git-bisect allows you to do, in a particularly intuitive way. Tell it a commit with the malfunction (HEAD), and a past commit, any one, that doesn't have it.

Git then places you in your development history, right in the middle between the two commits. After checking, you tell git whether the regression appears or not. And we start again. Each time, we divide by 2 the range of commits that potentially introduced the bug, until we catch the culprit. Great, isn't it?

Let's put it into practice:

```bash
> git bisect start
> git bisect bad <bad_version>
> git bisect good <good_version>

Bisecting: 8 revisions left to test after this
```

8 commits potentially introduced the regression. We indicate whether the current commit is correct or not:

```bash
> git bisect good # or git bisect bad
Bisecting : 4 revisions left to test after this
```

And we continue, again and again, until we find the bad commit that caused this wasted time (but it could have been worse):

```
37b4745bf75e44638b9fe796c6dc97c1fa349e8e is first bad commit
```

At any point, you can get a history of your journey:

```bash
git bisect log
```

If at a specific moment, you don't want to test a particular commit, for whatever reason, use the command:

```bash
git bisect skip
```

All this is fine and good, but there's better (no!? yes!). Suppose you're a fan of TDD. You surely have a script that automatically tests whether the current code is good or not. It's then possible to automate the search:

```bash
git bisect start HEAD <bad_commit> --
git bisect run script
```

Go get a coffee, git is working for you. Note: the script must return 0 if the code is correct, 1 if it's incorrect, and 125 if it's untestable (git skip). During bisection, git creates a special dedicated branch. Don't forget to switch back to your development branch when you've caught the regression. You can do it simply by typing:

```bash
git bisect reset
```

So, with all this, no more tedious regression searches.

### Knowing commits line by line

This is the solution for slapping the fingers of someone who did a commit. Among other things, you have the ability to know line by line who wrote what:

```bash
> git blame <file>
^a35889c (Deimos 2010-04-20 17:36:29 +0200   1) <?php
^a35889c (Deimos 2010-04-20 17:36:29 +0200   2) class DeleteHistory extends SpecialPage
^a35889c (Deimos 2010-04-20 17:36:29 +0200   3) {
724f724d (Deimos 2011-06-14 23:15:40 +0200   4) 	function __construct()
724f724d (Deimos 2011-06-14 23:15:40 +0200   5) 	{
724f724d (Deimos 2011-06-14 23:15:40 +0200   6) 		// Need to belong to Administor group
724f724d (Deimos 2011-06-14 23:15:40 +0200   7) 		parent::__construct( 'DeleteHistory', 'editinterface' );
724f724d (Deimos 2011-06-14 23:15:40 +0200   8) 		wfLoadExtensionMessages('DeleteHistory');
```

### Moving a folder and its history to another repository

Sometimes you may need to recreate another empty repository to contain a folder from another repository. This was my case for the [DeleteHistory](https://www.mediawiki.org/wiki/Extension:DeleteHistory) extension for Mediawiki that I created. At first, I had a repository called 'mediawiki_extensions' where I had 2 extensions. And then with time and advancing Mediawiki versions, it was preferable to separate by plugins while keeping the history. So I'll explain how I did it (I followed this documentation).

We'll work locally with my old repository which we'll call oldrepo and my new repository newrepo (how original).

{{< alert context="warning" text="Create a temporary repository if you already have one because it will be deleted at the end" />}}

So let's clone this repo:

```bash
git clone git://git.deimos.fr/oldrepo.git
git clone git://git.deimos.fr/newrepo.git
```

Then we'll filter for the folder we're interested in, let's call it folder2keep:

```bash
git filter-branch --subdirectory-filter folder2keep -- -- all
```

Now, all the files and folders from folder2keep are at the root of our repository. If you wish, you can create a folder and place everything inside, but this is not mandatory:

```bash
mkdir new_directory/
git mv * new_directory/
```

Replace \* with all elements you want to place inside and commit:

```bash
git commit -m "Collected the data I need to move"
```

Now, we'll go into our famous folder and make a local reference between the 2:

```bash
cd ../newrepo/
git remote add oldrepo ../oldrepo/
```

Then we'll bring in the sources, create the main branch, and merge everything:

```bash
git fetch oldrepo
git branch oldrepo remotes/oldrepo/master
git merge oldrepo
```

Now we'll transfer the history and do some cleaning:

```bash
git remote rm oldrepo
git branch -d oldrepo
git push origin master
```

And there you go, the temporary repository (old) can now be deleted.

### Hooks

There are hooks in Git allowing pre or post processing. Here's a use case for performing an action (running a command via ssh) when a tag arrives on the server. The idea is to be able to deploy a new version of software on X servers based on a tag (example: prod or preprod). Here are the options that will be necessary:

1. Create a hook in the repository
2. Create a dedicated user if you use Gitlab/GitHub/Gitolite
3. Generate a private SSH key with the 'git' user
4. Copy via ssh-copy-id to remote servers the 'git' user key (usually www-data)
5. Create a deployment script and set the proper permissions
6. Clone the repository on all necessary servers

```bash {linenos=table,hl_lines=[3,4,5,6,9,12]}
#!/bin/bash

# Create your environment with DNS/IP servers separated by spaces
name=( list array )
prod=( X.X.X.X Y.Y.Y.Y )
preprod=( Z.Z.Z.Z )

# Folder to deploy on client side
folder='/var/www/cloned_repository'

# Set SSH remote username
ssh_username='www-data'

# Log file
logfile='/var/log/git-deploy-hook.log'

########################################################################
# get infos
read oldrev newrev refname

# vars
tag=`echo $refname | cut -d/ -f3`
environment=`echo $tag | cut -d- -f1`

# function to call distant git deployments
deploy_new_version() {
    declare -a servers=("${!1}")

    echo $(date +'%Y/%m/%d - %H-%M-%S ') "Begin deployment hook for $environment env" >> $logfile
    for server in $(seq 0 $((${#servers[@]} - 1))) ; do
        echo "Deploying $tag" >> $logfile
        ssh $ssh_username@${servers[$server]} /usr/bin/git_deploy.sh $tag $folder 2>&1 >> $logfile
    done
    echo $(date +'%Y/%m/%d - %H-%M-%S ') "Deployment hook finished for $environment env" >> $logfile
}

# check if a tag has been pushed
if [[ -n $(echo $refname | grep -Eo "^refs/tags/$environment") ]]; then
    echo "Git hook is called for $environment environment"
    if [ ${!environment[0]} ] ; then
        # Deploy the correct environment
        deploy_new_version $environment[@]
    else
        echo "Error, $environment is not recognized as an existing environment"
        exit 1
    fi
fi
```

You need to modify the environments with the ones you want so that when a tag arrives with the name of the tag in question, an action is triggered behind it. For example, to deploy on production servers, you'll need to put a tag like "prod-v1.0".

Switch to the git user (or the one who has the rights) and copy the public key to all the servers on which the tags will need to act:

```bash
su - git
ssh-copy-id www-data@x.x.x.x
```

Then create the deployment script:

```bash
#!/bin/bash
logfile='/var/log/git-deploy-hook.log'
echo $(date +'%Y/%m/%d - %H-%M: ') 'Begin deployment hook' >>$logfile
cd $2
git fetch origin -v >> $logfile 2>&1
git fetch origin -v --tags >> $logfile 2>&1
git stash drop >> $logfile 2>&1
git checkout -f tags/$1 >> $logfile 2>&1
echo $(date +'%Y/%m/%d - %H-%M: ') 'Deployment finished' >> $logfile
```

You can add any commands you want. You obviously need to copy this script to all target servers and give it execution rights.

Create the log file and set the right permissions:

```bash
touch /var/log/git-deploy-hook.log
chow www-data. /var/log/git-deploy-hook.log
chmod 755 /usr/bin/git_deploy.sh
```

Then clone the repositories with the final users:

```bash
cd /var/www
git clone git@gitlab/cloned_repository.git
chown -Rf www-data. /var/www/cloned_repository
```

And there you go, all that's left is to push a tag on a commit:

```bash
git tag -a prod-v1.0 -m 'production version 1.0' 9fceb02
git push --tags
```

## Utilities

### Gitk

Gitk provides a graphical version of your git status (branches etc.).

![Gittk.jpg](/images/gittk.avif)

## FAQ

### warning: You did not specify any refspecs to push, and the current remote

If you get this kind of message when you do a push:

```
warning: You did not specify any refspecs to push, and the current remote
warning: has not configured any push refspecs. The default action in this
warning: case is to push all matching refspecs, that is, all branches
warning: that exist both locally and remotely will be updated.  This may
warning: not necessarily be what you want to happen.
warning:
warning: You can specify what action you want to take in this case, and
warning: avoid seeing this message again, by configuring 'push.default' to:
warning:   'nothing'  : Do not push anything
warning:   'matching' : Push all matching branches (default)
warning:   'tracking' : Push the current branch to whatever it is tracking
warning:   'current'  : Push the current branch
```

Simply run this command, selecting what corresponds best to you:

```bash
git config --global push.default matching
```

## References

http://www.kernel.org/pub/software/scm/git/docs/user-manual.html
http://git.wiki.kernel.org/index.php/GitFaq
http://alexgirard.com/git-book/index.html
[How To Install A Public Git Repository On A Debian Server](/pdf/how_to_install_a_public_git_repository_on_a_debian_server.pdf)
[Git it](/pdf/gitit.pdf)
http://www.unixgarden.com/index.php/administration-systeme/git-les-mains-dans-le-cambouis
http://stackoverflow.com/questions/1811730/how-do-you-work-with-a-git-repository-within-another-repository
http://chrisjean.com/2009/04/20/git-submodules-adding-using-removing-and-updating/
