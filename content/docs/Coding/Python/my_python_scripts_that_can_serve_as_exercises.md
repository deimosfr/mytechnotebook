---
weight: 999
url: "/Mes_scripts_Python_qui_peuvent_servir_d'exercices/"
title: "My Python Scripts That Can Serve as Exercises"
description: "A collection of Python scripts that can serve as learning exercises for beginners, including a script for automating package signing and pushing in Red Hat Satellite."
categories: ["Development", "Linux", "Red Hat"]
date: "2012-06-06T13:09:00+02:00"
lastmod: "2012-06-06T13:09:00+02:00"
tags: ["Python", "Red Hat", "Satellite", "Development", "Scripts", "Automation"]
toc: true
---

![Python](/images/python-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.7 / 3.2 |
| **Website** | [Python Website](https://www.python.org/) |
| **Last Update** | 06/06/2012 |
{{< /table >}}

## Introduction

It's not always easy to start learning a programming language, especially when you've never studied it in school! That's why I'm offering some small scripts that I created right after going through some books. The scripts should increase in difficulty more or less progressively.

## Red Hat Satellite

I created this script for the [Red Hat Satellite server](/Satellite_:_Deploiement_d'OS_Red_Hat_via_Red_Hat_Satellite/) to save time when signing and sending multiple packages to a Satellite server:

(`satellite_add_packages.py`)

```python
#!/usr/bin/env python
# Made by Pierre Mavro 14/03/2012
# Version : 0.1
# This script permit to automate in a non secure way, new packages for a custom repository on Red Hat Satellite
# Require : pexpect

import getopt, os, sys, glob, pexpect
from string import Template

# Help
cmd_name = sys.argv[0]
def help(code):
    print cmd_name, "[-h] [-r] [-s] [-l] [-p] [-d]"
    str = """
   -h, --help
      Show this help
   -s, --passphrase
      Passphrase to sign packages
   -r, --repository
      Select wished repository to push the packages
   -l, --login
      Red Hat Network username
   -p, --password
      Red Hat Network password
   -f, --folder
      folder were new packages should be added (default: /tmp/packages)
   -d, --debug
      Debug mode
"""
    print str
    sys.exit(code)

class bcolors:
    OK = '\033[92m'
    FAIL = '\033[91m'
    END = '\033[0m'

    def disable(self):
        self.OK = ''
        self.FAIL = ''
        self.END = ''

# Sign and push function
def sign_push(passphrase,repository,login,password,folder,debug):

    # Package signing
    def sign(rpm_files,passphrase,folder,debug,charspacing):
        if (debug == 1): print 80*'=' + "\n"
        print '[+] Signing packages :'
        # Sign all packages
        for package in rpm_files:
            # Formating
            charspace = Template("{0:<$space}")
            print charspace.substitute(space = charspacing).format(' - ' + package + '...'),
            # Launch resign
            child = pexpect.spawn('rpm --resign ' + package)
            if (debug == 1): child.logfile = sys.stdout
            child.expect ('Enter pass phrase|Entrez la phrase de passe')
            child.sendline (passphrase)
            if (debug == 1): child.logfile = sys.stdout
            child.expect(pexpect.EOF)
            child.close()
            # Check return status
            if (child.exitstatus == 0):
                print '[ ' + bcolors.OK + 'OK' + bcolors.END + ' ] '
            else:
                print '[ ' + bcolors.FAIL + 'FAIL' + bcolors.END + ']'

    # Package push
    def push(rpm_files,repository,login,password,folder,debug,charspacing):
        if (debug == 1): print 80*'=' + "\n"
        print '[+] Adding packages to satellite server :'
        for package in rpm_files:
            # Formating
            charspace = Template("{0:<$space}")
            print charspace.substitute(space = charspacing).format(' - ' + package + '...'),
            # RPM push command
            child = pexpect.spawn('rhnpush --force --no-cache -c ' + repository + ' ' + package)
            if (debug == 1): child.logfile = sys.stdout
            child.expect ('Red Hat Network username')
            child.sendline (login)
            child.expect ('Red Hat Network password')
            child.sendline (password)
            if (debug == 1): child.logfile = sys.stdout
            child.expect(pexpect.EOF)
            child.close()
            # Check return status
            if (child.exitstatus == 0):
                print '[ ' + bcolors.OK + 'OK' + bcolors.END + ' ] '
            else:
                print '[ ' + bcolors.FAIL + 'FAIL' + bcolors.END + ' ]'

    # Get rpm files list
    rpm_files=glob.glob(folder + '/*.rpm')
    if (debug == 1): print 80*'=' + "\n" + 'RPM found :'
    if (debug == 1): print rpm_files

    # Check if RPM were found
    if (len(rpm_files) == 0):
       print "No RPM were found in " + folder
       sys.exit(2)

    # Get maximum rpm size for visual answers (OK/FAIL)
    charspacing=0
    for package in rpm_files:
        count = len(package)
        if (count > charspacing):
            charspacing=count
    charspacing += 10

    # Sign packages
    sign(rpm_files,passphrase,folder,debug,charspacing)
    # Push packages
    push(rpm_files,repository,login,password,folder,debug,charspacing)

# Main
def main(argv):
    try:
        opts, args = getopt.getopt(argv, 'hs:r:l:p:f:d', ["passphrase=","repository=","login=","password=","folder=","help"])
    except getopt.GetoptError:
        # Print help and exit
        print "Unknow option, bad or missing argument\n"
        help(2)

    # Initialize vars
    # GPG passphrase for package sign in
    passphrase=None
    repository=None
    login=None
    password=None
    folder='/tmp/'
    debug=0

    # Check opts
    for opt, arg in opts:
        if opt in ("-h", "--help"):
            help(0)
            sys.exit(0)
        elif opt in ("-s", "--passphrase"):
            passphrase = str(arg)
        elif opt in ("-r", "--repository"):
            repository=str(arg)
        elif opt in ("-l", "--login"):
            login=str(arg)
        elif opt in ("-p", "--password"):
            password=str(arg)
        elif opt in ("-f", "--folder"):
            folder=str(arg)
        elif opt in ("-d", "--debug"):
            debug=1
        else:
            print "Unknow option, please see usage\n"
            help(2)

    # Checks
    if (passphrase or repository or login or password) is None:
        print "Unknow option, please see usage\n"
        help(2)

    sign_push(passphrase,repository,login,password,folder,debug)

if __name__ == "__main__":
   main(sys.argv[1:])
```
