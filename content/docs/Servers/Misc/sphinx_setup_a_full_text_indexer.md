---
weight: 999
url: "/Sphinx_\\:_setup_a_full_text_indexer/"
title: "Sphinx: Setup a Full Text Indexer"
description: "Guide on setting up Sphinx, an open source full text search server, with specific instructions for MediaWiki integration and configuration."
categories: ["Debian", "Storage", "Database"]
date: "2013-08-02T14:58:00+02:00"
lastmod: "2013-08-02T14:58:00+02:00"
tags: ["search", "sphinx", "indexing", "mediawiki", "database", "configuration"]
toc: true
---

![Sphinx search](/images/sphinx_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.0.4 |
| **Operating System** | Debian 7 |
| **Website** | [Sphinx Search Website](https://sphinxsearch.com/) |
| **Last Update** | 02/08/2013 |
{{< /table >}}

## Introduction

Sphinx is an open source full text search server, designed from the ground up with performance, relevance (aka search quality), and integration simplicity in mind. It's written in C++ and works on Linux (RedHat, Ubuntu, etc), Windows, MacOS, Solaris, FreeBSD, and a few other systems.

Sphinx lets you either batch index and search data stored in an SQL database, NoSQL storage, or just files quickly and easily — or index and search data on the fly, working with Sphinx pretty much as with a database server.
A variety of text processing features enable fine-tuning Sphinx for your particular application requirements, and a number of relevance functions ensures you can tweak search quality as well.
Searching via SphinxAPI is as simple as 3 lines of code, and querying via SphinxQL is even simpler, with search queries expressed in good old SQL.

Sphinx clusters scale up to tens of billions of documents and hundreds of millions search queries per day, powering top websites such as Craigslist, Living Social, MetaCafe and Groupon... to view a complete list of known users please visit our Powered-by page.
And last but not least, it's open-sourced under GPLv2, and the community edition is free to use.

## Installation

To install Sphinx on Debian, this is simple:

```bash
aptitude install sphinxsearch
```

## Configuration

First of all, you need to setup the application configuration for Sphinx before starting the indexation.

### MediaWiki

#### Install plugin

First of all, we need to get the MediaWiki plugin[^1] to get the sphinx config with and deploy:

```bash
cd /var/www/mediawiki/extension
git clone https://gerrit.wikimedia.org/r/p/mediawiki/extensions/SphinxSearch.git
cd SphinxSearch
cp sphinx.conf /etc/sphinxsearch/
```

You will also need to get the PHP client API. The biggest problem here is to find the one that match your Sphinx version. I had several problem regarding that point and I strongly suggest that you take time to choose the best corresponding one for your version: http://code.google.com/p/sphinxsearch/source/browse/#svn%2Ftags.

For Debian 7 and the version installed with, you need to take this one:

```bash
wget http://sphinxsearch.googlecode.com/svn/tags/REL_1_10/api/sphinxapi.php
```

You also can take the logo to add it just next to the Mediawiki logo (this is optional):

```bash
mkdir -p skins/images
wget -O skins/images/Powered_by_sphinx.png http://upload.wikimedia.org/wikipedia/mediawiki/8/8e/Powered_by_sphinx.png
```

#### Configure

Edit the configuration file and fulfill the informations with your mediawiki database (`/etc/sphinxsearch/sphinx.conf`):

```bash
#
# Sphinx configuration for MediaWiki
#
# Based on examples by Paul Grinberg at http://www.mediawiki.org/wiki/Extension:SphinxSearch
# and Hank at http://www.ralree.info/2007/9/15/fulltext-indexing-wikipedia-with-sphinx
# Modified by Svemir Brkic for http://www.newworldencyclopedia.org/
#
# Released under GNU General Public License (see http://www.fsf.org/licenses/gpl.html)
#
# Latest version available at http://www.mediawiki.org/wiki/Extension:SphinxSearch

# data source definition for the main index
source src_wiki_main
{
	# data source
	type		= mysql
	sql_host	= localhost
	sql_db		= wikidb
	sql_user	= user
	sql_pass	= password
	# these two are optional
	#sql_port	= 3306
	#sql_sock	= /var/lib/mysql/mysql.sock

	# pre-query, executed before the main fetch query
	sql_query_pre	= SET NAMES utf8

	# main document fetch query - change the table names if you are using a prefix
	sql_query	= SELECT page_id, page_title, page_namespace, page_is_redirect, old_id, old_text FROM wiki_page, wiki_revision, wiki_text WHERE rev_id=page_latest AND old_id=rev_text_id

	# attribute columns
	sql_attr_uint	= page_namespace
	sql_attr_uint	= page_is_redirect
	sql_attr_uint	= old_id

	# collect all category ids for category filtering
	sql_attr_multi  = uint category from query; SELECT cl_from, page_id AS category FROM wiki_categorylinks, wiki_page WHERE page_title=cl_to AND page_namespace=14

	# used by command-line search utility to display document information
	sql_query_info	= SELECT page_title, page_namespace FROM wiki_page WHERE page_id=$id
}

# data source definition for the incremental index
source src_wiki_incremental : src_wiki_main
{
	# adjust this query based on the time you run the full index
	# in this case, full index runs at 7 AM UTC
	sql_query	= SELECT page_id, page_title, page_namespace, page_is_redirect, old_id, old_text FROM wiki_page, wiki_revision, wiki_text WHERE rev_id=page_latest AND old_id=rev_text_id AND page_touched>=DATE_FORMAT(CURDATE(), '%Y%m%d070000')

	# all other parameters are copied from the parent source
}

# main index definition
index wiki_main
{
	# which document source to index
	source		= src_wiki_main

	# this is path and index file name without extension
	# you may need to change this path or create this folder
	path		= /var/lib/sphinxsearch/data/wiki_main

	# docinfo (ie. per-document attribute values) storage strategy
	docinfo		= extern

	# morphology (comment it if your wiki is not full english)
	# morphology	= stem_en

	# stopwords file
	#stopwords	= /var/data/sphinx/stopwords.txt

	# minimum word length
	min_word_len	= 1

	# allow wildcard (*) searches
	min_infix_len = 1
	enable_star = 1

	# charset encoding type
	charset_type	= utf-8

	# charset definition and case folding rules "table"
	charset_table	= 0..9, A..Z->a..z, a..z, \
		U+C0->a, U+C1->a, U+C2->a, U+C3->a, U+C4->a, U+C5->a, U+C6->a, \
		U+C7->c,U+E7->c, U+C8->e, U+C9->e, U+CA->e, U+CB->e, U+CC->i, \
		U+CD->i, U+CE->i, U+CF->i, U+D0->d, U+D1->n, U+D2->o, U+D3->o, \
		U+D4->o, U+D5->o, U+D6->o, U+D8->o, U+D9->u, U+DA->u, U+DB->u, \
		U+DC->u, U+DD->y, U+DE->t, U+DF->s, \
		U+E0->a, U+E1->a, U+E2->a, U+E3->a, U+E4->a, U+E5->a, U+E6->a, \
		U+E7->c,U+E7->c, U+E8->e, U+E9->e, U+EA->e, U+EB->e, U+EC->i, \
		U+ED->i, U+EE->i, U+EF->i, U+F0->d, U+F1->n, U+F2->o, U+F3->o, \
		U+F4->o, U+F5->o, U+F6->o, U+F8->o, U+F9->u, U+FA->u, U+FB->u, \
		U+FC->u, U+FD->y, U+FE->t, U+FF->s,

}

# incremental index definition
index wiki_incremental : wiki_main
{
	path		= /var/lib/sphinxsearch/data/wiki_incremental
	source		= src_wiki_incremental
}


# indexer settings
indexer
{
	# memory limit (default is 32M)
	mem_limit	= 64M
}

# searchd settings
searchd
{
	# IP address and port on which search daemon will bind and accept
	listen		= 127.0.0.1:9312

	# searchd run info is logged here - create or change the folder
	log		= /var/log/sphinxsearch/searchd.log

	# all the search queries are logged here
	query_log	= /var/log/sphinxsearch/query.log

	# client read timeout, seconds
	read_timeout	= 5

	# maximum amount of children to fork
	max_children	= 30

	# a file which will contain searchd process ID
	pid_file	= /var/run/sphinxsearch/searchd.pid

	# maximum amount of matches this daemon would ever retrieve
	# from each index and serve to client
	max_matches	= 1000

        # Remove warning of deprecated function
        compat_sphinxql_magics = 0
}

# --eof--
```

Others infos you need to know:

- Adapt the tables if you use a prefix (like you can see here with 'wiki\_') on the SQL requests
- I've also modified all paths to match in Debian's ones
- All highlighted lines are important, I've added a comment on each that needed to bring additional informations

Now you should perform an [indexation of the wiki](#indexation).

Then we are going to configure the MediaWiki plugin. add those lines to your LocalSettings.php:

```php
# Sphinx search
$wgEnableMWSuggest = true;
$wgSearchType = 'SphinxMWSearch';
$wgFooterIcons['poweredby']['sphinxsearch'] = array(
        'src' => "$wgScriptPath/extensions/SphinxSearch/skins/images/Powered_by_sphinx.png",
        'url' => 'http://www.mediawiki.org/wiki/Extension:SphinxSearch',
        'alt' => 'Search Powered by Sphinx',
);
require_once( "$IP/extensions/SphinxSearch/SphinxSearch.php" );
```

Then you'll be able to make search with sphinx :-)

### Sphinx default

To permit to the deamon to start on boot, simply edit that file and change 'no' to 'yes' (`/etc/default/sphinxsearch`):

```bash
# Settings for the sphinxsearch searchd daemon
# Please read /usr/share/doc/sphinxsearch/README.Debian for details.
#

# Should sphinxsearch run automatically on startup? (default: no)
# Before doing this you might want to modify /etc/sphinxsearch/sphinx.conf
# so that it works for you.
START=yes
```

## Indexation

### Index

You need to create a first indexation once you've configured your application. To prepare sphinx to search:

```bash
> indexer --config /etc/sphinxsearch/sphinx.conf --all
Sphinx 2.0.4-release (r3135)
Copyright (c) 2001-2012, Andrew Aksyonoff
Copyright (c) 2008-2012, Sphinx Technologies Inc (http://sphinxsearch.com)

using config file '/etc/sphinxsearch/sphinx.conf'...
indexing index 'wiki_main'...
collected 1987 docs, 5.1 MB
collected 37 attr values
sorted 0.0 Mvalues, 100.0% done
sorted 21.5 Mhits, 100.0% done
total 1987 docs, 5053410 bytes
total 5.859 sec, 862461 bytes/sec, 339.11 docs/sec
indexing index 'wiki_incremental'...
collected 8 docs, 0.0 MB
collected 37 attr values
sorted 0.0 Mvalues, 100.0% done
sorted 0.1 Mhits, 100.0% done
total 8 docs, 12505 bytes
total 0.016 sec, 751141 bytes/sec, 480.53 docs/sec
total 139 reads, 0.018 sec, 359.0 kb/call avg, 0.1 msec/call avg
total 133 writes, 0.076 sec, 781.3 kb/call avg, 0.5 msec/call avg
rotating indices: succesfully sent SIGHUP to searchd (pid=22053).
```

If this is the first time and it works (no configuration problem), start the deamon (you'll to have setup this before):

```bash
service sphinxsearch start
```

### Test your indexation

There are several way to test your indexation but you need to know that the search binary contains bugs. If it crash, it doesn't mean that you have a problem with. Anyway, here is how to test:

```bash
> search -q --config /etc/sphinxsearch/sphinx.conf "test"
Sphinx 2.0.4-release (r3135)
Copyright (c) 2001-2012, Andrew Aksyonoff
Copyright (c) 2008-2012, Sphinx Technologies Inc (http://sphinxsearch.com)

using config file '/etc/sphinxsearch/sphinx.conf'...
index 'wiki_main': query 'test ': returned 119 matches of 119 total in 0.000 sec


1. document=2020, weight=1670, page_namespace=0, page_is_redirect=0, old_id=12125, category=()
2. document=1024, weight=1669, page_namespace=0, page_is_redirect=0, old_id=9383, category=()
3. document=3323, weight=1668, page_namespace=0, page_is_redirect=0, old_id=10361, category=()
[...]
20. document=1776, weight=1639, page_namespace=0, page_is_redirect=0, old_id=9510, category=(2575,2577,2578)

words:
1. 'test': 119 documents, 346 hits

index 'wiki_incremental': query 'test ': returned 0 matches of 0 total in 0.000 sec

words:
1. 'test': 0 documents, 0 hits
```

As you can see, we have results here :-). The work "test" have been found 20 times.

### Incremental updates

We need to setup the incremental updates. Change it to a slower value if you need to have more often indexation. For my own usage, once by hour, is really enough. I've added the MediaWiki example here (`/etc/cron.d/sphinxsearch`):

```bash
# Rebuild all indexes daily and notify searchd.
@daily      root . /etc/default/sphinxsearch && if [ "$START" = "yes" ] && [ -x /usr/bin/indexer ]; then /usr/bin/indexer --quiet --rotate --all >/dev/null 2>&1
; fi

# Example for rotating only specific indexes (usually these would be part of
# a larger combined index).

# */5 * * * * root [ -x /usr/bin/indexer ] && /usr/bin/indexer --quiet --rotate postdelta threaddelta >/dev/null 2>&1

# Mediawiki
0 */1 * * * root [ -x /usr/bin/indexer ] && indexer wiki_incremental --quiet --rotate >/dev/null 2>&1
```

## Debug

What if you don't see any results or you want to be sure that Sphinx receive search requests? There is a console mode:

```bash {linenos=table,hl_lines=[10]}
> searchd --console --config /etc/sphinxsearch/sphinx.conf --pidfile
Sphinx 2.0.4-release (r3135)
Copyright (c) 2001-2012, Andrew Aksyonoff
Copyright (c) 2008-2012, Sphinx Technologies Inc (http://sphinxsearch.com)

using config file '/etc/sphinxsearch/sphinx.conf'...
listening on 127.0.0.1:9312
accepting connections

[Fri Aug  2 14:19:15.652 2013] 0.000 sec [ext2/2/rel 50 (0,20)] [*] test
```

I see my test search here :-). All is good! If you don't see anything, that should be a problem with the application API or a missmatch configuration. You can also check with tcpdump if you see network connections arriving on 9312 port.

## FAQ

### I don't see any search result on MediaWiki, why?

You certainly have a problem with your php API. Select another version that should match. Check also the [Debug part](#debug) to help you to see what's wrong.

## References

[^1]: http://www.mediawiki.org/wiki/Extension:SphinxSearch
