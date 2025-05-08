---
title: "Netcat: Creating a Listening Port"
description: "How to use Netcat to create a listening port for testing firewall configurations and network connections"
categories: ["Linux", "Network"]
date: "2012-04-27T07:55:00+02:00"
lastmod: "2012-04-27T07:55:00+02:00"
tags: ["Network", "Servers", "Linux"]
---

## Creating a Listening Port

Having a server that listens is good, but when testing with a firewall and you don't necessarily have the listening port behind it yet, you can simply use the Netcat command.

This command will act as if it were a server-type service that starts listening:

```bash
nc -l numero_du_port
```

or

```bash
nc -lp numero_du_port
```

The command varies depending on the version of netcat.

If you want the connection to remain open:

```bash
nc -lk numero_du_port
```

Then all you need to do is test it:

```bash
nc @IP numero_du_port
```

## Resources
- [Netcat: Remote Partition Backup](netcat_remote_partition_backup.md)
- [Netcat: File Transfer](netcat_file_transfer.md)
- [Netcat Documentation](../../../static/pdf/netcat.pdf)
- [Useful Uses Of netcat](../../../static/pdf/useful_uses_of_netcat.pdf)
