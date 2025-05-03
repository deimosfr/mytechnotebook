---
weight: 999
url: "/ElasticSearch\\:_powerful_search_and_analytics_engine/"
title: "ElasticSearch: Powerful Search and Analytics Engine"
description: "A guide to set up and use ElasticSearch, a flexible and powerful open source distributed real-time search and analytics engine."
categories: ["Programming", "Security", "Servers"]
date: "2014-08-12T14:07:00+02:00"
lastmod: "2014-08-12T14:07:00+02:00"
tags:
  [
    "Elasticsearch",
    "Search Engine",
    "Database",
    "Analytics",
    "Java",
    "JSON",
    "RESTful API",
    "Lucene",
  ]
toc: true
---

![Elastic Search](/images/elasticsearch_logo.avif)

## Introduction

Elasticsearch is a flexible and powerful open source, distributed, real-time search and analytics engine. Architected from the ground up for use in distributed environments where reliability and scalability are must haves, Elasticsearch gives you the ability to move easily beyond simple full-text search. Through its robust set of APIs and query DSLs, plus clients for the most popular programming languages, Elasticsearch delivers on the near limitless promises of search technology.

## Basics concepts

Here are a some Lucene information that you need to know:

- All the information of the structures are called inverted index.
- You can't modify, only delete then insert.
- Deletes (like on MariaDB XtraDB called "optimize") creates fragmentation. To merge data this process is called **segment merge**.

### Input data

Data analysis is made by the **analyser** which is built of a tokenizer and zero or more token filters, and it can also have zero or more character mappers. A **tokenizer** in Lucene is used to split the text into tokens and is built of zero or more token filters.

Filters are processed sequentially. The character mappers are used before the tokenizer. For example you can remove HTML tags with it.

{{< alert context="info" text="Remove all unnecessary fields like html tags to avoid mistaken scoring" />}}

### Index

A query may be not analyzed (you can decide). For example, the prefix and the term queries are not analyzed while the match query is!
In ElasticSearch, an index is like a table in MariaDB. Data is **stored in JSON format** called a "document".

### Architecture

ElasticSearch knows how to work in standalone mode or is able to work in cluster. **Cluster implies Sharding + Replication:**
![Es-cluster.png](/images/es-cluster.avif)

When you send a new document to the cluster, you specify a target index and send it to one node (any of available nodes). In cluster mode, ElasticSearch gateways forwards their data to the primary node. In a cluster, there is only one writing node that can switch to another node if this one falls down.

## Installation

To install ElasticSearch, you can take the last stable version available on the official repository. First of all install the repository key:

```bash
cd /tmp
wget -O - http://packages.elasticsearch.org/GPG-KEY-elasticsearch | apt-key add -
```

Add the repository file:

```bash
deb http://packages.elasticsearch.org/elasticsearch/1.3/debian stable main
```

Now install elasticsearch with the dependencies:

```bash
aptitude install elasticsearch openjdk-7-jre-headless openntpd
```

To finish configure the init file:

```bash
update-rc.d elasticsearch defaults 95 10
```

## Configuration

### File descriptors

To avoid reaching maximum file descriptor, you have to update the limits.conf file with those settings:

```bash
elasticsearch soft nofile 32000
elasticsearch hard nofile 32000
```

### JVM

Regarding the JVM parameters, it's recommended to use 1G (XMX) for small deployments. Check out your logs to see indications about OutOfMemoryError exceptions 'ES_HEAP_SIZE' variable size.

{{< alert context="info" text="You should avoid to allocate 50% of your total system memory to the JVM." />}}

### Cluster

Depending on the configuration you want to have (single or cluster), you have to edit 2 values in the default configuration file:

```bash
cluster.name: elasticsearch
node.name: "Node 1"
```

- cluster.name: set it if you want your server to join a cluster.
- node.name: set a hostname. If not set, it will take the server hostname.

### Dynamic scripting

You may want to enable dynamic scripting to do advanced query in cli. To enable it, add it in the configuration:

```bash
script.disable_dynamic: false
```

## Administration

### Check health

You can check your cluster health like this:

```javascript
> curl -XGET http://127.0.0.1:9200/_cluster/health?pretty
{
  "cluster_name" : "elasticsearch",
  "status" : "green",
  "timed_out" : false,
  "number_of_nodes" : 3,
  "number_of_data_nodes" : 3,
  "active_primary_shards" : 5,
  "active_shards" : 10,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 0
}
```

### Get nodes information

To get information regarding nodes, you can use 'cat':

```bash
> curl -XGET "http://127.0.0.1:9200/_cat/nodes?v&h=name,id,ip,port,v,m"
name  id   ip            port v     m
node1 YbCv 192.168.33.31 9300 1.2.2 m
node2 kXy7 192.168.33.32 9300 1.2.2 m
node3 VNK9 192.168.33.33 9300 1.2.2 *
```

The interesting things here are the master node (last column defined by '\*').

Or you can use this:

```javascript
> curl -XGET "http://127.0.0.1:9200/_nodes/process?pretty"
{
  "cluster_name" : "elasticsearch",
  "nodes" : {
    "c8AX1atwQ6C2hl13_S0r4g" : {
      "name" : "node3",
      "transport_address" : "inet[/192.168.33.33:9300]",
      "host" : "node3",
      "ip" : "192.168.33.33",
      "version" : "1.2.2",
      "build" : "9902f08",
      "http_address" : "inet[/192.168.33.33:9200]",
      "process" : {
        "refresh_interval_in_millis" : 1000,
        "id" : 3457,
        "max_file_descriptors" : 65535,
        "mlockall" : false
      }
    },
    "pmsqiKGHRMGEo3iWaxv3Gw" : {
      "name" : "node1",
      "transport_address" : "inet[/192.168.33.31:9300]",
      "host" : "node1",
      "ip" : "192.168.33.31",
      "version" : "1.2.2",
      "build" : "9902f08",
      "http_address" : "inet[/192.168.33.31:9200]",
      "process" : {
        "refresh_interval_in_millis" : 1000,
        "id" : 3480,
        "max_file_descriptors" : 65535,
        "mlockall" : false
      }
    },
    "xVXb40pgRNKdd9G6u8-7Uw" : {
      "name" : "node2",
      "transport_address" : "inet[/192.168.33.32:9300]",
      "host" : "node2",
      "ip" : "192.168.33.32",
      "version" : "1.2.2",
      "build" : "9902f08",
      "http_address" : "inet[/192.168.33.32:9200]",
      "process" : {
        "refresh_interval_in_millis" : 1000,
        "id" : 3886,
        "max_file_descriptors" : 65535,
        "mlockall" : false
      }
    }
  }
}
```

To get more information and options, look at the official documentation.

### Shutdown a node

To shutdown a specific node, use that curl command and replace the nodeid with the desired id number:

```javascript
> curl -XPOST http://127.0.0.1:9200/_cluster/nodes/<nodeid>/_shutdown?pretty
{
  "cluster_name" : "elasticsearch",
  "nodes" : {
    "zfnG3AKMShad0Ti9qgchFQ" : {
      "name" : "node2"
    }
  }
}
```

### Shutdown the cluster

If you want to shutdown the whole cluster at once:

```javascript
> curl -XPOST http://127.0.0.1:9200/_cluster/nodes/_shutdown?pretty
{
  "cluster_name" : "elasticsearch",
  "nodes" : {
    "FKCjz60DRgWCat7WE9NkBQ" : {
      "name" : "node3"
    },
    "IfQBC4VrRICLyO5pNsohHA" : {
      "name" : "node1"
    },
    "kzlYH_8rRBmWXCdZIjYrlQ" : {
      "name" : "node2"
    }
  }
}
```

## Usage

### Create a new entry

To create a new entry with it's automated index, you simply needs to insert like this:

```javascript
> curl -XPUT http://localhost:9200/vehicule/moto/1?pretty -d '{"vendor": "Kawazaki", "model": "Z1000", "tags": ["sports", "roadster"] }'
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 1,
  "created" : true
}
```

If everything was fine, you should have "created" value to true. Each time there will be an update on the document, the version will automatically increase. If you do not specify the id, it will automatically be generated:

```javascript
> curl -XPOST http://localhost:9200/vehicule/moto/?pretty -d '{"vendor": "Kawazaki", "model": "Z1000", "tags": ["sports", "roadster"] }'
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "q1mSSqHbSqCuOHLdUGVLYQ",
  "_version" : 1,
  "created" : true
}
```

### Get a document

To get a document (an entry), this is simple:

```javascript
> curl -XGET http://localhost:9200/vehicule/moto/1?pretty
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 1,
  "found" : true,
  "_source":{"vendor": "Kawazaki", "model": "Z1000", "tags": ["sports", "roadster"] }
}
```

You only have to know the id. If a document is not found:

```javascript {linenos=table,hl_lines=[6],anchorlinenos=true}
> curl -XGET http://localhost:9200/vehicule/moto/4?pretty
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "4",
  "found" : false
}
```

You'll get found value set to false

### Update a document

Lucene doesn't know how to update a document. So when you'll ask to ElasticSearch to update a document, you will in fact delete the current and create a new one. To modify a document (here the model value), you can do it like that:

```javascript
> curl -XPOST http://localhost:9200/vehicule/moto/1/_update?pretty -d '{"script": "ctx._source.model = \"Z800\""}'
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 2
}
```

As you can see the version number has been incremented.

To add a new field to a current document:

```javascript
curl -XPOST 'localhost:9200/vehicule/moto/1/_update?pretty' -d '{
>     "script" : "ctx._source.power = \"139cv\""
> }'
{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 11
}
```

If you want to add a tag in the current tag list of a document:

```javascript
> curl -XPOST 'localhost:9200/vehicule/moto/1/_update?pretty' -d '{
    "script" : "ctx._source.tags += tag",
    "params" : {
        "tag" : "white/orange"
    }
}'

{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 10
}
```

### Remove a document or it's content

To remove a complete document:

```javascript
> curl -XDELETE 'localhost:9200/vehicule/moto/4?pretty'
{
  "found" : true,
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "4",
  "_version" : 3
}
```

To remove a document field (here power):

```javascript
> curl -XPOST 'localhost:9200/vehicule/moto/1/_update?pretty' -d '{
    "script" : "ctx._source.remove(\"power\")"
}'

{
  "_index" : "vehicule",
  "_type" : "moto",
  "_id" : "1",
  "_version" : 13
}
```

ElasticSearch knows how to deal with concurrency, however if you really want to be sure to safely delete a document at a certain version, you can force it. It will fail if the document has changed in the meantime:

```javascript
> curl -XDELETE 'localhost:9200/vehicule/moto/4?version=15'
```

## References

1. [https://www.elasticsearch.org/overview/](https://www.elasticsearch.org/overview/)
2. [https://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-scripting.html](https://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-scripting.html)
3. [https://www.elasticsearch.org/guide/en/elasticsearch/reference/current/cat-nodes.html](https://www.elasticsearch.org/guide/en/elasticsearch/reference/current/cat-nodes.html)
