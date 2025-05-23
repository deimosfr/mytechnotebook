---
title: "Create a clean Git/GitHub pages branch"
slug: create-a-clean-gitgithub-pages-branch/
description: "Prepare a clean GitHub pages branch for your project"
categories: ["Git", "GitHub"]
tags: ["git", "github"]
---

## Introduction

This document provides a step-by-step guide to creating a clean GitHub Pages branch for your project. This is particularly useful when you want to deploy your project using GitHub Pages without including unnecessary files or branches.

## Hands on

This is the sequence of steps to follow to create a root gh-pages branch:

```bash
cd /path/to/repo-name
git symbolic-ref HEAD refs/heads/gh-pages
rm .git/index
git clean -fdx
echo "My GitHub Page" > index.html
git add .
git commit -a -m "First pages commit"
git push origin gh-pages
```

Why do we need to do all this, instead of just calling git branch gh-pages. Well, if you are at master and you do git branch gh-pages, gh-pages will be based off master.

Here, the intention is to create a branch for github pages, which is typically not associated with the history of your repo (master and other branches) and hence the usage of git symbolic-ref. This creates a "root branch", which is one without a previous history.

Note that it is also called an orphan branch and git checkout --orphan will now do the same thing as the git symbolic-ref that was being done before.

!!! warning

    Be careful with git clean -fdx, it will wipe out the files of the folder.

## Resources

- https://gist.github.com/ramnathv/2227408
