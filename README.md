# Elasticsearch client examples

[Official Elasticsearch client](https://www.elastic.co/guide/en/elasticsearch/client/index.html) examples of Go, Node.js and Python.

- [Go](https://github.com/elastic/go-elasticsearch)
- [Node.js](https://github.com/elastic/elasticsearch-js)
- [Python](https://github.com/elastic/elasticsearch-py)

## Envionment

- Elasticsearch v7.9.3
- Go 1.15
- Node.js v7.9.3
- Python v7.9.3

## Preparing

### Elastic Stack

Start Elasticsearch and Kibana docker container by docker-compose.

```sh
docker-compose up -d
```

Post [movie data](https://docs.aws.amazon.com/ja_jp/elasticsearch-service/latest/developerguide/samples/sample-movies.zip) to Elasticsearch.

```sh
wget https://docs.aws.amazon.com/ja_jp/elasticsearch-service/latest/developerguide/samples/sample-movies.zip
unzip sample-movies.zip
# post data by REST API
curl -H 'Content-Type: application/json' -XPOST 'localhost:9200/_bulk?pretty' --data-binary "@sample-movies.bulk"
rm sample-movies.bulk sample-movies.zip
```

[Access Kibana dev-tools](http://localhost:5601/app/dev_tools#/console) then test query.

```json
GET /movies/_search
{
    "size": 50,
    "query": {
        "multi_match": {
            "query": "thor",
            "fields": ["fields.title^4", "fields.plot^2", "fields.actors", "fields.directors"]
        }
    }
}
```

Response will be returned like.

```json
{
  "took" : 10,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 7,
      "relation" : "eq"
    },
    "max_score" : 40.92619,
    "hits" : [
      {
        "_index" : "movies",
        "_type" : "movie",
        "_id" : "tt0800369",
        "_score" : 40.92619,
        "_source" : {
          "fields" : {
            "directors" : [
              "Kenneth Branagh",
              "Joss Whedon"
            ],
            "release_date" : "2011-04-21T00:00:00Z",
            "rating" : 7,
            "genres" : [
              "Action",
              "Adventure",
              "Fantasy"
            ],
            "image_url" : "https://m.media-amazon.com/images/M/MV5BMTYxMjA5NDMzNV5BMl5BanBnXkFtZTcwOTk2Mjk3NA@@._V1_SX400_.jpg",
            "plot" : "The powerful but arrogant god Thor is cast out of Asgard to live amongst humans in Midgard (Earth), where he soon becomes one of their finest defenders.",
            "title" : "Thor",
            "rank" : 135,
            "running_time_secs" : 6900,
            "actors" : [
              "Chris Hemsworth",
              "Anthony Hopkins",
              "Natalie Portman"
            ],
            "year" : 2011
          },
          "id" : "tt0800369",
          "type" : "add"
        }
      },
      {
        "_index" : "movies",
        "_type" : "movie",
        "_id" : "tt1981115",
        "_score" : 25.360119,
        "_source" : {
          "fields" : {
              ...
          },
          "id" : "tt1981115",
          "type" : "add"
        }
      }
    ]
  }
}
```
