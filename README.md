# CRC Plugin for Steampipe

Use SQL to query [console.redhat.com APIs](console.redhat.com/docs/api).

Get started: https://hub.steampipe.io/plugins/juandspy/crc

Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/juandspy/crc/tables)

Community:

Get involved: https://github.com/juandspy/steampipe-plugin-crc/issues

##  Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install juandspy/crc
```

Authenticate with the plugin:

```
export CRC_CLIENT_ID="12345678-0000-1111-2222-123456789012"
export CRC_CLIENT_SECRET="abcdefghijklmnopqrstuvwxyz123456"
```

Run a query:

```sql
SELECT
    cluster_name, cluster_version,
    total_hit_count, hits_by_total_risk
FROM crc_openshift_insights_aggregator_v2_clusters
LIMIT 3
```

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
> SELECT version, conditions, gathering_functions FROM crc_openshift_insights_gcs_v1_gathering_rules;
```

You can check the plugin and steampipe logs using
```sh
tail -f ~/.steampipe/logs/*$(date "+%Y-%m-%d").log;                                                                           
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Examples

Visit the [tables folder](./docs/tables) for examples.
