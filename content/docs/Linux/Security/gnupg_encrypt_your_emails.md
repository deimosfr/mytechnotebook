---
weight: 999
url: "/GnuPG\\:_Crypter_vos_emails/"
title: "GnuPG: Encrypt Your Emails"
description: "A guide on how to encrypt your emails using GnuPG (GNU Privacy Guard) with instructions for key generation, server registration, revocation and usage."
categories: ["Linux"]
date: "2008-09-10T16:12:00+02:00"
lastmod: "2008-09-10T16:12:00+02:00"
tags: ["Security", "Email", "Encryption", "Privacy", "GnuPG", "PGP"]
toc: true
---

## Introduction

The purpose of this article is to show you how to encrypt your emails. For [Mozilla Thunderbird](https://www.mozilla.org), there is the [Enigmail](https://enigmail.mozdev.org/home/index.php) plugin. If you don't have an email client, don't want one and prefer webmail... that's a shame.

But wait! :-) There is a plugin called [FireGPG](https://firegpg.tuxfamily.org/) for [Mozilla Firefox](https://www.mozilla.org) that allows you to encrypt your emails :-)

To start, you need to install GPG on your computer. GPG (GnuPG) is the Open Source version of PGP which is paid software.

## Installation

Start by installing GPG like any other program:

```bash
apt-get install gnupg
```

## Generating Your Key

### Creating a Key Pair

"Hold on, what's a key pair?"

Let me give you some basic reminders about public key/private key encryption:

This system relies on two keys. The public key is used to encrypt a message, and the private key is used to decrypt a message that has been encrypted using the public key. In practice, everyone keeps their private key safely at home and gives their public key to all their contacts. So if I want to encrypt my email and send it to John, I need to have John's public key in my keyring. If I don't have it, it's either because John hasn't given it to me, or simply because John doesn't use GPG (which is his right), in which case it's impossible for me to send an encrypted email to John!

So to summarize, if one of my contacts wants to send me an encrypted email, they must first have my public key in their keyring. They will encrypt their email with it, and from then on, the encrypted email can only be decrypted with my private key (not someone else's private key).

The sender themselves would be unable to decrypt the email they just sent - it's absolutely irreversible. An email encrypted using a public key can only be decrypted by the private key generated at the same time as the public key. That's the brief theoretical reminder!

Now for the practical part, open a shell and type:

```bash
$ gpg --gen-key
gpg (GnuPG) 1.4.7; Copyright (C) 2006 Free Software Foundation, Inc.
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under certain conditions. See the file COPYING for details.

Select the type of key you want:
   (1) DSA and Elgamal (default)
   (2) DSA (sign only)
   (5) RSA (sign only)
```

Choose the default option:

```bash
$ Your choice? 1
The DSA key pair will have 1024 bits.
ELG-E keys may be between 1024 and 4096 bits long.
```

Personally, I'm a bit paranoid, so I put 4096, but 2048 is sufficient:

```bash
What key size do you want? (2048) 4096
The requested size is 4096 bits
Specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
```

I don't want the key to expire so I don't have to recreate it later, but sometimes it's better to do so. It depends on your needs.

```bash
Key is valid for? (0) 0
Key does not expire at all
Is this correct? (y/N) y

You need a user ID to identify your key; the software constructs the user ID
from the Real Name, Comment and Email address in this form:
   << Heinrich Heine (Der Dichter) <heinrichh@duesseldorf.de> >>

Real name: Pierre Mavro
Email address: username@domain.com
Comment:
You selected this USER-ID:
    "Pierre Mavro <username@domain.com>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? o
You need a Passphrase to protect your secret key.

A lot of random bytes need to be generated. You should do other work
(type on the keyboard, move the mouse, utilize the disks)
during prime generation; this gives the random number
generator a better chance to get enough entropy.
+++++++++++++++.+++++++++++++++++++++++++.+++++.++++++++++.+++++.+++++.+++++++++++++++++++++++++..+++++.++++++++++++++++++++++++++++++++++++++++.>.+++++............>+++++.....<+++++.....>+++++...<+++++...+++++
A lot of random bytes need to be generated. You should do other work
(type on the keyboard, move the mouse, utilize the disks)
during prime generation; this gives the random number
generator a better chance to get enough entropy.
+++++.+++++.++++++++++.++++++++++.++++++++++.+++++.++++++++++.+++++..++++++++++++++++++++.++++++++++.+++++.+++++..+++++++++++++++++++++++++++++++++++.+++++++++++++++.+++++.+++++...+++++>++++++++++......+++++...+++++...+++++..+++++.++++++++++++++++++++>+++++>+++++>.+++++...>+++++...<+++++................................>.+++++........................................................................................+++++^^^^
public and secret key created and signed.

gpg: checking the trustdb
gpg: 3 marginal(s) needed, 1 complete(s) needed, PGP trust model
gpg: depth: 0  valid:   1  signed:   0
trust: 0-. 0g. 0n. 0m. 0f. 1u
pub   1024D/A39D9E94 2008-01-06
    Key fingerprint = 1457 2EEC F76C 87CF B4A2  CB24 1405 33C6 A3DF 8093
uid                  Pierre Mavro <username@domain.com>
sub   4096g/E734E40D 2008-01-06
```

Answer these simple questions and wait for your key to be generated. To verify that your key has been properly generated:

```bash
$ gpg --list-keys
/Users/pmavro/.gnupg/pubring.gpg
--------------------------------
pub   1024D/A39D9E94 2008-01-06
uid                  Pierre Mavro <username@domain.com>
sub   4096g/E734E40D 2008-01-06
```

Here we can see that my public key has been created, and has the ID A39D9E94.

Note that the small size of this ID (8 characters) suggests a risk of duplication with another generated public key, so to remove any doubt about key identification (although technically it's also possible to have a duplicate here), we prefer to use the key's fingerprint:

```bash
$ gpg --fingerprint
/Users/pmavro/.gnupg/pubring.gpg
--------------------------------
pub   1024D/A39D9E94 2008-01-06
    Key fingerprint = 1457 2EEC F76C 87CF B4A2  CB24 1405 33C6 A3DF 8093
uid                  Pierre Mavro <username@domain.com>
sub   4096g/E734E40D 2008-01-06
```

There you go, you can now read your public key fingerprint (1457 2EEC F76C 87CF B4A2 CB24 1405 33C6 A3DF 8093), ending with the 8 digits of your previously seen ID.
Note this fingerprint on a piece of paper, we'll use it later.

### Registering with a Key Server

Once your keys are generated, you need to store your public key on a server so that your contacts can find your public key and send you encrypted emails!

To do this, we'll use the server [https://pgp.mit.edu](https://pgp.mit.edu), but there are certainly others. Most of them sync with each other, but not always, so it's best to always use the same server and tell your contacts which one you use to make sure they can find you.

Another detail, creating a revocation certificate is essential, for example to notify the key server that your public key is no longer valid and that your contacts should stop using it!

Let's create this revocation certificate:

```bash
$ gpg --gen-revoke username@domain.com > username@domain.com.txt

sec  1024D/A39D9E94 2008-01-06 Pierre Mavro <username@domain.com>

Create a revocation certificate for this key? (y/N) y
Select the reason for the revocation:
  0 = No reason specified
  1 = Key has been compromised
  2 = Key is superseded
  3 = Key is no longer used
  Q = Cancel
(You should probably select 1 here)
Your decision? 0
Enter an optional description; end it with an empty line:
> Just in case
>
Reason for revocation: No reason specified
Just in case
Is this okay? (y/N) y

You need a passphrase to unlock the secret key for
user: << Pierre Mavro <username@domain.com> >>
1024-bit DSA key, ID A39D9E94, created 2008-01-06

ASCII armored output forced.
Revocation certificate created.

Move it to a medium you can hide; if Mallory gets
access to this certificate, he can use it to make your key
unusable.
A good idea is to print this certificate and store it somewhere
else, in case the medium becomes unreadable. But be careful:
the printer system of your machine might store these
data and make them available to others!
```

There you go, your certificate has been created. Keep the text file carefully, because if someone uses it, they can invalidate your keys by notifying the server they no longer exist, even though that's false!! We'll see how to use it later...

To get back to our subject, we need to register on the key server. It's very simple:

```bash
$ gpg --keyserver pgp.mit.edu --send-keys A39D9E94
gpg: sending key A39D9E94 to hkp server pgp.mit.edu
```

There you go, your key has been exported to the key server. You can go admire it on pgp.mit.edu by typing your 8-digit ID preceded by '0x' in the String field, or your name, email address, etc. For former president Jacques Chirac:

```
http://pgp.mit.edu:11371/pks/lookup?op=get&search=0xA4723848
```

## Revoking Your Key on the Key Server

If, as mentioned above, you no longer use your keys for some reason (compromised by someone who has your private key, you lost your password, etc.), you must notify the key server of this revocation!

To do this, use the revocation file you created above, and import it into your keyring with the command:

```bash
gpg --import revoc_username@domain.com.txt
```

Check that your keyring has correctly registered the revocation by listing your keys:

```bash
gpg --list-keys
```

Your key should now be marked as [revoked].
You can then send it back to the key server to update it:

```bash
gpg --keyserver pgp.mit.edu --send-keys A39D9E94
```

There you go, your key is revoked, and therefore unusable. You can now delete your public and private keys from your keyring. Delete the secret key first:

```bash
gpg --delete-secret-keys username@domain.com
```

then the public key(s) attached:

```bash
gpg --delete-keys username@domain.com
```

## Usage

Now, use your favorite software to encrypt your emails. Send your public key to your friends as well so they can easily decrypt them :-)

## Resources
- [https://gpglinux.free.fr/](https://gpglinux.free.fr/)
- [GnuPG - for more confidentiality]({{< ref "docs/Linux/Security/gnupg_encrypt_your_emails.md" >}})
