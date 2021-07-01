# Grutil

A helpful GraphQL utility

Please log any issues to Github issues.

The current version works with [GraphCMS](https://graphcms.com/) and downloads schemas, content, backs up your project to disk (including asset files) and allows simple rendering.

**Functions**
* Download
* Backup
* Render

## Install

Once you have [Go installed](https://golang.org/doc/install) on your target environment, simply run the following to install

```
go get github.com/arroyo/grutil
```

## Config

In your home directory create a folder called ".grutil" and put a text file called "config.yaml" in it.

There is an example config file in the repo.  

It has been tested with YAML, but it is possible to use JSON or TOML (use the extensions .json and .toml respectively).

CMS_API_URL must start with https:// or http://

### Envars

You can override the settings in your yaml config with the following environment variables

CMS_API_URL

CMS_API_KEY

## Usage

Get a copy of your schema or content with the download command. Backup your entire project with the backup command.  Render a GraphQL query against a template.

You can add optional --verbose and --debug flags to any command.  Verbose will add more details of what is happening as the program runs.  Debug is useful if you are having issues and want to see what API calls are happening behind the scenes.  Debug will show the GraphQL query being made and the API response body.

### Download schema models and enumerations

```shell-script
grutil download schemas --verbose
```

### Download node content and assets

```shell-script
grutil download content
```

### Backup your site

download both schemas and content into a timestamped folder

```shell-script
grutil backup
```

### Render content

Query content with GraphQL and render it against a go template.

```shell-script
grutil render --template json.txt --query "query MyQuery { faq(where: {id: \"cknmrjsvw7yd\"}) { id title publishedAt updatedAt } }"
```

### Get help

Get general help or specific help with a command using.

```shell-script
grutil help
grutil download -h
grutil backup -h
```
