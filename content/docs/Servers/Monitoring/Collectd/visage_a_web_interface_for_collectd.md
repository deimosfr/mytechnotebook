---
weight: 999
url: "/Visage_\\:_Une_interface_web_pour_Collectd/"
title: "Visage: A Web Interface for Collectd"
description: "A guide for installing and configuring Visage, a web interface for Collectd that allows visualization and comparison of collected metrics."
categories: ["Linux", "Debian"]
date: "2011-06-08T15:40:00+02:00"
lastmod: "2011-06-08T15:40:00+02:00"
tags: ["Collectd", "Monitoring", "Visualization", "Web Interface", "Ruby"]
toc: true
---

## Introduction

[Visage](https://auxesis.github.com/visage/) is currently the best interface for [Collectd]({{< ref "docs/Servers/Monitoring/Collectd/collectd_installation_and_configuration.md">}}). It still needs many new features, but at present it's a superb interface that allows comparative analysis between metrics.

## Installation

On Debian, we'll install the prerequisites:

```bash
aptitude install build-essential librrd-ruby ruby ruby-dev rubygems libsinatra-ruby collectd
```

Then via gem (as it's not currently packaged), we'll install it:

```bash
gem install visage-app
```

## Launching

To launch it, it's very simple:

```bash
$(dirname $(dirname $(gem which visage-app)))/bin/visage-app start
```

or

```bash
/var/lib/gems/1.8/gems/visage-app-0.9.4/bin/visage-app start
```

Which will show you this:

```
 _    ___
| |  / (_)________ _____ ____
| | / / / ___/ __ `/ __ `/ _ \
| |/ / (__  ) /_/ / /_/ /  __/
|___/_/____/\__,_/\__, /\___/
                 /____/

will be running at http://nala.deimos.fr:9292/

Looking for RRDs in /var/lib/collectd/rrd

[2011-04-23 19:19:14] INFO  WEBrick 1.3.1
[2011-04-23 19:19:14] INFO  ruby 1.8.7 (2010-08-16) [x86_64-linux]
[2011-04-23 19:19:14] INFO  WEBrick::HTTPServer#start: pid=5409 port=9292
```

No further explanation needed, just connect to the web interface indicated.

## Automatic Host Updates

As of this writing, host profiles are not automatically created in Visage as they might be in competing interfaces. I'm not sure if this is intentional or a missing feature.

Whatever the reason, I have so many hosts and add so many that I need it automated, which is why I created a [small script](https://www.deimos.fr/gitweb/) that updates my Visage profiles:

```bash {linenos=table}
#!/bin/sh
# Visage Update Profiles
# This script will automatically update visage to profile to match with current available rrd graphs
# Made by Pierre Mavro

# Config folders and files
visage_profile_yaml='/usr/lib/ruby/gems/1.8/gems/visage-app-0.9.4/lib/visage-app/config/profiles.yaml'
collectd_rrd='/var/lib/collectd/rrd'

############# DO NOT EDIT #############
cd $collectd_rrd
echo "--- " > $visage_profile_yaml
for profile in `find . -maxdepth 1 -mindepth 1 -type d` ; do
    profile=`echo "$profile" | sed "s/\.\///g"`
    profile_lower=`echo "$profile" | tr "[:upper:]" "[:lower:]"`
    profile_yaml=`echo "$profile_lower" | sed "s/-/+/g"`
    cat <<EOT >>$visage_profile_yaml
$profile_yaml:
  :url: $profile_yaml
  :profile_name: $profile_lower
  :hosts: $profile
  :metrics: "*"
EOT
done
```

You only need to modify the first two parameters to make it work on your platform (default is for Debian Squeeze) and add it to crontab if you want it periodically updated.

## References

http://auxesis.github.com/visage/
