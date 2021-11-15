# Grutil 

A simple docker container for grutil. 

## Docker

docker pull arroyo/grutil

## Usage 

See github for commands and cli options 

https://github.com/arroyo/grutil

### Examples with docker run

Get help

```
docker run arroyo/grutil grutil help
```

Backup with environment variables passed for your api details

```
docker run -e CMS_API_URL=${CMS_API_URL} -e CMS_API_KEY=${CMS_API_KEY} arroyo/grutil grutil backup
```
