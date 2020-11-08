# TFSchema

TFSchema is a command line tool to inspect resource and data source schemas from terraform providers.
It does this by parsing the output from a `terraform providers schema` command. This has the following benefits:

- Can easily look up inputs for a resource/data source from the CLI without having to open a browser.
- The schema is pulled from your current providers, so the version will always be what you need.
- It can work on custom providers as well, which often lack the nice documentation of official hashicorp providers.

## Installation

TFSchema is distributed as a binary hosted on the github releases page. Simply download it and add it to your path, and invoke it with

```shell script
tfschema
```

## Usage - `get-schema`

This command will simply run `terraform providers schema`, and print the output to stdout.
This can be useful for saving the schema locally if you need to inspect it.
Most `tfschema` commands do this implicitly on start, but for very large providers this can sometimes cause a slight delay in the CLI startup.
If you're working heavily with tfschema, you can save the provider schema in a file and reference it with other commands.

```shell script
tfschema get-schema > aws.json

# Implicitly get schema on startup, or read it from a file
tfschema resources 
tfschema resources -f aws.json
```

## Usage - `resources`

The resources command can be used to browse resources in a provider.
When run without arguments, it will list all resources in a provider.
This is easily combined with grep to narrow your search.

```shell script
# List all resources in a provider.
tfschema resources
tfschema resources | grep ec2
```

Once you find a resource you want to view in more detail, you can re-run the resources command to see the resource in greater detail. This will print out information from the schema, like attributes and descriptions. By default computed values are omitted to reduce noise, but can be printed by using the verbose `-v` flag.

```shell script
# Print descriptions and attributes for s3 buckets
tfschema resources aws_s3_bucket 

# Include computed attributes, like arns
tfschema resources aws_s3_bucket -v 
```

![Resource Detail](docs/resource-detail.png?raw=true "Resource Detail") 

## Supported Providers
TFSchema should work with any terraform providers supporting terraform 0.12.X, testing and development has focused mainly on the following:

- [ ] AWS Provider
- [ ] Random Provider


## Roadmap

- [x] Retrieve the schema from the terraform CLI transparently for commands.
- [x] Add support for resources, under a `resource` command
- [ ] Add support for data sources, under a `datasource` command
- [ ] Add support for block level attributes, like lists and maps
- [ ] Makefile
- [ ] Formatting Improvements
    - [ ] Better display for multiple providers
