---
weight: 999
url: "/Les_différentes_boucles_du_shell_script/"
title: "Different Shell Script Loops"
description: "A guide to different loop structures in shell scripting including if, for, while, until, case statements and testing conditions."
categories: ["Linux", "Development"]
date: "2008-10-21T15:22:00+02:00"
lastmod: "2008-10-21T15:22:00+02:00"
tags: ["Shell", "Scripting", "Bash", "Programming"]
toc: true
---

## 1. if

Conditional binary structure: depending on whether a condition is true or false, we execute a block or we don't.

```bash
if condition ; then
    instruction
fi
```

If the condition is true, then instruction (which can be a block of instructions) is executed.

```bash
if condition ; then
    instruction1
else
    instruction2
fi
```

For the same construction but with multiple conditions:

```bash
if [condition1] || [condition2] ; then
    instruction1
else
    instruction2
fi
```

If the condition is true, then instruction1 is executed and the second is ignored, otherwise instruction2 is executed and the first is ignored.

```bash
if condition1 ; then
    instruction1
elif
    condition2
then
    instruction2
fi
```

elif is equivalent to else if. Thus instruction2 is only executed if condition1 and condition2 are both true at the same time.

## 2. for

Repetitive bounded structure: loop for which we know the total number of iterations before the first pass.

```bash
for variable in set ; do
    instructions
done
```

We therefore vary the variable by making it take all the values of the set successively.

## 3. while

Repetitive structure as long as the condition is true.

```bash
while condition ; do
    instructions
done
```

Loop whose instructions are executed as long as the condition expression is true.

Example:

```bash
nb_cdk=3;
i=0;
while [ "$i" -lt "$nb_cdk" ];do
        echo "salut bahan\n"
        i=$(expr $i + 1)
done
```

or

```bash
nb_cdk=3;
declare -i i; i=0
while [ "$i" -lt "$nb_cdk" ];do
        echo "salut bahan\n"
        i=$i+1
done
```

## 4. until

Repetitive structure until the condition is true (i.e. as long as it is false).

```bash
until condition ; do
    instructions
done
```

Loop whose instructions are executed as long as the condition expression is false.

## 5. case

Multiple choice conditional structure: depending on the value of the string expression, we can execute a wide range of instructions.

```bash
case string in
    pattern_1)
      instruction1
    ;;
    pattern_2)
      instruction2
    ;;
    *)
      the rest
    ;;
esac
```

If the string is similar to pattern_1 (which is a path expansion string, accepts meta characters * ?, [], {} and ~) then instruction1 is executed.

The string can perfectly verify several different patterns simultaneously, in this case the corresponding instructions (or blocks of instructions) will all be executed.

## 6. Tests

The conditions of the structures can be evaluated using the test command.

### 6.1 String Tests

```bash
test -z string
```

Returns true if string is empty.

```bash
test -n string
```

Returns true if string is not empty.

```bash
test string_1 = string_2
```

Returns true if the two strings are equal.

```bash
test string_1 != string_2
```

Returns true if the two strings are not equal.

### 6.2 Numeric Tests

```bash
test string_1 operator string_2
```

Returns true if strings string_1 and string_2 are in the relational order defined by the operator which can be one among those in the table below.

{{< table "table-hover table-striped" >}}
| Operator | Meaning |
|----------|---------|
| -eq | = |
| -ne | <> |
| -lt | < |
| -le | <= |
| -gt | > |
| -ge | >= |
{{< /table >}}

### 6.3 File Tests

```bash
test condition file
```

Returns true if the file meets the condition which can be one among those described in the table below.

{{< table "table-hover table-striped" >}}
| Operator | Returns true if |
|----------|----------------|
| -p | it's a named pipe |
| -f | regular file |
| -d | directory |
| -c | special file in character mode |
| -b | special file in block mode |
| -r | read access |
| -w | write access |
| -x | execution access |
| -s | non-empty file |
{{< /table >}}

### 6.4 Other Operators

It is possible to combine tests by using parentheses and the boolean operators of the following table:

{{< table "table-hover table-striped" >}}
| Operator | Meaning |
|----------|---------|
| ! | negation |
| -a | conjunction (and) |
| -o | disjunction (or) |
{{< /table >}}
