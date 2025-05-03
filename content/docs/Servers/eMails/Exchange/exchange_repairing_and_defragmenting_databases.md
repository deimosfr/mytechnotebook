---
weight: 999
url: "/Exchange_\\:_Réparer_et_défragmenter_les_bases_de_données/"
title: "Exchange: Repairing and Defragmenting Databases"
description: "Guide on how to repair and defragment Exchange databases to improve performance and reduce storage space."
categories: ["Windows"]
date: "2007-01-20T11:00:00+02:00"
lastmod: "2007-01-20T11:00:00+02:00"
tags: ["Servers", "Windows", "Exchange"]
toc: true
---

## Introduction

During the defragmentation process, database objects that are no longer useful are removed from the database to increase its free space. By defragmenting an Exchange database, you'll increase data access speed, compact the database, and thus reduce the used space.

## Practice

To perform an Offline defragmentation, go to your Exchange server folder `C:\Program Files\EXCHSRVR\BIN`. There you will find the EDBUTIL utility that allows you to dismount the database.
Then use the "ESEUTIL" command with the /D option. For example, to defragment the PRIV.EDB file, use the following command:

```bash
ESEUTIL /D PRIV.EDB.
```

- The /P switch (which is generally used to repair a database) can be used in combination with the /D switch to increase performance and reliability. Every time you repair a database, the original database file does not change.

Unlike ESEUTIL which creates another file and sends the repaired database to this file. In the case of a fully functional database, using the /P switch with the /D switch results in the creation of the defragmented database in a separate file.

There are two advantages to proceeding this way:

- First, you know you're not overwriting good database files. So if something goes wrong in the defragmentation process, you don't have to worry about whether the database has been destroyed.
- The second advantage of using this method is that it increases the speed of the defragmentation process since Exchange doesn't have to browse the data in a single file.

The /T option is not required but allows you to control the name and location of the new version of the database. Once defragmentation is complete, you can simply move this new database to the location where the old database is, delete or rename the old database, then rename the new database by giving it the same name as the old database.
The syntax is:

```bash
ESEUTIL /D /P "path and file" /T "path and file"
```

For example, if you are defragmenting the PRIV.EBD file and want to create a new file called PRIV2.EDB, you would use the following command:

```bash
ESEUTIL /D /P"C:\PROGRAM FILES\EXCHSRVR\MDBDATA\PRIV.EDB" /T"C:\PROGRAM FILES\EXCHSRVR\MDBDATA\PRIV2.MDB"
```
