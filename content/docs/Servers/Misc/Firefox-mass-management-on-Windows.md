---
weight: 999
url: "/Firefox_mass_management_on_Windows/"
title: "Firefox: Mass Management on Windows Environment"
description: "Learn how to manage Firefox deployments in a Windows Active Directory environment using GPO for Firefox."
categories: ["Windows", "Browser"]
date: "2013-01-22T06:40:00+02:00"
lastmod: "2013-01-22T06:40:00+02:00"
tags: ["Firefox", "GPO", "Windows", "Active Directory"]
toc: true
---

![Firefox](/images/firefox_icon.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 17+ |
| **Operating System** | Windows 2003 RC2<br />Windows 7 |
| **Website** | [Firefox Website](https://www.mozilla.org) |
| **Last Update** | 22/01/2013 |
| **Others** |  |
{{< /table >}}

## Introduction

You probably manage your Internet Explorer configuration via GPO under Active Directory. This is convenient, there are plenty of options, but logically there is nothing by default for alternative browsers.

The first issue you will encounter is the lack of MSI by default. You have 2 options:

1. Use [SCCM](https://fr.wikipedia.org/wiki/System_Center_Configuration_Manager)[^1] (formerly SMS) to deploy .exe files
2. Use [repackaged versions in MSI format](https://www.frontmotion.com/Firefox/download_firefox.htm)[^2] and deploy them via GPO.

We will use the first option here because SCCM allows you to do much more than classic GPO deployments. However, we won't cover how to set up deployment via SCCM or GPO as tutorials are plentiful on the internet, and we'll focus on the Firefox part.

We will use the Firefox plugin called [GPO for Firefox](https://addons.mozilla.org/en-us/firefox/addon/gpo-for-firefox/)[^3]. This plugin allows Firefox, when launched, to look at registry properties (which will have been pushed by GPO) to force default values in Firefox's "about:config" and make them non-modifiable by the user.

It's possible to go a bit further with Firefox repackaging by implementing additional default features via [the CCK plugin](https://addons.mozilla.org/fr/firefox/addon/cck/)[^4].

## Prerequisites

For the prerequisites, we'll need at minimum:

* An Active Directory + DNS server
* An SCCM server
* A Windows 7 client

## Installation

### Firefox

Add the Firefox extension "[GPO for Firefox](https://addons.mozilla.org/fr/firefox/addon/gpo-for-firefox/)"[^5] by repackaging it if possible so that it's installed by default.

### GPO on AD

We'll install the GPOs on the server where Active Directory is located. First, download the file containing the new GPO entries for Firefox:
[https://sourceforge.net/projects/gpofirefox/files/firefox.adm/download](https://sourceforge.net/projects/gpofirefox/files/firefox.adm/download)[^6]

Then launch the Users and Computers interface:
![Firefox AD](/images/firefox_ad.avif)

We'll install the GPO at the domain level, but you can also do it at the GPO level if you wish. Go to domain properties:
![Firefox AD properties](/images/firefox_ad_properties.avif)

Here's the procedure to follow:
![Firefox create GPO](/images/firefox_create_gpo.avif)

1. Click on the GPO tab (Group Policy)
2. Add a new GPO
3. Name it in a recognizable way
4. Click Edit to edit it

Now we'll add the firefox.adm file to access the new options. Position yourself on "Administrative Templates" and click on "Add/Remove Templates":
![Firefox add template](/images/firefox_add_template.avif)

Add the firefox.adm file:
![Firefox add firefox adm](/images/firefox_add_firefox_adm.avif)

Then you should now see it appear:
![Firefox adm added](/images/firefox_adm_added.avif)

## Configuration

### Configuration of the GPO

Configuring the GPO is quite simple once you understand the principle. It's possible to force certain Firefox properties from the GPOs on the user or computer side. You can choose what works best for you. Depending on what you choose, the user-side parameters that will be in the registry can be found here:

* HKLM or HKCU\Software\Policies\Mozilla\LockPref for locked preferences
* HKLM or HKCU\Software\Policies\Mozilla\defaultPref for default preferences

One of the first parameters to apply is the following: activating the GPOFirefox module:

![Firefox template added](/images/firefox_template_added.avif)

We activate it and hide it from the list of classic applications:

![Firefox set GPO](/images/firefox_set_gpo.avif)

You have the ability to access most of the functions present in "about:config" in the "Mozilla Advanced Options" folder.

## References

[^1]: http://fr.wikipedia.org/wiki/System_Center_Configuration_Manager
[^2]: http://www.frontmotion.com/Firefox/download_firefox.htm
[^3]: https://addons.mozilla.org/en-us/firefox/addon/gpo-for-firefox/
[^4]: https://addons.mozilla.org/fr/firefox/addon/cck
[^5]: https://addons.mozilla.org/fr/firefox/addon/gpo-for-firefox/
[^6]: http://sourceforge.net/projects/gpofirefox/files/firefox.adm/download
