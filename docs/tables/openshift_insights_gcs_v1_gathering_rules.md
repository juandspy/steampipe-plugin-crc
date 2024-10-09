---
title: "Steampipe Table: openshift_insights_gcs_v1_gathering_rules - List gathering rules using Openshift Insights Gathering Service (v1)"
description: "Allows users to query OpenShift Insights Gathering Service to retrieve information about gathering rules."
---

# Table: openshift_insights_gcs_v1_gathering_rules - Query Insights Gathering Service using SQL

The Insights Gathering Service is a service that is used by the Insights
Operator running on the connected clusters. It retrieves the configuration
from the Gathering Service and uses it to gather custom data from the clusters.

## Examples

### List the conditions and gathering functions for all versions

```sql
SELECT version, conditions, gathering_functions
FROM crc_openshift_insights_gcs_v1_gathering_rules;
```

### Get the count of gathering rules by version

```sql
SELECT version, COUNT(*) AS rule_count
FROM crc_openshift_insights_gcs_v1_gathering_rules
GROUP BY version;
```

### Get the gathering functions for a given condition

```sql
SELECT *
FROM crc_openshift_insights_gcs_v1_gathering_rules
WHERE conditions::jsonb @> '[{"alert":{"name":"KubePodCrashLooping"},"type":"alert_is_firing"}]';
```
