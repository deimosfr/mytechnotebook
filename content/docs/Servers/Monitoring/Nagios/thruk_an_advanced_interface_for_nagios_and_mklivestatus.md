---
weight: 999
url: "/Thruk_\\:_Une_interface_évoluée_pour_Nagios_et_MKlivestatus/"
title: "Thruk: An Advanced Interface for Nagios and MKlivestatus"
description: "Thruk is an advanced interface for Nagios and Shinken that allows connecting to multiple Nagios instances simultaneously and display more relevant information compared to the standard Nagios interface."
categories: ["Debian", "Linux", "Apache", "Monitoring", "Nagios"]
date: "2012-05-30T09:33:00+02:00"
lastmod: "2012-05-30T09:33:00+02:00"
tags: ["Thruk", "Nagios", "Monitoring", "MKlivestatus", "LDAP", "Apache"]
toc: true
---

![Thruk logo](/images/thruk_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.30 |
| **Operating System** | Debian 6 |
| **Website** | [Thruk Website](https://www.thruk.org) |
| **Last Update** | 30/05/2012 |
{{< /table >}}

## Introduction

Thruk is an advanced interface for Nagios (and Shinken) that allows connecting to multiple Nagios instances simultaneously and displaying more relevant information compared to the standard Nagios interface.

## Prerequisites

Here are the prerequisites:

```bash
aptitude install libapache2-mod-fcgid
```

## Installation

Let's download Thruk:

```bash
wget http://www.thruk.org/files/pkg/v1.30/debian6/amd64/thruk_1.30_debian6_amd64.deb
dpkg -i thruk_1.30_debian6_amd64.deb
```

You can now connect to the Thruk interface: http://nagios/thruk/  
Login and password: thrukadmin

## Configuration

First, we'll adjust some permissions to avoid future problems:

```bash
chown www-data. /etc/thruk/thruk* /etc/thruk/cgi.cfg
chown -Rf nagios. /var/lib/nagios3
touch /etc/thruk/checkconfig
chown nagios:www-data /etc/thruk/checkconfig
```

### Apache LDAP

You might want to use LDAP authentication with Thruk, just like your Nagios. It's quite simple:

```apache {linenos=table,hl_lines=["36-63"]}
<IfModule mod_fcgid.c>
  AddHandler fcgid-script .sh

  <Directory /usr/share/thruk>
    Options FollowSymLinks
    AllowOverride All
    order allow,deny
    allow from all
  </Directory>
  <Directory /etc/thruk/themes>
    Options FollowSymLinks
    allow from all
  </Directory>
  <Directory /etc/thruk/plugins>
    Options FollowSymLinks
    allow from all
  </Directory>

  # redirect to a startup page when there is no pidfile yet
  RewriteEngine On
  RewriteCond %{REQUEST_METHOD} GET
  RewriteCond %{REQUEST_URI} !^/thruk/startup.html
  RewriteCond %{REQUEST_URI} !^/thruk/side.html
  RewriteCond %{REQUEST_URI} !^/thruk/.*\.(css|png|js)
  RewriteCond %{REQUEST_URI} ^/thruk
  RewriteCond /var/cache/thruk/thruk.pid !-f
  RewriteRule ^(.*)$ /thruk/startup.html?$1 [R=302,L,NE,QSA]

  Alias /thruk/documentation.html /usr/share/thruk/root/thruk/documentation.html
  Alias /thruk/startup.html /usr/share/thruk/root/thruk/startup.html
  AliasMatch ^/thruk/(.*\.cgi|.*\.html)  /usr/share/thruk/fcgid_env.sh/thruk/$1
  AliasMatch ^/thruk/plugins/(.*?)/(.*)$  /etc/thruk/plugins/plugins-enabled/$1/root/$2
  Alias /thruk/themes/  /etc/thruk/themes/themes-enabled/
  Alias /thruk /usr/share/thruk/root/thruk

  <Location /thruk>
    Options ExecCGI
  </Location>

  <DirectoryMatch (/usr/share/thruk/root/thruk|/etc/thruk/themes/themes-enabled/)>
    Options FollowSymLinks
    #DirectoryIndex index.php

    AllowOverride AuthConfig
    Order Deny,Allow
    Deny From All
    Allow From 127.0.1.1
    Allow From 127.0.0.1

    AuthName "Thruk Access"
    AuthType Basic
    AuthBasicProvider ldap

    AuthzLDAPAuthoritative on
    AuthLDAPURL ldap://openldap.deimos.fr/dc=openldap,dc=deimos,dc=fr?uid?sub?(objectClass=posixAccount)
    AuthLDAPRemoteUserIsDN off
    AuthLDAPGroupAttribute memberUid
    AuthLDAPGroupAttributeIsDN off
    Require ldap-group ou=Groups,dc=openldap,dc=deimos,dc=fr
    Require ldap-user nagiosadmin

    Satisfy Any
  </DirectoryMatch>
</IfModule>
```

### Thruk

For Thruk configuration, there are 2 important files:

- thruk.conf: General configuration
- thruk_local.conf: Custom configuration

According to the official documentation, it's preferable not to touch the general configuration and override certain parameters with the custom configuration.

Here's a basic configuration to make it work with [MK Livestatus]({{< ref "docs/Servers/Monitoring/Nagios/check_mk_collect_nagios_info_and_extend_possibilities.md" >}}):

```text
############################################
# put your own settings into this file
# settings from this file will override
# those from the thruk.conf
############################################

# Interface Options
default_theme       = Exfoliation
use_timezone = CET
use_ajax_search = 1
use_new_search = 1
info_popup_event_type = onmouseover
show_modified_attributes = 1
show_long_plugin_output = inline
show_full_commandline = 2

# Backends to Mklivestatus
<Component Thruk::Backend>
    <peer>
        name   = Local Nagios
        type   = livestatus
        <options>
            peer          = /var/lib/nagios3/rw/live
            resource_file = /etc/nagios3/resource.cfg
       </options>
       <configtool>
            core_conf      = /etc/nagios3/nagios.cfg
            obj_check_cmd  = /usr/sbin/nagios3 -v /etc/nagios3/nagios.cfg
            obj_reload_cmd = /etc/init.d/nagios3 reload
       </configtool>
    </peer>
</Component>
```

You can add multiple Peers to have a unified interface with several Nagios servers. Once you have access to the Thruk interface, you can configure all its options directly from this interface, as well as those of Nagios.

Then reload Apache for the changes to take effect.

### Minimal Interface for Monitoring Screens

I had already discussed in [the Nagios documentation]({{< ref "docs/Servers/Monitoring/Nagios/nagios_installation_and_configuration.md" >}}) solutions for having a fairly minimal screen. The problem is that it's still not sufficient, and fortunately with Thruk there's a way to fix this issue without having to recode 3/4 of the program. So I created a small patch that I submitted to the Thruk team (which has just been accepted [https://github.com/sni/Thruk/commit/d1eefef82cd8fbab6ebebff8a570bb1e026d1a9f](https://github.com/sni/Thruk/commit/d1eefef82cd8fbab6ebebff8a570bb1e026d1a9f) but will only be available in the next release (1.28)), so in the meantime, here's how to have the most minimal interface possible.

Here is a first patch to modify the Thruk Perl modules so that they take into account a new parameter in the URL called "minimal":

```diff
diff -Nru old/Root.pm new/Root.pm
--- old/Root.pm 2012-04-06 13:05:08.000000000 +0200
+++ new/Root.pm 2012-04-06 13:05:16.000000000 +0200
@@ -215,6 +215,9 @@
     }
     $c->stash->{hidetop} = $c->{'request'}->{'parameters'}->{'hidetop'} || '';

+       # Add custom monitor screen function
+    $c->stash->{minimal} = $c->{'request'}->{'parameters'}->{'minimal'} || '';
+
     # initialize our backends
     unless ( defined $c->{'db'} ) {
         $c->{'db'} = $c->model('Thruk');
diff -Nru old/status.pm new/status.pm
--- old/status.pm       2012-04-06 13:05:08.000000000 +0200
+++ new/status.pm       2012-04-06 13:05:21.000000000 +0200
@@ -744,6 +744,9 @@
     $c->stash->{hidetop}    = 1 unless $c->stash->{hidetop} ne '';
     $c->stash->{hidesearch} = 1;

+    # Monitor screen interface
+    $c->stash->{minimal}    = 1 unless $c->stash->{minimal} ne '';
+
     # which host to display?
     my( $hostfilter)           = Thruk::Utils::Status::do_filter($c, 'hst_');
     my( undef, $servicefilter) = Thruk::Utils::Status::do_filter($c, 'svc_');
```

And then we will modify the interface so that the changes apply only to a part of the code:

```diff
diff -Nru old/_header.tt new/_header.tt
--- old/_header.tt      2012-04-06 13:39:31.000000000 +0200
+++ new/_header.tt      2012-04-06 13:37:44.000000000 +0200
@@ -192,14 +192,18 @@
         </tr>
       </table>
     </td>
+       [% UNLESS minimal == 1 %]
     <td>
       <input type="image" src="[% url_prefix %]thruk/themes/[% theme %]/images/arrow_refresh.png" class="top_refresh_button" onClick="refresh_button(this)" alt="refresh this page" title="refresh this page">
     </td>
     <td class="top_nav_pref">
       <input type="button" class="top_nav_pref_button" value="preferences" title="preferences" onMouseOver="button_over(this)" onMouseOut="button_out(this)" onClick="toggleElement('pref_pane'); return false;">
     </td>
+       [% END %]
   </tr>
 </table>
+[% UNLESS minimal == 1 %]
 [% IF page == 'status' || page == 'statusmap' %]
 <a href="#" onClick="toggleTopPane(); return false;"><img src="[% url_prefix %]thruk/themes/[% theme %]/images/icon_minimize.gif" class="btn_toggle_top_pane" id="btn_toggle_top_pane" alt="toggle header"></a>
 [% END %]
+[% END %]
diff -Nru old/status_detail.tt new/status_detail.tt
--- old/status_detail.tt        2012-04-06 13:37:10.000000000 +0200
+++ new/status_detail.tt        2012-04-06 13:37:24.000000000 +0200
@@ -12,6 +12,7 @@

 [% PROCESS _overdiv.tt %]
 [% PROCESS _status_cmd_pane.tt %]
+[% UNLESS minimal == 1 %]
     <table border="0" width="100%" cellspacing="0" cellpadding="0" id="top_pane"[% IF hidetop == 1 %]style="visibility:hidden; display:none;"[% END %]>
       <tr>
         <td align="left" valign="top" width="33%">
@@ -103,9 +104,11 @@
         </td>
       </tr>
     </table>
+       [% END %]

     [% PROCESS _status_detail_table.tt %]

+       [% UNLESS minimal == 1 %]
     [% UNLESS authorized_for_read_only %]
     <div class="hint">
         <a href="#" onclick="selectAllServices(true,'[% paneprefix %]');return false;">select all</a> (<a href="#" onclick="selectAllHosts(true,'[% paneprefix %]');">hosts</a>)
@@ -118,6 +121,7 @@
     <div align="center">[% PROCESS _pager.tt %]</div>

     <br>
+       [% END %]
     <div class='itemTotalsTitle'>[% IF !has_error && data.size %][% data.size %] of [% pager.total_entries %][% ELSE %]0[% END %] Matching Service Entries Displayed</div>

 [% PROCESS _footer.tt %]
```

And let's patch everything:

```bash
patch -d /usr/share/thruk/lib/Thruk/Controller/ < Thruk_perl_modules.patch
patch -d /usr/share/thruk/templates/ < Thruk_interface.patch
```

Now, all you have to do is add "&minimal=1" to the end of your thrunk URL to remove all unnecessary elements. Example:

```
http://nagios/thruk/cgi-bin/status.cgi?host=all&type=detail&hostprops=10&serviceprops=42&servicestatustypes=28&hidetop=1&minimal=1
```

### Adding a Custom CGI

In some cases, you may have certain checks that temporarily store information on the Nagios server, and you want to be able to execute actions from the Nagios interface. For this, there is the 'action_url' option where we can provide a URL to a CGI that will execute what we want, possibly with options.

Let's start by creating our CGI. Here is a minimalist example where I remove a temporary file:

```perl
#!/usr/bin/perl

use CGI;
$query = CGI::new();
$host = $query->param("host");

# Avoid inputing special characters that would crash the program
if ( $h =~ /\`|\~|\@|\#|\$|\%|\^|\&|\*|\(|\)|\:|\=|\+|\"|\'|\;|\<|\>/ ) {
    print "Illegal special chars detected. Exit\n";
    exit(1);
}

print "Content-type: text/html\n\n";
print "<HTML>\n";
print "<HEAD><Title>Removing $host temporary file</Title>\n";
print "<LINK REL='stylesheet' TYPE='text/css' HREF='/nagios/stylesheets/common.css'><LINK REL='stylesheet' TYPE='text/css' HREF='/nagios/stylesheets/status.css'>\n";
print "</HEAD><BODY>\n";
print "Removing $host Interface Network Flapping temporary file...";
if (-f "/tmp/iface_state_$host.txt")
{
    unlink("/tmp/iface_state_$host.txt") or print "FAIL<br />/tmp/iface_state_$host.txt : $!\n" and exit(1);
    print "OK\n";
}
else
{
    print "FAIL<br />/tmp/iface_state_$host.txt : No such file or directory\n";
}
print "</body></html>\n";
```

And then in the configuration of the service in question, I insert my 'action_url':

```text {linenos=table,hl_lines=["9-13"]}
define service{
         use                             generic-services-ulsysnet
         hostgroup_name                  network
         service_description             Interface Network Flapping
         check_period                    24x7
         notification_period             24x7
 	 _SNMP_PORT			 161
	 _SNMP_COMMUNITY		         public
         _DURATION			 86400
         check_command                   check_interface_flapping
         # For Thruk & Nagios
	 # action_url			 ../../cgi-bin/nagios3/remove.cgi?host=$HOSTADDRESS$
         # For Nagios only
         action_url			 remove.cgi?host=$HOSTADDRESS$
}
```

Now you just need to reload Nagios.

You'll notice that for Thruk, I found an easy but not very clean method, which consists of rewriting the URL to point to the Nagios links. For a clean method, you would need to write a dedicated plugin.

## FAQ

### OS Icons No Longer Display

Apparently in the Debian package of this version of Thruk, there is a small issue with displaying OS icons. To fix this, we'll create a symbolic link:

```bash
ln -s /usr/share/nagios/htdocs/images/logos/base /usr/share/thruk/themes/themes-available/Classic/images/logos/
```
