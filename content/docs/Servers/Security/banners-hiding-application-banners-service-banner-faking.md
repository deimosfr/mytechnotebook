---
weight: 999
url: "/Banières_\\:_Cacher_les_banières_de_ses_applications_(Service_banner_faking)/"
title: "Banners: Hiding Application Banners (Service Banner Faking)"
description: "A guide to hiding service banners and modifying version information to improve security by making it harder for attackers to fingerprint your services."
categories: ["Security", "Server", "Linux", "Windows"]
date: "2009-04-17T20:56:00+02:00"
lastmod: "2009-04-17T20:56:00+02:00"
tags: ["security", "banner", "hiding", "apache", "openssh", "proftpd", "postfix", "lighttpd"]
toc: true
---

## Introduction

This is a quick howto on faking service banners. Service banners often contain a lot of useful information for malicious script-kiddies, like the (real) running software on the remote host and its version number. Knowing this, they can better target their exploits. This howto deals with changing this information. Keep in mind that this won't make your system more secure against a known exploit when you run a vulnerable service, however it can provide some 'social engineering' security: script-kiddies often scan whole IP blocks for a known vulnerability, and only attack those who give back a banner telling that they run the vulnerable service. This howto aims to fake the service banners and in this way, fool the script-kiddies. However, your system will still be vulnerable to an exploit if you're running a vulnerable service! If a script-kiddie runs his exploit, even if he sees you don't send out the right banner, you can still be hacked. So, always keep your system up-to-date, see this as a way to decrease the amount of effective attacks on your system, not as a way to be invulnerable. Of course there's also the fun-factor: it's quite amusing to see script-kiddies attempt to break into your 'Microsoft-IIS/5.0' also known as Apache 1.3.27 *grin*.

In this howto we're going to hide some known services with banners from some other known (but worse) services. We've got five services running: ftp, ssh, smtp, http and pop3. Currently these services are running on: Proftpd 1.2.7, OpenSSH_3.5p1, Postfix 2.0.6, Apache/1.3.27 and Teapop 0.3.5. We're going to 'transform' these services to: Microsoft FTP Service, OpenSSH, Microsoft ESMTP MAIL Service, Microsoft-IIS/5.0 and Microsoft Exchange 2000 POP3 server. Of course we could had changed them to anything, but for the fun of it, we'll change it to Microsoft.

Probably, the service programmers don't want users to change the service banners. The only reason I can come up with is their ego. Statistics collectors on the Internet (example: Netcraft.com collecting HTTP/HEAD information), count the number of machines running service X. For programmers it's really a boost to see how many people are using their software. As long as the software is released under the GNU/GPL, you're completely free to modify anything from the source, and you're even allowed to re-distribute the changes. If you want to keep the programmers on the friendly side, you could change the banner to only advertise which software runs, not the version number or other information.

**Please note:**  
You are completely responsible for your own actions. I can never be held responsible for any damage this HOWTO has done to you, your systems or your life. This works for me, however that doesn't guarantee this will work for you.

**Terminology:**  
FQDN = Fully Qualified Domain Name (hostname.domain.tld)
Hostname = First part of the FQDN (example: localhost)
Text between the < and > should be replaced with the corresponding values.

## Services

### ProFTPd

Current banner:
```
220 ProFTPD 1.2.7 Server (FTP for: ) []
```

Banner lay-out: "response_code product_name product_version Server (ServerName) [hostname]"
Wanted banner:
```
220 Microsoft FTP Service (Version 5.0).
```

Howto:
Open src/main.c and search for
```
if((id = find_config(server->conf,CONF_PARAM,"ServerIdent",FALSE))
```

Comment:
```
(/* if-block */)
```

the whole if-block and add the following line under the if-block:
```
send_response("220", "%s", server->ServerName);
```

Now re-compile proftpd. After compiling edit proftpd.conf and change the "ServerName" directive to " Microsoft FTP Service (Version 5.0).".

### OpenSSH

Current banner:
```
SSH-2.0-OpenSSH_3.5p1
```

Banner lay-out: SSH-version-OpenSSH_version
Wanted banner:
```
SSH-2.0-OpenSSH
```

Howto:
Open version.h and cut the "_3.5p1" from the end:
```
My_SSH_version
```

Re-compile and it's done.

### Postfix

**Current banner:**
```
220  ESMTP Ready and Serving.
Banner lay-out: "response_code hostname ESMTP additional_information"
```

Wanted banner:
```
220 fire.deimos.fr Microsoft ESMTP MAIL Service, Version: 6.0.3790.1830 ready at Tue, 1 Apr 2008 16:57:50 +0200
```

**Howto:**

* Open postfix's main.cf (configuration file) and add this line:
```
smtpd_banner = fire.deimos.fr Microsoft ESMTP MAIL Service, Version: 6.0.3790.1830
```

* In the Postfix source edit the file src/global/mail_date.c and search the line:
```
#define STRFTIME_FMT "%a, %d %b %Y %H:%M:%S "
```

and replace with:
```
#define STRFTIME_FMT "%a,%d %b %Y %H:%M:%S "
```

There is only one space to delete.

Stay in this file and search this line:
```
while (strftime(vstring_end(vp), vstring_avail(vp), " (%Z)", lt) == 0)
```

and replace with:
```
while (strftime(vstring_end(vp), vstring_avail(vp), " ", lt) == 0)
```

This is for deleting the (EST) at the end of the line.

* Now to finish, open src/smtpd/smtpd.c and search for the line:
```
smtpd_chat_reply(state, "220 %s", var_smtpd_banner);
```

and replace it with this line:
```
smtpd_chat_reply(state, "220 %s ready at %s", var_smtpd_banner, mail_date(time((time_t *) 0)));
```

* Now recompile, restart postfix and you're done :-)

### Apache

#### Method 1

Current banner:
```
Apache/1.3.27 (Unix) mod_perl/1.25 PHP/4.2.3
```

Banner lay-out: "BASEPRODUCT/BASEREVISION (OS) Apache modules"  
Wanted banner:
```
Microsoft-IIS/5.0
```

Howto:
Open /src/include/httpd.h and search for:
```
#define SERVER_BASEVENDOR "Apache Group"
#define SERVER_BASEPRODUCT "Apache"
#define SERVER_BASEREVISION ""
```

Change this to the desired values:
```
define SERVER_BASEVENDOR: Microsoft
define SERVER_BASEPRODUCT: Microsoft-IIS
define SERVER_BASEREVISION: 5.0
```

Now re-compile apache.

You can continue with the second method to have a full banner faking.

#### Method 2

Open your httpd.conf (or apache2.conf) and search those directive. If it's not there, add it. Set ServerTokens to Min:
```
ServerSignature Off
ServerTokens Prod
```

More information about the ServerTokens directive is at: http://carnagepro.com/pub/Docs/Apache2/mod/core.html#servertokens.

### Teapop

Current banner:
```
+OK Teapop [v0.3.5] - Teaspoon stirs around again <1048009854.3FB15180@Llywellyn>
```

Banner lay-out: "POP_OK Teapop [version] - banner - "  
Wanted banner:
```
+OK Microsoft Exchange 2000 POP3 server version 6.0.6249.0 () ready.
```

Howto: The file /teapop/pop_hello.c contains the following line:
```
pop_socket_send(pinfo->out, "%s Teapop [v%s] - %s %s", POP_OK, POP_VERSION, POP_BANNER, pinfo->apopstr);"
```

Change this line to:
```
pop_socket_send(pinfo->out, "%s Microsoft Exchange 2000 POP3 server version 6.0.6249.0 () ready.", POP_OK);"
```

Now re-compile and it's done.

### Telnet (on Solaris)

The default banner displayed during a telnet login contains the Solaris version which can be useful to a potential attacker.

Create a plain text file called "/etc/default/telnetd" which contains a line such as:
```
BANNER="Unauthorized access prohibited\n\n"
```

The \n characters encode blank lines.

### Lighttpd

On Lighttpd, it's very easy, simply add this line to your configuration file and restart the server:
```
server.tag = "Apache/1.3.29 (Unix) mod_perl/1.29 PHP/4.4.1 mod_ssl/2.8.16 OpenSSL/0.9.7g"
```
