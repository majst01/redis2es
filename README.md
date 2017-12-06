# Redis to Elasticsearch

[![Build Status](https://travis-ci.org/majst01/redis2es.svg?branch=master)](https://travis-ci.org/majst01/redis2es)
[![codecov](https://codecov.io/gh/majst01/redis2es/branch/master/graph/badge.svg)](https://codecov.io/gh/majst01/redis2es)
[![Go Report Card](https://goreportcard.com/badge/majst01/redis2es)](https://goreportcard.com/report/github.com/majst01/redis2es)


In order to ship logs from applications to a cluster of elasticsearch servers, most of the time a combination of redis and logstash is in place.
Redis is used to decouple write load from receive performance of elastic.
Logstash is in place to filter the incoming logs with certain criteria.

Logstash is known to use a lot of memory to work properly.


## Deployment
Old Picture
```graphviz
fluent-bit/logstash --> stunnel --> redis <-- logstash --> elasticsearch
```

New Picture
```graphviz
fluent-bit/logstash --> stunnel --> redis <-- redis-to-elastic --> elasticsearch
```

In a Elastic Cluster
```graphviz
                                                                /-> elasticsearch
                                                               |
fluent-bit/logstash --> stunnel --> redis <-- redis-to-elastic --> elasticsearch
                                                               |
                                                               \-> elasticsearch
```

## Configuration

Configuration is done by environment variables, available variables are:


| Variable            | Description                                    | Default        |
| --------------------|------------------------------------------------|----------------|
| DEBUG               | Enable/Disable debug logging true|false        | false          |
| REDIS_KEY           | The key where the log entries are stored       | logstash       |
| REDIS_HOST          | The hostname/ip of redis                       | 127.0.0.1      |
| REDIS_Port          | the port of redis                              | 6379           |
| REDIS_DB            | the database number of redis                   | 0              |
| REDIS_PASSWORD      | the redis database password                    | ""             |
| REDIS_USETLS        | use tls encryption to redis                    | false          |
| REDIS_TLSSKIPVERIFY | if usetls is set verify certificates           | false          |
| ELASTIC_URLS        | the comma seperated urls of the elastic servers| "127.0.0.1:9200"|

## Filters

Filters can be implemented in go and loaded on startup with the go plugin mechanism.
For a start look at the sample plugin here [uppercase](filter/uppercase)