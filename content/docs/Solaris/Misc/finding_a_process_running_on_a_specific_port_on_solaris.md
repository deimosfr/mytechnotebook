---
weight: 999
url: "/Trouver_le_process_qui_tourne_sur_un_certain_port_sur_Solaris/"
title: "Finding a Process Running on a Specific Port on Solaris"
description: "Methods to find which process is using a specific port on Solaris operating system."
categories: 
  - Linux
  - Solaris
date: "2010-08-05T13:44:00+02:00"
lastmod: "2010-08-05T13:44:00+02:00"
tags: 
  - Network
  - Servers
  - Solaris
toc: true
---

## Introduction

Since `netstat -auntpl` doesn't exist on Solaris, I had to do some research to find out how to determine which process is listening on a specific port.

## Solutions 1

You'll need to create a small script:

(`get_process_from_port.sh`)

```bash
#!/bin/ksh

line='---------------------------------------------'
pids=$(/usr/bin/ps -ef | sed 1d | awk '{print $2}')

if [ $# -eq 0 ]; then
   read ans?"Enter port you would like to know pid for: "
else
   ans=$1
fi

for f in $pids
do
   /usr/proc/bin/pfiles $f 2>/dev/null | /usr/xpg4/bin/grep -q "port: $ans"
   if [ $? -eq 0 ]; then
      echo $line
      echo "Port: $ans is being used by PID:\c"
      /usr/bin/ps -ef -o pid -o args | egrep -v "grep|pfiles" | grep $f
   fi
done
exit 0
```

Run the script and enter the port number when prompted.

## Solution 2

```bash
fuser -n tcp port
```
