---
weight: 999
url: "/Proxychains_\\:_proxyfier_n'importe_quelle_connexion_vers_l'extérieur/"
title: "Proxychains: Proxy Any Outbound Connection"
description: "Guide to setup and use Proxychains to route any application's traffic through a proxy without modifying the application itself."
categories: ["Linux", "Networking", "Security"]
date: "2012-05-14T09:53:00+02:00"
lastmod: "2012-05-14T09:53:00+02:00"
tags: ["Proxy", "SSH", "Network", "Security", "Debian", "Mac OS X"]
toc: true
---

## Introduction

Ah, those proxy servers at work or school! They can be quite troublesome! Sometimes you really need temporary access to the outside world. Depending on the commands you use, you may or may not have the ability to configure proxy usage by modifying a configuration file, an environment variable, etc.

The benefit of [proxychains](https://proxychains.sourceforge.net/) is that configuration is done only once, in its own configuration file. Then you use the syntax `proxychains <command> <args>` and your command will use the proxy specified in the proxychains configuration file!

![Proxychains](/images/proxychains.avif)

## Installation

### Debian

```bash
aptitude install proxychains
```

### Mac OS X

On Mac OS X, it's not yet available in MacPorts, so we'll need to patch the source to make it work on Mac. Copy this diff to your machine:

```diff
diff -ruN proxychains-3.1/proxychains/Makefile.in proxychains-3.1_resolv/proxychains/Makefile.in
--- proxychains-3.1/proxychains/Makefile.in	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/Makefile.in	2011-06-16 13:17:20.000000000 -0500
@@ -121,7 +121,7 @@
 LIBS = @LIBS@
 libproxychains_la_DEPENDENCIES =
 libproxychains_la_OBJECTS =  libproxychains.lo core.lo
-CFLAGS = @CFLAGS@
+CFLAGS = @CFLAGS@ -arch x86_64 -arch i386
 COMPILE = $(CC) $(DEFS) $(INCLUDES) $(AM_CPPFLAGS) $(CPPFLAGS) $(AM_CFLAGS) $(CFLAGS)
 LTCOMPILE = $(LIBTOOL) --mode=compile $(CC) $(DEFS) $(INCLUDES) $(AM_CPPFLAGS) $(CPPFLAGS) $(AM_CFLAGS) $(CFLAGS)
 CCLD = $(CC)
diff -ruN proxychains-3.1/proxychains/core.c proxychains-3.1_resolv/proxychains/core.c
--- proxychains-3.1/proxychains/core.c	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/core.c	2011-06-16 13:17:19.000000000 -0500
@@ -35,12 +35,18 @@
 #include <fcntl.h>
 #include <time.h>
 #include <stdarg.h>
+#include <dlfcn.h>
 #include "core.h"

 extern int tcp_read_time_out;
 extern int tcp_connect_time_out;
 extern int proxychains_quiet_mode;
-
+extern connect_t true_connect;
+extern getaddrinfo_t true_getaddrinfo;
+extern freeaddrinfo_t true_freeaddrinfo;
+extern getnameinfo_t true_getnameinfo;
+extern gethostbyaddr_t true_gethostbyaddr;
+
 static const char base64[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

 static void encode_base_64(char* src,char* dest,int max_len)
@@ -159,13 +165,14 @@

 	pfd[0].fd=sock;
 	pfd[0].events=POLLOUT;
-	fcntl(sock, F_SETFL, O_NONBLOCK);
+	fcntl(sock, F_SETFL, O_NONBLOCK);
   	ret=true_connect(sock, addr,  len);
-//	printf("\nconnect ret=%d\n",ret);fflush(stdout);
+//	printf("\nconnect ret=%d\n",ret); fflush(stdout);
+
   	if(ret==-1 && errno==EINPROGRESS)
    	{
     		ret=poll(pfd,1,tcp_connect_time_out);
-//      		printf("\npoll ret=%d\n",ret);fflush(stdout);
+//	     		printf("\npoll ret=%d\n",ret);fflush(stdout);
       		if(ret==1)
         	{
            		value_len=sizeof(int);
@@ -388,14 +395,18 @@
 				inet_ntoa(*(struct in_addr*)&pd->ip),
 				htons(pd->port));
 	pd->ps=PLAY_STATE;
+
 	bzero(&addr,sizeof(addr));
+
 	addr.sin_family = AF_INET;
 	addr.sin_addr.s_addr = pd->ip;
 	addr.sin_port = pd->port;
+

 	if (timed_connect (*fd ,(struct sockaddr*)&addr,sizeof(addr))) {
 		pd->ps=DOWN_STATE;
 		goto error1;
 	}
+
 	pd->ps=BUSY_STATE;
 	return SUCCESS;
 error1:
@@ -641,7 +652,7 @@
 			dup2(pipe_fd[1],1);
 			//dup2(pipe_fd[1],2);
 		//	putenv("LD_PRELOAD=");
-			execlp("proxyresolv","proxyresolv",name,NULL);
+			execlp("./proxyresolv","proxyresolv",name,NULL);
 			perror("can't exec proxyresolv");
 			exit(2);

diff -ruN proxychains-3.1/proxychains/core.h proxychains-3.1_resolv/proxychains/core.h
--- proxychains-3.1/proxychains/core.h	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/core.h	2011-06-16 13:17:19.000000000 -0500
@@ -66,29 +66,28 @@
 int proxychains_write_log(char *str,...);
 struct hostent* proxy_gethostbyname(const char *name);

+typedef struct hostent* (*gethostbyname_t)(const char *);
+static gethostbyname_t true_gethostbyname;

 typedef int (*connect_t)(int, const struct sockaddr *, socklen_t);
-connect_t true_connect;
-
-typedef struct hostent* (*gethostbyname_t)(const char *);
-gethostbyname_t true_gethostbyname;
+// connect_t true_connect;

 typedef int (*getaddrinfo_t)(const char *, const char *,
 		const struct addrinfo *,
 		struct addrinfo **);
-getaddrinfo_t true_getaddrinfo;
+// getaddrinfo_t true_getaddrinfo;

 typedef int (*freeaddrinfo_t)(struct addrinfo *);
-freeaddrinfo_t true_freeaddrinfo;
+// freeaddrinfo_t true_freeaddrinfo;

 typedef int (*getnameinfo_t) (const struct sockaddr *,
 		socklen_t, char *,
 		socklen_t, char *,
 		socklen_t, unsigned int);
-getnameinfo_t true_getnameinfo;
+// getnameinfo_t true_getnameinfo;

 typedef struct hostent *(*gethostbyaddr_t) (const void *, socklen_t, int);
-gethostbyaddr_t true_gethostbyaddr;
+// gethostbyaddr_t true_gethostbyaddr;

 int proxy_getaddrinfo(const char *node, const char *service,
 		                const struct addrinfo *hints,
diff -ruN proxychains-3.1/proxychains/libproxychains.c proxychains-3.1_resolv/proxychains/libproxychains.c
--- proxychains-3.1/proxychains/libproxychains.c	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/libproxychains.c	2011-06-16 13:17:19.000000000 -0500
@@ -32,7 +32,6 @@
 #include <sys/fcntl.h>
 #include <dlfcn.h>

-
 #include "core.h"

 #define     satosin(x)      ((struct sockaddr_in *) &(x))
@@ -57,6 +56,13 @@
 	unsigned int *proxy_count,
 	chain_type *ct);

+connect_t true_connect;
+getaddrinfo_t true_getaddrinfo;
+freeaddrinfo_t true_freeaddrinfo;
+getnameinfo_t true_getnameinfo;
+gethostbyaddr_t true_gethostbyaddr;
+
+
 static void init_lib()
 {
 //	proxychains_write_log("ProxyChains-"VERSION
@@ -291,7 +297,7 @@
 int getnameinfo (const struct sockaddr * sa,
 			socklen_t salen, char * host,
 			socklen_t hostlen, char * serv,
-			socklen_t servlen, unsigned int flags)
+			socklen_t servlen, int flags)
 {
 	int ret = 0;
 	if(!init_l)
diff -ruN proxychains-3.1/proxychains/proxychains proxychains-3.1_resolv/proxychains/proxychains
--- proxychains-3.1/proxychains/proxychains	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/proxychains	2011-06-16 13:17:20.000000000 -0500
@@ -1,9 +1,11 @@
 #!/bin/sh
 echo "ProxyChains-3.1 (http://proxychains.sf.net)"
+echo "Mod for OSX - using dylib"
 if [ $# = 0 ] ; then
 	echo "	usage:"
 	echo "		proxychains <prog> [args]"
 	exit
 fi
-export LD_PRELOAD=libproxychains.so
+export DYLD_FORCE_FLAT_NAMESPACE=
+export DYLD_INSERT_LIBRARIES=./.libs/libproxychains.3.0.0.dylib
 exec "$@"
diff -ruN proxychains-3.1/proxychains/proxyresolv proxychains-3.1_resolv/proxychains/proxyresolv
--- proxychains-3.1/proxychains/proxyresolv	2006-03-15 10:16:59.000000000 -0600
+++ proxychains-3.1_resolv/proxychains/proxyresolv	2011-06-16 13:18:51.000000000 -0500
@@ -11,6 +11,6 @@
 	exit
 fi

-
-export LD_PRELOAD=libproxychains.so
-dig $1 @$DNS_SERVER +tcp | awk '/A.+[0-9]+\.[0-9]+\.[0-9]/{print $5;}'
+export DYLD_FORCE_FLAT_NAMESPACE=
+export DYLD_INSERT_LIBRARIES=./.libs/libproxychains.3.0.0.dylib
+dig $1 @$DNS_SERVER +tcp | awk '/^[^;].+A.+[0-9]+\.[0-9]+\.[0-9]/{print $5;}'
```

Now let's compile it:

```bash
wget -O proxychains-3.1.tar.gz "http://prdownloads.sourceforge.net/proxychains/proxychains-3.1.tar.gz?download"
tar -xzf proxychains-3.1.tar.gz
patch -p0 < proxychains-3.1_osx.diff
cd proxychains-3.1
./configure
make
cd proxychains
sudo cp proxychains proxyresolv /usr/sbin/
```

## Configuration

You can use your personal configuration file in '~/proxychains/proxychains.conf' or the one for the entire machine in '/etc/proxychains':

```bash {linenos=table,hl_lines=[1,3,5,11]}
strict_chain
# Quiet mode (no output from library)
quiet_mode
# Proxy DNS requests - no leak for DNS data
#proxy_dns
# Some timeouts in milliseconds
tcp_read_time_out 15000
tcp_connect_time_out 8000

[ProxyList]
socks5 127.0.0.1 12345
```

Use the port for the SOCKS proxy that we'll open next via SSH.

## Usage

To get an application that doesn't have SOCKS proxy options to access the outside world using TCP protocol, you first need to set up a SOCKS proxy using SSH:

To establish an SSH connection that opens a SOCKS proxy, run this from server A:

```bash
ssh -D <port> <user>@<destination>
```

For example:

```bash
ssh -D 12345 user@serverB
```

Then, to route an application's traffic through this proxy, use the proxychains command:

```bash
proxychains my_application
```

## Resources
- http://proxychains.sourceforge.net/
- http://lesdatabases.blogspot.com/2011/06/proxychains-ou-l-de-proxifier.html
- http://chrootlabs.org/bgt/proxychains_osx.html
