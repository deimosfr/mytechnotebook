---
weight: 999
url: "/Introduction_à_Perl/"
title: "Introduction to Perl"
description: "A comprehensive guide to Perl programming language covering basics, syntax, functions, modules, and advanced techniques"
categories: ["Database", "Linux"]
date: "2013-09-24T11:20:00+02:00"
lastmod: "2013-09-24T11:20:00+02:00"
tags: ["Perl", "Development"]
toc: true
---

## Introduction

This documentation is intended for people with good knowledge of software development, as it simply provides notes on Perl syntax and how it works.

**Perl** is a programming language created by Larry Wall in 1987, incorporating features from the C language and scripting languages like sed, awk, and shell.

Larry Wall gives two interpretations of the "PERL" acronym:

- _Practical Extraction and Report Language_
- _Pathetically Eclectic Rubbish Lister_

These names are retroactive acronyms.

The organization responsible for the development and promotion of Perl is The Perl Foundation. In France, "Les Mongueurs de Perl" promotes this language, notably through Perl Days events.

## Basics

Exponents:

```perl
7.5e-24 # gives 7.5 times 10 to the power of -24
```

Simplified integer literal notation:

```perl
123456789 # Can also be written as 123_45678_9
```

Here are the escape character sequences:

| Construct | Meaning                                                  |
| --------- | -------------------------------------------------------- |
| \n        | Newline                                                  |
| \r        | Return                                                   |
| \t        | Tab                                                      |
| \f        | FormFeed                                                 |
| \b        | Backspace                                                |
| \a        | Bell                                                     |
| \e        | Escape (ASCII escape character)                          |
| \007      | Any octal ASCII value (here, 007 = bell)                 |
| \x7F      | Any hex ASCII value (here, 7F = delete)                  |
| \cC       | A "control" character (here, Ctrl-C)                     |
| \\        | Backslash                                                |
| \"        | Double quote                                             |
| \l        | Lowercase next letter                                    |
| \L        | Lowercase all following letters until \E                 |
| \u        | Uppercase next letter                                    |
| \U        | Uppercase all following letters until \E                 |
| \Q        | Quote non-word characters by adding a backslash until \E |
| \E        | Terminate \L, \U, or \Q                                  |

String operators (concatenation):

```perl
"hello" . "world" # equivalent to "helloworld"
"hello" . ' ' . "world" # equivalent to 'hello world'
'hello world' . "\n" # equivalent to "helloworld\n"
"fred" x 3 # equivalent to "fredfredfred"
"5 * 3" # equivalent to 15
```

Be careful with string operators for numbers:

```perl
"5" x 4 # equivalent to "5555"
```

Binary assignment operator:

```perl
$fred = $fred + 5 # Can be written as $fred += 5
```

Avoid:

```perl
print "$fred" -> print $fred
```

Operator precedence and associativity:

| Associativity | Operators                                           |
| ------------- | --------------------------------------------------- |
| left          | parentheses and arguments to list operators         |
| left          | ->                                                  |
|               | ++ -- (autoincrement and autodecrement)             |
| right         | \*\*                                                |
| right         | \ ! ~ + - (unary operators)                         |
| left          | =~ !~                                               |
| left          | \* / % x                                            |
| left          | + - . (binary operators)                            |
| left          | << >>                                               |
|               | named unary operators (~X filetests, rand)          |
|               | < <= > >= lt le gt ge (the "unequal" ones)          |
|               | == != <=> eq ne cmp (the "equal" ones)              |
| left          | &                                                   |
| left          | \| ^                                                |
| left          | &&                                                  |
| left          | \|\|                                                |
|               | .. ...                                              |
| right         | ? : (ternary)                                       |
| right         | = += -= \*= etc. (and similar assignment operators) |
| left          | , =>                                                |
|               | list operators (rightward)                          |
| right         | not                                                 |
| left          | and                                                 |
| left          | or xor                                              |

Comparison operators:

| Comparison               | Numeric | String |
| ------------------------ | ------- | ------ |
| Equal                    | ==      | eq     |
| Not equal                | !=      | ne     |
| Less than                | <       | lt     |
| Greater than             | >       | gt     |
| Less than or equal to    | <=      | le     |
| Greater than or equal to | >=      | ge     |

## Loops

### If

Here is the structure of an if statement:

```perl
if ($mytruc gt 'toto') {
   print "'$mytruc' is greater than toto.\n";
} else {
   print "'$mytruc' is less than toto.\n";
}
```

To get STDIN as a scalar variable (Reminder: STDIN is used to capture what the user types on the keyboard):

```perl
$line = <STDIN>;
if ($line eq "\n") {
  print "This was just an empty line!\n";
} else {
  print "This input line was: $line";
}
```

### while

Here is the structure of a while loop:

```perl
$total = 0;
while ($total < 10) {
  $total += 1;
  print "total is now $total\n"; # gives values from 1 to 10
}
```

## Useful Functions

### Defined

Checks if a value is defined or not

```perl
$bre1 = <STDIN>;
if (defined($bre1) ) {
  print "the input is $bre1";
} else {
  print "no input available!\n";
}
```

### Chomp

Removes newline characters:

```perl
$text = "a line of text\n";
chomp($text);
```

## Arrays

Here's how to set the 0th entry of an array:

```perl
$tab[0] = "test";
```

And here's how you can calculate an array index:

```perl
$tab[5 - 0.4] = "test"; # This is $tab[4]
```

To create 99 undefined elements:

```perl
$tab[99] = "test";
```

To define the last element of the array, there are two methods:

```perl
$tab[$#tab] = 'test';
$tab[-1] = 'test';
```

To know the number of columns in an array:

```perl
my $numbers = @tab;
```

Now $numbers contains the number of columns.

To empty an array:

```perl
@tab = ();
```

### Lists

Here's an example of lists:

```perl
(1..5, 7) # from 1 to 5 and 7;
```

#### qw

If for instance we want to have a list like this:

```perl
("toto", "tata", "titi")
```

To simplify things, we can do this:

```perl
qw/ toto tata titi /
```

You can replace **/** with !, #, (), {}, [] or <>.

The advantage of using these notations is that malformed spaces or new lines will be automatically ignored, and your words will be taken into account. Here's another possible form of writing:

```perl
qw{
  /usr/dict/words;
  /home/root/.ispell_french;
}
```

We can also do:

```perl
@nombre = qw(un deux trois);
```

#### List Assignments

To assign lists:

```perl
($toto, $tata, $titi) = ("un", "deux", "trois");
```

Here toto equals un, tata equals deux...

We can swap this way:

```perl
($toto, $tata) = ($tata, $toto);
```

Be careful with ignored elements like:

```perl
($toto, $tata) = qw<un deux trois quatre>;
```

Here, only the first 2 will be considered. For the reverse:

```perl
($toto, $tata) = qw<un>;
```

$tata will have an undef value.

### pop

Pop is used to remove elements from the end of an array:

```perl
@tab = 5..9 # my array goes from 5 to 9
$toto = pop(@tab); # $toto will receive 9, and the array will only have 5,6,7 and 8
$tata = pop @tab # same thing, so $tata gets 8 and the array stops at 7
pop @tab; # 7 is thrown away
```

When the array is empty, pop returns undef.

### push

Push is used to add values to the end of arrays:

```perl
push (@tab, 0); # the array now contains 6, 7 and 0
push @tab, 1..10 # @tab now contains 10 more items
```

### shift

Shift is like pop but acts at the beginning of the array:

```perl
@tab = qw<tata toto titi>;
$myshift = shift @tab; # @tab now contains toto and titi.
```

### unshift

Unshift is like push, but acts at the beginning of the array:

```perl
unshift @tab, tata # Now my array is back to tata, toto, titi
```

### foreach

A foreach will allow you to iterate through a complete list:

```perl
@total = qw/toto titi tata/;
foreach $toutou (@total) {
  print "$toutou\n";
}
```

### $\_

This is probably the most used default variable in Perl. It allows you to not declare variables in a loop. For example:

```perl
foreach (1..10) { # Default use of $_
  print "I can count to $_!\n";
}
```

Another example:

```perl
$_ = "My test\n";
print; # prints the default variable $_
```

### reverse

Reverse takes a list of values and returns the list in reverse order:

```perl
@num = 6..10;
@tab1 = reverse(@num); # returns 10,9,8,7,6.
@tab2 = reverse 6..10; # Same but without fetching from the array.
@num = reverse(@num); # Directly replaces in the original array.
@inverse = reverse qw/ tata toto titi/; # gives tata titi toto
$inverse = reverse qw/ tata toto titi/; # gives ototititatat
```

### sort

Like the sort binary on Unix, it sorts alphabetically (**but in ASCII order**):

```perl
@planetes = qw/ terre saturne mars pluton /;
@triee = sort(@planetes); # gives mars, pluton, saturne and terre.
@triee = sort @planetes; # Same
@inverse = reverse sort(@planetes); # reverses the result above.
```

### List and Scalar Context

It's important to understand this section:

```perl
5 + x # x must be a scalar
sort x # x must be a list
```

```perl
@personnes = qw( toto tata titi);
@triee = sort @personnes;
$nombre = 5 + @personnes # gives 5 + 3
$n = @personnes; # gives 3
```

To force a scalar context:

```perl
print "I ", scalar @tab, " force the scalar\n";
```

## Functions

Define with sub:

```perl
sub way {
   $n += 1;
   print "The subway number is: $n\n";
}
```

And call the function with an ampersand (&):

```perl
&way; # Will display: The subway number is: 1
&way; # Will display: The subway number is: 2
```

Here's an example:

```perl
my @c;
sub liste_de_toto_a_titi {
   if ($toto < $titi) {
      # Count up from toto to titi
      push @c, $toto..$titi;
   } else {
      # Count down from toto to titi
      push @c, $titi..$toto;
   }
}
$toto = 11;
$titi = 6;
@c = &liste_de_toto_a_titi; # @c receives (11,10,9,8,7,6)
print @c;
```

### Function Arguments

Let's pass 2 arguments to a function:

```perl
$n = &max(10,15);
```

And in the function, we'll call the **2 arguments with $_[1] or $_[2]** (but only in the function!) which are **part of the @\_ list**. **These variables have nothing to do with the $\_ variable**. Here's an example:

```perl
sub max {
   if ($_[0] < $_[1]) {
      $_[0];
   } else {
      $_[1];
   }
}
```

If a called value is not defined, it will be set to undef.

To ensure receiving the exact number of arguments:

```perl
sub max {
   if (@_ != 2) {
      print "Warning! &max should receive 2 arguments\n";
   }
}
```

### my

The my value allows you to create private variables. Declared in a function, it will be called in the function and at the end, it will be removed from memory.

```perl
sub max {
   my ($a, $b); # new private variables for this block
   ($a, $b) = @_; # give names to parameters
   if ($a > $ b) { $a } else { $b }
}
```

We can simplify things:

```perl
my($a, $b) = @_;
```

### local

Local is the old name for my, except that local saves a copy of the variable's value in a secret place (stack). This value cannot be consulted, modified, or deleted while it's saved. Then local initializes the variable to an empty value (undef for scalars, or empty list for arrays), or the assigned value. When returning from the function, the variable is automatically restored to its original value. In short, the variable was borrowed for a time and returned before anyone noticed.

You cannot replace local with my in old Perl scripts, but remember to use my preferably when creating new scripts. Example:

```perl
local($a, $b) = @_;
```

### return

Return will return the value of a function. For example, return, placed in a foreach will return a value when it has been found.

## use strict

To write "clean" code, it's better to put this in your scripts:

```perl
use strict;
```

## Hashes

A hash table (or hash) is like an array except that instead of having numbers as references, we will have keys.

![Perl hash table](/images/perl_hash_tab.avif)

In case it's not clear enough, here's the difference:

- Array:

| Identifier or key (**not modifiable**) | 0    | 1    | 2    |
| -------------------------------------- | ---- | ---- | ---- |
| Value (**modifiable**)                 | tata | titi | toto |

- Hash table:

| Identifier or key (**modifiable**) | IP          | Machine       | Name |
| ---------------------------------- | ----------- | ------------- | ---- |
| Value (**modifiable**)             | 192.168.0.1 | toto-portable | Toto |

We declare a hash table this way:

```perl
$my_hash_table{"identifier"} = "value";
```

To copy a column:

```perl
$my_hash_table{"identifier1"} .= $my_hash_table{"identifier2"};
```

Perl decides the layout of keys in the hash table. You won't have the possibility to organize them as you wish! This allows Perl to access information you're looking for more quickly.

### List

To declare a hash list, we can do:

```perl
%list_hash = ("IP", "192.168.0.1", "Machine", "toto-portable", "Name", "Toto");
```

There's another way, much clearer, to make list declarations:

```perl
my %list_hash = (
   "IP" => "192.168.0.1",
   "Machine" => "toto-portable",
   "Name" => "Toto",
);
```

We can leave the comma on the last line without any problems. This can be convenient in some cases :-)

#### reverse

The reverse on hash lists will prove very useful! Indeed, as you know, to perform searches, we can only take keys to find values. If we want to search in the reverse direction (values become keys, and keys become values):

```perl
%list_hash_inversed = reverse %list_hash;
```

But don't think that this doesn't take resources, because contrary to what one might think, a copy of a list does a complete unwinding, then a one-by-one copy of all elements. **It's exactly the same for a simple copy of a list:**

```perl
%hash1 = %hash2;
```

Unfortunately, once again, this makes Perl work extremely hard! So if you have the possibility to avoid this kind of thing, that's good :-)

### Hash Functions

#### Keys and values

Keys and values are 2 hash functions that return only the keys or only the values:

```perl
%list_hash = ("IP", "192.168.0.1", "Machine", "toto-portable", "Name", "Toto");
my @keys = keys %list_hash;
my @values = values %list_hash;
```

Here @keys will contain IP, Machine and Name, while @values will contain the rest.

In scalar form, this will give us the number of elements in keys and values!
Note: If a value is empty, it will be considered false.

#### each

To iterate through an entire hash, we'll use each which returns a key/value pair in the form of a 2-element list:

```perl
while ( ($key, $value) = each %hash ) {
   print "$key => $value\n";
}
```

#### exists

To know if a key exists in a hash:

```perl
if (exists $name{"192.168.0.1"}) {
   print "There is a name for the IP address 192.168.0.1!\n";
}
```

#### delete

The delete function will delete the key in question:

```perl
my $dhcp{name}='toto';
delete $dhcp{$name};
```

Note: This is not equivalent to inserting undef into a hash element.

#### References

Here's how to return hash references:

```perl
sub foo
{
    my %hash_ref;

    $hash_ref->{ 'key1' } = 'value1';
    $hash_ref->{ 'key2' } = 'value2';
    $hash_ref->{ 'key3' } = 'value3';

    return $hash_ref;
}

my $hash_ref = foo();

print "the keys: %$hash_ref\n";
```

## Input/Output

### The Diamond Operator

This operator allows you to merge multiple inputs into one large file:

```perl
while (<>) {
   chomp;
   print "This line comes from $_!\n"
}
```

This operator is generally used to process all input. Using it multiple times in a program would be an error.

#### @ARGV

@ARGV is an array containing the calling arguments. The diamond operator will first read the @ARGV array; if the list is empty, it will read the standard input, otherwise the list of files it finds.

### Standard Output

Normally, you understand that there are differences between displaying and interpolating:

```perl
print @tab; # displays a list of elements
print "@table" # displays a string (containing an interpolated array)
```

### printf

Printf lets you better control the output. To display a number in a generally correct way, we will use **%g**, which automatically selects decimal, integer, or exponential notation:

```perl
printf "%g %g", 6/2, 51 ** 17; # This will give us 3 1.0683e+29
```

**%d** means a decimal integer:

```perl
printf "%d times more", 15.987690987; # Gives 15 times more
```

Printf is often used to present data in columns. We'll define spaces:

```perl
printf "%6d\n", 32; # Gives "    32"
```

Here we have 4 spaces then the number 32 which gives 6 characters.

Same with **%s** which is dedicated to strings:

```perl
printf "%10s\n", toto; # "      toto"
```

If we make the 10 negative, we'll have left alignment:

```perl
printf "%-10s\n", toto; # "toto      "
```

With numbers **%f**, it is able to truncate:

```perl
printf "%12f\n", 6 * 7 + 2/3; # "   42.666667"
printf "%12.3f\n", 6 * 7 + 2/3; # "      42.667", here I told it to round to 3 digits after the decimal
```

If we want the % sign to be mentioned then we need to have **%%**:

```perl
printf "My percentage is %.2f%%\n"; 5.25/12; # result: "0.44%"
```

In the case where we want to enter the value to truncate in STDIN, don't forget to interpolate it:

```perl
printf "%${value}s\n", $_
```

For more information: http://perldoc.perl.org/functions/sprintf.html

## Regex

Regex or regular expressions are like another language to learn, which very easily provide access to very sophisticated tools.

If for instance we want to do a search and display the matching expression (equivalent to grep):

```perl
$_ = "tata titi toto";
if (/tot/) {
   print "Match found!";
}
```

### Metacharacters

Metacharacters are used to push searches even further. If I take the example above:

```perl
if (/t.t/) {
```

The character "." means anything. So here, it matches tata, titi, and toto. If you really want to use "." and not have it pass as regex, you need to use the metacharacter "\".
Which gives "\.":

```perl
3\.15 # Gives 3.15
```

### Quantifiers

The \* sign means that the previous character is repeated x times or not at all:

```perl
/toto\t*titi/ # This search can have from 0 to x tabs between toto and titi
```

If we want anything between toto and titi, just add a ".":

```perl
/toto.*titi/ # Will search for "toto'anything'titi"
```

To repeat the previous element 1 or more times, use "+":

```perl
/toto +titi/ # Can give "toto       titi" but not "tototiti"
```

The ? indicates that it's not mandatory:

```perl
/toto-?titi/ # Gives tototiti or toto-titi
```

### Groupings

We can group with ():

```perl
/toto+/ # Can make totooooooooo
/(toto)+/ # Can give totototototototo
```

### Pipe

The pipe is used to designate one element or another:

```perl
/toto|titi|tata/ # Will search for toto or titi or tata
```

If one of them matches, it is taken. If we want to search for spaces or tabs:

```perl
/( +|\t+)/
```

### Character Classes

- To designate ranges or some letters or numbers, we'll use []:

```perl
[a-c] # Gives a,b and c
[a-cw-z] # Gives a,b,c,w,x,y and z
[1-6] # Gives 1 to 6
[a-zA-Z] # Alphabet in lowercase and uppercase
```

To avoid certain characters, add ^:

```perl
[^a-c] # Anything except a,b and c
```

Now, even better! Some classes appear so often that they have been further simplified:

| Operator | Description                                                                              | Example                           |
| -------- | ---------------------------------------------------------------------------------------- | --------------------------------- |
| ^        | **Beginning** of line                                                                    | **^Deimos** for 'Deimos Fr!!!'    |
| $        | **End** of line                                                                          | **!$** for 'Deimos Fr!!!'         |
| .        | Any character                                                                            | **d.im.s** for 'Deimos Fr!!!'     |
| \*       | Repetition of previous character **from 0 to x times**                                   | **!\*** for 'Deimos Fr **!!!**'   |
| +        | Repetition of previous character **from 1 to x times**                                   | **!+** for 'Deimos Fr **!!!**'    |
| ?        | Repetition of previous character **from 0 to 1 time**                                    | **F?** for 'Deimos **F**r!!!'     |
| \        | Escape character                                                                         | **\.** for 'Deimos Fr\.'          |
| a,b,...z | Specific character                                                                       | **D**eimos for 'Deimos Fr'        |
| \w       | **Alphanumeric** character (a...z,A...Z,0...9)                                           | **\w**eimos for 'Deimos Fr'       |
| \W       | **Anything except** an **alphanumeric** character                                        | I**\W**ll for 'I**'**ll be back'  |
| \d       | A **digit**                                                                              | \d for 1                          |
| \D       | **Anything except** a **digit**                                                          | \Deimos for **D**eimos            |
| \s       | A **spacing character** such as: space, tab, carriage return, or line feed (\f,\t,\r,\n) | 'Deimos**\s**Fr' for 'Deimos Fr'  |
| \S       | **Anything except** a **spacing character**                                              | 'Deimos**\S**Fr' for 'Deimos Fr'  |
| {x}      | **Repeats** the previous character exactly **x times**                                   | !{3} in 'Deimos Fr **!!!**'       |
| {x,}     | **Repeats** the previous character **at least x times**                                  | !{2} in 'Deimos Fr **!!!**'       |
| {, x}    | **Repeats between 0 and x times** the previous character                                 | !{3} in 'Deimos Fr **!!!**'       |
| {x, y}   | **Repeats between x and y times** the previous character                                 | !{1, 3} in 'Deimos Fr **!!!**'    |
| []       | Allows to put a **range** (from a to z[a-z], from 0 to 9[0-9]...)                        | [A-D][0-5] in 'A4'                |
| [^]      | Allows to specify **unwanted characters**                                                | [^0-9]eimos in 'Deimos'           |
| ()       | Allows to **record the content of parentheses** for later use                            | (Deimos) in 'Deimos Fr'           |
| \|       | Allows to make an **exclusive or**                                                       | (Org\|Fr\|Com) in 'Deimos **Fr**' |

There's a site that allows you to visualize a regex: [Regexper](https://www.regexper.com) (sources: https://github.com/javallone/regexper)

So:

```perl
/toto \w+ titi/ # Will give toto, space, any word, space and titi
```

Isn't that beautiful? To represent a space, we can also use \s:

```perl
/toto\s\w+\stiti/
```

- Now, if we want the opposites:

```perl
[^\d] # Anything except digits
```

or we can use uppercase:

```perl
\D
```

Pretty neat, right! :-)

We can also find [\d\D] which means any digit or non-digit (unlike the . which is identical except that the . doesn't accept new lines).

### General Quantifiers

If we want to match a pattern multiple times:

```perl
/a{5,15}/ # "a" can be repeated between 5 and 15 times
/a{3,}/ # "a" can be repeated from 3 to infinity
```

Which can give for example, if we're looking for an 8-character word:

```perl
/\w{8}/
```

### Anchors

As there are too few characters, some are reused:

| Searches            | Anchors          |
| ------------------- | ---------------- |
| Beginning of a line | /**^**My\sstart/ |
| End of a line       | /my\send**$**/   |

#### Word Anchors

To define a whole word, we'll use \b:

```perl
/\btoto\b/ # Will only search for "toto"
```

To reverse the order of things, so if we want anything except toto:

```perl
/\Btoto\B/
```

But we might want titine and titi:

```perl
/\btiti/ # Equivalent to titi.*
```

Super-timorous, even stronger:

```perl
/\bsearch\B/ # Will match searches, searching, and searched, but not search or research
```

### Memorization Parentheses

A good example and we understand better. If we use:

```perl
/./ # We locate any individual character except the new line
```

To memorize this regex, we'll use parentheses:

```perl
/(.)/ # There's memorization here
```

To reference it, we'll use:

```perl
/(.)\1/ # This will contain the first regex memory
```

This is not simple to understand but it will look for a character identical to the previous search

For example if we have HTML code:

```perl
<image source='toto.png'>
<image source="toto-l'artichaud.png">
```

We can make our search with this:

```perl
/<image source=(['"]).*\1>/ # the *\1 indicates that we use the first memory of the regular expression. It is repeated when called (\1).
```

Given the complexity of the thing, I'll try to show good examples. First of all, you need to count the opening parentheses, this will be our regex memory number (e.g.: ((... = 2, because 2 opening parentheses):

```perl
/((toto|titi) (tata)) \1/ # Can match toto tata toto tata or titi tata titi tata
/((toto|titi) (tata)) \2/ # Can match titi tata titi or toto tata toto
/((toto|titi) (tata)) \3/ # Can match toto tata tata or titi tata tata (even if this kind of expression is rare)
/\b((toto|titi)\s+tata\b\b/ # Allows to search for exactly these words (and not totototo for example
/^b[\w.]{1,12}\b$/ # Corresponds to strings of letters, never ending with a . and maximum 12 characters.
```

### Non-memorization Parentheses

If you want to use parentheses without them being stored in memory, you need to use these symbols "?:" like this:

```perl
if (/(?:toto|tata) est en train de jouer/)
{
   print "Someone is playing";
}
```

Here, toto or tata won't be stored.

### Using Regex

Just as we've seen with the qw// operator, it's possible to do the same thing for matches with m//:

```perl
m(toto)
m<titi>
m{tata}
```

In short, the possibilities are ^,!(<{[ ]}>)!,^. The m// shortcut is not mandatory with the double **//**

#### Ignoring Case

To ignore case, there's **/i**:

```perl
/toto\n/i # The match can be toto or TOTO, or even ToTo...etc...:-)
```

#### Ignoring Spaces and Comments

If you want to ignore all spaces and comments in your code, add **/x**:

```perl
/   toto # \n/x
```

#### Matching Any Character

The fact that the . doesn't match a new line character can be annoying. That's why **/s** is useful. For example:

```perl
$_ = "Here:\ntoto\ntiti\ntata.\n";
if (/\btoto\b.*\btiti\b/s) {
   print "This string indicates toto after titi!\n"
}
```

We can even make combinations:

```perl
if (/\btoto\b.*\btiti\b/si) { # Here we combine /s and /i
```

#### Searching up to a Specific Character

I struggled a lot with this regex before finding it. If for example, I have a line like:

```xml
<task id="3" name="Cacti" color="#8cb6ce" meeting="false" start="2010-04-26" duration="65" complete="0" priority="1" expand="true">
```

And I want to search for the content of name:

```perl
/name="(.*)"/
```

This regex won't be enough because it will give me:

```perl
Cacti" color="#8cb6ce" meeting="false" start="2010-04-26" duration="65" complete="0" priority="1" expand="true">
```

To fix the problem, here's the solution:

```perl
/name="(.*?)"/
```

I simply put a ? which will ask to search not to the last ", but to the first one!

#### The =~ Binding Operator

Matching with $\_ is the default:

```perl
my $someone = "The person in question is toto.";
if ($someone =~ /\btoto/) {
   print "Now things get complicated";
}
```

If no binding operator is indicated, it will work with $\_.

Here's another example of matching, but with regular expression memories:

```perl
$me = "For now I have lived 24 years!";
if ($me =~/(\d*) ans/) {
   print "So I am $1 years old\n";
}
```

One last one:

```perl
$_ = "toto tata, titi";
if (/(\S+) (\S+), (\S+)/) {
   print "Here are the names: $1 $2 $3";
}
```

The memory remains intact until there is a match, whereas a successful match resets all of them. If you start playing too much with memories, you may have surprises. It is therefore advised to store them in variables:

```perl
if {$toto =~ /(\(w+)/) {
   my $search_toto = $1;
}
```

Now, watch out, we'll see the kind of things that I find great with Perl:

```perl
if ("Here is toto, tata and titi" =~ /,\s(\w+)/) {
   print "The match is '$&'.\n"; # That's ", tata"
}
```

$&: is the match
$`: what is before the match
$': what is after the match

In conclusion, if we want the original string:

```perl
print "The original string is: ($`)($&)($').\n";
```

These "magic" variables have a price! They slow down subsequent regular expressions. This can make you lose a few milliseconds...minutes, depending on your program. Apparently, this problem is fixed in Perl 6.

If you have the possibility to use numbering instead of these regex, don't hesitate!

Know that you can get even more info on regex here: http://www.perl.com/doc/manual/html/pod/perlre.html

### Substitution

For substitution (replacement), we'll use **s///**:

```perl
$_ = "Toto plays on the Wii with titi tonight";
s/titi/tata/; # Replace titi with tata
s/with (\w+)/against $1/; # Replace "with titi" with "against titi"
```

Or, more simply:

```perl
my $toto = 'TOTO';
my $titi = "\L$toto";
print "$titi\n";
```

Now, some slightly more complex examples:

```perl
$_ = "bumbo, small car";
```

- For a global substitution, that is, on all occurrences found, simply add "g":

```perl
$_ = "bumbo, nice bumbo!";
s/bumbo/auto/g;
print "$_\n";
```

It's possible to use other delimiters (like for m and qw) such as "{}, [], <>, ##":

```perl
s<^https://>(http://);
```

It's also possible to combine "g" with another (like case sensitivity for example):

```perl
s{https://}[http://]gi;
```

- To replace with uppercase, use \U:

```perl
s/(toto|titi)/\U$1/gi; # Will give TOTO or TITI
```

- To replace with \L for all lowercase:

```perl
s/(toto|titi)/\L$1/gi; # Will give toto or titi
```

- You can disable case modification with \E:

```perl
s/(\w+) with (\w+)/\U$2\E with $1/i; will give TOTO with titi
```

- Written in lowercase (\l and \u), these escapes only affect the following character:

```perl
s/(toto|titi)/\u$1/gi; # Will give Toto or Titi
```

You can also combine them so that everything is lowercase except the first letter:

```perl
s/(toto|titi)/\u\L$1/gi; # Will give Toto or Titi
```

You can even do this in a simple print:

```perl
print "His name is \u\L$name\E\n";
```

### split

Split allows cutting based on a space, tab, period... pretty much anything except commas:

```perl
@fields = split /separator/, $string;
```

Split moves the pattern in a string and returns the list of fields separated by separators. Each match of the pattern corresponds to the end of one field and the beginning of another:

```perl
@fields = split /:/, "abc:def:ij:k::m"; # Will give ("abc", "def", "ij", "k", "", "m")
```

As I mentioned above, it's also possible to make separations with spaces:

```perl
my $phrase = "Here is my line\t of spaces.\n"
my @args = split /\s+/, $phrase; # Gives ("Here", "is", "my", "line", "of spaces."
```

By default, if no separation options are specified, "\s+" will be used.

### join

Join works exactly like split except that its result will give the inverse of split. It will join pieces (with or without) separators:

```perl
my @mac = join ":", 00, AB, EF, EF, EF, EF; # Will give 00:AB:EF:EF:EF:EF
```

or even:

```perl
my $mac = join ":", 00, AB, EF, EF, EF, EF; # Will give 00:AB:EF:EF:EF:EF
```

## More Complex Control Structures

### unless

Unless is the opposite of if, that is, we'll enter the loop if the searched expression is not the right one. It's actually equivalent to an else in an if loop. This also equates to making an if negative:

```perl
if ( ! ( $toto =~ /Â-Z_]\w*$/i) ) <--> unless ($toto =~ /Â-Z_]\w*$/i)
```

It is also, just like an if, possible to use else with unless, but I don't recommend it as it's often a source of errors.

### until

If you want to invert the condition of the while loop:

```perl
until ($j > $i) {
   $j *= 2;
}
```

This is actually a disguised while loop that repeats as long as the condition is false.

### Expression Modifiers

For a more compact notation, an expression can be followed by a modifier that controls it:

```perl
print "$n is a negative number.\n" if $n < 0;
&error("Invalid input") unless &valid($input);
$i *= 2 until $i > $j;
print " ", ($n += 2) while $n < 10;
&greet($_) foreach @person;
```

### Bare Block

This is a bare block:

```perl
{
body;
body;
body;
}
```

It's a block that will be executed only once. The advantage is that we can create variables, but they will only be kept in this block.

### elsif

In an if loop, if we want to have multiple solutions, we can use elsif:

```perl
if (! defined $toto) {
   print "The value is undef.\n";
} elsif ($toto eq '') {
   print "The value is a string.\n";
else {
   ...
}
```

We can put as many elsif as we want (see perlfaq to emulate case or switch).

### Auto Increment/Decrement

As in C, to increment a variable for example:

```perl
my $a = 5;
my $b = $a++; # can also be written as ++$a
```

Same with "--" for decrementing

### for

For is quite classic and resembles C again:

```perl
for ($i = 1; $i <=10; $i++) {
   print "I can count to $i!\n";
}
```

Another example. Imagine that we want to count from -150 to 1000 but in steps of 3:

```perl
for ($i = -150; $i <= 1000; $i +=3) {
   print "$i\n";
}
```

Otherwise, a simple for loop for a successful search:

```perl
for ($_ = "toto"; s/(.)//; ) {
   print "$1 is found here\n";
}
```

Be careful with infinite loops if you use variables:

```perl
for (;;) {
   print "Infinite loop\n";
}
```

If you really want to write an infinite loop, the best way is this:

```perl
for (1) {
   print "Infinite loop\n";
}
```

### foreach

The loop and for and foreach are identical except that if there's no ";", it's a foreach loop:

```perl
for (1..10) { # Actually a foreach loop from 1 to 10
   print "I can count to $_!\n";
 }
```

So it's a foreach loop but written with a for.

### Loop Controls

"last" allows to terminate a loop immediately (like break in C or shell):

```perl
while (<STDIN>) {
   if (/_END_/) {
      last;
   } elsif (/fred/) {
      print;
   }
}
```

As soon as a line contains the "**END**" marker, the loop ends.

### next

Sometimes, you're not prepared for the loop to end, but you've finished the current iteration. This is where "next" comes in! It jumps to the inside of the bottom of the loop, then it goes to the next iteration of the loop:

```perl
while (<>) {
   foreach (split) { # splits $_ into words, assigning each to $_ in turn
      $total++;
      next if /\W/; # strange words skip the rest of the loop
      $valid++;
      $count{$_}++; # count each individual word
      ## next leads here ##
   }
}

print "total = $total, valid words = $valid\n";
foreach $word (sort keys %count) {
   print "$word was encountered $count{$word} times.\n";
}
```

### redo

Redo indicates to go back to the beginning of the current loop block without going further in the loop:

```perl
my @words = qw{ toto tata titi };
my $errors = 0;

foreach (@words) {
   ## redo comes back here ##
   print "Enter the word '$_': ";
   chomp(my $try = <STDIN>);
   if ($try ne $_) {
      print "Sorry - There is an error.\n\n";
      $errors++;
      redo; # go back to the top of the loop
   }
}

print "You finished the test with $errors errors.\n";
```

A small test to understand well:

```perl
foreach (1..10) {
   print "Iteration number $_.\n\n";
   print "Choose: last, next, redo or none";
   chomp(my $choice = <STDIN>);
   print "\n";
   last if $choice =~ /last/i;
   next if $choice =~ /next/i;
   redo if $choice =~ /redo/i;
   print "That wasn't any of the choices... moving on!\n\n";
}
print "That's all, folks!\n";
```

### Labeled Blocks

Labeled blocks are used to work with loop blocks. They are made of letters, underscores, numbers but cannot start with the latter. It is advised to name them with capital letters. In reality, labels are rare. But here's an example:

```perl
LINE: while (<>) {
      foreach (split) {
         last LINE if /__END__/; # exits the LINE loop
         ...;
      }
}
```

### Logical Operators

Like in shell:

- && (AND): executes what follows if the previous condition is true. Also allows to say that the expression before and the one that will follow must be validated to perform what follows.
- || (OR): executes what follows if the previous condition is false. Also allows to say that if the expression before doesn't match, the following must match to be able to continue.

```perl
if ($dessert{'cake'} && $dessert{'ice cream'}) {
   # Both are true
   print "Hooray! Cake and ice cream!\n";
} elsif ($dessert{'cake'} || $dessert{'ice cream'}) {
   # At least one is true
   print "It's still good...\n";
} else {
   # Neither is true
}
```

We can also write like this:

```perl
$hour = 3;
if ( (9 <= $hour) && ($hour < 18) ) {
   print "Shouldn't you be working?\n";
```

They are also called short-circuit operators, because in the example below, the left operator needs to check the right one to avoid a division by 0:

```perl
if ( ($n != 0) && ($total/$n < 5) ) {
   print "The average is less than 5.\n";
}
```

**Unlike other languages, the value of a short-circuit operator is the last part evaluated, not a simple boolean value.**

### Ternary Operator

```perl
expression ? expr_if_true : expr_if_false
```

The ternary operator looks like an if-then-else. We first check if the expression is true or false:

- If it's true, the 2nd expression is used, otherwise the third. Each time, one of the 2 right expressions is evaluated and the other ignored.
- If the first expression is true the 2nd is evaluated and the 3rd ignored.
- If the 1st is false, the 2nd is ignored and the 3rd evaluated as the value of the whole.

```perl
my $place = &is_weekend($day) ? "home" : "work";
```

Here's another example:

```perl
my $average = $n ? ($total/$n) : "No average";
print "Average: $average\n";
```

A slightly more elegant example:

```perl
my $size =
   ($width < 10) ? "small" :
   ($width < 20) ? "medium" :
   ($width < 50) ? "large" :
                   "very large"; # Default
```

One last example:

```perl
($a < $b) ? ($a = $c) : ($b = $c);
```

## File Handles and File Tests

File handles are named like other Perl identifiers (with letters, digits, and underscores, without starting with a digit) but since they don't have a prefix, they can be confused with current or future reserved words. It is therefore advised to use only capital letters for a file handle name.

Today there are 6 file handle names used by Perl for its own use:

- STDIN
- STDOUT
- STDERR
- DATA
- ARGV
- ARGVOUT

You may see these handles written in lowercase in some scripts, but that doesn't always work, which is why it's advised to put everything in uppercase.

A program's output is called STDOUT, and we'll see how to redirect this output. You can take a look at the [Perlport documentation](https://perldoc.perl.org/perlport.html).

### Running Programs

Here are 2 ways to read your program:

```bash
$ ./my_soft <toto>titi
```

This will read the input toto and send the output to titi.

```bash
$ cat toto | ./my_soft | grep stuff
```

This will take the input toto, send it to my soft and grep the output of my soft.

For error redirection:

```bash
$ ps aux | ./my_soft 2>/var/log/my_soft.err
```

### Opening a Handle

Here's how to open files:

```perl
open FILER, "toto";
open FILER, "<toto";
```

Here are 2 methods of reading a file. For security reasons, I strongly advise you to use the last method.

```perl
open FILEW, ">titi";
```

This is for writing to a file.

```perl
open FILEW, ">>titi";
```

This is also for writing to a file, but added to the end of an already existing file. Here's an example of use:

```perl
my $soft_output = "my_output";
open FILE>, "> $soft_output";
```

It's clearer to make a scalar variable and say what we're going to do.

### Closing a Handle

To close a file handle, simply do:

```perl
close FILER;
```

By default Perl will close all your files when your script closes.

### Handle Problems

You may encounter problems when opening a file, such as permission issues or others. Here's an example to avoid mistakes:

```perl
#!/usr/bin/perl -w
my $success = open LOG, ">>journal_file"; # Capture the return value
unless ($success) {
   # Opening problem
   ...
}
```

For a fatal error, we have another solution that is more common (die):

```perl
unless (open LOG, ">>journal_file") {
   die "Cannot create journal_file : $!";
}
```

The sign "$!" is used to display the readable form of the complaint from the system.

Here are **generally** the possible return states:

- 0: everything went well
- 1: is a syntax error in the arguments
- 2: is an error produced during processing
- 3: the file was not found

"$!" can display the line numbers of errors. If you don't want them, write like this:

```perl
die "Problem\n" if @ARGV < 2;
```

This line will analyze the number of arguments. If it is greater than 2 this program will terminate.

We can also do like this, which is the most used solution:

```perl
open LOG, ">>journal_file" or die "Cannot create journal_file: $!";
```

You can also replace "or" with "||", but then you'll need to use open with parentheses.

Note: It's possible to use "||" instead of "or", however, it's older code that will sometimes use the higher precedence operator. The only difference is that when "open" is written without parentheses, the higher precedence operator will be linked to the filename argument, not to the return value. So the return value of open is not checked afterwards. If you use ||, make sure to indicate the parentheses. The ideal remains to use "or" to avoid encountering unwanted effects.

### Reading a Text Block

You may need to read a text block to interpret it afterwards. The idea is to make the input, then press CTRL+D afterwards to say that you've finished the input. Here's how to proceed:

```perl
my @userinput = <STDIN>;
foreach (@userinput) {
        print;
}
```

### Warnings with warn

You can use _warn_ just like _die_ except that instead of quitting the application radically, it will give you a warning message.

### Using File Handles

After opening a file for reading, you can read the lines as if it was the standard input STDIN:

```perl
open PASSWD, "/etc/passwd" or die "How are you logged in? ($!)";

while (<PASSWD>) {
   chomp;
   if (/^root:/) {
      print "root has been found!!!";
   }
}
```

We can also add to the end of a file:

```perl
open RC_LOCAL, ">> /etc/rc.local" or die "This file doesn't exist! ($!)";

while (<RC_LOCAL>) {
   chomp;
   print RC_LOCAL "/usr/sbin/nrpe -d\n";
}
```

Here the content of the print is added to the end of rc.local. Another example for an error:

```perl
printf STDERR "%d percent reached.\n", $done/$total * 100;
```

### HEREDOC

HEREDOC serves to include a multi-line text literally in the program, the syntax is as follows:

```perl
my $string = <<"ENDOFTEXT";
This is a HEREDOC
We can move to the next line without problems,
to end the HEREDOC, we need to put
the end marker (here ENDOFTEXT)
alone on a line, and really all alone!
    ENDOFTEXT
for example the HEREDOC is not finished yet,
because there were spaces before the marker.
ENDOFTEXT

print $string; # we are back in Perl code
```

There are subtleties, for example for a string between apostrophes (quote '), the variables are not interpolated while when we put it between quotation marks (double quote ") the variables are interpolated!
Well for a HEREDOC, the behavior of variables is determined by what you put around the marker:

```perl
my $hello = "Hello world";
print <<"EOD";
This is a HEREDOC with interpolation:
$hello
EOD

print <<'EOD';
This is a HEREDOC without interpolation:
$hello
EOD
```

### Replacing the Default Output File Handle

By default, printf prints on the default output. The default output can be modified with the select operator.

```perl
select LOG;
print "This is my log file\n";
print "ok done\n";
```

From the moment you've selected a file handle as the default output, it remains so (for more info man perlfunc).

Also by default, a file handle's output is buffered. By initializing the special variable "$|" to 1, you set the selected file (the one selected when the variable is modified) so that the buffer is always flushed after each output operation. Thus, to make sure that the log file will immediately receive all its entries (in case you read the log to monitor the progress of a long program) you can for example write:

```perl
select LOG;
$| = 1; # don't leave LOG entries in the buffer
select STDOUT;
# blablabla...
print LOG "this is written to LOG once!\n";
```

### Reopening a Standard File Handle

If you've already opened a handle (opening a LOG when a LOG is already open), the old one would be automatically closed. Just as you can't use the 6 standards (except in exceptional cases). Messages such as warn or die will go directly to STDERR.

```perl
open STDERR, ">>/home/pmavro/errors" or die "Problem opening file: $!\n";
```

### File Tests

Maybe you'd like to know how to check if a file exists, its age or other:

```perl
die "Sorry the file already exists\n" if -e $my_file;
```

You can see that we didn't put $! here, because it's not the system that rejected a request, but my file that already exists. For regular updating, if you want to check a file is not older than 4 days:

```perl
warn "The config file is too old\n" if -M CONFIG > 4;
```

For find users, this will please you, imagine that we don't want files larger than 100k, we move them to a folder if it has been unused for at least 90 days:

```perl
my @original_files = qw/ toto titi tata tutu tete /;
my @big_old_files; # What we want to transfer to the old folder
foreach my $my_files (@original_files) {
   push @big_old_files, $_ if -s $my_files > 100_000 and -A $my_files > 90;
}
```

Here are the possibilities:
![Files tests perl](/images/files_tests_perl.avif)
The tests -r, -w, -x and -o will only work for the user running the Perl script. Also be careful with certain system limitations, such as -w which doesn't prevent writing, only if it's on a CD because it's mounted read-only.

Another thing to be careful about is symbolic links which can be deceiving. That's why, it would be better to test for the presence of a symbolic link before testing what interests you.

For searches at the time level, it's possible that there are floating commas as well as negative numbers (if the execution is still ongoing for example).

The -t test returns true if the handle is a TTY. It's able to be interactive with the user with for example a "-t STDIN".

For -r, if you forget to specify a file, $\_ will be taken into account.

If you want to transform a size into KB be sure to put the parentheses:

```perl
my $size_in_KB = (-s) / 1024;
```

### The stat and lstat Functions

Stat allows to obtain a lot of information about a file:

```perl
my ($dev, $ino, $mode, $nlink, $uid, $gid, $rdev, $size, $atime, $mtime, $ctim, $blksize, $blocks) = stat($file_name);
```

- $dev: Device number
- $ino: Inode number
- $mode: Gives the file permissions like 'ls -l would give
- $nlink: Number of hard links. Always worth 2 or more for folders and 1 for files.
- $uid: User ID
- $gid: Group ID
- $size: the file size in bytes (like -s)
- $atime: Equivalent to -A
- $mtime: Equivalent to -M
- $ctime: Equivalent to -C

Invoking stat on a symbolic link returns info on the original file. To get information about the symbolic link, you can use lstat. However, if there's no info, it will give you the info of the original file instead.

### Localtime

This function allows to convert unix time to human readable time:

```perl
my $timestamp = 19934309;
my $date = localtime $timestamp;
```

or:

```perl
my($sec, $min, $hour, $day, $mon, $year, $wday, $yday, $isdst) = localtime $timestamp
```

For GMT time, you can use the gmtime function:

```perl
my $now = gmtime;
```

If you want a concrete example for instance to not be able to run a script during production hours:

```perl
sub production_time_check
{
	my ($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst) = localtime(time);
	{
		if (($hour >= 8) and ($hour <= 18) and ($wday >= 1) and ($wday <= 5))
		{
			exit;
		}
		else
		{
			print "You can continue :-)\n";
		}
	}
}
```

### Bit by Bit Operators

This is useful for doing binary calculations, such as the values returned by the stat function:
![Bit to bit operator](/images/operateur_bit_a_bit.avif)

```perl
# $mode is the mode value returned by a stat of CONFIG
warn "The configuration file is world-writable!\n" if $mode & 0002; # configuration security issue
my $classical_mode = 0777 & $mode;                # mask the additional upper bits
my $u_plus_x = $classical_mode | 0100;            # activate 1 bit
my $go_minus_r = $classical_mode & (~ 0044);      # deactivate 2 bits
```

### Using the Special "\_" File Handle

With each use of stat and lstat we make 2 system calls:

- 1: to know if it's possible to read
- 2: to know if it's possible to write

But we lose time when we make requests to the system. The goal is therefore to ask only once with "\_" and then reuse this to get the new information. It looks a bit ridiculous like this but if you make a lot of system calls, you'll see that with this your program will be much faster:

```perl
my @original_files = qw/ fred barney betty wilma pebbles dino bamm-bamm /;
my @big_old_files;                       # What we want to save on tape
foreach (@original_files) {
  push @big_old_files, $_
    if (-s) > 100_000 and -A _ > 90;     # More efficient than previously
}
```

We use $_ for the first test which is not more efficient, but we get info from the operating system.
Then we use the magic file handle: "_". The data left after retrieving the file size is used. Here we optimize the requests.

## Directory Operations

### Moving in the Directory Tree

To move (equivalent to "cd" in shell):

```perl
chdir "/etc/" or die "Unable to move to /etc: $!\n";
```

**As this is a system request, the value of $! is initialized in case of error.**

### Globalization

Here's an example of globalization (in shell):

```perl
$ echo *.pm
toto.pm titi.pm tata.pm
```

In perl it's quite similar actually, imagine we have an array with all sorts of files:

```perl
foreach $arg (@ARGV) {
   print "one of the arguments is $arg\n";
}
```

and here's how to bring out what interests us:

```perl
my @all_files = glob "*";
my @pm_files = glob "*.pm";
```

If we want multiple searches, just separate them with spaces:

```perl
my @all_files_including_dot = glob ".* *";
```

Here's another type of syntax:

```perl
my @all_files = <*>;
```

equivalent to:

```perl
my @all_files = glob "*";
```

another example that needs no comment:

```perl
my $rep = "/etc";
my @rep_files = <$rep/* $rep/.*>;
```

### Directory Handles

If we want to list the content of a folder:

```perl
my $folder="/etc";
opendir FOLDER, $folder or die "Unable to open $folder: $!";
foreach $file (readdir FOLDER) {
   print "$file is in the folder $folder\n";
}
closedir FOLDER;
```

The result will be unsorted files or folders (even files starting with a .). If now, we want to get only files ending with pm, we'll need to do like this:

```perl
while ($name = readdir DIR) {
   next unless $name =~ /\.pm$/;
   ...other processing...
}
```

If we had wanted everything except . and ..:

```perl
next if $name eq "." or $name eq "..";
```

Now, if you want to do recursive, then I advise you the [File::find](https://search.cpan.org/~rgarcia/perl-5.10.0/lib/File/Find.pm) library.

## Manipulating Files and Directories

### Deleting Files

The equivalent of rm is unlink:

```perl
unlink "file1", "$file2";
```

Another example, with glob:

```perl
unlink glob "*.bak";
```

To validate the deletion of files:

```perl
my $success = unlink "titi", "tata", "toto";
print "The $success files have been deleted\n";
```

We'll see here if 0 or 3 files have been deleted, but not 1 or 2. We'll need to make a loop in case you absolutely want to know:

```perl
foreach my $file (qw(toto tata titi)) {
   unlink $file or warn "The file $file could not be deleted: $!\n";
}
```

### Renaming Files

Here are 2 examples:

```perl
rename "toto", "titi";
rename "/my/file/toto", "titi";
```

We can also proceed with a loop and elegantly rename:

```perl
(my $newfile = $file) =~ s/\.old$/.new/;
```

### Links and Files

To find out which is the source file of a symbolic link:

```perl
my $where_to = readlink "toto"; # Will return tata (because ln -s tata toto
```

### Creating and Deleting Directories

To create a folder it's very simple:

```perl
mkdir "toto", 0755 or warn "Impossible: $!";
```

Be careful however if you want to assign permissions to variables because this won't work like this:

```perl
my $user = 'toto';
my $permissions = '0755';
mkdir $user, $permissions;
```

And there's the catastrophe because our folder has weird permissions like 01363. All because the string is by default in decimal and we need octal, to solve this thorny problem:

```perl
mkdir $user, oct($permissions);
```

Now if we want to delete a folder, it's simple:

```perl
rmdir glob "toto/*" or warn "Sorry $!";
```

This will delete all empty directories in the toto folder. Rmdir returns the number of elements deleted. Rmdir will only delete a directory if it is empty, so use unlink before rmdir:

```perl
unlink glob "toto/*" "toto/.*";
rmdir 'toto';
```

If this alternative seems too boring, use the [File::Path](https://search.cpan.org/~dland/File-Path-2.07/Path.pm) module.

### Determining the Process

To determine the currently running process, use the $$ variable. When you create a file for example, that gives:

```perl
mkdir "/tmp/temp_$$";
```

### Changing Permissions

To change permissions, simply, like in shell use chmod:

```perl
chmod 0755, "toto", "tata";
```

If you want to use u+x or that sort of thing, refer to http://search.cpan.org/~pinyan/File-chmod-0.32/chmod.pm

### Changing the Owner

Once again it's like in shell, we use chown:

```perl
my $user = 1004;
my $grp = 100;
chown $user, $grp, glob "*.o";
```

If you don't want to use the uid and guid to make the change and prefer names, then do like this:

```perl
defined(my $user = getpwnam "toto") or die "bad user";
defined(my $grp = getgrnam "users") or die "bad group";
chown $user, $grp, glob "/home/toto/*";
```

The defined function checks that the return value is not undef.

### Changing Date and Time

Sometimes you want to lie to certain programs, here's how to change access and modification time:

```perl
my $now = time;
my $before = now - 24 * 60 * 60 # seconds per day
utime $now, $before, glob "*"; # initializes access to now, modified from yesterday
```

This can be very useful in case of problems with backups.

### The File::Basename Module

If we want to get the path of a binary for example, we'll need this module, here's how it works:

```perl
use File::Basename;
my $name = "/usr/jdk/latest/bin/java";
my $base_name = basename $name; # gives 'java'
```

### Using Only Certain Functions from a Module

Imagine that you have a function with the same name as a function of one of your modules. To load only what you need, here's how to do it:

```perl
use File::Basename qw/ basename /; # Here we only load the basename function
use File::Basename qw/ /; # Here we don't ask for any functions
my $perlpath = "/usr/bin/perl";
my $folder_name = File::Basename::dirname $perlpath; # Will use dirname from the module
```

### The File::Spec Module

With the File::Basename module, it's convenient, you know what you need to get a file, but if you want to get the complete path where your file is, you'll need to use the [File::Spec](https://search.cpan.org/~smueller/PathTools-3.29/lib/File/Spec.pm) module.

## Process Management

Calling external programs can be very practical when you don't have time to rack your brain or simply when you have no choice.

### System

This one is my favorite, because it allows launching a child process of your Perl program. If you need to fork your Perl program with a command, the system function is very convenient. Personally I had to develop in Perl for a generic script (GDS) on SunPlex (Sun Cluster) for the company I work for, and I needed to fork at one point. I was delighted to see that system did it.

It is however important to understand that the system function will return the data to STDOUT and not to your Perl program:

```perl
system 'ls -l $HOME';
system "ls -l \$HOME";
```

Note that simple apostrophes are for shell values and double apostrophes for your perl program.

The problem with this command is also in the command you call because if it asks you questions (like asking for confirmation etc.), your Perl program will wait for the end of your command. To bypass this, add a &:

```perl
system "rm -R /tmp/ &";
```

If you are on Windows, here's the solution to adopt:

```perl
system(1, "del c:\toto");
```

Where the system function has been well written is that it doesn't require launching a shell when it's a small command. But if it contains characters with $, / or \ for example, then a shell will be launched for this execution.

#### Avoiding the Shell

If you can avoid the shell call, it's not bad. Here's an example:

```perl
my $file='toto';
my @folders=qw/titi, tata, tutu/;
system "rm -Rf $file @folders"; # And here's the disaster!
```

While there's a "cleaner" way that doesn't call the shell:

```perl
system "rm", "-Rf", $file, @folders;
```

If now, I want to use the return values to see if everything went well:

```perl
unless (system 'ls -l') {
   print "Command has been successfully executed\n";
}
```

or

```perl
!system "rm", "-Rf", $file, @folders or die "Rm command failed\n";
```

It's useless to use $! here because Perl cannot evaluate what happened since it doesn't happen in Perl.

### Exec

The functioning of the system function compared to exec is identical except for a very important point! Indeed, it won't create child processes, but will execute itself.

**When in your code we arrive at the exec part, it jumps into the command and the command takes the PID of your Perl program until the end of its execution and then gives back control (not to Perl, but generally to the shell, where you launched the Perl script).**

### Environment Variables

In Perl, there's a hash table called %ENV containing environment variables that contains values inherited from the previous shell that was launched. There's this:

```perl
$ENV{'PATH'} = "/usr/cluster/bin:$ENV{'PATH'}";
delete $ENV{'CVS'};
```

Which will let you modify your shell PATH when you're going to call child processes. It's quite convenient, but it obviously doesn't work for parent processes.

As a reminder under linux, the command to see environment variables is "env" and under Windows it's "set".

### Using ` to Capture Output

Here's an example of what we can do by recovering the output:

```perl
my $now = `date`;
print "The date of the day is $now";
```

Here's another example:

```perl
my @functions = qw( int rand sleep) {
   $documentations{$_} = `perldoc -t -f $_`;
}
```

**Don't have an abusive use of ` because Perl has to work a bit harder to recover the output. If there's no need for them, useless to use them, prefer system to that.**

If you want to recover errors, use this 2>&1:

```perl
my $err_out = `command 2>&1`;
```

If you need to use a command that might (without you wanting it) ask you a question, that will cause problems. To avoid this kind of inconvenience, send /dev/null to this command:

```perl
my $result = `command arg arg arg </dev/null`;
```

### Using ` in a List Context

If a program's output returns multiple lines, then the variable will contain all the lines one after another on a single one. For example with the who command, it's preferable to use an array:

```perl
my @persons = `who`;
```

Then for the analysis:

```perl
foreach (`who`) {
   my ($user, $tty, $date) = /(\S+)\s+(\S+)\s+(.*)/;
   $ttys{$user} .= "$tty at $date\n";
}
```

Notice that when =~ is not present, $\_ is automatically taken.

### Processes as File Handles

I really like this way of doing things, because it seems the clearest to me. You have to operate as if it was a file and not forget a | at the end of your command:

```perl
open (DATE, "date |") or die "Impossible to open the pipe from date: $!"\n;
close (DATE);
```

You can also put the | on the other side:

```perl
open (MAIL, "|mail -s deimos \"This is a test\"") or die "Impossible to open the pipe to: $!\n";
close (MAIL);
```

There's not even a need to go that far anyway to do this, a simple print will do:

```perl
my $now = <DATE>;
print MAIL "It is now $now";
close MAIL;
```

This use (with open) is more complex to use than with `, however, it allows to have a result arriving gradually. What I'm telling you makes sense with the find command for example, which gives little by little its results. This allows you to analyze and process in real time while with ` you would have to wait until the end of the find command to process:

```perl
open F, "find / -atime +90 -size +1000 -print|" or die "impossible fork: $!";
while (<F>) {
   chomp $_;
   printf "%s size %dk last access %s\n", $_, (1023 + -s $_)/1024, -A $_;
}
```

### Using fork

In addition to the high-level interfaces like above, we can in Perl make low-level system calls. For example fork, which is very practical. For example this:

```perl
system "date";
```

in forked version gives:

```perl
defined(my $pid = fork) or die "impossible fork: $!";
unless ($pid)  {
   # Child process:
   exec "date";
   die "Unable to execute date: $!";
}
# Parent process is:
waitpid($pid, 0);
```

The 0 means that the PID shouldn't be 0.

If you want more information, consult the perlipc manual.

### Sending and Receiving Signals

The different signals are identified by a name (for example SIGINT, for "interrupt signal") and an integer (ranging from 1 to 16, 1 to 32, or 1 to 63, depending on your flavor of Unix). For example, if we want to send a SIGHINT:

```perl
kill 2, 3094 or die "Unable to send a SIGHINT: $!";
```

The PID number here is 3094 and you can change 2 to INT if you want. Be careful about permissions, because if you don't have the authorization to kill the process, you'll get an error.

If the process no longer exists, you'll have a return to _false_, **which allows you to know if the process is still in use or not.**

If you simply want to check if a process is alive or not, you can test it with a kill 0:

```perl
unless (kill 0, $PID) {
   warn "$PID doesn't exist\n";
}
```

Signal interception may seem more interesting than sending. For example if someone does a Ctrl+C on your program and you still have temporary files that exist, maybe you would like these files to be deleted anyway. Here's how to proceed:

```perl
my $temp = "/tmp/mysoft.$$";
mkdir $temp, 0700 or die "Unable to create $temp: $!";

sub flush {
   unlink glob "$temp/*";
   rmdir $temp;
}

sub my_manager_int {
   &flush;
   die "Interrupted, exit...\n";
}

$SIG{'QUIT'} = 'my_manager_int';
$SIG{'INT'} = 'my_manager_int';
$SIG{'HUP'} = 'my_manager_int';
$SIG{'TERM'} = 'my_manager_int';
...
```

The temporary files are created, the program runs etc...
And at the end of the program we flush:

```perl
&flush;
```

If a Ctrl+C occurs, the program jumps directly to the &flush section.

## Strings and Sorts

When we need to search for text, regex are very convenient, but can sometimes be too complicated. That's why there are strings.

### Locating a Substring with index

For example, here we're looking for "mon":

```perl
my $substance = "Hello world!";
my $location = index($substance, "mon");
```

The return value here is 9, because there are 9 elements before finding "mon". If no occurrence has been found, the return value will be -1.

```perl
my $substance = "Hello world!";
my $location1 = index($substance, "e"); # location1 = 7
my $location2 = index($substance, "e", $location1 + 1); # location2 = 13
my $location3 = index($substance, "e", $location2 + 1); # location3 = -1
```

We can also reverse the search with rindex:

```perl
my $last_slash = rindex("/etc/passwd", "/"); # the value = 4
```

And finally there's a last optional parameter that allows to give the maximum authorized return value:

```perl
my $fred = "Yabba dabba doo!";
my $location1 = rindex($fred, "abba"); # location1 = 7
my $location2 = rindex($fred, "abba", $location1 - 1); # location2 = 1
my $location3 = rindex($fred, "abba", $location2 - 1); # location3 = -1
```

### Manipulating Substrings with substr

substr takes 3 arguments:

- A string value
- An initial location based on 0
- The length of the substring

The return value contains the substring:

```perl
my $mineral = substr("Fred J. Flinstone", 8, 5); # gets "Flint"
my $rock = substr "Fred J. Flinstone", 13, 1000; # gets "stone"
my $pebble = substr "Fred J. Flinstone", 13; # gets "stone"
```

If we don't put a 3rd parameter, it will go to the end of the string, regardless of the length.

To invert the selection, we'll use negative numbers:

```perl
my $output = substr("This is a long line", -3, 2); # output = in
```

index and substr work very well together. For example, here we'll extract a substring starting with the letter l:

```perl
my $long = "a very long line";
my $just = substr($long, index($long, "l"));
```

Here's now how to make the whole a little more flexible:

```perl
my $chain = "Hello world!";
substr($chain, 0, 4) = "Good bye"; # will give "Good bye world!"
```

Here's an even shorter way:

```perl
substr($chain, -20) =~ s/fred/barney/g;
```

In reality, we never use this kind of code. But you might need it.

Use most often the index and substr functions to regex because they don't have the regex engine overload:

- They are never case insensitive
- They don't care about metacharacters
- They don't initialize any of the variables in memory

### Formatting Data with sprintf

```perl
my $date_tag = sprintf "%4d/%02d/%02d %2d:%02d:%02d", $yr, $mo, $da, $h, $m, $s;
```

Here, $date_tag receives something like "2008/12/07 03:00:30". The format string (the first argument of sprintf) places a 0 at the beginning of certain numbers, which we hadn't mentioned when we studied the printf formats. This 0 asks to add leading 0s as requested to give the number the required width. Without this 0, there would be spaces instead: "2008/12/ 7 3: 0:30".

### Using sprintf with Monetary Numbers

To indicate a sum of money in the form 2.50, not 2.5 - and especially not 2.49997! The format "%2f" allows to easily get this result:

```perl
my $money = sprintf "%.2f", 2.49997;
```

**The complete implications of rounding are numerous and subtle, but most often it is desirable to keep the numbers in memory with all possible precision, only rounding for display.**

```perl
sub big_sum {
   my $number = sprintf "%.2f" shift @_;
   # We add a comma at each pass through the loop that does nothing
   1 while $number =~ s/^(-?\d+)(\d\d\d)/$1,$2/;
   # We place the dollar sign in the right place
   $number =~ s/^(-?)/$1\$/;
   $number;
}
```

```perl
while ($number =~ s/^(-?\d+)(\d\d\d)/$1,$2/) {
    1;
}
```

Why didn't we simply use the /g modifier to perform a "global" search and replace and avoid the pain and confusion of the 1 while? Because we're working backwards from the decimal point, not advancing from the beginning of the string. The placement of commas in a number like this cannot be accomplished by a s///g substitution alone.

### Advanced Sorting

You can use cmp to create a more complex sort order, for example case insensitive:

```perl
sub no_senscase {
 "\L$a" cmp "\L$b";
}
```

```perl
my @numbers = sort { $a <=> $b } @some_numbers;
my decreasing = reverse sort { $a <=> $b } @some_numbers;
```

If for example we want to sort by file modification/creation order:

```perl
# Get logs files
my @logs_files = glob "$path/logs/*$market_name*";
my @sorted_logs = sort { -M $a <=> -M $b } @logs_files;
```

### Multiple Key Sorting

```perl
my %score = (
   "toto" => 195,
   "tata" => 205,
   "titi" => 30,
   "tutu" => 195,
);
```

We want to rank the players above by score, and if the score is identical to another, alphabetically:

```perl
my @winners = sort by_score_and_name keys %score;

sub by_score_and_name {
   $score{b} <=> $score{a} # descending numerical score
   or
   $a cmp $b               # by ASCIIbetical name
}
```

If the spaceship operator sees 2 different scores, that's the desired comparison. It returns -1 or 1 (true value) and the short-circuit or indicates that the rest of the expression should be skipped and the desired comparison is returned. (Remember that the or shortcut returns the last evaluated expression).
But if the spaceship operator sees 2 identical scores, it returns 0 (false value); the cmp operator takes over and returns an appropriate ranking value considering the keys as strings. If the scores are identical, the string comparison ends the competition.

There's no reason for your sorting subroutine to be limited to 2 levels of sorting. Below, the Bedrock library program ranks a list of customer ID numbers according to a 5-level sort order: each customer's unpaid penalties (calculated by an absent subroutine here, &penalties), the number of items currently consulted (from %items), their name (in the order of last name then first name, both from hashes), finally the customer ID number, in case the rest is identical:

```perl
@clients_IDs = sort {
   &penalties($b) <=> &penalties($a) or
   $articles{$b} <=> $articles{$a} or
   $nom_de_familles{$a} cmp $nom_de_famille{$a} or
   $prenom{$a} cmp $nom_de_famille{$b} or
   $a <=> $b
} @clients_IDs;
```

## Simple Databases

### Opening and Closing DBM Hashes

It's quite simple to understand here:

```perl
dbmopen(%DATA, "my_database", 0644) or die "Can't create my database: $!";
dbmclose(%DATA);
```

### Using a DBM Hash

The DBM hash looks like any other hash, but instead of being stored in memory, it's stored on disk. Therefore, when your program opens it again, the hash already contains data from the previous call (some beginner docs will tell you to no longer use DBM bases and replace them with "tie", however there's no need to use complex methods when we wish to do something simple).

```perl
while (my($key, $value) = each(%DATA)) {
   print "Value of $key is $value\n";
}
```

### Manipulating Data with pack and unpack

Pack is used to pack data. We gather 3 numbers of different sizes into a 7-byte string using the formats c, s and l (reminiscent of char, short and long). The first number is packed into one byte, the second into 2 and the 3rd into 4:

```perl
 my $buffer = pack("c s l", 31, 4159, 265359);
```

It's possible to improve visibility by placing spacings in a format string. For example, you can convert "ccccccc" to "c7". Obviously in case of unpack, it's simpler to use "c\*" in case of unpack.

### Random Access Databases of Fixed Length

There are several available formats in the pack documentation. So you have to choose the appropriate format. The open function has another mode that we haven't presented yet. Placing "+<" before the file name parameter string is equivalent to using "<" to open the existing file for reading, while additionally requesting permission to write to the file:

```perl
open(READ, "<toto");
open(READ_WRITE, "<+toto");
```

Similarly, conversely, you can write to a file and then read it:

```perl
open (WRITE, ">toto");
open (WRITE_READ, ">+toto");
```

To summarize:

- "+<": Allows to read an existing file, then write to it
- ">+": Allows to write to a file, then read it

The latter is usually used for draft files.

Once the file is opened, we must browse it, we'll use the seek function.

```perl
seek(TOTO, 55 * $n, 0);
```

- 1st parameter: file handle
- 2nd parameter: the displacement in bytes from the beginning of the file
- 3rd parameter: The starting parameter, or 0

```perl
my $temp; # Input buffer variable
my $read_number = read(Toto, $temp, 55);
```

- 1st parameter: file handle
- 2nd parameter: buffer variable to receive the read data
- 3rd parameter: number of bytes to read

We asked for 55 bytes because that corresponds to the size of our record. But if in the buffer, you are 5 bytes from the end of the 55, then you'll only have 5.

Here's a small example of all this with a time function to give the time:

```perl
print TOTO pack("a40 C I5 L",
   $name, $age,
   $info, $info1, $info2, $info3, $info4,
   time);
```

On some systems, it's necessary to use seek at each transition from a read to a write, even when the current position in the file is correct. It's therefore advised to use seek immediately before any read or write.

Here's a small alternative to what we just saw:

```perl
my $pack_format = "a40 C I15 L";
my $pack_lenght = lenght pack($pack_format, "data", 0, 1, 2, 3, 4, 5, 6);
```

### Variable Length (text) Databases

The most common way to update a text file by program is to write an entirely new file similar to the old one, making the necessary changes along the way. This technique gives pretty much the same result as updating the file itself, but offers additional advantages.

That's where Perl is once again great, it's possible to do substitution directly on the file without having to recreate a new file with the operator "<>":

```perl
#!/usr/bin/perl -w

use strict;

chomp(my $date= `date`);
@ARGV = glob "toto*.dat" or die "no files were found\n";
$^I = ".bak";

while (<>) {
   s/^Toto:/^Titi:/;
   print;
}
```

To provide the list of files to the diamond operator, we read them in a glob. The next line initializes $^I. By default this is initialized to undef, but when it's initialized to a certain string, it makes the diamond operator even more magical.
From then on, the while loop reads a line from the old file, updates it then writes it to the new file. This program is able to update hundreds of files in a few seconds on a classic machine.
Once the program is finished, Perl has automatically created a backup file of the original as if by magic.

### In-place Editing from the Command Line

Imagine that you need to correct hundreds of files with a spelling mistake. Here's the solution in a single line:

```bash
perl -p -i.bak -w -e 's/toto/tata/g' titi*.dat
```

## Some Advanced Perl Techniques

### Intercepting Errors with eval

```perl
eval { $toto = $titi / $tutu };
```

Now, even if $tutu equals 0, this line won't crash the program because eval is an expression and not a control structure like while or foreach. If an error remains, you'll find it in $@.

```perl
eval {
   ...
}
if ($@) {
   print "An error occurred ($@), continuing...\n";
}
```

It's possible to nest eval blocks within other eval blocks.

### Selecting List Items with grep

Let's choose the odd numbers from a long list of numbers. We don't need any new features:

```perl
my $odd_numbers;

foreach (1..1000) {
   push @odd_numbers, $_ if $_ % 2;
}
```

This code uses the modulo operator. If a number is even, this number modulo 2 gives 0 which is worth false. But an odd number will give 1; as that equals true, only odd numbers will be placed in the array.

There's nothing wrong with this code, except that it's a bit long to write and **slow to execute**, if one knows that perl offers the grep operator:

```perl
my @odd_numbers = grep { $_ % 2 } 1..1000;
```

Another example:

```perl
my @matching_lines = grep /\btoto\b/i, <FILE>;
```

### Transforming List Elements with map

```perl
my @data = (2, 3, 123, 8.343, 29.84);
my @formated_data = map { &big_addition($_) } @data;
```

The map operator resembles grep because it has the same type of arguments: a block using $_ and a list of elements to process. It operates the same way, by evaluating the block once for each element of the list, $_ becoming the alias of a different original element of the list each time.

### Hash Keys Without Quotes

Naturally not on any key; indeed, a hash key can be any string. But keys are often simple. If a hash key consists only of letters, digits and underscores, without starting with a digit, you can omit the quotes (called simple word):

```perl
my %score = (
   toto => 10,
   tata => 15,
   tutu => 20,
);
```

### More Powerful Regular Expressions

#### Non-greedy Quantifiers

Regular expressions can consume more or less CPU depending on how you perform the searches. The more elements found (in chronological order), the longer your program's execution time will be.

Here's an example to remove tags, but there's an error:

```perl
s/<BOLD>(.*)\/BOLD>//$1/g;
```

The asterisk (\*) being greedy, we won't fall on what interests us, here's the good solution:

```perl
s/<BOLD>(.*?)\/BOLD>//$1/g;
```

The ? allows to say that it may be the case or not, and not necessarily always (as the first example would do).

#### Multiline Text Recognition

The ^ and $ anchors allow to match beginnings and ends of lines, but if we want to also match internal newline characters, we'll need to use "/m":

```perl
print "Found 'toto' at the beginning of the line\n" if /^toto\b/im;
```

That makes anchors of beginning and end at each line, not of the global string.

### Slices

If for example, we want to define a stat on variables with a certain number as undef:

```perl
my (undef, undef, undef, undef, undef, undef, undef, undef, undef, $mtime) = stat $file
```

If we indicate a wrong number of undef, we'll unfortunately reach an atime or ctime, which is quite difficult to debug. **However there is a better method!** Perl knows how to index a list as if it was an array.

Here, as mtime is element 9 of the list returned by stat, we get it by an index:

```perl
my $mtime = (stat $file)[9];
```

Now, let's see how with the list, we can split info:

```perl
my $card_number = (split /:/)[1];
my $total = (split /:/)[5];
```

This method is good but not efficient enough. Here's how to do better:

```perl
my($card_number, $total) = (split /:/)[1, 5];
```

If now we want to retrieve some elements from a list (-1 representing the last element):

```perl
my($first, $last) = (sort @names)[0, -1];
```

Here's an example of an extraction of 5 elements on a list of 10:

```perl
my @names = qw{ zero one two three four five six seven eight nine ten };
my @numbers = { @names }[9, 0, 2, 1, 0];
print "Bedrock @numbers\n"; # Gives "Bedrock nine zero two one zero"
```

#### Array Slice

The example above even simplified would have given:

```perl
my @numbers = @names[9, 0, 2, 1, 0];
```

To enter information into an array easily when it's currently variables, here's the solution:

```perl
my $phone = "08 36 65 65 65";
my $address = "Lapland";
@elements[2, 3] = ($phone, $address);
```

#### Hash Slice

Here's another technique that works:

```perl
my @three_scores = ($score{"toto"}, $score{"tata"}, $score{"titi"});
```

But again, this is not the most optimized solution. Here it is:

```perl
my @three_scores = @score{ qw/ toto tata titi/ };
```

Why isn't there a % when we're talking about hashing?
It's the mark indicating a global hash; a hash slice (like any slice) is always a list, not a hash (just like a house fire is a fire and not a house). In Perl, the $ sign means a single element, the @ a list of elements and the % an entire hash.

### Creating a Daemon

A daemon allows running an app in the background. Here's how to do it in perl:

```perl
use POSIX qw(setsid);
defined(my $pid=fork) or die "cannot fork process:$!";
exit if $pid;
setsid;
umask 0;

# Enter loop to do work
while (1) {
   sleep 1;
}
```

### Creating Sockets

This part explains how to make a socket (client/server) work.

We'll use a fork for each connection to not be limited. Otherwise, we would have to wait for the first connection to close before a second one is handled.
Thanks to fork, a new child is created for each connection. It will handle it while its parent will wait for a new connection, ready to delegate it to another of its children.

The following code will be properly commented to give all necessary explanations.

Here's what the server will look like:

```perl
#!/usr/bin/perl -w
use strict;
use POSIX qw(setsid);
use IO::Socket;
use Socket;

my ($client, $remote_machine, $port, $addr_i, $real_ip, $client_name, $child_pid);
my $server = new IO::Socket::INET ( LocalPort => '7070', Type => SOCK_STREAM, Reuse => 1, Listen => 10);
die "Could not create socket: $!\n" unless $server;

# We create a loop in which the parent process will return once a child is created
REQUEST:
# This while means a connection has been requested by a client.
while ($client = $server->accept())
{
    # Thanks to these first lines, we can display a welcome message
    # $remote_machine = getpeername($client);
    # ($port, $addr_i) = unpack_sockaddr_in($remote_machine);
    # $real_ip = inet_ntoa($addr_i);
    # $client_name = gethostbyaddr($addr_i, AF_INET);
    # print "$client_name is connecting!\n";


    # A new connection being requested, the parent tries to duplicate itself to be able to accept another client connection if needed
    if($child_pid = fork)
    {
        # If successful, then the unused handle (the parent's client one) is closed and the parent restarts listening while the child continues.
        # Thus a new child is created for each connection, the parent only generates them.
        close($client);
        next REQUEST;
    }
    defined $child_pid or die "Fork impossible: $!\n";

    # The child closes its parent's unused handle
    close($server);

    # Automatic buffer flushing
    $| = 1;

    # Real code for inputs and outputs with the client:
    while(my $input = <$client>)
    {
        print STDOUT "option sent: $input";

        # Gives a default handle for future displays
        select($client);
        # display on the client machine since the handle is specified
        print $client "display\n";


	# If we don't want to do "select($client)", we can choose which output will be sent to the client machine with these lines
        # open(STDIN, "<<&$client");
        # open(STDOUT, ">&$client");
        # open(STDERR, ">&$client");


        print "finished\n"; #In the client still since default handle

	# We close the handle to end the connection.
        close($client);
        # To prevent the child from going back into the loop. That would generate an error given that the server handle has been closed by the child.
        exit;
    }
}
```

The client, will look like this (much simpler):

```perl
#!/usr/bin/perl -w
use strict;
use IO::Socket;
use Getopt::Long;

my ($response, $option);
my $socket = new IO::Socket::INET (PeerAddr => 'localhost',PeerPort => '7070',Proto => 'tcp', Type => SOCK_STREAM);
die "Could not create socket: $!\n" unless $socket;

if($ARGV[0])
{
    # We send in the socket the option we gave to the script
    print $socket "$ARGV[0]\n";
}

while($response = <$socket>)
{
    # We display all the responses sent by the Server
    print $response;
}
# We close the socket to end the connection
close($socket);
exit;
```

It's possible to code functions in both parts. The easiest thing to implement to make the server perform actions is to have it analyze what the client sends it.
According to its correspondence with one regex or another, it will launch a certain function.

## Going Further in Perl

First of all, know that if you need additional documentation, you can look at the following mans:

- perltoc (table of contents)
- perlfaq
- perldoc
- perlrun

For regular expressions, you can also find what you need here:

- perlre
- perltut (for tutorials)
- perlrequick

### Modules

Modules are tools made to save you time. They allow you not to reinvent the wheel each time you want new functionalities. Abuse them! I invite you to go to the [CPAN](https://www.cpan.org) site and take a little tour.

To install a module, it's very simple, as root launch the cpan command, then:

```bash
install My::Module::to_install
```

And everything will be done automatically :-)

If you encounter compilation problems, check the errors, but in many cases, you're missing C development libraries:

```bash
apt-get install libc6-dev
```

In the very rare case where there's no module to do what you want, you can develop one yourself in Perl or C (and submit it afterwards, think of the community!). Consult:

- perlmod
- perlmodlib

#### Listing Installed Modules

It can be very practical to list installed modules. For this we'll need this package:

```bash
aptitude install perl-doc
```

Then, all that's left is to launch this command:

```bash
perldoc perllocal
```

#### Knowing if a Module is Integrated by Default in Perl

To know if a module is integrated into the Perl Core, this command exists:

```bash
corelist <perl_module>
```

Example:

```bash
> corelist strict
strict was first released with perl 5
```

### Some Important Modules

#### Cwd

Equivalent to the pwd command, it allows to know the directory where we are (.):

```perl
 use Cwd;
 my $folder = cwd;
```

And if we want to know the perl file that is currently running (the one we launch):

```perl
use Cwd 'abs_path';

my $current_file=$0;
print "$current_file\n";
```

#### Fatal

If you're tired of writing "or die" for each invocation of open or chdir for example, fatal is made for you:

```perl
 use Fatal qw/ open chdir /;
 chdir '/home/toto';
```

Here no need to indicate "or die", we'll get thrown out if we couldn't change folders.

#### Sys::Hostname

To know the name of the machine:

```perl
 use Sys::Hostname;
 my $host = hostname;
 print "This machine is called $host\n";
```

#### Time::Local

Conversion from Epoch time to human readable time:

```perl
 use Time::Local;
 my $now = timelocal($sec, $min, $hr, $day, $month, $year
```

#### diagnostics

To have more information about errors in your code:

```perl
 use diagnostics
```

#### Big Numbers

To have a calculation with big numbers, use Math::BigFloat or Math::BigInt.

#### splice

It allows to add or delete elements in the middle of an array.

#### Security

Perl has a functionality that will know exactly what Perl uses in memory (in case there would be data corruption at this level). The anti-pollution control module is also called taint checking. Also see perlsec.

#### Debugging

See the B::Lint module.

#### Converting from Other Languages

- To convert from sed: man s2p
- To convert from awk: a2p

#### find

To convert the find command to perl use the find2perl command:

```perl
 find2perl /tmp -atime +14 -eval unlink >toto.pl
```

#### Socket

If we want to know if a port is listening or not:

```perl
use Socket;

my $proto = getprotobyname('tcp');
my $port=44556;

socket(SOCK, PF_INET, SOCK_STREAM, $proto);
if (!connect(SOCK, sockaddr_in($port, inet_aton('localhost')))) {
    print "Nobody's listening\n";
}
close SOCK;
```

For more info:
http://perldoc.perl.org/perlipc.html#Sockets%3a-Client%2fServer-Communication

#### Getopt::Long and Getopt::Std

This small module is one of my favorites because it allows to manage the entire @ARGV part without racking your brain. A small example:

```perl
 use Getopt::Long;

 # Vars
 my ($hostname,$base,$filter);
 my $scope='sub';

 help if !(defined(@ARGV));

 # Set options
 GetOptions( "help|h"    => \&help,
     "H=s"      => \$hostname,
     "b=s"	=> \$base,
     "f=s"	=> \$filter);
```

Here, you define variables to indicate the arguments. For example:

- s: is of type string (here help is defined by the letter h or help)
- i: is of type integer.

The advantage is that this module handles the variables no matter where they're located, no matter if there's 1 or 2 '-' before. In short, only happiness for the end user.

For more info: http://perldoc.perl.org/Getopt/Long.html

#### Term::ANSIColor

This is to use colors. Personally, I use it to display OK in green and Failed in red. A little example:

```perl
use Term::ANSIColor;

# Set OK color in green and FAILED in red
# 1st arg is message line and 2nd is ok or failed
sub print_color
{
    # Print message
    printf "%-60s", "$_[0]";
    # Print OK in green
    if ($_[1] =~ /ok|reussi/i)
    {
        print color 'bold green';
    }
    else
    {
        # Print failed in red
        print color 'bold red';
    }
    printf "%20s", "$_[1]\n";
    # Resetting colors
    print color 'reset';
}

print color 'bold red';
print "# WARNING: BE AWARE THAT YOU ARE ACTUALLY WORKING DURING PRODUCTION TIME #\n";
print color 'reset';

print_color("Checks passed successfully",'[ OK ]');
print_color("Environment check failed",'[ FAILED ]');
```

This should give you an idea of how colors work.

#### Term::ReadKey

Here's a simple but useful thing, if we don't want passwords to be displayed when someone types something on STDIN:

```perl
use Term::ReadKey;

print "Please enter password:\n";
ReadMode 2;
<STDIN>;
ReadMode 0;
```

### Creating a Module

Creating a module can be convenient to separate parts of your code. It is preferable, most of the time, to only do functions for modules and call them only from the main code.

To create a module, it's simple. Imagine I have 2 scripts:

- main.pl (my main script)
- my_module.pm (my module)

So the module must have an extension in ".pm". Then, it must contain this at the beginning:

```perl
#!/usr/bin/perl -w -I .
package my_module;
...
```

And for it to be valid, the end of the module must end like this:

```perl
...
1;
```

That's it :-)

### Creating a Binary

I don't really like this policy but it can be very interesting to create binaries from perl source code. To do this, just use the perlcc command:

```bash
perlcc -o binary source.pl
```

And there you go :-), it's very simple and your code is no longer disclosed.

### Creating an exe Under Windows

There are several commercial tools to make exes. For my part, I chose to use, once again, free and free. First, you need to install [activeperl](https://www.activestate.com/activeperl/) or [Strawberry](https://strawberryperl.com/) to be able to run Perl under Windows ([cygwin](https://www.cygwin.com/) might also do the trick).
Personally, I have a slight preference for Strawberry because it's very similar to Linux and it's completely free.

#### With Strawberry

Open the Perl Command Line, run cpan and install the following modules:

```bash
> cpan
install PAR
install Getopt::ArgvFile
install Module::ScanDeps
install Parse::Binary
install Win32::Exe
install PAR::Packer
```

#### With ActivePerl

Now, use Perl Package Manager to install the following packages as it will allow us to install all the necessary dependencies (in View, click on "All packages" to see all available packages):
![Ppm1](/images/ppm1.avif)

We'll also take advantage to install these packages:

- PAR
- MinGW
- Getopt-ArgvFile
- Module-ScanDeps
- Parse-Binary
- Win32-Exe
- PAR-Packer

If the PAR-Packer package is not available for your version, you'll have to get it from an external site [PAR-Packer](https://search.cpan.org/~smueller/PAR-Packer-1.002/lib/pp.pm). To install it, make sure you have version 5.10 of Perl, otherwise adapt with the right file, and run this command under Windows:

```bash
ppm install http://www.bribes.org/perl/ppm/PAR-Packer-5101.ppd
```

Or else, add the repository in PPM preferences:

- Name: A repository of Bioperl packages
- Location: http://bioperl.org/DIST

#### Generate the exe

Now we can generate .exe from a Perl script very easily:

```bash
C:\Documents and Settings\deimos\Bureau\>pp test.pl -o test.exe
Set up gcc environment - 3.4.5 (mingw-vista special r3)
Set up gcc environment - 3.4.5 (mingw-vista special r3)
Copyright (C) 2004 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.


Set up gcc environment - 3.4.5 (mingw-vista special r3)
Copyright (C) 2004 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
```

So here I compiled a test.pl file into test.exe. Simple, right? :-)

If you want to add an icon add your icon at the end like this:

```bash
pp test.pl -C -x -o test.exe --icon icon.ico
```

If when launching, you encounter a problem of dependency libraries, you must specify a path during compilation to tell it where to find them (--lib):

```bash
pp test.pl -C -x -o test.exe --icon icon.ico --lib=C:\strawberry\perl\lib:C:\strawberry\perl\site\lib:C:\strawberry\perl\vendor\lib
```

## Memos

Here are some small memos that I use quite often.

### clear

To do the equivalent of a clear (clear the screen), here's the solution in perl:

```perl
print "\033[2J";
```

If you want to clear the current line to for example make a countdown:

```perl
$| = 1;
for (1 .. 5) {
print "\rStep $_";
sleep 1;
}
print "\rFinished.\n";
```

### Display the Paths of Your Perl Libraries

To display all available paths for libraries, we'll use this command:

```bash
$ perl -e 'print map { $_ . "\n" } @INC;'
/etc/perl
/usr/local/lib/perl/5.10.0
/usr/local/share/perl/5.10.0
/usr/lib/perl5
/usr/share/perl5
/usr/lib/perl/5.10
/usr/share/perl/5.10
/usr/local/lib/site_perl
.
```

### Getting the PID of Your Program

You can get your Perl program's PID very simply:

```perl
print "$$\n";
```

## Resources

- [Official Perl Site](https://www.perl.org/)
- [Perl Modules and Documentation](https://www.cpan.org)
- [To validate your code and claim to have "Clean" code](https://perlcritic.com)
- [Les Mongueurs de Perl](https://www.mongueurs.net/)
- [Perlport](https://perldoc.perl.org/perlport.html)
- [O'Reilly's Introduction to Perl](https://books.google.fr/books?id=NuKGSvzneeUC&pg=PP1&ots=cxREfzjaS7&dq=introduction+perl&sig=lRxH8nw5lWxPOyFA-A5_0aaxMgo) (I highly recommend it)
- [Another documentation on Introduction to Perl](/pdf/intro_perl.pdf)
- [New Features of Perl 5.10 Part 1](/pdf/perl_unixgarden_1.pdf)
- [New Features of Perl 5.10 Part 2](/pdf/perl_unixgarden_2.pdf)
- [Documentation to Create a Perl Module](/pdf/cpan_module.pdf)
