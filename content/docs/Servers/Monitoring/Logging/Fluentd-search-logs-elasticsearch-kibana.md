---
weight: 999
url: "/Fluentd\\:_quickly_search_in_your_logs_with_Elasticsearch,_Kibana_and_Fluentd/"
title: "Fluentd: Quickly Search in Your Logs with Elasticsearch, Kibana and Fluentd"
description: "How to set up a log management and search solution using Fluentd, Elasticsearch, and Kibana for efficient log collection, processing, and visualization"
categories: ["Servers", "Monitoring", "Linux"]
date: "2014-05-10T08:20:00+02:00"
lastmod: "2014-05-10T08:20:00+02:00"
tags: ["fluentd", "elasticsearch", "kibana", "logs", "monitoring", "syslog"]
toc: true
---

![Fluentd](/images/fluentd_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.1.19-1 |
| **Operating System** | Debian 7 |
| **Website** | [Fluentd Website](https://fluentd.org/) |
| **Last Update** | 10/05/2014 |
| **Others** | [Elasticsearch](https://www.elasticsearch.org) 1.1<br />[Kibana](https://www.elasticsearch.org/overview/kibana/) 3.0.1 |
{{< /table >}}

## Introduction

Managing logs is not a complicated task with classical syslog systems (syslog-ng, rsyslog...). However, being able to search in them quickly when you have several gigabits of logs, with scalability, with a nice graphical interface etc...is not the same thing.

Fortunately today, tools that permit us to do it very well exist. Here are the list of tools that we're going to use to achieve it:

* **Elasticsearch**[^1]: Elasticsearch is a flexible and powerful open source, distributed, real-time search and analytics engine. Architected from the ground up for use in distributed environments where reliability and scalability are must-haves, Elasticsearch gives you the ability to move easily beyond simple full-text search. Through its robust set of APIs and query DSLs, plus clients for the most popular programming languages, Elasticsearch delivers on the near limitless promises of search technology
* **Kibana**[^2]: Kibana is Elasticsearch's data visualization engine, allowing you to natively interact with all your data in Elasticsearch via custom dashboards. Kibana's dynamic dashboard panels are savable, shareable and exportable, displaying changes to queries into Elasticsearch in real-time. You can perform data analysis in Kibana's beautiful user interface using pre-designed dashboards or update these dashboards in real-time for on-the-fly data analysis.
* **Fluentd**[^3]: Fluentd is an open source data collector designed for processing data streams.

All those elements will be installed on the same machine to make it simpler at start. Fluentd is an alternative to Logstash. They both are data collectors, however Fluentd permits sending logs to other destinations:

![Fluentd explanation](/images/fluentd_explain.avif)

Here is what kind of infrastructure you can setup (no redundancy here, just a single instance):

![Architecture overview](/images/es_ki_fl.avif)

To avoid dependencies issues and make things simpler, we're going to use fluentd as forwarder here to transfer syslog and other kinds of logs to another fluentd instance. On the last one, Elasticsearch and Kibana will be installed.

## Installation

### Elasticsearch

The first thing to put in place is the backend that will store our logs. As we want the latest version, we're going to use the dedicated repository:

```bash
cd /tmp
wget -O - http://packages.elasticsearch.org/GPG-KEY-elasticsearch | apt-key add -
echo "deb http://packages.elasticsearch.org/elasticsearch/1.1/debian stable main" > /etc/apt/sources.list.d/elasticsearch.list
```

Then we're ready to install:

```bash
aptitude update
aptitude install elasticsearch openjdk-7-jre-headless openntpd
```

To finish, configure the auto start of the service and run it:

```bash
update-rc.d elasticsearch defaults 95 10
/etc/init.d/elasticsearch start
```

### Kibana

Regarding Kibana, there is unfortunately no repository at the moment. So we're going to use the git repository to make it simpler. First of all, install a web server like Nginx:

```bash
aptitude install nginx git
```

Now clone the repository and use the latest version (here 3.0.1):

```bash
cd /usr/share/nginx/www
git clone https://github.com/elasticsearch/kibana.git
cd kibana
git checkout v3.0.1
```

You can get the list of all versions with **git tag** command.

You now need to configure Nginx to get it provided properly:

```text
server {
        listen   80;
        server_name  kibana.deimos.fr;

        root /usr/share/nginx/www/kibana/src/;
        index index.html;

        access_log /var/log/nginx/kibana.deimos.fr_access.log;
        error_log /var/log/nginx/kibana.deimos.fr_error.log; 

        location / {
            # First attempt to serve request as file, then
            # as directory, then fall back to displaying a 404.
            try_files $uri $uri/ /index.html;
            # Uncomment to enable naxsi on this location
            # include /etc/nginx/naxsi.rules
        }
 }
```

To finish for Kibana, edit the configuration file and adapt the elasticsearch line to your need:

```javascript
     * ==== elasticsearch
     *
     * The URL to your elasticsearch server. You almost certainly don't
     * want +http://localhost:9200+ here. Even if Kibana and Elasticsearch are on
     * the same host. By default this will attempt to reach ES at the same host you have
     * kibana installed on. You probably want to set it to the FQDN of your
     * elasticsearch host
     *
     * Note: this can also be an object if you want to pass options to the http client. For example:
     *
     *  +elasticsearch: {server: "http://localhost:9200", withCredentials: true}+
     *
     */
    elasticsearch: "http://<elasticserver_dns_name>:9200",
```

Restart Nginx service to make the web interface available to http://<kibana_dns_name>:

![Kibana Interface](/images/kibana.avif)

### Fluentd

Fluentd is now the last part that will permit sending syslog to another Fluentd or Elasticsearch. So this has to be done on all Fluentd forwarders or servers.

First of all, we'll adjust system parameters to be sure we won't face performance issues due to it. First, edit the security limits and add these lines:

```text
root soft nofile 65536
root hard nofile 65536
* soft nofile 65536
* hard nofile 65536
```

Then we're going to add the sysctl tuning in that file:

```text
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.ip_local_port_range = 10240 65535
```

And apply the new configuration:

```bash
sysctl -p
```

We're going to add the official repository:

```bash
wget http://packages.treasure-data.com/debian/RPM-GPG-KEY-td-agent
apt-key add RPM-GPG-KEY-td-agent
echo 'deb http://packages.treasure-data.com/debian/ lucid contrib' > /etc/apt/sources.list.d/fluentd.list
```

However, during the time I'm writing this documentation, there are no Wheezy version available (squeeze only) and there is a missing dependency on the libssl. We're going to get it from squeeze and install it:

```bash
wget http://ftp.fr.debian.org/debian/pool/main/o/openssl/libssl0.9.8_0.9.8o-4squeeze14_amd64.deb
dpkg -i libssl0.9.8_0.9.8o-4squeeze14_amd64.deb
```

We're now ready to install Fluentd agent:

```bash
aptitude update
aptitude install td-agent openntpd
mkdir /etc/td-agent/config.d
```

Modify then the configuration to set the global configuration:

```apache
## match tag=debug.** and dump to console
<match debug.**>
  type stdout
</match>

# HTTP input
# POST http://localhost:8888/<tag>?json=<json>
# POST http://localhost:8888/td.myapp.login?json={"user"%3A"me"}
# @see http://docs.fluentd.org/articles/in_http
<source>
  type http
  port 8888
</source>

## live debugging agent
<source>
  type debug_agent
  bind 127.0.0.1
  port 24230
</source>

# glob match pattern
include config.d/*.conf
```

Restart td-agent service.

#### Elasticsearch plugin

By default, it doesn't know how to forward to Elasticsearch. So we will need to install a dedicated plugin for it on the server, **not on the forwarders**. Here is how to install it:

```bash
> aptitude install build-essential ruby-dev libcurl4-openssl-dev make
> /usr/lib/fluent/ruby/bin/fluent-gem install fluent-plugin-elasticsearch
Building native extensions.  This could take a while...
Fetching: multi_json-1.10.0.gem (100%)
Fetching: multipart-post-2.0.0.gem (100%)
Fetching: faraday-0.9.0.gem (100%)
Fetching: elasticsearch-transport-0.4.11.gem (100%)
Fetching: elasticsearch-api-0.4.11.gem (100%)
Fetching: elasticsearch-0.4.11.gem (100%)
Fetching: fluent-plugin-elasticsearch-0.3.0.gem (100%)
Successfully installed patron-0.4.18
Successfully installed multi_json-1.10.0
Successfully installed multipart-post-2.0.0
Successfully installed faraday-0.9.0
Successfully installed elasticsearch-transport-0.4.11
Successfully installed elasticsearch-api-0.4.11
Successfully installed elasticsearch-0.4.11
Successfully installed fluent-plugin-elasticsearch-0.3.0
8 gems installed
Installing ri documentation for patron-0.4.18...
Installing ri documentation for multi_json-1.10.0...
Installing ri documentation for multipart-post-2.0.0...
Installing ri documentation for faraday-0.9.0...
Installing ri documentation for elasticsearch-transport-0.4.11...
Installing ri documentation for elasticsearch-api-0.4.11...
Installing ri documentation for elasticsearch-0.4.11...
Installing ri documentation for fluent-plugin-elasticsearch-0.3.0...
Installing RDoc documentation for patron-0.4.18...
Installing RDoc documentation for multi_json-1.10.0...
Installing RDoc documentation for multipart-post-2.0.0...
Installing RDoc documentation for faraday-0.9.0...
Installing RDoc documentation for elasticsearch-transport-0.4.11...
Installing RDoc documentation for elasticsearch-api-0.4.11...
Installing RDoc documentation for elasticsearch-0.4.11...
Installing RDoc documentation for fluent-plugin-elasticsearch-0.3.0...
```

## Configuration

Here you will see how to configure multiple options of Fluentd. Choose the ones you want to add to your Fluentd instances (can have several). Here is a good example of what is needed in this kind of configuration:

![Fluentd architecture example](/images/fluentd_archi_example.avif)

### Forwarders

To make a Fluentd forwards data to a receiver, simply create that configuration file and set the Fluentd node to forward to:

```apache
<match **>
  type forward
  <server>
    host fluentd.deimos.fr
    port 24224
  </server>
</match>
```

### Receiver

If you want your node to be able to receive data from other Fluentd forwarders, you need to add this configuration:

```apache
## built-in TCP input
## @see http://docs.fluentd.org/articles/in_forward
<source>
  type forward
</source>
```

In that use case, you need to add this on the server role of Fluentd.

### Rsyslog

By default, Debian is using Rsyslog and we're going to see here how to forward syslog to Fluentd. First of all, on the Fluentd forwarders, create a syslog file containing the configuration as follows:

```apache
<source>
  type syslog
  port 5140
  bind 127.0.0.1
  tag syslog
</source>
```

And restart td-agent service. It will create a listening port for Syslog.

Then simply add this line to redirect (in addition to the local files) syslog to Fluentd:

```text
*.*          @127.0.0.1:5140
```

Restart Rsyslog service.

### Log files

You may want to be able to log files as well. Here is a way to do it for a single access file from Nginx logs:

```apache
<source>
  type tail
  # Select the file to watch
  path /var/log/nginx/access.log
  # Select a file to store offset position
  pos_file /tmp/td-agent.nginx.pos
  # Format type
  format syslog
  # Tag with a distinguish name
  tag system.nginx
</source>
```

Then restart the td-agent service.

#### Nginx

The problem of the basic example above is each element is passed on a single line. That means we can't filter accurately. To do it, you will need to split with regex each field and give them a field name. You also need to specify the time and date format. Here is how to do it for Nginx:

```apache
<source>
  type tail
  path /var/log/nginx/access.log
  pos_file /tmp/td-agent.nginx.pos
  format syslog
  tag nginx.access
  # Regex fields
  format /^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*) "(?<referer>[^\"]*)" "(?<agent>[^\"]*)"$/
  # Date and time format
  time_format %d/%b/%Y:%H:%M:%S %z 
</source>
```

It may be complicated to create a working regex the first time. That's why a website called Fluentular (http://fluentular.herokuapp.com) can help you to create the format line.

### Fluentd Elasticsearch

To send all incoming sources to Elasticsearch, simply create that configuration file:

```apache
<match **>
  type elasticsearch
  logstash_format true
  host localhost
  port 9200
</match>
```

Then restart the td-agent service.

## Usage

If you look at the web interface, you should have something like this:

![Kibana interface with data](/images/kibana2.avif)

You can now try to add other widgets, look at [the official documentation](https://www.elasticsearch.org/guide/en/kibana/current/)[^4].

## References

[^1]: http://www.elasticsearch.org
[^2]: http://www.elasticsearch.org/overview/kibana/
[^3]: http://fluentd.org/
[^4]: http://www.elasticsearch.org/guide/en/kibana/current/

---

* http://jasonwilder.com/blog/2013/11/19/fluentd-vs-logstash/
* http://repeatedly.github.io/2014/02/analyze-event-logs-using-fluentd-and-elasticsearch/
* http://www.devconsole.info/?p=917
* http://lifeandshell.com/install-elasticsearch-kibana-fluentd-opensource-splunk-with-syslog-clients/
* https://github.com/fluent/fluentd
