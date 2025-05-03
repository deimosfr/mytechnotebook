---
weight: 999
url: "/Installer_Rsync_sur_Windows/"
title: "Installing Rsync on Windows"
description: "A guide on how to install and configure Rsync on Windows using Cygwin, including SSH setup and BackupPC configuration."
categories: ["Windows", "Backup"]
date: "2008-02-26T15:24:00+02:00"
lastmod: "2008-02-26T15:24:00+02:00"
tags: ["Windows", "Cygwin", "Rsync", "SSH", "BackupPC"]
toc: true
---

## Installing Cygwin on Windows

Download: [https://www.cygwin.com/](https://www.cygwin.com/)

During installation, choose a download site in France (if available).

Packages to install: rsync, openssh, and cygwin (cygrunsrv, cygwin)

## Adding a System Variable on Windows (in Environment Variables)

```
variable: CYGWIN    value: ntsec tty
```

And add to the end of the path: `c:\cygwin\bin`

## Configuring sshd on Cygwin

Execute:

```bash
ssh-host-config
```

Answer "yes" to all questions and set a password when prompted.

Execute the sshd service:

```bash
cygrunsrv -S sshd
```

To verify it's working, try a local telnet connection on port 22. You can also check if the sshd service has been added to Windows services.

## Script for BackupPC

```bash
$Conf{XferMethod} = 'rsync';
$Conf{RsyncShareName} = [
  '/cygdrive/c/Shoreline Data'
];
$Conf{RsyncUserName} = 'administrateur';
$Conf{BlackoutPeriods} = [
  {
    'hourEnd' => '19.5',
    'weekDays' => [
      '1',
      '2',
      '3',
      '4',
      '5'
    ],
    'hourBegin' => '7'
  }
];
$Conf{FullPeriod} = '30';
$Conf{FullKeepCnt} = [
  '6',
  '0',
  '1'
];
$Conf{IncrAgeMax} = '30';
$Conf{IncrKeepCntMin} = '1';

$Conf{RsyncClientCmd} = '$sshPath -q -x -l Administrateur $host $rsyncPath $argList+';
$Conf{RsyncClientRestoreCmd} = '$sshPath -q -x -l Administrateur $host $rsyncPath $argList+';
```
