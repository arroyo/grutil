# cmsutil

Headless CMS CLI

*Alpha*

Log any issues to Github issues.

The current version works with [GraphCMS](https://graphcms.com/) and backs up your project to disk, including asset files.

## Config

In your home directory create a folder called ".cmsutil" and put a text file called "config.yaml" in it.

There is an example config file in the repo.  

It has been tested with YAML, but it is possible to use JSON or TOML (use the extensions .json and .toml respectively).

CMS_API_URL must start with http:// or https://

### Envars

You can override the settings in your yaml config with the following environment variables

CMS_API_URL

CMS_API_KEY

## Usage

Backup your content

```
cmsutil download
```

Get help

```
cmsutil help
```
