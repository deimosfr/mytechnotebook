---
weight: 999
url: "/Emails_récapitulatif_des_alertes_Nagios_en_cours/"
title: "Email Summary of Current Nagios Alerts"
description: "A script that sends an email summary of current Nagios alerts, allowing you to check server status before arriving at work."
categories: ["Monitoring", "Debian", "Linux"]
date: "2013-08-02T11:44:00+02:00"
lastmod: "2013-08-02T11:44:00+02:00"
tags: ["Servers", "Nagios", "Monitoring", "Email", "Perl", "Alerts"]
toc: true
---

{{< table "table-hover table-striped" >}}
| | |
| --- | --- |
| **Software version** | Core 3.2.1 |
| **Operating System** | Debian 6 |
| **Website** | [Nagios Website](https://www.nagios.org) |
| **Last Update** | 08/02/2013 |
{{< /table >}}

## Introduction

I've used this tool found on the internet for a long time, then further developed it with a former colleague. It allows you to receive an email with Nagios alerts. We used it to receive these emails before arriving at work to know what issues were happening. This helps you arrive either running or calmly :-)

This script parses the Nagios web interface, retrieves minimal information for good smartphone display, and sends an email. Just add it to crontab to receive this in the morning or whenever you want.

## Installation

You'll need Perl and some dependencies:

```bash
aptitude install liblwp-useragent-determined-perl libemail-sender-perl
```

## Configuration

```perl {linenos=table,hl_lines=[7,8,9,10],anchorlinenos=true}
#!/usr/bin/perl -w
# Inspired by Rob Moss, 2005-07-26, coding@mossko.com
# Created by Charles-Henri TURPIN
# Modified by Pierre Mavro

###########################################################################################

my $mailsmtp	=	'smtp.deimos.fr';    #	Fill these in!
my $mailfrom    =       'xxx@mycompany.com';  #   Coming from above
my $mailto	=	'xxx@mycompany.com';  # Set defautls receivers
my $mailsubject	=	'Morning Checks';

###########################################################################################

package Service;
sub const
{
	my ($classe, $hostname, $name, $duration, $message, $color) = @_; #la fonction reçoit comme premier paramètre le nom de la classe
   	my $this = {"hostname" => "hostname",
		   "name" => "name",
		   "duration" => "duration",
		   "message" => "message",
		   "color" => "color",};

	$this->{"hostname"} = $hostname if defined $hostname;
	$this->{"name"} = $name if defined $name;
	$this->{"duration"} = $duration if defined $duration;
	$this->{"message"} = $message if defined $message;
	$this->{"color"} = $color if defined $color;


   	bless ($this,$classe); #lie la référence à la classe
   	return $this; #on retourne la référence consacrée
}

sub DESTROY
{
	return 0;
}
1;

use strict;
use Getopt::Long;
use LWP::UserAgent;
use Mail::Sender;
$Mail::Sender::NO_X_MAILER = 1;


my @services = ();
my $mailbody	=	'';
my @finalmailbody = (
	'<html><head><style type="text/css">
    .titre {
	    color: #FFFFFF;
	    background-color: #357AB7;
	    font-size: 12;
    }
    .titre td{
	border: 2px solid #000000;
    }
    .el{
	font-size: 12;
    }
    .confluence{
	font-size: 10;
    }
    </style></head><body>'
);

my $debug		=	0;									#	Set the debug level to 1 or higher for information
my $webuser		=	'';							#	Set this to a read-only nagios user (not nagiosadmin!)
my $webpass		=	'';							#	Set this to a read-only nagios user (not nagiosadmin!)
my $full;
my $reporturl;
my @nagiosnames;


&GetOptions (
	"debug=s"	=>	\$debug,
	"help"		=>	\&help,
	"email=s"	=>	\$mailto,
	"names=s"	=>	\@nagiosnames,
	"full"		=>	\$full
);

if(!defined($nagiosnames[0]))
{
    @nagiosnames = ('nagios-prod');
}
program(@nagiosnames);
################################ FUNCTIONS #######################################
sub program
{
    my @nagiosnames = @_;

    foreach my $nagiosname (@nagiosnames)
    {
	launch("$nagiosname");
	main();
	check();
    }
    system("touch /tmp/nagios-report-htmlout.html; chmod 666 /tmp/nagios-report-htmlout.html");

    push @finalmailbody, '</body></html>';
    sendmail();
    system("rm /tmp/nagios-report-htmlout.html");
}
sub launch
{
    my $nagiosname = shift;
    my $title = "<p align=\"center\"><b>".uc($nagiosname)."</b></p><hr color=\"black\" align=\"center\" width=\"80%\"/>";
    push @finalmailbody, $title;
    #-- 20110926 - Match the monitoring screen's request
    #$reporturl = "http://$nagiosname/cgi-bin/nagios3/status.cgi?host=all&servicestatustypes=28&hoststatustypes=3&serviceprops=42&sorttype=1&sortoption=6&noheader";
    if($nagiosname =~ /internal/)
    {
        $reporturl = "http://$nagiosname/cgi-bin/nagios3/status.cgi?hostgroup=prod-srv&style=detail&servicestatustypes=28&serviceprops=8&sorttype=1&sortoption=6&noheader";
    }
    else
    {
        $reporturl = "http://$nagiosname/cgi-bin/nagios3/status.cgi?style=detail&servicestatustypes=28&serviceprops=8&sorttype=1&sortoption=6&noheader";
    }
	# Get all status
	if ($full)
	{
		$reporturl="http://$nagiosname/cgi-bin/nagios3/status.cgi?host=all&sorttype=2&sortoption=3";
	}
}
###############################################################################
sub main
{
    debug(1,"reporturl: [$reporturl]");

    $mailbody = http_request($reporturl);

    open(FILE, "> /tmp/nagios-report-htmlout.html") or warn "can't open file /tmp/nagios-report-htmlout.html: $!\n";

    #print "DEBUG $mailbody\n\n\n\n";
    print FILE $mailbody;
    close FILE;
}
###############################################################################
sub check
{
    open(FILE, "</tmp/nagios-report-htmlout.html");
    my $gotit = 0; # = 1 while having a service
    my $service;
    while(<FILE>)
    {
	chomp $_;
	my $line = $_;

	if($line =~ /\&host=(.*?)&service=(.*?)\'/)
	{
	    #print "DEBUG ENTERS!!!\n";
	    $service = Service->const();
	    $service->{hostname} = $1;
	    $service->{name} = $2;
	    $service->{name} =~ s/\+/ /g;
	    $service->{name} =~ s/\%2F/\//g;
	    $service->{name} =~ s/\%3A/:/g;
	}
	elsif($line =~ /CLASS=\'statusUNKNOWN/)
	{
	    #print "DEBUG UNKNOWN!!\n";
		$service->{color} = "FFDA9F";
		$gotit =1;
	}
	elsif($line =~ /CLASS=\'statusWARNING/)
	{
	    #print "DEBUG WARNING!!\n";
		$service->{color} = "FEFFC1";
		$gotit =1;
	}
	elsif($line =~ /CLASS=\'statusOK/)
	{
		if ($full)
		{
	    	#print "DEBUG WARNING!!\n";
			$service->{color} = "99FF99";
			$gotit =1;
		}
	}
	elsif($line =~ /CLASS=\'statusCRITICAL/)
	{
	    #print "DEBUG CRITICAL!!\n";
		$service->{color} = "FFBBBB";
		$gotit =1;
	}
	elsif($line =~ /\'center\'\&gt;(.*?)\&lt;/)
	{
	    #print "DEBUG MESSAGE!!\n";
		$service->{message} = $1;
	}
	elsif($line =~ /\&lt;\/TR\&gt;/)
	{
	    if($gotit == 1)
	    {
		#print "DEBUG PUSH! $service->{color}, $service->{name}, $service->{message}, $service->{hostname}, $service->{duration}\n";
		push(@services, $service);
		$gotit = 0;
	    }
	}
	elsif($line =~ /\&gt;\s*(\d+d\s*\d+h\s*\d+m\s*\d+s)/)
	{
	    #print "DEBUG DURATION!!\n";
		$service->{duration} = $1;
	}
    }
    close(FILE);

    #### mailbody creation ####

    push @finalmailbody, '<table width="100%"><tr class="titre"><td >Hostname</td><td>Service</td><td>Duration</td><td>Status</td></tr>';
    foreach my $el (@services)
    {
	#print "DEBUG  HAD A SERVICE!\n";
	push @finalmailbody, "<tr class=\"el\" bgcolor=\"#$el->{color}\"><td>$el->{hostname}</td><td bgcolor=\"#$el->{color}\">$el->{name}</td><td>$el->{duration}</td><td>$el->{message}</td></tr>\n";
    }
    push @finalmailbody, '</table>';

    $mailbody = ''; #Resetting values for next times
    @services = (); #Resetting values for next times
}

###############################################################################
sub help {
print <<_END_;

Nagios web->email reporter program.

$0 <args>

--help
	This screen

--email=<email>
	Send to this address instead of the default address
	"xxx@mycompany.com"

-names
    Names of the nagios servers you want to check (multiple names is ok)

-full
	Do not send only unhandled alerts but all (even ok)

_END_

exit 1;

}

###############################################################################
sub http_request {
	my $ua;
	my $req;
	my $res;

	my $geturl = shift;
	if (not defined $geturl or $geturl eq "") {
		warn "No URL defined for http_request\n";
		return 0;
	}
	$ua = LWP::UserAgent->new;
	$ua->agent("Nagios Report Generator " . $ua->agent);
	$req = HTTP::Request->new(GET => $geturl);
	$req->authorization_basic($webuser, $webpass);
	$req->header(	'Accept'			=>	'text/html',
				);

	# send request
	$res = $ua->request($req);

	# check the outcome
	if ($res->is_success) {
		debug(1,"Retreived URL successfully");
		return $res->decoded_content;
	}
	else {
		print "Error: " . $res->status_line . "\n";
		return 0;
	}
}

###############################################################################
sub debug {
	my ($lvl,$msg) = @_;
	if ( defined $debug and $lvl <= $debug ) {
		chomp($msg);
		print localtime(time) .": $msg\n";
	}
	return 1;
}

###############################################################################
sub sendmail {
	my $sender = new Mail::Sender {
		smtp => "$mailsmtp",
		from => "$mailfrom",
		on_errors => 'die',
	};
	$sender->Open({
		to => "$mailto",
	    subject => $mailsubject,
	    ctype => "text/html",
	    encoding => "quoted-printable"
	}) or die $Mail::Sender::Error,"\n";

	my $body_line;
	foreach $body_line (@finalmailbody)
	{
		$sender->SendEnc($body_line);
	};
	$sender->Close();
}

###############################################################################
```

## References

Add execution rights to the script.

## Usage

For usage, it's relatively simple. You need to fill in the information in the script if you want default values. Then you can specify these same options as parameters or not:

```bash
/usr/local/bin/morningchecks.pl -names nagios-prod
/usr/local/bin/morningchecks.pl -full -names nagios-prod --email=xxx@mycompany.com
```

Also check the help for more information:

```bash
> /usr/local/bin/morningchecks.pl --help

Nagios web->email reporter program.

/usr/local/bin/morningchecks.pl <args>

--help
	This screen

--email=<email>
	Send to this address instead of the default address
	"xxx@mycompany.com"

-names
    Names of the nagios servers you want to check (multiple names is ok)

-full
	Do not send only unhandled alerts but all (even ok)
```
