---
weight: 999
url: "/Vagrant_\\:_quickly_deploy_virtual_machines/"
title: "Vagrant: Quickly Deploy Virtual Machines"
description: "Learn how to use Vagrant to easily deploy and manage virtual machines for development and testing environments."
categories: ["Debian", "AWS", "Storage"]
date: "2014-03-13T16:42:00+02:00"
lastmod: "2014-03-13T16:42:00+02:00"
tags: ["Vagrant", "VirtualBox", "Virtualization", "Development"]
toc: true
---

![Vagrant](/images/vagrant-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7 |
| **Website** | [Vagrant Website](https://www.vagrantup.com/) |
| **Last Update** | 13/03/2014 |
| **Others** | VirtualBox |
{{< /table >}}

## Introduction

Vagrant provides easy to configure, reproducible, and portable work environments built on top of industry-standard technology and controlled by a single consistent workflow to help maximize the productivity and flexibility of you and your team.

To achieve its magic, Vagrant stands on the shoulders of giants. Machines are provisioned on top of VirtualBox, VMware, AWS, or any other provider. Then, industry-standard provisioning tools such as shell scripts, Chef, or Puppet, can be used to automatically install and configure software on the machine.

## Installation

We need to install these prerequisites in order to install VirtualBox in the desired version:

```bash
aptitude install libsdl-ttf2.0-0:amd64 gcc-4.6-base:amd64 cpp-4.6 dkms gcc-4.6 linux-headers-amd64 linux-kbuild-3.2
```

Then get this VirtualBox and vagrant version to avoid incompatibility issues:

```bash
wget http://download.virtualbox.org/virtualbox/4.2.12/virtualbox-4.2_4.2.12-84980~Debian~wheezy_amd64.deb
wget http://files.vagrantup.com/packages/7e400d00a3c5a0fdf2809c8b5001a035415a607b/vagrant_1.2.2_x86_64.deb
```

## Usage

Each created instance should have its own folder. For example you can do this kind of hierarchy:

```bash
.
|-- Vagrant
|   |-- vm1
|   |-- vm2
|   `-- vm3
```

### Add an image

The first thing you have to do is to add an image. You can [find several ones here](https://www.vagrantbox.es/) or on [the official Vagrant Cloud](https://vagrantcloud.com/discover/featured). Let's take a Debian Squeeze for example (change squeeze name if you want):

```bash
vagrant box add squeeze http://www.emken.biz/vagrant-boxes/debsqueeze64.box
```

This will download and store the image in ~/.vagrant.d/boxes/.

Here is a Wheezy image:

```bash
vagrant box add deimosfr/debian-wheezy
```

Or Jessie:

```bash
vagrant box add deimosfr/debian-jessie
```

### Deploy an image

To deploy a downloaded box image, simply run:

```bash
vagrant init <squeeze>
```

Replace squeeze by the name of the box. This will create a Vagrantfile file.

### Start the image

To start an image, it's simple:

```bash
vagrant up
```

### Stop the image

You can shutdown:

```bash
vagrant halt
```

### Connect to the image

To connect through ssh, it's simple:

```bash
vagrant ssh
```

### List all boxes machines

To list all available boxes on your machine:

```bash
> vagrant box list
squeeze (virtualbox)
```

## Plugins

### VirtualBox Guest Additions

To avoid reinstalling manually each new version of guests, here is a plugin that will do it for you each time you boot a VM! Install the plugin:

```bash
vagrant plugin install vagrant-vbguest
```

That's it :). Now start a VM and if VirtualBox Guests Additions are not at the latest version, they will automatically be updated.

## Example

### Ceph

Here is an example for 6 VMs with 2 interfaces (`Vagrantfile`):

```ruby
# -*- mode: ruby -*-
# vi: set ft=ruby :
ENV['LANG'] = 'C'

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Insert all your Vms with configs
boxes = [
    { :name => :mon1, :role => 'mon'},
    { :name => :mon2, :role => 'mon'},
    { :name => :mon3, :role => 'mon'},
    { :name => :osd1, :role => 'osd', :ip => '192.168.33.31'},
    { :name => :osd2, :role => 'osd', :ip => '192.168.33.32'},
    { :name => :osd3, :role => 'osd', :ip => '192.168.33.33'},
]

$install = <<INSTALL
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
echo deb http://ceph.com/debian/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
aptitude update
aptitude -y install ceph ceph-deploy openntpd
INSTALL

Vagrant::Config.run do |config|
  # Default box OS
  vm_default = proc do |boxcnf|
    boxcnf.vm.box       = "deimosfr/debian-wheezy"
  end

  # For each VM, add a public and private card. Then install Ceph
  boxes.each do |opts|
    vm_default.call(config)
    config.vm.define opts[:name] do |config|
        config.vm.network   :bridged, :bridge => "eth0"
        config.vm.host_name = "%s.vm" % opts[:name].to_s
        config.vm.provision "shell", inline: $install
        # Create 8G disk file and add private interface for OSD VMs
        if opts[:role] == 'osd'
            config.vm.network   :hostonly, opts[:ip]
            file_to_disk = 'osd-disk_' + opts[:name].to_s + '.vdi'
            config.vm.customize ['createhd', '--filename', file_to_disk, '--size', 8 * 1024]
            config.vm.customize ['storageattach', :id, '--storagectl', 'SATA', '--port', 1, '--device', 0, '--type', 'hdd', '--medium', file_to_disk]
        end
    end
  end
end
```

To boot them, simply "vagrant up".

## Conclusion

I've written here basic commands as an introduction but you can do more than that with Vagrant. Read [the documentation](https://docs.vagrantup.com) for more information.

## FAQ

### VirtualBox is complaining that the kernel module is not loaded

If you got this kind of error message while doing a "vagrant up" command:

```bash
VirtualBox is complaining that the kernel module is not loaded. Please
run `VBoxManage --version` to see the error message which should contain
instructions on how to fix this error.
```

You need to install your kernel sources or kernel headers package and launch modules compilation:

```bash
> /etc/init.d/vboxdrv setup
[ ok ] Stopping VirtualBox kernel modules:.
[ ok ] Recompiling VirtualBox kernel modules:.
[ ok ] Starting VirtualBox kernel modules:.
```

## References

1. http://www.vagrantbox.es/
2. https://vagrantcloud.com/discover/featured
3. http://docs.vagrantup.com
