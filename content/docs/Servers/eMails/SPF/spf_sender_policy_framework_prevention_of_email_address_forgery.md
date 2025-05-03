---
weight: 999
url: "/SPF_(Sender_Policy_Framework)_\\:_Prévention_de_la_contrefaçon_d'adresses_mails/"
title: "SPF (Sender Policy Framework): Prevention of Email Address Forgery"
description: "An overview of SPF (Sender Policy Framework), how it works, what it needs to function, and how to configure it to prevent email address forgery."
categories: ["Linux", "Network", "Servers"]
date: "2006-12-28T22:59:00+02:00"
lastmod: "2006-12-28T22:59:00+02:00"
tags: ["Network", "Servers", "Security", "Email", "DNS", "Anti-spam"]
toc: true
---

## Introduction

SPF stands for Sender Policy Framework. SPF aims to be an anti-counterfeiting standard to prevent email address forgery.

SPF was born in 2003. Its creator, Meng Weng Wong, took the best features of Reverse MX and DMP (Designated Mailer Protocol) to create SPF.

SPF uses the return path (MAIL FROM) present in the message header, since all MTAs work with these fields. However, there is a new concept proposed by Microsoft: PRA, which stands for Purported Responsible Address. The PRA corresponds to the end-user address that an MUA (like Thunderbird) uses.

Thus, when we combine SPF and PRA, we can get the so-called Sender ID that allows a user receiving an email to perform verifications of MAIL FROM fields (SPF verification) and PRA. In a way, it is said that MTAs will check the MAIL FROM field and MUAs will check the PRA field.

For now, SPF needs DNS to work properly. This means that "reverse MX" records need to be published. These records specify which machines send email for a given domain. This is different from MX records, used today, which specify the machines that receive email for a given domain.

## What Does SPF Need to Function?

To protect your system with SPF, you must:

* Configure your DNS to add the TXT record where the information that SPF requires is introduced.
* Configure your email system (qmail, sendmail) to use SPF; this means performing verification on each message received on your server.

The first step will be accomplished on the DNS server where the domain is located. In the next section, we will discuss the details of the records. One thing you need to keep in mind is the syntax your DNS server uses (bind or djbdns). But don't be afraid: the official SPF site provides excellent help that will guide you.

## The SPF TXT Record

The SPF record is contained in a TXT record and its format is as follows:

```
v=spf1 [[pre] type [ext] ] ... [mod]
```

The meaning of each parameter is as follows:

{{< table "table-hover table-striped" >}}
| Parameters | Descriptions |
|-----------|--------------|
| v=spf1 | SPF version. When using SenderID, you might see v=spf2 |
| pre | Defines a return code when a match occurs.<br>Possible values are:<br>{{< table >}}| Values | Descriptions |<br>|-------|----------------|<br>| + | Default. Means "pass" when a test is conclusive. |<br>| - | Means "fail a test". This value is normally applied to -all to say that there were no previous matches. |<br>| ~ | Means "soft fail". This value is normally applied when a test is not conclusive. |<br>| ? | Means "neutral". This value is normally applied when a test is not conclusive. |{{< /table >}} |
| type | Defines the type to use for verifications<br>Possible values are:<br>{{< table >}}| Values | Descriptions |<br>|-------|----------------|<br>| include | to include tests of a provided domain. It is written as: include:domain |<br>| all | to end the sequence of tests. For example, if it's -all, then all tests that haven't been met so far fail. But if there is uncertainty, it can be used in the form of ?all which means that the test will be accepted. |<br>| ip4 | Uses an IP version 4 for verification. This can be used in the form ipv4:ipv4 or ipv4:ipv4/cidr to define a range. This type is most recommended because it gives the smallest load on DNS servers. |<br>| ip6 | Uses an IP version 6 for verification. |<br>| a | Uses a domain name for verification. This will perform a lookup on the DNS for an A RR. It can be used in the form a:domain, a:domain/cidr, or a/cidr. |<br>| mx | Uses the MX RR of the DNS for verification. The MX RR defines the receiving MTA; for example, if it's not the same as the sending MTA, the tests based on MX will fail. It can be used in the form mx:domain, mx:domain/cidr, or mx/cidr. |<br>| ptr | Uses the PTR RR of the DNS for verification. In this case, a PTR RR is used, as well as a reverse map query. If the hostname returned is in the same domain, the communication is verified. It can be used in the form ptr:domain<br>exist | Tests the existence of a domain. It can be written in the form exist:domain. |{{< /table >}} |
| ext | Defines an optional extension to the type. If omitted, then a single record is used for the query. |
| mod | This is the last directive of type and it acts as a record modifier.<br>{{< table >}}| Modifiers | Descriptions |<br>|-----------|----------------|<br>| redirect | Redirects the verification to use SPF records of a defined domain. It is used in the form redirect=domain. |<br>| exp | This record must be the last one and it allows customizing the failure message.<br>```<br>IN TXT "v=spf1 mx -all exp=getlost.example.com"<br>getlost IN TXT "You are not allowed to send a message for the domain"<br>``` |{{< /table >}} |
{{< /table >}}

## In Case I Am an ISP

ISPs will have some "problems" with their roaming users if they use mechanisms like POP-before-Relay instead of SASL SMTP.

Well, if you are an ISP concerned about spam and forgery, you must consider your email policy and start using SPF.

Here are some steps you should consider:

* First, configure your MTA to use SASL; for example, you can enable it on ports 25 and 587.
* Warn your users about the policy you are implementing (spf.pobox.com provides an example, see the references).
* Give your users a grace period; this means you will publish your SPF records in DNS but with a soft fail (~all) instead of a fail (-all) for the tests.

And with that, you protect your servers, your clients, and the world against spam...

There is a lot of information for you on the official SPF site... What are you waiting for?

## What Things Should You Pay Attention To?

SPF is a perfect solution to protect yourself against fraud. However, it has a limitation: traditional email forwarding will no longer work. You cannot simply receive an email in your MTA and forward it. You must rewrite the sender's address. Patches for common MTAs are provided on [the SPF site](https://spf.pobox.com/downloads.html). In other words, if you start publishing SPF records in DNS, you should also update your MTA to rewrite sender addresses, even if you don't yet verify SPF records.

## Conclusion

You might think that implementing SPF could be somewhat confusing. Well, indeed, it's not complicated and, by the way, you have great help that helps you accomplish your mission (see the references section).

If you are concerned about spam, then SPF will help you by protecting your domain from forgeries, and all you have to do is add a line of text in your DNS server and configure your email server.

The advantages that SPF brings are enormous. However, as I told someone, it's not as big a difference as between day and night. The benefits of SPF will come with time, as others adopt it.

I referenced Sender ID and its relationship to SPF, but I didn't elaborate on explanations about it. You probably already know the reason: Microsoft's policy is still the same, namely software patenting. In the references, you can see openspf.org's position on SenderID.

In a future article, we will talk about MTA configuration. See you later!

I hope I have given you a brief introduction to SPF. If you want to learn more about it, simply use the references that were used to write this article.

## References

[The official SPF site](https://spf.pobox.com/)  
[The official SPF FAQ](https://spf.pobox.com/faq.html)  
[The official SPF help](https://spf.pobox.com/wizard.html)  
[The position of openspf.org about SenderID](https://www.openspf.org/OpenSPF_community_position_v101.html)  
[An excellent article about SenderID and SPF](https://trends.newsforge.com/article.pl?sid=04/08/26/1326244&tid=29)  
[Warn your users about the SASL conversion](https://spf.pobox.com/saslconversion.html)  
[HOWTO - Define an SPF record](https://www.zytrax.com/books/dns/ch9/spf.html)
