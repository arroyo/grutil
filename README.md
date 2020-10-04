# cmsutil

Headless CMS CLI Utility

*Pre-Alpha*

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

Get a copy of your schema or content with the download command. Backup your entire project with the backup command.

Download schema models and enumerations

```
cmsutil download schemas
```

Download node content and assets
```
cmsutil download content
```

Backup your site

```
cmsutil backup
```

Get general help or help with a specific command using.

```
cmsutil help
cmsutil download -h
```
