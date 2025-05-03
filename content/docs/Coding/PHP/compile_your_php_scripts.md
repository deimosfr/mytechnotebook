---
weight: 999
url: "/Compiler_vos_scripts_PHP/"
title: "Compile Your PHP Scripts"
description: "How to compile PHP scripts for both protection and performance improvements"
categories: ["Development", "Web", "PHP"]
date: "2013-05-08T18:54:00+02:00"
lastmod: "2013-05-08T18:54:00+02:00"
tags: ["php", "compilation", "bcompiler", "security", "optimization"]
toc: true
---

![PHP](/images/php_icon.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5.3 |
| **Operating System** | Debian 6 |
| **Website** | [PHP Website](https://www.php.net) |
| **Last Update** | 08/05/2013 |
{{< /table >}}

## Introduction

Sooner or later, we all face issues related to protecting the fruit of our intellectual labor. Here's a way to protect your precious PHP code and improve the execution time of your PHP scripts.

## Installation

The compiler is only provided in a development environment, which requires some installation steps before use:

```bash
aptitude update
aptitude upgrade
aptitude install make
aptitude install php5-dev
pecl install bcompiler
```

Then we'll modify the php.ini file. We need to add extension=bcompiler.so to the end of your php.ini file:

```bash
cp /etc/php5/apache2/php.ini /etc/php5/apache2/php.ini.old
echo "extension=bcompiler.so" >> /etc/php5/apache2/php.ini
```

## Configuration

```bash
> php --php-ini /etc/php5/apache2/php.ini -r "bcompiler_write_header();"
OK : PHP Warning: bcompiler_write_header() expects at least 1 parameter, 0 given in Command line code on line 1
KO : PHP Fatal error: Call to undefined function bcompiler_write_header() in Command line code on line 1
```

If you get a KO, you need to redo the procedure and make sure that pear-php is properly installed.

## Compiling PHP Code

Create a source file:

```php
<?php
echo "it codes and decodes";
?>
```

Then create a file called compiler.php:

```php
<?php
// the path of the bytecode file that will be created later
$bytecode = "bytecode.php";

// The source file
$codesource = "code.php";

// creation of the compiled file
$fichierbytecode = fopen($bytecode, "w");

// writing the file header;
bcompiler_write_header($fichierbytecode);

// writing the body of the file:
bcompiler_write_file($fichierbytecode, $codesource);

// writing the footer of the file:
bcompiler_write_footer($fichierbytecode);

?>
```

Run the PHP interpreter in CLI:

```bash
php --php-ini /etc/php5/apache2/php.ini compiler.php
```

Replace the links on your html index.php page from "code.php" to the bytecode.php file. The result doesn't change, but the recipe is unknown :-)

## References

[PHP Bcompiler documentation](https://php.net/manual/en/book.bcompiler.php)
