---
weight: 999
url: "/Introduction_au_C/"
title: "Introduction to C"
description: "A comprehensive introduction to the C programming language, covering basics, syntax, functions, data types, memory management, and more."
categories: ["Programming", "Linux"]
date: "2011-01-18T12:50:00+02:00"
lastmod: "2011-01-18T12:50:00+02:00"
tags: ["Development", "10.4 strlen", "4.3.2 Les variables", "GitHub", "3 Les macros", "6.1 Conversions", "7.1 if", "7.5 break", "4.1.4 Type char", "Mac OS X"]
toc: true
---

## Introduction

The history of the C language is intimately linked to that of the UNIX operating system. In 1965, Ken Thompson, from Bell Labs, developed an operating system that he called MULTICS (Multiplexed Information and Computing System) in order to run a game he had created, which gave birth in 1970 to the UNICS operating system (Uniplexed Information and Computing System), quickly renamed UNIX.

At the time, assembly language was the only language that allowed the development of an operating system. Ken Thompson then developed a higher-level language, the B language (whose name comes from BCPL, a subset of the CPL language, itself derived from Algol, a language that was popular at the time), to facilitate the writing of operating systems. It was a weakly typed language (an untyped language, as opposed to a typed language, is a language that manipulates objects in their binary form, without any notion of type (character, integer, real, etc.)) and too dependent on the PDP-7 (the machine on which UNIX was developed) to allow UNIX to be ported to other machines. So Denis Ritchie (who was, along with Ken Thompson, one of the creators of UNIX) and Brian Kernighan improved the B language to give birth to the C language. In 1973, UNIX was rewritten entirely in C. For 5 years, the C language was limited to internal use at Bell until the day Brian Kernighan and Denis Ritchie published a first definition of the language in a book entitled "The C Programming Language". This was the beginning of a revolution in the world of computing.

Thanks to its power, the C language quickly became very popular and in 1983, ANSI (American National Standards Institute) decided to standardize it by also adding some modifications and improvements, which gave birth in 1989 to the language as we know it today.

The characteristics of the C language are the following:

* Universality: the programming language par excellence, C is not confined to a particular field of application. It can be used both for writing operating systems and scientific or management programs, modern software, databases, compilers, assemblers or interpreters, etc.
* Flexibility: it is a concise, very expressive language, and programs written in this language are very compact thanks to a powerful set of operators. Your only limit is your imagination!
* Power: C is a high-level language but allows low-level operations and access to system functionalities, which is most of the time impossible in other high-level languages.
* Portability: it is a language that does not depend on any hardware or software platform. C also allows you to write portable programs, i.e. programs that can be compiled for any platform without any modification.

Moreover, its popularity but especially the elegance of programs written in C is such that its syntax has influenced many languages including C++ (which is considered a superset of C), JavaScript, Java, PHP and C#.

## Programs and functions

### First Program

Let's write a simple "Hello World" program:

```c
#include <stdio.h>

int main()
{
    printf("Hello, world\n");
    return 0;
}
```

**main** is the main function of the code you are going to write. This is where everything will start when you call your program.

According to the official C language standard, main is a function that **must return an integer (int)**. In many systems (including Windows and UNIX), this integer is called the application's error code. In C, although this is not necessarily the case for the operating system, we return 0 to say that everything went well.

As in many languages, we need to declare each variable we will use. Here it's simple, the preprocessor loader stdio.h will take care of it for us.

Next, we'll compile it to see any errors and to run our first program :-)

```bash
gcc -Wall hello_world.c -o hello_world
```

* gcc: the command corresponding to the compiler used.
* -Wall: enabling warning mode for possible errors. Very useful for debugging
* hello_world.c: source file
* hello_world: destination binary

#### Comments

For comments, here are the solutions:

```c
// Comment in C++ style but should work with all recent compilers
```

or

```c
/* You can start a paragraph
of comments without worrying
about anything until the end */
```

### Functions

In mathematics, a function is defined as follows:

```
f(x) = x² - 3
```

This means that f is a function that receives a real number x as an argument and returns a real number: x² - 3.

Let's write a C function that we'll call f, which takes an integer x as an argument and also returns an integer: x² - 3:

```c
int f(int x)
{
    return x*x - 3;
}
```

Now, let's use this function:

```c
#include <stdio.h>

int f(int); /* declaration of function f */

int main()
{
    int x = 4;
    printf("f(%d) = %d\n", x, f(x));
    return 0;
}

int f(int x)
{
    return x*x - 3;
}
```

int f(int x) says: f is a function that requires an int as an argument and returns an int.
The %d in the printf function is what we call a format specifier. It informs about the way we want to display the text. Here, we want to display the numbers 4 and 13 (f(4)). So we tell printf to use the "integer number" (%d) format to display them. The first %d corresponds to the format we want to use to display x and the second for f(x).

Also note that the variable x in the main function has absolutely nothing to do with the variable x in the parameter of the function f. Each function can have its own variables and completely ignores what happens in other functions.

*Note*: In a declaration, you can put the names of the function's arguments (good only for decoration and nothing else):

```c
int Surface(int Longueur, int largeur);
```

It is also strongly advised to declare what type of function argument you will use. Rather use:

```c
int Surface(void); /* void: 'empty', or 'nothing' if you prefer */
```

instead of not putting anything like:

```c
int Surface(); /* Surface is a function. Period. */
```

In a definition (implementation):

* You can omit the return type of a function. In this case, it is assumed to return an int.
* An empty pair of parentheses means that the function does not accept any arguments.

Other remarks:

* The declaration of a function is only necessary when its use precedes its definition. However, it is always advisable to define a function only after its use (which therefore requires a declaration) if only for the readability of the program (indeed, it is the program that we want to see at first sight, not the small details).
* A function may not return a value. Its return type is then void.

## Macros

### The preprocessor

Before being effectively compiled, C source files are processed by a preprocessor that resolves certain directives given to it, such as file inclusion for example. The preprocessor, although being a program independent of the compiler, is an indispensable element of the language.

A directive given to the preprocessor always starts with #. We have already encountered the include directive which allows to include a file. The define directive allows to define macros.

A macro, in its simplest form, is defined as follows:

```c
#define <macro> <the replacement text>
```

To replace for example all occurrences of PLUS by +:

```c
#define PLUS +
```

Similarly, it can even replace functions sometimes:

```c
#define carre(x) x * x
```

In this case, carre(3) will be replaced by 3 * 3. And finally:

```c
#define PI 3.14
```

The compiler replaces here each occurrence of PI with 3.14. You can also make a symbolic constant:

```c
#define USER " Toto "
```

### Global variables across multiple files

By default, a global variable is only accessible in the source file in which it is declared, because each source file will be compiled into an independent object file. However, it is possible to make a global variable accessible in all source files of a program using the extern keyword. This practice, however, should be avoided:

```c
extern my_extern_var
```

## Expressions and instructions

### Data types

{{< table "table-hover table-striped" >}}
|C Type|Corresponding Type|
|-|-|
|char|character (small integer)|
|int|integer|
|float|floating point number (real) in single precision|
|double|floating point number (real) in double precision|
{{< /table >}}

For example:

```c
char ch;
unsigned char c;
unsigned int n; /* or simply: unsigned n */
```

**The smallest possible value that can be assigned to an unsigned integer variable is 0, while signed integers accept negative values.**

#### int

You can also put short or long before int, in which case you would get a short integer (short int or simply short) respectively a long integer (long int or simply long). Here are examples of valid declarations:

```c
int n = 10, m = 5;
short a, b, c;
long x, y, z = -1;
unsigned long p = 2;
```

long can also be put before double, the resulting type is then long double (quadruple precision).

Before char or int, you can put the modifier **signed or unsigned** depending on whether you want to have a signed (by default) or unsigned integer. signed int is the signed value (negative or positive). unsigned int represents positive values.

{{< table "table-hover table-striped" >}}
|Type Name|Other Name|Range|Bytes|
|-|-|-|-|
|int|signed, unsigned int|-32,768 to 32,767|2|
|short|short_int, signed short, signed short_int|-32,768 to 32,767|2|
|long|long_int, signed long, signed long_int|-2,147,493,648 to 2,147,483,647|4|
|unsigned|unsigned int|0 to 65,535|2|
|unsigned short|unsigned short_int|0 to 65,535|2|
|unsigned long|unsigned long_int|0 to 4,294,967,295|4|
{{< /table >}}

#### Integer Numbers

* Any "pure" literal constant of integer type (e.g., 1, -3, 60, 40, -20, ...) is considered by the language to be of type int.
* To explicitly specify that an integer literal constant is of type unsigned, simply add the suffix u or U to the constant. For example: 2u, 30u, 40U, 50U, ...
* Similarly, just add the suffix l or L to explicitly specify that an integer literal constant is of long type (you can use the UL suffix for example for unsigned long).
* An integer literal constant can also be written in octal (base 8) or hexadecimal (base 16). The hexadecimal notation is obviously much more used.
* A literal constant written in octal must be preceded by 0 (zero). For example: 012, 020, 030UL, etc.
* A literal constant written in hexadecimal must start with 0x (zero x). For example 0x30, 0x41, 0x61, 0xFFL, etc.

#### Floating Point Numbers (floats)

It refers to decimal real numbers (positive or negative) (floating point). Any "pure" literal constant of floating point type (e.g., 0.5, -1.2, ...) is considered to be of **double type**.

The suffix **f or F allows to explicitly specify a float**. Be careful, 1f is not valid because 1 is an integer constant. However, **1.0f is perfectly correct**. The suffix l or L allows to explicitly specify a long double.

A literal constant of floating point type is composed, in this order:

* a sign (+ or -)
* a sequence of decimal digits: the integer part
* a point: the decimal separator
* a sequence of decimal digits: the decimal part
* either of the two letters e or E: symbol of power of 10 (scientific notation)
* a sign (+ or -)
* a sequence of decimal digits: the power of 10

For example, the following literal constants represent floating numbers: 1.0, -1.1f, 1.6E-19, 6.02e23L, 0.5 3e8

{{< table "table-hover table-striped" >}}
|Type Name|Other Name|Range|
|-|-|-|
|Floats|-|-3.4E38 to +3.4E38|
|double||-1.7E308 to +1.7E308|
|long double||-1.7E308 to +1.7E308|
{{< /table >}}

#### Type char

char is used to designate characters:

{{< table "table-hover table-striped" >}}
|Type Name|Other Name|Range|Bytes|
|-|-|-|-|
|char|signed char|-128 to 127|1|
|unsigned char||0 to 255|1|
{{< /table >}}

* Basic data types contain only one value.
* Derived or aggregate types can contain more than one value.
  * Example: strings, arrays, structures, enumerations, unions, pointers.

### Format specification in sprintf

Here is the list of format codes that we will use most often:

{{< table "table-hover table-striped" >}}
|Format Code|Usage|
|-|-|
|c|Display a character|
|d|Display an int|
|u|Display an unsigned int|
|x, X|Display an integer in hexadecimal format|
|f|Display a float or double in decimal notation|
|e|Display a float or double in scientific notation with a small e|
|E|Display a float or double in scientific notation with a capital E|
|g, G|Display a float or double (uses the most appropriate format)|
|%|Display the '%' character|
{{< /table >}}

Additionally:

* h before d or u indicates that the argument is a short
* l before d, u, x, or X indicates that the argument is of type long
* L before f, e, E, or g indicates that the argument is of type long double

There are many more possibilities (which you need to know!) that I will not detail here (but we will use some of them when the time comes). So please don't hesitate to consult the documentation.

### Variables and constants

#### Constants

In C, there are 4 basic types:

* Integer constants.
* Floating point constants (real constants).
* Character constants.
* String constants.

The first 2 constants are numeric constants.

* Integer numeric values are: decimal, octal, and hexadecimal.
* In base 10: a decimal point, an exponent (Note: precision depends on compilers (min: 6 digits, max: 18 digits))
* Single character between apostrophes (128 characters):
  * Character: 0...9 A...Z a...z
  * ASCII code: 48...57 65...90 97...122
* String of characters between quotes.
  * The compiler automatically places the null character '\0' at the end of the string (invisible). "hello" is actually "hello\0" which is useful in programs to mark the end of a string. 'A' is different from "A" because 'A' is a character with ASCII value 48 while "A" is a string ("A\0") without numeric value. "A" takes up more space than 'A'.

Here's how to create variables:

```c
int a, b, c;
int i = 0, j, n = 10;
double x, y, z = 0.5;
```

A declaration with initialization gives:

```c
int n = 0;
```

While here is a declaration, followed by an assignment:

```c
int n;
n = 0;
```

A constant is a variable whose particularity is to be read-only (the value of a constant cannot be modified):

```c
const int n = 10;
```

**Don't forget to declare all your variables or constants or you'll get errors.**

#### Variables

A variable is a value of a defined type that can be modified:

```c
int a;
a = 10;
int a = 9;
```

### Definition of new types

C has a very powerful mechanism that allows the programmer to create new data types using the typedef keyword:

```c
typedef int ENTIER;
```

Defines the ENTIER type as nothing other than the int type. Nothing prevents us therefore from writing:

```c
ENTIER a, b;
```

Although in this case a simple #define could have been sufficient, it is always recommended to use typedef which is much safer.

To declare functions without defining them, prototypes are used. Structures can be declared as follows:

```c
class myclass;
```

### Pointers

As we know very well, the place where a program executes is memory, so all the data of the program (variables, functions, ...) are in memory. The C language has an operator **&** allowing to retrieve the memory address of a variable or a function. For example, if n is a variable, **&n designates the address of n**.

C also has an operator ***** allowing to access the content of the memory whose address is given. For example, let's assume we have:

```c
int n;
```

Then the following instructions are strictly identical:

```c
n = 10;
```

```c
*( &n ) = 10;
```

A pointer (or a pointer type variable) is a variable intended to receive an address. It is then said to point to a memory location. Access to the memory content is done through the * operator.

Here's how to declare a variable p intended to receive the address of a variable of type int:

```c
int *p;
int * o;
```

I put another one with o, just to show that in C you can put as many spaces as you want. And now, here are other examples:

```c
int * p1, p2, p3; /* Only p1 is of type int *. The others are simply int.*/
int *p1, *p2, *p3; /* Here obviously, they are all int*/
/*Simplified use with typedef:*/
typedef int * PINT;
PINT p1, p2, p3;
```

However, this would not have worked if we had defined PINT using a #define because it would lead us to the first example.

## Inputs and outputs

### Enter data typed at the keyboard with the scanf function

You need to understand 2 things before you start:

* If we want to display an integer (with printf), we must display the integer (the value of the variable).
* If we want to ask the user (the one who uses our program) to type a number and then put the number thus entered in a variable, we must provide the address of the variable in which we wish to store the number entered.

An example will help understanding:

```c
#include <stdio.h>

int main()
{
    int a, b, c;

    printf("This program calculates the sum of 2 numbers.\n");

    printf("Enter the value of a: ");
    scanf("%d", &a);

    printf("Enter the value of b: ");
    scanf("%d", &b);

    c = a + b;
    printf("%d + %d = %d\n", a, b, c);

    return 0;
}
```

Beware of spaces and what you ask someone to type with scanf. For example, if you want the person to type 'years' in addition to a number:

```c
scanf("%d years", &a);
```

The user who enters their age will have to type "x years". So be careful what has to be typed with scanf.

This is what we call formatted input. Functions such as scanf are rather intended to be used to read data from a safe program (through a file for example), not those from a human, which are subject to error. The format codes used in scanf are roughly the same as in printf, except for floating points in particular.

{{< table "table-hover table-striped" >}}
|Format Code|Usage|
|-|-|
|f, e, g|float|
|lf, le, lg|double|
|Lf, Le, Lg|long double|
{{< /table >}}

Here is a program that calculates the volume of a right circular cone according to the formula: V = 1/3 * (B * h) where B is the base surface or for a circular base: B = PI*R², where R is the radius of the base.

```c
#include <stdio.h>

double Volume(double r_base, double hauteur);

int main()
{
    double R, h, V;

    printf("This program calculates the volume of a cone.\n");

    printf("Enter the radius of the base: ");
    scanf("%lf", &R);

    printf("Enter the height of the cone: ");
    scanf("%lf", &h);

    V = Volume(R, h);
    printf("The volume of the cone is: %f", V);

    return 0;
}

double Volume(double r_base, double hauteur)
{
    return (3.14 * r_base * r_base * hauteur) / 3;
}
```

### Example of permutation of the contents of two variables

This function must therefore be able to locate variables in memory, in other words we must pass to this function the addresses of the variables whose content we want to exchange:

```c
#include <stdio.h>

void permuter(int * addr_a, int * addr_b);

int main()
{
    int a = 10, b = 20;

    permuter(&a, &b);
    printf("a = %d\nb = %d\n", a, b);

    return 0;
}

void permuter(int * addr_a , int * addr_b)
/***************\
* addr_a <-- &a *
* addr_b <-- &b *
\***************/
{
    int c;

    c = *addr_a;
    *addr_a = *addr_b;
    *addr_b = c;
}
```

### Common arithmetic operators

The common arithmetic operators +, -, *, and / exist in the C language. However, integer division is a little bit tricky. Indeed, if a and b are integers, a / b equals the quotient of a and b, i.e., for example, 29 / 5 equals 5. The remainder of an integer division is obtained with the modulo operator %, i.e., taking up the previous example, 29 % 5 equals 4.

### Comparison operators

{{< table "table-hover table-striped" >}}
|Operator|Role|
|-|-|
|<|Less than|
|>|Greater than|
|==|Equal to|
|<=|Less than or equal to|
|>=|Greater than or equal to|
|!=|Not equal to|
{{< /table >}}

### Logical operators

{{< table "table-hover table-striped" >}}
|Operator|Role|
|-|-|
|&&|AND|
|\|\||OR|
|!|NOT|
{{< /table >}}

```c
int prop1, prop2, prop_ou, prop_et, prop_vrai;
prop1 = (1 < 1000);
prop2 = (2 == -6);
prop_ou = prop1 || prop2; /* TRUE, because prop1 is TRUE */
prop_et = prop1 && prop2; /* FALSE, because prop2 is FALSE */
prop_vrai = prop1 && !prop_2 /* TRUE because prop1 and !prop2 are TRUE */
```

### Sizeof: Data size

The size of a piece of data refers to the size, in bytes, that it occupies in memory. By extension of this definition, the size of a data type refers to the size of a piece of data of this type. Caution! byte refers here, by abuse of language, to the size of a memory element on the target machine (the abstract machine), i.e., the size of a memory cell (which equals 8 bits in most current architectures), and not a group of 8 bits. In the C language, a byte (a memory cell) is represented by a char. The size of a char is therefore not necessarily 8 bits, even if it is the case in many architectures, but dependent on the machine. The standard requires, however, that a char must be at least 8 bits and that the CHAR_BIT macro, declared in limits.h, indicates the exact size of a char on the target machine.

C has an operator, sizeof, allowing to know the size, in bytes, of a piece of data or a data type. The size of a char is therefore obviously 1 since a char represents a byte. Moreover, there can't be a type whose size is not a multiple of that of a char. The type of the value returned by the sizeof operator is size_t, declared in stddef.h, which is included by many header files including stdio.h.

As we have already said above, the size of the data is dependent on the target machine. In the C language, the size of the data is therefore not fixed. Nevertheless, the standard stipulates that we must have:

```
sizeof (char) <= sizeof (short) <= sizeof (int) <= sizeofd(long)
```

On an Intel (x86) 32-bit processor for example, a char is 8 bits, a short 16 bits, and int and long 32 bits.

### Increment and decrement operators

Just like in many languages:

```c
#include <stdio.h>

int main()
{
    int i = 1, j;
    j = ++i; /* j = 1+1 => j=2 */
    printf ("j (++i) = %d\n", j); 
    j = --i; /* j = 2-1 => j=1 */
    printf ("j (--i) = %d\n", j); 
    return 0;
}
```

Whether you write ++i or i++ doesn't matter, the effect is the same.

### Conditional expression

A conditional expression is an expression whose value depends on a condition. The expression:

```c
p ? a : b
```

equals a if p is true and b if p is false.

For assignment operations, these are the operators: +=, -=, *=, /=, ...

```c
x += a;
```

for example is equivalent to:

```c
x = x + a;
```

The operators are classified in order of priority. Here are the operators we have studied so far classified in this order.

{{< table "table-hover table-striped" >}}
|Operator|Associativity|
|-|-|
|Parentheses|left to right|
|! ++ -- - (sign) sizeof|left to right|
|* / %|left to right|
|+ -|left to right|
|< <= > >=|left to right|
|== !=|left to right|
|& (address of)|left to right|
|&&|left to right|
|\|\||left to right|
|Assignment operators (= += ...)|right to left|
|,|left to right|
{{< /table >}}

Just because this order exists doesn't mean you have to memorize it. For readable code, it's even advised not to rely on it too much and to use parentheses in ambiguous situations.

## Characters

The numerical representation of characters defines what is called a character set. For example, in the ASCII character set (American Standard Code for Information Interchange), which is a character set that uses only 7 bits and is the basis of many popular codes today, the character 'A' is represented by code 65, the character 'a' by 97 and '0' by 48. Alas, even ASCII does not define the C language. Indeed, if C depended on a particular character set, it would then not be totally portable. Nevertheless, the standard defines a certain number of characters that any environment compatible with C must possess, among which are the 26 letters of the Latin alphabet (actually 52 since we differentiate uppercase and lowercase), the 10 decimal digits, the characters # < > ( ) etc. The programmer (but not the compiler) does not need to know how these characters are represented in the character set of the environment. So the standard does not define a character set but only a set of characters that each compatible environment is free to implement in its own way (plus any characters specific to that environment). The only constraint imposed is that their value must be able to fit in a char.

Concerning the escape technique, also know that you can insert octal (starting with 0) or hexadecimal (starting with x) code after the escape character \ to get a character whose code in the character set is given. Hexadecimal is by far the most used. For example: '\x30', '\x41', '\x61', ... And finally for characters with code 0, 1, ... up to 7, we can use the shortcuts '\0', '\1', ... '\7'.

Overflow occurs when trying to assign to an lvalue a value larger than it can hold. For example, by assigning a 32-bit value to a variable that can only hold 16 bits.

### Conversions

#### Implicit

In the C language, implicit conversion rules apply to data that make up a complex expression when they are not of the same type (integer with a float, short integer with long integer, signed integer with an unsigned integer, etc.). For example, in the expression:

```c
'A' + 2
```

'A' is of type char and 2 of type int. In this case, 'A' is first converted to int before the expression is evaluated. The result of the operation is of type int (because an int + an int gives an int). Here, it equals 67 (65 + 2). In fact, char and short are always systematically converted to int, i.e., in adding two char for example, both are first converted to int before being added, and the result is an int (not a char). An unsigned char will be converted to an unsigned int, and so on.

As a general rule: the "weakest" type is converted into the "strongest" type. For example, integers are weaker than floats, so 1 mixed with a float for example will first be converted to 1.0 before the operation actually takes place.

The compiler converts from lower rank to higher rank (= promotion):

```
char < short < int < long < float < double
```

#### Explicit (cast)

Simply specify the destination type in parentheses in front of the expression to convert. For example:

```c
float f;
f = (float)3.1416;
```

In this example, we explicitly converted 3.1416, which is of type double, to float. When a float is assigned to an integer, only the integer part, if it can be represented, is retained.

## Instructions

### if

Allows for conditional choices. The syntax of the instruction is as follows:

```c
if ( <expression> )
    <one and only one instruction>
else
    <one and only one instruction>
```

An if instruction may not have an else. When there are several nested if instructions, an else always relates to the last if followed by one and only one instruction. For example: let's write a program that compares two numbers.

```c
#include <stdio.h>

int main()
{
    int a, b;

    printf("This program compares two numbers.\n");

    printf("Enter the value of a: ");
    scanf("%d", &a);

    printf("Enter the value of b: ");
    scanf("%d", &b);

    if (a < b)
        printf("a is smaller than b.\n");
    else
        if (a > b)
            printf("a is greater than b.\n");
        else
            printf("a equals b.\n");

    return 0;
}
```

### do

Allows to perform a loop. The syntax of the instruction is as follows:

```c
do
    <one and only one instruction>
while ( <expression> );
```

The do instruction allows to execute an instruction as long as <expression> is true. The test is done after each execution of the instruction. Here is a program that displays "Hello" 10 times.

```c
#include <stdio.h>

int main()
{
    int nb_lignes_affichees = 0;

    do
        {
            printf("Bonjour.\n");
            nb_lignes_affichees++;
        }
    while (nb_lignes_affichees < 10);

    return 0;
}
```

### while

Allows to perform a loop. The syntax of the instruction is as follows:

```c
while ( <expression> )
    <one and only one instruction>
```

The while instruction allows to execute an instruction as long as <expression> is true. The test is done before each execution of the instruction. So if the condition (<expression>) is false from the start, the loop will not be executed.

### for

Allows to perform a loop. The syntax of the instruction is as follows:

```c
for ( <init> ; <condition> ; <step>)
    <instruction>
```

It is practically identical to:

```c
<init>;
while ( <condition> )
{
    <instruction>
    <step>
}
```

For example, let's write a program that displays the multiplication table for 5.

```c
#include <stdio.h>

int main()
{
    int n;

    for(n = 0; n <= 10; n++)
        printf("5 x %2d %2d\n", n, 5 * n);

    return 0;
}
```

The %2d format displays an integer with a minimum of 2 characters (the remaining space will be filled with spaces).

### break

Allows to immediately exit a loop or a switch. The syntax of this instruction is:

```c
break;
```

### Switch and case

These instructions allow to avoid if instructions that are too nested as illustrated by the following example:

```c
#include <stdio.h>

int main()
{
    int n;

    printf("Enter an integer: ");
    scanf("%d", &n);

    switch(n)
    {
    case 0:
        printf("Case of 0.\n");
        break;

    case 1:
        printf("Case of 1.\n");
        break;

    case 2: case 3:
        printf("Case of 2 or 3.\n");
        break;

    case 4:
        printf("Case of 4.\n");
        break;

    default:
        printf("Unknown case.\n");
    }

    return 0;
}
```

### Continue

In a loop, allows to immediately move to the next iteration. For example, let's modify the multiplication table program so that we display nothing for n = 4 or n = 6.

```c
#include <stdio.h>

int main()
{
    int n;

    for(n = 0; n <= 10; n++)
    {
        if ((n == 4) || (n == 6))
            continue;

        printf("5 x %2d %2d\n", n, 5 * n);
    }

    return 0;
}
```

### Return

Allows to terminate a function. The syntax of this instruction is as follows:

```c
return <expression>; /* terminates the function and returns <expression> */
```

or:

```c
return; /* terminates the function without specifying a return value */
```

## Arrays, pointers, and character strings

An array is a variable that groups one or more pieces of data of the same type. Access to an element of the array is done by an index system, the index of the first element being 0. For example:

```c
int t[10];
```

declares an array of 10 elements (of type int) whose name is t. The elements of the array go from t[0], t[1], t[2] ... to t[9]. t is a variable of type array, more precisely (in our case), a variable of type array of 10 int (int [10]). The elements of the array are int. All the rules applying to variables also apply to the elements of an array.

```c
char msg[ ] = "bonjour";
char msg[8] = "bonjour"; // the 8 corresponds to the number of letters in the word and bonjour plus the '\0'.
```

### Initialization

You can initialize an array using braces. For example:

```c
int t[10] = {0, 10, 20, 30, 40, 50, 60, 70, 80, 90};
```

Of course, we are not obliged to initialize all the elements, we could have stopped after the 5th element for example, and in this case the other elements of the array will automatically be initialized to 0. Caution! an uninitialized local variable contains "anything", not 0!

When declaring an array with initialization, you can omit the number of elements because the compiler will calculate it automatically. Thus, the declaration:

```c
int t[] = {0, 10, 20, 30};
```

is strictly identical to:

```c
int t[4] = {0, 10, 20, 30};
```

#### Calculate the size of an array

The size of an array is obviously the number of elements in the array multiplied by the size of each element. Thus, the number of elements in an array is equal to its size divided by the size of an element. We then generally use the formula sizeof(t) / sizeof(t[0]) to know the number of elements of an array t. The macro defined below allows to calculate the size of an array:

```c
#define COUNT(t) (sizeof(t) / sizeof(t[0]))
```

### Multidimensions

You can also create a multidimensional array. For example:

```c
int t[10][3];
```

A multidimensional array is actually nothing but an array (one-dimensional array) whose elements are arrays. As in the case of one-dimensional arrays, the type of the elements of the array must be perfectly known. So in our example, t is an array of 10 arrays of 3 int, or to help you see more clearly:

```c
typedef int TRIPLET[3];
TRIPLET t[10];
```

The elements of t go from t[0] to t[9], each being an array of 3 int.

You can of course create arrays with 3 dimensions, 4, 5, 6, ...

You can also initialize a multidimensional array. For example:

```c
int t[3][4] = { {0, 1, 2, 3},
                {4, 5, 6, 7},
                {8, 9, 10, 11} };
```

Which we could also have simply written:

```c
int t[][4] = { {0, 1, 2, 3},
               {4, 5, 6, 7},
               {8, 9, 10, 11} };
```

## Pointer Arithmetic

### Introduction to address calculation

Here is an example of a pointer:

```c
t[5] = '*';
```

In practice, we use a pointer to an element of the array, usually the first. This allows access to any element of the array by simple address calculation. As we said above: t + 1 is equivalent to &(t[1]), t + 2 to &(t[2]), etc.

Here's an example that shows one way to traverse an array:

```c
#include <stdio.h>

#define COUNT(t) (sizeof(t) / sizeof(t[0]))

void Affiche(int * p, size_t nbElements);

int main()
{
    int t[10] = {0, 10, 20, 30, 40, 50, 60, 70, 80, 90};

    Affiche(t, COUNT(t));

    return 0;
}

void Affiche(int * p, size_t nbElements)
{
    size_t i;

    for(i = 0; i < nbElements; i++)
        printf("%d\n", p[i]);
}
```

### Pointer arithmetic

Pointer arithmetic was born from the facts we established earlier. Indeed, if p points to an element of an array, p + 1 must point to the next element. So if the size of each element of the array is for example 4, p + 1 moves the pointer 4 bytes (where the next element is) and not one.

Similarly, since we should have (p + 1) - p = 1 and not 4, the difference between two addresses gives the number of elements between these addresses and not the number of bytes between these addresses. The type of such an expression is **ptrdiff_t**, which is defined in the file **stddef.h**.

And finally, the notation p[i] is strictly equivalent to *(p + i).

This shows how important pointer typing is. However, there are so-called generic pointers capable of pointing to anything. Thus, the conversion of a generic pointer to a pointer of another type for example requires no cast and vice versa.

### Generic pointers

The type of generic pointers is void *. Since these pointers are generic, the size of the data pointed to is unknown and pointer arithmetic therefore does not apply to them. Similarly, since the size of the pointed data is unknown, the indirection operator * cannot be used with these pointers, a cast is then mandatory. For example:

```c
int n;
void * p;

p = &n;
*((int *)p) = 10; /* p being now seen as an int *, we can then apply the operator *. */
```

Given that the size of any data is a multiple of that of a char, the type char * can also be used as a universal pointer. Indeed, a variable of type char * is a pointer to a byte, in other words it can point to anything. This proves to be practical sometimes (when you want to read the content of a memory byte by byte for example) but in most cases, it is always better to use generic pointers. For example, the conversion of an address of different type to char * and vice versa always requires a cast, which is not the case with generic pointers.

In printf, the format specifier %p allows to print an address (void *) in the format used by the system.

And finally, there is a macro, NULL, defined in stddef.h, indicating that a pointer does not point anywhere. Its interest is therefore to allow to test the validity of a pointer and it is advised to always initialize a pointer to NULL.

### Example with a multidimensional array

Let's say:

```c
int t[10][3];
```

Let's define the TRIPLET type by:

```c
typedef int TRIPLET[3];
```

So as to have:

```c
TRIPLET t[10];
```

Seen from a pointer, t represents the address of t[0] (which is a TRIPLET) so the address of a TRIPLET. By doing t + 1, we move one TRIPLET, i.e., 3 int.

On the other hand, t can be seen as an array of 30 int (3 * 10 = 30), so we can access any element of t using a pointer to int.

Let p be a pointer to int and let's do:

```c
p = (int *)t;
```

We then have, numerically, the following equivalences:

```c
t   p
t + 1   p + 3
t + 2   p + 6
...      
t + 9   p + 27
```

Now let's take the 3rd TRIPLET of t, i.e., t[2].

Since the first element of t[2] is at address t + 2 or p + 6, the second is at p + 6 + 1 and the third and last at p + 6 + 2. After this integer, we find ourselves at the first element of t[3], at p + 9.

In conclusion, for an array declared:

```c
<type> t[N][M];
```

we have the formula:

```c
t[i][j] = *(p + N*i + j) /* or even p[N*i + j] */
```

Where obviously: p = (int *)t.

And we can of course extend this formula for any dimension.

## Character strings

By definition, a character string, or simply: string, is a finite sequence of characters. For example, "Hello", "3000", "Hi!", "EN 4", ... are character strings. In the C language, a constant of type character string is written between double quotes, exactly as in the examples given above.

The length of a string is the number of characters it contains. For example, the string "Hello" contains 5 characters ('H', 'e', 'l', 'l', 'o'). Its length is therefore 5. In C, the strlen function, declared in the file string.h, allows to obtain the length of a string passed as an argument. Thus, strlen("Hello") equals 5.

### Representation of character strings in C

As we have already mentioned above, constants of type character string are written in C between double quotes. In fact, the C language does not really have a character string type. A string is simply represented using an array of characters.

However, functions manipulating strings must be able to detect the end of a given string. In other words, every character string must end with a character indicating the end of the string. This character is the '\0' character and is called the null character or end of string character. Its ASCII code is 0. Thus the string "Hello" is actually an array of characters whose elements are 'H', 'e', 'l', 'l', 'o', '\0', in other words an array of 6 characters and we have "Hello"[0] = 'H', "Hello"[1] = 'e', "Hello"[2] = 'l', ... "Hello"[5] = '\0'. However, since it is a constant (character string constant), the content of the memory allocated for the string "Hello" cannot be modified.

**The string manipulation functions of the C language standard library are mainly declared in the file string.h. Here is an example of using one of these functions:**

```c
#include <stdio.h>
#include <string.h>

int main()
{
    char t[50];

    strcpy(t, "Hello, world!");
    printf("%s\n", t);

    return 0;
}
```

In this example, the string t can contain at most 50 characters, including the end of string character. In other words, **t can only contain 49 "normal" characters** because a place must always be reserved for the end of string character: '\0'. You can also of course initialize a string at the time of its declaration, for example:

```c
char s[50] = "Hello";
```

Which is strictly equivalent to:

```c
char s[50] = { 'H', 'e', 'l', 'l', 'o', '\0'};
```

Since, seen from a pointer, the value of a string literal expression is nothing other than the address of its first element, you can use a simple pointer to manipulate a string. For example:

```c
char * p = "Hello";
```

In this case, p points to the first element of the string "Hello". However, as we have already said above, the memory allocated for the string "Hello" is read-only so you cannot write for example:

```c
p[2] = '*'; /* Forbidden */
```

With an array, it is not the address in memory of the string that is stored, but the characters of the string, copied character by character. Since the memory used by the array is independent of that used by the source string, we can do what we want with our array. The strcpy function allows to copy a string to another memory location.

The following paragraph discusses string manipulation functions in C.

### strcpy, strncpy

```c
#include <stdio.h>
#include <string.h>

int main()
{
    char t1[50], t2[50];

    strcpy(t1, "Hello, world!");
    strcpy(t2, "*************");
    strncpy(t1, t2, 3);
    printf("%s\n", t1);

    return 0;
}
```

Caution! If t1 is not large enough to hold the string to copy, you will have a **buffer overflow**. A buffer is simply an area of memory used by a program to temporarily store data. For example, t1 is a buffer of 50 bytes. It is therefore the programmer's responsibility not to pass anything to it! Indeed in C, the compiler assumes that the programmer knows what he is doing!

The strncpy function is used in the same way as strcpy. The third argument indicates the number of characters to copy. No end of string character is automatically added.

### strcat, strncat

This will add characters to the end of an array:

```c
#include <stdio.h>
#include <string.h>

int main()
{
    char t[50];

    strcpy(t, "Hello, world");
    strcat(t, " from");
    strcat(t, " strcpy");
    strcat(t, " and strcat");
    printf("%s\n", t);

    return 0;
}
```

Which will give:

```
Hello, world from strcpy and strcat
```

### strlen

Returns the number of characters in a string.

### strcmp, strncmp

We don't use the == operator to compare strings because it's not the addresses we want to compare but the memory content. The strcmp function compares two character strings and returns:

* zero if the strings are identical
* a negative number if the first is "less than" (from a lexicographical point of view) the second
* and a positive number if the first is "greater than" (from the same point of view ...) the second

Thus, as an example, in the expression

```c
strcmp("clandestin", "clavier")
```

The function returns a negative number because, 'n' being smaller than 'v' (in the ASCII character set, it has nothing to do with the C language), "clandestin" is smaller than "clavier".

### Implementation of some string manipulation functions

We will therefore implement two string manipulation functions, namely str_len and str_cpy, which will be used in the same way as their twins strlen and strcpy.

```c
size_t str_len(char * t)
{
    size_t len;

    for(len = 0; t[len] != '\0'; len++)
        /* We do nothing, we just let it loop */ ;

    return len;
}
```

```c
char * str_cpy(char * dest, char * source)
{
    int i;

    for(i = 0; source[i] != '\0'; i++)
        dest[i] = source[i];

    dest[i] = '\0';

    return dest;
}
```

Notice the way we implemented the str_cpy function. You might have expected this function to return void and not a char *. Well, no! Many functions in the standard library also use this "convention", which allows code of the type:

```c
char s[50] = "Hello";
printf("%s\n", strcat(s, " everyone!"));
```

Input/output (I/O) is not really part of the C language because these operations are dependent on the system. So to perform input/output operations in C, you have to in principle go through the low-level functionalities of the system. Nevertheless its standard library is provided with functions allowing to perform such operations in order to facilitate the writing of portable code. The functions and data types related to inputs/outputs are mainly declared in the file stdio.h (standard input/output).

Input/output in C is done through logical entities, called streams, which represent objects external to the program, called files. A file can be opened in reading, in which case it is supposed to provide us with data (i.e., to be read) or opened in writing, in which case it is intended to receive data from the program. A file can be both open for reading and writing. Once a file is open, a stream is associated with it. An input stream is a stream associated with a file open for reading and an output stream a stream associated with a file open for writing.

When the data exchanged between the program and the file are of text type, the need to define what we call a line is paramount. In C, a line is a sequence of characters ended by the end of line character (included): '\n'. For example, when we type on the keyboard, a line corresponds to a sequence of characters ended by ENTER. Since the ENTER key ends a line, the character generated by this key is therefore, in standard C, the end of line character or '\n'.

When the system executes a program, three files are automatically opened:

* the standard input, by default the keyboard
* the standard output, by default the screen (or console)
* and the standard error, by default associated with the screen (or console)

You know the rest, for redirections to files (< >).

### Read a character, then display it

The getc macro allows you to read a character on an input stream. The putc macro allows you to write a character to an output stream. Here's a simple program that shows how to use the getc and putc macros:

```c
#include <stdio.h>

int main()
{
    int c; /* the character */

    printf("Please type a character: ");
    c = getc(stdin);

    printf("You typed: ");
    putc(c, stdout);

    return 0;
}
```

You're probably wondering why we used int rather than char in the declaration of c. Well simply because getc returns an int (same for putc which expects an argument of type int). But precisely: Why? Well because getc must be able not only to return the read character (a char) but also a value that must not be a char to signal that no character could be read. **This value is EOF. It is defined in the stdio.h file**. In these conditions, it is clear that we can use anything except a char as the return type of getc.

One of the most frequent cases where getc returns EOF is when we have encountered the end of the file. The end of a file is a point located beyond the last character of this file (if the file is empty, the beginning and end of the file are therefore confused). We are said to have encountered the end of a file after having attempted to read in this file while we are already at the end, not just after having read the last character. When stdin is associated with the keyboard, the notion of end of file loses its meaning a priori because the user can very well type anything at any time. However the execution environment (the operating system) generally offers a way to specify that there is no more character to provide (concretely, for us, this means that getc will return EOF). Under Windows for example, it is sufficient to type at the beginning of the line the combination of keys Ctrl + Z (inherited from DOS) and then validate by ENTER. Obviously, everything starts again at zero at the next read operation.

The getchar and putchar macros are used as getc and putc except that they operate only on stdin, respectively stdout. They are defined in stdio.h as follows:

```c
#define getchar() getc(stdin)
#define putchar(c) putc(c, stdout)
```

And finally fgetc is a function that does the same thing as getc (which can in fact be either a function or a macro...). Similarly fputc is a function that does the same thing as putc.

### Enter a character string

It is enough to read the characters present on the input stream (in our case: stdin) until we have reached the end of the file or the end of line character. We will have to provide as arguments of the function the address of the buffer intended to contain the character string entered and the size of this buffer to eliminate the risk of buffer overflow.

```c
#include <stdio.h>

char * saisir_chaine(char * lpBuffer, size_t nBufSize);

int main()
{
    char lpBuffer[20];

    printf("Enter a character string: ");
    saisir_chaine(lpBuffer, sizeof(lpBuffer));

    printf("You typed: %s\n", lpBuffer);

    return 0;
}

char * saisir_chaine(char * lpBuffer, size_t nBufSize)
{
    size_t nbCar = 0;    
    int c;

    c = getchar();
    while ((nbCar < nBufSize - 1) && (c != EOF) && (c != '\n'))
    {
        lpBuffer[nbCar] = (char)c;
        nbCar++;
        c = getchar();
    }

    lpBuffer[nbCar] = '\0';

    return lpBuffer;
}
```

The scanf function also allows to enter a character string not containing any space (space, tab, etc.) thanks to the format specifier %s. It will therefore stop reading when it encounters a space (but before reading, it will first advance to the first character that is not a space). scanf finally adds the end of string character. The template allows to indicate the maximum number of characters to read (end of string character not included). When using scanf with the %s specifier (which asks to read a string without space), you should never forget to also specify the maximum number of characters to read (to be placed just before the s) otherwise the program will be open to buffer overflow attacks. Here is an example that shows the use of scanf with the %s format specifier:

```c
#include <stdio.h>

int main()
{
    char lpBuffer[20];

    printf("Enter a character string: ");
    scanf("%19s", lpBuffer);

    printf("You typed: %s\n", lpBuffer);

    return 0;
}
```

And finally, there is also a function, gets, declared in stdio.h, which allows to read a character string on stdin. However this function is to be proscribed because it does not allow to specify the size of the buffer that will receive the read string.

### Read a line with fgets

The fgets function allows to read a line (i.e., including the '\n') on an input stream and to place the read characters in a buffer. This function then adds the '\0' character. Example:

```c
#include <stdio.h>

int main()
{
    char lpBuffer[20];

    printf("Enter a character string: ");
    fgets(lpBuffer, sizeof(lpBuffer), stdin);

    printf("You typed: %s", lpBuffer);

    return 0;
}
```

In this example, two cases can arise:

* the user enters a string containing at most 18 characters and then validates it all by ENTER, then all the characters of the line, including the end of line character, are copied into lpBuffer and then the end of string character is added
* the user enters a string containing more than 18 characters (i.e., >= 19) and then validates it all by ENTER, then only the first 19 characters are copied to lpBuffer and then the end of string character is added

### Input/output mechanism in C

Input/output in C is buffered, meaning that the data to be read (respectively to be written) are not directly read (respectively written) but are first placed in a buffer associated with the file. The proof, you have certainly noticed for example that when you enter data for the first time using the keyboard, these data will only be read once you have pressed the ENTER key. Then, all the read operations that follow will be done immediately as long as the '\n' character is still present in the read buffer, i.e., as long as it has not yet been read. When the '\n' character is no longer present in the buffer, you will have to press ENTER again to validate the input, and so on.

Write operations are less complicated, but there is still something that would be totally unfair not to talk about. As we have already said above, input/output in C is buffered i.e., passes through a buffer. In the case of a write operation, it may happen that we want at a certain moment to force the physical writing of the data present in the buffer associated with the file without waiting for the system to finally decide to do it. In this case, we will use the fflush function:
In the case of an output stream, this function causes the immediate physical writing of the buffer being filled. It returns EOF in case of error, zero otherwise.
According to the official C language standard, the effect of fflush on a stream that is not an output stream is undefined. But for most current libraries, calling this function on an input stream removes the characters available in the buffer. For example, in the frequent case where the standard input corresponds to the keyboard, the call fflush(stdin) makes all the characters already typed but not yet read by the program disappear.

Note. If the physical file that corresponds to the indicated stream is an interactive organ, for example the screen of a workstation, then the fflush function is implicitly called in two very frequent circumstances:

* the writing of the '\n' character which produces the emission of an end of line mark and the effective emptying of the buffer,
* the beginning of a read operation on the associated input unit (interactive input/output organs generally form pairs); thus, for example, a read at the keyboard causes the emptying of the write buffer to the screen. This allows a question to be effectively displayed before the user has to type the corresponding answer.

#### Read safely from standard input

First, let's analyze the very small program below:

```c
#include <stdio.h>

int main()
{
    char nom[12], prenom[12];

    printf("Enter your last name: ");
    fgets(nom, sizeof(nom), stdin);

    printf("Enter your first name: ");
    fgets(prenom, sizeof(prenom), stdin);

    printf("Your last name is: %s", nom);
    printf("And your first name: %s", prenom);

    return 0;
}
```

In this program, if the user enters a name containing less than 10 characters and then validates by ENTER, then all the characters fit in nom and the program proceeds as planned. On the other hand, if the user enters a name containing more than 10 characters, only the first 11 characters will be copied to nom and characters are therefore still present in the keyboard buffer. So, when reading the first name, the characters still present in the buffer will immediately be read without the user having been able to enter anything. Here is a second example:

```c
#include <stdio.h>

int main()
{
    int n;
    char c;

    printf("Enter a number (integer): ");
    scanf("%d", &n);

    printf("Enter a character: ");
    scanf("%c", &c);

    printf("The number you entered is: %d\n", n);
    printf("The character you entered is: %c\n", c);

    return 0;
}
```

When we ask scanf to read a number, it will move the pointer to the first non-blank character, read as long as it should read the characters that may appear in the expression of a number, then stop when it encounters an invalid character (space or letter for example).

So in the example above, the character will be read without the user's intervention because of the presence of the '\n' character (which will then be the read character) due to the ENTER key pressed during the number input.

These examples show us that in general, we must always empty the keyboard buffer after each input, unless it is already empty of course. To empty the keyboard buffer, simply eat all the characters present in the buffer until we have encountered the end of line character or reached the end of the file. As an example, here is an improved version (with input buffer emptying after reading) of our saisir_chaine function:

```c
char * saisir_chaine(char * lpBuffer, int nBufSize)
{
    char * p;

    fgets(lpBuffer, nBufSize, stdin);

    p = lpBuffer + strlen(lpBuffer) - 1;
    if (*p == '\n')
        *p = '\0'; /* we overwrite the \n */
    else
    {
        /* We empty the read buffer of the stdin stream */
        int c;

        c = getchar();
        while ((c != EOF) && (c != '\n'))
           c = getchar();
    }

    return lpBuffer;
}
```

## Dynamic memory allocation

The interest of dynamically allocating memory is felt when we want to create an array whose size we need is only known at execution time for example. In other words:

```c
int t[10];
...
/* END */
```

Can be replaced by:

```c
int * p;

p = malloc(10 * sizeof(int));
...
free(p); /* free the memory when we no longer need it */
/* END */
```

The malloc and free functions are declared in the stdlib.h file. malloc returns NULL in case of failure. Here is an example that illustrates a good way to use them:

```c
p = malloc(10 * sizeof(int));
if (p != NULL)
{
    ...
    free(p);
}
else
    /* FAILURE */
```

The realloc function:

```c
void * realloc(void * memblock, size_t newsize);
```

allows to "resize" a dynamically allocated memory (by malloc for example). If memblock equals NULL, realloc behaves like malloc. In case of success, this function then returns the address of the new memory, otherwise the NULL value is returned and the memory pointed to by memblock remains unchanged.

```c
int * p, * q; /* q: to test the return of realloc */

p = malloc(10 * sizeof(int));
if (p != NULL)
{
    ...

    q = realloc(p, 20 * sizeof(int));

    if (q != NULL)
    {
        p = q;

        ...
    }

    free(p);
}

...
```

## Resources
- http://c.developpez.com/cours
- [C in Action - O'Reilly](https://books.google.fr/books?id=dsfXx4ESnM8C&printsec=frontcover)
