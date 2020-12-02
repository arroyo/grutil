# cmsutil

Headless CMS CLI Utility

Please log any issues to Github issues.

The current version works with [GraphCMS](https://graphcms.com/) and backs up your project to disk, including asset files.

## Install

Once you have [Go installed](https://golang.org/doc/install) on your target environment, to install cmsutil run

```
go get github.com/arroyo/cmsutil
```

## Config

In your home directory create a folder called ".cmsutil" and put a text file called "config.yaml" in it.

There is an example config file in the repo.  

It has been tested with YAML, but it is possible to use JSON or TOML (use the extensions .json and .toml respectively).

CMS_API_URL must start with https:// or http://

### Envars

You can override the settings in your yaml config with the following environment variables

CMS_API_URL

CMS_API_KEY

## Usage

Get a copy of your schema or content with the download command. Backup your entire project with the backup command.

### Download schema models and enumerations

```shell-script
cmsutil download schemas
```

### Download node content and assets

```shell-script
cmsutil download content
```

### Backup your site

download both schemas and content into a timestamped folder

```shell-script
cmsutil backup
```

### Help

Get general help or help with a specific command using.

```shell-script
cmsutil help
cmsutil download -h
cmsutil backup -h
```
