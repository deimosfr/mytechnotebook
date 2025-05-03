---
weight: 999
url: "/Forcer_un_utilisateur_à_changer_son_mot_de_passe_à_la_première_connexion/"
title: "Force User to Change Password at First Login"
description: "How to force users to change their password on first login by setting account expiration"
categories: ["Linux", "Security"]
date: "2010-03-14T20:48:00+02:00"
lastmod: "2010-03-14T20:48:00+02:00"
tags: ["linux", "security", "password", "user management"]
toc: true
---

## Introduction

Indeed, I was looking for how to force a user to change their password during their first login session. Well, nothing obvious except that if we set an account to expire, the user will then be forced to change their password.

## Usage

If you are root, you can specify the user whose account you want to expire as follows:

```bash
chsh -s /bin/MySecureShell username
```

And otherwise, a user can change their shell themselves like this:

```bash
chsh -s /bin/MySecureShell
```
