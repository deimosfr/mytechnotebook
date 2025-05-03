---
weight: 999
url: "/Déployement_de_la_gestion_d'énergie_sur_un_serveur_pour_ses_clients/"
title: "Deploying Power Management on a Server for its Clients"
description: "Guide on how to deploy power management settings to Windows clients through scripts"
categories: ["Windows", "Servers"]
date: "2008-05-05T09:09:00+02:00"
lastmod: "2008-05-05T09:09:00+02:00"
tags: ["Windows", "Servers", "Scripts", "Power Management"]
toc: true
---

## Batch Script

Here is the small batch script:

```batch
Powercfg.exe /CREATE "myPowerScheme"
Powercfg.exe /CHANGE "myPowerScheme" /monitor-timeout-dc 20
Powercfg.exe /CHANGE "myPowerScheme" /monitor-timeout-ac 20
Powercfg.exe /CHANGE "myPowerScheme" /disk-timeout-dc 0
Powercfg.exe /CHANGE "myPowerScheme" /disk-timeout-ac 0
Powercfg.exe /CHANGE "myPowerScheme" /standby-timeout-dc 0
Powercfg.exe /CHANGE "myPowerScheme" /standby-timeout-ac 0
Powercfg.exe /CHANGE "myPowerScheme" /hibernate-timeout-dc 0
Powercfg.exe /CHANGE "myPowerScheme" /hibernate-timeout-ac 0
Powercfg.exe /SETACTIVE "myPowerScheme"
```

Here I just want the monitor to turn off automatically after 20 minutes.

## VBS Integration

For those who want to integrate it into the VBS script:

```vb
'*****************************************************************************
'### Fonction powermgmt ###
'crée un nouveau schema de gestion d'alimentation et l'applique à l'aide de l'éxecution d'un fichier batch
'Syntaxe : powermgmt
Function powermgmt

   Set objFSO = CreateObject("Scripting.FileSystemObject")
   Set WshShell = CreateObject("wscript.shell")
   WshShell.run "energy_save.bat", SH_WIDE

End Function
'*****************************************************************************
```
