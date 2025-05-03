---
weight: 999
url: "/MRTG_\\:_Monitoring_axé_réseau/"
title: "MRTG: Network-Focused Monitoring"
description: "A guide to installing and configuring MRTG (Multi Router Traffic Grapher) for network traffic monitoring using SNMP."
categories: ["Monitoring", "Linux", "Debian", "Network"]
date: "2011-12-25T20:44:00+02:00"
lastmod: "2011-12-25T20:44:00+02:00"
tags: ["MRTG", "Monitoring", "Network", "SNMP", "Cacti", "Debian"]
toc: true
---

## Introduction

MRTG Multi Router Traffic Grapher (MRTG) is software developed under the GNU/GPL license at the initiative of Tobias Oetiker. This software allows you to create graphs of network traffic. It uses the SNMP protocol to query network equipment such as routers, switches, or servers that have an MIB.

![MRTG](/images/mrtg.avif)

Today, the creators of MRTG have abandoned their project and moved on to [Cacti]({{< ref "docs/Servers/Monitoring/Cacti">}}), which is a much more powerful but also much more complex product.

## Installation

Before installing MRTG, enable SNMP on your router. Then, if you're using Mandrake, run `urpmi mrtg`, and if you're using Debian, run `aptitude install mrtg`. This will install everything MRTG needs.

```bash
aptitude install mrtg
```

On Debian, it asks if you want the configuration file to be readable only by the root user. Select yes to this question for security reasons.

## Configuration

Now we will tell MRTG where the MRTG destination folder is located, which will be used to display your graphs. You need to put this folder in your web server. To do this, we'll use the cfgmaker command. Here's an example that we'll detail:

```bash
cfgmaker --global 'WorkDir: /var/www/mrtg' --output /var/www/mrtg/routeur.cfg public@X.X.X.X
```

- Workdir: This is where your MRTG logs and graphs will be stored (here `/var/www` represents where the web server data is stored).
- `/var/www/mrtg/routeur.cfg`: This corresponds to the name and location where your MRTG configuration file will be stored.
- public@X.X.X.X: You must replace X.X.X.X with the IP of your router (public is the default name used for authentication).

### Creating the Index

Now we need to create the index pages for MRTG. We'll run the indexmaker command which will use the configuration file to generate an index page:

```bash
indexmaker /var/www/mrtg/routeur.cfg >/var/www/mrtg/index.html
```

- `/var/www/mrtg/routeurcfg`: You need to indicate the location of the MRTG configuration file (the same as the one you specified above).
- `/var/www/mrtg/index.html`: Here indicate the location where the index for the MRTG graphs will be located.

### Creating MRTG Pages

Now that the MRTG configuration and indexing is done, you need to run the mrtg command followed by the location of your configuration file so that it will retrieve the necessary information from your router and create graphs:

```bash
mrtg /var/www/mrtg/routeur.cfg
```

Once this command is launched, you might get some small errors. To fix this problem, run this line a few more times (2 or 3 should be enough) until it no longer displays error messages.

### Automating MRTG Graphs

To make the graphs automatically update, you just need to integrate it into the crontab. You should use the crontab of the Apache user.

Log in as www-data and then edit the crontab via the command `crontab -e`. We'll now automate all this by updating the graphs every 5 minutes:

```bash
5/* * * * * mrtg /var/www/mrtg/routeur.cfg
```

Your crontab is now configured. You now need to wait about 30 minutes to see something appear on the graphs.

## Modifications

If you want to change the direction of the MRTG graphs, edit your configuration file (here, routeur.cfg) and at the location Options[_]: growright, bits, uncomment the line (remove the # symbol). This way, the graphs will not go from left to right but from right to left. Edit the file `/var/www/mrtg/routeur.cfg`:

```bash
#  to get bits instead of bytes and graphs growing to the right
Options[_]: growright, bits
```
