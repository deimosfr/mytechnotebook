---
weight: 999
url: "/Jenkins_\\:_Mise_en_place_d'un_outil_d'intégration_continue/"
title: "Jenkins: Setting up a continuous integration tool"
description: "Learn how to install and configure Jenkins, an open source continuous integration tool, with Nginx as a reverse proxy."
categories: ["Server", "Development", "Continuous Integration"]
date: "2013-04-12T09:10:00+02:00"
lastmod: "2013-04-12T09:10:00+02:00"
tags: ["jenkins", "nginx", "continuous integration", "development", "debian"]
toc: true
---

![Jenkins](/images/jenkins_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.447 |
| **Operating System** | Debian 7 |
| **Website** | [Jenkins](https://jenkins-ci.org/) |
| **Last Update** | 12/04/2013 |
{{< /table >}}

## Introduction

Jenkins is an open source continuous integration tool, forked from the Hudson tool after disagreements between its author, Kohsuke Kawaguchi, and Oracle. Written in Java, Jenkins runs in a servlet container such as Apache Tomcat, or standalone with its own embedded web server.

It interfaces with version control systems such as CVS and Subversion, and executes projects based on Apache Ant and Apache Maven as well as arbitrary scripts in Unix shell or Windows batch.

Project builds can be initiated in various ways, such as cron-like scheduling mechanisms, dependency systems between builds, or through requests to specific URLs.

Recently, Jenkins has become a popular alternative to the reference tool CruiseControl.

On January 11, 2011, a proposal to rename Hudson was announced to avoid problems with a possible trademark registration of the name by Oracle. After failed negotiations with Oracle, a vote in favor of renaming was ratified on January 29, 2011.

Here we'll see how to set up a Jenkins server that will control a Selenium server for unit testing on a PHP application like [Limesurvey](/Limesurvey_:_Mise_en_place_d'une_solution_de_Sondages/).

## Installation

For the installation part, it's easy:

```bash
aptitude install jenkins nginx
```

Note that we're using an Nginx server as a frontend to forward requests to Jenkins and absorb the load of requests.

## Configuration

### Nginx

For Nginx, we will use the reverse proxy function to redirect the flow:

```bash {linenos=table,hl_lines=[2,8,15]}
upstream app_server {
    server 127.0.0.1:8080 fail_timeout=0;
}

server {
    listen 80;
    listen [::]:80 default ipv6only=on;
    server_name jenkins.deimos.fr;
    location / {
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_redirect off;

        if (!-f $request_filename) {
            proxy_pass http://127.0.0.1:8080;
            break;
        }
    }
}
```

Adapt these lines according to your configuration.

Then we activate this by default:

```bash
rm -f /etc/nginx/sites-enabled/default
ln -s /etc/nginx/sites-available/jenkins /etc/nginx/sites-enabled/
```

Then restart Nginx to access via http://jenkins-server
