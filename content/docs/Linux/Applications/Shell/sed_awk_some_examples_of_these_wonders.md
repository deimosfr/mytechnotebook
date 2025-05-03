---
weight: 999
url: "/Sed_\\&_Awk_\\:_Quelques_exemples_de_ces_merveilles/"
title: "Sed & Awk: Some Examples of These Wonders"
description: "Examples and usage patterns for Sed, Awk, and Grep commands in Linux for text manipulation, filtering, and processing."
categories: ["Linux", "Database"]
date: "2014-01-01T15:45:00+02:00"
lastmod: "2014-01-01T15:45:00+02:00"
tags:
  [
    "sed",
    "awk",
    "grep",
    "text processing",
    "command line",
    "linux",
    "data manipulation",
  ]
toc: true
---

## Introduction

Sed & Awk, but what are they? These are two binaries for content filtering. For example, if you have a file and you want to output only specific lines, or even replace things in a file. Of course, what I'm telling you isn't extraordinary, but combined with regex, it becomes very powerful.

## Sed

Display certain lines of a large file:

```bash
sed -n "$LINE1,${LINE2}p;${LINEA2}q;" "$FILE"
```

Replace all occurrences of the word "stain" with the word "whiteness" in the file my_laundry, redirecting everything to the file my_clean_laundry:

```bash
sed s/tache/blancheur/g mon_linge > mon_linge_propre
```

Replace all line beginnings with a "." in the file my_laundry, redirecting everything to my_well_arranged_laundry:

```bash
sed s/^/./ > mon_linge_bien_range
```

Replace all line endings with a "." in the file my_laundry, redirecting everything to my_well_arranged_laundry:

```bash
sed s/$/./ > mon_linge_bien_range
```

Delete all lines containing the word "stain", redirecting everything to the file my_clean_laundry:

```bash
sed s/.*tache.*/\\0/ > mon_linge_propre
```

To delete lines matching a pattern (here we delete empty lines):

```bash
sed '/^$/d' monfichier
```

Delete all lines containing the word "stain" and replace all line endings with a ".", redirecting to my_good_smelling_laundry:

```bash
sed -e 's/$/./' -e  's/.*tache.*/\\0/' > mon_linge_qui_sent_bon
```

To make modifications directly without going through another file (on the fly):

```bash
sed -i 's/sysconfig/defaults/' cgconfig
```

For purists, we can also do a substitution in perl:

```bash
perl -pe 's/toto/tata/g' < monfichier
```

Replace the content of a file on the fly:

```bash
perl -pi -e "s/HOSTNAME=.*/HOSTNAME=node3.deimos.fr/" /etc/sysconfig/network
```

Or even with the replace command (which replaces in all \*.ini of the current folder):

```bash
replace -v from to from to "persistence=FILE" "persistence=BRIDGEPERSISTENCE" -- *.ini
```

or also:

```bash
find . -name "*.c" -exec sed -i "s/oldWord/newWord/g" '{}' \;
```

Delete the 10th line of a file:

```bash
sed -i '10d'
```

Delete x lines from a line number:

```bash
sed '1,55d' fichier.old > fichier.new
```

This will delete lines 1 to 55

## Awk

Display the 2nd column (PID) of the result of a ps aux:

```bash
ps aux|awk '{print $2}'
```

Display the 2nd column (PID) of the result of a ps aux, preceded by a label:

```bash
ps aux|awk '{printf("PID: %s\n"), $2}'
```

Display the 2nd column (PID) of the result of a ps aux, preceded by a label and showing the user:

```bash
ps aux|awk '{printf("PID: %s used by %s\n"), $2, $1}'
```

Display /etc/passwd well arranged:

```bash
awk -F: '{printf("User: %s\nUID: %s\nGID: %s\nName: %s\nHome: %s\nShell: %s\n", $1, $3, $4, $5, $6, $7)}' /etc/passwd
```

Use mathematical functions:

```bash
echo "3 4" | awk '{print int($1+$2)}'
```

(normally displays 7, if all goes well)
It is possible to have the equivalent of a grep and an awk in the same awk:

```bash
ip a | awk '/inet/{print $2}'
```

This allows me to get all IPs present on my machine.

For more examples: [AWK - the reference scripting language for file processing](/pdf/awk.pdf)

## Grep

Grep was not planned for the program, but here it is, I'll add it anyway :-)

To display only the root line in /etc/passwd:

```bash
grep ^root < /etc/passwd
```

or

```bash
cat /etc/passwd | grep ^root
```

But this last line calls another instruction (cat), which is not necessary.

To display all lines except those containing root:

```bash
grep -v root < /etc/passwd
```

Combine greps:

```bash
cat /etc/passwd | grep -e toto -e tata -e titi
```

It's possible to see the lines before or after the grep result. To see 2 lines before:

```bash
grep -A2 string filename
```

To see 2 lines after:

```bash
grep -B2 string filename
```

To see 2 lines before and after:

```bash
grep -C2 string filename
```

## FAQ

### Binary database matches

If you encounter this error during a grep, your file may be too large, so grep considers it binary. But if this is not the case, you can force it to read anyway:

```bash
grep -a monfichier
```

The -a bypasses the restriction

## Resources
- [Tutorial for Sed](/pdf/sed_tuto.pdf)
- [Awk cheat sheet](/pdf/awk.cheat.sheet.pdf)
