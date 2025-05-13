---
title: 'Prusa maintenance'
slug: prusa-maintenance/
description: 'Make sure your Prusa is always in good condition'
categories: ['Prusa', '3D Print']
tags: ['Prusa', '3D Print']
date: '2025-05-04T21:36:09+02:00'
---

## Introduction

This page is dedicated to the maintenance of your Prusa printer. It is important to keep your printer in good condition to ensure optimal performance and longevity.

I've been using my Prusa MK4, MK4S and now Core One for a while now, and I have learned a few things about maintenance that I would like to share with you.

## Regular Maintenance

### Every 600/800 hours

- Air blow the fans and everywhere else to avoid dust accumulation.
- Clean the heatbed with isopropyl alcohol.
- Clean the nozzle: wait 5min at 250°C, then use a brass brush to clean it.
- Clean every axes with antonomase.
- Lubricate the axes with a thin layer of grease.
- Check screws of the axes and the frame.
- Check the belts: they should be tight, but not too much: [Prusa belt tensioning](https://belt.connect.prusa3d.com/).

Useful links:

- [MK4/S maintenance](https://help.prusa3d.com/article/regular-printer-maintenance-mk4-s-mk3-9-s_419000)
- [MMU3 maintenance](https://help.prusa3d.com/fr/article/maintenance-reguliere-du-mmu3_682693)

## Issues

### Cloffed hotend

Perform an [assisted Cold Pull](https://help.prusa3d.com/article/cold-pull-26702-mk4s-13702-mk4-28702-mk3-9s-21702-mk3-9-17702-xl_445071) or do it manually, it works well. You can also check this:

- https://help.prusa3d.com/article/clogged-hotend-mk4_411823

### X or Y layer shift

If you encounter a layer shift, check:

- the belts: they should be tight
- the pulleys: they should be tight
- the screws: they should be tight
- the axes: they should be clean and lubricated
- it can come from the slicer and piece to print: try the gyroid infill

Useful links:

- https://help.prusa3d.com/article/layer-shifting_2020
