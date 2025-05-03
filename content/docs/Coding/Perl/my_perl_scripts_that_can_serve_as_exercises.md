---
weight: 999
url: "/Mes_scripts_Perl_qui_peuvent_servir_d'exercices/"
title: "My Perl Scripts That Can Serve as Exercises"
description: "A collection of Perl scripts for practical training, including scripts for monitoring Java applications, managing SFTP accounts, watching log files, and more."
categories: ["Linux", "Development", "Database"]
date: "2012-06-06T13:08:00+02:00"
lastmod: "2012-06-06T13:08:00+02:00"
tags: ["Perl", "Development", "Scripting", "Linux", "Solaris", "Network", "Monitoring"]
toc: true
---

## Introduction

It's not always easy to start with a programming language, especially when you've never studied it in school! That's why I'm offering some small scripts I made right after going through some books. These scripts should increase in difficulty as you go along.

## PS, Java and XMX

For work, I had to create a script for Solaris that lists users, PIDs, launch dates for Java processes, and some other information. The goal was to display XMX values without the entire Java command inserted, or other unnecessary information for the people who were going to use this script. Here is my first Perl script:

```perl
#!/usr/bin/perl -w
## Java Status Script for Solaris ##
## Made by Pierre Mavro ##

use strict;

my $line;
my $ups;

print "\n############## Bridge Status ##############\n\n";
print "  USER\t   PID\t%MEM   TIME  XMX\n";

open FILE, "ps -edfo user,pid,pmem,stime,args \| grep java \| grep -v grep |";

while ($line = <FILE>){
        if ($line =~ /^(.*[_|:]\d+).*/) {
                $ups = $1; 
                if ($line =~ /.*-Xmx(\d+)/) {
                        print "$ups $1M\n";
                } else {
                        print "$ups NO XMX DEFINED\n";
                }
        }
}

close FILE;
print "\n"
```

## Creating SFTP User Accounts

This script was designed to create and delete SFTP user accounts. It then sends an email to the admin, and copies to the concerned person. There is also automatic password generation.

```perl
#!/usr/bin/perl -w
## SFTP Mangament Script ##
## Made by Pierre Mavro ##

## Needed componants
# - apg
# - Perl "Mail::Sender::Easy" module

use strict;

# Load Modules
use Mail::Sender::Easy qw(email);

die "USAGE: manage_sftp_user.pl [create|delete] username login(\@mycompany.com of the consultant mail)\n" if ($ARGV[0] !~ /(create|delete|remove)/);

my $create_account;
my $gen_pass;
my $crypted_pass;
my $user_login = $ARGV[1];
my $mail_cc = $ARGV[2];

# Functions
# User deletion
sub user_del {
    if ($ARGV[0] eq "delete") {
        system "sftp-kill", $user_login;
        system "userdel", "-r", $user_login;
        system "rm", "-Rf", "/home/clients/$user_login";
    }   
}

# User creation
sub user_add {
    if ($ARGV[1]) {
        # Check if user exist
        foreach $radius_user_files (@all_existing_radius_users) {
            open (FILER, "<$radius_user_files") || die("Can't open file: $!");
                while (<FILER>) {
                if ($_ =~ /^(\w+).+Auth-Type.+".*"/) {
                    if ($1 eq $ARGV[1]) {
                      die("Sorry but user already exist\n");
                    }
                }
            }
            close (FILER);
            # Generate password
            $password = `$apg_exec -a0 -n1 -x10 -m10 -MCN`;
            chomp $password;
            # Write to raddius files
            open (FILEW, ">>$radius_user_files") || die("Can't write file: $!");
            printf FILEW "%-15s", "$ARGV[1]";
            print FILEW "\t\tAuth-Type := Local, User-Password == \"$password\"\n";
            print FILEW "\t\t\tFilter-Id = \"vpnhome\"\n\n";
            close (FILEW);
            # Push to arrays for sending mails
            $user_pass_id = 0;
            push @tab_users, $ARGV[1];
            push @tab_pass, $password;
        }
    } else {
        &help;
    }
}

sub send_mail {
    # Just a verification to be sure that the username folder has been created
    if ( -d "/home/clients/$user_login" ) { 
        email({
            'from'         => 'it-system@mycompany.com',
            'to'           => 'it-system@mycompany.com',
            'cc'           => "$mail_cc\@mycompany.com",
            'subject'      => "New SFTP Account for $user_login",
            'priority'     => 2,
            'confirm'      => '', 
            'smtp'         => 'localhost',
            'port'         => '25',
            'auth'         => '',
            'authid'       => '',
            'authpwd'      => '',
            '_text'        => '',
            '_html'        => "Informations on configuration of an SFTP Client to connect to mycompany SFTP Server:<br /><br />
                               Host address: sftp.mycompany.com<br />
                               Port: 22<br />
                               Login: $user_login<br />
                               Password: $gen_pass<br /><br />
                               If you are on windows, we recommand this client: <a href=\"http://winscp.net\">WinSCP</a><br /><br />
                               Those informations are confidentials, please keep them safetly.",
        }) or die "email() failed: $@";
    }
}

# Launch
die "USAGE: manage_sftp_user.pl [create|delete] username login(\@mycompany.com of the consultant mail)\n" if ((@ARGV > "3") or (@ARGV < "2"));

if ($ARGV[0] eq "delete") {
    &user_del;
} elsif ($ARGV[0] eq "create") {
    die "USAGE: manage_sftp_user.pl [create|delete] username login(\@mycompany.com of the consultant mail)\n" if (@ARGV ne "3");
    &user_add;
    &send_mail;
}
```

## Controlling the Number of Users for MySecureShell

As one of the founders of MySecureShell, I had to create scripts for Nagios agents. Here's a script to monitor the number of users compared to the maximum:

```perl
#!/usr/bin/perl -w
#
# Usage: check_mss_users ([warn] [critical] are optionals)
#
# Nagios script to check MySecureShell users (can be checked only on MySecureShell server)
#
# Made by Pierre Mavro / MySecureShell Team

use strict;

#### Vars can be touched ####
my $sftp_who = "/usr/bin/sftp-who"; # sftp-who binary

#### Do not edit now ####
my $warn_mss_users = $ARGV[0] or 0.75; # warning MySecureShell users or using 75% by default
my $crit_mss_users = $ARGV[1] or 0.90; # critical MySecureShell users or using 90% by default
my $max_mss_users;
my $curr_mss_users;

sub get_infos {
        # Check config file and read config file to check LimitConnection value
        open (SFTP_WHO, "$sftp_who |") or die "Couldn't execute $sftp_who binary: $!\n";
        while (my $line = <SFTP_WHO>) {
                # Find the currently and maximum connected users
                if ($line =~ /^---.(d+).\/.(\d+).clients/) {
                        $curr_mss_users = $1; 
                        $max_mss_users = $2; 
                }
        }
        close SFTP_WHO;
}

sub check {
        if ($curr_mss_users < $warn_mss_users) {
                print "USERS OK - currently connected $curr_mss_users / $max_mss_users |users=$curr_mss_users;$warn_mss_users;$crit_mss_users\n";
                exit(0);
        } elsif ($curr_mss_users < $crit_mss_users) {
                print "USERS WARNING - currently connected $curr_mss_users / $max_mss_users\n";
                exit(1);
        } else {
                print "USERS CRITICAL - currently connected $curr_mss_users / $max_mss_users\n";
                exit(2);
        }
}

# If the twice values are defined
if ((defined($warn_mss_users)) and (defined($crit_mss_users))) {
        &get_infos;
        &check;
# If only one value is defined
} elsif ((defined($warn_mss_users)) or (defined($crit_mss_users))) {
        print "Usage: check_mss_users ([warn] [critical] are optionals)\n";
        exit(-1);
# If no values are defined, defaults are used
} else {
        &get_infos;
        &check;
}
```

## Managing Freeradius Accounts

Here is a script that manages (add/delete) radius accounts. It also allows you to change passwords, remind users that their password will change by email, and resend credentials for forgetful people.

You can also list people with their credentials in a nice table and replicate to secondary radius servers via SSH.

For Radius to work correctly afterwards, you'll need to add this to `/etc/freeradius/users`:

```
 $INCLUDE  vpn_users
```

For the usage of this script, you'll need to add the script to crontab with:

- generate: to generate new passwords in a new file and send an email to each user
- reminder: to remind users that their password will change
- switch: to use the new file with the new credentials

```perl
#!/usr/bin/perl -w
## VPN Users Management v0.3 ##
## Made by Pierre Mavro ##

## DEPENDANCIES ##
# For this working script, you need:
# - apg command
# - File::Copy perl module
# - Mail::Sender perl module
# - Net::SSH perl module
# - Net::SCP perl module

use strict;

# Load Modules
use File::Copy;
use Mail::Sender;
use Net::SSH qw(ssh);
use Net::SCP qw(scp);

# Verify syntax
unless ($ARGV[0]) {
    &help;
}

## Vars ##
my $apg_exec = "/usr/bin/apg"; # APG executable location
my $radius_user_file = "/etc/freeradius/vpn_users"; # Radius users file
my $radius_user_file_bak = "/etc/freeradius/vpn_users_bak"; # Radius backup user file
my $radius_user_file_new = "/etc/freeradius/vpn_users_new"; # Radius new users temp file
my $radius_user_file_tmp = "/etc/freeradius/vpn_users_tmp"; # Radius users file
my @all_existing_radius_users; # Get all radius users files
my $mail_fqdn = "myconpany.com"; # Will be use for your company fqdn mail
my @radius_nodes = qw(tasmania); # Radius server node list (must have ssh keys)

# Do not touch
my @tab_users_unsorted;
my @tab_users;
my @tab_pass;
my $user;
my $password;
my $number_of_users;
my $user_pass_id;
my $scp;
my $node;
my $radius_user_files;

# Add the new radius user file to array to upgrade new file too
$all_existing_radius_users[0] = $radius_user_file;
if (-f $radius_user_file_new) {
    push @all_existing_radius_users, $radius_user_file_new;
}

## Testing dependancies ##
die "Sorry but APG could not be found, please install it or change the location (actually: $apg_exec)\n" if (! -f $apg_exec);

## Functions ##

# Help
sub help {
    print "USAGE: vpn_user_management.pl [list|create|delete|generate|switch|reminder|send|replicate|help] [user_login]\n";
    print "\t- list: list all users with their password\n";
    print "\t- create: create user (eg. vpn_user_management.pl create username)\n";
    print "\t- delete: delete user (eg. vpn_user_management.pl delete username)\n";
    print "\t- reminder: resend new credentials by mail (eg. vpn_user_management.pl reminder (WARNING, this will send to all users))\n";
    print "\t- generate: generate a file with new passwords and send them by mail (eg. vpn_user_management.pl generate)\n";
    print "\t- switch: switch to new credentials (new to last generated credentials file)\n";
    print "\t- send: send credentials to one user or all users (eg. vpn_user_management.pl send [username|all])\n";
    print "\t- replicate: replicate the configuration to other nodes and restart (eg. vpn_user_management.pl replicate)\n";
    print "\t- help: print this page\n";
    exit (1);
}
```

## GC Log Analyzer

Here's a small script to monitor GC logs and alert after a certain percentage of overrun. It also warns if a Full GC occurs:

```perl
#!/usr/bin/perl -w
## GC Analyzer tool v0.1 ##
## Made by Pierre Mavro ##

## DEPENDANCIES ##
# To make this script working, you need:
# - Mail::Sender perl module

use strict;
use Mail::Sender;

## Vars ##
# Global vars
my $product="Confluence"; # Give the product name

# Mails parameters
my $mail_smtp_server="localhost"; # Set the ip or mail name server
my $mail_user_from="admin"; # Set the mail sender name
my $mail_user_to="admin"; # Set the mail receiver name
my $mail_fqdn="company.com"; # Set the FQDN for incomming mails (eg. user@

# Others vars
my $log_gc_file=$ARGV[0];
my $percent_gc_warn=$ARGV[1];
my $total_curr_gc;
my $percent_gc;
my $mail_subject;
my $problem_line;
my $mail_status;
my $curpos;

# Verifications
&help if (@ARGV < 2);

### Starting Analyze ###

## Funtions ##

# Help
sub help {
    print "USAGE: gc_analyzer.pl [LOG GC File] [Percent for warning]\n";
    print "\t- eg: gc_analyzer.pl /home/pmavro/loggc 20\n";
    print "\t- help: print this page\n";
    exit (1);
}

# Send mail
sub send_mail {
    if ($mail_status eq "warning") {
        $mail_subject="WARNING ($product): $percent_gc_warn% of GC memory has been reached!";
    } elsif ($mail_status eq "critical") {
        $mail_subject="CRITICAL ($product): A Full has proceed";
    }
        eval {
                (new Mail::Sender)
                ->OpenMultipart({
                        smtp => "$mail_smtp_server",
                        from => "$mail_user_from\@$mail_fqdn",
                        to => "$mail_user_to\@$mail_fqdn",
                        subject => "$mail_subject",
                        multipart => 'mixed',
                })
                        ->Part({ctype => 'text/html', disposition => 'NONE', msg => <<END})
<html><body>
This line has been extracted from LOG GC file ($log_gc_file) on $product product :<br /><br />
$problem_line
</body></html>
END
                        ->EndPart("multipart/alternative")
                ->Close();
        } or print "Error sending mail: $Mail::Sender::Error\n";
}

open GCFILE, "<$log_gc_file" or die "Sorry but the file wasn't found\n: $!";

# Watchdog on the file
seek(GCFILE, 0, 1);
for (;;) {
    for ($curpos = tell(GCFILE); <GCFILE>; $curpos = tell(GCFILE)) {
        chomp $_;
        $problem_line=$_;
        # If a Full GC is detected
        if ($_ =~ /.*Full GC.*/) {
            $mail_status="critical";
            &send_mail;
        # If GC, percentage should be calculated
        } elsif ($_ =~ /.+\[GC (\d+)K->(\d+)K\((\d+)K\).+/) {
            # Prepare for calcul
            $total_curr_gc=$2;
            $percent_gc=($3*$percent_gc_warn)/100;
            # Check the percentage before Full GC
            if ($total_curr_gc ge $percent_gc) {
                $mail_status="warning";
                &send_mail;
            }
        } else {
            print "Unknow line: $_\n";
        }
    }
    sleep(1);
    seek(GCFILE, $curpos, 0);
}

close ($log_gc_file);
```

## Multi-threading an MD5 Hash

Here's a simple example that can serve as a basis for a multi-threaded script. It calculates the MD5 hash of a list of files, which is rather trivial, but it's up to you to imagine what processing to perform.

What's interesting is:

- How threads are launched
- How threads lock the array, remove an element, then release it
- How threads are stopped

```perl
#!/usr/bin/perl

use strict;
use warnings;
use threads;
use threads::shared;
use Digest::MD5;

our $VERSION = 1.0;
our $SLASH = q{/};

# Getting command line arguments
my @arguments = @ARGV;
if (!$arguments[0])
{
    usage ();
}
my $directory = $arguments[0];
my $max_threads = 2;
if ($arguments[1] && $arguments[1] =~ /^\d{1,2}$/xm)
{
    $max_threads = $arguments[1];
}

# Getting regular files in the directory
my $DIR;
opendir $DIR , $directory or die "unable to opendir: $!\n";
my @files = grep { -f "$directory/$_" } readdir $DIR or die "unable to readdir: $!\n";
closedir $DIR or die "unable to closedir: $!\n";
my @tableau : shared = sort @files;

# Launching threads
my @threads;
for (1 .. $max_threads)
{
    my $thread = threads->create ('thread_code');
    push @threads , $thread;
}

# Waiting for the threads to finish
while (@threads)
{
    my $thread = shift @threads;
    $thread->join ();
}
print "fini\n" or die "unable to print: $!\n";
exit 0;


##################################


sub thread_code
{
    BOUCLE : while (1)
    {
	my $element;
	{
	    # In this block, we only manage the shared array
	    # We don't want it to be locked for too long
	    lock @tableau;
	    $element = shift @tableau;
	    if (!$element)
	    {
		last BOUCLE;
	    }
	    if ($element =~ /^\.{1,2}$/xm)
	    {
		next BOUCLE;
	    }
	}

	# We calculate the md5sum of the file, and we display it
	my $ctx = Digest::MD5->new;
	my $FILE;
	open $FILE , '<' , $directory.$SLASH.$element or next BOUCLE;
	$ctx->addfile ($FILE);
	close $FILE or next BOUCLE;
	my $digest = $ctx->hexdigest;
	print threads->self->tid () . ": $element -> $digest\n" or next BOUCLE;
    }
    return;
}

sub usage
{
    print <<"EOM" or die "unable to print: $!\n";
Usage: $0 directory [max_threads]
\tdirectory: the directory you want to scan.
\tmax_threads: the maximum number of threads you vant to run (1 .. 99). Defaults to 2.
EOM
    exit 1;
}
```

*Thanks to the GCU Squad team for this script*

## SUN Cluster Script for Java Applications

For my work, I needed to create this kind of script. I won't go into details, but it optimally manages Java applications that lack cluster functions:

```perl
#!/usr/bin/perl -w

# Copyright (c) 2008 by mycompany, Inc.  All rights reserved.
# Use is subject to license terms.

# Sun Cluster Script v1.3

# Pierre Mavro @ mycompany
# Last modification: 27/01/2009 11:45

# Supported Sun Solaris cluster Versions:
# 3.0, 3.1, 3.2

###################################
#  How to use Sun Cluster script  #
###################################

# Requierements:
# - Syslog server on localhost (on each nodes)
# - Perl 5.8 or newer should be installed on the server
# - Perl Module Unix::Syslog (use cpan to install it and Sun Studio Express to compile with cc)

# To configure the SUN Cluster mycompany Service, you'll need to add this in the configuration:
# - Start : PATH_OF_THE_SCRIPT start
# - Stop  : PATH_OF_THE_SCRIPT stop
# - Status: PATH_OF_THE_SCRIPT status

# To configure this script, edit the User section.
# If you need advanced functions, edit the advanced user section (by default you needn't)
```
