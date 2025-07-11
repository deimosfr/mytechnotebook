---
title: 'Prusa: make your own custom firmware'
slug: prusa-make-your-own-custom-firmware/
description: 'Make your custom firmware for your Prusa printer'
categories: ['Prusa', '3D Print']
tags: ['Prusa', '3D Print', 'Development']
date: '2025-07-11T21:30:09+02:00'
---

## Introduction

As a Prusa printer owner, I’ve been exploring ways to customize the firmware to add new features. At the time of writing, I’m using a Core One.

I also have an MMU3, but its current integration isn't ideal, and Prusa’s official solutions are quite limited. While waiting for a better hardware MMU3 release from Prusa, I wanted to use my existing setup with the MMU3—without going as far as building a fully custom solution like the [Coreboxx](https://www.printables.com/model/1168032-prusa-core-one-coreboxx).

I attempted to connect the MMU3 directly to the Core One without modifying the hardware, but it didn’t work. The configurable Bowden tube length was too short—capped at 1000mm (in the printer configuration), whereas my setup requires 1600mm.

So, I decided to develop custom firmware for the Core One. Here’s how I got started.

## Setup pre-requisites

On the official [Prusa Firmware Buddy](https://github.com/prusa3d/Prusa-Firmware-Buddy), you'll find instructions to setup the environment.

I first installed on my Mac the following software:

```bash
brew install jimctl
brew install openocd --HEA
```

Then I cloned the repository, created a Python virtual environment and installed the requirements:

```bash
git clone https://github.com/prusa3d/Prusa-Firmware-Buddy.git
cd Prusa-Firmware-Buddy
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

I also checkout on the tag I wanted to be to avoid any issue with the main branch (set the version you want to build):

```bash
git checkout v6.3.4-RC
```

## Patch the firmware

On my side, it was really simple to patch the firmware as it was mostly int values to change. I had to update the following line:

```diff
diff --git a/lib/Prusa-Firmware-MMU/src/config/config.h b/lib/Prusa-Firmware-MMU/src/config/config.h
index 313d64174..f95104c8a 100644
--- a/lib/Prusa-Firmware-MMU/src/config/config.h
+++ b/lib/Prusa-Firmware-MMU/src/config/config.h
@@ -94,7 +94,7 @@ static constexpr U_mm couplerToBowden = 3.5_mm; /// FINDA Coupler screw to bowde
 // Min, max and default bowden length setup
 static constexpr U_mm defaultBowdenLength = 360.0_mm; /// ~360.0_mm - Default Bowden length.
 static constexpr U_mm minimumBowdenLength = 341.0_mm; /// ~341.0_mm - Minimum bowden length.
-static constexpr U_mm maximumBowdenLength = 1000.0_mm; /// ~1000.0_mm - Maximum bowden length.
+static constexpr U_mm maximumBowdenLength = 5000.0_mm; /// ~5000.0_mm - Maximum bowden length.
 static_assert(minimumBowdenLength.v <= defaultBowdenLength.v);
 static_assert(maximumBowdenLength.v > defaultBowdenLength.v);
 
@@ -120,7 +120,7 @@ static constexpr AxisConfig pulley = {
 
 /// Pulley motion limits
 static constexpr PulleyLimits pulleyLimits = {
-    .lenght = 1000.0_mm, // TODO
+    .lenght = 5000.0_mm, // Updated to match maximum bowden length
     .jerk = 4.0_mm_s,
     .accel = 800.0_mm_s2,
 };
diff --git a/src/gui/menu_item/specific/menu_items_hw_mmu.cpp b/src/gui/menu_item/specific/menu_items_hw_mmu.cpp
index 2791a5b25..3f9fe9283 100644
--- a/src/gui/menu_item/specific/menu_items_hw_mmu.cpp
+++ b/src/gui/menu_item/specific/menu_items_hw_mmu.cpp
@@ -8,7 +8,7 @@ using namespace buddy;
 
 static constexpr NumericInputConfig mmu_bowden_limits {
     .min_value = 341,
-    .max_value = 1000,
+    .max_value = 5000,
     .unit = Unit::millimeter,
 };
```

## Build the firmware

Then it's time to build the firmware:

```bash
$ ./utils/build.py --preset coreone --build-type release --final
...
-- Build files have been written to: /Volumes/workspace/github/Prusa-Firmware-Buddy/build/coreone_release_noboot/xbuddy_extension-build
[199/1258] Performing build step for 'xbuddy_extension'
[150/150] Linking CXX executable firmware
Memory region         Used Size  Region Size  %age Used
             RAM:        7264 B        32 KB     22.17%
FLASH_BOOTLOADER:          0 GB         8 KB      0.00%
           FLASH:       42664 B     122752 B     34.76%
   FW_DESCRIPTOR:         128 B        128 B    100.00%
Puppy fingerprint: 262204c611d5039ab4e7e0e7015a925b502c7cf8286ba152d20101489574c22a

   text    data     bss     dec     hex filename
  38328    4456    2936   45720    b298 /Volumes/workspace/github/Prusa-Firmware-Buddy/build/coreone_release_noboot/xbuddy_extension-build/firmware
[234/1258] Performing install step for 'xbuddy_extension'
Skipping install step
[868/1258] Building CXX object CMakeFiles/firmware.dir[868/1258] Building CXX object CMakeFiles/firmware.dir[1258/1258] Linking CXX executable firmware
Memory region         Used Size  Region Size  %age Used
           FLASH:     1841792 B         2 MB     87.82%
             RAM:      115116 B       192 KB     58.55%
       MEMORY_B1:          0 GB         0 GB
          CCMRAM:       64728 B        64 KB     98.77%

   text    data     bss     dec     hex filename
1841568     220  174764 2016552  1ec528 /Volumes/workspace/github/Prusa-Firmware-Buddy/build/coreone_release_noboot/firmware

Processing  /Volumes/workspace/github/Prusa-Firmware-Buddy/build/coreone_release_noboot/firmware.bin
        version:  6.3.4+10503
        board: 0
        printer: COREONE
        sha256sum: ad0dbe3569dc30f769c1fab5c94b64485d107f10daa064355f7bb84dcd09a9f4
        sign:      00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
Done

Building finished: 2 success, 0 failure(s).
 coreone_release_boot   SUCCESS
 coreone_release_noboot SUCCESS
```

Now you'll find the firmware in the `build` directory. Here my firware was in `./build/products/coreone_release_boot_6.3.4.bbf`.

## Prepare the printer

Before flashing the firmware, you need to prepare the printer. Official firmware are signed with a Prusa private key. To be able to flash the firmware, you need to break the appendix seal on the xBuddy board. I know it can be frightening, but it's ok if you ([more info here](https://help.prusa3d.com/article/flashing-custom-firmware-core-one-mk4-s-mk3-9-s-mk3-5-s_814967)).

!!! note "Doing this doesn't break the warranty!"

## Flash the firmware

Now you can flash the firmware on your printer. Send the firmware through the USB stick or via [Prusa Connect](https://connect.prusa3d.com/fr).

Reboot to flash with your custom firmware and it's done!

## MMU custom firmware

I noticed having issues while saving the settings from the control panel. There was another control on the MMU3 forcing me to perform the same modification as the Core One. Hopefully Prusa did once again a great job [with the MMU](https://github.com/prusa3d/Prusa-Firmware-MMU) and the build is working the same way as for the Core One firmware.

So I did this modification on the MMU firmware (similar to the Core One firmware):

```diff
diff --git a/src/config/config.h b/src/config/config.h
index 313d641..93c0520 100644
--- a/src/config/config.h
+++ b/src/config/config.h
@@ -94,7 +94,7 @@ static constexpr U_mm couplerToBowden = 3.5_mm; /// FINDA Coupler screw to bowde
 // Min, max and default bowden length setup
 static constexpr U_mm defaultBowdenLength = 360.0_mm; /// ~360.0_mm - Default Bowden length.
 static constexpr U_mm minimumBowdenLength = 341.0_mm; /// ~341.0_mm - Minimum bowden length.
-static constexpr U_mm maximumBowdenLength = 1000.0_mm; /// ~1000.0_mm - Maximum bowden length.
+static constexpr U_mm maximumBowdenLength = 5000.0_mm; /// ~5000.0_mm - Maximum bowden length.
 static_assert(minimumBowdenLength.v <= defaultBowdenLength.v);
 static_assert(maximumBowdenLength.v > defaultBowdenLength.v);
 
@@ -120,7 +120,7 @@ static constexpr AxisConfig pulley = {
 
 /// Pulley motion limits
 static constexpr PulleyLimits pulleyLimits = {
-    .lenght = 1000.0_mm, // TODO
+    .lenght = 5000.0_mm, // TODO
     .jerk = 4.0_mm_s,
     .accel = 800.0_mm_s2,
 };
diff --git a/src/modules/permanent_storage.cpp b/src/modules/permanent_storage.cpp
index 9c5fb0a..3f9a296 100644
--- a/src/modules/permanent_storage.cpp
+++ b/src/modules/permanent_storage.cpp
@@ -48,7 +48,7 @@ static eeprom_t *const eepromBase = reinterpret_cast<eeprom_t *>(0); ///< First
 constexpr const uint16_t eepromEmpty = 0xffffU; ///< EEPROM content when erased
 constexpr const uint16_t eepromBowdenLenDefault = config::defaultBowdenLength.v; ///< Default bowden length (~360 mm)
 constexpr const uint16_t eepromBowdenLenMinimum = config::minimumBowdenLength.v; ///< Minimum bowden length (~341 mm)
-constexpr const uint16_t eepromBowdenLenMaximum = config::maximumBowdenLength.v; ///< Maximum bowden length (~1000 mm)
+constexpr const uint16_t eepromBowdenLenMaximum = config::maximumBowdenLength.v; ///< Maximum bowden length (~5000 mm)
 
 namespace ee = hal::eeprom;
 
diff --git a/src/registers.cpp b/src/registers.cpp
index 58512d2..d858a51 100644
--- a/src/registers.cpp
+++ b/src/registers.cpp
@@ -165,7 +165,7 @@
 | 0x1fh 31 | uint16  |Set/Get Selector iRun current| 0-31         | 1fh 31      | 31->530mA: see TMC2130 current conversion| Read / Write | M707 A0x1f | M708 A0x1f Xn
 | 0x20h 32 | uint16   | Set/Get Idler iRun current | 0-31         | 1fh 31      | 31->530mA: see TMC2130 current conversion| Read / Write | M707 A0x20 | M708 A0x20 Xn
 | 0x21h 33 | uint16   | Reserved for internal use  | 225          |             | N/A                                      | N/A          | N/A        | N/A
-| 0x22h 34 | uint16   | Bowden length              | 341-1000     | 168h 360    | unit mm                                  | Read / Write Persistent | M707 A0x22 | M708 A0x22 Xn
+| 0x22h 34 | uint16   | Bowden length              | 341-5000     | 168h 360    | unit mm                                  | Read / Write Persistent | M707 A0x22 | M708 A0x22 Xn
 | 0x23h 35 | uint8    | Cut length                 | 0-255        | 8           | unit mm                                  | Read / Write | M707 A0x23 | M708 A0x23 Xn
 */
```

Then ran a build to get a new firmware (`.dex` extension) in the `build` directory.

The only change is on the update, I did not find a way to be able to update it directly from the main firmware so I used a USB cable connected to the MMU3 to update it [as described here](https://help.prusa3d.com/article/how-to-update-firmware-mmu3_701584).

## Conclusion

I'm really happy with the result. I can now use my MMU3 with the Core One.  I did not expect it to be so easy to build a custom firmware for my printer.

![MMU3 with Core One](../../static/images/prusa_coreone_mmu3_cust_firm.avif)

Prusa doing Open Source is not just words. They did a lot of things to make it possible for newcomers. I was already loving Prusa since the MK4, but now I'm even more impressed.