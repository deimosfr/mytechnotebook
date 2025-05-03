---
weight: 999
url: "/SABnzbd_\\:_Une_interface_web_pour_gérer_les_newsgroups/"
title: "SABnzbd: A Web Interface for Managing Newsgroups"
description: "How to setup SABnzbd, a web interface for managing newsgroups downloads on Debian systems"
categories: ["Debian", "Linux", "Ubuntu"]
date: "2013-06-19T11:48:00+02:00"
lastmod: "2013-06-19T11:48:00+02:00"
tags: ["Network", "Servers", "Firefox", "Apache", "Python"]
toc: true
---

![Sabnzbd](/images/sabnzbd_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 6/7 |
| **Website** | [Sabnzbd Website](https://sabnzbd.org/) |
| **Last Update** | 19/06/2013 |
{{< /table >}}

## Introduction

You may not want to have heavy software running on your machine for newsgroups. That's when a friend told me about [SABnzbd](https://sabnzbd.org/). After looking at how it works, it turns out that it's particularly powerful.

I decided to write a small documentation for setting it up on Debian Squeeze.

## Installation

Here's the minimum to install for Debian 6:

```bash
aptitude install python-cheetah python-pyopenssl python-yenc
```

For Debian 7:

```bash
aptitude install python-cheetah python-openssl python-yenc
```

Then we'll install SABnzbd:

```bash
cd /var/www
wget "http://downloads.sourceforge.net/project/sabnzbdplus/sabnzbdplus/sabnzbd-0.6.15/SABnzbd-0.6.15-src.tar.gz?r=http%3A%2F%2Fsabnzbd.org%2Fdownload%2F&ts=1332934860&use_mirror=freefr" -o sab.tgz
tar -xzf sab.tgz
mv SABnzbd-0.6.15 sabnzbd
cd sabnzbd
```

## Configuration

### Apache

You may not want to have to type another port to access the web interface. The goal is to put it behind Apache which will redirect requests itself. First, we need to activate the Apache proxy module:

```bash
aptitude install apache2 apache2-utils apache2.2-common libapache2-mod-proxy-html
```

Then activate modules:

```bash
a2enmod proxy_connect
a2enmod proxy_http
a2enmod proxy_html
```

And restart Apache.

Here's the configuration to apply and adapt to your needs:

```apache
[...]
<Location /sabnzbd>
order deny,allow
deny from all
allow from all
ProxyPass http://localhost:8080/sabnzbd
ProxyPassReverse http://localhost:8080/sabnzbd
</Location>
```

## Launch

To start the server:

```bash
python ./SABnzbd.py -d
```

Then connect to http://localhost/SABnzbd (or http://localhost:8080/SABnzbd if you haven't set up Apache) and configure as you wish.

## Post-processing Scripts

You might want an action to happen automatically after each download. In my case, I want the files present to not stay more than 15 days. So I have a script that cleans up, which I placed in crontab:

```bash
#!/bin/sh
folder="/mnt/sabnzbd/downloads"
max_day=15

for user in `ls $folder/` ; do
    cd $folder
    find . -mtime +$max_day ! -name '.' ! -exec rm -f {} \;
    find . -type d ! -name '.' ! -name '..' -exec rmdir {} \;
done
```

The problem is that it removes files that were created more than 15 days ago but that I downloaded less than 15 days ago.

Example: I download an Ubuntu ISO that is 5 months old. It's in a compressed file (zip for example) that I download, Sabnzbd takes care of decompressing it, but the final file dates from 5 months ago and not from the date of decompression. So we need a solution to modify the file date to today. Let's create a scripts folder for example:

```bash
mkdir -p /etc/scripts/sabnzbd
```

Then we'll add a post-processing script:

```bash
#!/bin/sh
tmp_renamed_folder=`echo "$1" | perl -pe 's/\ |-/_/g'`
renamed_folder=`echo "$tmp_renamed_folder" | perl -pe 's/_+/_/g'`
mv "$1" $renamed_folder 2>/dev/null
find $renamed_folder -exec touch {} \;
exit 0
```

This script renames folders containing a - or spaces with _. And it concatenates multiple _ into a single one. Then it changes the date of files and folders. Why do this when renaming options are available directly in SABnzbd? Simply because these operations are done when importing an nzb and not when renaming. But since it causes problems for the find command, I chose to do renaming.

Set the proper permissions:

```bash
chmod 755 /etc/scripts/sabnzbd/postscript.sh
```

To understand how it works, I invite you to read [this link](https://wiki.sabnzbd.org/user-scripts)[^1]

Then, configure Sabnzbd (from the web interface) so it knows where your scripts are and uses them automatically:

1. Configuration -> Directories -> Post-Processing Scripts Folder -> /etc/scripts/sabnzbd
2. Configuration -> Categories -> Default -> Script -> postscript.sh

And there you go :-). Your downloads will be post-processed and will have the date and time of the end of the file download :-)

## Improving Performance

### Sabnzbd

By default, downloading and yencode happens on disk. It is possible to do everything in RAM. Obviously this consumes RAM, but the difference in speed is phenomenal! That means you will at least double the performance! Personally, I quadrupled the download speeds :-). Here's how I did it, go to the graphical interface then:
Configuration -> General -> Settings -> Article Cache Limit -> 250M

### System

At the system level, [check out the optimization of extX filesystems and RAID under Linux]({{< ref "docs/Linux/FilesystemsAndStorage/Raid/optimization_of_extx_filesystems_and_raid_under_linux.md" >}})

## Third-party Tools

### Firefox

I recommend [nzbdStatus](https://addons.mozilla.org/en-US/firefox/addon/nzbdstatus/?src=external-sabfront) for Firefox and [SABMobile for iPhone](https://itunes.apple.com/fr/app/sabmobile/id392506842?mt=8) or for [Android](https://play.google.com/store/apps/details?id=com.patey.SABMobile&hl=fr).

### sabnzbd_api.py

I wrote a Python script to work with the API provided by Sabnzbd. I created it with the idea that it should be easy to add functionality via the API. This script allows:

- To see the elements provided by the XML API
- To delete a certain number of elements in the history

Here's the script. Edit the beginning of the script and insert the URL, as well as the API key:

```python
#!/usr/bin/env python
# Made by Deimos <xxx@mycompany.com>
# Under GPL2 licence

# Set Sabnzbd URL
sabnzbd_url = 'http://127.0.0.1:8080/sabnzbd'
# Set Sabnzbd API
sabnzbd_api = '2d872eb6257123a9579da906b487e1de'

#############################################################

# Load modules
import argparse
from urllib2 import Request, urlopen, URLError
from lxml import etree
import sys

def Debug(debug_text):
    """Debug function"""
    if debug_mode == True:
        print debug_text

def args():
    """Command line parameters"""
    # Define globla vars
    global sabnzbd_url, sabnzbd_api, debug_mode, mode, keep_history

    # Main informations
    parser = argparse.ArgumentParser(description="Sabnzbd API usage in python")
    subparsers = parser.add_subparsers(title='Available sub-commands', help='Choose a subcommand (-h for help)', description='Set a valid subcommand', dest="subcmd_name")
    # Default args
    parser.add_argument('-u', '--url', action='store', dest='url', type=str, default=sabnzbd_url, help='Set Sabnzbd URL (default : %(default)s)')
    parser.add_argument('-a', '--api', action='store', dest='api', type=str, default=sabnzbd_api, help='Set Sabnzbd API key')
    parser.add_argument('-v', '--version', action='version', version='v0.1 Licence GPLv2', help='Version 0.1')
    parser.add_argument('-d', '--debug', action='store_true', default=False, help='Debug mode')
    # Show XML API
    parser_xa = subparsers.add_parser('xa', help='Show XML Elements from API')
    # Delete history
    parser_dh = subparsers.add_parser('dh', help='Delete old history')
    parser_dh.add_argument('-k', '--keep', action='store', dest='num', type=int, default=150, help='Number of items to keep in history (default : %(default)s)')

    result = parser.parse_args()

    # Set debug to True if requested by command line
    if (result.debug == True):
        debug_mode=True
    else:
        debug_mode=False
    # Send debug informations
    Debug('Command line : ' + str(sys.argv))
    Debug('Command line vars : ' + str(parser.parse_args()))

    # Check defaults options
    if result.url != sabnzbd_url:
        sabnzbd_url = result.url
    elif result.api != sabnzbd_api:
        sabnzbd_api = result.api

    # Managing options
    mode=result.subcmd_name
    if (mode == 'dh'):
        keep_history=result.num

def GetSabUrl(sabnzbd_url, sabnzbd_api):
    """Concat Sabnzbd URL
    Check that the connection to URL is correct and TCP is OK"""
    # Concat
    sab_url = sabnzbd_url + '/api?mode=history&output=xml&apikey=' + sabnzbd_api
    Debug('Concatened Sabnzbd URL : ' + str(sab_url))

    # Connectivity test
    request = Request(sab_url)
    try:
        response = urlopen(request)
    except URLError, e:
        if hasattr(e, 'reason'):
            print 'Can\'t connect to server : ' + str(e.reason)
            sys.exit(1)
        elif hasattr(e, 'code'):
            print 'The server couldn\'t fulfill the request.' + str(e.code)
            sys.exit(1)
    else:
        Debug('Sabnzbd TCP connection OK')
    return sab_url

def GetXmlHistory(sab_url):
    """Get XML nzo_id entries"""
    xml_history =  []
    # Parse XML from given url
    xml_content = etree.parse(sab_url)
    # Select only the wished tag and push it in xml_history list
    for node in xml_content.xpath("///*"):
        if (node.tag == 'nzo_id'):
            # Reencode value to avoid errors
            tag_value = unicode(node.text).encode('utf8')
            #Debug(node.tag + ' : ' + tag_value)
            xml_history.append(tag_value)
    Debug('XML history has ' + str(len(xml_history)) + ' entries')

    # If there were no entry in the list, check if the API failed
    if len(xml_history) == 0:
        Debug('Checking why there is no datas')
        xml_content = etree.parse(sab_url)
        for node in xml_content.xpath("//*"):
            if (node.tag == 'error'):
                print 'Can\'t connect to server : ' + unicode(node.text).encode('utf8')
                sys.exit(1)

    return xml_history

def DeleteHistory(xml_history,keep_history):
    """Delete old history"""
    # Create a new list contaning elements to remove
    elements2delete = []
    element_number=1
    for element in xml_history:
        if element_number > keep_history:
            elements2delete.append(element)
        element_number += 1
    queue_elements2remove=','.join(elements2delete)
    Debug('Elements to delete (' + str(len(xml_history)) + ' - ' + str(keep_history)  + ' = ' + str(len(elements2delete))  + ')')

    # Concat URL with elements to remove and then remove
    sab_url_remove = sabnzbd_url + '/api?mode=history&name=delete&value=' + queue_elements2remove + '&apikey=' + sabnzbd_api
    urlopen(sab_url_remove)

def ShowXMLElements(sab_url):
    """Show XML code from URL"""
    xml_content = etree.parse(sab_url)
    for node in xml_content.xpath("//*"):
        print node.tag + ' : ' + unicode(node.text).encode('utf8')

def main():
    """Main function that launch all of them"""
    # Set args and get debug mode
    args()
    # Get Sabnzbd URL and check connectivity
    sab_url = GetSabUrl(sabnzbd_url, sabnzbd_api)
    # Delete old history
    if mode == 'dh':
        xml_history = GetXmlHistory(sab_url)
        DeleteHistory(xml_history,keep_history)
    elif mode == 'xa':
        ShowXMLElements(sab_url)

if __name__ == "__main__":
    main()
```

For the usage part, it's divided into 2:

- The main commands
- The subcommands (requested functions)

Here's how to use it:

```bash
> sabnzbd_api.py -h
usage: sabnzbd_api.py [-h] [-u URL] [-a API] [-v] [-d] {xa,dh} ...

Sabnzbd API usage in python

optional arguments:
  -h, --help         show this help message and exit
  -u URL, --url URL  Set Sabnzbd URL (default : http://127.0.0.1:8080/sabnzbd)
  -a API, --api API  Set Sabnzbd API key
  -v, --version      Version 0.1
  -d, --debug        Debug mode

Available sub-commands:
  Set a valid subcommand

  {xa,dh}            Choose a subcommand (-h for help)
    xa               Show XML Elements from API
    dh               Delete old history
```

To use the dh subcommand that allows you to delete history:

```bash
> sabnzbd_api.py dh -h
usage: sabnzbd_api.py dh [-h] [-k NUM]

optional arguments:
  -h, --help          show this help message and exit
  -k NUM, --keep NUM  Number of items to keep in history (default : 150)
```

Here, if I remove -h, it will delete the history queue larger than 150 entries.

## References

[^1]: http://wiki.sabnzbd.org/user-scripts
