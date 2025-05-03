---

weight: 999
url: "/SSH_\\:_Mise_en_place_du_serveur_SSH_Solaris_sur_une_installation_minimale/"
title: "SSH: Setting up SSH Server on a Minimal Solaris Installation"
description: "How to install SSH server packages on a minimal Solaris 10 installation"
categories: ["Solaris", "Servers", "Network"]
date: "2008-04-05T10:16:00+02:00"
lastmod: "2008-04-05T10:16:00+02:00"
tags: ["SSH", "Solaris", "Server", "Installation"]
toc: true

---

If you need to manage a Solaris 10 box with a minimal install, and SSH is not available, you can install it off of the 2nd CD. Rather than figure out the path to your CDROM (see this article), it was easier in our case to just tar up the needed packages and FTP them to our Solaris box:

```bash
root@srv-3 Product # cp -R SUNWsshcu SUNWsshdr SUNWsshdu SUNWsshr 
SUNWsshu /home/srv-1/sshpkg/
root@srv-3 Product # cd /home/srv-1/sshpkg/
root@srv-3 sshpkg # ls
SUNWsshcu  SUNWsshdr  SUNWsshdu  SUNWsshr  SUNWsshu
root@srv-3 sshpkg # tar -cf ../ssh.tar *
root@srv-3 sshpkg # tar -tf ../ssh.tar
SUNWsshcu/
SUNWsshcu/archive/
.
.
.
SUNWsshu/reloc/
SUNWsshu/reloc/usr/
SUNWsshu/reloc/usr/bin/
SUNWsshu/reloc/usr/lib/
SUNWsshu/reloc/usr/lib/ssh/
root@srv-3 sshpkg #
```

On the Solaris side, FTP these to /tmp, then from tmp:

```bash
tar -xf ssh.tar
pkgadd -d .
svcadm enable ssh
svcadm restart ssh
```
