---
weight: 999
url: "/Rancid_Search_\\:_Mise_en_place_d'un_outil_de_recherche_pour_Rancid/"
title: "Rancid Search: Setting Up a Search Tool for Rancid"
description: "How to implement a search tool for Rancid to easily search through network device configurations."
categories: ["Linux", "Network"]
date: "2011-03-02T06:16:00+02:00"
lastmod: "2011-03-02T06:16:00+02:00"
tags: ["rancid", "search", "network", "cisco", "configuration", "web", "perl"]
toc: true
---

## Introduction

You may already be familiar with [Rancid](https://www.shrubbery.net/rancid/) which allows you to backup your precious Cisco equipment and store them in a [VCS]({{< ref "docs/Servers/Versionning/">}}). At work, my network team complained about not being able to search through all Cisco equipment at once. Imagine multiple devices with some containing more than 25,000 lines of ACLs and other VPN configurations. To search for information, they first had to know which equipment to connect to, then search through the configuration. In short, it was not always an easy task to read through, time-consuming, and without regex search capability.

I therefore decided to create a web interface to meet their needs with the help of a colleague. The interface is in Perl CGI where I use jQuery for some "flashy" or "Web 2.0" features, but it looks nice :-)
