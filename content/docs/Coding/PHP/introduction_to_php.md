---
weight: 999
url: "/Introduction_au_PHP/"
title: "Introduction to PHP"
description: "A comprehensive introduction to PHP programming language covering basic syntax, variables, functions, forms, databases, and more."
categories: ["Development", "Database"]
date: "2010-12-27T21:26:00+02:00"
lastmod: "2010-12-27T21:26:00+02:00"
toc: true
---

# Introduction

[PHP](https://fr.wikipedia.org/wiki/Php) (recursive acronym for PHP: Hypertext Preprocessor) is a free scripting language primarily used to produce dynamic web pages via an HTTP server, but it can also function like any interpreted language locally by executing programs on the command line. PHP is an imperative language that has had complete object-oriented model capabilities since version 5. Due to the richness of its library, PHP is sometimes referred to as a platform rather than just a language.

# Basics

- In pages containing mainly HTML, it is possible to insert small pieces of PHP to make life easier. Just use tags like:

```html
<html>
  <head>
    <title>PHP says hello</title>
  </head>
  <body>
    <b>
      <?php
print "Hello, World!";
?>
    </b>
  </body>
</html>
```

This will produce:

```php
<html>
<head><title>PHP says hello</title></head>
<body>
<b>
Hello, World!
</b>
</body>
</html>
```

- We can also reverse the example above to put HTML code in a PHP page:

```php
<?php
print <<<_HTML_
<form method="post" action="$_SERVER[PHP_SELF]">
 Your Name: <input type="text" name="user">
 <br/>
 <input type="submit" value="Say Hello">
 </form>
_HTML_;
?>
```

The `$_SERVER[PHP_SELF]` variable is special to PHP. It is used to get the URL of the called link but without the protocol or the hostname. For example:

```
https://www.deimos.fr/blocnotesinfo/index.php?title=Introduction_au_PHP
```

Only what is in bold is kept.

To give you an idea of what's available, you can create a file with this:

```bash
<?php
echo "<table border=\"1\">";
echo "<tr><td>" .$_SERVER['argv'] ."</td><td>argv</td></tr>";
echo "<tr><td>" .$_SERVER['argc'] ."</td><td>argc</td></tr>";
echo "<tr><td>" .$_SERVER['GATEWAY_INTERFACE'] ."</td><td>GATEWAY_INTERFACE</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_ADDR'] ."</td><td>SERVER_ADDR</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_NAME'] ."</td><td>SERVER_NAME</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_SOFTWARE'] ."</td><td>SERVER_SOFTWARE</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_PROTOCOL'] ."</td><td>SERVER_PROTOCOL</td></tr>";
echo "<tr><td>" .$_SERVER['REQUEST_METHOD'] ."</td><td>REQUEST_METHOD</td></tr>";
echo "<tr><td>" .$_SERVER['REQUEST_TIME'] ."</td><td>REQUEST_TIME</td></tr>";
echo "<tr><td>" .$_SERVER['QUERY_STRING'] ."</td><td>QUERY_STRING</td></tr>";
echo "<tr><td>" .$_SERVER['DOCUMENT_ROOT'] ."</td><td>DOCUMENT_ROOT</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_ACCEPT'] ."</td><td>HTTP_ACCEPT</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_ACCEPT_CHARSET'] ."</td><td>HTTP_ACCEPT_CHARSET</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_ACCEPT_ENCODING'] ."</td><td>HTTP_ACCEPT_ENCODING</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_ACCEPT_LANGUAGE'] ."</td><td>HTTP_ACCEPT_LANGUAGE</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_CONNECTION'] ."</td><td>HTTP_CONNECTION</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_HOST'] ."</td><td>HTTP_HOST</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_REFERER'] ."</td><td>HTTP_REFERER</td></tr>";
echo "<tr><td>" .$_SERVER['HTTP_USER_AGENT'] ."</td><td>HTTP_USER_AGENT</td></tr>";
echo "<tr><td>" .$_SERVER['HTTPS'] ."</td><td>HTTPS</td></tr>";
echo "<tr><td>" .$_SERVER['REMOTE_ADDR'] ."</td><td>REMOTE_ADDR</td></tr>";
echo "<tr><td>" .$_SERVER['REMOTE_HOST'] ."</td><td>REMOTE_HOST</td></tr>";
echo "<tr><td>" .$_SERVER['REMOTE_PORT'] ."</td><td>REMOTE_PORT</td></tr>";
echo "<tr><td>" .$_SERVER['SCRIPT_FILENAME'] ."</td><td>SCRIPT_FILENAME</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_ADMIN'] ."</td><td>SERVER_ADMIN</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_PORT'] ."</td><td>SERVER_PORT</td></tr>";
echo "<tr><td>" .$_SERVER['SERVER_SIGNATURE'] ."</td><td>SERVER_SIGNATURE</td></tr>";
echo "<tr><td>" .$_SERVER['PATH_TRANSLATED'] ."</td><td>PATH_TRANSLATED</td></tr>";
echo "<tr><td>" .$_SERVER['SCRIPT_NAME'] ."</td><td>SCRIPT_NAME</td></tr>";
echo "<tr><td>" .$_SERVER['REQUEST_URI'] ."</td><td>REQUEST_URI</td></tr>";
echo "<tr><td>" .$_SERVER['PHP_AUTH_DIGEST'] ."</td><td>PHP_AUTH_DIGEST</td></tr>";
echo "<tr><td>" .$_SERVER['PHP_AUTH_USER'] ."</td><td>PHP_AUTH_USER</td></tr>";
echo "<tr><td>" .$_SERVER['PHP_AUTH_PW'] ."</td><td>PHP_AUTH_PW</td></tr>";
echo "<tr><td>" .$_SERVER['AUTH_TYPE'] ."</td><td>AUTH_TYPE</td></tr>";
echo "</table>"
?>
```

- For HTML forms corresponding to this:

```php
<form method="POST" action="sayhello.php">
Your Name: <input type="text" name="user">
<br/>
<input type="submit" value="Say Hello">
</form>
```

This will produce in PHP:

```php
<?php
print "Hello, ";
// Print what was submitted in the form parameter called 'user'
print $_POST['user'];
print "!";
?>
```

**$\_POST** contains the values of the form parameters as they were submitted. In programming terms, it's a variable named as such because you can modify the values it contains. In fact, it's a **array** variable because it can hold multiple values simultaneously.

When $\_SERVER[PHP_SELF] is the form's action, you can place both the form display code and the form processing code on the same page:

```php
<?php
// Print a greeting if the form was submitted
if ($_POST['user']) {
    print "Hello, ";
    // Print what was submitted in the form parameter called 'user'
    print $_POST['user'];
    print "!";
} else {
    // Otherwise, print the form
    print <<<_HTML_
<form method="post" action="$_SERVER[PHP_SELF]">
Your Name: <input type="text" name="user">
<br/>
<input type="submit" value="Say Hello">
</form>
_HTML_;
}
?>
```

- For comments, you can use:

  - `//`: For a single line to comment
  - `/* to */`: To comment a paragraph, we will come back to this later.

- Now here is an example of PHP code that connects to a database server and retrieves a list of dishes with their prices based on the value of the "rpas" parameter, then displays these dishes and their prices in an HTML table:

```php
<?php
require 'DB.php';
// Connect to MySQL running on localhost with username "menu"
// and password "good2eaT", and database "dinner"
$db = DB::connect('mysql://menu:good2eaT@localhost/dinner');
// Define what the allowable meals are
$meals = array('breakfast','lunch','dinner');
// Check if submitted form parameter "meal" is one of
// "breakfast", "lunch", or "dinner"
if (in_array($meals, $_POST['meal'])) {
    // If so, get all of the dishes for the specified meal
    $q = $dbh->query("SELECT dish,price FROM meals WHERE meal LIKE '" .
                     $_POST['meal'] ."'");
    // If no dishes were found in the database, say so
    if ($q->numrows == 0) {
        print "No dishes available.";
    } else {
        // Otherwise, print out each dish and its price as a row
        // in an HTML table
        print '<table><tr><th>Dish</th><th>Price</th></tr>';
        while ($row = $q->fetchRow( )) {
            print "<tr><td>$row[0]</td><td>$row[1]</td></tr>";
        }
        print "</table>";
    }
} else {
    // This message prints if the submitted parameter "meal" isn't
    // "breakfast", "lunch", or "dinner"
    print "Unknown meal.";
}
?>
```

# Manipulation of textual and numerical data

Strings can be delimited by "". They work like strings surrounded by ', but allow more special characters:

{{< table "table-hover table-striped" >}}
| Character | Meaning |
|-----------|---------|
| \n | Newline (ASCII 10) |
| \r | Carriage return (ASCII 13) |
| \t | Tab (ASCII 9) |
| \\ | \ |
| \$ | $ |
| \" | " |
| \0 .. \777 | Octal (base 8) number |
| \x0 .. \xFF | Hexadecimal (base 16) number |
{{< /table >}}

A small example to clarify everything. Imagine $user is Pierre:

- 'Hello $user': Will display: Hello $user
- "Hello $user": Will display: Hello Pierre

To combine 2 strings, use a "." **which allows concatenation of strings**:

```
print 'bread' . 'fruit';
print "It's a beautiful day " . 'in the neighborhood.';
print "The price is: " . '$3.95';
print 'Inky' . 'Pinky' . 'Blinky' . 'Clyde';
```

This will display:

```
breadfruit
It's a beautiful day in the neighborhood.
The price is: $3.95
InkyPinkyBlinkyClyde
```

## Text manipulation

### String validation

The **trim()** function removes spaces at the beginning or end of a string. Combined with **strlen()**, which tells you the length of a string, you can know the length of a value submitted by a form while ignoring leading and trailing spaces.

```
// $_POST['zipcode'] holds the value of the submitted form parameter
// "zipcode"
$zipcode = trim($_POST['zipcode']);
// Now $zipcode holds that value, with any leading or trailing spaces
// removed
$zip_length = strlen($zipcode);
// Complain if the ZIP code is not 5 characters long
if ($zip_length != 5) {
   print "Please enter a ZIP code that is 5 characters long.";
}
```

You can do it shorter by combining functions. Similar to Final Fantasy when you combine mana:

```
if (strlen(trim($_POST['zipcode'])) != 5) {
   print "Please enter a ZIP code that is 5 characters long.";
}
```

To compare strings with the equality operator, use (==):

```
if ($_POST['email'] == 'president@whitehouse.gov') {
  print "Welcome, Mr. President.";
}
```

To compare strings without taking case into account, use the strcasecmp() function. It returns 0 if the 2 strings provided to strcasecmp() are equal (regardless of case):

```
if (strcasecmp($_POST['email'], 'president@whitehouse.gov') == 0) {
   print "Welcome back, Mr. President.";
}
```

### Text formatting

printf() gives more control than print over the output appearance:

```
$price = 5; $tax = 0.075;
printf('The dish costs $%.2f', $price * (1 + $tax));
```

This will display:

```
The dish costs $5.38
```

The "%.2f" is replaced by the value of $price \* (1 + $tax)) and formatted to have 2 digits after the decimal point.

Format rules start with %. You can then place optional modifiers that affect the rule's behavior:

- A fill character: If the string replacing the format rule is too short, this character will fill it. Use a space to fill with spaces or a 0 to fill with zeros.
- A sign: For numbers, a + will place a + before positive numbers (they are normally displayed without a sign). For strings, a - will make the string right-justified (by default, strings are left-justified).
- A minimum width: the minimum size that the value replacing the format rule should have. If it's shorter, the fill character will be used to fill the void.
- A point and a number of decimal places: For floating point numbers, this controls the number of digits after the decimal point. In the example above, the .2 formats $price \* (1 + $tax)); with 2 decimal places.

Here's an example of zero-padding with printf():

```
$zip = '6520';
$month = 2;
$day = 6;
$year = 2007;
printf("ZIP is %05d and the date is %02d/%02d/%d", $zip, $month, $day, $year);
```

Will display:

```
ZIP is 06520 and the date is 02/06/2007
```

Displaying signs with printf():

```
$min = -40;
$max = 40;
printf("The computer can operate between %+d and %+d degrees Celsius.", $min, $max);
```

Will display:

```
The computer can operate between -40 and +40 degrees Celsius.
```

To discover other rules with printf, visit http://www.php.net

The functions **strtolower()** and **strtoupper()** produce, respectively, all lowercase or all uppercase versions of a string.

```
print strtolower('Beef, CHICKEN, Pork, duCK');
print strtoupper('Beef, CHICKEN, Pork, duCK');
```

Will display:

```
beef, chicken, pork, duck
BEEF, CHICKEN, PORK, DUCK
```

The function **ucwords()** capitalizes the first letter of each word in a string:

```
print ucwords(strtolower('JOHN FRANKENHEIMER'));
```

Will display:

```
John Frankenheimer
```

Truncating a string with **substr()**:

```
// Grab the first thirty characters of $_POST['comments']
print substr($_POST['comments'], 0, 30);
// Add an ellipsis
print '...';
```

Will display:

```
The Fresh Fish with Rice Noodle was delicious, but I didn't like the Beef Tripe.
```

The three parameters of substr() are respectively, the string concerned, the starting position of the substring to extract and the number of characters to extract.
The string starts at position 0, not 1: substr($\_POST['comments'], 0, 30) therefore means "extract 30 characters from $\_POST['comments'] starting from the beginning of this string".

Extracting the end of a string with substr():

```
print 'Card: XX';
print substr($_POST['card'],-4,4);
```

If the form parameter is worth 4000-1234-5678-9101, this will display:

```
Card: XX9101
```

So this is a very practical example for credit cards.

For brevity, you can use substr($_POST['card'],-4) instead of substr($\_POST['card'], -4,4). If you don't provide the last parameter, substr() will return everything between the starting position (whether positive or negative) and the end of the string.

Using **str_replace()**:

```php
print str_replace('{class}',$my_class,
                  '<span class="{class}">Fried Bean Curd<span>
                   <span class="{class}">Oil-Soaked Fish</span>');
```

Will display:

```php
<span class="lunch">Fried Bean Curd<span>
<span class="lunch">Oil-Soaked Fish</span>
```

## Arithmetic operators

Some basic operations in PHP:

```
print 2 + 2;
print 17 - 3.5;
print 10 / 3;
print 6 * 9;
```

Will display:

```
4
13.5
3.3333333333333
54
```

In addition to the +, -, /, and \* signs, PHP has the modulo % which returns the remainder of a division:

```
print 17 % 3;
```

will display:

```
2
```

## Variables

Here are examples of acceptable variables:

```
$size
$drinkSize
$my_drink_size
$_drinks
$drink4you2
```

And unacceptable:

{{< table "table-hover table-striped" >}}
| Variable name | Flaw |
|--------------|------|
| Begins with a number | $2hot4u |
| Unacceptable character: - | $drink-size |
| Unacceptable characters: @ and . | $drinkmaster@example.com |
| Unacceptable character:! | $drink!lots |
| Unacceptable character: + | $drink+dinner |
{{< /table >}}

**Variable names are case sensitive!!!**

## Operations on variables

Some operations on variables:

```php
<?php
$price = 3.95;
$tax_rate = 0.08;
$tax_amount = $price * $tax_rate;
$total_cost = $price + $tax_amount;
$username = 'james';
$domain = '@example.com';
$email_address = $username . $domain;
print 'The tax is ' . $tax_amount;
print "\n"; // this prints a linebreak
print 'The total cost is ' .$total_cost;
print "\n"; // this prints a linebreak
print $email_address;
?>
```

Addition combined with assignment:

```php
// Add 3 the regular way
$price = $price + 3;
// Add 3 with the combined operator
$price += 3;
```

Incrementation and decrementation:

```php
// Add one to $birthday
$birthday = $birthday + 1;
// Add another one to $birthday
++$birthday;
// Subtract 1 from $years_left
$years_left = $years_left - 1;
// Subtract another 1 from $years_left
```

## Variables placed in strings

In-place document interpolation:

```php
$page_title = 'Menu';
$meat = 'pork';
$vegetable = 'bean sprout';
print <<<MENU
<html>
<head><title>$page_title</title></head>
<body>
<ul>
<li> Barbecued $meat
<li> Sliced $meat
<li> Braised $meat with $vegetable
</ul>
</body>
</html>
MENU;
```

Will display:

```php
<html>
<head><title>Menu</title></head>
<body>
<ul>
<li> Barbecued pork
<li> Sliced pork
<li> Braised pork with bean sprout
</ul>
</body>
```

Interpolation with braces:

```php
$preparation = 'Braise';
$meat = 'Beef';
print "{$preparation}d $meat with Vegetables";
```

Will display:

```
Braised Beef with Vegetables
```

Without the braces, the print instruction in the example above would have been:

```
print "$preparationd $meat with Vegetables";
```

# Decision making and repetitions

In this chapter, we'll see how to:

- Display a special menu if a user with admin rights is logged in.
- Display a different page header depending on the time of day.
- Notify a user if they've received new messages since their last connection.

When making decisions, the PHP interpreter reduces an expression to a true or false value.

The value of an assignment is the assigned value. The expression $price = 5 equals 5: this price has been assigned to $price. Since assignment produces a result, you can chain assignments to assign the same value to multiple variables:

```
$price = $quantity = 5;
```

## Decision making

Using **elseif()**:

```php
if ($logged_in) {
    // This runs if $logged_in is true
    print "Welcome aboard, trusted user.";
} elseif ($new_messages) {
    // This runs if $logged_in is false but $new_messages is true
    print "Dear stranger, there are new messages.";
} elseif ($emergency) {
    // This runs if $logged_in and $new_messages are false
    // But $emergency is true
    print "Stranger, there are no new messages, but there is an emergency.";
} else {
    // You can put what you want
}
```

## Creating complex decisions

The equality operator:

```
if ($new_messages == 10) {
   print "You have ten new messages.";
}
if ($new_messages == $max_messages) {
   print "You have the maximum number of messages.";
}
if ($dinner == 'Braised Scallops') {
   print "Yum! I love seafood.";
}
```

The difference operator:

```
if ($new_messages != 10) {
   print "You don't have ten new messages.";
}
if ($dinner != 'Braised Scallops') {
   print "I guess we're out of scallops.";
}
```

**Be careful not to use = when you want to use ==. A single equals sign assigns a value and returns the assigned value, while two equal signs test equality and return true if the values are equal. If you forget the second equals sign, you will usually get an if() test that will always be true:**

Example of an assignment that should have been a comparison:

```
if ($new_messages = 12) {
   print "It seems you now have twelve new messages.";
}
```

**One way to avoid inadvertently using = instead of == is to place the variable on the right side of the comparison and the literal on the left:**

```
if (12 == $new_messages) {
   print "You have twelve new messages.";
}
```

Less than and greater than operators:

```php
if ($age > 17) {
    print "You are old enough to download the movie.";
}
if ($age >= 65) {
    print "You are old enough for a discount.";
}
if ($celsius_temp <= 0) {
    print "Uh-oh, your pipes may freeze.";
}
if ($kelvin_temp < 20.3) {
    print "Your hydrogen is a liquid or a solid now.";
}
```

Floating-point numbers have an internal representation that may be slightly different from the assigned value (the internal representation of 50.0, for example, might be 50.000000002). To test if 2 floating-point numbers are equal, check if the difference between these 2 numbers is less than a reasonably small threshold, instead of using the equality operator. If you compare monetary values, for example, an acceptable threshold would be 0.00001.

Floating-point number comparisons:

```php
if(abs($price_1 - $price_2) < 0.00001) {
    print '$price_1 and $price_2 are equal.';
} else {
    print '$price_1 and $price_2 are not equal.';
}
```

The **abs()** function returns the absolute value of its parameter. With abs(), the comparison works correctly, whether $price_1 is greater than $price_2 or not.

**The ASCII codes for digits are less than those for uppercase letters, which are themselves less than the codes for lowercase letters; the codes for accented letters are in the "extended" part, thus higher than all codes in international ASCII.**

Generally, strings are compared alphabetically (caramel < chocolate).

String comparison:

```php
if ($word < 'baa') {
    print "Your word isn't cookie.";
}
if ($word >= 'zoo') {
    print "Your word could be zoo or zymurgy, but not zone.";
}
```

If you want to make sure the PHP interpreter compares strings without performing numeric conversion behind the scenes, use the **strcmp()** function which always compares its parameters according to the ASCII table order.

The **strcmp()** function takes two strings as parameters. It returns a positive number if the first string is greater than the second or a negative number if the first string is less than the second. The order is that of the extended ASCII table. This function returns 0 if the 2 strings are equal.

String comparisons with strcmp():

```php
$x = strcmp("x54321","x5678");
if ($x > 0) {
    print 'The string "x54321" is greater than the string "x5678".';
} elseif ($x < 0) {
    print 'The string "x54321" is less than the string "x5678".';
}
```

In the example below, strcmp() finds that the string "54321" is less than "5678" because the second characters of the 2 strings differ and "4" comes before "6". **The alphabetical order is all that matters to the strcmp() function, no matter that numerically it's the opposite**.

```php
// These values are compared using numeric order
if ("54321" > "5678") {
    print 'The string "54321" is greater than the string "5678".';
} else {
    print 'The string "54321" is not greater than the string "5678".';
}
```

Using the negation operator:

```
if (! strcasecmp($first_name,$last_name)) {
   print '$first_name and $last_name are equal.';
}
```

With logical operators, you can combine multiple expressions in the same if() statement. The logical AND operator, &, tests if 2 expressions are both true.
The logical OR operator, ||, tests if at least one of the 2 expressions is true:

```php
if (($age >= 13) && ($age < 65)) {
   print "You are too old for a kid's discount and too young for the senior's discount.";
}
if (($meal == 'breakfast') || ($dessert == 'souffle')) {
   print "Time to eat some eggs.";
}
```

## Repetitions

Producing a <select> menu with while():

```php
$i = 1;
print '<select name="people">';
while ($i <= 10) {
    print "<option>$i</option>\n";
    $i++;
}
print '</select>';
```

Will display:

```php
<select name="people"><option>1</option>
<option>2</option>
<option>3</option>
<option>4</option>
<option>5</option>
<option>6</option>
<option>7</option>
<option>8</option>
<option>9</option>
<option>10</option>
</select>
```

Producing a <select> menu with for():

```php
print '<select name="people">';
for ($i = 1; $i <= 10; $i++) {
    print "<option>$i</option>";
}
print '</select>';
```

The third expression ($i++) is the iteration expression. It is executed after each execution of the loop body.

Multiple expressions in for():

```php
print '<select name="doughnuts">';
for ($min = 1, $max = 10; $min < 50; $min += 10, $max += 10) {
    print "<option>$min - $max</option>\n";
}
print '</select>';
```

# Using arrays

In the examples that follow, here's the typical array we'll be working with:
![PHP tab](/images/php_tab.avif)

## Creating an array

- Creating arrays:

```php
// An array called $vegetables with string keys
$vegetables['corn'] = 'yellow';
$vegetables['beet'] = 'red';
$vegetables['carrot'] = 'orange';
// An array called $dinner with numeric keys
$dinner[0] = 'Sweet Corn and Asparagus';
$dinner[1] = 'Lemon Chicken';
$dinner[2] = 'Braised Bamboo Fungus';
// An array called $computers with numeric and string keys
$computers['trs-80'] = 'Radio Shack';
$computers[2600] = 'Atari';
$computers['Adam'] = 'Coleco';
```

Array keys and values are written exactly like other strings and numbers in a PHP program: **with apostrophes around strings, but not around numbers.**

- Creating arrays with **array()**:

```php
$vegetables = array('corn' => 'yellow',
                    'beet' => 'red',
                    'carrot' => 'orange');
$dinner = array(0 => 'Sweet Corn and Asparagus',
                1 => 'Lemon Chicken',
                2 => 'Braised Bamboo Fungus');
$computers = array('trs-80' => 'Radio Shack',
                   2600 => 'Atari',
                   'Adam' => 'Coleco');
```

With **array()**, you specify a list of key/value pairs separated by commas. The key and value are separated by =>. The array() syntax is more concise when adding multiple elements to an array at once. The bracket syntax is preferable when adding elements one by one.

## Creating an array indexed by numbers

If you create an array with array() by specifying only a list of values instead of key/value pairs, the PHP interpreter automatically assigns a numeric key to each value. The keys start at 0 and increase by 1 for each element. The example below uses this technique to create the $dinner array:

```php
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
print "I want $dinner[0] and $dinner[1].";
```

Will display:

```
I want Sweet Corn and Asparagus and Lemon Chicken.
```

PHP automatically uses incremented numbers for keys when creating an array or adding elements to an array with the empty brackets syntax:

```php
// Create $lunch array with two elements
// This sets $lunch[0]
$lunch[] = 'Dried Mushrooms in Brown Sauce';
// This sets $lunch[1]
$lunch[] = 'Pineapple and Yu Fungus';
// Create $dinner with three elements
$dinner = array('Sweet Corn and Asparagus', 'Lemon Chicken',
                'Braised Bamboo Fungus');
// Add an element to the end of $dinner
// This sets $dinner[3]
$dinner[] = 'Flank Skin with Spiced Flavor';
```

**If the array doesn't exist yet, empty brackets create it by adding an element with key 0.**

## Finding the size of an array

The **count()** function returns the number of elements in an array:

```php
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
$dishes = count($dinner);
print "There are $dishes things for dinner.";
```

Will display:

```
There are 3 things for dinner.
```

When passed an empty array (i.e., an array containing no elements), count() returns 0. An empty array is also evaluated as false in a test expression.

## Traversing arrays

- Traversal with foreach():

```php
$meal = array('breakfast' => 'Walnut Bun',
              'lunch' => 'Cashew Nuts and White Mushrooms',
              'snack' => 'Dried Mulberries',
              'dinner' => 'Eggplant with Chili Sauce');
print "<table>\n";
foreach ($meal as $key => $value) {
    print "<tr><td>$key</td><td>$value</td></tr>\n";
}
print '</table>';
```

Will display:

```php
<table>
<tr><td>breakfast</td><td>Walnut Bun</td></tr>
<tr><td>lunch</td><td>Cashew Nuts and White Mushrooms</td></tr>
<tr><td>snack</td><td>Dried Mulberries</td></tr>
<tr><td>dinner</td><td>Eggplant with Chili Sauce</td></tr>
</table>
```

- To alternate table row colors:

```php
$row_color = array('red','green');
$color_index = 0;
$meal = array('breakfast' => 'Walnut Bun',
              'lunch' => 'Cashew Nuts and White Mushrooms',
              'snack' => 'Dried Mulberries',
              'dinner' => 'Eggplant with Chili Sauce');
print "<table>\n";
foreach ($meal as $key => $value) {
    print '<tr bgcolor="' . $row_color[$color_index] . '">';
    print "<td>$key</td><td>$value</td></tr>\n";
    // This switches $color_index between 0 and 1
    $color_index = 1 - $color_index;
}
print '</table>';
```

Will display:

```php
<table>
<tr bgcolor="red"><td>breakfast</td><td>Walnut Bun</td></tr>
<tr bgcolor="green"><td>lunch</td><td>Cashew Nuts and White Mushrooms</td></tr>
<tr bgcolor="red"><td>snack</td><td>Dried Mulberries</td></tr>
<tr bgcolor="green"><td>dinner</td><td>Eggplant with Chili Sauce</td></tr>
</table>
```

- Modifying an array with foreach():

```php
$meals = array('Walnut Bun' => 1,
               'Cashew Nuts and White Mushrooms' => 4.95,
               'Dried Mulberries' => 3.00,
               'Eggplant with Chili Sauce' => 6.50);
foreach ($meals as $dish => $price) {
    // $price = $price * 2 does NOT work
    $meals[$dish] = $meals[$dish] * 2;
}
// Iterate over the array again and print the changed values
foreach ($meals as $dish => $price) {
    printf("The new price of %s is \$%.2f.\n",$dish,$price);
}
```

- Using foreach() with indexed arrays:

```php
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
foreach ($dinner as $dish) {
    print "You can eat: $dish\n";
}
```

- Traversing an indexed array with for():

```php
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
for ($i = 0, $num_dishes = count($dinner); $i < $num_dishes; $i++) {
  print "Dish number $i is $dinner[$i]\n";
}
```

- Alternating table row colors with for():

```php
$row_color = array('red','green');
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
print "<table>\n";
for ($i = 0, $num_dishes = count($dinner); $i < $num_dishes; $i++) {
    print '<tr bgcolor="' . $row_color[$i % 2] . '">';
    print "<td>Element $i</td><td>$dinner[$i]</td></tr>\n";
}
print '</table>';
```

- To ensure access to elements in **numeric order** of their keys, use for():

```php
for ($i = 0, $num_letters = count($letters); $i < $num_letters; $i++) {
    print $letters[$i];
}
```

**To check if an element with a given key exists, use array_key_exists():**

```
meals = array('Walnut Bun' => 1,
               'Cashew Nuts and White Mushrooms' => 4.95,
               'Dried Mulberries' => 3.00,
               'Eggplant with Chili Sauce' => 6.50,
               'Shrimp Puffs' => 0); // Shrimp Puffs are free!
$books = array("The Eater's Guide to Chinese Characters",
               'How to Cook and Eat in Chinese');
// This is true
if (array_key_exists('Shrimp Puffs',$meals)) {
    print "Yes, we have Shrimp Puffs";
}
// This is false
if (array_key_exists('Steak Sandwich',$meals)) {
    print "We have a Steak Sandwich";
}
// This is true
if (array_key_exists(1, $books)) {
    print "Element 1 is How to Cook in Eat in Chinese";
}
```

- The **array_search()** function is similar to in_array(), but if it finds an element it returns its key rather than true. In the example below, array_search() returns the name of the dish that costs 6.50 euros:

```php
$meals = array('Walnut Bun' => 1,
               'Cashew Nuts and White Mushrooms' => 4.95,
               'Dried Mulberries' => 3.00,
               'Eggplant with Chili Sauce' => 6.50,
               'Shrimp Puffs' => 0);
$dish = array_search(6.50, $meals);
if ($dish) {
    print "$dish costs \$6.50";
}
```

## Modifying arrays

- Interpolation of array elements in strings with double apostrophes:

```php
$meals['breakfast'] = 'Walnut Bun';
$meals['lunch'] = 'Eggplant with Chili Sauce';
$amounts = array(3, 6);
print "For breakfast, I'd like $meals[breakfast] and for lunch, ";
print "I'd like $meals[lunch]. I want $amounts[0] at breakfast and ";
print "$amounts[1] at lunch.";
```

The **unset()** function allows you to remove an element from an array:

```
unset($dishes['Roast Duck']);
```

Deleting an element with unset() is different from simply assigning it 0 or an empty string. When you use unset(), the element is no longer there when traversing the array or counting the number of its elements.
Using unset() on an array representing an inventory is like saying that the store no longer offers a product. Setting the element's value to 0 or assigning an empty string means that this element is no longer in stock for the moment.

- Creating a string from an array with implode():

```
$dimsum = array('Chicken Bun','Stuffed Duck Web','Turnip Cake');
$menu = implode(', ', $dimsum);
print $menu;
```

Will display:

```
Chicken Bun, Stuffed Duck Web, Turnip Cake
```

To convert an array without adding a delimiter, use an empty string as the first parameter of implode():

```
$letters = array('A','B','C','D');
print implode(,$letters);
```

Will display:

```
ABCD
```

- Generating HTML table rows with implode():

```
$dimsum = array('Chicken Bun','Stuffed Duck Web','Turnip Cake');
print '<tr><td>' . implode('</td><td>',$dimsum) . '</td></tr>';
```

Will display:

```
<tr><td>Chicken Bun</td><td>Stuffed Duck Web</td><td>Turnip Cake</td></tr>
```

As you've probably seen, this avoids having to write a loop.

- Converting a string to an array with **explode()**:

```php
$fish = 'Bass, Carp, Pike, Flounder';
$fish_list = explode(', ', $fish);
print "The second fish is $fish_list[1]";
```

Will display:

```
The second fish is Carp
```

## Sorting arrays

The **sort()** function sorts an array by the values of its elements. It should only be used on arrays with numeric keys because it resets these keys during sorting. Here are some arrays before and after a call to sort().

Sorting with sort():

```php
$dinner = array('Sweet Corn and Asparagus',
                'Lemon Chicken',
                'Braised Bamboo Fungus');
$meal = array('breakfast' => 'Walnut Bun',
              'lunch' => 'Cashew Nuts and White Mushrooms',
              'snack' => 'Dried Mulberries',
              'dinner' => 'Eggplant with Chili Sauce');
print "Before Sorting:\n";
foreach ($dinner as $key => $value) {
    print " \$dinner: $key $value\n";
}
foreach ($meal as $key => $value) {
    print "   \$meal: $key $value\n";
}
sort($dinner);
sort($meal);
print "After Sorting:\n";
foreach ($dinner as $key => $value) {
    print " \$dinner: $key $value\n";
}
foreach ($meal as $key => $value) {
    print "   \$meal: $key $value\n";
}
```

Will display:

```
Before Sorting:
$dinner: 0 Sweet Corn and Asparagus
$dinner: 1 Lemon Chicken
$dinner: 2 Braised Bamboo Fungus
  $meal: breakfast Walnut Bun
  $meal: lunch Cashew Nuts and White Mushrooms
  $meal: snack Dried Mulberries
  $meal: dinner Eggplant with Chili Sauce
After Sorting:
$dinner: 0 Braised Bamboo Fungus
$dinner: 1 Lemon Chicken
$dinner: 2 Sweet Corn and Asparagus
  $meal: 0 Cashew Nuts and White Mushrooms
  $meal: 1 Dried Mulberries
  $meal: 2 Eggplant with Chili Sauce
  $meal: 3 Walnut Bun
```

- To sort an associative array, **use asort(), which preserves the keys and values**:

```php
$meal = array('breakfast' => 'Walnut Bun',
              'lunch' => 'Cashew Nuts and White Mushrooms',
              'snack' => 'Dried Mulberries',
              'dinner' => 'Eggplant with Chili Sauce');
print "Before Sorting:\n";
foreach ($meal as $key => $value) {
    print "   \$meal: $key $value\n";
}
asort($meal);
print "After Sorting:\n";
foreach ($meal as $key => $value) {
    print "   \$meal: $key $value\n";
}
```

Will display:

```php
Before Sorting:
   $meal: breakfast Walnut Bun
   $meal: lunch Cashew Nuts and White Mushrooms
   $meal: snack Dried Mulberries
   $meal: dinner Eggplant with Chili Sauce
After Sorting:
   $meal: lunch Cashew Nuts and White Mushrooms
   $meal: snack Dried Mulberries
   $meal: dinner Eggplant with Chili Sauce
   $meal: breakfast Walnut Bun
```

- While sort() and asort() sort arrays by the values of their elements, ksort() allows sorting by their keys: the key/value pairs remain identical, but are ordered by keys.

The functions **rsort(), arsort() and krsort()** are the respective counterparts of sort(), asort() and ksort() for sorting in descending order**. They work exactly the same way, except that the largest key or value (or the last from an alphabetical point of view) will appear first in the sorted array and the following elements will be placed in descending order.**

## Using multidimensional arrays

- Creating multidimensional arrays with array():

```php
$meals = array('breakfast' => array('Walnut Bun','Coffee'),
               'lunch'     => array('Cashew Nuts', 'White Mushrooms'),
               'snack'     => array('Dried Mulberries','Salted Sesame Crab'));
$lunches = array( array('Chicken','Eggplant','Rice'),
                  array('Beef','Scallions','Noodles'),
                  array('Eggplant','Tofu'));
$flavors = array('Japanese' => array('hot' => 'wasabi',
                                     'salty' => 'soy sauce'),
                 'Chinese' => array('hot' => 'mustard',
                                     'pepper-salty' => 'prickly ash'));
```

You access the elements of these arrays by using additional pairs of brackets to identify them: each pair goes down one level in the complete array.

Accessing elements of a multidimensional array:

```php
print $meals['lunch'][1]; // White Mushrooms
print $meals['snack'][0];           // Dried Mulberries
print $lunches[0][0];               // Chicken
print $lunches[2][1];               // Tofu
print $flavors['Japanese']['salty'] // soy sauce
print $flavors['Chinese']['hot'];   // mustard
```

- Manipulating multidimensional arrays:

```php
$prices['dinner']['Sweet Corn and Asparagus'] = 12.50;
$prices['lunch']['Cashew Nuts and White Mushrooms'] = 4.95;
$prices['dinner']['Braised Bamboo Fungus'] = 8.95;
$prices['dinner']['total'] = $prices['dinner']['Sweet Corn and Asparagus'] +
                             $prices['dinner']['Braised Bamboo Fungus'];
$specials[0][0] = 'Chestnut Bun';
$specials[0][1] = 'Walnut Bun';
$specials[0][2] = 'Peanut Bun';
$specials[1][0] = 'Chestnut Salad';
$specials[1][1] = 'Walnut Salad';
// Leaving out the index adds it to the end of the array
// This creates $specials[1][2]
$specials[1][] = 'Peanut Salad';
```

- Traversing a multidimensional array with foreach():

```php
$flavors = array('Japanese' => array('hot' => 'wasabi',
                                     'salty' => 'soy sauce'),
                 'Chinese' => array('hot' => 'mustard',
                                     'pepper-salty' => 'prickly ash'));
// $culture is the key and $culture_flavors is the value (an array)
foreach ($flavors as $culture => $culture_flavors) {
    // $flavor is the key and $example is the value
    foreach ($culture_flavors as $flavor => $example) {
        print "A $culture $flavor flavor is $example.\n";
    }
}
```

Will display:

```php
A Japanese hot flavor is wasabi.
A Japanese salty flavor is soy sauce.
A Chinese hot flavor is mustard.
A Chinese pepper-salty flavor is prickly ash.
```

The first foreach() loop traverses the first dimension of $flavors. The keys stored in $culture are the strings Japanese and Chinese and the values stored in $culture_flavors are the arrays that are the elements of this dimension.
The next foreach() traverses these arrays by copying keys like hot and salty into $flavors and values like wasabi and soy sauce into $example. The code block of the second foreach() uses the variables from both foreach() instructions to produce a complete message.

Just as nested foreach() loops traverse a multidimensional associative array, nested for() loops allow traversing a multidimensional array with numeric indices:

```php
$specials = array( array('Chestnut Bun', 'Walnut Bun', 'Peanut Bun'),
                   array('Chestnut Salad','Walnut Salad', 'Peanut Salad') );
// $num_specials is 2: the number of elements in the first dimension of $specials
for ($i = 0, $num_specials = count($specials); $i < $num_specials; $i++) {
    // $num_sub is 3: the number of elements in each sub-array
    for ($m = 0, $num_sub = count($specials[$i]); $m < $num_sub; $m++) {
        print "Element [$i][$m] is " . $specials[$i][$m] . "\n";
    }
}
```

Will display:

```php
Element [0][0]   is  Chestnut Bun
Element [0][1]   is  Walnut Bun
Element [0][2]   is  Peanut Bun
Element [1][0]   is  Chestnut Salad
Element [1][1]   is  Walnut Salad
Element [1][2]   is  Peanut Salad
```

To interpolate the value of a multidimensional array in a string between double apostrophes or an in-place document, we must use the syntax to produce the same result as the previous example: in fact, the only difference is the print instruction:

Interpolating a multidimensional array:

```php
$specials = array( array('Chestnut Bun', 'Walnut Bun', 'Peanut Bun'),
                   array('Chestnut Salad','Walnut Salad', 'Peanut Salad') );
// $num_specials is 2: the number of elements in the first dimension of $specials
for ($i = 0, $num_specials = count($specials); $i < $num_specials; $i++) {
    // $num_sub is 3: the number of elements in each sub-array
    for ($m = 0, $num_sub = count($specials[$i]); $m < $num_sub; $m++) {
        print "Element [$i][$m] is {$specials[$i][$m]}\n";
    }
}
```

# Functions

## Function declarations and calls

- Declaration of a function called **page_header()**:

```
function page_header() {
   print '<html><head><title>Welcome to my site</title></head>';
   print '<body bgcolor="#ffffff">';
}
```

- Calling a function:

```php
page_header();
print "Welcome, $user";
print "</body></html>";
```

- Function definitions before or after their call:

```php
function page_header( ) {
    print '<html><head><title>Welcome to my site</title></head>';
    print '<body bgcolor="#ffffff">';
}

page_header( );
print "Welcome, $user";
page_footer( );

function page_footer( ) {
    print '<hr>Thanks for visiting.';
    print '</body></html>';
}
```

**I strongly recommend you to code "clean". Indeed, if you place functions anywhere in your code, it will quickly become messy. You should group functions together and put them at the beginning of your code. I'm not a purist (...well, maybe), but people reviewing your code after you, or even yourself if you take it up a few months later, will appreciate being able to quickly understand the code.**

## Passing parameters to functions

- Declaration of a function with a parameter (here $color):

```
function page_header2($color) {
   print '<html><head><title>Welcome to my site</title></head>';
   print '<body bgcolor="#' . $color . '">';
}
```

By calling the function like this:

```
page_header2('cc00cc');
```

You will get:

```php
<html><head><title>Welcome to my site</title></head><body bgcolor="#cc00cc">
```

When you define a function that takes a parameter, you must pass a parameter to this function when calling it. Otherwise, the PHP interpreter will produce a warning message to complain about it. If for example, you call **page_header2()** as follows:

```
page_header2();
```

It will display:

```
PHP Warning:  Missing argument 1 for page_header2()
```

- Using a default value for a parameter:

**To avoid this warning**, make it so this function can be called without parameters by providing a default value in the function declaration. In this case, if the function is called without parameters, this default value will be taken as the parameter value. To provide a default value, place it after the parameter name. Here this value is cc3399:

```
function page_header3($color = 'cc3399') {
   print '<html><head><title>Welcome to my site</title></head>';
   print '<body bgcolor="#' . $color . '">';
}
```

Default parameter values must be literals, such as 12, cc3399, or a string; **they cannot be variables**:

```php
$my_color = '#000000';

// This is incorrect: the default value can't be a variable.
function page_header_bad($color = $my_color) {
    print '<html><head><title>Welcome to my site</title></head>';
    print '<body bgcolor="#' . $color . '">';
}
```

- Defining a function with **2 parameters**:

```
function page_header4($color, $title) {
   print '<html><head><title>Welcome to ' . $title . '</title></head>';
   print '<body bgcolor="#' . $color . '">';
}
```

- Calling a function with 2 parameters:

```
page_header4('66cc66','my homepage');
```

- Defining multiple optional parameters:

```
// One optional argument: it must be last
function page_header5($color, $title, $header = 'Welcome') {
   print '<html><head><title>Welcome to ' . $title . '</title></head>';
   print '<body bgcolor="#' . $color . '">';
   print "<h1>$header</h1>";
}
// Acceptable ways to call this function:
page_header5('66cc99','my wonderful page'); // uses default $header
page_header5('66cc99','my wonderful page','This page is great!');
// no defaults

// Two optional arguments: must be last two arguments
function page_header6($color, $title = 'the page', $header = 'Welcome') {
   print '<html><head><title>Welcome to ' . $title . '</title></head>';
   print '<body bgcolor="#' . $color . '">';
   print "<h1>$header</h1>";
}
// Acceptable ways to call this function:
page_header6('66cc99'); // uses default $title and $header
page_header6('66cc99','my wonderful page'); // uses default $header
page_header6('66cc99','my wonderful page','This page is great!');
// no defaults

// All optional arguments
function page_header6($color = '336699', $title = 'the page', $header = 'Welcome') {
   print '<html><head><title>Welcome to ' . $title . '</title></head>';
   print '<body bgcolor="#' . $color . '">';
   print "<h1>$header</h1>";
}
// Acceptable ways to call this function:
page_header7( ); // uses all defaults
page_header7('66cc99'); // uses default $title and $header
page_header7('66cc99','my wonderful page'); // uses default $header
page_header7('66cc99','my wonderful page','This page is great!');
// no defaults
```

## Return values from functions

Capturing a return value:

```
$number_to_display = number_format(285266237);
print "The population of the US is about: $number_to_display";
```

This will display:

```
The population of the US is about: 285,266,237
```

To return values from your own functions, use the return keyword, followed by the value to return. The execution of a function stops as soon as the return keyword is encountered and then returns the specified value. The example below defines a function that returns the total amount of a meal after adding VAT and tip. Returning a value from a function:

```php
function restaurant_check($meal, $tax, $tip) {
    $tax_amount = $meal * ($tax / 100);
    $tip_amount = $meal * ($tip / 100);
    $total_amount = $meal + $tax_amount + $tip_amount;

    return $total_amount;
}
```

A return statement can only return one value: syntax like return 15, 23 is not allowed. **If you want a function to return multiple values, place them in an array and return it**.

Here's a modified version of the restaurant_check() function that returns a 2-element array, corresponding to the total amount before and after adding the tip.

- **Returning an array from a function**:

```php
function restaurant_check2($meal, $tax, $tip) {
    $tax_amount  = $meal * ($tax / 100);
    $tip_amount  = $meal * ($tip / 100);
    $total_notip = $meal + $tax_amount;
    $total_tip   = $meal + $tax_amount + $tip_amount;
    return array($total_notip, $total_tip);
  }
```

- **Using an array returned by a function**:

```php
$totals = restaurant_check2(15.22, 8.25, 15);

if ($totals[0] < 20) {
    print 'The total without tip is less than $20.';
}
if ($totals[1] < 20) {
    print 'The total with tip is less than $20.';
}
```

- Using return values with an if():

```php
if (restaurant_check(15.22, 8.25, 15) < 20) {
    print 'Less than $20, I can pay cash.';
} else {
    print 'Too expensive, I need my credit card.';
}
```

- Functions returning true or false:

```php
function can_pay_cash($cash_on_hand, $amount) {
    if ($amount > $cash_on_hand) {
        return false;
    } else {
        return true;
    }
}

$total = restaurant_check(15.22,8.25,15);
if (can_pay_cash(20, $total)) {
    print "I can pay in cash.";
} else {
    print "Time for the credit card.";
}
```

From a function, there are 2 ways to access a global variable. The most direct method is to look them up in a special array called $GLOBALS, because any global variable is accessible as an element of this array. Here's how to use the $GLOBALS array:

```php
$dinner = 'Curry Cuttlefish';

function macrobiotic_dinner( ) {
    $dinner = "Some Vegetables";
    print "Dinner is $dinner";
    // Succumb to the delights of the ocean
    print " but I'd rather have ";
    print $GLOBALS['dinner'];
    print "\n";
}

macrobiotic_dinner( );
print "Regular dinner is: $dinner";
```

This will display:

```
Dinner is Some Vegetables but I'd rather have Curry Cuttlefish
Regular dinner is: Curry Cuttlefish
```

The 2nd way to access a global variable from a function is to use the _global_ keyword, which tells the PHP interpreter that **the subsequent use, in the function, of the indicated variable will refer to this local variable, and not to a local variable**: this is called "placing a variable in the local scope". The _global_ keyword:

```php
$dinner = 'Curry Cuttlefish';

function vegetarian_dinner( ) {
    global $dinner;
    print "Dinner was $dinner, but now it's ";
    $dinner = 'Sauteed Pea Shoots';
    print $dinner;
    print "\n";
}

print "Regular Dinner is $dinner.\n";
vegetarian_dinner( );
print "Regular dinner is $dinner";
```

Will display:

```
Regular Dinner is Curry Cuttlefish.
Dinner was Curry Cuttlefish, but now it's Sauteed Pea Shoots
Regular dinner is Sauteed Pea Shoots
```

You can also use global with multiple variable names, separated by commas:

```
global $dinner, $lunch, $breakfast;
```

# Creating HTML forms

Displaying a "Hello":

```php
if (array_key_exists('my_name',$_POST)) {
    print "Hello, " . $_POST['my_name'];
} else {
    print<<<_HTML_
<form method="post" action="$_SERVER[PHP_SELF]">
 Your name: <input type="text" name="my_name">
<br/>
<input type="submit" value="Say Hello">
</form>
_HTML_;
}
```

- Form retrieval and display:

![Form retrieval in PHP](/images/recup_form_php.avif)

- Form submission and result display:

![Form submission in PHP](/images/soum_form_php.avif)

The form is sent back to the same URL as the one that initially requested it because the action attribute of the <form> tag is initialized with the special variable $\_SERVER[PHP_SELF]. The super-global $\_SERVER array contains a number of pieces of information about your server as well as the request being processed by the PHP interpreter. The PHP_SELF element of this array, in particular, contains the path component of the URL of the request: if for example, a PHP script is accessed through http://www.example.com/magazin/catalogue.php, $\_SERVER['PHP_SELF'] will be /magasin/catalogue.php in this page.

The $\_POST array is a super-global variable containing the data submitted by the form. Its keys are the names of the form elements and the corresponding values are those of these elements. When clicking on the submit button, the value of $\_POST['my_name'] is initialized with what was entered in the field whose name attribute is my_name.

Thus, by testing if there is a my_name key in the $\_POST array, we can know if a parameter named my_name has been submitted by the form. Even if this parameter has been left empty, *array_key_exists() will return *true\* and the welcome message will be displayed.

**Warning: If you start making forms accessible from outside, calling my_name is dangerous if you don't protect it because JavaScript or HTML code can be inserted into it. We'll see later how to protect against this. You can also look at [other PHP documentation including those on security]({{< ref "/docs/coding/php/" >}}).**

## Useful server variables

Besides PHP_SELF, the $\_SERVER super-global array contains a number of useful elements that provide information about the web server and the request being processed, here are some:

{{< table "table-hover table-striped" >}}
| Element | Example | Description |
|---------|---------|-------------|
| QUERY_STRING | http://www.example.com/catalog/store.php?**category=kitchen&price=5** | The part of the URL after the ? and containing the URL parameters |
| PATH_INFO | http://www.example.com/catalog/store.php/**browse** | Additional path information, placed after the / at the end of the URL. This allows passing information to a script without using the query string |
| SERVER_NAME | http://**www.example.com**/catalog/store.php | FQDN of the site. If the server hosts multiple virtual hosts, it's the one currently being used that will be taken |
| DOCUMENT_ROOT | /var/www | Directory containing the web site documents on the server. If /var/www is the root of the web server, then http://www.example.com/catalog/store.php corresponds to /var/www/catalog/store.php |
| REMOTE_ADDR | 175.56.28.3 | IP of the client (source IP) |
| REMOTE_HOST | pool0560.cvx.dialup.verizon.net | If the web server is configured to do name resolution, this request is useful, however very few web servers use it (**too time-consuming processing**) |
| HTTP_REFERER | http://directory.google.com/Top/Shopping/Clothing/ | If the current URL was reached using a link, HTTP_REFER contains the URL of the page containing this link. **Warning: this value can be falsified!** |
| HTTP_USER_AGENT | Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1) | Information about the client's browser, as well as the OS used |
{{< /table >}}

## Accessing form parameters

The parameters of URLs and forms that use the GET method are placed in $\_GET, while the parameters of forms using the POST method are put in $\_POST.
The URL *http://www.example.com/catalog.php?product_id=21&category=fryingpan* will therefore place 2 values in $\_GET: $\_GET['product_id'] will be initialized with 21 and $\_GET['category'] with fryingpan.
The submission of the form below will place these same values in $\_POST, assuming that 21 has been entered in the input field and that Frying Pan has been chosen from the menu.

- Form containing 2 elements:

```php
<form method="POST" action="catalog.php">
<input type="text" name="product_id">
<select name="category">
<option value="ovenmitt">Pot Holder</option>
<option value="fryingpan">Frying Pan</option>
<option value="torch">Kitchen Torch</option>
</select>
<input type="submit" name="submit"> </form>
```

- Display of submitted form parameters:

```php
<form method="POST" action="catalog.php">
<input type="text" name="product_id">
<select name="category">
<option value="ovenmitt">Pot Holder</option>
<option value="fryingpan">Frying Pan</option>
<option value="torch">Kitchen Torch</option>
</select>
<input type="submit" name="submit">
</form>
Here are the submitted values:

product_id: <?php print $_POST['product_id']; ?>
<br/>
category: <?php print $_POST['category']; ?>
```

Which gives:
![Form display PHP](/images/form_affiche_php.avif)
The name of a form element that can contain multiple values must end with [] to tell the PHP interpreter that it must treat the various values as elements of an array. Thus, the values submitted for the <select> menu in the example below will be placed in $\_POST['lunch'].

- Multivalued form elements:

```php
<form method="POST" action="eat.php">
<select name="lunch[  ]" multiple>
<option value="pork">BBQ Pork Bun</option>
<option value="chicken">Chicken Bun</option>
<option value="lotus">Lotus Seed Bun</option>
<option value="bean">Bean Paste Bun</option>
<option value="nest">Bird-Nest Bun</option>
</select>
<input type="submit" name="submit">
</form>
```

- Accessing values of a multivalued element:

```php
<form method="POST" action="eat.php">
<select name="lunch[  ]" multiple>
<option value="pork">BBQ Pork Bun</option>
<option value="chicken">Chicken Bun</option>
<option value="lotus">Lotus Seed Bun</option>
<option value="bean">Bean Paste Bun</option>
<option value="nest">Bird-Nest Bun</option>
</select>
<input type="submit" name="submit">
</form>
Selected buns:
<br/>
<?php
foreach ($_POST['lunch'] as $choice) {
    print "You want a $choice bun. <br/>";
}
?>
```

Which gives:
![PHP multival](/images/php_multival.avif)

## Processing forms with functions

To know if a form has been submitted, you can use hidden parameters. If a hidden parameter is found in $\_POST, the form will be processed; otherwise, it will be displayed. This technique is used in the example below with the hidden parameter "\_submit_check":

```php
// Logic to do the right thing based on
// the hidden _submit_check parameter
if ($_POST['_submit_check']) {
    process_form( );
} else {
    show_form( );
}

// Do something when the form is submitted
function process_form( ) {
    print "Hello, " . $_POST['my_name'];
}

// Display the form
function show_form( ) {
    print<<<_HTML_
<form method="POST" action="$_SERVER[PHP_SELF]">
Your name: <input type="text" name="my_name">
<br/>
<input type="submit" value="Say Hello">
<input type="hidden" name="_submit_check" value="1">
</form>
_HTML_;
}
```

- Form data validation:

```php
// Logic to do the right thing based on
// the hidden _submit_check parameter
if ($_POST['_submit_check']) {
    if (validate_form( )) {
        process_form( );
    } else {
        show_form( );
    }
} else {
    show_form( );
}

// Do something when the form is submitted
function process_form( ) {
    print "Hello, " . $_POST['my_name'];
}

// Display the form
function show_form( ) {
    print<<<_HTML_
<form method="POST" action="$_SERVER[PHP_SELF]">
Your name: <input type="text" name="my_name">
<br/>
<input type="submit" value="Say Hello">
<input type="hidden" name="_submit_check" value="1">
</form>
_HTML_;
}

// Check the form data
function validate_form( ) {
    // Is my_name at least 3 characters long?
    if (strlen($_POST['my_name']) < 3) {
        return false;
    } else {
        return true;
   }
}
```

**You should ALWAYS validate information submitted by a form before processing it. NEVER throw it directly into execution. As code injection is possible, it is much safer to check.**

- Displaying form error messages:

```php
// Logic to do the right thing based on
// the hidden _submit_check parameter
if ($_POST['_submit_check']) {
    // If validate_form( ) returns errors, pass them to show_form( )
    if ($form_errors = validate_form( )) {
        show_form($form_errors);
    } else {
        process_form( );
    }
} else {
    show_form( );
}

// Do something when the form is submitted
function process_form( ) {
    print "Hello, " . $_POST['my_name'];
}

// Display the form
function show_form($errors = '') {
    // If some errors were passed in, print them out
    if ($errors) {
        print 'Please correct these errors: <ul><li>';
        print implode('</li><li>', $errors);
        print '</li></ul>';
    }

    print<<<_HTML_
<form method="POST" action="$_SERVER[PHP_SELF]">
Your name: <input type="text" name="my_name">
<br/>
<input type="submit" value="Say Hello">
<input type="hidden" name="_submit_check" value="1">
</form>
_HTML_;
}

// Check the form data
function validate_form( ) {
    // Start with an empty array of error messages
    $errors = array( );

    // Add an error message if the name is too short
    if (strlen($_POST['my_name']) < 3) {
        $errors[  ] = 'Your name must be at least 3 letters long.';
    }

    // Return the (possibly empty) array of error messages
    return $errors;
}
```

## Required elements

To check that a required element has been entered, check its length with the strlen() function:

```
if (strlen($_POST['email']) =  = 0) {
   $errors[  ] = "You must enter an email address."; }
}
```

## Numeric or string elements

To ensure that an entered value is an integer or float, use the conversion functions intval() or floatval(). They will return the number (integer or float) contained in a string, removing any superfluous text or other numerical formats.

To use these functions to validate forms, compare the submitted value with what you would get if you passed it to intval() or floatval(), then to strval(). This last function converts the cleaned number back to a string, allowing comparison with the concerned element of $\_POST. If the submitted string and the cleaned string don't match, there's something wrong with the submitted value: you should reject it. In the example below, we see how to check that a submitted element is indeed an integer.

- Checking if an element is an integer:

```php
if ($_POST['age'] != strval(intval($_POST['age']))) {
    $errors[  ] = 'Please enter a valid age.';
}
```

If $_POST['age'] is an integer like 59, 0 or -32, intval($POST['age']) will return, respectively 59, 0 or -32. The 2 values will match and nothing will be added to $errors. However, if $_POST['age'] is worth 52-dust, intval($\_POST['age']) will be worth 52. Since these 2 values are different, the if() test expression will be verified and a message will be added to $errors. If $_POST['age'] contains no digits, intval($POST['age']) will return 0. If, for example, old is the value submitted for $_POST['age'], intval($\_POST['age']) will return 0.

- Checking if an element is a real number:

```php
if ($_POST['price'] != strval(floatval($_POST['price']))) {
    $errors[  ] = 'Please enter a valid price.';
}
```

When validating elements - especially strings, it's often useful to remove leading and trailing spaces with the trim() function.
So, to prevent entering a string containing only spaces in a required element, you can combine the call to trim() with that of strlen(), as shown in the example below:

```php
if (strlen(trim($_POST['name'])) =  = 0) {
    $errors[  ] = "Your name is required.";
}
```

## Numeric ranges

- Testing if a value belongs to a numeric range:

```php
if ($_POST['age'] != strval(intval($_POST['age']))) {
    $errors[  ] = "Your age must be a number.";
} elseif (($_POST['age'] < 18) || ($_POST['age'] > 65)) {
    $errors[  ] = "Your age must be at least 18 and no more than 65.";
}
```

## Email addresses

If you don't want to embark on this verification operation by sending a message, you can still perform some syntax tests in the form validation code to eliminate malformed addresses.
Thus, the following regular expression:

```
^[^@\s]+@([-a-z0-9]+\.)+[a-z]{2,}$
```

matches the most common email addresses and will reject most bad ones. As shown in the example below, you can use it in a call to preg_match():

```php
if (! preg_match('/^[^@\s]+@([-a-z0-9]+\.)+[a-z]{2,}$/i',
                 $_POST['email'])) {
    $errors[  ] = 'Please enter a valid e-mail address';
}
```

## <select> menu

When using a <select> menu in a form, you need to ensure that the value submitted for it is one of the possible choices in the menu. Although with a classic browser like Firefox or IE (and yes unfortunately), a user cannot submit a value that is not part of the menu, a hacker can very well not use a browser and construct a request containing any value.

- Displaying a <select> menu

```php
$sweets = array('Sesame Seed Puff','Coconut Milk Gelatin Square',
                 'Brown Sugar Cake','Sweet Rice and Meat');

// Display the form
function show_form( ) {
    print<<<_HTML_
<form method="post" action="$_SERVER[PHP_SELF]">
Your Order: <select name="order">

_HTML_;
foreach ($GLOBALS['sweets'] as $choice) {
    print "<option>$choice</option>\n";
}
print<<<_HTML_
</select>
<br/>
<input type="submit" value="Order">
<input type="hidden" name="_submit_check" value="1">
</form>
_HTML_;
}
```

This will display:

```php
<form method="post" action="order.php">
Your Order: <select name="order">
<option>Sesame Seed Puff</option>
<option>Coconut Milk Gelatin Square</option>
<option>Brown Sugar Cake</option>
<option>Sweet Rice and Meat</option>
</select>
<br/>
<input type="submit" value="Order">
<input type="hidden" name="_submit_check" value="1">
</form>
```

For security concerns, it's important to check your code, here's how to do it:

```php
if (! array_key_exists($_POST['order'], $GLOBALS['sweets'])) {
    $errors[  ] = 'Please choose a valid order.';
}
```

## HTML and Javascript

The guestbook makes it easy for a malicious user to put HTML or Javascript code on the server that will then be executed by a browser without the user realizing it. This type of problem is called cross scripting attack because this poorly written guestbook allows code from one source (malicious user) to make it seem like it comes from elsewhere (the guestbook site).

**To prevent this type of attack, never display data coming from outside without having previously checked it.** Remove suspicious parts (HTML tags, for example) or encode special characters so that browsers cannot act on any HTML or Javascript code that might have been integrated. PHP offers you 2 functions to simplify these operations: strip_tags() removes HTML markers from a string and htmlentities() encodes special and HTML characters.

- Removing HTML tags from a string:

```
// Remove HTML from comments
$comments = strip_tags($_POST['comments']);
// Now it's OK to print $comments
print $comments;
```

- Encoding HTML entities in a string:

```
$comments = htmlentities($_POST['comments']);
// Now it's OK to print $comments
print $comments;
```

**In most applications, you should use htmlentities() to clean texts coming from outside**, as this function doesn't remove any content while still protecting you from cross-site scripting attacks.

## Displaying default values

When re-displaying a form due to an error, preserving the information already entered by the user can be practical.

- Building a default values array:

```php
if ($_POST['_submit_check']) {
    $defaults = $_POST;
} else {
    $defaults = array('delivery'  => 'yes',
                      'size'      => 'medium',
                      'main_dish' => array('taro','tripe'),
}
```

If $\_POST['_submit_check'] is initialized, it means that the form has been submitted: in this case, the default values should be those sent by the user. Otherwise, we can set our own default values.

- Setting a default value for an input field:

```php
print '<input type="text" name="my_name" value="' .
htmlentities($defaults['my_name']). '">';
```

- Setting a default value for a multi-line text area:

```php
print '<textarea name="comments">';
print htmlentities($defaults['comments']);
print '</textarea>';
```

- Setting a default value in a `<select>` menu:

```php
$sweets = array('puff' => 'Sesame Seed Puff',
                'square' => 'Coconut Milk Gelatin Square',
                'cake' => 'Brown Sugar Cake',
                'ricemeat' => 'Sweet Rice and Meat');

print '<select name="sweet">';
// $val is the option value, $choice is what's displayed
foreach ($sweets as $option => $label) {
    print '<option value="'.$option.'"';
    if ($option =  = $defaults['sweet']) {
        print ' selected="selected"';
    }
    print "> $label</option>\n";
}
print '</select>';
```

# Using databases for information storage

Most of this chapter uses PEAR DB, an abstraction layer for database access that plugs into PHP to simplify communications between your programs and database software.

Here's what we'll be using:

{{< table "table-hover table-striped" >}}
| ID | Name | Price | Is Spicy? |
|----|------|-------|-----------|
| 1 | Fried Bean Curd | 5.50 | 0 |
| 2 | Braised Sea Cucumber | 9.95 | 0 |
| 3 | Walnut Bun | 1.00 | 0 |
| 4 | Eggplant with Chilli Sauce | 6.50 | 1 |
{{< /table >}}

## Connecting to a database system

- Loading an external file with require:

```
require 'DB.php';
```

This line asks the PHP interpreter to execute all the code in DB.php, which is the main file of the PEAR DB package and defines the functions you'll use to communicate with your database.
_Include_ can be used instead of _require_. The difference is that include will not quit your program if the file does not exist, unlike require.

After loading the DB module, you need to establish a connection to the database using the DB::connect() function:

```php
require 'DB.php';
$db = DB::connect('mysql://penguin:top^hat@db.example.com/restaurant');
```

Explanation of the last line (DSN: Data Source Name):

```
db_program://user:password@hostname/database
```

- Here are the possible values for the DBMS:

{{< table "table-hover table-striped" >}}
| db_program | Database program |
|------------|------------------|
| dbase | dBase |
| fbsql | FrontBase |
| ibase | InterBase |
| ifx | Informix |
| msql | Mini SQL |
| mssql | Microsoft SQL Server |
| mysql | MySQL (version <= 4.0) |
| mysqli | MySQL (verison >= 4.1.2) |
| oci8 | Oracle (version 7,8 and 9) |
| odbc | ODBC |
| pgsql | PostgreSQL |
| sqlite | SQLite |
| sybase | Sybase |
{{< /table >}}

A call to DB:connect() returns an object that will allow you to interact with the database. In case of a connection problem, this function will return a different type of object, containing information about the reasons for the failure. Before going further, use the DB:isError() function.

- Connection error checking:

```php
require 'DB.php';
$db = DB::connect('mysql://penguin:top^hat@db.example.com/restaurant');
if (DB::isError($db)) { die("Can't connect: " . $db->getMessage( )); }
```

In the example above, the message is "Can't connect:" followed by "$db->getMessage( )" which returns more information about the error.

After calling DB:connect(), you can use the functions of the obtained object to interact with the database, so $db->getMessage( ) means "call the getMessage( ) function that is in the $db object". Here, $db contains the error information and the getMessage() function displays some of it. For example:

```
Can't connect: DB Error: connect failed
```

## Creating a table

- Creating a dishes table:

```php
CREATE TABLE dishes (
    dish_id INT,
    dish_name VARCHAR(255),
    price DECIMAL(4,2),
    is_spicy INT
)
```

Some column types include a length or format information in parentheses:

- VARCHAR(255): Variable length string with a maximum of 255 characters.
- DECIMAL(4,2): Decimal number with four digits, of which 2 are after the decimal point.

- Frequently used types for table columns:

{{< table "table-hover table-striped" >}}
| Column type | Description |
|-------------|-------------|
| VARCHAR(length) | Variable length string. length=maximum characters. |
| INT | integer |
| BLOB (PostgreSQL calls this type BYTEA) | Binary data string, up to 64 KB |
| DECIMAL(total*digits,decimal_places) | Decimal number with \_total_digits* digits, of which _decimal_places_ are after the decimal point (The decimal point is represented by a period, not a comma). |
| DATETIME (Oracle calls this type DATE) | Date and time: 2008-09-07 21:38:48 for example |
{{< /table >}}

- Sending a CREATE TABLE command to the database system:

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("connection error: " . $db->getMessage( )); }
$q = $db->query("CREATE TABLE dishes (
        dish_id INT,
        dish_name VARCHAR(255),
        price DECIMAL(4,2),
        is_spicy INT
)");
```

- Deleting a table:

```
DROP TABLE dishes
```

## Storing information in a database

### Insert

To store data in the database, simply pass an INSERT instruction to the object's query() function:

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("connection error: " . $db->getMessage( )); }
$q = $db->query("INSERT INTO dishes (dish_name, price, is_spicy)     VALUES ('Sesame Seed Puff', 2.50, 0)");
```

For more explanations on the INSERT function of SQL, [refer to the SQL introduction]({{< ref "docs/Coding/SQL/introduction_to_sql.md#create---data-insertion">}}).

Instead of calling DB:isError after each query, it's more convenient to use the _setErrorHandling()_ function to set a default behavior for error handling: just pass the PEAR*ERROR_DIE constant to \_setErrorHandling()* for your program to automatically display an error message and stop in case a query fails.

- Automatic error handling with setErrorHandling():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("Can't connect: " . $db->getMessage( )); }
// print a message and quit on future database errors
$db->setErrorHandling(PEAR_ERROR_DIE);
$q = $db->query("INSERT INTO dishes (dish_size, dish_name, price, is_spicy)
    VALUES ('large', 'Sesame Seed Puff', 2.50, 0)");
print "Query Succeeded!";
}
```

This will display:

```
DB Error: syntax error
```

Since the program ends as soon as it encounters the error, the last line of the example will never execute: the "Query Success!" message will not be displayed.

**The setErrorHandling() function belongs to the $db object: you must first call DB::connect() to get this object.**

### Update

For more explanations on the UPDATE function of SQL, [refer to the SQL introduction]({{< ref "docs/Coding/SQL/introduction_to_sql.md#update---data-update">}}).

Here's an example of modifying data with query():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("connection error: " . $db->getMessage( )); }
// Decrease the price some some dishes
$db->query("UPDATE dishes SET price=price - 5 WHERE price > 20");
print 'Changed the price of ' . $db->affectedRows( ) . 'rows.';
```

### Delete

For more explanations on the DELETE function of SQL, [refer to the SQL introduction]({{< ref "docs/Coding/SQL/introduction_to_sql.md#delete---delete-data">}}).

- Deleting data with query():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("connection error: " . $db->getMessage( )); }
// remove expensive dishes
if ($make_things_cheaper) {
    $db->query("DELETE FROM dishes WHERE price > 19.95");
} else {
    // or, remove all dishes
    $db->query("DELETE FROM dishes");
}
```

The affectedRows() function indicates the number of rows that have been modified or deleted by an UPDATE or DELETE instruction. To know the number of rows affected by these queries, call it immediately after calling query().

- Knowing the number of rows affected by an UPDATE or DELETE:

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die("connection error: " . $db->getMessage( )); }
// Decrease the price some some dishes
$db->query("UPDATE dishes SET price=price - 5 WHERE price > 20");
print 'Changed the price of ' . $db->affectedRows( ) . 'rows.';
```

## Secure insertion of data from forms

Displaying unverified data from forms can expose you and your users to a cross-site scripting attack. In SQL queries, this data can pose a similar problem, called "**SQL Injection Attack**".

By submitting a cleverly constructed value through a form, a malicious user can inject arbitrary SQL queries into your database system. To prevent this type of attack, you must protect special characters (**especially the apostrophe**) in SQL queries. PEAR DB provides a very useful feature called _placeholders_, which makes all this very easy.

**The 'magic quotes' functionality is activated by the "magic_quotes_gpc" directive in the PHP configuration**. For more efficient and readable management of form parameters, disable magic_quotes_gpc and use placeholders or a protection function instead when you need to prepare external data for use in a database query.

To use a placeholder in a query, put a "?" in it wherever you want to put a value.

- Example of inserting unsafe data from a form:

```
$db->query("INSERT INTO dishes (dish_name) VALUES ('$_POST[new_dish_name]')");
```

- Example of safe insertion of data from a form:

```
$db->query('INSERT INTO dishes (dish_name) VALUES (?)', array($_POST['new_dish_name']));
```

You don't need to surround placeholders with apostrophes in the query. DB takes care of this for you as well.

- Using multiple placeholders:

```
$db->query('INSERT INTO dishes (dish_name,price,is_spicy) VALUES (?,?,?)', array($_POST['new_dish_name'], $_POST['new_price'], $_POST['is_spicy']));
```

## Generating unique identifiers

PEAR DB sequences will help you produce unique integer identifiers. Indeed, when requesting the next identifier of a particular sequence, you get a number that you can be sure will be unique in that sequence. Even if 2 PHP scripts are running simultaneously and request the next identifier in a sequence at the same time, each of them will get a different value.

- Getting a sequence identifier:

```php
$dish_id = $db->nextID('dishes');
$db->query("INSERT INTO orders (dish_id, dish_name, price, is_spicy)
    VALUES ($dish_id, 'Fried Bean Curd', 1.50, 0)");
```

## Retrieving data stored in the database

When the _query()_ function successfully executes a SELECT command, query() returns an object giving access to the obtained rows and, each time you call the _fetchRow()_ function on this object, you will get the next row of the query result. When there are no more rows, _fetchRow()_ returns a value considered false, which allows it to be used in a while() loop.

- Retrieving rows with query() and fetchRow():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
$q = $db->query('SELECT dish_name, price FROM dishes');
while ($row = $q->fetchRow( )) {
    print "$row[0], $row[1] \n";
}
```

When there are no more rows to return, fetchRow() will return a value that will be evaluated as false and the while() loop will terminate.
To know the number of rows returned by a SELECT query (without traversing them all), use the numrows() function of the object returned by query().

- Counting rows with numrows():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
$q = $db->query('SELECT dish_name, price FROM dishes');
print 'There are ' . $q->numrows( ) . ' rows in the dishes table.';
```

The getAll() function executes a SELECT query and returns an array containing all the rows of the result.

- Retrieving rows with getAll():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
$rows = $db->getAll('SELECT dish_name, price FROM dishes');
foreach ($rows as $row) {
    print "$row[0], $row[1] \n";
}
```

This will display:

```php
Walnut Bun, 1.00
Cashew Nuts and White Mushrooms, 4.95
Dried Mulberries, 3.00
Eggplant with Chili Sauce, 6.50
```

### Select

For more explanations on the SELECT function of SQL, [refer to the SQL introduction]({{< ref "docs/Coding/SQL/introduction_to_sql.md##select---retrieve-data">}}).

If you only expect one row as a result of a query, use getRow() instead: this function executes SELECT and returns only one row.

- Retrieving a row with getRow():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
$cheapest_dish_info = $db->getRow('SELECT dish_name, price
                                   FROM dishes ORDER BY price LIMIT 1');
print "$cheapest_dish_info[0], $cheapest_dish_info[1]";
```

If you're only interested in a single column of a single row, use the getOne() function. It sends a SELECT query and returns a simple value: the first column of the first row returned.

- Retrieving a value with getOne():

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
$cheapest_dish = $db->getOne('SELECT dish_name, price FROM dishes ORDER print "The cheapest dish is $cheapest_dish";
```

## Modifying the format of result rows

For more explanations on the Order By and Limit functions of SQL, [refer to the SQL introduction]({{< ref "docs/Coding/SQL/introduction_to_sql.md#order-by-and-limit---data-sorting">}}).

Until now, _fetchRow(), getAll() and getOne()_ have returned rows from the database in the form of arrays indexed by integers. This makes it easy and quick to interpolate values into strings between double apostrophes, but knowing, for example, which column of the SELECT query corresponds to element 6 of the array is a tricky and error-prone operation. To avoid this, PEAR DB allows retrieving each row of the result in an array indexed by strings or as an object.

The fetch mode, set by the setFetchMode() function, controls the format of the result rows.

By passing DB_FETCHMODE_ASSOC as a parameter, you get the result rows as string-indexed arrays.

- Retrieving rows as string-indexed arrays:

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');

// Change the fetch mode to string-keyed arrays
$db->setFetchMode(DB_FETCHMODE_ASSOC);

print "With query( ) and fetchRow( ): \n";
// get each row with query( ) and fetchRow( );
$q = $db->query("SELECT dish_name, price FROM dishes");
while($row = $q->fetchRow( )) {
    print "The price of $row[dish_name] is $row[price] \n";
}

print "With getAll( ): \n";
// get all the rows with getAll( );
$dishes = $db->getAll('SELECT dish_name, price FROM dishes');
foreach ($dishes as $dish) {
    print "The price of $dish[dish_name] is $dish[price] \n";
}

print "With getRow( ): \n";
$cheap = $db->getRow('SELECT dish_name, price FROM dishes
    ORDER BY price LIMIT 1');
print "The cheapest dish is $cheap[dish_name] with price $cheap[price]";
```

This will display:

```php
With query( ) and fetchRow( ):
The price of Walnut Bun is 1.00
The price of Cashew Nuts and White Mushrooms is 4.95
The price of Dried Mulberries is 3.00
The price of Eggplant with Chili Sauce is 6.50
With getAll( ):
The price of Walnut Bun is 1.00
The price of Cashew Nuts and White Mushrooms is 4.95
The price of Dried Mulberries is 3.00
The price of Eggplant with Chili Sauce is 6.50
With getRow( ):
The cheapest dish is Walnut Bun with price 1.00
```

The rows returned by these functions are now indexed by strings representing the names of the columns in the dishes table.

**To get rows as objects, pass the constant DB_FETCHMODE_OBJECT as a parameter to setFetchMode(). Each result row will then be an object whose attributes will have the names of the table columns (like the keys of the array when using the DB_FETCHMODE_ASSOC fetch mode).**
Compared to string-indexed arrays, the DB_FETCHMODE_OBJECT fetch mode allows using a more concise and easier to interpolate syntax to designate the data: the name of the object is given followed by -> then the name of the desired information. $dish->dish_name, for example, designates the dish_name information found in the $dish object.

- Retrieving rows as objects:

```php
require 'DB.php';
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');

// Change the fetch mode to objects
$db->setFetchMode(DB_FETCHMODE_OBJECT);

print "With query( ) and fetchRow( ): \n";
// get each row with query( ) and fetchRow( );
$q = $db->query("SELECT dish_name, price FROM dishes");
while($row = $q->fetchRow( )) {
    print "The price of $row->dish_name is $row->price \n";
}

print "With getAll( ): \n";
// get all the rows with getAll( );
$dishes = $db->getAll('SELECT dish_name, price FROM dishes');
foreach ($dishes as $dish) {
    print "The price of $dish->dish_name is $dish->price \n";
}

print "With getRow( ): \n";
$cheap = $db->getRow('SELECT dish_name, price FROM dishes
    ORDER BY price LIMIT 1');
print "The cheapest dish is $cheap->dish_name with price $cheap->price";
```

## Secure retrieval of form data

**When using form data or any other external source in the WHERE clause of a SELECT, UPDATE or DELETE command, you need to ensure that SQL wildcard characters are properly disabled in this data.**

To ensure that SQL wildcards entered in a form are disabled in queries, you must give up the comfort and ease of placeholders and use 2 other functions:

- quoteSmart() and DB
- strtr(), provided by PHP

First apply quoteSmart() to the submitted value, which will have the same effect on apostrophes as a placeholder: this call will replace for example, 'Homard à l'américaine' with 'Homard à l\'américaine'. The next step is to use strtr() to disable the SQL wildcard characters, % and \_.
The value with protected apostrophes and disabled wildcard characters can then be safely used in a query.

- SELECT command without placeholders:

```php
// First, do normal quoting of the value
$dish = $db->quoteSmart($_POST['dish_search']);
// Then, put backslashes before underscores and percent signs
$dish = strtr($dish, array('_' => '\_', '%' => '\%'));
// Now, $dish is sanitized and can be interpolated right into the query
$matches = $db->getAll("SELECT dish_name, price FROM dishes WHERE dish_name LIKE $dish");
```

Here, you can't use a placeholder, because disabling the wildcard characters has to come after protecting the apostrophes. Indeed, this places backslashes before simple apostrophes, but also before backslashes: if strtr() or with placeholders) it would transform it back into '\\%homard\\%' This value will be interpreted by the database as a literal backslash, followed by a wildcard %. If quoteSmart() is called first, however, %homard% will be transformed into '%homard%', then strtr() will transform it into '\%homard\%'.

This value will be interpreted by the database as a literal percent sign, followed by homard, followed by another literal percent sign, which is what the user had entered.

Undeactivated wildcard characters can have an even more terrible effect in the WHERE clause of an UPDATE or DELETE command. **The example below is an example of what should absolutely not be done**. It shows a query incorrectly using placeholders to allow an external value to control which dishes will be sold for 1 euro.

- **Incorrect** use of placeholders in an UPDATE command

```php
$db->query('UPDATE dishes SET price = 1 WHERE dish_name LIKE ?',
            array($_POST['dish_name']));
```

## Complete consultation form

```php
<?php

// Load PEAR DB
require 'DB.php';
// Load the form helper functions.
require 'formhelpers.php';

// Connect to the database
$db = DB::connect('mysql://hunter:w)mp3s@db.example.com/restaurant');
if (DB::isError($db)) { die ("Can't connect: " . $db->getMessage( )); }

// Set up automatic error handling
$db->setErrorHandling(PEAR_ERROR_DIE);

// Set up fetch mode: rows as objects
$db->setFetchMode(DB_FETCHMODE_OBJECT);

// Choices for the "spicy" menu in the form
$spicy_choices = array('no','yes','either');

// The main page logic:
// - If the form is submitted, validate and then process or redisplay
// - If it's not submitted, display
if ($_POST['_submit_check']) {
    // If validate_form( ) returns errors, pass them to show_form( )
    if ($form_errors = validate_form( )) {
        show_form($form_errors);
    } else {
        // The submitted data is valid, so process it
        process_form( );
    }
} else {
    // The form wasn't submitted, so display
    show_form( );
}

function show_form($errors = '') {
    // If the form is submitted, get defaults from submitted parameters
    if ($_POST['_submit_check']) {
        $defaults = $_POST;
    } else {
        // Otherwise, set our own defaults
        $defaults = array('min_price' => '5.00',
                      'max_price' => '25.00');
}

    // If errors were passed in, put them in $error_text (with HTML markup)
    if ($errors) {
        $error_text = '<tr><td>You need to correct the following errors:';
        $error_text .= '</td><td><ul><li>';
        $error_text .= implode('</li><li>',$errors);
        $error_text .= '</li></ul></td></tr>';
    } else {
        // No errors? Then $error_text is blank
        $error_text = '';
    }

    // Jump out of PHP mode to make displaying all the HTML tags easier
?>
<form method="POST" action="<?php print $_SERVER['PHP_SELF']; ?>">
<table>
<?php print $error_text ?>

<tr><td>Dish Name:</td>
<td><?php input_text('dish_name', $defaults) ?></td></tr>

<tr><td>Minimum Price:</td>
<td><?php input_text('min_price', $defaults) ?></td></tr>

<tr><td>Maximum Price:</td>
<td><?php input_text('max_price', $defaults) ?></td></tr>

<tr><td>Spicy:</td>
<td><?php input_select('is_spicy', $defaults, $GLOBALS['spicy_choices']); ?>
</td></tr>

<tr><td colspan="2" align="center"><?php input_submit('search','Search'); ?>
</td></tr>

</table>
<input type="hidden" name="_submit_check" value="1"/>
</form>
<?php
      } // The end of show_form( )

function validate_form( ) {
    $errors = array( );

    // minimum price must be a valid floating point number
    if ($_POST['min_price'] != strval(floatval($_POST['min_price']))) {
        $errors[  ] = 'Please enter a valid minimum price.';
    }

    // maximum price must be a valid floating point number
    if ($_POST['max_price'] != strval(floatval($_POST['max_price']))) {
        $errors[  ] = 'Please enter a valid maximum price.';
    }

    // minimum price must be less than the maximum price
    if ($_POST['min_price'] >= $_POST['max_price']) {
        $errors[  ] = 'The minimum price must be less than the maximum price.';
    }

    if (! array_key_exists($_POST['is_spicy'], $GLOBALS['spicy_choices'])) {
        $errors[  ] = 'Please choose a valid "spicy" option.';
    }
    return $errors;
}

function process_form( ) {
    // Access the global variable $db inside this function
    global $db;

    // build up the query
    $sql = 'SELECT dish_name, price, is_spicy FROM dishes WHERE
            price >= ? AND price <= ?';

    // if a dish name was submitted, add to the WHERE clause
    // we use quoteSmart( ) and strtr( ) to prevent user-entered wildcards from working
    if (strlen(trim($_POST['dish_name']))) {
        $dish = $db->quoteSmart($_POST['dish_name']);
        $dish = strtr($dish, array('_' => '\_', '%' => '\%'));
        $sql .= " AND dish_name LIKE $dish";
    }

    // if is_spicy is "yes" or "no", add appropriate SQL
    // (if it's "either", we don't need to add is_spicy to the WHERE clause)
    $spicy_choice = $GLOBALS['spicy_choices'][ $_POST['is_spicy'] ];
    if ($spicy_choice =  = 'yes') {
        $sql .= ' AND is_spicy = 1';
    } elseif ($spicy_choice =  = 'no') {
        $sql .= ' AND is_spicy = 0';
    }

    // Send the query to the database program and get all the rows back
    $dishes = $db->getAll($sql, array($_POST['min_price'],
                                      $_POST['max_price']));

    if (count($dishes) =  = 0) {
        print 'No dishes matched.';
    } else {
        print '<table>';
        print '<tr><th>Dish Name</th><th>Price</th><th>Spicy?</th></tr>';
        foreach ($dishes as $dish) {
            if ($dish->is_spicy =  = 1) {
                $spicy = 'Yes';
            } else {
                $spicy = 'No';
            }
            printf('<tr><td>%s</td><td>$%.02f</td><td>%s</td></tr>',
                   htmlentities($dish->dish_name), $dish->price, $spicy);
        }
    }
}
?>
```

The example above contains an additional line in the database configuration code: a call to setFetchMode().

### MySQL without PEAR DB

PEAR DB rounds a lot of angles for database access from a PHP program, but there are 2 cases where it's not always the best choice: it may not be available on some systems and a program using PHP's predefined functions, designed specifically for a particular database, is faster than if it used PEAR DB.

The differences are in the details: the available functions and how they work vary depending on the database and, in general, you have to retrieve the results row by row because you don't have the convenience provided by getAll(). There is also no unified error handling.

{{< table "table-hover table-striped" >}}
| PEAR DB function | mysqli function | Comments |
|------------------|----------------|----------|
| $db = DB::connect( DSN) | $db = mysqli_connect(hostname, username, password, database) | |
| $q = $db->query(SQL) | $q = mysqli_query($db,SQL) | There is no placeholder support in mysqli_query( ). |
| $row = $q->fetchRow( ) | $row = mysqli_fetch_row($q) | mysqli_fetch_row( ) always returns numerically indexed arrays. Use mysqli_fetch_assoc( ) for string-indexed arrays or mysqli_fetch_object( ) for objects. |
| $db->affectedRows( ) | mysqli_affected_rows($db) | |
| $q->numRows( ) | mysqli_num_rows($q) | |
| $db->setErrorHandling(ERROR_MODE) | None | You can't set automatic error handling with mysqli, but mysqli_connect_error( ) gives you the error message if something goes wrong connecting to the database program, and mysqli_error($db)gives you the error message after a query or other function call fails. |
{{< /table >}}

- Function process_form() using mysqli:

```php
function process_form( ) {
    // Access the global variable $db inside this function
    global $db;

    // build up the query
    $sql = 'SELECT dish_name, price, is_spicy FROM dishes WHERE ';

    // add the minimum price to the query
    $sql .= "price >= '" .
            mysqli_real_escape_string($db, $_POST['min_price']) . "' ";

    // add the maximum price to the query
    $sql .= " AND price <= '" .
            mysqli_real_escape_string($db, $_POST['max_price']) . "' ";

    // if a dish name was submitted, add to the WHERE clause
    // we use mysqli_real_escape_string( ) and strtr( ) to prevent
    // user-entered wildcards from working
    if (strlen(trim($_POST['dish_name']))) {
        $dish = mysqli_real_escape_string($db, $_POST['dish_name']);
        $dish = strtr($dish, array('_' => '\_', '%' => '\%'));
        // mysqli_real_escape_string( ) doesn't add the single quotes
        // around the value so you have to put those around $dish in
        // the query:
        $sql .= " AND dish_name LIKE '$dish'";
    }

    // if is_spicy is "yes" or "no", add appropriate SQL
    // (if it's either, we don't need to add is_spicy to the WHERE clause)
    $spicy_choice = $GLOBALS['spicy_choices'][ $_POST['is_spicy'] ];
    if ($spicy_choice =  = 'yes') {
        $sql .= ' AND is_spicy = 1';
    } elseif ($spicy_choice =  = 'no') {
        $sql .= ' AND is_spicy = 0';
    }

    // Send the query to the database program and get all the rows back
    $q = mysqli_query($db, $sql);

    if (mysqli_num_rows($q) =  = 0) {
        print 'No dishes matched.';
    } else {
        print '<table>';
        print '<tr><th>Dish Name</th><th>Price</th><th>Spicy?</th></tr>';
        while ($dish = mysqli_fetch_object($q)) {
            if ($dish->is_spicy
                $spicy = 'Yes';
            } else {
                $spicy = 'No';
            }
            printf('<tr><td>%s</td><td>$%.02f</td><td>%s</td></tr>',
                   htmlentities($dish->dish_name), $dish->price, $spicy);
        }
    }
}
```

Since PEAR DB's placeholders are not available here, the minimum and maximum prices are directly placed in the variable $sql that contains the query, but after being protected by mysqli_real_escape_string(). This same function is also used to protect $\_POST['dish_name']. Finally, the functions used to pass the query to the database and retrieve the results are different: the mysqli_query() function sends the query, mysqli_num_rows() returns the number of rows in the result and mysqli_fetch_object() retrieves each row of this result in an object.

## Remembering users with cookies and sessions

A single cookie is perfect for memorizing the same piece of information but, often, you'll need to store several pieces of data for the same user (the contents of their virtual shopping cart, for example). Using multiple cookies for this purpose involves heavy processing: PHP sessions can solve this problem.

A session uses a cookie to differentiate users and facilitates the memorization of a temporary set of data for each of them on the server. These data persist between requests: during a request, you can add a variable to a user's session (like when adding an item to the virtual cart) and, at the next request, you can retrieve what's in this session (like on an order validation page, when you need to list everything in the cart). In the Information Storage and Retrieval section, we'll see how to use them.

## Resources
- [Introduction to PHP 5 Book](https://books.google.fr/books?id=fAlN55KwKjcC)
- [Pear and PHP libraries](/pdf/pear_et_les_librairies_php.pdf)
