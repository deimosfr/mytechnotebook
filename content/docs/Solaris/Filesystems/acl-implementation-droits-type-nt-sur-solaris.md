---
weight: 999
url: "/ACL_Implementation_droits_type_NT_sur_Solaris/"
title: "ACL: Implementing NT-Style Permissions on Solaris"
description: "How to implement and use NT-style ACLs (Access Control Lists) on Solaris systems for more flexible file permissions management."
categories: ["Solaris", "Security", "Administration"]
date: "2010-02-10T13:13:00+02:00"
lastmod: "2010-02-10T13:13:00+02:00"
tags:
  ["acl", "permissions", "solaris", "zfs", "security", "nt", "access control"]
toc: true
---

## Introduction

With respect to a computer filesystem, an access control list ([ACL](https://en.wikipedia.org/wiki/Access_control_list)) is a list of permissions attached to an object. An ACL specifies which users or system processes are granted access to objects, as well as what operations are allowed to be performed on given objects. In a typical ACL, each entry in the list specifies a subject and an operation (e.g. the entry (Alice, delete) on the ACL for file WXY gives Alice permission to delete file WXY).

This documentation is a quick reference. If you need more detailed explanations, please refer to the SUN ACL documentation on their website.

## Enabling ACL

By default, on ZFS, ACLs are automatically enabled. However, there are different modes to choose from based on your usage requirements.

You can see the current default configuration with the "zfs get all" command:

```bash
$ zfs get all zfs_volume
NAME                     PROPERTY         VALUE                         SOURCE
zfs_volume  type             filesystem                    -
...
zfs_volume  aclmode          groupmask                     default
zfs_volume  aclinherit       restricted                    default
...
```

### Inheritance mode

aclinherit - This property determines the behavior of ACL inheritance. Values include the following:

- discard - For new objects, no ACL entries are inherited when a file or directory is created. The ACL on the file or directory is equal to the permission mode of the file or directory.
- noallow - For new objects, only inheritable ACL entries that have an access type of deny are inherited.
- restricted - For new objects, the write_owner and write_acl permissions are removed when an ACL entry is inherited.
- passthrough - When property value is set to passthrough, files are created with a mode determined by the inheritable ACEs. If no inheritable ACEs exist that affect the mode, then the mode is set in accordance to the requested mode from the application.
- passthrough-x - Has the same semantics as passthrough, except that when passthrough-x is enabled, files are created with the execute (x) permission, but only if execute permission is set in the file creation mode and in an inheritable ACE that affects the mode.

The default mode for the aclinherit is restricted.

### Rights on creation mode

aclmode - This property modifies ACL behavior when a file is initially created or whenever a file or directory's mode is modified by the chmod command. Values include the following:

- discard - All ACL entries are removed except for the entries needed to define the mode of the file or directory.
- groupmask - User or group ACL permissions are reduced so that they are no greater than the group permission bits, unless it is a user entry that has the same UID as the owner of the file or directory. Then, the ACL permissions are reduced so that they are no greater than owner permission bits.
- passthrough - During a chmod operation, ACEs other than owner@, group@, or everyone@ are not modified in any way. ACEs with owner@, group@, or everyone@ are disabled to set the file mode as requested by the chmod operation.

The default mode for the aclmode property is groupmask.

### Changing mode

You can change mode with commands like these:

```bash
zfs set aclmode=passthrough zfs_volume
zfs set aclinherit=passthrough zfs_volume
```

Simply choose the one you prefer for your needs.

## ACL Properties

### Example

You can use the ls command with special arguments to see current ACL rights. Choose the format that's easier for you to read.

- ls -dv:

```bash
$ ls -dv zfs_volume
drwxrwxr-x  11 myuser   mygroup        11 oct  14 12:06 zfs_volume
     0:owner@::deny
     1:owner@:list_directory/read_data/add_file/write_data/add_subdirectory
         /append_data/write_xattr/execute/write_attributes/write_acl
         /write_owner:allow
     2:group@::deny
     3:group@:list_directory/read_data/add_file/write_data/add_subdirectory
         /append_data/execute:allow
     4:everyone@:add_file/write_data/add_subdirectory/append_data/write_xattr
         /write_attributes/write_acl/write_owner:deny
     5:everyone@:list_directory/read_data/read_xattr/execute/read_attributes
         /read_acl/synchronize:allow
```

- ls -dV:

```bash
$ ls -dV zfs_volume
drwxrwxr-x  11 myuser   mygroup        11 oct  14 12:06 zfs_volume
            owner@:--------------:------:deny
            owner@:rwxp---A-W-Co-:------:allow
            group@:--------------:------:deny
            group@:rwxp----------:------:allow
         everyone@:-w-p---A-W-Co-:------:deny
         everyone@:r-x---a-R-c--s:------:allow
```

### Complete properties list

#### ACL Entry Types

| ACL Entry Type | Global | Description                                                                                                                                                                    |
| -------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| owner@         | yes    | Specifies the access granted to the owner of the object.                                                                                                                       |
| group@         | yes    | Specifies the access granted to the owning group of the object.                                                                                                                |
| everyone@      | yes    | Specifies the access granted to any user or group that does not match any other ACL entry. With a user name, specifies the access granted to an additional user of the object. |
| user           | no     | Must include the ACL-entry-ID, which contains a username or userID. If the value is not a valid numeric UID or username, the ACL entry type is invalid.                        |
| group          | no     | Must include the ACL-entry-ID, which contains a groupname or groupID. If the value is not a valid numeric GID or groupname, the ACL entry type is invalid.                     |

#### ACL Access Privileges

| Access Privilege | Compact Access Privilege | Description                                                                                                           |
| ---------------- | ------------------------ | --------------------------------------------------------------------------------------------------------------------- |
| add_file         | w                        | Permission to add a new file to a directory.                                                                          |
| add_subdirectory | p                        | On a directory, permission to create a subdirectory.                                                                  |
| append_data      | p                        | Placeholder. Not currently implemented.                                                                               |
| delete           | d                        | Permission to delete a file.                                                                                          |
| delete_child     | D                        | Permission to delete a file or directory within a directory.                                                          |
| execute          | x                        | Permission to execute a file or search the contents of a directory.                                                   |
| list_directory   | r                        | Permission to list the contents of a directory.                                                                       |
| read_acl         | c                        | Permission to read the ACL (ls).                                                                                      |
| read_attributes  | a                        | Permission to read basic attributes (non-ACLs) of a file.                                                             |
| read_data        | r                        | Permission to read the contents of the file.                                                                          |
| read_xattr       | R                        | Permission to read the extended attributes of a file or perform a lookup in the file's extended attributes directory. |
| synchronize      | s                        | Placeholder. Not currently implemented.                                                                               |
| write_xattr      | W                        | Permission to create extended attributes or write to the extended attributes directory.                               |
| write_data       | w                        | Permission to modify or replace the contents of a file.                                                               |
| write_attributes | A                        | Permission to change the times associated with a file or directory to an arbitrary value.                             |
| write_acl        | C                        | Permission to write the ACL or the ability to modify the ACL by using the chmod command.                              |
| write_owner      | o                        | Permission to change the file's owner or group.                                                                       |

#### ACL Inheritance Flags

| Inheritance Flag | Compact Inheritance Flag | Description                                                                                                                               |
| ---------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- |
| file_inherit     | f                        | Only inherit the ACL from the parent directory to the directory's files.                                                                  |
| dir_inherit      | d                        | Only inherit the ACL from the parent directory to the directory's subdirectories.                                                         |
| inherit_only     | i                        | Inherit the ACL from the parent directory but applies only to newly created files or subdirectories and not the directory itself.         |
| no_propagate     | n                        | Only inherit the ACL from the parent directory to the first-level contents of the directory, not the second-level or subsequent contents. |
| -                | N/A                      | No permission granted.                                                                                                                    |

## Rights Management

### Adding rights

To add rights to a folder or file using ACLs:

```bash
$ chmod A+user:myuser:read_data/execute:allow directory
```

- A+: A means use ACL and + means add
- user:myuser: add username (here myuser)
- read_data/execute:allow: allowing these rights
- directory: the directory to change

You can verify the user has been added with their rights:

```bash
$ ls -dv test.dir
drwxr-xr-x+ 2 root      root           2 Aug 31 12:02 directory
    0:user:myuser:list_directory/read_data/execute:allow
    1:owner@::deny
    2:owner@:list_directory/read_data/add_file/write_data/add_subdirectory
        /append_data/write_xattr/execute/write_attributes/write_acl
        /write_owner:allow
    3:group@:add_file/write_data/add_subdirectory/append_data:deny
    4:group@:list_directory/read_data/execute:allow
    5:everyone@:add_file/write_data/add_subdirectory/append_data/write_xattr
        /write_attributes/write_acl/write_owner:deny
    6:everyone@:list_directory/read_data/read_xattr/execute/read_attributes
        /read_acl/synchronize:allow
```

For a faster alternative, you can use:

```bash
$ chmod A+user:myuser:rx:allow directory
```

### Deleting rights

To remove the previously added user (ID 0):

```bash
$ chmod A0- directory
```

- A0-: A for ACL, 0 for ID 0, and - for deleting

Verify the user has been removed:

```bash
$ ls -dv test.dir
drwxr-xr-x+ 2 root      root           2 Aug 31 12:02 directory
    0:owner@::deny
    1:owner@:list_directory/read_data/add_file/write_data/add_subdirectory
        /append_data/write_xattr/execute/write_attributes/write_acl
        /write_owner:allow
    2:group@:add_file/write_data/add_subdirectory/append_data:deny
    3:group@:list_directory/read_data/execute:allow
    4:everyone@:add_file/write_data/add_subdirectory/append_data/write_xattr
        /write_attributes/write_acl/write_owner:deny
    5:everyone@:list_directory/read_data/read_xattr/execute/read_attributes
        /read_acl/synchronize:allow
```

You can delete another right by changing the number (e.g., A4-).

To completely remove all ACLs:

```bash
$ chmod A- directory
```

### Replacement

To replace an existing right with another:

```bash
$ chmod A0=user:myuser:execute:deny directory
```

This changes the specified ACL entry (ID 0) to deny execute permission for myuser.

For a faster alternative:

```bash
$ chmod A0=user:myuser:x:deny directory
```

{{< alert context="warning" text="DO NOT FORGET TO SPECIFY ID OR IT WILL REPLACE ALL YOUR CURRENT RIGHTS WITH THIS SINGLE ONE" />}}

To replace all rights with only one user permission:

```bash
chmod A=user:myuser:read_data:allow directory
```

This removes all other rights, including owner permissions:

```bash
$ ls -v directory
----------+ 1 root      root        2455 Dec 25 12:08 directory
    0:user:myuser:read_data:allow
```

You can also reset rights using standard chmod:

```bash
chmod 755 directory
```

This restores the standard permission set with ACLs.

### Inheritance

Remember that file and directory inheritance depends on the ACL mode you've chosen. To add inheritance:

```bash
$ chmod A+user:myuser:read_data/execute:file_inherit:allow directory
```

This works only for files. Use dir_inherit for directories.

## References

[https://docs.sun.com/app/docs/doc/819-5461?l=en](https://docs.sun.com/app/docs/doc/819-5461?l=en)
