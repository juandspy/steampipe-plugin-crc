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
