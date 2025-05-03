---
weight: 999
url: "/Puppet_Dashboard_\\:_Mise_en_place_d'une_interface_graphique_pour_Puppet/"
title: "Puppet Dashboard: Setting up a Graphical Interface for Puppet"
description: "How to install, configure and use Puppet Dashboard to create a graphical interface to monitor Puppet nodes and reports."
categories: ["RHEL", "Nginx", "Debian"]
date: "2013-01-05T22:41:00+02:00"
lastmod: "2013-01-05T22:41:00+02:00"
tags: ["Puppet", "Dashboard", "MySQL", "Crontab", "Monitoring"]
toc: true
---

![Puppet Dashboard](/images/puppet-short.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.2.17 |
| **Operating System** | Debian 7 |
| **Website** | [Puppet Dashboard Website](https://puppetlabs.com/puppet/related-projects/dashboard/) |
| **Last Update** | 05/01/2013 |
| **Others** | Clients OS:<br />Debian 6/7<br />RHEL 6 |
{{< /table >}}

## Introduction

[Puppet](./puppet_:_solution_de_gestion_de_fichier_de_configuration.html) is great, but a web interface would be amazing to see the status of machines, synchronizations, etc.

So I suggest setting up [Puppet Dashboard](https://puppetlabs.com/puppet/related-projects/dashboard/).

## Installation

For Debian, we'll use the official repository:

```bash
wget http://apt.puppetlabs.com/puppetlabs-release-stable.deb
dpkg -i puppetlabs-release-stable.deb
```

And then, we update:

```bash
aptitude update
```

Next, we can simply install the dashboard with all its dependencies:

```bash
aptitude install puppet-dashboard mysql-server
```

## Configuration

### Daemon

We'll enable Puppet Dashboard to run automatically at startup by uncommenting the "start" line:

```bash {linenos=table,hl_lines=[6]}
# IMPORTANT: Be sure you have checked the values below, appropriately
# configured 'config/database.yml' in your DASHBOARD_HOME, and
# created and migrated the database.

# Uncomment the line below to start Puppet Dashboard.
START=yes

# Location where puppet-dashboard is installed:
DASHBOARD_HOME=/usr/share/puppet-dashboard

# User which runs the puppet-dashboard program:
DASHBOARD_USER=www-data

# Ruby version to run the puppet-dashboard as:
DASHBOARD_RUBY=/usr/bin/ruby

# Rails environment in which puppet-dashboard runs:
DASHBOARD_ENVIRONMENT=production

# Network interface which puppet-dashboard web server is running at:
DASHBOARD_IFACE=0.0.0.0

# Port on which puppet-dashboard web server is running at, note that if the
# puppet-dashboard user is not root, it has to be a > 1024:
DASHBOARD_PORT=3000
```

We'll also use Delayed Job Workers to ensure data arrives in full even when there's high demand on the Dashboard:

```bash {linenos=table,hl_lines=[5]}
# IMPORTANT: Be sure you have checked the values below, appropriately
# configured 'config/database.yml' in your DASHBOARD_HOME, and
# created and migrated the database.
. /etc/default/puppet-dashboard
START=yes

#  Number of dashboard workers to start.  This will be the number of jobs that
#  can be concurrently processed.  A simple recommendation would be to start
#  with the number of cores you have available.
NUM_DELAYED_JOB_WORKERS=2
```

You can also adjust the NUM_DELAYED_JOB_WORKERS parameter if needed.

### MySQL

To start, we need to initialize the database with the mysql*secure_installation command ([for more information, see this documentation](./mysql*:\_installation_et_configuration.html)). Now we can create a MySQL database and a dedicated user:

```
CREATE DATABASE puppet_dashboard CHARACTER SET utf8;
CREATE USER 'puppetdash_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON puppet_dashboard.* TO 'puppetdash_user'@'localhost';
flush privileges;
```

We'll modify the MySQL configuration to increase the maximum packet size by adjusting a value. Edit your MySQL configuration in the 'mysqld' section and set max_allowed_packet to at least 32M:

```bash {linenos=table,hl_lines=[3,4,5]}
[...]
[mysqld]
# Puppet Dashboard requirements:
# Allowing 32MB allows an occasional 17MB row with plenty of spare room
max_allowed_packet = 32M
[...]
```

Then restart MySQL.

Now we'll modify the default configuration to match our new MySQL database and user. Edit the production section in the following configuration file:

```
...
production:
  database: puppet_dashboard
  username: puppetdash_user
  password: password
  host: localhost
  encoding: utf8
  adapter: mysql
...
```

Now we can initialize the database:

```bash
cd /usr/share/puppet-dashboard
rake RAILS_ENV=production db:migrate
```

You can now start the puppet-dashboard service if you want, and the console is accessible at [http://puppet-dashboard:3000](http://puppet-dashboard:3000)

### Puppet Master

For Puppet Master, we need to tell it to send reports not just as files (in the /var/lib/puppet/reports directory), but also to the MySQL database. To do this, edit the configuration file and add this:

```bash {linenos=table,hl_lines=[5,7]}
[...]
[master]
reportdir = /var/lib/puppet/reports
reporturl = http://localhost:3000/reports/upload
reports = http,store,log
node_terminus = exec
external_nodes = /usr/bin/env PUPPET_DASHBOARD_URL=http://localhost:3000 /usr/share/puppet-dashboard/bin/external_node
[...]
```

Replace _localhost_ with your server name.

Then restart the puppetmaster and puppet-dashboard services.

### Puppet Clients

For the client part, we need to tell it to send a report to the server:

```bash
...
[agent]
report = true
...
```

### Nginx

If like me you use your Puppet Dashboard on the same machine as Puppet Master, it's cleaner to hide Puppet Dashboard's port 3000. For this, we'll use a proxy with Nginx:

```bash {linenos=table,hl_lines=[1,11]}
upstream puppet-prd-dash.deimos.fr:3000 {
  server unix:/usr/share/puppet-dashboard/tmp/sockets/dashboard.0.sock;
  server unix:/usr/share/puppet-dashboard/tmp/sockets/dashboard.1.sock;
  server unix:/usr/share/puppet-dashboard/tmp/sockets/dashboard.2.sock;
}

server {
  root /usr/share/puppet-dashboard/public;

  location / {
    proxy_pass http://puppet-prd-dash.deimos.fr:3000;
  }
}
```

And we enable the new configuration:

```bash
cd /etc/nginx/sites-available
ln -s /etc/nginx/sites-enabled/puppet-dashboard .
/etc/init.d/nginx restart
```

Then restart Nginx.

### Crontab

To avoid performance issues as the database fills up, it's better to purge and optimize it afterwards. We'll create a crontab for this work to be done every month:

```bash
#!/bin/sh
#
# puppet-dashboard cron monthly

set -e

PUPPETDASH_HOME=/usr/share/puppet-dashboard

cd $PUPPETDASH_HOME
# Flush DB old reports older than 1 month
su - www-data -c "cd $PUPPETDASH_HOME ; rake RAILS_ENV=production reports:prune upto=1 unit=mon"
# Optmize table
rake RAILS_ENV=production db:raw:optimize

exit 0
```

We assign it the correct permissions:

```bash
chmod 755 /etc/cron.monthly/puppet-dashboard
```

## Importations

### Puppet Nodes

If you want to import all Puppet nodes at once into your Dashboard without waiting for synchronization:

```bash
cd /usr/share/puppet-dashboard
for i in $(puppetca -la | awk -F\" '{ print $2 }' | grep -v `hostname`) ; do rake RAILS_ENV=production node:add name=$i ; done
```

### Reports

If you want to import the file reports you currently have into the database, it's simple:

```bash
cd /usr/share/puppet-dashboard
rake RAILS_ENV=production reports:import REPORT_DIR=/var/lib/puppet/reports
```

## FAQ

### cannot load such file -- ftools

If you get errors like this when importing the database schema:

```bash
> rake RAILS_ENV=production db:migrate
NOTE: Gem.source_index is deprecated, use Specification. It will be removed on or after 2011-11-01.
Gem.source_index called from /usr/share/puppet-dashboard/vendor/rails/railties/lib/rails/gem_dependency.rb:21.
NOTE: Gem::SourceIndex#initialize is deprecated with no replacement. It will be removed on or after 2011-11-01.
Gem::SourceIndex#initialize called from /usr/share/puppet-dashboard/vendor/rails/railties/lib/rails/vendor_gem_source_index.rb:100.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
NOTE: Gem::SourceIndex#add_spec is deprecated, use Specification.add_spec. It will be removed on or after 2011-11-01.
Gem::SourceIndex#add_spec called from /usr/lib/ruby/1.9.1/rubygems/source_index.rb:91.
rake aborted!
cannot load such file -- ftools

(See full trace by running task with --trace)
```

It's because the Ruby version you're using doesn't match what Puppet Dashboard requires. On Debian 7, you're using 1.9.1 by default and you need to switch to 1.8 (hopefully they'll update the dashboard soon). We'll install some prerequisites:

```bash
aptitude install -y build-essential irb libmysql-ruby libmysqlclient-dev libopenssl-ruby libreadline-ruby mysql-server rake rdoc ri ruby ruby-dev
```

Then we switch to Ruby 1.8:

```bash
update-alternatives --install /usr/bin/gem gem /usr/bin/gem1.8 1
rm /etc/alternatives/ruby
ln -s /usr/bin/ruby1.8 /etc/alternatives/ruby
```

Now we'll compile a version of rubygems:

```bash
URL="http://production.cf.rubygems.org/rubygems/rubygems-1.3.7.tgz"
PACKAGE=$(echo $URL | sed "s/\.[^\.]*$//; s/^.*\///")
cd $(mktemp -d /tmp/install_rubygems.XXXXXXXXXX) && \
wget -c -t10 -T20 -q $URL && \
tar xfz $PACKAGE.tgz && \
cd $PACKAGE && \
sudo ruby setup.rb
```

Now you can run [db::migrate](#mysql) again.

### Caught TERM; calling stop

If you encounter this type of error message in Puppet Dashboard when launching Puppet runs from [Mcollective](./mcollective_:_lancez_des_actions_en_parallèle_sur_des_machines_distante.html), you need to work on the puppet manifest, to comment this line:

```ruby {linenos=table,hl_lines=[5,6]}
[...]
    service {
        'puppet-srv' :
            name => 'puppet',
            # Let this line commented if you're using Puppet Dashboard
            #ensure => stopped,
            enable => false
    }
[...]
```

## Resources
- http://www.puppetlabs.com/blog/a-tour-of-puppet-dashboard-0-1-0/
- http://bitcube.co.uk/content/puppet-dashboard-v101-install
- http://www.mogilowski.net/lang/en-us/2011/01/20/puppet-dashboard-reports-ubuntu/
- http://www.craigdunn.org/2010/08/part-3-installing-puppet-dashboard-on-centos-puppet-2-6-1/
