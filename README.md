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

This application is configured via the environment. The following environment
variables can be used:

| KEY                      | TYPE                             | DEFAULT                 | REQUIRED   | DESCRIPTION
|--------------------------|----------------------------------|-------------------------|------------|------------
| REDIS2ES_KEY             | String                           | logstash                | False      | the redis key where to listen to
| REDIS2ES_HOST            | String                           | 127.0.0.1               | False      | the redis server host/ip
| REDIS2ES_PORT            | Integer                          | 6379                    | False      | the redis server port
| REDIS2ES_DB              | Integer                          | 0                       | False      | the redis database
| REDIS2ES_PASSWORD        | String                           |                         | False      | the redis password
| REDIS2ES_USETLS          | True or False                    | false                   | False      | connect to redis using tls
| REDIS2ES_TLSSKIPVERIFY   | True or False                    | false                   | False      | if connection to redis via tls, skip tls certificate verification
| REDIS2ES_ELASTICURLS     | Comma-separated list of String   | http://127.0.0.1:9200   | False      | the elasticsearch connection url, seperated by comma for many es servers
| REDIS2ES_BULKSIZE        | Integer                          | 1000                    | False      | writes to elastic are done in bulks of bulkSize
| REDIS2ES_DEBUG           | True or False                    | false                   | False      | turn on debug log


## Filters

Filters can be implemented in go and loaded on startup with the go plugin mechanism.
For a start look at the sample plugin here [uppercase](filter/uppercase)