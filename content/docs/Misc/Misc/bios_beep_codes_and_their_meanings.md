---
weight: 999
url: "/Signification_des_bips_émis_par_le_Bios/"
title: "BIOS Beep Codes and Their Meanings"
description: "A comprehensive guide to understanding the meaning of various BIOS beep codes during computer startup problems."
categories: ["Linux"]
date: "2008-04-09T10:03:00+02:00"
lastmod: "2008-04-09T10:03:00+02:00"
tags: ["bios", "hardware", "troubleshooting", "pc", "motherboard"]
toc: true
---

## Introduction

Many of you have probably wondered about the meaning of beeps when a computer fails to boot properly. Here's a list of what those beeps mean.

## Meanings

- **1 short beep**: Refresh failure. The system has problems accessing RAM for refresh operations. The issue is often due to a RAM or motherboard problem.

- **2 short beeps**: Parity error. Parity error in the first 64KB of memory. Problem with either RAM or the motherboard.

- **3 short beeps**: Base 64k memory failure. Memory failure in the first 64KB.

- **4 short beeps**: Timer not operational. Problem with one of the timers used to control motherboard functions. Cause: defective motherboard.

- **5 short beeps**: Processor error. Try removing and reinserting the processor, making sure it's properly seated. Note that this doesn't mean the processor is dead, or your system wouldn't boot at all.

- **6 short beeps**: 8042 - gate A20 failure. Keyboard problem. Try changing the keyboard or the keyboard controller chip (on the motherboard).

- **7 short beeps**: Processor exception interrupt error. Error in the virtual mode, which is one of the different modes in which the processor operates. The problem comes from the processor or motherboard.

- **8 short beeps**: Display memory read/write failure. The video controller is missing or the graphics card RAM is defective.

- **9 short beeps**: ROM checksum error. Error in the BIOS ROM. Either replace the BIOS chip or change the motherboard. Note that it may be another issue with the motherboard.

- **10 short beeps**: CMOS shutdown register read/write error. Error accessing CMOS memory. Problem with the motherboard.

- **11 short beeps**: Cache memory bad. External cache memory error. Try to properly reseat the cache memory.

- **1 long beep and 2 short beeps or 1 long beep and 3 short beeps**: Video error. Try reinserting the graphics card or its memory extension. If the problem persists, try another graphics card.

- If there are no beeps and nothing on screen, the first thing to check is the power supply. To determine if the power supply is OK, connect an LED to the motherboard power LED, for example. If the LED lights up and the hard drive or CD-ROM/burner start, the power supply is normally OK.

- If problems persist, remove all components except the graphics card. If the system starts in this case, try plugging in (not while the system is on, of course) the other components one by one to see what's causing the problem.

- If these explanations and their solutions haven't fixed anything: Contact technical support.
