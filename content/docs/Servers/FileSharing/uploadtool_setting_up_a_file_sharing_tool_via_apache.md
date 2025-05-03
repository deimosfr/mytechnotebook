---
weight: 999
url: "/UploadTool_\\:_Mise_en_place_d'un_outil_d'échange_de_fichiers_via_Apache/"
title: "UploadTool: Setting up a file sharing tool via Apache"
description: "A guide to install and configure UploadTool, a web-based file sharing tool that uses Apache for authenticated file uploads."
categories: ["Linux", "Apache"]
date: "2007-09-04T08:22:00+02:00"
lastmod: "2007-09-04T08:22:00+02:00"
tags: ["upload", "file sharing", "Apache", "PHP", "web application"]
toc: true
---

## Introduction

I was looking for a software for the company I work for, that allows file transfers with authentication but just at the upload level. So I found [UploadTool](https://uploadtool.sourceforge.net/).

## Installation

Download the [upload.tar.gz 1.0](https://belnet.dl.sourceforge.net/sourceforge/uploadtool/upload.tar.gz) archive and extract it:

```bash
wget http://belnet.dl.sourceforge.net/sourceforge/uploadtool/upload.tar.gz
tar -xzvf upload.tar.gz
```

Next, we'll move everything to our Apache directory, then assign the correct permissions:

```bash
mv upload /var/www
chown -Rf www-data. /var/www/upload
```

## Configuration

To configure it, it's quite simple, just go to the URL of your site followed by upload. Example: [https://www.mydomain.com/upload](https://www.mydomain.com/upload).

Create your root account and then the users who will have access to upload.

### Protections for public folders

To protect yourself from listing uploaded files, I suggest a small [HTML redirector]({{< ref "docs/Servers/Web/Apache/apache_2_installation_and_configuration.md#redirecteur-html" >}}).

#### Enhancement patch

```bash
*** bin/common.php      2006-05-08 02:23:34.000000000 +0200
--- bin/common.php      2007-09-03 11:46:22.000000000 +0200
***************
*** 59,67 ****
--- 59,70 ----
                echo '<td class="report">File Name</td>';
                echo '<td class="report">Size</td>';
                echo '<td class="report">Date</td>';
+               echo '<td class="report">URL to give</td>';
                echo '</tr>';
                foreach ($files as $filename)
                {
+                       if ($filename !== "index.html")
+                       {
                        $url = filepath2url($cwd . "/" . $filename);
                        echo "<tr>";
                        echo "<td>";
***************
*** 78,84 ****
--- 81,91 ----
                        echo "<td>";
                        echo date ("Y-m-d H:i:s", filemtime($cwd . "/" . $filename));
                        echo "</td>";
+                       echo "<td>";
+                       echo "$url";
+                       echo "</td>";
                        echo "</tr>\n";
+                       }
                }
                echo '</table>';
        }
```

Here's a patch I created to improve the interface a bit. To apply it, create an "upload.patch" file that you put in your upload folder, then run this command:

```bash
patch -p0 < upload.patch
```

Admire the result :-)

### Upload size limitations

Here's what you need to add in your Apache VirtualHost to limit the size of uploaded files:

```bash
    <Directory />
      php_value max_execution_time 300
      php_value upload_max_filesize 40M
      php_value post_max_size 40M
    </Directory>
    LimitRequestBody 40000000
```
