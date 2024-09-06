# CRC Plugin for Steampipe

Use SQL to query [console.redhat.com APIs](console.redhat.com/docs/api).

Get started: // TODO: Link to https://hub.steampipe.io/plugins/jdiazsua/crc
Documentation: Table definitions & examples // TODO: Link to https://hub.steampipe.io/plugins/jdiazsua/crc/tables
Community: // TODO: link to a Slack channel
Get involved: // TODO: link to issues

##  Quick start

// TODO

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
> SELECT version, conditions, gathering_functions FROM crc.openshift_insights_gcs_v1_gathering_rules;
```

You can check the plugin and steampipe logs using
```sh
tail -f ~/.steampipe/logs/*$(date "+%Y-%m-%d").log;                                                                           
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Examples

### openshift_insights_aggregator_v2_clusters

The queries involving this table takes some time because the API endpoint is quite slow.

#### List your clusters

```sql
SELECT
    cluster_id, cluster_name, cluster_version, managed, last_checked_at,
    total_hit_count, hits_by_total_risk
FROM crc.openshift_insights_aggregator_v2_clusters
```

#### Find problematic clusters

```sql
SELECT
    cluster_id, cluster_name, cluster_version, managed, last_checked_at,
    total_hit_count, hits_by_total_risk
FROM crc.openshift_insights_aggregator_v2_clusters
WHERE total_hit_count > 4
```

### openshift_insights_gcs_v1_gathering_rules

```sql
SELECT version, conditions, gathering_functions
FROM crc.openshift_insights_gcs_v1_gathering_rules;
```

### openshift_insights_gcs_v2_gathering_rules

#### Get the gathering rules for a valid version

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc.openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = '4.17.0';
```

#### Get the gathering rules for a version that is not available

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc.openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = '3.0.0';
```

This will print a 404 error.

#### Get the gathering rules for wrong version

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc.openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = 'foo';
```

This will print a 400 error.
