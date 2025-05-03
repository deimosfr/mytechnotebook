---
weight: 999
url: "/Page_Speed\\:_optimize_on_the_fly_your_rendered_code/"
title: "PageSpeed: Optimize Your Rendered Code On The Fly"
description: "Learn how to install and configure PageSpeed with Nginx to optimize web content and improve page load times."
categories: ["Nginx", "Debian", "Linux"]
date: "2014-02-19T17:14:00+02:00"
lastmod: "2014-02-19T17:14:00+02:00"
tags: ["PageSpeed", "Nginx", "Performance", "Optimization", "Web", "Google"]
toc: true
---

![PageSpeed](/images/pagespeed.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.7.30.3 |
| **Operating System** | Debian 7 |
| **Website** | [PageSpeed Website](https://developers.google.com/speed/pagespeed/module) |
| **Last Update** | 19/02/2014 |
| **Others** | Nginx 1.4.4-1 |
{{< /table >}}

## Introduction

PageSpeed[^1] speeds up your site and reduces page load time. This open-source webserver module automatically applies web performance best practices to pages and associated assets (CSS, JavaScript, images) without requiring that you modify your existing content or workflow.

## Installation

As Wheezy has an outdated version of Nginx, we're going to use Nginx extras package with PageSpeed built-in. To get it, we will use the well-known DotDeb repository.

To install DotDeb repository, add a preference file to avoid unwanted override packages:

```bash
Package: *
Pin: release o=packages.dotdeb.org
Pin-Priority: 100
```

Then add the repository:

```bash
deb http://packages.dotdeb.org wheezy all
deb-src http://packages.dotdeb.org wheezy all
```

Add the GPG key:

```bash
cd /tmp
wget http://www.dotdeb.org/dotdeb.gpg
sudo apt-key add dotdeb.gpg
```

And run an update:

```bash
apt-get update
```

Then you can install nginx with pagespeed integrated:

```bash
aptitude install nginx-extras
```

## Configuration

The configuration is quite easy but there are lots of options that need to be tested one by one to be sure they have the correct effect:

```bash
# PageSpeed
# Enable ngx_pagespeed
pagespeed on;
pagespeed FileCachePath /usr/share/nginx/pagespeed;

# Ensure requests for pagespeed optimized resources go to the pagespeed handler
# and no extraneous headers get set.
location ~ "\.pagespeed\.([a-z]\.)?[a-z]{2}\.[^.]{10}\.[^.]+\" {
  add_header "" "";
}
location ~ "^/ngx_pagespeed_static/" { }
location ~ "^/ngx_pagespeed_beacon$" { }
location /ngx_pagespeed_statistics { allow 127.0.0.1; deny all; }
location /ngx_pagespeed_global_statistics { allow 127.0.0.1; deny all; }
location /ngx_pagespeed_message { allow 127.0.0.1; deny all; }

# Defer and minify Javascript
pagespeed EnableFilters defer_javascript;
pagespeed EnableFilters rewrite_javascript;
pagespeed EnableFilters combine_javascript;
pagespeed EnableFilters canonicalize_javascript_libraries;

# Inline and minimize css
pagespeed EnableFilters rewrite_css;
pagespeed EnableFilters fallback_rewrite_css_urls;
# Loads CSS faster
#pagespeed EnableFilters move_css_above_scripts;
pagespeed EnableFilters move_css_to_head;

# Rewrite, resize and recompress images
pagespeed EnableFilters rewrite_images;

# remove tags with default attributes
pagespeed EnableFilters elide_attributes;

# To enable Varnish
pagespeed DownstreamCachePurgeLocationPrefix http://127.0.0.1:80/;
pagespeed DownstreamCachePurgeMethod PURGE;
pagespeed DownstreamCacheRewrittenPercentageThreshold 95;
```

Then apply this configuration to your desired virtual host:

```bash
server {
    listen 80;
    include pagespeed.conf;
[...]
```

And reload Nginx for changes to take effect. You can also enable pagespeed on the server side directly to apply the configuration to all vhosts at once.

### File path cache optimization

The file path cache is used to store rewritten elements. So it's preferable to have good performance on it, and to achieve this, you can use a tmpfs. Add this line to your fstab:

```bash
tmpfs                                       /usr/share/nginx/pagespeed  tmpfs   rw,mode=1777,size=512M  0   0
```

And then create the folder and mount it:

```bash
mkdir /usr/share/nginx/pagespeed
mount /usr/share/nginx/pagespeed
```

Restart nginx and it's ready!

## Benchmark

Several websites exist to benchmark your site and give you recommendations:

- GTMetrix: http://gtmetrix.com
- Page Speed Insight: http://developers.google.com/speed/pagespeed/insights

Then try to correlate with [PageSpeed options](https://developers.google.com/speed/pagespeed/optimization) and you're now ready to perform tests and see the awesome results!

## References

[^1]: https://developers.google.com/speed/
