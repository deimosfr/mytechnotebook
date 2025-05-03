---
weight: 999
url: "/Rebooter_sa_Freebox_Server_6_en_ligne_de_commande/"
title: "Reboot Your Freebox Server 6 via Command Line"
description: "A guide on how to reboot your Freebox Server 6 from the command line, including installation instructions for pfSense and other systems."
categories: 
  - "Linux"
  - "FreeBSD"
  - "Network"
date: "2012-03-24T23:25:00+02:00"
lastmod: "2012-03-24T23:25:00+02:00"
tags:
  - "Freebox"
  - "Command line"
  - "Network"
  - "pfSense"
  - "bash"
toc: true
---

## Introduction

Since version 6 of the Freebox, there is a web interface ([https://mafreebox.freebox.fr](https://mafreebox.freebox.fr)) that allows you to control several options, including rebooting it. We will use this capability to restart it via the command line. For this purpose, I'm using [a script that was obtained from admin-linux.fr](https://www.admin-linux.fr/?p=5049).

I tested it on Debian and it works natively. On pfSense 2.0.1, I encountered a few small issues. That's why I will detail here how I resolved them.

## Installation

### Prerequisites

You'll need bash and wget to run the script that follows. On most distributions, there are no issues as they are already installed. But on minimal or embedded versions, this can cause problems. I will explain how to proceed on pfSense.

#### pkg_add

Install wget and bash:

```bash
pkg_add -vr wget
pkg_add -vr bash
```

#### Manually

If you encounter problems because the repository no longer exists, take the one closest to your version and upload the files to your machine (pfSense for example):

```bash
wget ftp://ftp.freebsd.org/pub/FreeBSD/ports/i386/packages-8.2-release/ftp/wget-1.12_2.tbz
wget ftp://ftp.beastie.tdk.net/pub/FreeBSD/ports/i386/packages-8.2-release/shells/bash-4.1.9.tbz
```

Then we'll decompress everything:

```bash
cd /
tar -xzf wget-1.12_2.tbz
tar -xzf bash-4.1.9.tbz
```

### Script

Let's put this script in `/bin/freebox` for example:

```bash
#!/bin/bash

# Description : Script for rebooting the FreeBox
#	from the command line.
# Author : FHH <fhh@admin-linux.fr>
# $Id: freebox,v 1.6 2011/08/28 22:59:56 fhh Exp $

# Global script variables:
USERNAME="freebox" ;
URL="http://mafreebox.freebox.fr" ;

WORKING_DIRECTORY="${HOME}/.freebox" ;
CONFIG_FILE="${WORKING_DIRECTORY}/config" ;
TMP_DIRECTORY="${WORKING_DIRECTORY}/tmp" ;
COOKIES="cookies" ;

FREE_LOGIN_SCRIPT="/login.php" ;
FREE_REBOOT_CGI="/system.cgi" ;
WGET_OPTS="-T1 -t1 -qO /dev/null" ;

# usage : Display help
usage () {
	echo -e "$0 <options> [command]\n" ;
	echo "Options" ;
	echo " -h		: Display this help message." ;
	echo " -p <pass>	: Specify the FreeBox password." ;
	echo ;
	echo "Commands" ;
	echo " restart	: Restart the freebox." ;
	echo " savepass	: Save the password." ;
	echo ;
	echo "In case of malfunction, check your " ;
	echo "FreeBox password and/or the file \"~/.freebox/config\"." ;
}

# Script exit function:
die() {
        echo $@ >&2 ;
        exit 1 ;
}
# Authentication function on the FreeBox
auth() {
#       Authentication and cookie registration:
	wget ${WGET_OPTS} --save-cookies "${TMP_DIRECTORY}/${COOKIES}" --post-data "login=${USERNAME}&passwd=${password}" "${URL}/${FREE_LOGIN_SCRIPT}" ;
#	If the cookie is not downloaded, authentication failed:
	[ -e "${TMP_DIRECTORY}/${COOKIES}" ] || die "> Authentication failure. Check the password." ;
	[ ! -z "$(sed -e 's/[[:blank:]]//g' -e '/^#\|^$/d' "${TMP_DIRECTORY}/${COOKIES}")" ] || die "> Authentication failed. Check the password." ;
}

# Password input function
askpwd() {
	cpt=0 ;
	while [ -z "${password}" ] ; do
		(( ${cpt} >= 3 )) && die "> You must provide your password. Stop!" ;
		(( cpt = cpt + 1 )) ; 
		read -sp "Password: " password ;
		echo ;
	done
}

init () {
	[ ! -d "${WORKING_DIRECTORY}" ] && { mkdir "${WORKING_DIRECTORY}" || die "> Unable to create \"${WORKING_DIRECTORY}\"" ; }
	[ ! -d "${TMP_DIRECTORY}" ] && { mkdir "${TMP_DIRECTORY}" || die "> Unable to create \"${TMP_DIRECTORY}\"" ; }
	[ -e "${TMP_DIRECTORY}/${COOKIES}" ] && rm "${TMP_DIRECTORY}/${COOKIES}" ;
}

connection () {
	wget ${WGET_OPTS} ${URL} ;
	case $? in
		4)
			die "> Problem accessing the Freebox. Try to connect manually" ;;
		5)
			die "> SSL problem. Try to connect manually" ;;
	esac
}

cleancook () {
#       Cleanup of the cookie:
        [ -e "${TMP_DIRECTORY}/${COOKIES}" ] && rm "${TMP_DIRECTORY}/${COOKIES}" ;
}

# Function to restart the FreeBox
restart () {
	init ;
	[ -z "${password}" -a -r "${CONFIG_FILE}" ] && . ${CONFIG_FILE} ;
	connection ;
	askpwd ;
	auth ;
#	Restart the freebox:
	wget ${WGET_OPTS} --load-cookies "${TMP_DIRECTORY}/${COOKIES}" --post-data 'method=system.reboot&redirect_after=/reboot.php&timeout=1' -p "${URL}/${FREE_REBOOT_CGI}" ;
	cleancook ;
}

savepass () {
	init ;
	askpwd ;
	auth ;
	echo "password=${password}" > "${CONFIG_FILE}" ;
	chmod 0600 "${CONFIG_FILE}"
	cleancook ;
}

# Main program:
(( $# > 0 )) || { usage ; exit 0; }

while getopts ":p:h" opt ; do
	case ${opt} in
		p)
			password=${OPTARG} ;;
		h)
			usage ; 
			exit 0 ;;
		:)
			echo -e "> Option -$OPTARG requires an argument.\n" ;
			usage ;
			exit 1 ;;
		*)
			echo -e "> Invalid option \"-$OPTARG\".\n" ; 
			usage ;
			exit 1 ;;
	esac
done
shift $((OPTIND-1))

(( $# > 1 )) && { echo -e "> Only one command at a time.\n" ; usage ; exit 1 ; } ;
(( $# < 1 )) && { echo -e "> No command specified.\n" ; usage ; exit 1 ; } ;

case $1 in
	restart)
		restart ;;
	savepass)
		savepass ;;
	*)
		echo -e "> $1 : Unknown command!\n" ;
		usage ;
		exit 1 ;;
esac
```

Then give it the correct permissions:

```bash
chmod 555 /bin/freebox
```

## Usage

We'll save the password for the Freebox interface like this:

```bash
> freebox savepass
Password:
```

Here is the method to restart the Freebox:

```bash
freebox reboot
```

## Resources
- [https://www.admin-linux.fr/?p=5049](https://www.admin-linux.fr/?p=5049)
