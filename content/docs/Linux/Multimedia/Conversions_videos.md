---
weight: 999
url: "/Conversions_videos/"
title: "Video Conversions"
description: "Guide on how to convert videos between different formats using mencoder and ffmpeg on Linux"
categories: ["Ubuntu", "Linux"]
date: "2010-05-02T21:43:00+02:00"
lastmod: "2010-05-02T21:43:00+02:00"
tags: ["Video", "Conversion", "ffmpeg", "mencoder", "Multimedia"]
toc: true
---

## Introduction

mencoder is one of the most used applications along with ffmpeg to convert or encode videos.

## Installation

Depending on what you will use, select only the one you wish:

```bash
apt-get install mencoder ffmpeg
```

If you want to create OGV (OGG/Theora) videos for HTML5, you need to install this package:

```bash
aptitude install ffmpeg2theora
```

## Usage

### Convert a .wmv to a .avi

```bash
mencoder "/path/to/file.wmv" -ofps 23.976 -ovc lavc -oac copy -o "/path/to/file.avi"
```

### Convert a .mp4 to a .avi

```bash
ffmpeg -i "/path/to/file.mp4" "/path/to/file.avi"
```

And to convert an avi to mp4 and change to iPhone resolution:

```bash
ffmpeg -s 480x320 -i /path/to/file.avi /path/to/file_iphone.mp4
```

### Convert a .flv to a .mpg

```bash
fmpeg -i get_video.flv -ab 56 -ar 22050 -b 500 -s 320x240 test.mpg
```

The options are:

- b bitrate: set the video bitrate in kbit/s (default = 200 kb/s)
- ab bitrate: set the audio bitrate in kbit/s (default = 64)
- ar sample rate: set the audio samplerate in Hz (default = 44100 Hz)
- s size: set frame size. The format is WxH (default 160x128)

### Convert to ogv

```bash
ffmpeg2theora video.mov
```

This will keep the same video information and create an ogv file with the same name.

If you want better quality:

```bash
ffmpeg2theora --optimize video.mov
```

## Resources
- http://soukie.net/degradable-html5-audio-and-video-plugin/#instr
- http://www.paperblog.fr/3023554/ffmpeg2theora-guide-par-l-exemple/
- http://doc.ubuntu-fr.org/ffmpeg
