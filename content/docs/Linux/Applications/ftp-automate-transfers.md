---
weight: 999
url: "/FTP_\\:_Automatiser_des_transferts/"
title: "FTP: Automate Transfers"
description: "How to automate FTP transfers using shell scripts and command line tools like ftp and ncftp."
categories: ["Network", "FTP", "Automation"]
date: "2008-07-22T07:30:00+02:00"
lastmod: "2008-07-22T07:30:00+02:00"
tags: ["ftp", "automation", "script", "ncftp", "file transfer"]
toc: true
---

## Introduction

It's sometimes convenient to automate certain tasks like uploading files or other operations.

## FTP

With the universal ftp command, here's an example that you can place in a shell script:

```bash
transfertFile()
{
   inputFile=$1
 
   ftp -n <<end
      prompt
      open $Hostname $Port
      user $Login $Password
      ascii
      put $inputFile
      bye
end
}
```

## NCFTP

With the ncftp utility, it's even simpler since it works with just one line:

```bash
ncftpput -u $LOGIN -p $PASSWORD $ADDRESS DESTINATION-DIRECTORY FILE-TO-UPLOAD
```
