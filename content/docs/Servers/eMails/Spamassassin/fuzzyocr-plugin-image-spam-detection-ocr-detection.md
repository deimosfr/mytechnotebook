---
weight: 999
url: "/FuzzyOcr_Plugin_\\:_Détection_des_spams_image_(OCR_Detection)/"
title: "FuzzyOcr Plugin: Image Spam Detection (OCR Detection)"
description: "Learn how to install and configure FuzzyOcr Plugin for SpamAssassin to detect image-based spam using OCR technology."
categories: ["Servers", "Security", "AntiSpam"]
date: "2008-04-10T10:36:00+02:00"
lastmod: "2008-04-10T10:36:00+02:00"
tags: ["SpamAssassin", "FuzzyOCR", "OCR", "Anti-Spam", "Email"]
toc: true
---

## Introduction

The most widely used anti-spam tool in the open-source world is SpamAssassin. Its operation is similar to that of the amavisd-new service. It collects results obtained from a collection of other tools and passes them to the amavisd-new service, which makes a decision about the message based on its total score.

From a system perspective, implementing SpamAssassin is very simple. Just use the Debian/testing package, which is well maintained. This situation is particularly interesting because the collection of libraries and tools dependent on SpamAssassin is very large. The command:

```bash
apt-cache show spamassassin
```

lists these dependencies.

## SpamAssassin Configuration

There are 2 levels to consider when configuring this tool:

* At the system level, the `/etc/spamassassin` directory contains configuration files common to all users.
* For the amavis user, the `.spamassassin` directory located under the user's home directory contains dedicated configuration files as well as databases created during execution: automatic whitelist, calculation tokens, etc.

Since the scenarios described in this document use a dedicated gateway server for email processing, all configuration parameters will be in the general `/etc/spamassassin` directory and all databases will be in the `/var/lib/amavis/.spamassassin` directory.

## Optical Character Recognition

The trend since fall 2006 has been an explosion of spam containing animated GIF, PNG, or JPEG images. These images most often contain advertisements or dubious stock market offers.

The FuzzyOcr Plugin is a fairly effective weapon against image-containing spam. It uses optical character recognition to search for keywords in these images. This plugin also has some optimizations for processing deliberately corrupted image files.

The basic operation of the plugin, as described on the presentation web page, covers the following steps:

* Search for images in different parts of the message
* Each image is analyzed to identify its format (GIF, PNG, JPEG)
* Depending on the detected image format, different tools are called to convert the image to PNM format
* The optical character recognition program gocr is called to extract text from the PNM file
* The obtained strings are scanned for predefined words, scores are calculated, and results are transmitted to SpamAssassin

## Installing the FuzzyOCR Plugin

Unfortunately, a Debian/testing package is not yet available. Therefore, we need to do a manual installation based on the packages of the tools used by the plugin.

First, download the plugin code and place the sources in the email gateway server's directory tree:

```bash
wget http://users.own-hero.net/~decoder/fuzzyocr/fuzzyocr-latest.tar.gz
mv fuzzyocr-latest.tar.gz /usr/local/src/
cd /usr/local/src/
tar xf fuzzyocr-latest.tar.gz
chown -R root.src FuzzyOcr*
```

Next, make sure these packages are installed:

```bash
apt-get install spamassassin netpbm imagemagick libungif4g libungif-bin gocr libjpeg-progs libstring-approx-perl
```

## Configuring the FuzzyOCR Plugin

For this part, we don't follow the recommendations in the INSTALL file. Indeed, the default plugin configuration expects to use the `/etc/mail/spamassassin` directory, which doesn't exist on a Debian system. Furthermore, it's not desirable to mix code and configuration elements in the `/etc/` directory.

For the plugin code installation, the command `# dpkg -L spamassassin | grep Plugin/` helps identify the directory where other plugins used by SpamAssassin are stored: `/usr/share/perl5/Mail/SpamAssassin/Plugin/`. So we copy the code to this directory.

```bash
cp FuzzyOcr.pm /usr/share/perl5/Mail/SpamAssassin/Plugin/
```

For the configuration file, simply copy the following files to the general SpamAssassin configuration directory:

```bash
cp FuzzyOcr.cf /etc/spamassassin/
cp FuzzyOcr.words.sample /etc/spamassassin/FuzzyOcr.words
cd /etc/spamassassin
```

Several configuration files need to be edited to adapt to the configuration context.
We need to indicate the new plugin in the list traversed by SpamAssassin by adding a line at the end of the v310.pre file:

```bash
echo loadplugin FuzzyOcr /usr/share/perl5/Mail/SpamAssassin/Plugin/FuzzyOcr.pm >> v310.pre
```

We need to edit the FuzzyOcr.cf file to adapt it to the context (`/etc/spamassassin/FuzzyOcr.cf`):

```perl
#loadplugin FuzzyOcr FuzzyOcr.pm
body FUZZY_OCR eval:fuzzyocr_check()
describe FUZZY_OCR Mail contains an image with common spam text inside
body FUZZY_OCR_WRONG_CTYPE eval:dummy_check()
describe FUZZY_OCR_WRONG_CTYPE Mail contains an image with wrong content-type set
body FUZZY_OCR_CORRUPT_IMG eval:dummy_check()
describe FUZZY_OCR_CORRUPT_IMG Mail contains a corrupted image
body FUZZY_OCR_KNOWN_HASH eval:dummy_check()
describe FUZZY_OCR_KNOWN_HASH Mail contains an image with known hash

priority FUZZY_OCR             900

########### Plugin Configuration #############

#### Logging options #####
# Verbosity level (see manual) Attention: Don't set to 0, but to 0.0 for quiet operation. (Default value: 1)
focr_verbose 4
#
# Logfile (make sure it is writable by the plugin) (Default value: /etc/mail/spamassassin/FuzzyOcr.log)
focr_logfile /var/log/fuzzyocr.log
##########################

##### Wordlists #####
# Here we defined the words to scan for (Default value: /etc/mail/spamassassin/FuzzyOcr.words)
focr_global_wordlist /etc/spamassassin/FuzzyOcr.words
#
# This is the path RELATIVE to the respektive home directory for the personalized list
# This list is merged with the global word list on execution (Default value: .spamassassin/fuzzyocr.words)
#focr_personal_wordlist .spamassassin/fuzzyocr.words
#####################

# Set this to 1 if you are running a version < 3.1.4.
# This will disable a function used in conjunction with animated gifs that isn't available in earlier versions (Default value: 0.0)
#focr_pre314 0.0

# These parameters can be used to change other detection settings
# If you leave these commented out, the defaults will be used.
# Do not use " " around any parameters!
#
##### Location of helper applications (path + binary) (Default values: /usr/bin/<app>) #####
#focr_bin_giffix /usr/bin/giffix
#focr_bin_giffix /usr/bin/giffix
#focr_bin_giftext /usr/bin/giftext
#focr_bin_gifasm /usr/bin/gifasm
#focr_bin_gifinter /usr/bin/gifinter
#focr_bin_giftopnm /usr/bin/giftopnm
#focr_bin_jpegtopnm /usr/bin/jpegtopnm
#focr_bin_pngtopnm /usr/bin/pngtopnm
#focr_bin_ppmhist /usr/bin/ppmhist
#focr_bin_convert /usr/bin/convert
#focr_bin_identify /usr/bin/identify
#focr_bin_gocr /usr/bin/gocr
############################################################################################

##### Scansets, comma seperated (Default value: $gocr -i -, $gocr -l 180 -d 2 -i -) #####
# Each scanset consists of one or more commands which make text out of pnm input.
# Each scanset is run seperately on the PNM data, results are combined in scoring.
#focr_scansets $gocr -i -, $gocr -l 180 -d 2 -i -
#
# To use only one scan with default values, uncomment the next line instead
#focr_scansets $gocr -i -
#
# Some example for more advanced sets
# Thisone uses the first the standard scan, then a scanset which first reduces the image to 3 colors and then scans it with custom settings
# and then it scans again only with these custom settings
# NOTE: This is for advanced users only, if you have questions how to use this, ask on the ML or on IRC
#focr_scansets $gocr -i -, pnmnorm 2>$errfile 
```

If your version is significantly different from the one above, here's a patch between the gateway configuration and the distributed version:

```bash
diff -uBb /usr/local/src/FuzzyOcr-2.3b/FuzzyOcr.cf    2006-08-25 22:56:00.000000000 +0200
+++ FuzzyOcr.cf 2006-09-10 23:23:39.000000000 +0200
@@ -1,4 +1,4 @@
-loadplugin FuzzyOcr FuzzyOcr.pm
+#loadplugin FuzzyOcr FuzzyOcr.pm
body FUZZY_OCR eval:fuzzyocr_check()
describe FUZZY_OCR Mail contains an image with common spam text inside
body FUZZY_OCR_WRONG_CTYPE eval:dummy_check()
@@ -14,15 +14,15 @@
 
#### Logging options #####
# Verbosity level (see manual) Attention: Don't set to 0, but to 0.0 for quiet
# operation.
-#focr_verbose 1
+focr_verbose 4
#
# Logfile (make sure it is writable by the plugin)
-focr_logfile /etc/mail/spamassassin/FuzzyOcr.log
+focr_logfile /var/log/fuzzyocr.log
##########################
 
##### Wordlists #####
# Here we defined the words to scan for
-focr_global_wordlist /etc/mail/spamassassin/FuzzyOcr.words
+focr_global_wordlist /etc/spamassassin/FuzzyOcr.words
#
# This is the path RELATIVE to the respektive home directory for the personalized list
# This list is merged with the global word list on execution
@@ -90,6 +90,7 @@
#
# This is used to disable the OCR engine if the message has already more points
# than this value (Default value: 10)
#focr_autodisable_score 10
+#focr_autodisable_score 50
#
# Number of minimum matches before the rule scores (Default value: 2)
#focr_counts_required 2
```

Finally, we need to implement the configuration needed for logging:

```bash
touch /var/log/fuzzyocr.log
chown amavis.amavis /var/log/fuzzyocr.log
```

To limit the disk usage of the plugin's logging while maintaining a 90-day history, create a configuration file for the logrotate service. Create this file and insert the following lines (`/etc/logrotate.d/fuzzyocr`):

```
$ /etc/logrotate.d/fuzzyocr
/var/log/fuzzyocr.log {
    rotate 90
    daily
    compress
    delaycompress
    create 666 amavis amavis
    }
```

## Validation Tests

Before any manipulation, ensure that SpamAssassin is functioning correctly. The two commands:

```bash
spamassassin --lint
```

and

```bash
su - amavis -- spamassassin --lint
```

should not produce any output. If this is not the case, you should review your SpamAssassin or amavis configuration before testing the plugin's operation.

Next, use the spam samples provided with the plugin sources and analyze the results:

```bash
$ spamassassin -t /usr/local/src/FuzzyOcr-2.3b/samples/animated-gif.eml
<snipped>
 20 FUZZY_OCR              BODY: Mail contains an image with common spam text inside
                           Words found:
                           "alert" in 4 lines
                           "charts" in 1 lines
                           "symbol" in 1 lines
                           "alert" in 4 lines
                           "stock" in 2 lines
                           "company" in 3 lines
                           "trade" in 1 lines
                           "meridia" in 1 lines
                           "growth" in 1 lines
                           (18 word occurrences found)
```

```bash
$ spamassassin -t /usr/local/src/FuzzyOcr-2.3b/samples/corrupted-gif.eml
<snipped>
1.5 FUZZY_OCR_WRONG_CTYPE  BODY: Mail contains an image with wrong
                           content-type set
                           Image has format "GIF" but content-type is
                           "image/jpeg"
2.5 FUZZY_OCR_CORRUPT_IMG  BODY: Mail contains a corrupted image
                           Corrupt image: GIF-LIB error: Image is
                           defective, decoding aborted.
 14 FUZZY_OCR              BODY: Mail contains an image with common spam text inside
                           Words found:
                           "alert" in 1 lines
                           "alert" in 1 lines
                           "stock" in 2 lines
                           "investor" in 1 lines
                           "company" in 1 lines
                           "price" in 2 lines
                           "trade" in 1 lines
                           "target" in 1 lines
                           "service" in 1 lines
                           "recommendation" in 1 lines
                           (12 word occurrences found)
```

```bash
$ spamassassin -t /usr/local/src/FuzzyOcr-2.3b/samples/jpeg.eml
<snipped>
 6.0 FUZZY_OCR             BODY: Mail contains an image with common spam text inside
                           Words found:
			   "viagra" in 2 lines
			   "cialis" in 1 lines
			   "levitra" in 1 lines
			   (4 word occurrences found)
```

We can see that the reports on the three tested samples show that keywords were detected in the images and that the scores were increased accordingly.

Now we just need to validate the solution in real operation! Without waiting, we get a sample that shows that optical character recognition added 18 points to the score.

Extract from the system log produced by the amavisd-new service:

```
Sep 10 18:08:57 MailGw amavis[13387]: (13387-04) SPAM, \
  <lwrceuqp@demax.sk> -> <xxxxxxxxxxxxx@xxx.xxx-xxxx.fr>, Yes, \
  score=27.449 tag=-999 tag2=6.31 kill=6.31 tests=[BAYES_50=2.5, \
  DATE_IN_FUTURE_03_06=1.961, EXTRA_MPART_TYPE=1.091, FUZZY_OCR=18.000, \
  HTML_MESSAGE=0.001, RCVD_IN_XBL=3.897, SPF_PASS=-0.001], \
  autolearn=spam, quarantine lkwO3NF5SAMT (spam-quarantine)
```

Extract from the log produced by the plugin for the same message:

```
[2006-09-10 18:08:57] Debug mode: Message is spam (score 18)...
[2006-09-10 18:08:57] Debug mode: Words found:
                      "charts" in 1 lines
                      "symbol" in 1 lines
                      "stock" in 2 lines
	              "international" in 3 lines
		      "company" in 2 lines
		      "million" in 1 lines
		      "buy" in 1 lines
                      "trade" in 3 lines
                      "target" in 1 lines
                      "meridia" in 1 lines
                      (16 word occurrences found)
```

## FAQ

### I still receive image-type spam, what should I do?

You should know that if the spam is very well camouflaged, detection will fail. But also check what the logs say. I had a problem with permissions at the log level. Make a small correction at that level if that's what's blocking things.

## Other Documentation

[Here is another documentation in English](/pdf/fuzzyocr_debian.pdf)
