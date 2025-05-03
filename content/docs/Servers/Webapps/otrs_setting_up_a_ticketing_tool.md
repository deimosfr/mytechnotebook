---
weight: 999
url: "/OTRS_\\:_mise_en_place_d'un_outil_de_ticketing/"
title: "OTRS: Setting up a ticketing tool"
description: "Guide for installing and configuring OTRS ticketing system on Debian 6 with PostgreSQL and LDAP integration"
categories: ["Debian", "Storage", "Database"]
date: "2012-10-31T15:12:00+02:00"
lastmod: "2012-10-31T15:12:00+02:00"
tags: ["OTRS", "Ticketing", "PostgreSQL", "LDAP", "Apache", "API"]
toc: true
---

![OTRS](/images/otrs_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.4.14 |
| **Operating System** | Debian 6 |
| **Website** | [OTRS Website](https://www.otrs.com) |
| **Last Update** | 31/10/2012 |
{{< /table >}}

## Introduction

Open-source Ticket Request System (OTRS, literally "open source ticket request system"), is an open source software for customer relationship management or support service management. A company, organization or institution can use it to assign "tickets" to requests made via the help desk or troubleshooting service. This system facilitates the processing of support or troubleshooting requests and any requests made by phone or email. OTRS is distributed under GNU Affero General Public License.[^1]

For the implementation of this version, I used the official documentation[^2] and made some small adjustments.

## Prerequisites

Check that your locales are correctly defined:

```bash
dpkg-reconfigure locales
```

## Installation

We will need these packages:

```bash
aptitude install libapache2-mod-perl2 libdbd-pg-perl libnet-dns-perl libnet-ldap-perl libio-socket-ssl-perl libpdf-api2-perl libsoap-lite-perl libgd-text-perl libgd-graph-perl libapache-dbi-perl postgresql aspell aspell-en dbconfig-common dictionaries-common javascript-common libalgorithm-diff-perl libalgorithm-diff-xs-perl libaspell15 libauthen-sasl-perl libbit-vector-perl libcarp-clan-perl libcrypt-passwdmd5-perl libdate-pcalc-perl libemail-valid-perl libio-socket-inet6-perl libjs-prototype libjs-yui libmail-pop3client-perl libnet-domain-tld-perl libnet-imap-simple-perl libnet-imap-simple-ssl-perl libnet-smtp-ssl-perl libsocket6-per libtext-csv-perl libtext-csv-xs-perl libtext-diff-perl libxml-feedpp-perl libxml-treepp-perl procmail wwwconfig-common
```

Note that I'm using PostgreSQL and not MySQL. Next, we install OTRS (I deliberately didn't take the latest version, but you can do it without any worries):

```bash
useradd -r -d /opt/otrs/ -c 'OTRS user' otrs
usermod -g www-data otrs
cd /opt
wget http://ftp.otrs.org/pub/otrs/otrs-2.4.14.tar.gz
tar -xzf otrs-2.4.14.tar.gz
rm -f otrs-2.4.14.tar.gz
mv otrs-* otrs && cd otrs
cp Kernel/Config.pm.dist Kernel/Config.pm
cp Kernel/Config/GenericAgent.pm.dist Kernel/Config/GenericAgent.pm
perl bin/SetPermissions.pl --otrs-user=otrs --otrs-group=otrs --web-user=www-data --web-group=www-data /opt/otrs
```

## Configuration

### Apache

We'll create this configuration in Apache:

```apache
# --
# added for OTRS (http://otrs.org/)
# $Id: apache2-httpd-new.include.conf,v 1.5 2008/11/10 11:08:55 ub Exp $
# --

# agent, admin and customer frontend
ScriptAlias /otrs/ "/opt/otrs/bin/cgi-bin/"
Alias /otrs-web/ "/opt/otrs/var/httpd/htdocs/"

# if mod_perl is used
<IfModule mod_perl.c>

    # load all otrs modules
    Perlrequire /opt/otrs/scripts/apache2-perl-startup.pl

    # Apache::Reload - Reload Perl Modules when Changed on Disk
    PerlModule Apache2::Reload
    PerlInitHandler Apache2::Reload
    PerlModule Apache2::RequestRec

    # set mod_perl2 options
    <Location /otrs>
        ErrorDocument 403 /otrs/index.pl
        ErrorDocument 404 /otrs/index.pl
        SetHandler  perl-script
        PerlResponseHandler ModPerl::Registry
        Options +ExecCGI +FollowSymLinks
        PerlOptions +ParseHeaders
        PerlOptions +SetupEnv
        Order allow,deny
        Allow from all
       #<IfModule mod_rewrite.c>
       #   RewriteEngine On
       #   RewriteCond /usr/share/otrs/var/httpd/htdocs/maintenance.html -l
       #   RewriteRule ^.*$ /otrs-web/maintenance.html
       #</IfModule>
    </Location>

</IfModule>

# directory settings
<Directory "/opt/otrs/bin/cgi-bin/">
    AllowOverride None
    Options +ExecCGI -Includes
    Order allow,deny
    Allow from all
</Directory>
<Directory "/opt/otrs/var/httpd/htdocs/">
    AllowOverride None
    Order allow,deny
    Allow from all
</Directory>
```

Then we restart Apache:

```bash
service apache2 restart
```

### PostgreSQL

We'll configure the authentication part:

```bash {linenos=table,hl_lines=["7-9"]}
[...]
# Database administrative login by UNIX sockets
local   all         postgres                          ident

# TYPE  DATABASE    USER        CIDR-ADDRESS          METHOD

# OTRS
local   otrs    otrs    md5
host    otrs    otrs    127.0.0.1/32    md5
# "local" is for Unix domain socket connections only
local   all         all                               ident
# IPv4 local connections:
host    all         all         127.0.0.1/32          md5
# IPv6 local connections:
host    all         all         ::1/128               md5
```

And create users, databases, grant access:

```bash
su postgres
psql
create user otrs password 'otrs' nosuperuser;
create database otrs owner otrs;
```

Then we restart everything to ensure the new configuration is active:

```bash
service postgresql restart
```

And finally, we import the schemas:

```bash
psql -U otrs -W -f scripts/database/otrs-schema.postgresql.sql otrs
psql -U otrs -W -f scripts/database/otrs-initial_insert.postgresql.sql otrs
psql -U otrs -W -f scripts/database/otrs-schema-post.postgresql.sql otrs
```

### OTRS

All that's left is to configure OTRS. Here I've added LDAP:

```perl
# --
# Kernel/Config.pm - Config file for OTRS kernel
# Copyright (C) 2001-2009 OTRS AG, http://otrs.org/
# --
# $Id: Config.pm.dist,v 1.21 2009/02/16 12:01:43 tr Exp $
# --
# This software comes with ABSOLUTELY NO WARRANTY. For details, see
# the enclosed file COPYING for license information (AGPL). If you
# did not receive this file, see http://www.gnu.org/licenses/agpl.txt.
# --
#  Note:
#
#  -->> OTRS does have a lot of config settings. For more settings
#       (Notifications, Ticket::ViewAccelerator, Ticket::NumberGenerator,
#       LDAP, PostMaster, Session, Preferences, ...) see
#       Kernel/Config/Defaults.pm and copy your wanted lines into "this"
#       config file. This file will not be changed on update!
#
# --

package Kernel::Config;

sub Load {
    my $Self = shift;
    # ---------------------------------------------------- #
    # ---------------------------------------------------- #
    #                                                      #
    #         Start of your own config options!!!          #
    #                                                      #
    # ---------------------------------------------------- #
    # ---------------------------------------------------- #

    # ---------------------------------------------------- #
    # database settings                                    #
    # ---------------------------------------------------- #
    # DatabaseHost
    # (The database host.)
    $Self->{DatabaseHost} = 'localhost';
    # Database
    # (The database name.)
    $Self->{Database} = 'otrs';
    # DatabaseUser
    # (The database user.)
    $Self->{DatabaseUser} = 'otrs';
    # DatabasePw
    # (The password of database user. You also can use bin/CryptPassword.pl
    # for crypted passwords.)
    $Self->{DatabasePw} = 'otrs';
    # DatabaseDSN
    # (The database DSN for MySQL ==> more: "man DBD::mysql")
    #$Self->{DatabaseDSN} = "DBI:mysql:database=$Self->{Database};host=$Self->{DatabaseHost};";

    # (The database DSN for PostgreSQL ==> more: "man DBD::Pg")
    # if you want to use a local socket connection
    $Self->{DatabaseDSN} = "DBI:Pg:dbname=$Self->{Database};";
    # if you want to use a tcpip connection
#    $Self->{DatabaseDSN} = "DBI:Pg:dbname=$Self->{Database};host=$Self->{DatabaseHost};";

    # ---------------------------------------------------- #
    # fs root directory
    # ---------------------------------------------------- #
    $Self->{Home} = '/opt/otrs';

    # ---------------------------------------------------- #
    # insert your own config settings "here"               #
    # config settings taken from Kernel/Config/Defaults.pm #
    # ---------------------------------------------------- #
    # $Self->{SessionUseCookie} = 0;
    # $Self->{CheckMXRecord} = 0;

    # ---------------------------------------------------- #

    # ---------------------------------------------------- #
    # data inserted by installer                           #
    # ---------------------------------------------------- #
    # $DIBI$

    ##############
    # AGENT PART #
    ##############
    # Allows to define the Agent view (accessible via index.pl)

    $Self->{SecureMode} = 1;
    $Self->{FQDN} = 'ldap.deimos.fr';
    $Self->{AdminEmail} = 'xxx@mycompany.com';
    $Self->{ProductName} = 'Deimosfr';
    $Self->{DefaultLanguage} = 'fr';
    $Self->{DefaultCharset} = 'utf-8';
    $Self->{'LogModule::SysLog::Charset'} = 'utf-8';
    $Self->{'AuthModule::LDAP::Charset'} = 'utf-8';
    $Self->{'AuthSyncModule::LDAP::Charset'} = 'utf-8';

    $Self->{'AuthModule'} = 'Kernel::System::Auth::LDAP';
    $Self->{'AuthModule::LDAP::Host'} = 'ldap.deimos.fr';
    $Self->{'AuthModule::LDAP::BaseDN'} = 'ou=users,dc=deimos,dc=fr';
    $Self->{'AuthModule::LDAP::UID'} = 'uid';
    $Self->{'AuthModule::LDAP::UserLowerCase'} = 1;
    $Self->{'AuthModule::LDAP::Params'} = {
        port    => 389,
        timeout => 120,
        async   => 0,
        version => 3,
    };

    #############
    # SYNC PART #
    #############
    # Sync allows to synchronize customer information in the database once authentication is successful

    $Self->{'AuthSyncModule'} = 'Kernel::System::Auth::Sync::LDAP';
    $Self->{'AuthSyncModule::LDAP::Host'} = 'ldap.deimos.fr';
    $Self->{'AuthSyncModule::LDAP::BaseDN'} = 'ou=users,dc=deimos,dc=fr';
    $Self->{'AuthSyncModule::LDAP::UID'} = 'uid';
    $Self->{'AuthSyncModule::LDAP::UserSyncMap'} = {
        # DB -> LDAP
        UserFirstname => 'givenName',
        UserLastname  => 'sn',
        UserEmail     => 'mail',
    };
    $Self->{'AuthSyncModule::LDAP::Params'} = {
        port    => 389,
        timeout => 120,
        async   => 0,
        version => 3,
    };
    $Self->{'AuthSyncModule::LDAP::UserSyncInitialGroups'} = [
        'users',
    ];

    #################
    # CUSTOMER PART #
    #################
    # Allows to define the Customer view (accessible via customer.pl)

    $Self->{'Customer::AuthModule'} = 'Kernel::System::CustomerAuth::LDAP';
    $Self->{'Customer::AuthModule::LDAP::Host'} = 'localhost';
    $Self->{'Customer::AuthModule::LDAP::BaseDN'} = 'dc=deimos,dc=fr';
    $Self->{'Customer::AuthModule::LDAP::UID'} = 'uid';

    $Self->{CustomerUser} = {
        Name => 'LDAP Backend',
        Module => 'Kernel::System::CustomerUser::LDAP',
        Params => {
            # ldap host
            Host => 'localhost',
            # ldap base dn
            BaseDN => 'dc=deimos,dc=fr',
            # search scope (one|sub)
            SSCOPE => 'sub',
            SourceCharset => 'utf-8',
            DestCharset => 'iso-8859-1',
            Die => 0,
            # Net::LDAP new params (if needed - for more info see perldoc Net::LDAP)
            Params => {
                port    => 389,
                timeout => 120,
                async   => 0,
                version => 3,
            },
        },
        # customer uniq id
        CustomerKey => 'uid',
        # customer #
        CustomerID => 'mail',
        CustomerUserListFields => ['cn', 'mail'],
        CustomerUserSearchFields => ['uid', 'cn', 'mail'],
        CustomerUserSearchPrefix => '',
        CustomerUserSearchSuffix => '*',
        CustomerUserSearchListLimit => 250,
        CustomerUserPostMasterSearchFields => ['mail'],
        CustomerUserNameFields => ['givenname', 'sn'],
        # show not own tickets in customer panel, CompanyTickets
        CustomerUserExcludePrimaryCustomerID => 0,
        # add a ldap filter for valid users (expert setting)
        # CustomerUserValidFilter => '(!(description=gesperrt))',
        # admin can't change customer preferences
        AdminSetPreferences => 0,
        # cache time to life in sec. - cache any ldap queris
        CacheTTL => 0,
        Map => [
            # note: Login, Email and CustomerID needed!
            # var, frontend, storage, shown (1=always,2=lite), required, storage-type, http-link, readonly
            [ 'UserFirstname',  'Firstname',  'givenname',       1, 1, 'var', '', 0 ],
            [ 'UserLastname',   'Lastname',   'sn',              1, 1, 'var', '', 0 ],
            [ 'UserLogin',      'Username',   'uid',             1, 1, 'var', '', 0 ],
            [ 'UserEmail',      'Email',      'mail',            1, 1, 'var', '', 0 ],
       #    [ 'UserCustomerID', 'CustomerID', 'mail',            0, 1, 'var', '', 0 ],
            [ 'UserCustomerID', 'CustomerID', 'o',               0, 1, 'var', '', 0 ],
            [ 'UserCustomerIDs', 'CustomerIDs', 'customerids',   1, 0, 'var', '', 0 ],
            [ 'UserPhone',      'Phone',      'telephoneNumber', 1, 0, 'var', '', 0 ],
        ],
    };

    $Self->{'SendmailModule'} = 'Kernel::System::Email::SMTP';
    $Self->{'SendmailModule::Host'} = 'localhost';
    $Self->{'SendmailModule::Port'} = '25';
    $Self->{'CustomerGroupSupport'} = '1';

    # ---------------------------------------------------- #
    # ---------------------------------------------------- #
    #                                                      #
    #           End of your own config options!!!          #
    #                                                      #
    # ---------------------------------------------------- #
    # ---------------------------------------------------- #
}

# ---------------------------------------------------- #
# needed system stuff (don't edit this)                #
# ---------------------------------------------------- #
use strict;
use warnings;

use vars qw(@ISA $VERSION);
use Kernel::Config::Defaults;
push (@ISA, 'Kernel::Config::Defaults');

use vars qw(@ISA $VERSION);
$VERSION = qw($Revision: 1.21 $)[1];

# -----------------------------------------------------#
```

Make sure to configure the PostgreSQL part and comment out the MySQL part, as well as fill in all the fields correctly. We'll also use Apache::DBI to boost performance. Edit this file to uncomment these lines:

```bash {linenos=table,hl_lines=["2-3","6-8","11-12"]}
[...]
use Apache::DBI ();
Apache::DBI->connect_on_init('DBI:mysql:otrs', 'otrs', 'password');
use DBI ();

# enable this if you use mysql
#use DBD::mysql ();
#use Kernel::System::DB::mysql;

# enable this if you use postgresql
use DBD::Pg ();
use Kernel::System::DB::postgresql;
[...]
```

All that's left is to initialize the crontab:

```bash
/opt/otrs/bin/Cron.sh start otrs
```

Your OTRS server is accessible via: http://server-otrs/otrs :-)

## API

You might want to use the API for various reasons. Here's an example for creating a Ticket and an Article (used for ticket responses):

```perl
#!/usr/bin/perl

use strict;
use warnings;

# Load lib folders from OTRS directory
use lib "/opt/otrs";

# Load modules
use Kernel::Config;
use Kernel::System::Encode;
use Kernel::System::Log;
use Kernel::System::Time;
use Kernel::System::Main;
use Kernel::System::DB;
use Kernel::System::Ticket;
use DBI;
use Encode ;

# Load objects
my $ConfigObject = Kernel::Config->new();
my $EncodeObject = Kernel::System::Encode->new( ConfigObject => $ConfigObject, );

my $LogObject = Kernel::System::Log->new(
    ConfigObject => $ConfigObject,
    EncodeObject => $EncodeObject,
);

my $TimeObject = Kernel::System::Time->new(
    ConfigObject => $ConfigObject,
    LogObject    => $LogObject,
);

my $MainObject = Kernel::System::Main->new(
    ConfigObject => $ConfigObject,
    EncodeObject => $EncodeObject,
    LogObject    => $LogObject,
);

my $DBObject = Kernel::System::DB->new(
    ConfigObject => $ConfigObject,
    EncodeObject => $EncodeObject,
    LogObject    => $LogObject,
    MainObject   => $MainObject,
);

my $TicketObject = Kernel::System::Ticket->new(
    ConfigObject => $ConfigObject,
    LogObject    => $LogObject,
    DBObject     => $DBObject,
    MainObject   => $MainObject,
    TimeObject   => $TimeObject,
    EncodeObject => $EncodeObject,
);

# Encoding to avoir errors
my $subject = "My subject:-)";
Encode::encode( 'UTF-8', $subject ) or die("Cannot encode subject into UTF-8");

# Create Ticket
my $TicketID = $TicketObject->TicketCreate(
    Title			=> $subject, # now encoded in UTF-8
    CustomerID		        => 'pmavro',
    Queue			=> 'queue',
    Lock			=> 'unlock',
    Priority		        => '1 demande d\'information',
    State			=> 'new',
    Type			=> 'default',
    CustomerUser	        => 'pmavro',
    OwnerID			=> 1,
    UserID			=> 1,
);

# Create Article
my $ArticleID = $TicketObject->ArticleCreate(
    TicketID    => $TicketID,
    ArticleType => 'note-internal',    # email-external|email-internal|phone|fax|...
    SenderType  => 'agent',                           # agent|system|customer4
    From        => 'Some Agent <email@example.com>',  # not required but useful
    To          => 'Some Customer A <customer-a@example.com>',  # not required but useful
    Subject     => $subject,           # required
    Body        => 'the message text',                 # required
    Charset     => 'UTF-8',
    MimeType    => 'text/plain',
    HistoryType => 'OwnerUpdate',
    HistoryComment => 'Some free text!',
    UserID         => 1,
);
```

## References

[^1]: http://fr.wikipedia.org/wiki/OTRS
[^2]: http://wiki.otterhub.org/index.php?title=Installation_on_Debian_6_with_Postgres
