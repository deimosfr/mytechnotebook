---
weight: 999
url: "/Heymon_\\:_Une_interface_web_pour_Collectd/"
title: "Heymon: A Web Interface for Collectd"
description: "Guide on how to install and configure Heymon, a web interface for Collectd that allows comparing metrics between different machines."
categories: ["Linux", "Debian", "Monitoring"]
date: "2010-07-28T15:03:00+02:00"
lastmod: "2010-07-28T15:03:00+02:00"
tags: ["Collectd", "Monitoring", "Web Interface", "Ruby", "Database"]
toc: true
---

## Introduction

[Heymon](https://github.com/newobj/heymon) is one of the most advanced interfaces currently available for Collectd. In my opinion, it's complementary to other interfaces since it allows comparisons between different machines. However, it's quite complicated to set up.

## Installation

Let's install everything we need via Debian packages:

```bash
aptitude install unzip librrd-ruby rubygems1.9 libyaml-ruby libzlib-ruby libdbd-sqlite3-ruby mongrel libopenssl-ruby1.8
```

Then we download the project source code:

```bash
cd /var/www
wget "http://github.com/newobj/heymon/zipball/master"
unzip newobj-heymon-25ceb0e.zip
mv newobj-heymon-25ceb0e heymon
cd heymon
```

Next, we'll need to install gem if you don't have it already (or update it if it's already the case):

```bash
wget "http://rubyforge.org/frs/download.php/70697/rubygems-1.3.7.zip"
unzip rubygems-1.3.7.zip
cd rubygems-1.3.7
ruby1.9 ./setup.rb
```

And install some Ruby modules:

```bash
gem install rake
gem install right_aws
gem install haml
gem install -v=2.3.5 rails
```

## Configuration

Now we'll configure what's needed to run Heymon. Let's generate what's necessary for the SQLite database:

```bash
cd ..
rake db:migrate
```

Edit the configuration file and adapt it if needed (`/var/www/heymon/config/environment.rb`):

```bash
...
RAILS_GEM_VERSION = '2.3.5' unless defined? RAILS_GEM_VERSION
COLLECTD_RRD = '/var/lib/collectd/rrd/'
RRDTOOL_BIN = '/usr/bin/rrdtool'
...
```

## Launching

Now all that's left is to launch the application:

```bash
/var/www/heymon/script/server -d
```

Now try to access the following URL: http://192.168.0.48:3000.
