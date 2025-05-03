---
weight: 999
url: "/Créer_un_package_RedHat_depuis_un_tar/"
title: "Creating a RedHat Package from a Tar File"
description: "How to easily create an RPM package from a tar file in RedHat environments in just 5 minutes."
categories: ["Linux", "Red Hat"]
date: "2012-08-02T11:51:00+02:00"
lastmod: "2012-08-02T11:51:00+02:00"
tags: ["Servers", "Software Packaging", "RPM"]
toc: true
---

{{< table "table-striped table-hover" >}}
| | |
|---|---|
| Operating System | Red Hat 6.3 |
| Website | [Website](https://www.redhat.com) |
| Last Update | 02/08/2012 |
{{< /table >}}

## Introduction

There are many methods for creating packages available on the web, but many of them are obsolete and often don't allow for simple implementation. That's why I'll show you a technique to build a package in 5 minutes.

## Installation

We'll install the necessary tools for package creation:

```bash
yum install rpmbuild
```

## Creating the Tar Archive

Let's say I want to package a Ruby library that I installed via gems. Everything is installed in specific directories, so I'll create an archive containing what I need:

```bash
tar -czvPf rubygemsysproctable-0.9.tar.gz /usr/lib/ruby/gems/1.8/doc/sys-proctable-0.9.0-x86-linux /usr/lib/ruby/gems/1.8/gems/sys-proctable-0.9.0-x86-linux
```

## Creating the RPM

For this script, I'm using [a script that I found on this site](https://www.mindtwist.de/main/downloads.html?task=summary&cid=3&catid=3), which automatically creates a .spec file with the main information. Adapt the beginning with your desired information:

```bash
#!/bin/bash
# This script creates an RPM from a tar file.
# $1 : tar file

NAME=$(echo ${1%%-*} | sed 's/^.*\///')
VERSION=$(echo ${1##*-} | sed 's/[^0-9]*$//')
RELEASE=0
VENDOR="Deimos"
EMAIL="<xxx@mycompany.com>"
SUMMARY="Summary"
LICENSE="GPL"
GROUP="System"
ARCH="noarch"
DESCRIPTION="My description"

######################################################
# users should not change the script below this line.#
######################################################

# This function prints the usage help and exits the program.
usage(){
    /bin/cat << USAGE

This script has been released under BSD license. Copyright (C) 2010 Reiner Rottmann <rei..rATrottmann.it>

$0 creates a simple RPM spec file from the contents of a tarball. The output may be used as starting point to create more complex RPM spec files.
The contents of the tarball should reflect the final directory structure where you want your files to be deployed. As the name and version get parsed
from the tarball filename, it has to follow the naming convention "<name>-<ver.si.on>.tar.gz". The name may only contain characters from the range
[A-Z] and [a-z]. The version string may only include numbers seperated by dots.

Usage: $0  [TARBALL]

Example:
  $ $0 sample-1.0.0.tar.gz

  $ /usr/bin/rpmbuild -ba /tmp/sample-1.0.0.spec

USAGE
    exit 1
}

if echo "${1##*/}" | sed 's/[^0-9]*$//' | /bin/grep -q  '^[a-zA-Z]\+-[0-9.]\+$'; then
   if /usr/bin/file -ib "$1" | /bin/grep -q "application/x-gzip"; then
      echo "INFO: Valid input file '$1' detected."
   else
      usage
   fi
else
    usage
fi

OUTPUT=/tmp/${NAME}-${VERSION}.spec

FILES=$(/bin/tar -tzf $1 | /bin/grep -v '^.*/$' | sed 's/^/\//')

/bin/cat > $OUTPUT << EOF
Name: $NAME
Version: $VERSION
Release: $RELEASE
Vendor: $VENDOR
Summary: $SUMMARY
License: $LICENSE
Group: $GROUP
Source0: %{name}-%{version}.tar.gz
BuildRoot: /var/tmp/%{name}-buildroot
BuildArch: $ARCH

%description
$DESCRIPTION

%prep

%setup -c -n %{name}-%{version}

%build

%install
[ -d \${RPM_BUILD_ROOT} ] && rm -rf \${RPM_BUILD_ROOT}
/bin/mkdir -p \${RPM_BUILD_ROOT}
/bin/cp -axv \${RPM_BUILD_DIR}/%{name}-%{version}/* \${RPM_BUILD_ROOT}/


%post

%postun

%clean

%files
%defattr(-,root,root)
$FILES

%define date    %(echo \`LC_ALL="C" date +"% a % b % d % Y"\`)

%changelog

* %{date} User $EMAIL
- first Version

EOF

echo "INFO: Spec file has been saved as '$OUTPUT':"
echo "---------%<----------------------------------------------------------------------"
/bin/cat $OUTPUT
echo "---------%<----------------------------------------------------------------------"
```

Now simply run this script on your tar.gz file:

```bash
chmod 755 tgz2rpm.sh
./tgz2rpm.sh rubygemsysproctable-0.9.tar.gz
```

Then build the RPM:

```bash
rpmbuild -ba /tmp/rubygemsysproctable-0.9.spec
```

The RPM is now created in ~/rpmbuild/RPMS/noarch/ :-)

## Verification

Let's verify our RPM:

```bash
> rpm -qip ~/rpmbuild/RPMS/noarch/rubygemsysproctable-0.9-0.noarch.rpm
Name        : rubygem-sysproctable         Relocations: (not relocatable)
Version     : 0.9                               Vendor: Deimos
Release     : 0                             Build Date: Thu 02 Aug 2012 01:19:08 PM CEST
Install Date: (not installed)               Build Host: server1.deimos.fr
Group       : System Environment/Libraries   Source RPM: rubygem-sysproctable-0.9-0.src.rpm
Size        : 159304                           License: GPL
Signature   : (none)
Packager    : Ruby sys-proctable
Summary     : Ruby sys-proctable library
Description :
Ruby sys-proctable library
```

## References

http://fedoraproject.org/wiki/How_to_create_an_RPM_package  
http://www.mindtwist.de/main/linux/3-linux-tipps/32-how-to-convert-tar-gz-archive-to-rpm-.html
