---
weight: 999
url: "/Introduction_au_SQL/"
title: "Introduction to SQL"
description: "A comprehensive introduction to SQL, covering basic commands like CREATE, INSERT, UPDATE, DELETE and SELECT, as well as advanced topics like sorting, filtering and wildcards."
categories: ["Development"]
date: "2010-04-10T22:10:00+02:00"
lastmod: "2010-04-10T22:10:00+02:00"
tags: ["SQL", "Database", "Development"]
toc: true
---

## Introduction

[Structured Query Language (SQL)](https://fr.wikipedia.org/wiki/SQL) is a standardized and normalized pseudo-language (query type), designed to query or manipulate a relational database with:

- A Data Definition Language (DDL)
- A Data Manipulation Language (DML), the most common and visible part of SQL
- A Data Control Language (DCL)
- A Transaction Control Language (TCL)
- And other modules designed to write routines (procedures, functions, or triggers) and interact with external languages.

SQL is part of the same family as SEQUEL (of which it is a descendant), QUEL, or QBE (Zloof) languages.

## Create - Data Insertion

The INSERT command adds a row to a database table:

```sql
INSERT INTO TABLE (column1[, column2, column3, ...]) VALUES (value1[, value2, value3, ...])
```

The number of columns appearing in parentheses before VALUES must match the number of values in the parentheses after. To insert a row that contains values for only certain columns, simply indicate the relevant columns and their corresponding values:

```sql
INSERT INTO dishes (dish_name, is_spicy) VALUES ('Salt Baked Scallops', 0)
```

The column list can be omitted if you are inserting values for all columns in a row:

```sql
INSERT INTO dishes VALUES (1, 'Braised Sea Cucumber', 6.50, 0)
```

## UPDATE - Data Update

The UPDATE command modifies data already present in a table:

```sql
UPDATE tablename SET column1=value1[, column2=value2, column3=value3, ...] [WHERE where_clause]
```

Here's an example of initializing a column with a string or number:

```sql
; CHANGE price TO 5.50 IN ALL ROWS OF the TABLE 
UPDATE dishes SET price = 5.50 
; CHANGE is_spicy TO 1 IN ALL ROWS OF the TABLE 
UPDATE dishes SET is_spicy = 1
```

Another example, using a column name in an UPDATE expression:

```sql
UPDATE dishes SET price = price * 2
```

All UPDATE queries we have just presented modify each row in the dishes table. To have UPDATE change only certain rows, simply add a **WHERE** clause which is a logical expression that specifies which rows to modify (in this example).
Here's the use of a WHERE clause with UPDATE:

```sql
; CHANGE the spicy STATUS OF Eggplant WITH Chili Sauce 
UPDATE dishes SET is_spicy = 1 WHERE dish_name = 'Eggplant with Chili Sauce' 
; Decrease the price OF General Tso's Chicken 
UPDATE dishes SET price = price - 1 WHERE dish_name = 'General Tso\'s Chicken'
```

Another example: I wanted to replace words in MediaWiki (source by syntaxhighlight) for updating Geshi (Syntax Highlight), here are the commands I used:

```sql
UPDATE `blocnotesinfo`.`wiki_text` SET `old_text` = REPLACE(`old_text`,"<source","<syntaxhighlight");
UPDATE `blocnotesinfo`.`wiki_text` SET `old_text` = REPLACE(`old_text`,"</source","</syntaxhighlight");
```

## DELETE - Delete Data

The DELETE command removes rows from a table:

```sql
DELETE FROM tablename [WHERE where_clause]
```

Without a WHERE clause, DELETE removes all rows from the table:

```sql
DELETE FROM dishes
```

Deleting specific rows from a table:

```sql
; DELETE ROWS IN which price IS greater than 10.00 
DELETE FROM dishes WHERE price > 10.00 
; DELETE ROWS IN which dish_name IS exactly "Walnut Bun" 
DELETE FROM dishes WHERE dish_name = 'Walnut Bun'
```

**As there is no UNDELETE command in SQL, be careful when using DELETE!**

## SELECT - Retrieve Data

* Data retrieval:

```sql
SELECT column1[, column2, column3, ...] FROM tablename
```

* Retrieving dish_name and price:

```sql
SELECT dish_name, price FROM dishes
```

The * symbol is a shortcut that refers to all columns in the table.

* Using * in a SQL query:

```sql
SELECT * FROM dishes
```

* Constraints on rows returned by SELECT:

```sql
SELECT column1[, column2, column3, ...] FROM tablename WHERE where_clause
```

* Retrieving specific dishes:

```sql
; Dishes WITH price greater than 5.00 
SELECT dish_name, price FROM dishes WHERE price > 5.00 

; Dishes whose name exactly matches "Walnut Bun" 
SELECT price FROM dishes WHERE dish_name = 'Walnut Bun' 

; Dishes WITH price more than 5.00 but less than OR equal TO 10.00 
SELECT dish_name FROM dishes WHERE price > 5.00 AND price <= 10.00 

; Dishes WITH price more than 5.00 but less than OR equal TO 10.00, 
; OR dishes whose name exactly matches "Walnut Bun" (at any price)  
SELECT dish_name, price FROM dishes WHERE price > 5.00 AND price <= 10.00 OR dish_name = 'Walnut Bun'
```

* SQL Operators for the WHERE clause:

{{< table "table-hover table-striped" >}}
| Operator | Description |
|---------|-------------|
| = | Equal to (like == in PHP) |
| <> | Not equal to (like != in PHP) |
| > | Greater than |
| < | Less than |
| >= | Greater than or equal to |
| <= | Less than or equal to |
| AND | Logical AND (like && in PHP) |
| OR | Logical OR (like \|\| in PHP) |
| () | Grouping |
{{< /table >}}

## ORDER BY and LIMIT - Data Sorting

* Sort rows returned by a SELECT query:

```sql
SELECT dish_name FROM dishes ORDER BY price
```

* To sort in descending order, add DESC after the column to sort:

```sql
SELECT dish_name FROM dishes ORDER BY price DESC
```

You can sort by multiple columns: if 2 rows have the same value for the first column, they will be sorted by the value of the second, etc.
The query in the example below sorts the rows in the dishes table by descending price; if rows have the same price, they will then be sorted alphabetically by name.

* Sorting by multiple columns:

```sql
SELECT dish_name FROM dishes ORDER BY price DESC, dish_name
```

It is sometimes convenient to retrieve only a certain number of rows: you may want to know the cheapest dish, or display only 10 results for example.

* Limiting the number of rows returned by SELECT:

```sql
SELECT * FROM dishes ORDER BY price LIMIT 1
```

* Limiting the number of rows returned by SELECT:

```sql
SELECT dish_name, price FROM dishes ORDER BY dish_name LIMIT 10
```

**As a general rule, you should only use LIMIT in queries with an ORDER BY clause**, because otherwise, the database system can return rows in any order: "the first" row of a query result can then be different from that returned by the same query at another time.

## Wildcards or SQL Regex

Wildcards or Regex characters allow you to search for strings matching certain patterns: you can find strings ending with .edu or containing @. SQL has 2 wildcard characters: _ corresponds to one and only one character, while % corresponds to a sequence of characters of any length (possibly empty). These wildcards are active when they are in strings used with the LIKE operator in a WHERE clause.

* Using wildcards with SELECT:

```sql
; Retrieve ALL ROWS IN which dish name begins WITH D 
SELECT * FROM dishes WHERE dish_name LIKE 'D%' 

; Retrieve ROWS IN which dish name IS Fried Cod, Fried Bod, 
; Fried Nod, AND so ON. 
SELECT * FROM dishes WHERE dish_name LIKE 'Fried _od'
```

Wildcards can also be used in WHERE clauses of UPDATE and DELETE commands. The query in the example below doubles the price of all dishes containing "chili" in their names:

* Using wildcards with UPDATE:

```sql
UPDATE dishes SET price = price * 2 WHERE dish_name LIKE '%chili%'
```

* Using wildcards with DELETE:

```sql
DELETE FROM dishes WHERE dish_name LIKE '%Shrimp'
```

To literally use a % or _ character with the LIKE operator, precede it with a backslash (\).
