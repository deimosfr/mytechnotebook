---
weight: 999
url: "/Cacher_une_application_ouverte_du_dock_et_du_switcher/"
title: "Hide an Open Application from the Dock and Switcher"
description: "Learn how to hide Mac applications from the Dock and App Switcher while keeping them running in the background."
categories: ["Mac OS X", "Customization"]
date: "2008-04-08T07:16:00+02:00"
lastmod: "2008-04-08T07:16:00+02:00"
tags: ["macos", "dock", "switcher", "background apps", "plist"]
toc: true
---

## Introduction

If you want to use software that must remain open all the time but you don't want to have its icon in your Dock or in the switcher (Cmd + Tab), it's possible, though it doesn't always work in all cases.

## Prerequisites

There aren't many prerequisites since the only thing you'll need to be careful about is that before doing the following steps, your application must be closed.

## Info.plist

Now, let's get to the serious business. Go to your Applications folder and find the application you want to modify. In my case, I'm using [XRG (a tool that creates graphs)](https://www.gauchosoft.com/Software/X%20Resource%20Graph/).

Right-click on the application, then "Show Package Contents". Go to the "Contents" folder and edit the Info.plist file. There are only 2 lines to add at the end before `</dict>`:

```xml
<key>LSUIElement</key>
<string>1</string>
```

Which gives me in the end:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>CFBundleDevelopmentRegion</key>
        <string>English</string>
        <key>CFBundleDocumentTypes</key>
        <array>
               <dict>
                       <key>CFBundleTypeExtensions</key>
                       <array>
                               <string>xtf</string>
                       </array>
                       <key>CFBundleTypeIconFile</key>
                       <string>xtf2</string>
                       <key>CFBundleTypeName</key>
                       <string>XRG Theme File</string>
                       <key>CFBundleTypeOSTypes</key>
                       <array>
                               <string>XTF</string>
                       </array>
                </dict>
       </array>
       <key>CFBundleExecutable</key>
       <string>X Resource Graph</string>
       <key>CFBundleGetInfoString</key>
       <string>XRG v1.1u</string>
       <key>CFBundleHelpBookFolder</key>
       <string>Online Help</string>
       <key>CFBundleHelpBookName</key>
       <string>XRG Help</string>
       <key>CFBundleIconFile</key>
       <string>icon4.icns</string>
       <key>CFBundleIdentifier</key>
       <string>com.piatekjimenez.xrg</string>
       <key>CFBundleInfoDictionaryVersion</key>
       <string>6.0</string>
       <key>CFBundleName</key>
       <string>X Resource Graph</string>
       <key>CFBundlePackageType</key>
       <string>APPL</string>
       <key>CFBundleShortVersionString</key>
       <string>1.1u</string>
       <key>CFBundleSignature</key>
       <string>XRGA</string>
       <key>CFBundleVersion</key>
       <string>1.1u</string>
       <key>NSMainNibFile</key>
       <string>MainMenu</string>
       <key>NSPrincipalClass</key>
       <string>NSApplication</string>
       <key>LSUIElement</key>
       <string>1</string>
</dict>
</plist>
```

Now, all that's left is to save the changes and restart the application :-)

## Issues

### I made the changes but nothing different is happening

You need to log out and back into your session and try again.

### The application no longer opens

You need to move the added lines elsewhere in the code; something is not working correctly. In the worst case, you can always delete the 2 added lines and the application will work as before.
