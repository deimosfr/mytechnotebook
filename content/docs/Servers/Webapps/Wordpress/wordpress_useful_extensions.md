---
weight: 999
url: "/Wordpress_\\:_les_extentions_pratiques/"
title: "WordPress: Useful Extensions"
description: "A collection of useful WordPress plugins and extensions that can enhance your website's functionality."
categories: ["CMS", "Web Development"]
date: "2010-04-12T20:35:00+02:00"
lastmod: "2010-04-12T20:35:00+02:00"
tags: ["WordPress", "Plugins", "Web Development", "CMS"]
toc: true
---

## Introduction

WordPress is, in my opinion, THE ultimate blogging platform. Like Firefox, it really shines with its plugins. Here's my list of plugins that I currently use or have used in the past that I find interesting.

## Tips

### Adding File Extensions for Upload

You may have noticed that you can't upload just anything to WordPress, which can be frustrating. By looking into the source code, you can add more extensions. For example, I needed to add the ogv format:

```php
/**
 * Retrieve list of allowed mime types and file extensions.
 *
 * @since 2.8.6
 *
 * @return array Array of mime types keyed by the file extension regex corresponding to those types.
 */
function get_allowed_mime_types() {
    static $mimes = false;

    if ( !$mimes ) {
        // Accepted MIME types are set here as PCRE unless provided.
        $mimes = apply_filters( 'upload_mimes', array(
        'jpg|jpeg|jpe' => 'image/jpeg',
        'gif' => 'image/gif',
        'png' => 'image/png',
        'bmp' => 'image/bmp',
        'tif|tiff' => 'image/tiff',
        'ico' => 'image/x-icon',
        'asf|asx|wax|wmv|wmx' => 'video/asf',
        'avi' => 'video/avi',
        'divx' => 'video/divx',
        'flv' => 'video/x-flv',
        'ogv' => 'video/ogg',
...
```

## The Extensions

### iWPhone

This extension automatically resizes the site to the correct format (iPhone) when accessing the WordPress site. It's very convenient because everything is automatic, with no need to modify the pages :-)

http://iwphone.contentrobot.com/

PS: Don't forget to install the WordPress plugin on your iPhone via the App Store, which will allow you to post messages very easily.

### Contact Form 7

This handy plugin allows you to create forms easily. If you encounter an error like "The database table for Contact Form 7 does not exist", simply temporarily add all rights, long enough for it to create the table, then reset everything.

### NextGen Gallery

This plugin is truly amazing. Flash slideshows, 3D photo walls with Cooliris, etc. It's perfect for displaying photos and creating impressive slideshows.

#### Adding Flash to the Header

Modify your header file for your current theme and add this:

```php
<?php
$showgallery = '[slideshow=1]';
$showgallery = apply_filters('the_content', $showgallery);
echo $showgallery;
?>
```

You obviously need to have a photo gallery to display and set the corresponding number for the slideshow value.
