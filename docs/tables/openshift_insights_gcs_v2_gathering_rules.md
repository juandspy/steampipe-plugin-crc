---
title: "Steampipe Table: openshift_insights_gcs_v2_gathering_rules - List gathering rules using Openshift Insights Gathering Service (v2)"
description: "Allows users to query OpenShift Insights Gathering Service to retrieve information about gathering rules."
---

# Table: openshift_insights_gcs_v2_gathering_rules - Query Insights Gathering Service using SQL

The Insights Gathering Service is a service that is used by the Insights
Operator running on the connected clusters. It retrieves the configuration
from the Gathering Service and uses it to gather custom data from the clusters.

## Examples

### Get the gathering rules for a valid version

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc_openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = '4.17.0';
```

### Get the gathering rules for a version that is not available

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc_openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = '3.0.0';
```

This will print a 404 error.

### Get the gathering rules for wrong version

```sql

SELECT version, conditional_gathering_rules, container_logs
FROM crc_openshift_insights_gcs_v2_gathering_rules
WHERE ocp_version = 'foo';
```

This will print a 400 error.
