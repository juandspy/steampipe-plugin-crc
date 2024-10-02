---
organization: Juandspy
category: ["public cloud"]
icon_url: "/images/plugins/juandspy/crc.svg"
brand_color: "#EE0000"
display_name: Red Hat Console Dot
name: crc
description: Steampipe plugin for querying console.redhat.com data.
og_description: Query console.redhat.com with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/juandspy/crc-social-graphic.png"
---

# Red Hat Console Dot (crc)

The Red Hat Console Dot (crc) plugin allows you to query console.redhat.com
data including clusters, vulnerabilities, recommendations and more.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

[Console Dot](console.redhat.com) provides a UI to manage your clusters and
workstations. It also provides [APIs](console.redhat.com/docs/api) to access
this data.

For example:

```sql
SELECT
    cluster_name, cluster_version,
    total_hit_count, hits_by_total_risk
FROM crc_openshift_insights_aggregator_v2_clusters
LIMIT 3
```
```
+--------------------------------------+-----------------+-----------------+---------------------------+
| cluster_name                         | cluster_version | total_hit_count | hits_by_total_risk        |
+--------------------------------------+-----------------+-----------------+---------------------------+
| 00000000-1111-2222-3333-444444444444 |                 | 0               | {"1":0,"2":0,"3":0,"4":0} |
| My Testing Cluster                   | 4.16.8          | 2               | {"1":0,"2":0,"3":2,"4":0} |
| My CI/CD Cluster                     | 4.18.0-ec.0     | 0               | {"1":0,"2":0,"3":0,"4":0} |
+--------------------------------------+-----------------+-----------------+---------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/juandspy/crc/tables)**

## Get started

### Install

Download and install the latest CRC plugin:

```bash
steampipe plugin install juandspy/crc
```

### Credentials

You must specify the `client ID` and `client secret` in the credentials file
(`~/.steampipe/config/crc.spc`) or via environment variables. You can run
```
cp config/crc.spc ~/.steampipe/config
vi ~/.steampipe/config/crc.spc
```

to copy the default credentials file and edit it or set the environment
variables:
```
export CRC_CLIENT_ID="12345678-0000-1111-2222-123456789012"
export CRC_CLIENT_SECRET="abcdefghijklmnopqrstuvwxyz123456"
```

### Configuration

Installing the latest crc plugin will create a config file
(`~/.steampipe/config/crc.spc`) with a single connection named `crc`:
```hcl
connection "crc" {
  plugin = "crc"

  # The baseUrl (prod or stage) for the console.redhat.com APIs
  # Can also be set with the CRC_URL environment variable.
  base_url = "https://console.redhat.com/"

  # The tokenUrl (prod or stage) for updating the token used to communicate
  # with console.redhat.com
  # Can also be set with the CRC_TOKEN_URL environment variable.
  token_url = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"

  # The client ID to access the console.redhat.com cloud instance
  # Can also be set with the `CRC_CLIENT_ID` environment variable.
  # client_id = "12345678-0000-1111-2222-123456789012"

  # The client secret to access the console.redhat.com cloud instance
  # Can also be set with the `CRC_CLIENT_SECRET` environment variable.
  # client_secret = "abcdefghijklmnopqrstuvwxyz123456"
}
```

You can configure the base URL (and use console.stage.redhat.com),
or the token URL (and use sso.stage.redhat.com) for development.
