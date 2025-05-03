---
weight: 999
url: "/PAM_(Pluggable_Authentification_Module)_\\:_Choisir_ses_méthodes_d'authentifications/"
title: "PAM (Pluggable Authentication Module): Choosing Authentication Methods"
description: "A comprehensive guide to PAM (Pluggable Authentication Module) configuration on Linux systems, explaining how to customize authentication methods for various services."
categories: ["Linux", "Security"]
date: "2008-04-06T07:42:00+02:00"
lastmod: "2008-04-06T07:42:00+02:00"
tags: ["PAM", "Authentication", "Linux", "Security", "System Administration"]
toc: true
---

## Introduction

PAM is THE standard authentication system used in Linux. The power of this tool is unlimited, but it's not always well documented which tends to work against it. As an introduction, let's first look at the motivations that led to the creation of PAM. Originally, in Unix (and in early versions of Linux), the file that centralized user management was `/etc/passwd`. It contained many sensitive pieces of information, including encrypted passwords. To use a Unix machine, the first thing to do was to authenticate via the login program (the last program launched by init). This program was developed to parse the `/etc/passwd` file.

However, over time, it was realized that storing user passwords in a file readable by everyone could represent a security hole (as personal machines became increasingly capable of brute-forcing passwords). This realization led to the creation of a separate password file called `/etc/shadow`.

However, login had to be reimplemented to take this change into account. The same was true for all programs requiring authentication (ftp, su, sudo...).

Later, it became apparent that having a flat file for authentication could be limiting when handling tens of thousands of accounts. Other databases were therefore used, such as directories. The idea was revolutionary, but the problem was that, once again, certain parts of the authentication code for the concerned applications had to be rewritten.

After these observations (a few hours of software engineering and hundreds of thousands of lines of code later), Linux kernel developers had the idea to move the entire authentication layer outside the programs that needed it. They therefore created the Pluggable Authentication Modules.

## Mode of Operation

Linux PAM is a set of dynamic libraries managing specific points of the authentication process. As we will see below, when we talk about authentication, we are not limited to the simple login/password challenge, but rather to a set of authorization points that can affect both authentication itself and sessions or password management.

An application can be developed to dynamically link to these libraries and thus access already implemented and stable functionalities. As such, Linux PAM provides developers with a detailed API of available functions.

For information, these libraries are generally found in the `/lib/security` directory.

From the system administrator's perspective, this allows configuring the behavior to adopt for all applications requiring authentication. The only restriction is that the application must be "PAM enabled", meaning it was developed to use PAM libraries. Refer to the application documentation in case of doubt.

The administrator can therefore, if desired, define that for a given service (for example xdm), authentication will be done in 3 distinct steps (materialized by the use of 3 dynamic libraries, therefore 3 modules) to which configuration arguments will be passed.

One can therefore choose exactly their own authentication policy for this application, independently of the application itself.

In PAM terminology, configuring an application amounts to configuring access to the service. Thus, if you install proftpd or wuftpd, you will provide the ftp service to your users.

All services are configured in the `/etc/pam.d/` directory where each file details the authentication policies related to that service.

## Hierarchical Organization of Authentication Tasks

To simplify the use and understanding of the role of each module, authentication tasks are divided into 4 independent groups:

* Account
* Authentication
* Password
* Session

A service is therefore divided into 4 groups. One or more modules are assigned to each group.

Although in theory these groups seem to clearly delineate each task of the authentication process, the reader should keep in mind that for some modules, distinguishing the placement of the module among these 4 groups is not always easy, especially since a module can intervene in several groups. To know precisely the possibilities offered by a module for each group, one must read the documentation of the developer of the module in question. Generally, the synopsis of a module's usage is provided in its man page.

To summarize:

* A module is a piece of code capable of being dynamically linked to an application providing a service.
* A service is configured in a file (named according to its service type) contained in the `/etc/pam.d/` directory.
* A module offers a certain number of functionalities that are organized into 4 groups, defining their scope of action.
* A module can provide functionalities for one or more groups. For example, the pam_xauth module only provides functionalities in the session group whereas pam_unix provides them for each group.

## Example of Linux PAM Usage

The first idea that generally comes to an administrator's mind regarding security is to define who can physically access a machine. For this, we will configure 2 accesses:

* login which is used to connect in console mode
* gdm which is used to connect in graphical mode with the Gnome environment (the reader will use their own session manager, according to the work environment)

These applications being "PAM enabled", we can freely define the security policy we desire.

```bash
[seb@localhost seb]$ cat /etc/pam.d/login
#%PAM-1.0 : login service
auth required pam_nologin.so
auth required pam_access.so
#auth required pam_securetty.so
auth required pam_stack.so
service=system-auth
account required pam_stack.so service=system-auth password required pam_stack.so
service=system-auth session required pam_stack.so
service=system-auth session optional pam_console.so
```

As we mentioned earlier, each service (defined by the file name) uses modules (3rd column), distributed into groups (1st column). Similarly, some modules can be used multiple times, in different groups, their behavior being different.

In the examples provided above, we can see the use of the pam_stack module in almost all authentication groups. This module is a bit special since its only objective is to refer authentication to the module passed as an argument. It thus allows defining a common behavior for many services, ideal for simple or homogeneous configurations.

If you modify the configuration of the system-auth module, all other modules that use it through pam_stack will be modified.

Another important point in the configuration of service files concerns the stacking of modules in the order of reading. Thus, in the example above, in the authentication group, the pam_nologin module will be evaluated before the pam_access module, itself evaluated before pam_stack. This is important when one considers that one can define what the consequence of the successful use of a module will be.

We have the possibility to define 4 types of success obligations to hold:

* required
* requisite
* sufficient
* optional

If we synthesize all the concepts we have just discussed, the configuration of a security policy that applies to a service involves the use of one or more modules, informed in a file located in `/etc/pam.d`. In this file, each line informs PAM to use a functionality of a module according to the group, written at the beginning of the line (auth, session...).

PAM will sequentially read the lines of each group and will confront the success obligation defined by the administrator (required, requisite...) with the return code sent back by the library and decide whether to continue or not, according to the result.

It is important to identify what gives rise to an interaction with the user, for example a module of the type pam_unix.so, and what only applies a treatment (pam_env) to set the control flag.

Let's take a simple example. Imagine that we want to define the following scenario for authentication with gdm:

* The user must have a local account on the machine
* Their password must not be null
* We wish to define environment variables for them
* The user must have the right to log in

We are therefore only interested in the authentication part, each line that we will fill in will therefore be preceded by "auth":

```bash
auth required pam_env.so // Define the use of user environment variables
auth required pam_unix.so // Use the standard Linux authentication procedure
auth required pam_nologin.so // The user must be among the users with login rights
auth required pam_deny.so
```

When Boulay wants to log in to the machine, PAM will check the following points: The first entry in the stack concerns the pam_env module. This defines environment variables specified in the `/etc/environment` file. A priori, no error can be generated at this stage, stack reading continues.

The pam_unix module will handle Unix-style authentication, that is, it will read `/etc/passwd` and `/etc/shadow` to validate the login/password challenge. In addition, we provide the argument that a null password is not acceptable. If Boulay succeeds in proving who he really is by this mechanism, PAM will continue reading the stack.

The last module used for authentication is pam_deny which will systematically refuse the connection request. As with pam_env, the result will always be valid. The total result of the stack is valid, so the user can access the service. Login will fork and execute a shell for them, giving them access to the machine.

## Example of Advanced Configuration

After seeing how to organize the configuration of PAM's behavior for a simple service, we will now detail some configuration examples to significantly increase the security of certain programs.

By default, distributions accept that any user can request access to a super user session by invoking the su command. To only accept users who are members of the wheel group, add the following line in the su service (`/etc/pam.d/su`):

```bash
auth required pam_wheel.so use_uid
```

And since PAM is modular, it will be the same for each service. You can therefore define that for such a service, only users who are members of the wheel group will be able to authenticate.

With the above example, it's easy to understand that it is now very simple to apply execution rights according to a given service and user.

One of the recurring problems under Unix is the use of certain minimalist C programs that can cause the machine to malfunction. Thus, as is often the case in scientific schools or universities, students may use the famous fork Bomb (while(1) fork() ), leading the machine into a deluge of process creation, no longer allowing it to respond to legitimate actions.

Fortunately, the pam_limits.so module allows limiting the resources used by users. This module, once defined in the service one wishes to use, must be configured by editing the `/etc/security/limits.conf` file.

This module works in the session group, so we edit `/etc/pam.d/{gdm,kdm,login}` according to the manager used:

```bash
session required /lib/security/pam_limits.so // Without arguments, uses the file /etc/security/limits.conf
```

Then, by editing the limits.conf file:

```bash
#<domain> <type> <item> <value> boulay hard fsize 100000 // Defines a maximum disk usage size for the user boulay
@etudiant hard nproc 30 // Defines a maximum number of ongoing processes for the student user group
```

Thus, during each session the user boulay will not be able to create a file with a size greater than 10MB and students will not be able to launch more than 30 processes. However, be careful with the limits you wish to define, knowing that with environments such as KDE or Gnome, many background processes are executed...

## Conclusion

Due to its power, PAM may seem quite complex to use at first, especially since an error can quickly prevent you from accessing your machine (only a reboot in single user mode can help you in this case). Nevertheless, the large number of official and unofficial modules allow making authentication tasks under Linux very versatile. You can easily integrate authentication through LDAP, NDS, or with an NT server or even by biometrics without modifying the code of the services.
