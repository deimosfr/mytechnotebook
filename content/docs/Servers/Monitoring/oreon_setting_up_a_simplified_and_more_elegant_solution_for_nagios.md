---
weight: 999
url: "/Oreon_\\:_Mise_en_place_d'une_solution_simplifiée_et_plus_élégante_pour_Nagios/"
title: "Oreon: Setting up a simplified and more elegant solution for Nagios"
description: "Overview of Oreon (now Centreon), a monitoring and network supervision software based on Nagios with an improved interface and additional features."
categories: ["Linux", "Servers", "Network"]
date: "2009-05-11T12:43:00+02:00"
lastmod: "2009-05-11T12:43:00+02:00"
tags: ["Nagios", "Oreon", "Centreon", "Monitoring", "Network", "Supervision"]
toc: true
---

## Introduction

Oreon is a network monitoring and supervision software based on the most efficient Open Source information retrieval engine: Nagios.

Today Oreon no longer exists and is called Centreon.

The objective of this project is to provide a new interface for Nagios, capable of giving it new functionalities while keeping the logic of its existing mechanisms. All within a modern, scalable interface designed for everyone.

The Oreon project benefits from a constantly growing community and unwavering investment, which allows it to make its solution increasingly robust and constantly offer innovative new tools.

You can find on the website, www.oreon-project.org, all the necessary information for integrating Oreon into your infrastructure:

* Packages to download
* Necessary documentation
* Appropriate community tools
* Professional support service

## Installation of Nagios

To install Nagios, refer to the docs in this wiki, or read this documentation written specifically for Oreon integration.  
[documentation](/pdf/nagiosoreonfr.pdf)

## Installation of Oreon

You can follow this documentation.  
[documentation](/pdf/oreon13fr.pdf)

However, during the installation, I encountered some dependency issues. That's why at the end of it, before proceeding to the web installation, enter these lines:

```bash
pear install -o -f --alldeps Mail Mail_Mime Net_SMTP Net_Socket Net_Traceroute Net_Ping Validate Image_Graph Image_GraphViz HTML_Table HTML_QuickForm_advmultiselect Auth_SASL HTTP Numbers_Roman Numbers_Words MDB2 DB_DataObject_FormBuilder DB_DataObject DB Date
```
