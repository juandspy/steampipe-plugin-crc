# CRC Plugin for Steampipe

Use SQL to query [console.redhat.com APIs](console.redhat.com/docs/api).

Get started: // TODO: Link to https://hub.steampipe.io/plugins/jdiazsua/crc
Documentation: Table definitions & examples // TODO: Link to https://hub.steampipe.io/plugins/jdiazsua/crc/tables
Community: // TODO: link to a Slack channel
Get involved: // TODO: link to issues

##  Quick start

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone git@github.com:juandspy/steampipe-plugin-crc
cd steampipe-plugin-crc
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make install
```

Configure the plugin:

```
cp config/crc.spc ~/.steampipe/config
vi ~/.steampipe/config/crc.spc
```

Note that there are some variables that need to be defined there, used for authentication.

Try it!

```
steampipe query
> .inspect crc
> 
```

You can check the plugin and steampipe logs using
```sh
tail -f ~/.steampipe/logs/*$(date "+%Y-%m-%d").log;                                                                           
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)
