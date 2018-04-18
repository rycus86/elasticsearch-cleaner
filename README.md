# Elasticsearch index cleaner

Small Go application to clear old Elasticsearch indices.

## Motivation

As descibed [on my blog](https://blog.viktoradam.net/2018/02/06/home-lab-part5-monitoring-madness/#logging), if I want to keep only a certain amount of logs collected into Elasticsearch, I need to periodically delete the old indices that correspond to days too far in the past. This can be done manually with a simple `curl` command.

```shell
$ curl -X DELETE http://elasticsearch:9200/fluentd-20180412
```

## Usage

The application is available on [Docker Hub](https://hub.docker.com/r/rycus86/elasticsearch-cleaner/), and to run it, you can:

```shell
docker run --rm -it \
	-e BASE_URL=http://elasticsearch:9200 	\
	-e PATTERN=logstash-.* 			\
	-e MAX_INDICES=14 			\
	-e INTERVAL=12h 			\
	-e TIMEOUT=30s 				\
	rycus86/elasticsearch-cleaner
```

The available environment variables for the configuration are:

- __BASE_URL__: The URL of the Elasticsearch server
- __PATTERN__: The regular expression to look for matching indices to delete
- __MAX_INDICES__: The number of (matching) indices to keep *(default: 20)*
- __INTERVAL__: The interval between checking for old indices *(default: 12h)*
- __TIMEOUT__: The timeout for HTTP calls to the Elasticsearch server *(default: 30s)*

## License

MIT
