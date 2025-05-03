---
weight: 999
url: "/Convertir_du_WMA_en_MP3/"
title: "Converting WMA to MP3"
description: "A guide on how to convert WMA audio files to MP3 format using Linux tools like lame, mplayer and perl."
categories: ["Linux", "Debian"]
date: "2007-01-30T14:36:00+02:00"
lastmod: "2007-01-30T14:36:00+02:00"
tags: ["Audio", "Conversion", "MP3", "WMA"]
toc: true
---

## Introduction

To convert WMA format (thanks Microsoft) to MP3 format, you will need "lame", "mplayer", and "perl". You must have the unofficial package repositories configured in Debian to install lame.

## Installation

To install everything once you have the unofficial package repositories configured:

```bash
apt-get install lame mplayer perl
```

Then create a file named convert.pl and insert these lines:

```perl
#! /usr/bin/perl
### WMA TO MP3 CONVERTER
###
$dir=`pwd`;
chop($dir);
opendir(checkdir,"$dir");

while ($file=readdir(checkdir)) {
$orig_file=$file;

if ($orig_file !~ /\.wma$/i) {next};

print "Conversion in progress: $orig_file\n";

$new_wav_file=$orig_file;$new_wav_file=~s/\.wma/\.wav/;
$new_mp3_file=$orig_file;$new_mp3_file=~s/\.wma/\.mp3/;

$convert_to_wav="mplayer \"./$orig_file\" -ao pcm -aofile \"./$new_wav_file\"";
$convert_to_mp3="lame -h \"./$new_wav_file\" \"./$new_mp3_file\"";
$remove_wav="rm -rf \"./$new_wav_file\"";

print "EXEC 1: $convert_to_wav\n";
$cmd=`$convert_to_wav`;
print "EXEC 2: $convert_to_mp3\n";
$cmd=`$convert_to_mp3`;
print "REMOVE WAV: $remove_wav\n";
$cmd=`$remove_wav`;
print "\n\n";

}

print "Done....";
```

Then set the appropriate permissions:

```bash
chmod 755 convert.pl
```

Execute it in the folder containing your WMA files:

```bash
./convert.pl
```
