---
weight: 999
url: "/RRDtool_\\:_créer_ses_propre_graphiques_avec_RRDtool/"
title: "RRDtool: Create Your Own Graphics with RRDtool"
description: "Learn how to create, manage and generate graphs with RRDtool to visualize your data such as disk usage, temperature, and more."
categories: ["Database", "Linux", "Monitoring"]
date: "2010-01-21T12:44:00+02:00"
lastmod: "2010-01-21T12:44:00+02:00"
tags: ["RRDtool", "Monitoring", "Graphs", "Performance", "Database"]
toc: true
---

![RRDtool Logo](/images/rrdtool-logo.avif)

## Introduction

[RRDtool](https://fr.wikipedia.org/wiki/RRdTool) is a Round-Robin Database (RRD) management tool created by Tobi Oetiker. It is used by many open source tools, such as Cacti, collectd, Lighttpd, and Nagios, for saving cyclical data and plotting chronological data graphs. This tool was created to monitor server data, such as bandwidth and CPU temperature. The main advantage of an RRD database is its fixed size.

RRDTool also includes a tool to graphically represent the data contained in the database.
RRDTool is free software distributed under the terms of the GNU GPL.

I had to use RRD to graph disk usage of several user folders (e.g., `/home/users/*`).

- Overview of how RRDtool works:

1. Create an empty RRD database that will contain the data to graph
2. Update RRD data using the "rrdtool update" command or via a script
3. Generate graphs with the "rrdtool graph" command

## Managing an RRD Database

The RRD database uses a defined number of records. Each addition is placed at the head of the database and the others gradually shift, and so on. There can therefore be no overflow since it's controlled. The only issue with all this is that you need to know the number of records you want to keep. For example, you may want a day or a month, which are not the same.

### Creating a Database

First, let's create the database:

```bash
> rrdtool create temptrax.rrd \
--start N --step 300 \
DS:probe1-temp:GAUGE:600:55:95 \
DS:probe2-temp:GAUGE:600:55:95 \
DS:probe3-temp:GAUGE:600:55:95 \
DS:probe4-temp:GAUGE:600:55:95 \
RRA:MIN:0.5:12:1440 \
RRA:MAX:0.5:12:1440 \
RRA:AVERAGE:0.5:1:1440
```

With the create argument, RRDtool will create the database that will contain all the necessary fields. This database doesn't contain any data yet. Then:

- temptrax.rrd: this is the name of the database and its location
- --start N: gives an indication of when the graph starts. Here I use N to say now. But I can use a date in epoch format.
- --step 300: indicates the interval time in seconds when data will arrive in the database (here 5 min)

---

- DS: specifies the different data sources. Here I have 4 temperature probes followed by the name (DS-Name) I want to assign them.
  - GAUGE: This is a DST (Data Source Type). There are several of them:

{{< table "table-hover table-striped" >}}
| DST Type | Description |
|----------|-------------|
| GAUGE | This is the most common, and generally the best choice |
| COUNTER | This is a counter that will increment continuously |
| DERIVE | Will record the drift of the previous and next values |
| ABSOLUTE | Records values and resets them after each reading |
{{< /table >}}

- 600:55:95: These last three fields mean:
  - 600: The minimum heartbeat in seconds (after this delay, the value will become unknown if the database has not received anything during this period)
  - 55: Minimum possible value (outside of which the value will be unknown)
  - 95: Maximum possible value (outside of which the value will be unknown)

For these two values above, if you don't know what to put, use 'U' for unknown (e.g., DS:probe1-temp:GAUGE:600:U:U)

---

- RRA: RRA stands for Round Robin Archives. These are like views in which data will be stored. In each RRD database, RRAs are stored separately with a defined number of records. With each new record in the database, a PDP (Primary Data Point) is added which will be combined with it and placed in our RRA in a CF (Consolidation Function). It will determine the current value to write.
  - MIN: This is the type of CF we use. There are others such as:

{{< table "table-hover table-striped" >}}
| CF Type |
|---------|
| AVERAGE |
| MIN |
| MAX |
| LAST |
{{< /table >}}

- 0.5: This is an XFF (XFiles Factor) which is a percentage of PDPs that can be unknown without receiving unknown values.
- 12: This is the number of PDPs that will make up the recorded value.
- 1440: This is the number of records that the RRA should contain.

---

To summarize: I create an RRD database called temptrax.rrd which will:

- start now
- be updated every 5 minutes

Additionally:

- I have 4 different data sources with 1 probe of type GAUGE.
- If my data source is not updated **at least** every 10 minutes and the value is not between 55 and 95, then the value will be unknown.

I also have 3 types of RRA:

- 2 are for min and max values using 12 PDPs allowing 50% of them to be unknown
- We can record up to 1440 records

Knowing that we normally update every 5 min, and we use 12 PDPs (each update is a PDP). This means that:

- we add one record every hour (5 mins \* 12)
- we make 1440 records

This gives us 60 days (1440/24h) of RRA. Each min and max value collected in PDP will be used as a value for the RRA.

In the last RRA, we use an average of collected PDPs, allowing 50% unknown, but we use only one PDP per record (so each update). We make 1440 records which means this RRA will record (1440/12 updates per hour/24h) 5 days of data.
In a future case, we will see that with a simple PDP, the CF is not very important and you will probably have used LAST.

## Updating a Database

Updating can be very simple or very complicated depending on the type of graphs. For usage, however, it is very simple and takes this form:

```bash
rrdtool update <file.rrd> timestamp:val1[:val2:...]
```

Generally this command is placed at the end of a script that has retrieved all the data via SNMP or other means to pass a return result as an argument. If we break down the command, it gives us:

- \<file.rrd\>: the location with the rrd file
- timestamp: the update time (as above, "N" means now)
- val1[:val2:...]: the values are separated by ":" and **there must be as many as declared DS and in order!**

Here is a crude example to update data:

```perl
#!/usr/bin/perl
# RRD Update Script: update_rrd_temps.pl

$HOST = "10.10.0.90";
$PATH = "/home/benr/RRD/TempTrax-RRD";
$NumProbes = 4;

for($i=1; $i <= $NumProbes; $i++) {
  $x = `${PATH}/check_temptraxe -H ${HOST} -p ${i}`;
  $x =~ s/^.*([0-9.]{4}).*$/$1/;
  chomp($x);
  push(@TEMPS,$x);
}

`/usr/local/rrdtool-1.0.48/bin/rrdtool update ${PATH}/temptrax.rrd
 "N:$TEMPS[0]:$TEMPS[1]:$TEMPS[2]:$TEMPS[3]"`;
```

Additionally, you will need to set up this kind of script in crontab to automatically update the data. Otherwise, you can create a loop to update (a bit cruder).

### Modifying the Data Insertion Order

If you want to change the order of data insertion, you can do this:

```bash
rrdtool update --template ds2:ds1:ds3
```

However, remember that data already inserted has been entered in a specific order, if you change it, the graph history will be incorrect.

# Generating Graphs

We will use rrdtool graph to generate a graph from our rrd database. There are tons of options for this command. We will only see the basics here. Here's what we're going to run:

```bash
> rrdtool graph mygraph.png -a PNG --title="TempTrax" \
--vertical-label "Deg F" \
'DEF:probe1=temptrax.rrd:probe1-temp:AVERAGE' \
'DEF:probe2=temptrax.rrd:probe2-temp:AVERAGE' \
'DEF:probe3=temptrax.rrd:probe3-temp:AVERAGE' \
'DEF:probe4=temptrax.rrd:probe4-temp:AVERAGE' \
'LINE1:probe1#ff0000:Switch Probe' \
'LINE1:probe2#0400ff:Server Probe' \
'AREA:probe3#cccccc:HVAC' \
'LINE1:probe4#35b73d:QA Lab Probe' \
'GPRINT:probe1:LAST:Switch Side Last Temp\: %2.1lf F' \
'GPRINT:probe3:LAST:HVAC Output Last Temp\: %2.1lf F\j' \
'GPRINT:probe2:LAST:Server Side Last Temp\: %2.1lf F' \
'GPRINT:probe4:LAST:QA Lab Last Temp\: %2.1lf F\j'
```

To eventually get this:

![Temptraxrrd](/images/temptraxrrd.avif)

Let's break down what this does:

- mygraph.png: path and name of the graph to generate
  - -a PNG: type of image file to generate. By default it's gif, but you can force PNG or GD.
  - --title: The title to display at the top of the graph
  - --vertical-label: Name to give the Y axis

---

- DEF: These are just virtual names (vname) that we give to DS. It is strongly recommended to use vnames! Indeed, since we can use multiple RRD databases to make a graph, if we have DS with the same name, there will be undesirable effects.
  - probe1: name of the vname
  - temptrax.rrd: rrd database to process
  - rrd:probe1-temp: the name of the original DS
  - AVERAGE: The type of CF we are interested in specified in the RRA

---

- LINE1|AREA: These are the type of graph we want to make. For more examples [click here](#different-types-of-graphs) to see available graphs.
  - probe4: the corresponding vname
  - #35b73d: Curve color in hexadecimal format
  - QA Lab Probe: Legend name displayed at the bottom of the graph

---

- GPRINT: These lines allow additional information to be placed at the bottom of the graph. It is generally nice to see the
  - probe1: the corresponding vname
  - LAST: This is the last CF data in the database (because I use LAST)
  - Switch Side Last Temp\:: This is the line to be displayed at the bottom of the graph
    - %2.1lf: used to display a numeric value with 1 decimal place (see the documentation for printf or sprintf for all possibilities)
    - \j: allows right alignment

### Different Types of Graphs

- LINE type graph:

![Graph-line](/images/graph-line.avif)

You can specify several types of lines (LINE1, LINE2, LINE3, LINE4). The higher this number, the more the corresponding line will be above lines of lower numbers.

- AREA type graph:

![Graph-area](/images/graph-area.avif)

AREA allows you to fill the lower part of the graphs.

- STACK type graph:

![Graph-stack](/images/graph-stack.avif)

Allows you to stack the graphs.

## Generating Historical Graphs

The concept is quite simple. Our initial RRD database must collect sufficient time for us to make the graphs we are interested in (e.g., we need at least 3 weeks in the database if we want to make graphs over 3 weeks of history).

However, we will use 2 new arguments: --start and --end. By default, graphs are made over 24 hours. By adding these parameters, we can specify a start and end:

- --end: by default it's now. And it is generally practical to leave the default.
- --start: This is the start date of the graph which can be specified in several formats:
  - epoch: you can specify the date in seconds since January 1, 1970
  - days: you can specify a day of the week "monday" for example
  - weeks: you can ask for the last 2 weeks "-2weeks"
  - month: you can ask for last month "-1month"
  - year: And finally, the year "-1year"

In short, the notation is relatively easy as you can see. Check the man page for rrdfetch if you want more info.

## Example

Here's an example I propose of something I developed to have graphs on the size occupied by users on their home directory (this is very similar to [my documentation on OpenChart](./open_flash_chart_:_créer_des_graphiques_flash.html)). This will give me something like this:

![Rrd day](/images/rrd_day.avif)
![Rrd week](/images/rrd_week.avif)

### Creating the RRD Database

All my scripts are stored in a folder `/etc/scripts/`. So if you want to play copy/paste, create this folder.

So we'll create the database like this:

```bash
#!/bin/sh

rrdtool create /etc/scripts/rrd_db.rrd \
--start N --step 300 \
DS:user1:GAUGE:750:U:U \
DS:user2:GAUGE:750:U:U \
DS:user3:GAUGE:750:U:U \
DS:user4:GAUGE:750:U:U \
RRA:MIN:0.5:12:720 \
RRA:MAX:0.5:12:720 \
RRA:AVERAGE:0.5:12:720
```

Adapt the part of user names and RRD database (`/stats/rrd_db.rrd`) to what you want.

### Generating a Data File

Next we'll make a small script that will generate a data file that can be picked up by OpenChart and RRD (well, something a bit generic and easy to parse).

Again, adapt the users, user colors, source and destination:

```perl
#!/usr/bin/perl -w

use strict;
no strict "refs";

# Set folder containing all users folder
my $source='/home';
my $destination='/var/www/.stats/datas';

# Set all users
my @users = (
            'user1',
            'user2',
            'user3',
            'user4'
);

# Set color per user
my @colors = (
            '#000000',
            '#2d9409',
            '#4e5fff',
            '#de1000'
);

# Size per user
my @ksize;
my @printable_size;

# Get total size
sub get_total_size
{
    my $good_size;

    open (GET_TOTAL_SIZE, "df -k $source |") or die ("Can't get full size of $source\n");
    while (<GET_TOTAL_SIZE>)
    {
        if (/\S+\s*(\d*)\s*(\d*).*$source/)
        {
            push @users, 'disponible';
            push @colors, '#ffffff';
            return ($1, $2);
        }
    }
    close (GET_TOTAL_SIZE);
}

# Calcul userspace size for each users
sub get_users_size
{
    my $good_size;
    my $total_size=shift;
    my $ref_users=shift;
    my @users=@$ref_users;

    foreach (@users)
    {
        if (-d "$source/$_")
        {
            open (GET_SIZE, "du -sk $source/$_ |");
            while (<GET_SIZE>)
            {
                chomp $_;
                if (/(\d*)\s*(\w*)/)
                {
                    if (($1 / 1024) >= 1000)
                    {
                        # Set Go
                        $good_size = sprintf ("%.1fGo", ($1 / 1048576));
                        push @printable_size, $good_size;
                        push @ksize, $1;
                    }
                    else
                    {
                        # Set Mo
                        $good_size = sprintf ("%.0fMo", ($1 / 1024));
                        push @printable_size, $good_size;
                        push @ksize, $1;
                    }
                }
                else
                {
                    $good_size = sprintf ("%.0f", ($1 / $total_size) * 100);
                    push @ksize, $1;
                }
            }
            close (GET_SIZE);
        }
        else
        {
            push @ksize, 0;
            push @printable_size, '0Mo';
        }
    }
}

sub get_free_size
{
    my $total_size=shift;
    my $busy_size=shift;
    my $ref_ksize=shift;
    my @ksize=@$ref_ksize;
    my $used_ksize=0;

    sub print_size
    {
        my $good_size;
        if (($_[0] / 1024) >= 1000)
        {
            # Set Go
            $good_size = sprintf ("%.1fGo", ($_[0] / 1048576));
            push @printable_size, $good_size;
        }
        else
        {
            # Set Mo
            $good_size = sprintf ("%.0fMo", ($_[0] / 1024));
            push @printable_size, $good_size;
        }
    }

    foreach (@ksize)
    {
        $used_ksize += $_;
    }

    # Delete free space values
    pop @ksize;
    pop @printable_size;

    # Add others if needed
    my $unknow_ksize = ($busy_size - $used_ksize) / 1024;
    if ($unknow_ksize > 1024)
    {
        push @users, 'autres';
        $unknow_ksize = sprintf ("%.0f", $unknow_ksize);

        # Free size
        my $free_size = ($total_size - $busy_size);
        push @ksize, $free_size;
        &print_size($free_size);

        # Other size
        push @ksize, $unknow_ksize;
        &print_size($unknow_ksize);

        push @colors, '#ff6600';
    }
    return @ksize;
}

# Write datas
sub write_datas_file
    {
    my $ref_users=shift;
    my $ref_ksize=shift;
    my $ref_printable_size=shift;
    my $ref_colors=shift;
    my @users=@$ref_users;
    my @ksize=@$ref_ksize;
    my @printable_size=@$ref_printable_size;
    my @colors=@$ref_colors;

    open (DATAW, ">$destination") or die "Can't write $destination file : $!
    my $i=0;
    foreach (@users)
    {
        print DATAW "$users[$i]:$ksize[$i]:$printable_size[$i]:$colors[$i]\n
        $i++;
    }
    close (DATAW);
}

# Get total partition size and busy size
my ($total_size, $busy_size)=&get_total_size;
# Get all userspace size
&get_users_size($total_size,\@users);
# Get free size and unknow size
@ksize=&get_free_size($total_size,$busy_size,\@ksize);
# Write datas file
&write_datas_file(\@users,\@ksize,\@printable_size,\@colors);
```

### Updating the Database

Now for updating the database, we'll need to put a script in crontab to update it:

```perl
#!/usr/bin/perl -w

use strict;

my $source='/var/www/.stats/datas';
my $destination='/var/www/.stats/rrd_db.rrd';
my @rrd_arg;
my $gb_size;

open (RRD_DATA, "<$source");
    while (<RRD_DATA>)
    {
        chomp $_;
        unless (/disponible|autre/i)
        {
            if (/^.*:(\d*):/)
            {
                $gb_size=sprintf ("%.1f", ($1 / 1048576));
                push @rrd_arg, $gb_size;
            }
        }
    }
close (RRD_DATA);

@rrd_arg = join ':', @rrd_arg;
system "rrdtool update $destination N:@rrd_arg\n";
```

Adapt the source and destination here.

### Generating Graphs

To generate the graphs, we'll need to define where the images will be stored. For this, adapt the source, rrd_db and the 2 destinations. There will be one image for the day graph and one image for the 2-week graph:

```perl
#!/usr/bin/perl -w

use strict;

my $source='/var/www/.stats/datas';
my $rrd_db='/var/www/.stats/rrd_db.rrd';
my $destination_day='/var/www/.stats/usage_day.png';
my $destination_week='/var/www/.stats/usage_week.png';
my @users_lines;
my @colors_lines;
my @user_size;

open (RRD_DATA, "<$source");
    my $i=1;
    while (<RRD_DATA>)
    {
        chomp $_;
        unless (/disponible|autres/i)
        {
            if (/(.*):(.*):(.*):(.*)/
)            {
                push @users_lines, "'DEF:user$i=$rrd_db:$1:AVERAGE'";
                push @colors_lines, "'LINE1:user$i$4:$1 ($3)'";
                $i++;
            }
        }
    }
close (RRD_DATA);

# Generate day graph
system "rrdtool graph $destination_day -a PNG --title=\"Utilisation disque par jour\" --vertical-label \"Giga Octets\" @users_lines @colors_lines";
# Generate 2 weeks graphs
system "rrdtool graph $destination_week --start -2weeks --end N -a PNG --title=\"Utilisation disque (2 semaines)\" --vertical-label \"Giga Octets\" @users_lines @colors_lines";
```

Now, we have everything we need to create, update and generate graphs. Let's just set the proper execution rights:

```bash
chmod u+rx /etc/scripts/*
```

And finish by defining the crontab:

```bash
# Generate graphs
*/5 * * * * /etc/scripts/gen-piechart-db.pl 2>/dev/null
*/5 * * * * /etc/scripts/gen-datas.pl ; /etc/scripts/update-rrd.pl ; /etc/scripts/gen-rrd-graph.pl 1> /dev/null
```

## Resources
- http://www.cuddletech.com/articles/rrd/index.html
- http://oss.oetiker.ch/rrdtool/
- http://www.vandenbogaerdt.nl/rrdtool/tutorial/graph.php
