---
weight: 999
url: "/Puppet\_\\:\_Solution_de_gestion_de_fichier_de_configuration/"
title: "Puppet: Configuration File Management Solution"
description: "A comprehensive guide to installing and configuring Puppet, a powerful configuration management tool that helps automate system administration tasks."
categories: ["Linux"]
date: "2013-08-02T11:43:00+02:00"
lastmod: "2013-08-02T11:43:00+02:00"
tags: ["Puppet", "Configuration manager"]
toc: true
---

## Introduction

Puppet is a very practical application... It's what you would find in companies with large volumes of servers, where the information system is "industrialized".
Puppet allows you to automate many administration tasks, such as software installation, services deployment, or file modifications.
Puppet allows you to do this in a centralized way, which helps to better manage and control a large number of heterogeneous or homogeneous servers.
Puppet works in Client/Server mode.

On each machine, a client will be installed, which will contact the PuppetMaster (the server) through HTTPS communication, and therefore SSL, with a provided PKI system.
Puppet was developed in Ruby, making it multi-platform: BSD (free, macOS...), Linux (RedHat, Debian, SUSE...), Sun (OpenSolaris...)
[Reductive Labs](https://reductivelabs.com/), the company publishing Puppet, has developed a complementary product called Facter.
This application lists specific elements of the managed systems, such as hostname, IP address, distribution, and environment variables that can be used in Puppet templates.
As Puppet manages templates, you can quickly understand the usefulness of Facter. For example, if you manage a farm of mail servers that require configuration containing the machine name, a template combined with environment variables proves quite useful.
In short, Puppet combined with Facter seems like a very interesting solution to simplify system administration.

Here is a diagram showing how Puppet works:
![Puppet Star.png](/images/puppet_star.avif)

For Puppet configuration, if you want to use an IDE, there is [Geppetto](https://www.puppetlabs.com/blog/geppetto-a-puppet-ide/). I recommend it, as it will save you a lot of syntax troubles.

Documentation for previous versions is available here:

- [Puppet 2.7](/pdf/puppet_2.7.pdf)
- [Puppet 0.25.4](/pdf/puppet_0.25.4.pdf)

## Puppet Hierarchy

Before going further, I've borrowed from the official site how the Puppet directory structure works:

All Puppet data files (modules, manifests, distributable files, etc) should be maintained in a Subversion or CVS repository (or your favorite Version Control System). The following hierarchy describes the layout one should use to arrange the files in a maintainable fashion:

- `/manifests/`: this directory contains files and subdirectories that determine the manifest of individual systems but do not logically belong to any particular module. Generally, this directory is fairly thin and alternatives such as the use of LDAP or other external node tools can make the directory even thinner. This directory contains the following special files:
  - site.pp: first file that the Puppet Master parses when determining a server's catalog. It imports all the underlying subdirectories and the other special files in this directory. It also defines any global defaults, such as package managers. See sample site.pp.
  - templates.pp: defines all template classes. See also terminology:template classes. See sample templates.pp.
  - nodes.pp: defines all the nodes if not using an external node tool. See sample nodes.pp.
- `/modules/{modulename}/`: houses puppet modules in subdirectories with names matching that of the module name. This area defines the general building blocks of a server and contains modules such as for openssh, which will generally define classes openssh::client and openssh::server to setup the client and server respectively. The individual module directories contains subdirectories for manifests, distributable files, and templates. See modules organization, terminology:module.
- `/modules/user/`: A special module that contains manifests for users. This module contains a special subclass called user::virtual which declares all the users that might be on a given system in a virtual way. The other subclasses in the user module are classes for logical groupings, such as user::unixadmins, which will realize the individual users to be included in that group. See also naming conventions, terminology:realize.
- `/services/`: this is an additional modules area that is specified in the module path for the puppetmaster. However, instead of generic modules for individual services and bits of a server, this module area is used to model servers specific to enterprise level infrastructure services (core infrastructure services that your IT department provides, such as www, enterprise directory, file server, etc). Generally, these classes will include the modules out of /modules/ needed as part of the catalog (such as openssh::server, postfix, user::unixadmins, etc). The files section for these modules is used to distribute configuration files specific to the enterprise infrastructure service such as openldap schema files if the module were for the enterprise directory. To avoid namespace collision with the general modules, it is recommended that these modules/classes are prefixed with s\_ (e.g. s_ldap for the enterprise directory server module)
- `/clients/`: similar to the /services/ module area, this area is used for modules related to modeling servers for external clients (departments outside your IT department). To avoid namespace collision, it is recommended that these modules/classes are prefixed with c\_.
- `/notes/`: this directory contains notes for reference by local administrators.
- `/plugins/`: contains custom types programmed in Ruby. See also terminology:plugin-type.
- `/tools/`: contains scripts useful to the maintenance of Puppet.

## Installation

### Puppet Server

The master version used must be the same as that of the client machines. It is highly recommended to use a version greater than or equal to 0.25.4 (which fixes numerous performance issues).
For this, on Debian, you'll need to install the version available in squeeze/lenny-backport or higher, and lock it to prevent an accidental upgrade from changing its version (use "pin locks"). Here we will opt for the version given on the official Puppet site.

For now, we need to configure the /etc/hosts file with the server IP:

```bash
...
192.168.0.93    puppet-prd.deimos.fr puppet
...
```

_Note: Check that the puppetmaster's clock (and client clocks as well) is up-to-date/synchronized. There can be an issue with certificates not being recognized/accepted if there is a time discrepancy (run `dpkg-reconfigure tz-data`)._

Configure the official Puppet repository if you want the latest version, otherwise skip this step to install the version provided by your distribution:

```bash
wget http://apt.puppetlabs.com/puppetlabs-release-stable.deb
dpkg -i puppetlabs-release-stable.deb
```

And then update:

```bash
aptitude update
```

Then install puppetmaster:

```bash
aptitude install puppetmaster
```

You can verify that puppetmaster is installed correctly by running 'facter' (see if it returns something) or checking for SSL files (in /var/lib/puppet).

### Web Server

Now we need to configure a web server on the same machine as the Puppet server (Puppet Master). Why? Simply because the default server is Webrick and it collapses if 10 nodes access it simultaneously.

{{< alert context="warning" text="You can keep Webrick if you want to test on a few nodes, but not for production!" />}}

The choice is yours between Passenger and Nginx. Passenger is the recommended solution since Puppet 3.

#### Passenger

If you've chosen to use Passenger as recommended by PuppetLab, you need to disable automatic daemon startup, as it would start a web server and conflict with Passenger:

```bash
# Defaults for puppetmaster - sourced by /etc/init.d/puppetmaster

# Start puppetmaster on boot? If you are using passenger, you should
# have this set to "no"
START=no
# Startup options
DAEMON_OPTS=""

# What port should the puppetmaster listen on (default: 8140).
PORT=8140
```

Now disable the service:

```bash
/etc/init.d/puppetmaster stop
```

Then, install Passenger:

```bash
aptitude install puppetmaster-passenger
```

You don't need to do any configuration. Everything is provided by the official Puppet packages. In case you installed it without the official repository, here is the generated configuration:

```apache
# you probably want to tune these settings
PassengerHighPerformance on
PassengerMaxPoolSize 12
PassengerPoolIdleTime 1500
# PassengerMaxRequests 1000
PassengerStatThrottleRate 120
RackAutoDetect Off
RailsAutoDetect Off

Listen 8140

<VirtualHost *:8140>
        SSLEngine on
        SSLProtocol -ALL +SSLv3 +TLSv1
        SSLCipherSuite ALL:!ADH:RC4+RSA:+HIGH:+MEDIUM:-LOW:-SSLv2:-EXP

        SSLCertificateFile      /var/lib/puppet/ssl/certs/puppet.deimos.lan.pem
        SSLCertificateKeyFile   /var/lib/puppet/ssl/private_keys/puppet.deimos.lan.pem
        SSLCertificateChainFile /var/lib/puppet/ssl/certs/ca.pem
        SSLCACertificateFile    /var/lib/puppet/ssl/certs/ca.pem
        # If Apache complains about invalid signatures on the CRL, you can try disabling
        # CRL checking by commenting the next line, but this is not recommended.
        SSLCARevocationFile     /var/lib/puppet/ssl/ca/ca_crl.pem
        SSLVerifyClient optional
        SSLVerifyDepth  1
        # The `ExportCertData` option is needed for agent certificate expiration warnings
        SSLOptions +StdEnvVars +ExportCertData

        # This header needs to be set if using a loadbalancer or proxy
        RequestHeader unset X-Forwarded-For

        RequestHeader set X-SSL-Subject %{SSL_CLIENT_S_DN}e
        RequestHeader set X-Client-DN %{SSL_CLIENT_S_DN}e
        RequestHeader set X-Client-Verify %{SSL_CLIENT_VERIFY}e

        DocumentRoot /usr/share/puppet/rack/puppetmasterd/public/
        RackBaseURI /
        <Directory /usr/share/puppet/rack/puppetmasterd/>
                Options None
                AllowOverride None
                Order allow,deny
                allow from all
        </Directory>
</VirtualHost>
```

#### NGINX and Mongrel

##### Installation

If you've chosen NGINX and Mongrel, you'll need to start by installing them:

```bash
aptitude install nginx mongrel
```

##### Configuration

Modify the /etc/default/puppetmaster file:

```bash
# Defaults for puppetmaster - sourced by /etc/init.d/puppet

# Start puppet on boot?

START=yes

# Startup options

DAEMON_OPTS=""

# What server type to run

# Options:

# webrick (default, cannot handle more than ~30 nodes)

# mongrel (scales better than webrick because you can run

# multiple processes if you are getting

# connection-reset or End-of-file errors, switch to

# mongrel. Requires front-end web-proxy such as

# apache, nginx, or pound)

# See: http://reductivelabs.com/trac/puppet/wiki/UsingMongrel

SERVERTYPE=mongrel

# How many puppetmaster instances to start? Its pointless to set this

# higher than 1 if you are not using mongrel.

PUPPETMASTERS=4

# What port should the puppetmaster listen on (default: 8140). If

# PUPPETMASTERS is set to a number greater than 1, then the port for

# the first puppetmaster will be set to the port listed below, and

# further instances will be incremented by one

#

# NOTE: if you are using mongrel, then you will need to have a

# front-end web-proxy (such as apache, nginx, pound) that takes

# incoming requests on the port your clients are connecting to

# (default is: 8140), and then passes them off to the mongrel

# processes. In this case it is recommended to run your web-proxy on

# port 8140 and change the below number to something else, such as

# 18140.

PORT=18140
```

After (re-)starting the daemon, you should be able to see the attached sockets:

```bash
 > netstat -pvltpn
 Connexions Internet actives (seulement serveurs)
 Proto Recv-Q Send-Q Adresse locale          Adresse distante        Etat        PID/Program name
 tcp        0      0 0.0.0.0:41736           0.0.0.0:*               LISTEN      2029/rpc.statd
 tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN      2018/portmap
 tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      2333/sshd
 tcp        0      0 127.0.0.1:18140         0.0.0.0:*               LISTEN      10059/ruby       tcp        0      0 127.0.0.1:18141         0.0.0.0:*               LISTEN      10082/ruby       tcp        0      0 127.0.0.1:18142         0.0.0.0:*               LISTEN      10104/ruby       tcp        0      0 127.0.0.1:18143         0.0.0.0:*               LISTEN      10126/ruby
 tcp6       0      0 :::22                   :::*                    LISTEN      2333/sshd
```

Add the following lines to the /etc/puppet/puppet.conf file:

```ini
[main]
logdir=/var/log/puppet
vardir=/var/lib/puppet
ssldir=/var/lib/puppet/ssl
rundir=/var/run/puppet
factpath=$vardir/lib/facter
templatedir=$confdir/templates
pluginsync = true
[master]
# These are needed when the puppetmaster is run by passenger
# and can safely be removed if webrick is used.
ssl_client_header = HTTP_X_SSL_SUBJECT
ssl_client_verify_header = SSL_CLIENT_VERIFY
report = true

[agent]
server=puppet-srv.deimos.fr
```

Modify the following configuration in /etc/nginx.conf:

```nginx
user www-data;
worker_processes  4;

error_log  /var/log/nginx/error.log;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
    # multi_accept on;
}

http {
    #include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    access_log /var/log/nginx/access.log;

    sendfile        on;     tcp_nopush     on;

    # Look at TLB size in /proc/cpuinfo (Linux) for the 4k pagesize
    large_client_header_buffers     16      4k;
    proxy_buffers                   128     4k;

    #keepalive_timeout  0;
    keepalive_timeout  65;
    tcp_nodelay        on;

    gzip  on;
    gzip_disable "MSIE [1-6]\.(?!.*SV1)";

    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
```

And add this configuration for puppet:

```nginx
upstream puppet-prd.deimos.fr {
   server 127.0.0.1:18140;
   server 127.0.0.1:18141;
   server 127.0.0.1:18142;
   server 127.0.0.1:18143;
}

server {
   listen                  8140;

   ssl                     on;
   ssl_certificate         /var/lib/puppet/ssl/certs/puppet-prd.pem;
   ssl_certificate_key     /var/lib/puppet/ssl/private_keys/puppet-prd.pem;   ssl_client_certificate  /var/lib/puppet/ssl/ca/ca_crt.pem;
   ssl_ciphers             SSLv2:-LOW:-EXPORT:RC4+RSA;
   ssl_session_cache       shared:SSL:8m;
   ssl_session_timeout     5m;

   ssl_verify_client       optional;

   # obey to the Puppet CRL
   ssl_crl                 /var/lib/puppet/ssl/ca/ca_crl.pem;

   root                    /var/empty;
   access_log              /var/log/nginx/access-8140.log;
   #rewrite_log             /var/log/nginx/rewrite-8140.log;

   # Variables
   # $ssl_cipher returns the line of those utilized it is cipher for established SSL-connection
   # $ssl_client_serial returns the series number of client certificate for established SSL-connection
   # $ssl_client_s_dn returns line subject DN of client certificate for established SSL-connection
   # $ssl_client_i_dn returns line issuer DN of client certificate for established SSL-connection
   # $ssl_protocol returns the protocol of established SSL-connection

   location / {
       proxy_pass          http://puppet-prd.deimos.fr;       proxy_redirect      off;
       proxy_set_header    Host             $host;
       proxy_set_header    X-Real-IP        $remote_addr;
       proxy_set_header    X-Forwarded-For  $proxy_add_x_forwarded_for;
       proxy_set_header    X-Client_DN      $ssl_client_s_dn;
       proxy_set_header    X-Client-Verify  $ssl_client_verify;
       proxy_set_header    X-SSL-Subject    $ssl_client_s_dn;
       proxy_set_header    X-SSL-Issuer     $ssl_client_i_dn;
       proxy_read_timeout  65;

   }
}
```

Then create the symbolic link to apply the configuration:

```bash
cd /etc/nginx/sites-available
ln -s /etc/nginx/sites-enabled/puppetmaster .
```

And then restart the Nginx server.

To verify that the daemons are running correctly, you should have the following sockets open:

```bash
> netstat -vlptn
Connexions Internet actives (seulement serveurs)
Proto Recv-Q Send-Q Adresse locale          Adresse distante        Etat        PID/Program name
tcp        0      0 0.0.0.0:41736           0.0.0.0:*               LISTEN      2029/rpc.statd
tcp        0      0 0.0.0.0:8140            0.0.0.0:*               LISTEN      10293/nginx
tcp        0      0 0.0.0.0:8141            0.0.0.0:*               LISTEN      10293/nginx
tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN      2018/portmap
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      2333/sshd
tcp        0      0 127.0.0.1:18140         0.0.0.0:*               LISTEN      10059/ruby
tcp        0      0 127.0.0.1:18141         0.0.0.0:*               LISTEN      10082/ruby
tcp        0      0 127.0.0.1:18142         0.0.0.0:*               LISTEN      10104/ruby
tcp        0      0 127.0.0.1:18143         0.0.0.0:*               LISTEN      10126/ruby
tcp6       0      0 :::22                   :::*                    LISTEN      2333/sshd
```

### Puppet Clients

For clients, it's also simple. But first add the server line to the hosts file:

```bash
...
192.168.0.93        puppet-prd.deimos.fr       puppet
```

This is not mandatory if your DNS names are correctly configured.

#### Debian

If you want to use the latest version:

```bash
wget http://apt.puppetlabs.com/puppetlabs-release-stable.deb
dpkg -i puppetlabs-release-stable.deb
```

And then update:

```bash
aptitude update
```

Verify that the /etc/hosts file contains the hostname of the client machine, then install puppet:

```bash
aptitude install puppet
```

#### Red Hat

Just like Debian, there is a yum repo on Red Hat and we'll install a package that will configure it for us:

```bash
rpm -ivh http://yum.puppetlabs.com/el/6/products/x86_64/puppetlabs-release-6-6.noarch.rpm
```

Then install:

```bash
yum install puppet
```

#### Solaris

The stable Puppet client in blastwave is too old (0.23). Therefore, Puppet (and Facter) will need to be installed through Ruby's standard manager: gem.
To do this, you'll first need to install ruby with the following command:

```bash
pkg-get -i ruby
```

{{< alert context="warning" text="Check that rubygems is not already installed, otherwise remove it:\n\n```bash\npkg-get -r rubygems\n```" />}}

Then install a more up-to-date version from the sources:

```bash
wget http://rubyforge.org/frs/download.php/45905/rubygems-1.3.1.tgz
gzcat rubygems-1.3.1.tgz | tar -xf -
cd rubygems-1.3.1
ruby setup.rb
gem --version
```

Install puppet with the command and the -p argument if you have a proxy:

```bash
gem install puppet --version '0.25.4' -p http://proxy:3128/
```

You need to modify/add some commands that are not available by default on Solaris for Puppet to work better:

- Create a link for uname and puppetd:

```bash
ln -s /usr/bin/uname /usr/bin/
ln -s /opt/csw/bin/puppetd /usr/bin/
```

- Create a script called /usr/bin/dnsdomainname:

```bash
#!/usr/bin/bash
DOMAIN="`/usr/bin/domainname 2> /dev/null`"
if [ ! -z "$DOMAIN" ]; then
    echo $DOMAIN | sed 's/^[^.]*.//'
fi
```

```bash
chmod 755 /usr/bin/dnsdomainname
```

Then, the procedure is the same as for other OSes, that is, modify /etc/hosts to include puppet-prd.deimos.fr, and run:

```bash
puppetd --verbose --no-daemon --test --server puppet-prd.deimos.fr
```

At this point, if it doesn't work, it's simply because you need to modify the configuration (puppet.conf) of your client.

## Configuration

### Server

For the server part, here's how the directory structure is organized (in /etc/puppet):

```
.
|-- auth.conf
|-- autosign.conf
|-- fileserver.conf
|-- manifests
|   |-- common.pp
|   |-- modules.pp
|   `-- site.pp
|-- modules
|-- puppet.conf
`-- templates
```

#### auth.conf

This is where we set all the permissions:

```
# This is an example auth.conf file, it mimics the puppetmasterd defaults
#
# The ACL are checked in order of appearance in this file.
#
# Supported syntax:
# This file supports two different syntax depending on how
# you want to express the ACL.
#
# Path syntax (the one used below):
# ---------------------------------
# path /path/to/resource
# [environment envlist]
# [method methodlist]
# [auth[enthicated] {yes|no|on|off|any}]
# allow [host|ip|*]
# deny [host|ip]
#
# The path is matched as a prefix. That is /file match at
# the same time /file_metadat and /file_content.
#
# Regex syntax:
# -------------
# This one is differenciated from the path one by a '~'
#
# path ~ regex
# [environment envlist]
# [method methodlist]
# [auth[enthicated] {yes|no|on|off|any}]
# allow [host|ip|*]
# deny [host|ip]
#
# The regex syntax is the same as ruby ones.
#
# Ex:
# path ~ .pp$
# will match every resource ending in .pp (manifests files for instance)
#
# path ~ ^/path/to/resource
# is essentially equivalent to path /path/to/resource
#
# environment:: restrict an ACL to a specific set of environments
# method:: restrict an ACL to a specific set of methods
# auth:: restrict an ACL to an authenticated or unauthenticated request
# the default when unspecified is to restrict the ACL to authenticated requests
# (ie exactly as if auth yes was present).
#

### Authenticated ACL - those applies only when the client
### has a valid certificate and is thus authenticated

# allow nodes to retrieve their own catalog (ie their configuration)
path ~ ^/catalog/([^/]+)$
method find
allow $1

# allow nodes to retrieve their own node definition
path ~ ^/node/([^/]+)$
method find
allow $1

# allow all nodes to access the certificates services
path /certificate_revocation_list/ca
method find
allow *

# allow all nodes to store their reports
path /report
method save
allow *

# inconditionnally allow access to all files services
# which means in practice that fileserver.conf will
# still be used
path /file
allow *

### Unauthenticated ACL, for clients for which the current master doesn't
### have a valid certificate; we allow authenticated users, too, because
### there isn't a great harm in letting that request through.

# allow access to the master CA
path /certificate/ca
method find
allow *

path /certificate/
method find
allow *

path /certificate_request
method find, save
allow *

# this one is not stricly necessary, but it has the merit
# to show the default policy which is deny everything else
path /
auth any
allow *.deimos.fr
```

In case you encounter access problems, for **insecure but simple testing**, add this line at the end of your configuration file:

```
allow *
```

#### autosign.conf

You can auto-sign certain certificates to save time. This can be a bit dangerous, but if your node filtering is done correctly behind, no worries :-)

```
*.deimos.fr
```

Here I'll auto-sign all my nodes with the deimos.fr domain.

#### fileserver.conf

Give permissions for client machines in the /etc/puppet/fileserver.conf file:

```
# This file consists of arbitrarily named sections/modules
# defining where files are served from and to whom

# Define a section 'files'
# Adapt the allow/deny settings to your needs. Order
# for allow/deny does not matter, allow always takes precedence
# over deny
[files]
  path /etc/puppet/files
  allow *.deimos.fr
#  allow *.example.com
#  deny *.evil.example.com
#  allow 192.168.0.0/24

[plugins]
  allow *.deimos.fr
#  allow *.example.com
#  deny *.evil.example.com
#  allow 192.168.0.0/24
```

#### manifests

Let's create the missing files:

```bash
touch /etc/puppet/manifests/{common.pp,modules.pp,site.pp}
```

##### common.pp

The common.pp is empty, but you can insert things that will be taken as global configuration.

```ruby

```

##### modules.pp

Then I'll define my base module(s) here. For example, in my future configuration I'll declare a "base" module that will contain everything that any machine that's part of puppet will inherit:

```ruby
 import "base"
```

##### site.pp

We ask to load all modules present in the modules folder:

```ruby
# /etc/puppet/manifests/site.pp
import "common.pp"

# The filebucket option allows for file backups to the server
filebucket { main: server => 'puppet-prod-nux.deimos.fr' }

# Backing up all files and ignore vcs files/folders
File {
	backup => '.puppet-bak',
	ignore => ['.svn', '.git', 'CVS' ]
}

# Default global path
Exec { path => "/usr/bin:/usr/sbin/:/bin:/sbin" }

# Import base module
import "modules.pp"
```

Here I tell it to use the filebucket on the puppet server and to rename files that will be replaced by puppet to <file>.puppet-bak.
I also ask it to ignore any directories or files created by VCS like SVN, Git or CVS.
And finally, I indicate the default path that puppet will have when it runs on clients.

All this configuration is specific to the puppet server, so it's global. Anything we put in it can be inherited.
Restart the puppetmaster to make sure that the server-side changes have been properly taken into account.

#### puppet.conf

I intentionally skipped the modules folder as it's the big piece of puppet and will be given special attention later in this article.

So we're going to move on to the puppet.conf configuration file that you probably already configured during the mongrel/nginx installation...:

```ini
[main]
logdir=/var/log/puppet
vardir=/var/lib/puppet
ssldir=/var/lib/puppet/ssl
rundir=/var/run/puppet
factpath=$vardir/lib/facter
templatedir=$confdir/templates
pluginsync = true
[master]
# These are needed when the puppetmaster is run by passenger
# and can safely be removed if webrick is used.
ssl_client_header = HTTP_X_SSL_SUBJECT
ssl_client_verify_header = SSL_CLIENT_VERIFY
report = true

[agent]
server=puppet-srv.deimos.fr
```

For the templates folder, I have nothing in it.

### Client

Each client must have its entry in the DNS server (just like the server)!

#### puppet.conf

##### Debian / Red Hat

The configuration file must contain the server address:

```ini
[main]
    # The Puppet log directory.
    # The default value is '$vardir/log'.
    logdir = /var/log/puppet

    # Where Puppet PID files are kept.
    # The default value is '$vardir/run'.
    rundir = /var/run/puppet

    # Where SSL certificates are kept.
    # The default value is '$confdir/ssl'.
    ssldir = $vardir/ssl

    # Puppet master server
    server = puppet-prd.deimos.fr

    # Add custom facts
    pluginsync = true
    pluginsource = puppet://$server/plugins
    factpath = /var/lib/puppet/lib/facter

[agent]
    # The file in which puppetd stores a list of the classes
    # associated with the retrieved configuratiion.  Can be loaded in
    # the separate ``puppet`` executable using the ``--loadclasses``
    # option.
    # The default value is '$confdir/classes.txt'.
    classfile = $vardir/classes.txt

    # Where puppetd caches the local configuration.  An
    # extension indicating the cache format is added automatically.
    # The default value is '$confdir/localconfig'.
    localconfig = $vardir/localconfig

    # Reporting
    report = true
```

##### Solaris

For Solaris, the configuration needed quite a bit of adaptation:

```ini
[main]
    logdir=/var/log/puppet
    vardir=/var/opt/csw/puppet
    rundir=/var/run/puppet
    # ssldir=/var/lib/puppet/ssl
    ssldir=/etc/puppet/ssl
    # Where 3rd party plugins and modules are installed
    libdir = $vardir/lib
    templatedir=$vardir/templates
    # Turn plug-in synchronization on.
    pluginsync = true
    pluginsource = puppet://$server/plugins
    factpath = /var/puppet/lib/facter

[puppetd]
    report=true
    server=puppet-prd.deimos.fr
    # certname=puppet-prd.deimos.fr
    # enable the marshal config format
    config_format=marshal
    # different run-interval, default= 30min
    # e.g. run puppetd every 4 hours = 14400
    runinterval = 14400
    logdest=/var/log/puppet/puppet.log
```

## The Language

Before starting to create modules, we need to know a bit more about the syntax/language used for puppet. Its syntax is close to ruby and it's even possible to write complete modules in ruby. I'll explain here some techniques/possibilities to allow you to create advanced modules later.

We'll also use types, I won't go into detail on these, because the doc on the site is clear enough: http://docs.puppetlabs.com/references/latest/type.html

### Functions

Here's how to define a function with multiple arguments:

```ruby
define network_config( $ip, $netmask, $gateway ) {
    notify {"$ip, $netmask, $gateway":}
}
network_config { "eth0":
    ip      => '192.168.0.1',
    netmask => '255.255.255.0',
    gateway => '192.168.0.254,
}
```

### Installing packages

We'll see here how to install a package, then use a function to easily install many more. For a single package, it's simple:

```ruby
# Install kexec-tools
package { 'kexec-tools':
    ensure => 'installed'
}
```

Here we're asking for a package (kexec-tool) to be installed. If we want several to be installed, we'll need to create an array:

```ruby
# Install kexec-tools
package {
   [
      'kexec-tools',
      'package2',
      'pacakge3'
   ]:
   ensure => 'installed'
}
```

This is quite practical and arrays often work this way for pretty much any type used. We can also create a function for this in which we'll send each element of the array:

```ruby
# Validate that pacakges are installed
define packages_install () {
    notice("Installation of ${name} package")
    package {
        "${name}":
            ensure => 'installed'
    }
}
# Set all custom packages (not embended in distribution) that need to be installed
packages_install
{ [
   'puppet',
   'tmux'
]: }
```

Some of you will say that for this specific case, it's pointless, since the method above allows it to be done while others will find this method more elegant and easier to comprehend for a novice coming to puppet. The function name used is packages_install, the $name variable is always the first element sent to a function, which corresponds here to each element in our array.

### Including / Excluding modules

You've just seen functions, we'll push them a bit further with a solution to include and exclude the loading of certain modules. Here, I have a functions file:

```ruby
# Load or not modules (include/exclude)
define include_modules () {
    if ($exclude_modules == undef) or !($name in $exclude_modules) {
        include $name
    }
}
```

Here I have another file representing the roles of my servers (we'll get to that later):

```ruby
# Load modules
$minimal_modules =
[
    'puppet',
    'resolvconf',
    'packages_defaults',
    'configurations_defaults'
]
include_modules{ $minimal_modules: }
```

And finally a file containing the name of a server in which I'm going to ask it to load certain modules, but also exclude some:

```ruby
node 'srv.deimos.fr' {
    $exclude_modules = [ 'resolvconf' ]
    include base::minimal
}
```

Here I use an array '$exclude_modules' (with a single element, but you can add several separated by commas), which will allow me to specify which modules to exclude. Because by the next line it will load everything it will need through the include_modules function.

### Templates

When you write manifests, you call a directive named 'File' when you want to send a file to a server. But if the content of that file needs to change based on certain parameters (name, ip, timezone, domain...), then you need to use templates! And that's where it gets interesting as it's possible to script within a template to generate its content. Templates use a language very close to ruby.

In a template, the following syntax is used:

- For facts, it's simple, you need to prefix the variable with a "@". For example <%= fqdn %> becomes <%= @fqdn %>.
- For your variables, if it's declared in the manifest that calls the template, also use "@".
- My variable myvar defined in the manifest that calls this template has the value <%= myvar %>.

If you want to access a variable defined outside the current manifest, outside the local scope, use the scope.lookupvar function:

```ruby
<%= scope.lookupvar('common::config::myvar') %>
```

You can validate your template via:

```bash
erb -P -x -T '-' mytemplate.erb | ruby -c
```

Here's an example with OpenSSH so you understand. I've taken the configuration that will vary according to certain parameters:

```
# Package generated configuration file
# See the sshd(8) manpage for details

# What ports, IPs and protocols we listen for

<% ssh_default_port.each do |val| -%>
Port <%= val -%>
<% end -%>

# Use these options to restrict which interfaces/protocols sshd will bind to

#ListenAddress ::
#ListenAddress 0.0.0.0
Protocol 2

# HostKeys for protocol version 2

HostKey /etc/ssh/ssh_host_rsa_key
HostKey /etc/ssh/ssh_host_dsa_key
#Privilege Separation is turned on for security
UsePrivilegeSeparation yes

# Lifetime and size of ephemeral version 1 server key

KeyRegenerationInterval 3600
ServerKeyBits 768

# Logging

SyslogFacility AUTH
LogLevel INFO

# Authentication:

LoginGraceTime 120
PermitRootLogin yes
StrictModes yes

RSAAuthentication yes
PubkeyAuthentication yes
#AuthorizedKeysFile %h/.ssh/authorized_keys

# Don't read the user's ~/.rhosts and ~/.shosts files

IgnoreRhosts yes

# For this to work you will also need host keys in /etc/ssh_known_hosts

RhostsRSAAuthentication no

# similar for protocol version 2

HostbasedAuthentication no

# Uncomment if you don't trust ~/.ssh/known_hosts for RhostsRSAAuthentication

#IgnoreUserKnownHosts yes

# To enable empty passwords, change to yes (NOT RECOMMENDED)

PermitEmptyPasswords no

# Change to yes to enable challenge-response passwords (beware issues with

# some PAM modules and threads)

ChallengeResponseAuthentication no

# Change to no to disable tunnelled clear text passwords

#PasswordAuthentication yes

# Kerberos options

#KerberosAuthentication no
#KerberosGetAFSToken no
#KerberosOrLocalPasswd yes
#KerberosTicketCleanup yes

# GSSAPI options

#GSSAPIAuthentication no
#GSSAPICleanupCredentials yes

X11Forwarding yes
X11DisplayOffset 10
PrintMotd no
PrintLastLog yes
TCPKeepAlive yes
#UseLogin no

#MaxStartups 10:30:60
#Banner /etc/issue.net

# Allow client to pass locale environment variables

AcceptEnv LANG LC_*

Subsystem sftp /usr/lib/openssh/sftp-server

UsePAM yes
# AllowUsers <%= ssh_allowed_users %>
```

Here we're using two types of template usage. A multi-line repetition, and the other with a simple variable replacement:

- ssh_default_port.each do: allows us to put a line of "Port num_port" for each specified port
- ssh_allowed_users: allows us to give a list of users

These variables are usually declared either in the [node part](#servers.pp) or in the [global configuration](#vars.pp). We've just seen how to put a variable or a loop in a template, but know that it's also possible to use if statements! In short, a complete language exists and allows you to modulate a file as you wish.

These methods prove simple and very effective. Small subtlety:

- -%>: When a line ends like this, there won't be a line break thanks to the - at the end.
- %>: There will be a line break here.

#### Inline-templates

This is a small subtlety that may seem unnecessary, but is actually very useful for executing small methods within a manifest! Take for example the 'split' function that exists in puppet today, it would seem normal that the 'join' function exists, right? Well, no... at least not in the current version at the time of writing this (2.7.18). So I can use in the same way as templates code in my manifests, see for yourself:

```ruby
$ldap_servers = [ '192.168.0.1', '192.168.0.2', '127.0.0.1' ]
$comma_ldap_servers = inline_template("<%= (ldap_servers).join(',') %>")
```

- $ldap_servers: this is a simple array with my list of LDAP servers
- $comma_ldap_servers: we use the inline_template function, which will call the join function, pass it the ldap_servers array and join the content with commas.

I would finally have:

```ruby
$comma_ldap_servers = '192.168.0.1,192.168.0.2,127.0.0.1'
```

### Facters

"Facts" are scripts (see /usr/lib/ruby/1.8/facter for standard facts) that allow building dynamic variables, which change depending on the environment in which they are executed.
For example, we could define a "fact" that determines if we are on a "cluster" type machine based on the presence or absence of a file:

```ruby
 # is_cluster.rb

 Facter.add("is_cluster") do
   setcode do
      FileTest.exists?("/etc/cluster/nodeid")
   end
 end
```

You can also use functions that allow you to directly use facter-type functions in templates (downcase or upcase to change case):

```ruby
#
# Config file for collectd(1).
# Please read collectd.conf(5) for a list of options.
# http://collectd.org/
#

Hostname     <%= hostname.downcase %>
FQDNLookup   true
```

Be careful, if you want to test the fact on the destination machine, don't forget to specify the path where the facts are located on the machine:

```bash
export FACTERLIB=/var/lib/puppet/lib/facter
```

or for Solaris:

```bash
export FACTERLIB=/var/opt/csw/puppet/lib/facter
```

To see the list of facts currently on the system, simply type the facter command:

```bash
> facter
facterversion => 1.5.7
hardwareisa => i386
hardwaremodel => i86pc
hostname => PA-OFC-SRV-UAT-2
hostnameldap => PA-OFC-SRV-UAT
id => root
interfaces => lo0,e1000g0,e1000g0_1,e1000g0_2,e1000g1,e1000g2,e1000g3,clprivnet0
...
```

See http://docs.puppetlabs.com/guides/custom_facts.html for more details.

### Dynamic Information

It's possible to use server-side scripts and retrieve their content in a variable. Here's an example:

```ruby
#!/usr/bin/ruby
require 'open-uri'
page = open("http://www.puppetlabs.com/misc/download-options/").read
print page.match(/stable version is ([\d\.]*)/)[1]
```

And in the manifest:

```ruby
$latestversion = generate("/usr/bin/latest_puppet_version.rb")
notify { "The latest stable Puppet version is ${latestversion}. You're using ${puppetversion}.": }
```

Magical, isn't it? :-). Know that it's even possible to pass arguments with a comma between each one!!!

### Parsers

Parsers are the creation of special functions usable in manifests (server-side). For example, I created a parser that will allow me to do a reverse DNS lookup:

```ruby
# Dns2IP for Puppet
# Made by Pierre Mavro
# Does a DNS lookup and returns an array of strings of the results
# Usage : need to send one string dns servers separated by comma. The return will be the same

require 'resolv'

module Puppet::Parser::Functions
  newfunction(:dns2ip, :type => :rvalue) do |arguments|
    result = [ ]
    # Split comma sperated list in array
    dns_array = arguments[0].split(',')
    # Push each DNS/IP address in result array
    dns_array.each do |dns_name|
      result.push(Resolv.new.getaddresses(dns_name))
    end
    # Join array with comma
    dns_list = result.join(',')
    # Delete last comma if exist
    good_dns_list = dns_list.gsub(/,$/, '')
    return good_dns_list
  end
end
```

We'll be able to create this variable and then insert it into our manifests:

```ruby
$comma_ldap_servers = 'ldap1.deimos.fr,ldap2.deimos.fr,127.0.0.1'
$ip_ldap_servers = dns2ip("${comma_ldap_servers}")
```

Here I send a list of LDAP servers and their IP addresses will be returned to me. Now you understand that it's a call, a bit like [inline_templates](#Les_inline_templates), but much more powerful.

{{< alert context="info" text="I've noticed a rather annoying cache behavior with this type of function! Indeed, when you develop a parser and test it, you are likely to use 'Notify' functions in your manifests for debugging. However, the changes you make to your parser won't necessarily apply until you've cleared the caches. After some research and IRC queries, it turns out that the only working method is to **restart the Puppet Master and the web server (Nginx in our case)**. It works very well, but it's a bit annoying during the debug phase." />}}

### Ruby in your manifests

It's entirely possible to write Ruby in your manifests. See for yourself:

```ruby
notice( "I am running on node %s" % scope.lookupvar("fqdn") )
```

This looks a lot like sprintf.

### Adding a Ruby variable in manifests

If we want to retrieve the current time in a manifest, for example:

```ruby
require 'time'
scope.setvar("now", Time.now)
notice( "Here is the current time : %s" % scope.lookupvar("now") )
```

### Classes

You can use classes with arguments like this:

```ruby
class mysql( $package, $socket, $port = "3306" ) {
    …
}
class { "mysql":
    package => "percona-sql-server-5.0",
    socket  => "/var/run/mysqld/mysqld.sock",
    port    => "3306",
}
```

### Using hash tables

Just like arrays, it's also possible to use hash tables, look at this example:

```ruby
$interface = {
    name => 'eth0',
    address => '192.168.0.1'
}
notice("Interface ${interface[name]} has address ${interface[address]}")
```

### Regex

It's possible to use regex and retrieve patterns:

```ruby
$input = "What a great tool"
if $input =~ /What a (\w+) tool/ {
   notice("You said the tool is : '$1'. The complete line is : $0")
}
```

### Substitution

Substitution is possible:

```ruby
$ipaddress = '192.168.0.15'
$class_c = regsubst($ipaddress, "(.*)\\..*", "\\1.0")
notify { $ipaddress: }
notify { $class_c: }
```

This will give me 192.168.0.15 and 192.168.0.0.

### Notify and Require

These two functions are very useful once inserted into a manifest. This allows, for example, a service to say that it requires (require) a Package to function and a configuration file to notify (notify) a service if it changes so that it restarts the daemon. You can also write something like this:

```ruby
Package["ntp"] -> File["/etc/ntp.conf"] ~> Service["ntp"]
```

- ->: means 'require'
- ~>: means 'notify'

It's also possible to do requires on classes :-)

### The +> operator

Here's a great operator that will save us some time. The example below:

```ruby
file { "/etc/ssl/certs/cookbook.pem":
    source => "puppet:///modules/apache/deimos.pem",
}
Service["apache2"] {
    require +> File["/etc/ssl/certs/deimos.pem"],
}
```

Corresponds to:

```ruby
service { "apache2":
       enable  => true,
       ensure  => running,
       require => File["/etc/ssl/certs/deimos.pem"],
}
```

### Checking software version number

If you need to check the version number of a software to make a decision, here's a good example:

```ruby
$app_version = "2.7.16"
$min_version = "2.7.18"
if versioncmp( $app_version, $min_version ) >= 0 {
   notify { "Puppet version OK": }
} else {
   notify { "Puppet upgrade needed": }
}
```

### Virtual resources

Useful for test writings, you can, for example, create a resource by preceding it with an '@'. It will be read but not executed until you explicitly tell it to (realize). Example:

```ruby
@package {
    'postfix':
        ensure => installed
}
realize( Package[''postfix] )
```

One of the big advantages of this method is that you can declare the realize in several places in your puppet master without having conflicts!

### Advanced file deletion

You can request the deletion of a file after a given time or from a certain size:

```ruby
tidy { "/var/lib/puppet/reports":
    age     => "1w",
    size    => "512k",
    recurse => true,
}
```

This will cause the deletion of a folder after a week with its content.

## Modules

It's recommended to create modules for each service to make the configuration more flexible. This is part of certain best practices.

I'll cover different techniques here trying to keep an increasing order of difficulty.

### Initializing a module

So we'll create the appropriate directory structure on the server. For this example, we'll start with "sudo", but you can choose something else if you want:

```bash
mkdir -p /etc/puppet/modules/sudo/manifests
touch /etc/puppet/modules/sudo/manifests/init.pp
```

Note that this is necessary for each module. The init.pp file is the first file that will load when the module is called.

### The initial module (base)

We need to create an initial module that will manage the list of servers, the functions we'll need, the roles, global variables... in short, it may seem a bit abstract at first but just know that we need a module to then manage all the others. We'll start with this one which is one of the most important for the future.

As you now know, we need an init.pp file for the first module to be loaded. So we'll create our directory structure which we'll call "base":

```bash
mkdir -p /etc/puppet/modules/base/{manifests,puppet/parser/functions}
```

#### init.pp

Then we'll create and fill the init.pp file:

```ruby
################################################################################
#                                BASE MODULES                                  #
################################################################################

# Load defaults vars
import "vars.pp"
# Load functions
import "functions.pp"
# Load sysctl module
include "sysctl"
# Load network module
include "network"
# Load roles
import "roles.pp"
# Set servers properties
import "servers.pp"
```

The lines corresponding to import are equivalent to an "include" (in services like ssh or nrpe) of my other .pp files that we'll create later. While the includes will load other modules that I'll create later.

#### vars.pp

We'll then create the vars.pp file which will contain all my global variables for my future modules or manifests (\*.pp):

```ruby
################################################################################
#                                    VARS                                      #
################################################################################

# Default admins emails
$root_email = 'xxx@mycompany.com'

# NTP Timezone. Usage :
# Look at /usr/share/zoneinfo/ and add the continent folder followed by the town
$set_timezone = 'Europe/Paris'

# Define empty exclude modules
$exclude_modules = [ ]

# Default LDAP servers
$ldap_servers = [ ]

# Default DNS servers
$dns_servers = [ '192.168.0.69', '192.168.0.27' ]
```

#### functions.pp

Now, we'll create functions that will allow us to add some features not currently present in puppet or simplify some:

```ruby
/*
Puppet Functions
Made by Pierre Mavro
*/
################################################################################
#                             GLOBAL FUNCTIONS                                 #
################################################################################

# Load or not modules (include/exclude)
define include_modules () {
    if ($exclude_modules == undef) or !($name in $exclude_modules) {
        include $name
    }
}

# Validate that pacakges are installed
define packages_install () {
    notice("Installation of ${name} package")
    package {
        "${name}":
            ensure => present
    }
}

# Check that those services are enabled on boot or not
define services_start_on_boot ($enable_status) {
    service {
        "${name}":
            enable => "${enable_status}"
    }
}

# Add, remove, comment or uncomment lines
define line ($file, $line, $ensure = 'present') {
    case $ensure {
        default : {
            err("unknown ensure value ${ensure}")
        }
        present : {
            exec {
                "echo '${line}' >> '${file}'" :
                    unless => "grep -qFx '${line}' '${file}'",
                    logoutput => true
            }
        }
        absent : {
            exec {
                "grep -vFx '${line}' '${file}' | tee '${file}' > /dev/null 2>&1" :
                    onlyif => "grep -qFx '${line}' '${file}'",
                    logoutput => true
            }
        }
        uncomment : {
            exec {
                "sed -i -e'/${line}/s/#\+//' '${file}'" :
                    onlyif => "test `grep '${line}' '${file}' | grep '^#' | wc -l` -ne 0",
                    logoutput => true
            }
        }
        comment : {
            exec {
                "/bin/sed -i -e'/${line}/s/\(.\+\)$/#\1/' '${file}'" :
                    onlyif => "test `grep '${line}' '${file}' | grep -v '^#' | wc -l` -ne 0",
                    logoutput => true
            }
        }

        # Use this resource instead if your platform's grep doesn't support -vFx;
        # note that this command has been known to have problems with lines containing quotes.
        # exec { "/usr/bin/perl -ni -e 'print unless /^\Q${line}\E$/' '${file}'":
        #     onlyif => "grep -qFx '${line}' '${file}'"
        # }

    }
}

# Validate that softwares are installed
define comment_lines ($filename) {
    line {
        "${name}" :
            file => "${filename}",
            line => "${name}",
            ensure => comment
    }
}

# Sysctl managment
class sysctl {
    define conf ($value) {
    # $name is provided by define invocation
    # guid of this entry
        $key = $name
        $context = "/files/etc/sysctl.conf"
        augeas {
            "sysctl_conf/$key" :
                context => "$context",
                onlyif => "get $key != '$value'",
                changes => "set $key '$value'",
                notify => Exec["sysctl"],
        }
    }
    file {
        "sysctl_conf" :
            name => $::operatingsystem ? {
                default => "/etc/sysctl.conf",
            },
    }
    exec {
        "sysctl -p" :
            alias => "sysctl",
            refreshonly => true,
            subscribe => File["sysctl_conf"],
    }
}

# Function to add ssh public keys
define ssh_add_key ($user, $key) {
   # Create users home directory if absent
   exec {
      "mkhomedir_${name}" :
         path => "/bin:/usr/bin",
         command => "cp -Rfp /etc/skel ~$user; chown -Rf $user:group ~$user",
         onlyif => "test `ls ~$user 2>&1 >/dev/null | wc -l` -ne 0"
   }

   ssh_authorized_key {
      "${name}" :
         ensure => present,
         key => "$key",
         type => 'ssh-rsa',
         user => "$user",
         require => Exec["mkhomedir_${name}"]
   }
}

# Limits.conf managment
define limits_conf ($domain = "root", $type = "soft",$item = "nofile",	$value = "10000") {
    # guid of this entry
    $key = "$domain/$type/$item"

    # augtool> match /files/etc/security/limits.conf/domain[.="root"][./type="hard" and ./item="nofile" and ./value="10000"]
    $context = "/files/etc/security/limits.conf"
    $path_list = "domain[.=\"$domain\"][./type=\"$type\" and ./item=\"$item\"]"
    $path_exact = "domain[.=\"$domain\"][./type=\"$type\" and ./item=\"$item\" and ./value=\"$value\"]"

    augeas {
        "limits_conf/$key" :
            context => "$context",
            onlyif => "match $path_exact size==0",
            changes => [
                # remove all matching to the $domain, $type, $item, for any $value
                "rm $path_list",
                # insert new node at the end of tree
                "set domain[last()+1] $domain",
                # assign values to the new node
                "set domain[last()]/type $type",
                "set domain[last()]/item $item",
                "set domain[last()]/value $value",],
    }
}
```

So we have:

- Line 10: The ability to load or not load modules via an array sent as a function argument (as described earlier in this documentation)
- Line 17: The ability to verify that packages are installed on the machine
- Line 26: The ability to verify that services are correctly loaded at machine boot
- Line 34: The ability to ensure that a line in a file is present, absent, commented, or not commented
- Line 78: The ability to comment multiple lines via an array sent as a function argument
- Line 88: The ability to manage the sysctl.conf file
- Line 117: The ability to easily deploy SSH public keys
- Line 128: The ability to simply manage the limits.conf file

All these functions are of course not mandatory but greatly help with the use of puppet.

#### roles.pp

Then we have a file containing the roles of the servers. See it as groups to which we'll subscribe the servers:

```ruby
################################################################################
#                                   ROLES                                      #
################################################################################

# Level 1 : Minimal
class base::minimal
{
    # Load modules
    $minimal_modules =
    [
        'stdlib',
        'puppet',
        'resolvconf',
        'packages_defaults',
        'configurations_defaults',
        'openssh',
        'selinux',
        'grub',
        'kdump',
        'tools',
        'timezone',
        'ntp',
        'mysecureshell',
        'openldap',
        'acl',
        'sudo',
        'snmpd',
        'postfix',
        'nrpe'
    ]
    include_modules{ $minimal_modules: }
}

# Level 2 : Cluster
class base::cluster inherits minimal
{
    # Load modules
    $cluster_modules =
    [
        'packages_cluster',
        'configurations_cluster'
    ]
    include_modules{ $cluster_modules: }
}

# Level 2 : Low Latency
class base::low_latency inherits minimal
{
    # Load modules
    $lowlatency_modules =
    [
        'low_latency'
    ]
    include_modules{ $lowlatency_modules: }
}

# Level 3 : Low Latency + Cluster
class base::low_latency_cluster inherits minimal
{
    include base::cluster
    include base::low_latency
}
```

I've defined classes here that inherit from each other to varying degrees. It's actually defined by levels. Level 3 depends on 2 and 1, 2 depends on 1, and 1 has no dependencies. This gives me a certain flexibility. For example, I know that if I load my cluster class, my minimal class will also be loaded. You'll notice the 'base::minimal' annotation. It's recommended to load your classes by calling the module, followed by '::'. This makes it much easier to read the manifests.

#### servers.pp

And finally, I have a file where I make my server declaration:

```ruby
/*##############################################################################
#                                     SERVERS                                  #
################################################################################

== Automated Dependancies Roles ==
* cluster -> minimal
* low_latency -> minimal
* low_latency_cluster -> low_latency + cluster + minimal

== Template for servers ==
node /regex/
{
    #$exclude_modules = [ ]
    #$ldap_servers = 'x.x.x.x'
    #$set_timezone = 'Europe/Paris'
    #$dns_servers = [ ]
    #include base::minimal
    #include base::cluster
    #include base::low_latency
    #include base::low_latency_cluster
}

##############################################################################*/

# One server
node 'srv1.deimos.fr'  {
    $ldap_servers = [ '127.0.0.1' ]
    include base::minimal
}

# Multiple servers
node 'srv2.deimos.fr' 'srv3.deimos.fr'  {
    $ldap_servers = [ '127.0.0.1' ]
    include base::minimal
}

# Multiple regex based servers
node /srv-prd-\d+/ {
    include base::minimal
    include base::low_latency
    $set_timezone = 'Europe/London'
}
```

Here I've put a server as an example or a regex for multiple servers. For info, the configuration can be integrated into [LDAP](https://reductivelabs.com/trac/puppet/wiki/LDAPNodes).

#### Parser

Let's create the necessary directory structure:

```bash
mkdir -p /etc/puppet/modules/base/puppet/parser/functions
```

Then add an empty parser that will allow us to detect if an array/variable is empty or not:

```ruby
#
# empty.rb
#
# Copyright 2011 Puppet Labs Inc.
# Copyright 2011 Krzysztof Wilczynski
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

module Puppet::Parser::Functions
  newfunction(:empty, :type => :rvalue, :doc => <<-EOS
Returns true if given array type or hash type has no elements or when a string
value is empty and false otherwise.

Prototype:

empty(x)

Where x is either an array, a hash or a string value.

For example:

Given the following statements:

$a = ''
$b = 'abc'
$c = []
$d = ['d', 'e', 'f']
$e = {}
$f = { 'x' => 1, 'y' => 2, 'z' => 3 }

notice empty($a)
notice empty($b)
notice empty($c)
notice empty($d)
notice empty($e)
notice empty($f)

The result will be as follows:

notice: Scope(Class[main]): true
notice: Scope(Class[main]): false
notice: Scope(Class[main]): true
notice: Scope(Class[main]): false
notice: Scope(Class[main]): true
notice: Scope(Class[main]): false
EOS
  ) do |*arguments|
    #
    # This is to ensure that whenever we call this function from within
    # the Puppet manifest or alternatively form a template it will always
    # do the right thing ...
    #
    arguments = arguments.shift if arguments.first.is_a?(Array)

    raise Puppet::ParseError, "empty(): Wrong number of arguments " +
      "given (#{arguments.size} for 1)" if arguments.size < 1

    value = arguments.shift

    unless [Array, Hash, String].include?(value.class)
      raise Puppet::ParseError, 'empty(): Requires either array, hash ' +
        'or string type to work with'
    end

    value.empty?
  end
end

# vim: set ts=2 sw=2 et :
# encoding: utf-8
```

## Module Examples

### stdlib

This [stdlib](https://forge.puppetlabs.com/puppetlabs/stdlib) module is not essential, but it's useful if you're missing features in Puppet. Indeed, it brings a fairly large set of functions:

```
abs		     ensure_resource	  include	       loadyaml		    reverse		 to_bytes
bool2num	     err		  info		       lstrip		    rstrip		 type
capitalize	     extlookup		  inline_template      md5		    search		 unique
chomp		     fail		  is_array	       member		    sha1		 upcase
chop		     file		  is_domain_name       merge		    shellquote		 validate_absolute_pa
create_resources     flatten		  is_float	       notice		    size		 validate_array
crit		     fqdn_rand		  is_hash	       num2bool		    sort		 validate_bool
debug		     fqdn_rotate	  is_integer	       parsejson	    squeeze		 validate_hash
defined		     generate		  is_ip_address	       parseyaml	    str2bool		 validate_re
defined_with_params  get_module_path	  is_mac_address       prefix		    str2saltedsha512	 validate_slength
delete		     getvar		  is_numeric	       range		    strftime		 validate_string
delete_at	     grep		  is_string	       realize		    strip		 values
downcase	     has_key		  join		       regsubst		    swapcase		 values_at
emerg		     hash		  keys		       require		    time		 zip
empty
```

First create the directory structure:

```bash
mkdir -p /etc/puppet/modules/stdlib
```

Download the latest version and simply decompress it:

```bash
cd /etc/puppet/modules/stdlib
wget http://forge.puppetlabs.com/puppetlabs/stdlib/3.2.0.tar.gz
tar -xzf 3.2.0.tar.gz
mv puppetlabs-stdlib-3.2.0/stdlib stdlib
rm -f 3.2.0.tar.gz
```

### Puppet

This one is quite funny because it's simply the configuration of the Puppet client. However, it can be very useful for managing its own updates. So let's create the directory structures:

```bash
mkdir -p /etc/puppet/modules/puppet/{manifests,files}
```

#### init.pp

We create the init.pp module here that will allow us to choose the file to load according to the OS.

```ruby
/*
Puppet Module for Puppet
Made by Pierre Mavro
*/
class puppet {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include ::puppet::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

```ruby
/*
Puppet Module for Puppet
Made by Pierre Mavro
*/
class puppet::redhat {
    # Change default configuration
    file {
        '/etc/puppet/puppet.conf' :
            ensure => present,
            source => "puppet:///modules/puppet/${::osfamily}.puppet.conf",
            mode => 644,
            owner => root,
            group => root
    }

    # Disable service on boot and be sure it is not started
    service {
        'puppet-srv' :
            name => 'puppet',
            # Let this line commented if you're using Puppet Dashboard
            #ensure => stopped,
            enable => false
    }
}
```

On line 9, we use a variable available in the facts (client-side) so that depending on the response, we load a file associated with the OS. So we'll have a configuration file accessible via Puppet in the form 'RedHat.puppet.conf'.
Then, for the service, we make sure it's properly stopped at startup and that it's in an off state for now. In fact, I don't want it to trigger and synchronize every 30 minutes (default value), I find it too dangerous and prefer to decide via other mechanisms (SSH, Mcollective...) when I want a synchronization to be done.

#### files

In files, we'll have the basic configuration file that should apply to all RedHat type machines:

```ini
[main]
    # The Puppet log directory.
    # The default value is '$vardir/log'.
    logdir = /var/log/puppet

    # Where Puppet PID files are kept.
    # The default value is '$vardir/run'.
    rundir = /var/run/puppet

    # Where SSL certificates are kept.
    # The default value is '$confdir/ssl'.
    ssldir = $vardir/ssl

    # Puppet master server
    server = puppet-prd.deimos.fr

    # Add custom facts
    pluginsync = true
    pluginsource = puppet://$server/plugins
    factpath = /var/lib/puppet/lib/facter

[agent]
    # The file in which puppetd stores a list of the classes
    # associated with the retrieved configuratiion.  Can be loaded in
    # the separate ``puppet`` executable using the ``--loadclasses``
    # option.
    # The default value is '$confdir/classes.txt'.
    classfile = $vardir/classes.txt

    # Where puppetd caches the local configuration.  An
    # extension indicating the cache format is added automatically.
    # The default value is '$confdir/localconfig'.
    localconfig = $vardir/localconfig

    # Reporting
    report = true

    # Inspect reports for a compliance workflow
    archive_files = true
```

### resolvconf

I made this module to manage the resolv.conf configuration file. The usage is quite simple, it will retrieve the information of the DNS servers filled in the array available in [vars.pp of the base module](#vars.pp). So fill in the default DNS servers:

```ruby
# Default DNS servers
$dns_servers = [ '192.168.0.69', '192.168.0.27' ]
```

You can override these values directly at the level of one or more nodes if you need to have specific configurations for certain nodes (in the [servers.pp file of the base module](#servers.pp)):

```ruby
# One server
node 'srv.deimos.fr' {
    $dns_servers = [ '127.0.0.1' ]
    include base::minimal
}
```

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/resolvconf/{manifests,templates}
```

#### init.pp

```ruby
/*
Resolv.conf Module for Puppet
Made by Pierre Mavro
*/
class resolvconf {
    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include resolvconf::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

Here's the configuration for Red Hat, I use a template file here, which will be filled with the information present in the $dns_servers array:

```ruby
/*
Resolvconf Module for Puppet
Made by Pierre Mavro
*/
class resolvconf::redhat {
    # resolv.conf file
    file {
        "/etc/resolv.conf" :
            content => template("resolvconf/resolv.conf"),
            mode => 744,
            owner => root,
            group => root
    }
}
```

#### templates

And finally my resolv.conf template file:

```ruby
# Generated by Puppet
domain deimos.fr
search deimos.fr deimos.lan
<% dns_servers.each do |dnsval| -%>
nameserver <%= dnsval %>
<% end -%>
```

Here we have a ruby loop that will go through the $dns_servers array and build the resolv.conf file by inserting line by line 'nameserver' with the associated server.

## Puppet: Configuration File Management Solution

### packages_defaults

I use this module to install or uninstall packages that I absolutely need on all my machines. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/packages_defaults/manifests
```

#### init.pp

```ruby
class packages_defaults {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include ::packages_defaults::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

I could have grouped everything into a single block, but for the sake of readability on the packages included in the distribution and those I added in a custom repository, I preferred to make a separation:

```ruby
# Red Hat Defaults packages
class packages_defaults::redhat
{
    # Set all default packages (embended in distribution) that need to be installed
    packages_install
    { [
        'nc',
        'tree',
        'telnet',
        'dialog',
        'freeipmi',
        'glibc-2.12-1.80.el6.i686'
    ]: }

    # Set all custom packages (not embended in distribution) that need to be installed
    packages_install
    { [
        'puppet',
        'tmux'
    ]: }
}
```

### configurations_defaults

This module, like the previous one, is used for the configuration of the OS delivered as standard. I actually want to make adjustments to parts of the pure system here, without really getting into a particular software. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/configuration_defaults/{manifests,templates,lib/facter}
```

#### init.pp

```ruby
class configurations_defaults {
    import '*.pp'

    # Configure common security parameters
    include configurations_defaults::common

    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include configurations_defaults::redhat
        }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

Here I load all .pp files at startup, then call and import common configurations (common), then apply configurations specific to each OS.

#### common.pp

Here I want to have the same base motd file for all my machines. You'll see later why it appears as a template:

```ruby
class configurations_defaults::common
{
    # Motd banner for all servers
    file {
        '/etc/motd':
            ensure => present,
            content => template("configurations_defaults/motd"),
            mode => 644,
            owner => root,
            group => root
    }
}
```

#### redhat.pp

I will load security options here, automatically configure bonding on my machines, and a sysctl option:

```ruby
class configurations_defaults::redhat
{
    # Security configurations
    include 'configurations_defaults::redhat::security'

    # Configure bonding
    include 'configurations_defaults::redhat::network'

    # Set sysctl options
    sysctl::conf
    {
        'vm.swappiness': value => '0';
    }
}
```

#### security.pp

You'll see that this file does quite a lot:

```ruby
class configurations_defaults::redhat::security inherits configurations_defaults::redhat
{
    # Manage Root passwords    $sha512_passwd='$6$lhkAz...'
    $md5_passwd='$1$Fcwy...'

    if ($::passwd_algorithm == sha512)
    {
        # sha512 root password
        $root_password="$sha512_passwd"
    }
    else
    {
        # MD5 root password
        $root_password="$md5_passwd"
    }
    user {
        'root':
            ensure   => present,
            password => "$root_password"
    }

    # Enable auditd service    service {
        "auditd" :
            enable => true,
            ensure => 'running',
    }

    # Comment unwanted sysctl lines    $sysctl_file = '/etc/sysctl.conf'
    $sysctl_comment_lines =
    [
        "net.bridge.bridge-nf-call-ip6tables",
        "net.bridge.bridge-nf-call-iptables",
        "net.bridge.bridge-nf-call-arptables"
    ]
    comment_lines {
        $sysctl_comment_lines :
            filename => "$sysctl_file"
    }

    # Add security sysctl values    sysctl::conf
    {
        'vm.mmap_min_addr': value => '65536';
        'kernel.modprobe': value => '/bin/false';
        'kernel.kptr_restrict': value => '1';
        'net.ipv6.conf.all.disable_ipv6': value => '1';
    }

    # Deny kernel read to others users    case $::kernel_security_rights {
        '0': {
            exec {'chmod_kernel':
                command => 'chmod o-r /boot/{vmlinuz,System.map}-*'
            }
        }
        '1' : {
            notice("Kernel files have security rights")
        }
    }

    # Change opened file descriptor value and avoid fork bomb by limiting number of process    limits_conf {
        "open_fd": domain => '*', type => '-', item => nofile, value => 2048;
        "fork_bomb_soft": domain => '@users', type => soft, item => nproc, value => 200;
        "fork_bomb_hard": domain => '@users', type => hard, item => nproc, value => 300;
    }
}
```

Some explanations are needed:

- Manage Root passwords: We define the desired root password in md5 and sha1 form. Depending on what's configured on the machine, it will configure the desired password. For this detection I use a facter (passwd_algorithm.rb)
- Enable auditd service: we make sure that the auditd service will start at boot and is currently running
- Comment unwanted sysctl lines: we ask that certain lines present in sysctl be commented if they exist
- Add security sysctl values: we add sysctl rules, and assign them a value
- Deny kernel read to others users: I created a facter here, which checks the rights of kernel files (kernel_rights.rb)
- Change opened file descriptor value: allows to use the limits_conf function to manage the limits.conf file. Here I've changed the default value of file descriptors and added a small security to avoid fork bombs.

#### facter

We'll insert here the facters that will be used for certain functions requested above.

##### passwd_algorithm.rb

This facter will determine the algorithm used for authentication:

```ruby
# Get Passwd Algorithm
Facter.add("passwd_algorithm") do
  setcode do
    Facter::Util::Resolution.exec("grep ^PASSWDALGORITHM /etc/sysconfig/authconfig | awk -F'=' '{ print $2 }'")
  end
end
```

##### kernel_rights.rb

This facter will determine if any user has the right to read the kernels installed on the current machine:

```ruby
# Get security rights
Facter.add(:kernel_security_rights) do
  # Get kernel files where rights will be checked
  kernel_files = Dir.glob("/boot/{vmlinuz,System.map}-*")
  current_rights=1

  # Check each files
  kernel_files.each do |file|
    # Get file mode
    full_rights = sprintf("%o", File.stat(file).mode)
    # Get last number (correponding to other rights)
    other_rights = Integer(full_rights) % 10
    # Check if other got read rights
    if other_rights >= 4
      current_rights=0
    end
  end
  setcode do
    # Set kernel_security_rights to 1 if read value is detected
    current_rights
  end
end
```

##### get_network_infos.rb

This facter allows to retrieve the current ip on eth0, the netmask and the gateway:

```ruby
# Get public IP address
Facter.add(:public_ip) do
  setcode do
    # Get bond0 ip if exist
    if File.exist? "/proc/sys/net/ipv4/conf/bond0"
      Facter::Util::Resolution.exec("ip addr show dev bond0 | awk '/inet/{print $2}' | head -1 | sed 's/\\/.*//'")
    else
      # Or eth0 ip if exist
      if File.exist? "/proc/sys/net/ipv4/conf/eth0"
        Facter::Util::Resolution.exec("ip addr show dev eth0 | awk '/inet/{print $2}' | head -1 | sed 's/\\/.*//'")
      else
        # Else return error
        'unknow (puppet issue)'
      end
    end
  end
end

# Get netmask on the fist interface
Facter.add(:public_netmask) do
  setcode do
    # Get bond0 netmask if exist
    if File.exist? "/proc/sys/net/ipv4/conf/bond0"
      Facter::Util::Resolution.exec("ifconfig bond0 | awk '/inet/{print $4}' | sed 's/.*://'")
    else
      # Or eth0 netmask if exist
      if File.exist? "/proc/sys/net/ipv4/conf/eth0"
        Facter::Util::Resolution.exec("ifconfig eth0 | awk '/inet/{print $4}' | sed 's/.*://'")
      else
        # Else set a default netmask
        '255.255.255.0'
      end
    end
  end
end

# Get default gateway
Facter.add(:default_gateway) do
  setcode do
    Facter::Util::Resolution.exec("ip route | awk '/default/{print $3}'")
  end
end
```

#### network.pp

Here we will load the bonding configuration and other network-related things:

```ruby
class configurations_defaults::redhat::network inherits configurations_defaults::redhat
{
    # Disable network interface renaming    augeas {
        "grub_udev_net" :
            context => "/files/etc/grub.conf",
            changes => "set title[1]/kernel/biosdevname 0"
    }

    # Load bonding module at boot    line {
        'load_bonding':
            file => '/etc/modprobe.d/bonding.conf',
            line => 'alias bond0 bonding',
            ensure => present
    }

    # Bonded master interface - static    network::bond::static {
        "bond0" :
            ipaddress => "$::public_ip",
            netmask => "$::public_netmask",
            gateway => "$::default_gateway",
            bonding_opts => "mode=active-backup",
            ensure => "up"
    }

    # Bonded slave interface - static
    network::bond::slave {
        "eth0" :
            macaddress => $::macaddress_eth0,
            master => "bond0",
    }

    # Bonded slave interface - static
    network::bond::slave {
        "eth1" :
            macaddress => $::macaddress_eth1,
            master => "bond0",
    }
}
```

- Disable network interface renaming: we add an argument and set its value to 0 in grub so that it doesn't rename the interfaces and leaves them as ethX. I wrote an article about this if you're interested.
- Load bonding module at boot: We make sure that the bonding module will be loaded at boot time and that an alias on bond0 exists
- Bonded interfaces: I refer you to the bonding module available on Puppet Forge, as well as the documentation on bonding if you don't know what it is. I also created a facter (get_network_infos.rb) for this to retrieve the public interface (eth0, on which I would connect), the netmask and gateway already present and configured on the machine

#### templates

We had talked about it earlier, we manage a template for motd to display, in addition to a text, the hostname of the machine you connect to (line 16):

```
================================================================================
This is an official computer system and is the property of Deimos. It
is for authorized users only. Unauthorized users are prohibited. Users
(authorized or unauthorized) have no explicit or implicit expectation of
privacy. Any or all uses of this system may be subject to one or more of
the following actions: interception, monitoring, recording, auditing,
inspection and disclosing to security personnel and law enforcement
personnel, as well as authorized officials of other agencies, both domestic
and foreign. By using this system, the user consents to these actions.
Unauthorized or improper use of this system may result in administrative
disciplinary action and civil and criminal penalties. By accessing this
system you indicate your awareness of and consent to these terms and
conditions of use. Discontinue access immediately if you do not agree to
the conditions stated in this notice.
================================================================================

                             <%= hostname %>
```

### OpenSSH - 1

Here's a first example for OpenSSH. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/openssh/{manifests,templates,lib/facter}
```

#### init.pp

```ruby
/*
OpenSSH Module for Puppet
Made by Pierre Mavro
*/
class openssh {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include openssh::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

Here we load everything we need, then load the common because OpenSSH needs to be installed and configured before moving on to the common part:

```ruby
/*
OpenSSH Module for Puppet
Made by Pierre Mavro
*/
class openssh::redhat {
    # Install ssh package
    package {
        'openssh-server' :
            ensure => present
    }

    # SSHd config file
    file {
        "/etc/ssh/sshd_config" :
            source => "puppet:///modules/openssh/sshd_config.$::operatingsystem",
            mode => 600,
            owner => root,
            group => root,
            notify  => Service["sshd"]
    }

    service {
        'sshd' :
            enable => true,
            ensure => running,
            require => File['/etc/ssh/sshd_config']
    }

    include openssh::common
}
```

In the service part, there is important information (require), which will restart the service if the configuration file changes.

#### common.pp

Here, we make sure that the folder where SSH keys are stored is present with the right permissions, then we call another file that will contain all the keys we want to export:

```ruby
/*
OpenSSH Module for Puppet
Made by Pierre Mavro
*/
class openssh::common {
    # Check that .ssh directory exist with correct rights
    file {
        "$::home_root/.ssh" :
            ensure => directory,
            mode => 0700,
            owner => root,
            group => root
    }

    # Load all public keys
    include openssh::ssh_keys
}
```

The 'home_root' directive is generated from a facter provided below.

#### facter

Here's the facter that allows you to retrieve the home of the root user:

```ruby
# Get Home directory
Facter.add("home_root") do
  setcode do
    Facter::Util::Resolution.exec("echo ~root")
  end
end
```

#### ssh_keys.pp

Here I add the keys, both for access from other servers or users:

```ruby
/*
OpenSSH Module for Puppet
Made by Pierre Mavro
*/
class openssh::ssh_keys inherits openssh::common {
    ###################################
    # Servers
    ###################################

    # Puppet master
    ssh_add_key {
        'puppet_root' :
            user => 'root',
            key => 'AAAA...'
    }

    ###################################
    # Sys-Admins
    ###################################

    # Pierre Mavro
    ssh_add_key {
        'pmavro_root' :
            user => 'root',
            key => 'AAAA...'
    }
}
```

#### files

The OpenSSH configuration file for Red Hat:

```
# $OpenBSD: sshd_config,v 1.80 2008/07/02 02:24:18 djm Exp $

# This is the sshd server system-wide configuration file.  See
# sshd_config(5) for more information.

# This sshd was compiled with PATH=/usr/local/bin:/bin:/usr/bin

# The strategy used for options in the default sshd_config shipped with
# OpenSSH is to specify options with their default value where
# possible, but leave them commented.  Uncommented options change a
# default value.

Port 22
#AddressFamily any
#ListenAddress 0.0.0.0
#ListenAddress ::

# Disable legacy (protocol version 1) support in the server for new
# installations. In future the default will change to require explicit
# activation of protocol 1
Protocol 2

# HostKey for protocol version 1
#HostKey /etc/ssh/ssh_host_key
# HostKeys for protocol version 2
#HostKey /etc/ssh/ssh_host_rsa_key
#HostKey /etc/ssh/ssh_host_dsa_key

# Lifetime and size of ephemeral version 1 server key
#KeyRegenerationInterval 1h
#ServerKeyBits 1024

# Logging
# obsoletes QuietMode and FascistLogging
#SyslogFacility AUTH
SyslogFacility AUTHPRIV
#LogLevel INFO

# Authentication:

LoginGraceTime 2m
PermitRootLogin without-password
#StrictModes yes
#MaxAuthTries 6
#MaxSessions 10

#RSAAuthentication yes
#PubkeyAuthentication yes
#AuthorizedKeysFile	.ssh/authorized_keys
#AuthorizedKeysCommand none
#AuthorizedKeysCommandRunAs nobody

# For this to work you will also need host keys in /etc/ssh/ssh_known_hosts
#RhostsRSAAuthentication no
# similar for protocol version 2
#HostbasedAuthentication no
# Change to yes if you don't trust ~/.ssh/known_hosts for
# RhostsRSAAuthentication and HostbasedAuthentication
#IgnoreUserKnownHosts no
# Don't read the user's ~/.rhosts and ~/.shosts files
#IgnoreRhosts yes

# To disable tunneled clear text passwords, change to no here!
#PasswordAuthentication yes
#PermitEmptyPasswords no
PasswordAuthentication yes

# Change to no to disable s/key passwords
#ChallengeResponseAuthentication yes
ChallengeResponseAuthentication no

# Kerberos options
#KerberosAuthentication no
#KerberosOrLocalPasswd yes
#KerberosTicketCleanup yes
#KerberosGetAFSToken no
#KerberosUseKuserok yes

# GSSAPI options
# Disable GSSAPI to avoid login slowdown
GSSAPIAuthentication no
#GSSAPIAuthentication yes
#GSSAPICleanupCredentials yes
#GSSAPICleanupCredentials no
#GSSAPIStrictAcceptorCheck yes
#GSSAPIKeyExchange no

# Set this to 'yes' to enable PAM authentication, account processing,
# and session processing. If this is enabled, PAM authentication will
# be allowed through the ChallengeResponseAuthentication and
# PasswordAuthentication.  Depending on your PAM configuration,
# PAM authentication via ChallengeResponseAuthentication may bypass
# the setting of "PermitRootLogin without-password".
# If you just want the PAM account and session checks to run without
# PAM authentication, then enable this but set PasswordAuthentication
# and ChallengeResponseAuthentication to 'no'.
#UsePAM no
UsePAM yes

# Accept locale-related environment variables
AcceptEnv LANG LC_CTYPE LC_NUMERIC LC_TIME LC_COLLATE LC_MONETARY LC_MESSAGES
AcceptEnv LC_PAPER LC_NAME LC_ADDRESS LC_TELEPHONE LC_MEASUREMENT
AcceptEnv LC_IDENTIFICATION LC_ALL LANGUAGE
AcceptEnv XMODIFIERS

# For security reasons, deny tcp forwarding
AllowTcpForwarding no
X11Forwarding no
# Disable DNS usage to avoid login slowdown
UseDNS no
# Disconnect client if they are idle
ClientAliveInterval 600
ClientAliveCountMax 0

#AllowAgentForwarding yes
#GatewayPorts no
#X11Forwarding no
#X11DisplayOffset 10
#X11UseLocalhost yes
#PrintMotd yes
#PrintLastLog yes
#TCPKeepAlive yes
#UseLogin no
#UsePrivilegeSeparation yes
#PermitUserEnvironment no
#Compression delayed
#ShowPatchLevel no
#PidFile /var/run/sshd.pid
#MaxStartups 10
#PermitTunnel no
#ChrootDirectory none

# no default banner path
#Banner none

# override default of no subsystems
Subsystem	sftp	/usr/libexec/openssh/sftp-server

# Example of overriding settings on a per-user basis
#Match User anoncvs
#	X11Forwarding no
#	AllowTcpForwarding no
#	ForceCommand cvs server
```

### OpenSSH - 2

Here's a second example for OpenSSH, slightly different. You need to initialize the module before continuing.

#### init.pp

Here I want my sshd_config file to be of interest:

```ruby
#ssh.pp

# SSH Class with all includes
class ssh {
   $ssh_default_port = ["22"]
   $ssh_allowed_users = "root"
   include ssh::config, ssh::key, ssh::service
}

# SSHD Config file
class ssh::config {
  File {
    name => $operatingsystem ? {
      Solaris  => "/etc/ssh/sshd_config",
      default => "/etc/ssh/sshd_config"
    },
  }

  # Using templates for sshd_config
  file { sshd_config:
    content => $operatingsystem ? {
       default => template("ssh/sshd_config"),
       Solaris => template("ssh/sshd_config.solaris"),
    }
  }
}

# SSH Key exchange
class ssh::key {
	$basedir = $operatingsystem ? {
        Solaris => "/.ssh",
        Debian => "/root/.ssh",
		Redhat => "/root/.ssh",
	}

        # Make sur .ssh exist in root home dir
	file { "$basedir/":
		ensure => directory,
		mode => 0700,
		owner => root,
		group => root,
		ignore => '.svn'
	}

    # Check if authorized_keys key file exist or create empty one
    file { "$basedir/authorized_keys":
        ensure => present,
    }

    # Check this line exist
    line { ssh_key:
  	file => "$basedir/authorized_keys",
  	line => "ssh-dss AAAAB3NzaC1....zG3ZA== root@puppet",
  	    ensure => present;
	}
}

# Check servoce status
class ssh::service {
  service { ssh:
    name => $operatingsystem ? {
      Solaris => "svc:/network/ssh:default",
      default => ssh
    },
    ensure    => running,
    enable    => true
  }
}
```

Then, compared to sudo, I have a notify which automatically restarts the service when the file is replaced by a new version. It's the ssh service with the "ensure => running" option, which will allow detection of the version change and restart.

#### templates

Since we use templates, we'll need to create a templates folder:

```bash
mkdir -p /etc/puppet/modules/ssh/templates
```

Then we'll create 2 files (sshd_config and sshd_config.solaris) because the configurations don't behave the same way (OpenSSH vs Sun SSH). However, I'll only cover the OpenSSH part here:

```
# Package generated configuration file
# See the sshd(8) manpage for details

# What ports, IPs and protocols we listen for
<% ssh_default_port.each do |val| -%>
Port <%= val -%>
<% end -%>
# Use these options to restrict which interfaces/protocols sshd will bind to
#ListenAddress ::
#ListenAddress 0.0.0.0
Protocol 2
# HostKeys for protocol version 2
HostKey /etc/ssh/ssh_host_rsa_key
HostKey /etc/ssh/ssh_host_dsa_key
#Privilege Separation is turned on for security
UsePrivilegeSeparation yes

# Lifetime and size of ephemeral version 1 server key
KeyRegenerationInterval 3600
ServerKeyBits 768

# Logging
SyslogFacility AUTH
LogLevel INFO

# Authentication:
LoginGraceTime 120
PermitRootLogin yes
StrictModes yes

RSAAuthentication yes
PubkeyAuthentication yes
#AuthorizedKeysFile	%h/.ssh/authorized_keys

# Don't read the user's ~/.rhosts and ~/.shosts files
IgnoreRhosts yes
# For this to work you will also need host keys in /etc/ssh_known_hosts
RhostsRSAAuthentication no
# similar for protocol version 2
HostbasedAuthentication no
# Uncomment if you don't trust ~/.ssh/known_hosts for RhostsRSAAuthentication
#IgnoreUserKnownHosts yes

# To enable empty passwords, change to yes (NOT RECOMMENDED)
PermitEmptyPasswords no

# Change to yes to enable challenge-response passwords (beware issues with
# some PAM modules and threads)
ChallengeResponseAuthentication no

# Change to no to disable tunnelled clear text passwords
#PasswordAuthentication yes

# Kerberos options
#KerberosAuthentication no
#KerberosGetAFSToken no
#KerberosOrLocalPasswd yes
#KerberosTicketCleanup yes

# GSSAPI options
#GSSAPIAuthentication no
#GSSAPICleanupCredentials yes

X11Forwarding yes
X11DisplayOffset 10
PrintMotd no
PrintLastLog yes
TCPKeepAlive yes
#UseLogin no

#MaxStartups 10:30:60
#Banner /etc/issue.net

# Allow client to pass locale environment variables
AcceptEnv LANG LC_*

Subsystem sftp /usr/lib/openssh/sftp-server

UsePAM yes
# AllowUsers <%= ssh_allowed_users %>
```

### SELinux

If you want to know more about SELinux, I invite you to look at this documentation.
Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/selinux/manifests
```

#### init.pp

```ruby
class selinux {
    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include selinux::redhat
        }
    }
}
```

#### redhat.pp

Here, I'll use a function (augeas), which will allow me to change the value of the 'SELINUX' variable to 'disabled', because I want to disable this module:

```ruby
class selinux::redhat
{
    # Disable SELinux
    augeas {
        "selinux" :
            context => "/files/etc/sysconfig/selinux/",
            changes => "set SELINUX disabled"
    }
}
```

### Grub

There are many Grub options and its use can vary from system to system. I recommend this documentation before tackling this.
Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/grub/manifests
```

#### init.pp

```ruby
/*
Grub Module for Puppet
Made by Pierre Mavro
*/
class grub {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include ::grub::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

Things will get a bit more complicated here. I'm still using the Augeas module to remove the quiet and rhgb arguments from the available kernels in grub.conf:

```ruby
/*
Grub Module for Puppet
Made by Pierre Mavro
*/
class grub::redhat
{
    # Remove unwanted parameters parameter
    augeas {
        "grub" :
            context => "/files/etc/grub.conf",
            changes => [
                "remove title[*]/kernel/quiet",
                "remove title[*]/kernel/rhgb"
            ]
    }
}
```

### Kdump

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/kdump/manifests
```

#### init.pp

```ruby
/*
Kdump Module for Puppet
Made by Pierre Mavro
*/
class kdump {
    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include kdump::redhat
        }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

Here's the configuration I want:

```ruby
/*
Kdump Module for Puppet
Made by Pierre Mavro
*/
class kdump::redhat
{
    # Install kexec-tools    package { 'kexec-tools':
        ensure => 'installed'
    }

    # Be sure that service is set to start at boot    service {
        'kdump':
            enable => true
    }

    # Set crashkernel in grub.conf to the good size (not auto)    augeas {
        "grub_kdump" :
            context => '/files/etc/grub.conf',
            changes => [
                'set title[1]/kernel/crashkernel 128M'
            ]
    }

    # Set location of crash dumps    line {
        'var_crash':
            file => '/etc/kdump.conf',
            line => 'path \/var\/crash',
            ensure => uncomment
    }
}
```

- Install kexec-tools: I make sure the package is installed
- Be sure that service is set to start at boot: I make sure it's enabled at boot. For information, I don't check if it's running since this would require a restart if it wasn't present.
- Set crashkernel in grub.conf to the good size (not auto): Sets the value 128M to the crashkernel argument of the first kernel found in the grub.conf file. It's not possible to specify '\*' as for a remove in augeas. However, given that with each kernel update, all others will inherit, this is not a problem :-)
- Set location of crash dumps: we specify in the kdump configuration that a line is indeed uncommented and has the desired path for crash dumps.

### Tools

This is not software, but rather a module that I use to send all my admin scripts, my tools, etc... Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/tools/{manifests,files}
```

#### init.pp

```ruby
class tools {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include ::tools::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

We'll see different ways to add files here:

```ruby
class tools::redhat {
    # Check that scripts folder exist    file {
        "/etc/scripts" :
            ensure => directory,
            mode => 0755,
            owner => root,
            group => root
    }

    # Synchro admin-scripts    file {
        "/etc/scripts/admin-scripts" :
            ensure => directory,
            mode => 0755,
            owner => root,
            group => root,
            source => "puppet:///modules/tools/admin-scripts/",
            purge => true,
            force => true,
            recurse => true,
            ignore => '.svn',
            backup => false
    }

    # Fast reboot command    file {
        "/usr/bin/fastreboot" :
            source => "puppet:///modules/tools/fastreboot",
            mode => 744,
            owner => root,
            group => root
    }
}
```

- Check that scripts folder exist: I make sure my folder exists with the right permissions before placing files in it.
- Synchro admin-scripts: I copy an entire directory with its contents:
  - purge => true: I make sure that anything not in my puppet files should disappear server-side. So if you've manually added a file to the /etc/scripts/admin-scripts folder, it will disappear.
  - force => true: Forces in case of deletion or replacement.
  - recurse => true: This allows me to say to copy all the content
  - backup => false: We don't back up the files before replacing them
- Fast reboot command: I add an executable file

#### files

Here's what my directory structure looks like in the files folder:

```
.
|-- admin-scripts
|   |-- script1.pl
|   `-- script2.pl
`-- fastreboot
```

For information, the fastreboot command is available here.

### Timezone

Managing timezones is more or less easy depending on the OS. We'll see how to handle it on Red Hat. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/timezone/{manifests,templates}
```

#### init.pp

```ruby
/*
Timezone Module for Puppet
Made by Pierre Mavro
*/
class timezone {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include timezone::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

We use a variable '$set_timezone' stored in the global variables or in the configuration of a node with the continent and country. Here's an example:

```ruby
# NTP Timezone. Usage :
# Look at /usr/share/zoneinfo/ and add the continent folder followed by the town
$set_timezone = 'Europe/Paris'
```

And the manifest:

```ruby
/*
Timezone Module for Puppet
Made by Pierre Mavro
*/
class timezone::redhat
{
    # Usage : set a var called set_timezone with required informations. Ex :
    # set_timezone = "Europe/Paris"

    # Set timezone file
    file { '/etc/sysconfig/clock':
        content => template("timezone/clock.$::operatingsystem"),
        mode => 644,
        owner => root,
        group => root
    }

    # Create required Symlink
    file { '/etc/localtime':
        ensure => link,
        target => "/usr/share/zoneinfo/${set_timezone}"
    }
}
```

#### templates

And the template file needed for Red Hat:

```
ZONE="<%= set_timezone %>"
```

### NTP

If you want to understand how to configure an NTP server, I invite you to read this documentation.

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/ntp/{manifests,files}
```

#### init.pp

```ruby
/*
NTP Module for Puppet
Made by Pierre Mavro
*/
class ntp {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include ntp::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

I want the package to be installed, configured at boot, to retrieve a standard configuration file, and to configure the crontab so that the service only starts outside production hours (so that production logs have consistency):

```ruby
/*
NTP Module for Puppet
Made by Pierre Mavro
*/
class ntp::redhat {
    # Install NTP service
    package {
        'ntp' :
            ensure => 'installed'
    }

    # Be sure that service is set to start at boot
    service {
        'ntpd' :
            enable => false
    }

    # Set configuration file
    file {
        '/etc/ntp.conf' :
            ensure => present,
            source => "puppet:///modules/ntp/${::osfamily}.ntp.conf",
            mode => 644,
            owner => root,
            group => root
    }

    # Enable ntp service during off production hours
    cron {
        'ntp_start' :
            command => '/etc/init.d/ntpd start',
            user => root,
            minute => 0,
            hour => 0
    }

    # Disable ntp service during on production hours
    cron {
        'ntp_stop' :
            command => '/etc/init.d/ntpd stop',
            user => root,
            minute => 3,
            hour => 0
    }
}
```

You'll notice the use of 'cron' directives in puppet which allow the management of crontab lines for a given user.

#### files

Here's the configuration file for Red Hat:

```
# For more information about this file, see the man pages
# ntp.conf(5), ntp_acc(5), ntp_auth(5), ntp_clock(5), ntp_misc(5), ntp_mon(5).

driftfile /var/lib/ntp/drift

# Permit time synchronization with our time source, but do not
# permit the source to query or modify the service on this system.
restrict default kod nomodify notrap nopeer noquery
restrict -6 default kod nomodify notrap nopeer noquery

# Permit all access over the loopback interface.  This could
# be tightened as well, but to do so would effect some of
# the administrative functions.
restrict 127.0.0.1
restrict -6 ::1

# Hosts on local network are less restricted.
#restrict 192.168.1.0 mask 255.255.255.0 nomodify notrap

# Use public servers from the pool.ntp.org project.
# Please consider joining the pool (http://www.pool.ntp.org/join.html).
server 0.rhel.pool.ntp.org
server 1.rhel.pool.ntp.org
server 2.rhel.pool.ntp.org

#broadcast 192.168.1.255 autokey	# broadcast server
#broadcastclient			# broadcast client
#broadcast 224.0.1.1 autokey		# multicast server
#multicastclient 224.0.1.1		# multicast client
#manycastserver 239.255.254.254		# manycast server
#manycastclient 239.255.254.254 autokey # manycast client

# Undisciplined Local Clock. This is a fake driver intended for backup
# and when no outside source of synchronized time is available.
#server	127.127.1.0	# local clock
#fudge	127.127.1.0 stratum 10

# Enable public key cryptography.
#crypto

includefile /etc/ntp/crypto/pw

# Key file containing the keys and key identifiers used when operating
# with symmetric key cryptography.
keys /etc/ntp/keys

# Specify the key identifiers which are trusted.
#trustedkey 4 8 42

# Specify the key identifier to use with the ntpdc utility.
#requestkey 8

# Specify the key identifier to use with the ntpq utility.
#controlkey 8

# Enable writing of statistics records.
#statistics clockstats cryptostats loopstats peerstats
```

### MySecureShell

The configuration of MySecureShell is not very complex, only it generally differs from one machine to another, especially if you have a large fleet. You probably want to have identical global management and be able to do custom on certain users, groups, virtualhost, etc... That's why we'll manage this in the simplest way possible, that is to say a global configuration, then includes.

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/mysecureshell/{manifests,files}
```

#### init.pp

```ruby
/*
MySecureShell Module for Puppet
Made by Pierre Mavro
*/
class mysecureshell {
    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include mysecureshell::redhat
        }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

Nothing magical here, I make sure the package is installed and the configuration correctly pushed:

```ruby
/*
MySecureShell Module for Puppet
Made by Pierre Mavro
*/
class mysecureshell::redhat
{
     # Install MySecureShell
     package { 'mysecureshell':
         ensure => 'installed'
     }

     # MySecureShell default configuration
     file {
          '/etc/ssh/sftp_config' :
               ensure => present,
               source => "puppet:///modules/mysecureshell/${::osfamily}.sftp_config",
               mode => 644,
               owner => root,
               group => root
     }
}
```

#### files

Here's my global file, at the end of which you'll notice there's an include. Feel free to put what you want in it, have multiple ones, and send the right one according to criteria:

```
# HEADER: This file is managed by puppet.
# HEADER: While it can still be managed manually, it is definitely not recommended.

#Default rules for everybody
<Default>
	GlobalDownload		0	#total speed download for all clients
	GlobalUpload		0	#total speed download for all clients (0 for unlimited)
	Download 			0 	#limit speed download for each connection
	Upload 				0	#unlimit speed upload for each connection
	StayAtHome			true	#limit client to his home
	VirtualChroot		true	#fake a chroot to the home account
	LimitConnection		100	#max connection for the server sftp
	LimitConnectionByUser	2	#max connection for the account
	LimitConnectionByIP	2	#max connection by ip for the account
	Home			/home/$USER	#overrite home of the user but if you want you can use
	IdleTimeOut		5m	#(in second) deconnect client is idle too long time
	ResolveIP		false	#resolve ip to dns
	IgnoreHidden	true	#treat all hidden files as if they don't exist
#	DirFakeUser		true	#Hide real file/directory owner (just change displayed permissions)
#	DirFakeGroup	true	#Hide real file/directory group (just change displayed permissions)
#	DirFakeMode		0400	#Hide real file/directory rights (just change displayed permissions)
#	HideFiles		"^(lost\+found|public_html)$"	#Hide file/directory which match
	HideNoAccess		true	#Hide file/directory which user has no access
#	MaxOpenFilesForUser	20	#limit user to open x files on same time
#	MaxWriteFilesForUser	10	#limit user to x upload on same time
#	MaxReadFilesForUser	10	#limit user to x download on same time
	DefaultRights		0640 0750	#Set default rights for new file and new directory
#	MinimumRights		0400 0700	#Set minimum rights for files and dirs
#	PathDenyFilter		"^\.+"	#deny upload of directory/file which match this extented POSIX regex
	ShowLinksAsLinks	false	#show links as their destinations
#	ConnectionMaxLife	1d	#limits connection lifetime to 1 day
#	Charset			"ISO-8859-15"	#set charset of computer
#	GMTTime			+1	#set GMT Time (change if necessary)
</Default>

# Include another custom file
Include /etc/ssh/deimos_sftp_config	#include this valid configuration file
```

### OpenLDAP client & Server

OpenLDAP can quickly become complicated and you need to understand how it works before deploying it. I invite you to look at the documentation before diving into this module.

This module is a big chunk, I spent a lot of time on it so that it works completely automatically. I'll be as explicit as possible so that everything is clear. I've adapted it for a nomenclature of names that allows this kind of thing in a computer fleet, but we can refer to pretty much anything.

In the case of this module, to understand it well, I need to explain how the infrastructure is constituted. We can assume that we have a batch or a cluster with 4 nodes. On these 4 nodes, only machines numbered 1 and 2 are OpenLDAP servers. All nodes are clients of the servers:

```
         +-------------------+      +-------------------+
         |  TO-DC-SRV-PRD-1  |      |  TO-DC-SRV-PRD-2  |
         |   OpenLDAP SRV1   |      |   OpenLDAP SRV2   |
         |  OpenLDAP Client  |      |  OpenLDAP Client  |
         +--------+----------+      +---------+---------+
                  |                           |
                  +---------------------------+
                  |                           |
         +--------+----------+       +--------+---------+
         |  TO-DC-SRV-PRD-3  |       | TO-DC-SRV-PRD-4  |
         |  OpenLDAP Client  |       | OpenLDAP Client  |
         +-------------------+       +------------------+
```

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/openldap/{manifests,files,lib/facter,puppet/parser/functions}
```

#### init.pp

Here I load all my manifests to be able to call their classes later:

```ruby
/*
OpenLDAP Module for Puppet
Made by Pierre Mavro
*/
class openldap inherits openssh {
    import "*.pp"

    # Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include openldap::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

I'll handle all OpenLDAP management here:

- Lines 6 to 17: Preparation
- Line 20: Client configuration
- Lines 22 to 27: Server configuration

```ruby
/*
OpenLDAP Module for Puppet
Made by Pierre Mavro
*/
class openldap::redhat inherits openldap {
    # Select generated ldap servers if no one is specified or join specified ones
    if empty("${ldap_servers}") == true
    {
        $comma_ldap_servers = "$::generated_ldap_servers"
    }
    else
    {
        $comma_ldap_servers = inline_template("<%= (ldap_servers).join(',') %>")
    }

    # DNS lookup them to use IP address instead
    $ip_ldap_servers = dns2ip("${comma_ldap_servers}")

    # OpenLDAP Client
    include openldap::redhat::ldapclient

    # Check if the current host is a server or not
    $ip_hostname = dns2ip($::hostname)
    $array_ldap_servers = split($ip_ldap_servers, ',')
    if $ip_hostname in $array_ldap_servers {
        # OpenLDAP Server
        include openldap::redhat::ldapserver
    }
}
```

We use a variable '$ldap_servers' stored in the global variables or in the configuration of a node with the list of default LDAP servers desired:

```ruby
# Default LDAP servers
#$ldap_servers = [ '192.168.0.1', '192.168.0.2' ]
$ldap_servers = [ ]
```

If we use, as here, an empty array, the empty parser (line 7) will detect it and we will then use the automatic generation of the list of LDAP servers to use. This method uses a facter to detect, based on the machine name, its equivalent 1 and 2.

Then, it will do a reverse lookup to convert to IP, any DNS names that would appear in the list of LDAP servers. For this, we use a parser named dns2ip. Then there will be the call to the OpenLDAP client manifest.

And finally, if the current server appears in the list of OpenLDAP servers, then we apply its configuration.

#### facter

##### generated_ldap_servers.rb

This facter retrieves the hostname of the machine in question, if it corresponds to the desired nomenclature, it sends in an array the auto-generated names of the LDAP servers. Then it returns as a result the array joined by commas:

```ruby
# Generated ldap servers
Facter.add("generated_ldap_servers") do
  ldap_srv = []

  # Get current hostname and add -1 and -2 to define default LDAP servers
  hostname = Facter.value('hostname')
  if hostname =~ /(.*)-\d+$/
    ldap_srv.push($1 + '-1', $1 + '-2')
  end

  setcode do
    # Join them with a comma
    ldap_srv.join(",")
  end
end
```

##### current_ldap_servers.rb

This facter reads the OpenLDAP configuration file to see what configuration is applied to it and returns the result separated by a comma:

```ruby
# Current ldap servers
Facter.add("current_ldap_servers") do
  ldap_srv = []
  conf_ldap = []
  config_found = 0

  # Add in an array all ldap servers lines
  f = File.open('/etc/openldap/ldap.conf', 'r')
  f.each_line do |line|
    if line =~ /^URI (.*)/
      config_found = 1
      conf_ldap = $1.split(" ")
    end
  end
  f.close

  # Get all LDAP servers names/IP
  if conf_ldap.each do |line|
      if line =~ /ldap:\/\/(.*)\//
        ldap_srv.push($1)
      end
    end
  end

  setcode do
    if config_found == 1
      # Join them with a comma
      ldap_srv.join(",")
    else
      config_found
    end
  end
end
```

#### Parser

##### dns2ip.rb

This parser returns, once processed, a list of machines separated by commas. The process is actually the conversion of DNS name to IP:

```ruby
# Dns2IP for Puppet
# Made by Pierre Mavro
# Does a DNS lookup and returns an array of strings of the results
# Usage : need to send one string dns servers separated by comma. The return will be the same

require 'resolv'

module Puppet::Parser::Functions
  newfunction(:dns2ip, :type => :rvalue) do |arguments|
    result = [ ]
    # Split comma sperated list in array
    dns_array = arguments[0].split(',')
    # Push each DNS/IP address in result array
    dns_array.each do |dns_name|
      result.push(Resolv.new.getaddresses(dns_name))
    end
    # Join array with comma
    dns_list = result.join(',')
    # Delete last comma if exist
    good_dns_list = dns_list.gsub(/,$/, '')
    return good_dns_list
  end
end
```

#### redhat_ldapclient.pp

This manifest will install all the necessary packages for the server to become an OpenLDAP client, then make sure that 2 services (you'll notice the use of arrays) are running and that they are started at boot. Then we'll use the previously mentioned facters, as well as other custom variables to validate that the configuration on the server exists or not. In the negative case, we execute a command that will configure the OpenLDAP client:

```ruby
/*
OpenLDAP Module for Puppet
Made by Pierre Mavro
*/
class openldap::redhat::ldapclient inherits openldap::redhat {
	# Install required clients pacakges
	packages_install
	{ [
		'nss-pam-ldapd',
		'openldap',
      	'openldap-clients',
      	'openssh-ldap'
	]: }

	# Be sure that service is set to start at boot
	service {
		['nscd', 'nslcd'] :
			enable => true,
			ensure => running,
			require => Package['nss-pam-ldapd', 'openldap', 'openldap-clients', 'openssh-ldap']
	}

	# Configure LDAP client if current configuration doesn't match
	if ($::current_ldap_servers != $ip_ldap_servers or $::current_ldap_servers == 0)
	{
		exec {
			"authconfig --enableldap --enableldapauth --disablenis --enablecache --passalgo=sha512 --disableldaptls --disableldapstarttls --disablesssdauth --enablemkhomedir --enablepamaccess --enablecachecreds --enableforcelegacy --disablefingerprint --ldapserver=${ip_ldap_servers} --ldapbasedn=dc=openldap,dc=deimos,dc=fr --updateall" :
				logoutput => true,
				require => Service['nslcd']
		}
	}
}
```

#### redhat_ldapserver.pp

For the server part, we'll check the existence of all folders, verify that packages and services are correctly configured, send the custom schemas and OpenLDAP configuration:

```ruby
/*
OpenLDAP Module for Puppet
Made by Pierre Mavro
*/
class openldap::redhat::ldapserver inherits openldap::redhat {
# Install required clients pacakges
    packages_install {
        ['openldap-servers'] :
    }

    # Copy schemas
    file {
        "/etc/openldap/schema" :
            ensure => directory,
            mode => 755,
            owner => root,
            group => root,
            source => "puppet:///openldap/schema/",
            recurse => true
    }

    # Copy configuration folder (new format)
    file {
        "/etc/openldap/slapd.d" :
            ensure => directory,
            mode => 700,
            owner => ldap,
            group => ldap,
            source => "puppet:///openldap/$::operatingsystem/slapd.d/",
            recurse => true
    }

    # Copy configuration file (old format)
    file {
        "/etc/openldap/slapd.conf" :
            ensure => present,
            mode => 700,
            owner => ldap,
            group => ldap,
            source => "puppet:///openldap/$::operatingsystem/slapd.conf"
    }

    # Ensure rights are ok for folder LDAP database
    file {
        "/var/lib/ldap" :
            ensure => directory,
            mode => 700,
            owner => ldap,
            group => ldap
    }

    # Ensure rights are ok for folder pid and args
    file {
        "/var/run/openldap" :
            ensure => directory,
            mode => 755,
            owner => ldap,
            group => ldap
    }

    # Be sure that service is set to start at boot
    service {
        'slapd' :
            enable => true,
            ensure => running,
            require => File['/etc/openldap/slapd.conf']
    }
}
```

#### files

To give you an idea of the content of files:

```
.
`-- RedHat
    |-- slapd.conf
    `-- slapd.d
        |-- cn=config
        |   |-- cn=module{0}.ldif
        |   |-- cn=schema
        |   |   |-- cn={0}core.ldif
        |   |   |-- cn={1}cosine.ldif
        |   |   |-- cn={2}nis.ldif
        |   |   |-- cn={3}inetorgperson.ldif
        |   |   |-- cn={4}microsoft.ldif
        |   |   |-- cn={5}microsoft.ldif
        |   |   |-- cn={6}microsoft.ldif
        |   |   |-- cn={7}samba.ldif
        |   |   `-- cn={8}deimos.ldif
        |   |-- cn=schema.ldif
        |   |-- olcDatabase={-1}frontend.ldif
        |   |-- olcDatabase={0}config.ldif
        |   `-- olcDatabase={1}bdb.ldif
        `-- cn=config.ldif
```

And for the LDAP service configuration:

```
# Schema and objectClass definitions
include         /etc/openldap/schema/core.schema
include         /etc/openldap/schema/cosine.schema
include         /etc/openldap/schema/nis.schema
include         /etc/openldap/schema/inetorgperson.schema
include         /etc/openldap/schema/microsoft.schema
include         /etc/openldap/schema/microsoft.sfu.schema
include         /etc/openldap/schema/microsoft.exchange.schema
include         /etc/openldap/schema/samba.schema
include	    /etc/openldap/schema/deimos.schema


# Where the pid file is put. The init.d script
# will not stop the server if you change this.
pidfile         /var/run/openldap/slapd.pid

# List of arguments that were passed to the server
argsfile        /var/run/openldap/slapd.args

# Read slapd.conf(5) for possible values
loglevel        0

# Where the dynamically loaded modules are stored
modulepath	/usr/libexec/openldap
moduleload	back_bdb
#moduleload	back_meta

# The maximum number of entries that is returned for a search operation
sizelimit 500

# The tool-threads parameter sets the actual amount of cpu's that is used
# for indexing.
#tool-threads 1

#######################################################################
# Specific Backend Directives for bdb:
# Backend specific directives apply to this backend until another
# 'backend' directive occurs
backend		bdb

#######################################################################
# Specific Directives for database #1, of type bdb:
# Database specific directives apply to this databasse until another
# 'database' directive occurs
database        bdb

# The base of your directory in database #1
suffix          "dc=openldap,dc=deimos,dc=fr"

# rootdn directive for specifying a superuser on the database. This is needed
# for syncrepl.
rootdn          "cn=admin,dc=openldap,dc=deimos,dc=fr"
rootpw		{SSHA}Sh+Yd...

# Where the database file are physically stored for database #1
directory       "/var/lib/ldap"

# For the Debian package we use 2MB as default but be sure to update this
# value if you have plenty of RAM
dbconfig set_cachesize 0 2097152 0

# Number of objects that can be locked at the same time.
dbconfig set_lk_max_objects 1500
# Number of locks (both requested and granted)
dbconfig set_lk_max_locks 1500
# Number of lockers
dbconfig set_lk_max_lockers 1500

# Indexing options for database #1
index       objectClass eq
index       cn,sn,uid,mail  pres,eq,sub
index       mailnickname,userprincipalname,proxyaddresses  pres,eq,sub
index       entryUUID,entryCSN eq


# Save the time that the entry gets modified, for database #1
lastmod         on

access to attrs=userPassword,shadowLastChange
        by dn="cn=replica,ou=Gestion Admin,ou=Utilisateurs,dc=openldap,dc=deimos,dc=fr" write
        by anonymous auth
        by self write
        by * none

access to *
        by dn="cn=replica,ou=Gestion Admin,ou=Utilisateurs,dc=openldap,dc=deimos,dc=fr" write
        by * read

access to dn.base="" by * read
```

### sudo

Sudo quickly becomes indispensable when you want to give privileged rights to groups or users. Check out this documentation if you need more information.

#### init.pp

The init.pp file is the heart of our module, fill it in like this for now:

```ruby
# sudo.pp
class sudo {
    # OS detection
    $sudoers_file = $operatingsystem ? {
	Solaris => "/opt/csw/etc/sudoers",
	default => "/etc/sudoers"
    }

    # Sudoers file declaration
    file { "$sudoers_file":
        owner   => root,
        group   => root,
        mode    => 640,
        source  => $operatingsystem ? {
        	Solaris => "puppet:///modules/sudo/sudoers.solaris",
        	default => "puppet:///modules/sudo/sudoers"
        }
    }

    # Symlink for solaris
    case $operatingsystem {
	Solaris: {
            file {"/opt/sfw/bin/sudo":
                ensure => "/opt/csw/bin/sudo"
            }
        }
    }
}
```

I'll try to be clear in the explanations:

- We declare a sudo class
- We make an OS detection. This declaration is inherited from a variable ($sudoers_file) which allows us to declare different paths for the sudoers file.
- Which includes a configuration file located at (name) /opt/csw/etc/sudoers on the target system (for this config, it's on Solaris, adapt as needed)
- The file must belong to user and group root with permissions 440.
- The source of this file (which we haven't deposited yet) is available (source) at puppet:///modules/sudo/sudoers (or /etc/puppet/modules/sudo/files/sudoers). It simply indicates the location of the files in your puppet tree which uses a mechanism of file system internal to Puppet.

##### Improving a module

Let's edit the file to add some requests:

```ruby
# sudo.pp
class sudo {

    # Check if sudo is the latest version
    package { sudo: ensure => latest }

    # OS detection
    $sudoers_file = $operatingsystem ? {
	Solaris => "/opt/csw/etc/sudoers",
	default => "/etc/sudoers"
    }

    # Sudoers file declaration
    file { "$sudoers_file":
        owner   => root,
        group   => root,
        mode    => 640,
        source  => $operatingsystem ? {
        	Solaris => "puppet:///modules/sudo/sudoers.solaris",
        	default => "puppet:///modules/sudo/sudoers"
        },
        require => Package["sudo"]
    }
}
```

- require => Package["sudo"]: we ask it to install the sudo package if it's not installed
- package { sudo: ensure => latest }: we ask it to check that it's the latest version

#### files

Now, we need to add the files we want to publish in the files folder (here only sudoers):

```bash
cp /etc/sudoers /etc/puppet/modules/sudo/files/sudoers
cp /etc/sudoers /etc/puppet/modules/sudo/files/sudoers.solaris
```

To keep it simple, I simply copied the sudoers from the server to the destination that corresponds to my config described above.

If you haven't installed sudo on the server, put a sudoers file that suits you in /etc/puppet/modules/sudo/files/sudoers.

In future versions, puppet will support http and ftp protocols to fetch these files.

### snmpd

The snmpd service is quite simple to use. That's why this module is too :). Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/snmpd/{manifests,files}
```

#### init.pp

```ruby
/*
Snmpd Module for Puppet
Made by Pierre Mavro
*/
class snmpd {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include snmpd::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

```ruby
/*
Snmpd Module for Puppet
Made by Pierre Mavro
*/
class snmpd::redhat {
    # Install snmp packages
    package {
        'net-snmp':
            ensure => present
    }

    # Snmpd config file
    file {
        "/etc/snmp/snmpd.conf":
            source => "puppet:///modules/snmpd/snmpd.conf.$::operatingsystem",
            mode => 644,
            owner => root,
            group => root,
            notify  => Service["snmpd"]
    }

    # Service should start on boot and be running
    service {
        'snmpd':
            enable => true,
            ensure => running,
            require => File['/etc/snmp/snmpd.conf']
    }
}
```

#### files

And the configuration file:

```
# Generated by Puppet
###########################################################################
#
# snmpd.conf
#
#   - created by the snmpconf configuration program
#

###########################################################################
# SECTION: System Information Setup
#
#   This section defines some of the information reported in
#   the "system" mib group in the mibII tree.

# syslocation: The [typically physical] location of the system.
#   Note that setting this value here means that when trying to
#   perform an snmp SET operation to the sysLocation.0 variable will make
#   the agent return the "notWritable" error code.  IE, including
#   this token in the snmpd.conf file will disable write access to
#   the variable.
#   arguments:  location_string

syslocation Unknown (edit /etc/snmp/snmpd.conf)

# syscontact: The contact information for the administrator
#   Note that setting this value here means that when trying to
#   perform an snmp SET operation to the sysContact.0 variable will make
#   the agent return the "notWritable" error code.  IE, including
#   this token in the snmpd.conf file will disable write access to
#   the variable.
#   arguments:  contact_string

syscontact Root <xxx@mycompany.com> (configure /etc/snmp/snmp.local.conf)

###########################################################################
# SECTION: Access Control Setup
#
#   This section defines who is allowed to talk to your running
#   snmp agent.

# rocommunity: a SNMPv1/SNMPv2c read-only access community name
#   arguments:  community [default|hostname|network/bits] [oid]

rocommunity  lookdeimos

#
# Unknown directives read in from other files by snmpconf
#
com2sec notConfigUser  default       public
group   notConfigGroup v1           notConfigUser
group   notConfigGroup v2c           notConfigUser
view    systemview    included   .1.3.6.1.2.1.1
view    systemview    included   .1.3.6.1.2.1.25.1.1
access  notConfigGroup ""      any       noauth    exact  systemview none none
dontLogTCPWrappersConnects yes
```

### Postfix

For postfix, we'll use templates to facilitate certain configurations. If you need more information on how to configure Postfix, I recommend these documentations.

#### init.pp

```ruby
/*
Postix Module for Puppet
Made by Pierre Mavro
*/
class postfix {
# Check OS and request the appropriate function
    case $::operatingsystem {
        'RedHat' : {
            include postfix::redhat
        }
        #'sunos':  { include packages_defaults::solaris }
        default : {
            notice("Module ${module_name} is not supported on ${::operatingsystem}")
        }
    }
}
```

#### redhat.pp

I use transport maps here that need to be rebuilt if their configuration changes. That's why we use 'Subscribe' in the 'Exec' directive:

```ruby
/*
Postfix Module for Puppet
Made by Pierre Mavro
*/
class postfix::redhat {
# Install postfix packages
    package {
        'postfix' :
            ensure => present
    }

    # Postfix main config file
    file {
        "/etc/postfix/main.cf" :
            content => template("postfix/main.cf.$::operatingsystem"),
            mode => 644,
            owner => root,
            group => root,
            notify => Service["postfix"]
    }

    # Postfix transport map
    file {
        "/etc/postfix/transport" :
            source => "puppet:///modules/postfix/transport",
            mode => 644,
            owner => root,
            group => root,
            notify => Service["postfix"]
    }

    # Rebuild the transport map
    exec {
        'build_transport_map' :
            command => "/usr/sbin/postmap /etc/postfix/transport",
            subscribe => File["/etc/postfix/transport"],
            refreshonly => true
    }

    # Service should start on boot and be running
    service {
        'postfix' :
            enable => true,
            ensure => running,
            require => File['/etc/postfix/main.cf']
    }
}
```

#### templates

Here I have my fqdn which is unique depending on the machine:

```
# Generated by Puppet

# Postfix directories
queue_directory = /var/spool/postfix
command_directory = /usr/sbin
daemon_directory = /usr/libexec/postfix
data_directory = /var/lib/postfix
sendmail_path = /usr/sbin/sendmail.postfix
newaliases_path = /usr/bin/newaliases.postfix
mailq_path = /usr/bin/mailq.postfix

# Inet configuration
inet_interfaces = all
inet_protocols = all

# Reject unknow recipents
unknown_local_recipient_reject_code = 550

# Do not set relayhost. Postfix must use transport_maps
relayhost =

# Destinations
mydomain = deimos.fr
myorigin = <%= fqdn %>
mydestination = $myorigin, localhost.$mydomain, localhost

# Transport_maps permits to route email using Google exchangers
transport_maps = dbm:/etc/opt/csw/postfix/transport
# Add TLS to emails if possible
smtp_tls_security_level = may

# Masquerade_domains hides hostnames from addresses
masquerade_domains = $mydomain

# Aliases
alias_maps = hash:/etc/aliases
alias_database = hash:/etc/aliases

# SMTP banner
smtpd_banner = $myhostname ESMTP $mail_name

# Debug
debug_peer_level = 2
debugger_command =
	 PATH=/bin:/usr/bin:/usr/local/bin:/usr/X11R6/bin
	 ddd $daemon_directory/$process_name $process_id & sleep 5

# Postfix rights
mail_owner = postfix
setgid_group = postdrop

# Helps
html_directory = no
manpage_directory = /usr/share/man
sample_directory = /usr/share/doc/postfix-2.6.6/samples
readme_directory = /usr/share/doc/postfix-2.6.6/README_FILES
```

#### files

And in my files, I have my transport_map file:

```
deimos.fr      :google.com
.deimos.fr     :google.com
```

### Nsswitch

For nsswitch, we'll use an advanced technique which consists of using Facter (a built-in that saves a lot of time). Facter offers, in short, specific environment variables for Puppet which allow conditions to be made based on this. For example, I want to check for the presence of a file that will tell me if my server is in cluster mode (for Sun Cluster) or not and modify the nsswitch file accordingly. For this I'll use facter.

Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/nsswitch/{manifests,lib/facter}
```

#### facter

Let's create our fact file:

```ruby
# is_cluster.rb

Facter.add("is_cluster") do
        setcode do
          #%x{/bin/uname -i}.chomp
		   FileTest.exists?("/etc/cluster/nodeid")
        end
end
```

#### init.pp

Then we'll specify the use of a template:

```ruby
class nsswitch {
	case $operatingsystem {
		Solaris:  { include nsswitch::solaris }
			default: { }
		}
}

class nsswitch::solaris {
	$nsfilename = $operatingsystem ? {
		Solaris => "/etc/nsswitch.conf",
		default => "/etc/nsswitch.conf"
    }

	config_file { "$nsfilename":
		content => template("nsswitch/nsswitch.conf"),
	}
}
```

#### templates

And now the configuration file that will call the facter:

```
#
# Copyright 2006 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#

#
# /etc/nsswitch.dns:
#
# An example file that could be copied over to /etc/nsswitch.conf; it uses
# DNS for hosts lookups, otherwise it does not use any other naming service.
#
# "hosts:" and "services:" in this file are used only if the
# /etc/netconfig file has a "-" for nametoaddr_libs of "inet" transports.

# DNS service expects that an instance of svc:/network/dns/client be
# enabled and online.

passwd:     files ldap
group:      files ldap

# You must also set up the /etc/resolv.conf file for DNS name
# server lookup.  See resolv.conf(4).
#hosts:      files dns

hosts:      <% if is_cluster -%>cluster<% end -%> files dns

# Note that IPv4 addresses are searched for in all of the ipnodes databases
# before searching the hosts databases.
#ipnodes:   files dns
#ipnodes:   files dns [TRYAGAIN=0]
ipnodes:   files dns [TRYAGAIN=0 ]

networks:   files
protocols:  files
rpc:        files
ethers:     files
#netmasks:   files
netmasks:   <% if is_cluster -%>cluster<% end -%> files
bootparams: files
publickey:  files
# At present there isn't a 'files' backend for netgroup;  the system will
#   figure it out pretty quickly, and won't use netgroups at all.
netgroup:   files
automount:  files
aliases:    files
services:   files
printers:       user files

auth_attr:  files
prof_attr:  files
project:    files

tnrhtp:     files
tnrhdb:     files
```

### Nagios NRPE + plugins

You must initialize the module before continuing.

#### init.pp

Let's say I have a folder of nagios plugins that I want to deploy everywhere. Sometimes I add them, I remove them, in short, I live my life with them and I want it to be fully synchronized. Here's how:

```ruby
#nagios_plugins.pp
class nrpe {
    # Copy conf and send checks
    include nrpe::common

    # Check operating system
    case $operatingsystem {
        Solaris: {
            include nrpe::service::solaris
        }
        default: { }
    }
}

class nrpe::common {
    class nrpe::configs {
        # Used for nrpe-deimos.cfg templates
        $nrpe_distrib = $operatingsystem ? {
            'Solaris' => "/opt/csw/libexec/nagios-plugins",
            'redhat' => "/home/deimos/checks/nrpe_scripts",
            'debian' => "/etc/nrpe",
        }

        # Used for nrpe-deimos.cfg templates
        $deimos_script = $operatingsystem ? {
                'Solaris' => "/opt/deimos/libexec",
                'redhat' => "/etc/nrpe.d",
                'debian' => "/etc/nrpe",
        }

        # Copy NRPE config file
        file { '/opt/csw/etc/nrpe.cfg':
            mode => 644, owner => root, group => root,
            content => template('nrpe/nrpe.cfg'),
        }

        ## Copy and adapt NRPE deimos Config file ##
        file { '/opt/csw/etc/nrpe-deimos.cfg':
            mode => 644, owner => root, group => root,
            content => template('nrpe/nrpe-deimos.cfg'),
        }
    }

    class nrpe::copy_checks {
        ## Copy deimos Production NRPE Checks ##

        file { "nrpe_prod_checks":
            name => $operatingsystem ? {
                'Solaris' => "/opt/deimos/libexec",
                'redhat' => "/home/deimos/checks/nrpe_scripts",
                'debian' => "/etc/nrpe",
            },
            ensure => directory,
            mode => 755, owner => root, group => root,
            sourceselect => all,
            source => $operatingsystem ? {
                'Solaris' => [ "puppet:///modules/nrpe/Nagios/nrpechecks/deimos", "puppet:///modules/nrpe/Nagios/nrpechecks/system" ],
            recurse => true,
            force => true,
            ignore => '.svn'
        }
    }

    include nrpe::configs
    include nrpe::copy_checks
}

class nrpe::service::solaris {
    ## Check service and restart if needed
    file { '/var/svc/manifest/network/nagios-nrpe.xml':
        source => "puppet:///modules/nrpe/nagios-nrpe.xml",
        mode => 644, owner => root, group => sys,
        before => Service["svc:/network/cswnrpe:default"]
    }

    # Restart service if smf xml has changed
    exec { "`svccfg import /var/svc/manifest/network/nagios-nrpe.xml`" :
    	subscribe => File['/var/svc/manifest/network/nagios-nrpe.xml'],
    	refreshonly => true
	}

	# Restart if one of those files have changed
    service { "svc:/network/cswnrpe:default":
        ensure => running,
        manifest => "/var/svc/manifest/network/nagios-nrpe.xml",
        subscribe => [ File['/opt/csw/etc/nrpe.cfg'], File['/opt/csw/etc/nrpe-deimos.cfg'] ]
    }
}
```

Here:

- purge allows to delete files that no longer exist
- recurse: recursive
- force: allows to force
- before: allows to be executed before something else
- subscribe: subscribes to a dependency with respect to the value of it
- refreshonly: refreshes only if there are changes

With the subscribe, you can see here, we make lists like this: [ element_1, element_2, element_3... ]. However, small precision, the change operates (here a restart) only if one of the elements in the list is modified and not all.

You can also see from lines 54 to 56, we do multi sourcing which allows us to specify multiple sources and based on these to send to one place the content of these various files.

#### files

For the file part, I linked with an external svn which allows me to free myself from the integration of plugins in the puppet part (which between us has nothing to do here).

#### templates

Just copy the nagios_plugins folder to files and make the templates here (I'll only do the nrpe.cfg):

```
...
command[check_load]=<%= nrpe_distrib %>/check_load -w 15,10,5 -c 30,25,20
command[sunos_check_rss_mem]=<%= deimos_script %>/check_rss_mem.pl -w $ARG1$ -c $ARG2$
...
```

### Munin

Disclaimer: this work is mostly based upon DavidS work, available on his [git repo]. In the scope of my work I needed to have munin support for freeBSD & Solaris. I also wrote a class for snmp_plugins & custom plugins. Some things are quite dependant from my infrastructure, like munin.conf generation script but it can easily be adapted to yours, by extracting data from your CMDB.

It requires the munin_interfaces fact published here (and merged into DavidS repo, thanks to him), and [Volcane's extlookup function] to store some parameters. Enough talking, this is the code:

```ruby
# Munin config class
# Many parts taken from David Schmitt's http://git.black.co.at/
# FreeBSD & Solaris + SNMP & custom plugins support by Nicolas Szalay <nico@gcu.info>

class munin::node {
	case $operatingsystem {
		openbsd: {}
		debian: { include munin::node::debian}
		freebsd: { include munin::node::freebsd}
		solaris: { include munin::node::solaris}
		default: {}
	}
}

class munin::node::debian {

	package { "munin-node": ensure => installed }

	file {
	"/etc/munin":
		ensure => directory,
		mode => 0755,
		owner => root,
		group => root;

	"/etc/munin/munin-node.conf":
		source => "puppet://$fileserver/files/apps/munin/munin-node-debian.conf",
		owner => root,
		group => root,
		mode => 0644,
		before => Package["munin-node"],
		notify => Service["munin-node"],
	}

	service { "munin-node": ensure => running }

	include munin::plugins::linux
}

class munin::node::freebsd {
	package { "munin-node": ensure => installed, provider => freebsd }

        file { "/usr/local/etc/munin/munin-node.conf":
                source => "puppet://$fileserver/files/apps/munin/munin-node-freebsd.conf",
                owner => root,
                group => wheel,
                mode => 0644,
                before => Package["munin-node"],
                notify => Service["munin-node"],
        }

	service { "munin-node": ensure => running }

	include munin::plugins::freebsd
}

class munin::node::solaris {
	# "hand made" install, no package.
	file { "/etc/munin/munin-node.conf":
		source => "puppet://$fileserver/files/apps/munin/munin-node-solaris.conf",
                owner => root,
                group => root,
                mode => 0644
	}

	include munin::plugins::solaris
}

class munin::gatherer {
	package { "munin":
		ensure => installed
	}

	# custom version of munin-graph : forks & generates many graphs in parallel
	file { "/usr/share/munin/munin-graph":
		owner => root,
		group => root,
		mode => 0755,
		source => "puppet://$fileserver/files/apps/munin/gatherer/munin-graph",
		require => Package["munin"]
	}

	# custon version of debian cron file. Month & Year cron are generated once daily
	file { "/etc/cron.d/munin":
		owner => root,
		group => root,
		mode => 0644,
		source => "puppet://$fileserver/files/apps/munin/gatherer/munin.cron",
		require => Package["munin"]
	}

	# Ensure cron is running, to fetch every 5 minutes
	service { "cron":
		ensure => running
	}

	# Ruby DBI for mysql
	package { "libdbd-mysql-ruby":
		ensure => installed
	}

	# config generator
	file { "/opt/scripts/muningen.rb":
		owner => root,
		group => root,
		mode => 0755,
		source => "puppet://$fileserver/files/apps/munin/gatherer/muningen.rb",
		require => Package["munin", "libdbd-mysql-ruby"]
	}

	# regenerate munin's gatherer config every hour
	cron { "munin_config":
		command => "/opt/scripts/muningen.rb > /etc/munin/munin.conf",
		user => "root",
		minute => "0",
		require => File["/opt/scripts/muningen.rb"]
	}

	include munin::plugins::snmp
	include munin::plugins::linux
	include munin::plugins::custom::gatherer
}


# define to create a munin plugin inside the right directory
define munin::plugin ($ensure = "present") {

	case $operatingsystem {
		freebsd: {
			$script_path = "/usr/local/share/munin/plugins"
			$plugins_dir = "/usr/local/etc/munin/plugins"
		}
		debian: {
			$script_path = "/usr/share/munin/plugins"
			$plugins_dir = "/etc/munin/plugins"
		}
		solaris: {
			$script_path = "/usr/local/munin/lib/plugins"
			$plugins_dir = "/etc/munin/plugins"
		}
		default: { }
	}

	$plugin = "$plugins_dir/$name"

	case $ensure {
		"absent": {
			debug ( "munin_plugin: suppressing $plugin" )
			file { $plugin: ensure => absent, }
		}

		default: {
			$plugin_src = $ensure ? { "present" => $name, default => $ensure }

			file { $plugin:
				ensure => "$script_path/${plugin_src}",
				require => Package["munin-node"],
				notify => Service["munin-node"],
			}
		}
	}
}

# snmp plugin define, almost same as above
define munin::snmp_plugin ($ensure = "present") {
	$pluginname = get_plugin_name($name)

	case $operatingsystem {
		freebsd: {
			$script_path = "/usr/local/share/munin/plugins"
			$plugins_dir = "/usr/local/etc/munin/plugins"
		}
		debian: {
			$script_path = "/usr/share/munin/plugins"
			$plugins_dir = "/etc/munin/plugins"
		}
		solaris: {
			$script_path = "/usr/local/munin/lib/plugins"
			$plugins_dir = "/etc/munin/plugins"
		}
		default: { }
	}

	$plugin = "$plugins_dir/$name"

	case $ensure {
		"absent": {
			debug ( "munin_plugin: suppressing $plugin" )
			file { $plugin: ensure => absent, }
		}

		"present": {
			file { $plugin:
				ensure => "$script_path/${pluginname}",
				require => Package["munin-node"],
				notify => Service["munin-node"],
			}
		}
	}
}

class munin::plugins::base
{
	case $operatingsystem {
		debian: { $plugins_dir = "/etc/munin/plugins" }
		freebsd: { $plugins_dir = "/usr/local/etc/munin/plugins" }
		solaris: { $plugins_dir = "/etc/munin/plugins" }
		default: {}
	}

	file { $plugins_dir:
		source => "puppet://$fileserver/files/empty",
		ensure => directory,
		checksum => mtime,
		ignore => ".svn*",
		mode => 0755,
		recurse => true,
		purge => true,
		force => true,
		owner => root
	}
}

class munin::plugins::interfaces
{
	$ifs = gsub(split($munin_interfaces, " "), "(.+)", "if_\\1")
	$if_errs = gsub(split($munin_interfaces, " "), "(.+)", "if_err_\\1")
	plugin {
		$ifs: ensure => "if_";
		$if_errs: ensure => "if_err_";
	}

	include munin::plugins::base
}

class munin::plugins::linux
{
	plugin { [ cpu, load, memory, swap, irq_stats, df, processes, open_files, ntp_offset, vmstat ]:
		ensure => "present"
	}

	include munin::plugins::base
	include munin::plugins::interfaces
}

class munin::plugins::nfsclient
{
	plugin { "nfs_client":
		ensure => present
	}
}

class munin::plugins::snmp
{
	# initialize plugins
	$snmp_plugins=extlookup("munin_snmp_plugins")
	snmp_plugin { $snmp_plugins:
		ensure => present
	}

	# SNMP communities used by plugins
	file { "/etc/munin/plugin-conf.d/snmp_communities":
		owner => root,
		group => root,
		mode => 0644,
		source => "puppet://$fileserver/files/apps/munin/gatherer/snmp_communities"
	}

}

define munin::custom_plugin($ensure = "present", $location = "/etc/munin/plugins") {
	$plugin = "$location/$name"

	case $ensure {
		"absent": {
			file { $plugin: ensure => absent, }
		}

		"present": {
			file { $plugin:
				owner => root,
				mode => 0755,
				source => "puppet://$fileserver/files/apps/munin/custom_plugins/$name",
				require => Package["munin-node"],
				notify => Service["munin-node"],
			}
		}
	}
}

class munin::plugins::custom::gatherer
{
	$plugins=extlookup("munin_custom_plugins")
	custom_plugin { $plugins:
		ensure => present
	}
}

class munin::plugins::freebsd
{
	plugin { [ cpu, load, memory, swap, irq_stats, df, processes, open_files, ntp_offset, vmstat ]:
		ensure => "present",
	}

	include munin::plugins::base
	include munin::plugins::interfaces
}

class munin::plugins::solaris
{
	# Munin plugins on solaris are quite ... buggy. Will need rewrite / custom plugins.
	plugin { [ cpu, load, netstat ]:
		ensure => "present",
	}

	include munin::plugins::base
	include munin::plugins::interfaces
}
```

#### Munin Interfaces

Everyone using puppet knows DavidS awesome git repository: git.black.co.at. Unfornately for me, his puppet infrastructure seems to be almost only linux based. I have different OS in mine, including FreeBSD & OpenSolaris. Looking at his module-munin I decided to reuse it (and not recreate the wheel) but he used a custom fact that needed some little work. So this is a FreeBSD & (Open)Solaris capable version, to know what network interfaces have link up.

```ruby
# return the set of active interfaces as an array
# taken from http://git.black.co.at
# modified by nico <nico@gcu.info> to add FreeBSD & Solaris support

Facter.add("munin_interfaces") do

	setcode do
		# linux
		if Facter.value('kernel') == "Linux" then
			`ip -o link show`.split(/\n/).collect do |line|
					value = nil
					matches = line.match(/^\d*: ([^:]*): <(.*,)?UP(,.*)?>/
					if !matches.nil?
						value = matches[1]
						value.gsub!(/@.*/, '')
					end
					value
			end.compact.sort.join(" ")
		#end

		# freebsd
		elsif Facter.value('kernel') == "FreeBSD" then
			Facter.value('interfaces').split(/,/).collect do |interface|
				status = `ifconfig #{interface} | grep status`
				if status != "" then
					status=status.strip!.split(":")[1].strip!
					if status == "active" then # I CAN HAZ LINK ?
						interface.to_a
					end
				end
			end.compact.sort.join(" ")
		#end

		# solaris
		elsif Facter.value('kernel') == "SunOS" then
			Facter.value('interfaces').split(/,/).collect do |interface|
				if interface != "lo0" then # /dev/lo0 does not exists
					status = `ndd -get /dev/#{interface} link_status`.strip!
					if status == "1" # ndd returns 1 for link up, 0 for down
						interface.to_a
					end
				end
			end.compact.sort.join(" ")
		end
	end
end
```

### Mcollective

Mcollective is a bit like a gas plant, but it's very powerful and works very well with Puppet. If you don't know or want to learn more, follow this link.

This Mcollective module for Puppet allows you to install mcollective, as well as modules on the client side (Mcollective servers). With the modules, here's what my directory structure looks like:

```
.
|-- files
|   |-- agent
|   |   |-- filemgr.rb
|   |   |-- nrpe.rb
|   |   |-- package.rb
|   |   |-- process.rb
|   |   |-- puppetd.rb
|   |   |-- service.rb
|   |   `-- shell.rb
|   |-- facts
|   |   `-- facter_facts.rb
|   `-- RedHat.server.cfg
`-- manifests
    |-- common.pp
    |-- init.pp
    `-- redhat.pp
```

For the modules (agents, facts), I invite you to look at my doc on Mcollective to know where to get them. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/mcollective/{manifests,files}
```

#### init.pp

The init.pp file is the heart of our module, fill it in like this:

```ruby
/*
Mcollective Module for Puppet
Made by Pierre Mavro
*/
class mcollective {
	# Check OS and request the appropriate function
	case $::operatingsystem {
		'RedHat' : {
			include mcollective::redhat
		}
		default : {
			notice("Module ${module_name} is not supported on ${::operatingsystem}")
		}
	}
}
```

#### common.pp

```ruby
/*
Mcollective Module for Puppet
Made by Pierre Mavro
*/
class mcollective::common {
	# Used for nrpe.cfg file
	$mco_agent_dir = $operatingsystem ? {
		'RedHat' => "/usr/libexec/mcollective/mcollective",
		'Solaris' => "/usr/libexec/mcollective/mcollective"
	}

	# Transfert agents
	file {
		'mco_agent_folder' :
			name => "$mco_agent_dir/agent",
			ensure => directory,
			mode => 644,
			owner => root,
			group => root,
			source => ["puppet:///modules/mcollective/agent"],
			recurse => true,
			ignore => '.svn',
			backup => false,
			notify => Service['mcollective'],
			require => Package['mcollective']
	}

	# Transfert facts
	file {
		'mco_facts_folder' :
			name => "$mco_agent_dir/facts",
			ensure => directory,
			mode => 644,
			owner => root,
			group => root,
			source => ["puppet:///modules/mcollective/facts"],
			recurse => true,
			ignore => '.svn',
			backup => false,
			notify => Service['mcollective'],
			require => Package['mcollective']
	}
}
```

#### redhat.pp

```ruby
/*
Mcollective Module for Puppet
Made by Pierre Mavro
*/
class mcollective::redhat {
# Install Mcollective client
	package { [
		'mcollective',
		'rubygem-stomp',
		'rubygem-sysproctable'
		]:
		ensure => 'installed'
	}

	# Be sure that service is set to start at boot and running
	service {
		'mcollective' :
			enable => true,
			ensure => running
	}

	# Set configuration file
	file {
		'/etc/mcollective/server.cfg' :
			ensure => present,
			source => "puppet:///modules/mcollective/${::osfamily}.server.cfg",
			mode => 640,
			owner => root,
			group => root,
			notify => Service['mcollective']
	}

	# Include common
	include mcollective::common
}
```

#### files

##### RedHat.server.cfg

Here's the configuration for Mcollective:

```
########################
# GLOCAL CONFIGURATION #
########################

topicprefix = /topic/
main_collective = mcollective
collectives = mcollective
libdir = /usr/share/mcollective/plugins
logfile = /var/log/mcollective.log
loglevel = info
daemonize = 1
classesfile = /var/lib/puppet/classes.txt

###########
# MODULES #
###########

# Security
securityprovider = psk
plugin.psk = unset

# Stomp
connector = stomp
plugin.stomp.host = mcollective.deimos.fr
plugin.stomp.port = 61613
plugin.stomp.user = mcollective
plugin.stomp.password = marionette

# AgentPuppetd
plugin.puppetd.puppetd = /usr/sbin/puppetd
plugin.puppetd.lockfile = /var/lib/puppet/state/puppetdlock
plugin.puppetd.statefile = /var/lib/puppet/state/state.yaml
plugin.puppet.pidfile = /var/run/puppet/agent.pid
plugin.puppetd.splaytime = 100
plugin.puppet.summary = /var/lib/puppet/state/last_run_summary.yaml

#########
# FACTS #
#########

factsource = facter
plugin.yaml = /etc/mcollective/facts.yaml
plugin.facter.facterlib = /var/lib/puppet/lib/facter
fact_cache_time = 300
```

### Bind

It can sometimes be useful to have a local DNS cache server to speed up certain processes and not be dependent on a third-party DNS if it fails for a few moments. For more info on Bind, follow this link. Let's create the directory structure:

```bash
mkdir -p /etc/puppet/modules/bind/{manifests,files,templates}
```

#### init.pp

```ruby
/*
Bind Module for Puppet
Made by Pierre Mavro
*/
class bind {
# Check OS and request the appropriate function
	case $::operatingsystem {
		'RedHat' : {
			include ::bind::redhat
		}
		#'sunos':  { include packages_defaults::solaris }
		default : {
			notice("Module ${module_name} is not supported on ${::operatingsystem}")
		}
	}
}
```

#### redhat.pp

```ruby
/*
Bind Module for Puppet
Made by Pierre Mavro
*/
class bind::redhat
{
	# Install bind package
	package {
		'bind' :
			ensure => present
	}

	# resolv.conf file
	file {
		"/etc/resolv.conf" :
			source => "puppet:///modules/bind/resolvconf",
			mode => 744,
			owner => root,
			group => root,
			require => File['/etc/named.conf']
	}

	# Bind main config file for cache server
	file {
		"/etc/named.conf" :
			content => template("bind/named.conf.$::operatingsystem"),
			mode => 640,
			owner => root,
			group => named,
			notify => Service["named"]
	}

	# Service should start on boot and be running
	service {
		'named' :
			enable => true,
			ensure => running,
			require => [ Package['bind'], File['/etc/named.conf', '/etc/resolv.conf'] ]
	}
}
```

#### files

We'll manage resolv.conf here:

```
# Generated by Puppet
domain deimos.fr
search deimos.fr deimos.lan
nameserver 127.0.0.1
```

#### templates

And finally, the configuration of the template that will act with the information filled in vars.pp:

```
// Generated by Puppet
//
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
        listen-on port 53 { 127.0.0.1; };
        listen-on-v6 port 53 { ::1; };
        directory       "/var/named";
        dump-file       "/var/named/data/cache_dump.db";
        statistics-file "/var/named/data/named_stats.txt";
        memstatistics-file "/var/named/data/named_mem_stats.txt";
        allow-query     { localhost; };
        recursion yes;

        dnssec-enable no;
        dnssec-validation no;
        dnssec-lookaside auto;

        forwarders {
            <% dns_servers.each do |dnsval| -%>            <%= dnsval %>;
            <% end -%>        };

        /* Path to ISC DLV key */
        bindkeys-file "/etc/named.iscdlv.key";

        managed-keys-directory "/var/named/dynamic";
};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

zone "." IN {
        type hint;
        file "named.ca";
};

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
```

### Importing a module

Modules need to be imported into puppet for them to be handled. Here, we don't have to worry about it because according to the server configuration we've made, it will automatically load the modules in /etc/puppet/modules. However, if you want to authorize module by module, you need to import them manually:

```ruby
# /etc/puppet/manifests/modules.pp

import "sudo"
import "ssh"
```

**WARNING**

Be careful because from the moment this is filled in, there's no need to restart puppet for the changes to be taken into account. So you'll need to be careful about what you make available at time t.

#### Importing all modules at once

This can have serious consequences, but know that it's possible to do it:

```ruby
# /etc/puppet/manifests/modules.pp

import "*.pp"
```

## Usage

### Certificates

Puppet works with certificates for client/server exchanges. So we'll need to generate SSL keys on the server and then exchange with all clients. The use of SSL keys therefore requires a good configuration of your DNS server. So check that:

- The server name is definitive
- The server name is accessible via puppet.mydomain.com (I'll continue the configuration for myself with puppet-prd.deimos.fr)

#### Creating a certificate

- So we'll create our certificate (**normally already done when the installation is done on Debian**):

```
puppet cert generate puppet-prd.deimos.fr
```

**WARNING**

It's imperative that **ALL names filled in the certificate are reachable** otherwise clients won't be able to synchronize with the server

- If you want to validate a certificate with multiple domain names, you'll need to proceed like this:

```
puppet cert cert generate --dns_alt_names puppet:scar.deimos.fr puppet-prd.deimos.fr
```

Insert the names you need one after the other. I remind you that by default, clients will look for **puppet.**mydomain.com, so don't hesitate to add names if needed.
You can then verify that the certificate contains both names:

```
> openssl x509 -text -in /var/lib/puppet/ssl/certs/puppet-prd.deimos.fr.pem | grep DNS
DNS:puppet, DNS:scar.deimos.fr, DNS:puppet-prd.deimos.fr
```

You can see possible certificate errors by connecting:

```
openssl s_client -host puppet -port 8140 -cert /path/to/ssl/certs/node.domain.com.pem -key /path/to/ssl/private_keys/node.domain.com.pem -CAfile /path/to/ssl/certs/ca.pem
```

_Note: Don't forget to restart your web server if you change certificates_

#### Adding a puppet client to the server

For the certificate, it's simple, we'll make a certificate request:

```
$ puppet agent --server puppet-prd.deimos.fr --waitforcert 60 --test
notice: Ignoring cache
info: Caching catalog at /var/lib/puppet/state/localconfig.yaml
notice: Starting catalog run
notice: //information_class/Exec[echo running on debian-puppet-prd.deimos.fr is a Debian with ip 192.168.0.106. Message is 'puppet c'est super']/returns: executed successfully
notice: Finished catalog run in 0.45 seconds
```

Now we'll connect **to the server** and run this command:

```
$ puppet cert -l
debian-puppet-prd.deimos.fr
```

Here I see, for example, that I have a host that wants to exchange keys in order to then obtain the configurations due to it. Only it needs to be authorized. To do this:

```
puppet cert -s debian-puppet-prd.deimos.fr
```

The machine is now accepted and can fetch configs from the Puppet server.

If you want to refuse a waiting node:

```
puppet cert -c debian-puppet-prd.deimos.fr
```

If you want to see the list of all authorized nodes, puppetmaster keeps them registered here:

```
0x0001 2010-03-08T15:45:48GMT 2015-03-07T15:45:48GMT /CN=puppet-prd.deimos.fr
0x0002 2010-03-08T16:36:16GMT 2015-03-07T16:36:16GMT /CN=node1.deimos.fr
0x0003 2010-03-08T16:36:25GMT 2015-03-07T16:36:25GMT /CN=node2.deimos.fr
0x0004 2010-04-14T12:41:24GMT 2015-04-13T12:41:24GMT /CN=node3.deimos.fr
```

#### Synchronizing a client

Now that our machines are connected, we'll want to synchronize them from time to time. Here's how to do a manual synchronization from a client machine:

```
puppet agent -t
```

If you want to synchronize only one module, you can use the --tags option:

```
puppet agent -t --tags module1 module2...
```

##### Simulating

If you need to test before actually deploying, there's the --noop option:

```
puppet agent -t --noop
```

You can also add the 'audit' directive in a manifest if you just want to audit a resource:

```ruby
file { "/etc/passwd":
    audit => [ owner, mode ],
}
```

Here we ask to audit the user and the rights, but know that it's possible to replace with 'all' to audit everything!

#### Revoking a certificate

- If you want to revoke a certificate, here's how to proceed on the server:

```
puppet cert clean my_machine
```

- Otherwise there's this method:

```
puppet cert -r my_machine
puppet cert -c my_machine
```

Simple, right? :-)

If you want to reassign again, delete the ssl folder on the client side in '/etc/puppet/' or '/var/lib/puppet/'. Then you can restart a certificate request.

#### Revoking all server certificates

If your puppet master is broken everywhere and you want to regenerate new keys and delete all the old ones:

```
puppet cert clean --all
rm -Rf /var/lib/puppet/ssl/certs/*
/etc/init.d/puppetmaster restart
```

### Process monitoring

On clients, it's not necessary to use monitoring software, since the puppetd daemons are only launched manually or by CRON task. On the server, it's important that the puppetmaster process is always present. You can use the 'monit' software, which can automatically restart the process in case of problems.

#### Determining node status

To see if there are problems on a node, you can manually run the following command:

```
puppetd --no-daemon --verbose --onetime
```

If you want to know the last state/result of a puppet update, you can use the 'report' system once activated (on the clients):

```
...
report = true
...
```

By enabling reporting, each time the puppetd daemon is executed, a report will be sent to the puppetmaster in a YAML format file, in the /var/lib/puppet/reports/MACHINE_NAME directory.

Here's an example of a report file, easily transformable into a dictionary with the yaml module of python (or ruby):

```ruby
--- !ruby/object:Puppet::Transaction::Report
host: deb-puppet-client.deimos.fr
logs:
- !ruby/object:Puppet::Util::Log
  level: :info
  message: Applying configuration version '1275982371'
  source: Puppet
  tags:
  - info
  time: 2010-06-08 11:37:59.521132 +02:00
- !ruby/object:Puppet::Util::Log
  file: /etc/puppet/modules/sudo/manifests/init.pp
  level: :err
  line: 11
  message: "Failed to retrieve current state of resource: Could not retrieve information from source(s) puppet:///sudo/sudoers at /etc/puppet/modules/sudo/manifests/init.pp:11"
  source: //sudo/File[/etc/sudoers]
```

A prettier method by graphical interface also exists, it's the Puppet Dashboard.

## Advanced usage

### Monitoring Passenger

There are binaries that allow you to get information about your Passenger and tune it if needed:

```
----------- General information -----------
max      = 12
count    = 0
active   = 0
inactive = 0
Waiting on global queue: 0

----------- Application groups -----------
```

And here's another piece of information:

```
--------- Apache processes ---------
PID   PPID  VMSize    Private  Name
------------------------------------
1841  1     91.9 MB   0.6 MB   /usr/sbin/apache2 -k start
1846  1841  91.1 MB   0.6 MB   /usr/sbin/apache2 -k start
1903  1841  308.0 MB  0.5 MB   /usr/sbin/apache2 -k start
1904  1841  308.0 MB  0.5 MB   /usr/sbin/apache2 -k start
### Processes: 4
### Total private dirty RSS: 2.14 MB


-------- Nginx processes --------

### Processes: 0
### Total private dirty RSS: 0.00 MB


---- Passenger processes ----
PID   VMSize   Private  Name
-----------------------------
1847  22.9 MB  0.3 MB   PassengerWatchdog
1852  31.6 MB  0.3 MB   PassengerHelperAgent
1857  43.1 MB  5.8 MB   Passenger spawn server
1860  79.9 MB  0.9 MB   PassengerLoggingAgent
### Processes: 4
### Total private dirty RSS: 7.26 MB
```

### Checking the syntax of your .pp files

When you create/edit a puppet module, it can be quickly very practical to check the syntax. Here's how:

```
puppet parser validate init.pp
```

And if you want to do it on a larger scale:

```
find /etc/puppet/ -name '*.pp' | xargs -n 1 -t puppet parser validate
```

### Overriding restrictions

If, for example, we've defined a class and for certain hosts we want to modify this configuration, we need to do it like this:

```ruby
class somehost_postfix inherits postfix {
    # blah blah blah
}

node somehost {
    include somehost_postfix
}
```

Let's say we have a defined postfix module. We want to apply a particular config to certain hosts defined here by 'somehost'. To bypass the configuration, you need to create a class 'somehost_postfix'. I insist here on the nomenclature of the name to give for this class since it's only like this that Puppet will recognize that you want to apply a particular configuration.

### Disabling a resource

To temporarily disable a resource, just set noop to true:

```ruby
file { "/etc/passwd":
    noop => true}
```

### Pre and Post puppet run

It's possible to run scripts before and after a Puppet run. This can be useful in the case of backing up certain files via etckeeper for example. Add this to your puppet.conf file on your clients:

```
[...]
prerun_command = /usr/local/bin/before-puppet-run.sh
postrun_command = /usr/local/bin/after-puppet-run.sh
```

### CFT

CFT (pronounced shift) is a small software that will look at what you do during a given period to generate a manifest for you. For example, you launch it just before doing an installation with its configuration and it will generate the manifest once completed. An example is better than a long speech:

```
cft begin apache
[...]
cft finish apache
cft manifest apache
```

### Generating a manifest from an existing system

Here's a simple solution to generate a manifest from an already installed system by specifying a resource. Here are some examples:

```
puppet resource user root
puppet resource service httpd
puppet resource package postfix
```

### Puppet Push

Puppet works in client -> server mode. It's the client that contacts the server every 30 minutes (by default) and asks for synchronization. When you're in dumb mode or when you want to trigger at a specific time the update of your client machines to the server, there's a special mode (listen). The problem is that the puppet client can't run in client and listen mode. So you need 2 instances... in short, the nightmares begin.

To overcome this problem, I've developed Puppet push which allows asking clients (via SSH) to synchronize. As you'll have understood, it's more than necessary to have a key exchange done beforehand with the clients. How to do it? No worries, we've seen that above.

To access the latest version, follow this link: http://www.deimos.fr/gitweb/?p=puppet_push.git;a=tree

### MCollective

MCollective is a bit of a gas plant, but it's very powerful and works very well with Puppet. It allows you like Puppet Push to do several actions on nodes, but uses a dedicated protocol, it doesn't need SSH to communicate with the nodes. To learn more, look at this article.

## FAQ

### err: Could not retrieve catalog from remote server: hostname was not match with the server certificate

You have a problem with your certificates. The best thing is to regenerate a certificate with all the server's hostnames in this certificate: Creating a certificate

## Resources

- http://reductivelabs.com/products/puppet/
- http://www.rottenbytes.info
- [Modules for Puppet](https://git.black.co.at/)
- [Puppet recipes](https://reductivelabs.com/trac/puppet/wiki/Recipes)
- [Types of objects for Puppet](https://reductivelabs.com/trac/puppet/wiki/TypeReference)
- [Puppet SSL Explained](https://www.masterzen.fr/2010/11/14/puppet-ssl-explained/) [(/pdf)](/pdf/puppet_ssl_explained.pdf)
- http://puppetcookbook.com/
