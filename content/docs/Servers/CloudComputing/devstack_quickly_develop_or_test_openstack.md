---
weight: 999
url: "/DevStack_\\:_développer_ou_tester_rapidement_OpenStack/"
title: "DevStack: Quickly Develop or Test OpenStack"
description: "Learn how to use DevStack to quickly set up and test OpenStack for development and testing purposes"
categories: ["Ubuntu", "Linux"]
date: "2013-04-14T07:21:00+02:00"
lastmod: "2013-04-14T07:21:00+02:00"
tags: ["Servers", "OpenStack", "DevStack", "Development", "Testing"]
toc: true
---

|                  |                                          |
| ---------------- | ---------------------------------------- |
| Software version | Grizzly                                  |
| Operating System | Ubuntu Server 12.04                      |
| Website          | [DevStack Website](https://devstack.org/) |
| Last Update      | 14/04/2013                               |

## Introduction

If you want to quickly test [OpenStack](https://www.openstack.org/)[^1], whether for playing around or developing with it, [DevStack](https://devstack.org/)[^2] is currently the fastest method.

I started with an Ubuntu Server 12.04 base as this distribution is one of the versions recommended by DevStack. The goal here is to have a VM with all OpenStack services installed and functional.

{{< alert context="warning" text="<b>WARNING</b>: Do NOT use DevStack in production!!!" />}}

## Installation

To set up DevStack, we'll need git:

```bash
aptitude install git git-core
```

Now let's get the current version of DevStack, and switch to the desired version (Grizzly):

```bash
git clone git://github.com/openstack-dev/devstack.git
cd devstack
git checkout stable/grizzly
```

Then launch the installation and provide all requested passwords:

```bash
./stack.sh
```

## Utilization

To launch a development stack:

```bash
./stack.sh
```

When the installation is complete, you can access the different services like this:

- Horizon: http://server/
- Keystone: http://server:5000/v2.0/

The default users are admin and demo

For more information, check the [GitHub](https://github.com/openstack-dev/devstack/tree/stable/grizzly)[^3] page for DevStack.

## References

[^1]: http://www.openstack.org/
[^2]: http://devstack.org/
[^3]: https://github.com/openstack-dev/devstack/tree/stable/grizzly
