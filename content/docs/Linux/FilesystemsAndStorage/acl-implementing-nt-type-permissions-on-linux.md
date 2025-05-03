---
weight: 999
url: "/ACL\\:_Implémentation_des_droits_de_type_NT_sur_Linux/"
title: "ACL: Implementing NT-type Permissions on Linux"
description: "Learn how to implement NT-type permissions (ACLs) on Linux systems to extend file access control beyond traditional Unix permissions."
categories: ["Linux", "Security", "System Administration"]
date: "2013-05-06T15:00:00+02:00"
lastmod: "2013-05-06T15:00:00+02:00"
tags: ["ACL", "Permissions", "Linux", "Samba", "Security"]
toc: true
---

An ACL, or Access Control List, is simply defined as a list of permissions on a file, directory, or tree structure, added to the "classic" permissions (technically, POSIX.1 permissions) of that file. These permissions concern defined users and/or groups. ACL management under GNU/Linux is inspired by the POSIX 1003.1e standard (project 17) but does not fully comply with it.

With ACLs, you can extend the number of users and groups that have rights to the same file. Remember that in the UNIX world, each file can normally only indicate permissions for a single user and a single group, which are opposed to a single category corresponding to "all others" (or "the rest of the world"). With ACLs, you can (among other things) add other users and groups to a file and define their rights separately. This brings the system closer to the permission system used on NT platforms (although many differences remain).

ACLs are very useful (and even essential) in collaborative and shared computing environments; similarly, their use with SAMBA extends its capabilities.

However, be careful not to confuse them! Unix ACLs are not identical to those of NT (Microsoft). Indeed, there are variants; for example, only the owner of a file can change the owner of that file, even if other users have all rights to the file.

If, following this, you want to set up a Samba server and use it as under Windows with ACLs, there is a small difference in operation:

- Under Samba, ACLs can be modified by the owner of the object, members of the group that owns the object, or the administrator of the share
- Under Linux, ACLs can be modified by the owner of the object or root (to simplify)

Here are some useful resources:

[Documentation on NT-type ACLs](/pdf/acl_nt.pdf)  
[ACL and EA under Linux](/pdf/acl_et_ea_sous_linux.pdf)

## Useful Commands

Copy the rights from one file to another:

```bash
getfacl <file-with-acl> | setfacl -f - <file-with-no-acl>
```
