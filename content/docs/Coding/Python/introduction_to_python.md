---
weight: 999
url: "/Introduction_au_Python/"
title: "Introduction to Python"
description: "A comprehensive introduction to Python programming language, covering syntax, data types, structures, functions, modules, and useful libraries."
categories: ["Linux"]
date: "2012-06-06T12:45:00+02:00"
lastmod: "2012-06-06T12:45:00+02:00"
tags:
  [
    "Linux",
    "10.2 Les Tuples",
    "13.1 input",
    "12.1.1 dump",
    "navigation",
    "13.2 raw_input",
    "17 String",
    "18.4.1 Groups",
    "20 argparse",
    "4.1 sep",
  ]
toc: true
---

![Python](/images/python-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|---|---|
| **Software version** | 2.7 / 3.2 |
| **Website** | [Python Website](https://www.python.org/) |
| **Last Update** | 06/06/2012 |
{{< /table >}}

## Introduction

[Python](https://fr.wikipedia.org/wiki/Python_%28langage%29) is an interpreted, multi-paradigm programming language. It supports structured imperative programming and object-oriented programming. It features strong dynamic typing, automatic memory management through garbage collection, and an exception handling system; thus it is similar to Perl, Ruby, Scheme, Smalltalk, and Tcl.

The Python language is under a free license similar to the BSD license and runs on most computing platforms, from supercomputers to mainframes, from Windows to Unix including Linux and MacOS, with Java or even .NET. It is designed to optimize programmer productivity by providing high-level tools and a simple syntax. It is also appreciated by educators who find in it a language where syntax, clearly separated from low-level mechanisms, allows for an easier introduction to basic programming concepts.

Among all the programming languages currently available, Python is one of the easiest to learn. Python was created in the late '80s and has matured enormously since then. It comes pre-installed in most Linux distributions and is often one of the most overlooked when choosing a language to learn. We'll confront command-line programming and play with GUI (Graphical User Interface) programming. Let's dive in by creating a simple application.

In this documentation, we'll see how to write Python code and we'll use the interpreter (via the python command). The lines below corresponding to the interpreter will be visible via elements of this type: '>>>' or '...'

When you write a file that should understand Python, it should contain this at the beginning (the Shebang and encoding):

```python
#!/usr/bin/env python
# -*- coding:utf-8 -*-
```

The encoding line is optional but necessary if you use accents.

## Syntax

In Python, we must indent our lines to make them readable and especially for them to work. You need to indent using tabs or spaces. **Be careful not to mix the two, Python doesn't like that at all: your code won't work and you won't get an explicit error message.**

You can use comments at the end of lines if you wish. Here's an example:

```python
bloc1
   blabla # comment1
Suite du bloc1
   blibli # comment2
```

The end of a block is automatically done with a line break.

## help

Know that at any time you have the possibility to ask for help thanks to the help command. For example for help with the input method:

```python
>>> help(input)
Help on built-in function input in module __builtin__:

input(...)
    input([prompt]) -> value

    Equivalent to eval(raw_input(prompt)).
```

You can also access this help from the shell using the pydoc command:

```bash
pydoc input
```

## Displaying text

Let's see the first most basic command, displaying text:

```python
>>> print 'Deimos Fr'
Deimos Fr
```

Then here's how to concatenate 2 elements:

```python
>>> print 'Deimos ' + 'Fr'
Deimos Fr
>>> print 'Deimos', 'Fr'
Deimos Fr
```

When using the comma, strings are automatically separated by a space character, whereas when using concatenation, you have to manually manage this issue.
If you concatenate non-string variables, you'll need to convert them before you can display them:

```python
>>> c = 3
>>> print 'Value: ' + c Traceback (most recent call last):
File "<stdin>", line 1, in <module>
TypeError: cannot concatenate 'str' and 'int' objects
>>> print 'Value: ' + str(c)
Value: 3
>>> print 'Value:', c
Value: 3
```

**If you don't want to have an automatic line return (the equivalent of "\n")**, you simply **put a comma at the end of your line**:

```python
>>> print 'No \n',
```

In Python 3.2, here's how to print:

```python
>>> print('Deimos', 'Fr')
Deimos Fr
```

### sep

In Python 3.2, with sep we can separate strings with characters:

```python
>>> print('Deimos', 'Fr', sep='--')
Deimos--Fr
>>> print('a', 'b', 'c', 'd', sep=',')
a,b,c,d
```

### end

In Python 3.2, with end we can add a character at the end of a string:

```python
>>> print('Deimos', 'Fr', sep='--', end='!')
Deimos--Fr!
```

## Data types

There are several data types in Python:

- Integers: allow representing integers on 32 bits (from -2,147,483,648 to 2,147,483,647).

```python
>>> type(2)
<type 'int'>
>>> type(2147483647)
<type 'int'>
```

- Long integers: any integer that cannot be represented on 32 bits. An integer can be forced to a long integer by following it with the letter L or l.

```python
>>> type(2147483648)
<type 'long'>
>>> type(2L)
<type 'long'>
```

- Floats: The symbol used to determine the decimal part is the dot. You can also use the letter E or e to indicate an exponent in base 10.

```python
>>> type(6.15)
<class 'float'>
>>> 3e1
30.0
>>> type(3e1)
<class 'float'>
```

## Calculations

When doing divisions, you need to be careful, results depend on the context! See for yourself:

```python
>>> 7/6
1
>>> 7.0/6.0
1.1666666666666667
>>> 7//6
1
```

Be careful with operations! Here's an example:

```python
>>> ((0.7+0.1)*10)
7.999999999999999
```

You see, there's a rounding issue!

It's possible to perform an operation using the old value of a variable with +=, -=, \*=, and /= (variable += value is equivalent to variable = variable + value). For example:

```python
>>> a = 1
>>> a
1
>>> a += 2
>>> a
3
```

And you can also perform multiple variable assignments with different values by separating variable names and values with commas:

```python
>>> a, b = 1, 2
>>> a
1
>>> b
2
```

To raise a number to any power, use the \*\* operator:

```python
>>> 2*3
8
```

## String manipulation

Concatenation (assembling several strings to produce only one) is done using the + operator. The \* operator allows repeating a string. Examples:

```python
>>> a = "Hello"
>>> b = 'World !'
>>> a+b
'Hello World !'
>>> 3*a
'Hello Hello Hello '
```

Access to a particular character in the string is done by indicating its position in brackets (the first character is at position 0):

```python
>>> a[0]
'H'
```

## Booleans

Boolean values are noted as True or False. In a test context, values 0, _empty string_, and None are considered as False. Here are the comparison operators:

- ==: for equality
- !=: for difference
- <: less than
- > : greater than
- <=: less than or equal
- > =: greater than or equal

```python
>>> 1 == 2
False
>>> 1 <= 2
True
```

To combine tests, we use boolean operators:

- and
- or
- not

```python
>>> 0 and True
0
>>> 1 and True
True
>>> not (0 or True)
False
```

## Structures

### if

Here's how an if statement is constructed:

```python
if condition:
# Process block 1 # ...
else:
# Process block 2 # ...
```

### for

Here's a for loop with a list:

```python
>>> for e in ['tata', 'toto', 'titi']:
... print e
...
tata
toto
titi
```

We'll use a list here with two new parameters: continue and break:

```python
>>> range(0, 10)
[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
>>> for i in range(0,10):
...   if i==2:
...      continue
...   if i==4:
...      break
...   print i
...
0
2
3
```

- continue: allows moving to the next element
- break: stops the loop it's in

### While

Here's an example of a while loop:

```python
>>> i = 0
>>> while i<5:
...   i += 1
...   if i==2:
...      continue
...   if i==4:
...      break
...   print i
1
3
```

- continue: allows moving to the next element
- break: stops the loop it's in

## Lists

### String characters

As with a list, we can access its elements (the characters) by specifying an index:

```python
>>> chaine = 'Pierre Mavro'
>>> chaine[0]
'P'
```

However, it's impossible to modify a string! We say it's an immutable list:

```python
>>> chaine[1] = 'a'
Traceback (most recent call last):
   File "<stdin>", line 1, in <module>
TypeError: 'str' object does not support item assignment
```

In fact, here we're making a new assignment: the variable chaine is overwritten (erased then replaced) with the value chaine + '!'. This is not a modification in the proper sense.

In Python 2.7, we'll use formatting expressions from C:

- %d: for an integer
- %f: for a float
- %s: for a string
- %.3f: to force display of 3 digits after the decimal point
- %02d: to force display of an integer on two characters

```python
>>> c = 1
>>> d = 2
>>> '%d * %d = %d' % (c, d, c*d)
'1 * 2 = 2'
```

In Python 3.2:

```python
>>> '{} * {} = {}'.format(c, d, c*d)
'2 * 3 = 6'
>>> chaine = 'An integer on two digits: {:02d}'
>>> chaine.format(6)
'An integer on two digits: 06'
```

### Tuples

Tuples are also immutable lists (string characters are therefore actually special tuples):

```python
>>> t = (1, 2, 3)
>>> t[0]
1
>>> l = [1, 2, 3]
>>> t = tuple(l)
>>> t
(1, 2, 3)
>>> l
[1, 2, 3]
>>> t = tuple('1234')
>>> t
('1', '2', '3', '4')
```

As a reminder, a tuple containing only one element is noted in parentheses, but with a comma before the last parenthesis:

```python
>>> t = (1,)
>>> t
(1,)
```

### Using lists

#### Modifying lists

To modify lists:

```python
>>> l = [ 1, 2, 3, 4]
>>> l
[1, 2, 3, 4]
>>> l[0] = 0
>>> l
[0, 2, 3, 4]
```

#### Adding an element to the end of a list

The append() method adds an element to the end of a list:

```python
>>> l = [ 1, 2, 3, 4]
>>> l.append(5)
>>> l
[0, 2, 3, 4, 5]
```

#### Adding an element at a defined location in a list

The insert() method allows specifying the index where to insert an element:

```python
>>> l = [ 0, 2, 3, 4]
>>> l.insert(1, 1)
>>> l
[0, 1, 2, 3, 4, 5]
```

#### Removing an element at a defined location in a list

To delete an element, you can use the remove()/del() method which removes the first occurrence of the value passed as a parameter:

```python
>>> l = [ 1, 2, 3, 1 ]
>>> l.remove(1)
>>> l
[2, 3, 1]
>>> l.remove(1)
>>> l
[2, 3]
>>> l.remove(1)
Traceback (most recent call last):
File "<stdin>", line 1, in <module> ValueError: list.remove(x): x not in list
>>> l = [ 1, 2, 3, 1 ]
```

#### Concatenating 2 lists

The extend() method allows concatenating two lists without reassignment. This operation is achievable with reassignment using the + operator:

```python
>>> l1 = [ 1, 2, 3 ]
>>> l2 = [ 4, 5, 6 ]
>>> l1.extend(l2)
>>> l1
[1, 2, 3, 4, 5, 6]
>>> l1 = [ 1, 2, 3 ]
>>> l1 = l1 + l2
>>> l1
[1, 2, 3, 4, 5, 6]
```

#### Checking if an element belongs to a list

If the tested element is in the list, the returned value will be True and otherwise it will be False:

```python
>>> guys = [ 'tata', 'toto', 'titi' ]
>>> 'tata' in guys
True
>>> 'tutu' in guys
False
```

#### Getting the size of a list

If you want to know the number of elements present in a list:

```python
>>> l = [ 1, 2, 3 ]
>>> len(l)
3
```

#### List comprehension

Python implements a mechanism called "list comprehension", allowing use of a function that will be applied to each element of a list:

```python
>>> l = [ 0, 1, 2, 3, 4, 5 ]
>>> carre = [ x**2 for x in l ]
>>> carre
[0, 1, 4, 9, 16, 25]
```

Thanks to the instruction for x in l, we retrieve each element contained in the list l and place them in the variable x. We then calculate all the square values (x\*\*2) in the list, which produces a new list. This mechanism can produce more complex results. It can also be used with dictionaries.

### Slicing

Slicing is a method applicable to all list-type objects (except dictionaries). It's a "slicing into pieces" of list elements to retrieve certain objects. This is translated in this form:

```
d[start:end:step]
```

Here's an example:

```python
>>> c = 'Deimos Fr'
>>> c[0:6]
'Deimos'
>>> c[7:8]
'Fr'
```

- No step indication was given, the default value is then used, that is, 1.
- If the start value is absent, the default value used will be 0
- If the end value is omitted, the default value used will be the length of the string + 1

#### Reading a string

To retrieve the entire string:

```python
>>> c[:]
'Deimos Fr'
```

To invert the direction:

```python
>>> c[::-1]
'rF somieD'
```

#### Accessing the last element of the list

To access the last element of the list:

```python
>>> c[-2:]
'Fr'
```

If you give the interval [-2:-1], **you will not have the last letter**, and if you give [-2:0] you will get nothing since it's impossible to go from -2 to 0 with a step of 1 (the last letter being -1).

#### On Tuples

Applied to lists and tuples, slicing reacts in the same way. The difference is that we no longer manipulate only characters but any type of data:

```python
>>> guys = [ 'toto', 'tata', 'titi', 'tutu' ]
>>> guys[1:3]
['tata', 'titi']
>>> guys[-1::-2]
['tata', 'tata']
>>> t_guys = ('toto', 'tata', 'titi', 'tutu')
>>> t_guys[::-1]
('tutu', 'titi', 'tata', 'toto')
```

#### Deleting elements from a list

Here's slicing on a list to remove elements:

```python
>>> guys = [ 'toto', 'tata', 'titi', 'tutu' ]
>>> guys[1:3] = []
('tutu', 'toto')
```

#### Copying lists

Look at this if you want to make copies of a list:

```python
>>> liste_a = [ 1, 2, 3 ]
>>> liste_b = liste_a
>>> liste_a
[1, 2, 3]
>>> liste_b
[1, 2, 3]
>>> liste_b[0] = 0
>>> liste_b
[0, 2, 3]
>>> liste_a
[0, 2, 3]
```

You'll notice that **both lists are affected because in fact, this acts as an alias!** To make a copy, there are two solutions! Here's the first, you need to use [:]:

```python
>>> liste_a = [ 1, 2, 3 ]
>>> liste_b = liste_a[:]
>>> liste_a
[1, 2, 3]
>>> liste_b
[1, 2, 3]
>>> liste_b[0] = 0
>>> liste_b
[0, 2, 3]
>>> liste_a
[1, 2, 3]
```

## Dictionaries

### Deleting keys

Dictionaries are key => values lists. Like for lists, the del command allows removing an element and the word key in allows checking the existence of a key:

```python
>>> d = { 'key_1': 1, 'key_2': 2, 'key_3': 3 }
>>> del d['key_2']
>>> d
{'key_3': 3, 'key_1': 1}
>>> 'key_1' in d True
>>> 'key_2' in d False
```

If the lists are not in the desired order, it's simply because dictionaries are not ordered!

### Getting keys and values

In Python 2.7, to read the keys and values of a list:

```python
>>> guys = { 'tata': 1, 'toto': 2, 'titi': 3 }
>>> guys.keys()
['toto', 'titi', 'tata']
>>> guys.values()
[2, 3, 1]
>>> guys.items()
[('toto', 2), ('titi', 3), ('tata', 1)]
```

In Python 3.2, it's a little different:

```python
>>> guys = { 'tata': 1, 'toto': 2, 'titi': 3 }
>>> guys.keys()
dict_keys(['toto', 'titi', 'tata'])
>>> guys.values()
dict_values([2, 3, 1])
>>> guys.items()
dict_items([('toto', 2), ('titi', 3), ('tata', 1)])
```

Dictionaries don't consume as much memory as storing a list, only a pointer to the current element is stored.

### Adding an error if an element is not found

There's a solution to return a substitute value if no value is found during a search in a list thanks to the get() method:

```python
>>> guys = { 'tata': 1, 'toto': 2, 'titi': 3 }
>>> guys.get('toto', 'not found')
2
>>> guys.get('salades', 'not found')
'not found'
```

### Concatenating dictionaries

You can, just like lists, concatenate dictionaries with the update() method:

```python
>>> guys = { 'tata': 1, 'toto': 2 }
>>> guys_2 = { 'titi': 3, 'tutu': 4 }
>>> guys.update(guys_2)
>>> guys
{'tata': 1, 'toto': 2, 'titi': 3, 'tutu': 4}
```

### Copying a dictionary

To copy an already existing dictionary, use the copy() method:

```python
>>> dico_a = {'key_1': 1, 'key_2': 2, 'key_3': 3 }
>>> dico_b = dico_a.copy()
>>> dico_a
{'key_2': 2, 'key_3': 3, 'key_1': 1}
>>> dico_b
{'key_2': 2, 'key_3': 3, 'key_1': 1}
>>> dico_b['key_1'] = 0
>>> dico_b
{'key_2': 2, 'key_3': 3, 'key_1': 0}
>>> dico_a
{'key_2': 2, 'key_3': 3, 'key_1': 1}
```

## Complex Dictionaries and Lists

### pickle

There is a solution to simply save and restore complex structures (multidimensional lists/dictionaries for example) through serialization, called Pickle.
**Using pickle remains very dangerous**: only use files for which data can be verified, because loading a corrupted file can lead to execution of malicious code!

```python
>>> try:
...    import cPickle as pickle
... except:
...    import pickle
...
```

#### dump

To write data, we'll open a file in binary mode (wb). We'll then use the dump() function to write a variable by specifying the file object as a parameter:

```python
>>> file = open('data', 'wb')
>>> pickle.dump(data, file)
>>> file.close()
```

To read data from a file (rb), we'll use the load() function by passing the file object as a parameter:

```python
>>> file = open('data', 'rb')
>>> var = pickle.load(file)
>>> var
{'tata': 1, 'toto': 2, 'titi': 3, 'tutu': 4, 'tete': [a, b, c, d]}
file.close()
```

#### dumps

If you want to see the serialization of a variable, you can use the dumps() function, which performs the same task as dump(), but into a string rather than a file:

```python
>>> pickle.dumps(data)
"(dp1\nS'tata'\np2..."
```

### ConfigParser

This function allows managing ini-type files. It's very convenient for managing configuration files. For those who don't see what it looks like:

```ini
[section1]
   option_1 = value_1
   option_2 = value_2
[section2]
   option_3 = value_3
   option_4 = value_4
```

We can access the value of an option in a particular section:

```python
>>> from ConfigParser import SafeConfigParser
>>> parser = SafeConfigParser()
>>> parser.read('file.ini')
['file.ini']
>>> print parser.get('section1', 'option_1')
value_1
```

To get everything from a section:

```python
>>> from ConfigParser import SafeConfigParser
>>> parser = SafeConfigParser()
>>> parser.read('file.ini')
['file.ini']
for name in parser.options('section1'):
... print 'Option: ' + name
... print ' ' + parser.get('section1', name)
...
Option: option_1
value_1
Option: option_2
value_2
```

## Stdin

STDIN allows displaying a message on the screen so a user can type characters on the keyboard and press the enter key to validate. In Python 2.7 there are 2 solutions:

- input(): retrieves an integer or a float
- raw_input(): retrieves a string

### input

With the input() function, entering a string will result in an attempt to interpret it as a variable and thus cause an error:

```python
>>> c = input('Give an integer: ')
Give an integer: 1
>>> c
1
>>> c = input('Give an integer: ')
Give an integer: hello
Traceback (most recent call last):
File "<stdin>", line 1, in <module>
File "<string>", line 1, in <module> NameError: name 'hello' is not defined
>>> a = 1
>>> c = input('Give an integer: ') Give an integer: a
>>> c
1
```

In Python 3.2, there is only one input() function.

The main conversion functions are:

- int() to convert to integer
- float() to convert to float
- complex() to convert to complex
- str() to convert to string (useless when using with input()).

### raw_input

raw_input() is simpler since any input will be considered a string:

```python
>>> c = raw_input('Give a string:')
Give a string:
hello
>>> c
'hello'
>>> c = raw_input('Give a string:')
Give a string:
2
>>> c
'2'
```

## Handles

### Moving to a directory

It's possible to move to a directory via the chdir() function:

```python
from os import chdir
chdir('/etc')
```

### Reading a file

It's possible to open files for reading(r), writing(w), and appending(a):

```python
file = open('file.txt', 'mode')
file.close()
```

Here for example, we'll open a file named fichier.txt and we need to enter the mode (r/w/a).

#### Read

We can read the entire file to work on its content afterward (be careful about the file size which will be stored in memory):

```python
>>> file_content = file.read()
>>> file_content
'1st line\n2nd line \n3rd Line\n'
```

You'll note that line breaks are not interpreted! To interpret them, we'll need to use print:

```python
>>> print file_content
1st line
2nd line
3rd Line
```

It's possible to read x characters from the reading location (the beginning by default or from another location if you've already started reading the file):

```python
>>> file = open('file.txt', 'r')
>>> file_content = file.read(8)
>>> print file_content
1st line
>>> file_content = file.read(8)
>>> print file_content file_content
2nd line
```

We read the first 8 characters, then the next 8.

#### Readline

readlines() is identical to the read() function, except that the data will be sent in a list where each element will contain a line (with a \n at the end of each line):

```python
>>> file_content = file.readlines()
>>> file_content
['1st line\n', '2nd line\n', '3rd Line\n']
```

### Writing to a file

To write to a file, it's very simple, we'll call the write function:

```python
file = open('file.txt', 'w')
file.write('Write this text')
file.close()
```

### Strip

You may know [chomp in Perl]({{< ref "docs/Coding/Perl/introduction_to_perl.md" >}}#chomp), the strip() function allows doing the same thing, that is removing invisible characters (spaces, tabs, line breaks) from the beginning and end of a string. This function is part of the string module and has two variants:

- lstrip(): only removes characters at the beginning of the string ('l' for left)
- rstrip(): only removes characters at the end of the string ('r' for right)

```python
>>> from string import strip, lstrip, rstrip
>>> line = ' Deimos Fr\n '
>>> strip(line)
'Deimos Fr'
>>> lstrip(line)
'Deimos Fr\n '
>>> rstrip(line)
' Deimos Fr'
```

## Functions

Here's how we define a function (def) and we call it with its name and parentheses at the end:

```python
>>> def fonction:
...   print "Hello World"
>>> fonction()
```

### Function documentation

Documentation related to a function is written using triple quotes:

```python
>>> def fonction(x):
...   """Here I write my documentation
...   on the function present here.
...   And I finish like this"""
...   return x**2
...
>>> help(fonction)
>>> Help on function fonction in module __main__:
fonction()
Here I write my documentation
on the function present here.
And I finish like this
```

You saw then, we can ask for help on a function directly from Python.

### Function parameters

In Python, all parameters passed make calls to their memory address! However, since some types are not modifiable, it will then seem like they are passed by value. **The non-modifiable types are simple types (integers, floats, complexes, etc.), strings, and tuples**:

```python
>>> def plusone(a):
...    a = a + 1
...
>>> value = 3
>>> value
3
>>> plusone(value)
>>> value
3
```

With a modifiable parameter, such as a list, the modifications will be visible when exiting the function.

Here we ask it to take the function argument (x) and raise it to the power. The function will return the result thanks to the return. If a function doesn't have a return specified, the return value will be None.

```python
>>> def fonction(x):
...    return x**2
...
>>> fonction(3)
9
```

If parameters are not specified when calling the function, the default values will be used. **The only requirement is that parameters having a default value must be specified at the end of the parameter list. Here's an example with its calls:**

```python
>>> def fonction(a, b, c=0, d=1):
...    print a, b, c, d
...
>>> fonction(1, 2)
1 2 0 1
>>> fonction(1, 2, 3)
1 2 3 1
```

You've understood, we can therefore declare default values and override them on demand by passing them as arguments. The order of a function's parameters is not fixed if we name them:

```python
>>> def fonction_ordre(a, b):
...    return a
...
>>> fonction_ordre(b=1, a=2)
2
```

### The number of arguments of a function

In Python, you don't have to define beforehand the number of arguments that will be used. For this, we use \*args and \*\*kwargs:

- **\*args**: tuple containing the list of parameters passed by the user
- \***\*kwargs**: dictionary of parameters

The difference between the two syntaxes (\* or \*\*) is the type of the return variable:

```python
>>> def exemple_args(*args):
...    print args
>>> exemple_args(1, 2, "hello")
(1, 2, 'hello')
>>> def exemple_kwargs(**kwargs):
...   print kwargs
>>> exemple_kwargs(a=1, b=2, text="hello")
{'a': 1, 'text': 'hello', 'b': 2}
```

## Modules

Modules are very useful, as they correspond to ready-made functions, saving a lot of time when we use them. Imports can be placed anywhere in the code and can be integrated into conditional loops. But it's much easier to find them at the top of the file.

To load a module with all its functions:

```python
from module_name1 import *
import module_name2
```

Here we ask Python to load all functions (_) of the module_name1 module into memory.
If we don't include (_), we'll need to prefix the module name before calling its function:

```python
module_name2.function1()
```

If you wish to import only a few functions from a module (takes less memory space and allows using only what we want):

```python
from module_name import function1, function2
```

_Note: In case of identical function names, the last import takes precedence over those before!_

If you're handling modules with a very long name, you can define an alias with the keyword 'as'. For example, a module named ModuleWithAVeryLongName containing the func() function. With each call to func() you're not going to write: ModuleWithAVeryLongName.func()... So you need to use as:

```python
import ModuleWithAVeryLongName as myModule
myModule.func()
```

If during a module import, it contains instructions for immediate execution (which are not in functions), these will be executed at the time of import. To differentiate the behavior of the interpreter during direct execution of a module or during its loading, there is a specific test to add which allows determining a sort of main program in the module:

```python
>>> def function():
... print 'Deimos Fr'
...
>>> if __name__ == '__main__':
...    print 'Test function function()'
>>> function()
```

To define the body of the main program, insert this line:

```python
if __name__ == '__main__':
```

### Module Paths

By default, Python starts by searching:

1. In the current directory
2. In the directory(ies) specified by the PYTHONPATH environment variable (if defined)
3. In the Python library directory: /usr/lib/python<version>

**So be very careful when naming your files: if they have the name of an existing Python module, they will be imported instead of the latter since the import first looks for modules in the current directory**.

Here's an example of the PYTHONPATH (environment variables):

```bash
export PYTHONPATH=$PYTHONPATH:~/.python_libs
```

To indicate that A is a "module directory", we'll create an additional file in the folder containing the module: '\_ _init_ \_.py'. This file may contain nothing or contain functions. Modules are also called packages.

## String

string is installed by default with Python. It provides many methods to search for text and replace strings. It helps avoid using regex which can in some cases be very CPU intensive if poorly written.

### split

The split() function allows cutting a string passed as a parameter following one or more separator characters and returns a list of the cut strings. If nothing is specified, the space character will be used:

```python
>>> import string
>>> s = 'deimos:x:1000:1000::/home/deimos:/usr/bin/zsh'
>>> string.split(s, ':')
['deimos', 'x', '1000', '1000', '', '/home/deimos', '/usr/bin/zsh']
>>> p = 'Bloc Notes Info'
>>> string.split(p)
['Bloc', 'Notes', 'Info']
```

### Join

Join is the opposite of split and allows joining several elements together via one or more separator characters:

```python
>>> l = [ 'a', 'b', 'c', 'd' ]
>>> ' => '.join(l)
'a => b => c => d'
```

### Lower

Lower allows converting a string to lowercase:

```python
>>> s = 'Bloc Notes Info'
>>> s.lower()
'bloc notes info'
>>> s
'Bloc Notes Info'
```

**Attention: We talk about conversion but these methods do not modify the original string but return a new string!!!**

### Upper

Upper allows converting a string to uppercase:

```python
>>> s = 'Bloc Notes Info'
>>> s.upper()
'BLOC NOTES INFO'
>>> s
'Bloc Notes Info'
```

**Attention: We talk about conversion but these methods do not modify the original string but return a new string!!!**

### Capitalize

This function allows capitalizing only the first letter of a string:

```python
>>> s = 'bloc Notes Info'
>>> string.capitalize(s)
'Bloc notes info'
```

### Capwords

Capwords allows capitalizing the beginning of each word:

```python
>>> s = 'bloc Notes Info'
>>> string.capwords(s)
'Bloc Notes Info'
```

### Count

The count() function allows counting the number of occurrences of a substring in a string. The first parameter is the string in which to perform the search and the second parameter is the substring:

```python
>>> p = 'one guy, two guys, three guys'
>>> string.count(p, 'guy')
3
>>> string.count(p, 'guys')
2
```

### Find

The find() function allows finding the index of the first occurrence of a substring. The parameters are the same as for the count() function:

```python
>>> s = 'one guy, two guys, three guys'
>>> string.find(s, 'guy')
3
>>> string.find(s, 'four')
-1
```

In case of failure, find() returns the value -1 (0 corresponds to the index of the first character).

### Replace

The replace() function allows, as its name suggests, replacing a substring with another inside a string of characters. The parameters are, in order:

1. The string to modify
2. The substring to replace
3. The replacement substring
4. The maximum number of occurrences to replace (if not specified, all occurrences will be replaced), this is optional

```python
>>> s = 'one guy, two guys, three guys'
>>> string.replace(s, 'guy', 'girl')
'one girl, two girls, three girls'
>>> string.replace(s, 'guy', 'girl', 1)
'one girl, two guys, three guys'
```

### Maketrans

maketrans() creates the translation table and thanks to the translate() function, applies the translations to a string of characters:

```python
>>> hack = string.maketrans('aei', '43!')
>>> s = 'Ecriture de Hacker'
>>> s.translate(hack)
'3cr!tur3 d3 h4ck3r'
```

When creating the table with maketrans(), each character in the first position is transformed with characters in the second position. So e's are replaced with 3's, a with 4, and i with !.

### Substitute and Safe_substitute

There's a simpler approach than maketrans(). We can write this in a more readable and more flexible way:

```python
>>> s = string.Template("""
... Show one var: $var
... In another one: $var2
... In a text: ${var}iable
... Vars names: $$var et $$var2
... """)
>>> values = { 'var': 'variable', 'var2': '123456' }
>>> s.substitute(values)
'\nShow one var: variable\nIn another one: 123456\nIn a text: variableiable\nVars names: $var et $var2\n'
```

- $$ allows displaying the character $
- ${var_name} allows isolating a variable included in a word

You should use the safe_substitute() function if you don't want Python to cause an error in case of non-substitution (because absent from the dictionary):

```python
[...]
>>> s.safe_substitute(values)
[...]
```

## Regex

A regex allows for example finding elements within a line/phrase that could correspond to certain elements but for which we don't always have certainty.

### Search

We'll use the re module which offers the search() function:

```python
>>> pattern = '\d{4}'
>>> date = 'Year: 1982'
>>> match = re.search(pattern, date)
>>> print 'Pattern', match.re.pattern, 'is located:', match.start(), '-', match.end()
Pattern \d{4} is located: 8 - 11
```

**If the pattern is not found, the search() function returns the value None.**

If you frequently use the same search pattern, it's better to compile it to have more efficient code:

```python
>>> pattern = re.compile('\d{4}')
>>> match = pattern.search(pattern)
```

When you use complex patterns, it's recommended to comment them and therefore write a "verbose" regular expression. In this mode, multiple spaces and comments need to be ignored:

```python
>>> pattern="""
... ^                       # Start of line
... Deimos\s                # First name
... Fr                      # Last name
... $                       # End of line
... """
>>> match = re.search(pattern, 'Deimos Fr', re.VERBOSE)
```

**Attention: the search() function only allows finding the first substring matching the searched pattern!**

### List of regex

Here are the most common regex:

{{< table "table-hover table-striped" >}}
| Operator | Description | Example |
|----------|-------------|---------|
| ^ | **Beginning** of line | **^Deimos** for 'Deimos Fr!!!' |
| $ | **End** of line | **!$** for 'Deimos Fr!!!' |
| . | Any character | **d.im.s** for 'Deimos Fr!!!' |
| \* | Repetition of the previous character from **0 to x times** | **!\*** for 'Deimos Fr **!!!**' |
| + | Repetition of the previous character from **1 to x times** | **!+** for 'Deimos Fr **!!!**' |
| ? | Repetition of the previous character from **0 to 1 time** | **F?** for 'Deimos **F**r!!!' |
| \ | Escape character | **.** for 'Deimos Fr\.' |
| a,b,...z | Specific character | **D**eimos for 'Deimos Fr' |
| \w | **Alphanumeric** character (a...z,A...Z,0...9) | **\w**eimos for 'Deimos Fr' |
| \W | **Anything except** an **alphanumeric** character | I**\W**ll for 'I**'**ll be back' |
| \d | A **digit** | \d for 1 |
| \D | **Anything except** a **digit** | \Deimos for **D**eimos |
| \s | A **spacing character** such as: space, tab, carriage return or line break (\f,\t,\r,\n) | 'Deimos**\s**Fr' for 'Deimos Fr' |
| \S | **Anything except** a **spacing character** | 'Deimos**\S**Fr' for 'Deimos Fr' |
| {x} | **Repeats** the previous character exactly **x times** | !{3} in 'Deimos Fr **!!!**' |
| {x,} | **Repeats** the previous character **at least x times** | !{2} in 'Deimos Fr **!!!**' |
| {, x} | **Repeats between 0 and x times** the previous character | !{3} in 'Deimos Fr **!!!**' |
| {x, y} | **Repeats between x and y times** the previous character | !{1, 3} in 'Deimos Fr **!!!**' |
| [] | Allows setting a **range** (from a to z[a-z], from 0 to 9[0-9]...) | [A-D][0-5] in 'A4' |
| [^] | Allows specifying **unwanted characters** | [^0-9]eimos in 'Deimos' |
| () | Allows **recording the content of parentheses** for later use | (Deimos) in 'Deimos Fr' |
| \| | Allows doing an **exclusive or** | (Org\|Fr\|Com) in 'Deimos **Fr**' |
{{< /table >}}

There's a site allowing visualizing a regex: [Regexper](https://www.regexper.com/) (sources: https://github.com/javallone/regexper)

### Searching all patterns

If you want to find all occurrences of a string, use findall(). It works the same way as search() but returns a list of strings corresponding to the pattern searched:

```python
>>> pattern = "\d{4}"
>>> years = "Years: 1981, 1982, 1983"
>>> for m in re.findall(pattern, years):
...    print "Year:", m
...
Year: 1981
Year: 1982
Year: 1983
```

If you want to get an element providing you with the same information as search(), it's the finditer() function you'll need to use:

```python
>>> pattern = "\d{4}"
>>> years = "Years: 1981, 1982, 1983"
>>> for m in re.finditer(pattern, years):
...    print 'Pattern', m.re.pattern, 'found on position:', m.start(), ',', m.end()
...
Pattern \d{4} found on position: 8 , 12
Pattern \d{4} found on position: 14 , 18
Pattern \d{4} found on position: 20 , 24
```

### Memorizing and using found patterns

Regular expressions allow referencing an element noted in parentheses. We can then call them by referring (from left to right) to parenthesis number n (preceded by a backslash). To summarize, for the first parenthesis encountered, you use \1, for the second \2, etc...:

```python
>>> pattern = r'^Nom: ([A-Z][a-z]+) Prenom: ([A-Z][a-z]+) Mail: \2\1@deimos.fr$'
>>> match = re.search(pattern, 'Nom: Deimos Prenom: Fr Mail: DeimosFr@deimos.fr')
```

We preceded the pattern r'...'. This allows determining a raw string that will interpret backslashes differently (they don't protect any character and must therefore be considered as characters in their own right).

For better readability, it's possible to name the captured elements:

```python
(?P<name>pattern)
```

And use them:

```python
(?P=name)
```

Example:

```python
>>> pattern = r'^Nom: (?P<name>[A-Z][a-z]+) Prenom: (?P<firstname>[A-Z][a-z]+) Mail: (?P=firstname)(?P=name)@deimos.fr$'
>>> match = re.search(pattern, 'Nom: Deimos Prenom: Fr Mail: DeimosFr@deimos.fr')
```

#### Groups

We can also retrieve the captured value after a call to a search()-type function. Elements are retrieved in a list thanks to the groups() function and can be obtained using the group() function which takes as parameter the number of the element (starts at 1, and for index 0 the function returns the whole string).

```python
>>> pattern = '^Nom: ([A-Z][a-z]+) Prenom: ([A-Z][a-z]+)'
>>> import re
>>> match = re.search(pattern, 'Nom: Fr Prenom: Deimos$')
>>> match.groups()
('Fr', 'Deimos')
>>> match.group(0)
'Nom: Fr Prenom: Deimos'
>>> match.group(1)
'Fr'
>>> match.group(2)
'Deimos'
```

We can also get them in a dictionary using the groupdict() function (if you've named the captured elements, this is the method you'll need to choose).
Let's first see an example of application with the groups() and group() functions:

```python
>> pattern = '^Nom: (?P<name>[A-Z][a-z]+) Prenom: (?P<firstname>[A-Z][a-z]+)$'
>>> import re
>>> match = re.search(pattern, 'Nom: Fr Prenom: Deimos')
>>> match.groupdict()
{'name': 'Fr', 'firstname': 'Deimos'}
>>> m = match.groupdict()
>>> m['name']
'Fr'
>>> m['firstname']
'Deimos'
```

### Search parameters

It's possible in a search to indicate certain parameters, we'll see them here:

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| IGNORECASE | Performs a case-insensitive search |
| MULTILINE | The search can contain several lines separated by a line break (character \n) |
| DOTALL | The search can contain several lines separated by a dot (character .) |
| UNICODE | Allows using Unicode encoding for the search (useful with accented characters). Your strings will need to be given in Unicode format by prefixing them with u (Useless in Python 3.2): u'Unicode string' |
{{< /table >}}

You can also specify several options with a pipe:

```python
>>> pattern = '^deimos$'
>>> match = re.search(pattern, 'deimos\nfr', re.MULTILINE | re.IGNORECASE)
>>> match.start()
0
```

### Parameters in patterns

The previous parameters can be indicated directly in the patterns (?option):

{{< table "table-hover table-striped" >}}
| Pattern option | Correspondence |
|---------------|----------------|
| i | IGNORECASE |
| m | MULTILINE |
| s | DOTALL |
| u | UNICODE |
| x | VERBOSE |
{{< /table >}}

Here's an example:

```python
>>> pattern = '(?i)(?m)^deimos$'
>>> match = re.search(pattern, 'deimos\nfr')
>>> match.start()
0
```

### Substitution

Substitution is done using the sub() method and returns the modified string:

```python
>>> pattern = '(org|com)$'
>>> re.sub(pattern, 'fr', 'deimos com')
'deimos fr'
```

For named captured elements, we'll use \g<name>:

```python
>>> pattern = '(?P<name>org|com)$'
>>> re.sub(pattern, r'\g<name>', 'deimos fr')
'deimos fr'
```

**sub() performs substitutions for all occurrences of the pattern found!**

To limit substitutions to the first n strings matching the pattern, use the subn() function:

```python
>>> pattern = '(org|com)$'
>>> re.subn(pattern, r'\g<name>', 'deimos fr deimos fr deimos fr', count=2)
('deimos fr deimos fr deimos org', 2)
```

#### Split

Finally, let's note the existence of a split() function where the separation characters are determined by a pattern:

```python
>>> pattern = '-\*-|--|\*'
>>> re.split(pattern, '1-*-2--3--4*5*6-*-7')
['1', '2', '3', '4', '5', '6', '7']
```

## argv

In Python, the simplest way to work taking into account the input of arguments on command lines is to use the sys module.

```python
import sys
print 'Number of arguments:', str(len(sys.argv) – 1)
print 'Script name:', sys.argv[0]
print 'Parameters:'

for i in range(1, len(sys.argv)):
   print 'sys.argv[' + str(i) + '] = ' + sys.argv[i]
```

**You'll notice that you need to do -1 to know the number of arguments!**

When launching the script:

```bash
> python example.py 1 2 'Deimos Fr'
Number of arguments: 3
Script name: example.py
Parameters:
sys.argv[1] = 1
sys.argv[2] = 2
sys.argv[3] = Deimos Fr
```

## argparse

The argparse module is the new version of the optparse module which is now deprecated. We'll create an ArgumentParser object that will contain the list of arguments:

```python
import argparse
parser = argparse.ArgumentParser(description='Description', epilog='Epilog')
```

- description: allows indicating by a small comment what your script does
- epilog: the text that will be displayed at the end of the automatically generated help. These texts will be used when displaying the automatically generated help.

### Adding arguments to argparse

Let's define a -a option that collects no information and will only indicate if the user used it when calling the command (True) or not (False):

```python
parser.add_argument('-a', action='store_true', default=False, help='Description of the argument a')
```

Don't forget that this module automatically handles help :-). However, we can declare the command version:

```python
parser.add_argument('-v', '--version', action='store_true', help='Command version')
```

We can use here -v or --version, knowing that it's the first argument containing -- that will be the one to use when retrieving arguments.

- If you want to make an argument mandatory:

```python
parser.add_argument('-a', action='store', required=True)
```

- You can specify the number of necessary arguments:

```python
parser.add_argument('-a', action='store', nargs=2)
```

Here the number of expected arguments is 2 and will be returned in a list (args.a). You can set **the possible presence of this variable with a ?** to nargs.
Just like regex, it's also possible to put **\* to nargs to specify 0 or more arguments**. And finally the **+ for 1 to more arguments**.

#### Argument actions

There are several possible actions for arguments:

- store: This is the default action that saves the value the user will have entered
- store_const: Allows defining a default value if the user enters the argument. Otherwise the args.v variable will contain the value 'None':

```python
parser.add_argument('-v', action='store_const', const='Default value', help='Give help')
args = parser.parse_args()
```

- store_true/store_false: If the argument is specified, the value will be true, otherwise false. In case of absence of argument, the args.v variable takes the boolean value opposite to that indicated by the store action.
- append: Records the value in a list. If several calls to the same argument are specified, then they are added to this same list:

```python
parser.add_argument('-l', action='append', help='Values list')
args = parser.parse_args()
```

- append_const: It's a mix between store_const and append. It allows saving predefined elements in a list:

```python
parser.add_argument('-l', action='append_const', const='Default value', help='Values list')
args = parser.parse_args()
```

- version: Allows indicating the software version.

```python
parser.add_argument('-v', action='version', version='%(prog)s v0.1 License GPLv2', help='Version')
args = parser.parse_args()
```

%(prog)s refers to the name of the running script (like sys.argv[0]).

#### Choosing your argument prefix characters

If we want to use characters other than - or --, it's possible for example to add +:

```python
parser = argparse.ArgumentParser(description='Description', prefix_chars='-+')
```

We can for example have a use case where we add options (+option) or remove them (-option).

#### Forcing values

It's possible to force values entered by the user by specifying a type (int, float, bool, file...). This allows converting data on the fly and checking that the values entered by the user are of the expected type:

```python
parser.add_argument('-i', action='store', type=int, help='Give an integer')
args = parser.parse_args()
```

We can also define a default value in all circumstances:

```python
parser.add_argument('-i', action='store', type=int, default=10, help='Give an integer')
args = parser.parse_args()
```

**Note: A call to the command without specifying an -a argument will still initialize the args.i variable with the value 10.**

#### Managing duplicates

If we declare the same argument twice, it's possible (but not clean) to have Python resolve this automatically. The last value will then overwrite the first one(s):

```python
parser.add_argument('--other_a', '-a', action='store')
parser.add_argument('-a', action='store')
parser = argparse.ArgumentParser(description='Description', conflict_handler='resolve')
```

#### Option groups

- For better readability, it's possible to group certain options when displaying help:

```python
group1 = parser.add_argument_group('login')
group1.add_argument('-u', action='store', metavar='name', help='Login')
group1.add_argument('-p', action='store', metavar='password', help='Password')

group2 = parser.add_argument_group('io')
group2.add_argument('-i', action='store', metavar='filename', help='Input file')
group2.add_argument('-o', action='store', metavar='filename', help='Output file')
```

- It's also possible to define groups where only one option is selectable under penalty of being rejected if more than one is requested:

```python
group = parser.add_mutually_exclusive_group('Filesystems')
group.add_argument('-z', action='store', metavar='zfs', help='Choose ZFS')
group.add_argument('-b', action='store', metavar='btrfs', help='Choose BTRFS')
```

### Using the command

The script we just created has the -a option and **-h for help which is automatically generated**:

```bash
> python example.py
usage: example.py [-h] [-a]

Description of my command

optional arguments:
   -h, --help show this help message and exit
   -a Description of the argument a
```

### Using arguments

To display the list of arguments:

```python
print args
```

To display the value of the -a argument:

```python
print args.a
```

### Retrieving arguments

To retrieve the arguments that the user will have entered:

```python
args = parser.parse_args()
```

## References

http://inforef.be/swi/python.htm
http://diamond.izibookstore.com/produit/50/9786000048501/Linux%20Pratique%20HS%20n23%20%20La%20programmation%20avec%20Python
http://docs.python.org/library/string.html
