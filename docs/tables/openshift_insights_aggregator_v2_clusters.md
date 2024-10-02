---
title: "Steampipe Table: openshift_insights_aggregator_v2_clusters - List clusters using OpenShift Insights"
description: "Allows users to query OpenShift Insights to retrieve information about clusters."
---

# Table: openshift_insights_aggregator_v2_clusters - Query Insights Aggregator clusters using SQL

The Insights Aggregator is a service that identifies OpenShift clusters
in your organization and propose a set of recommendations to improve the
configuration of the clusters.

## Examples

### List your clusters

```sql
SELECT
    cluster_id, cluster_name, cluster_version, managed, last_checked_at,
    total_hit_count, hits_by_total_risk
FROM crc_openshift_insights_aggregator_v2_clusters
LIMIT 10
```

### Find problematic clusters

```sql
SELECT
    cluster_id, cluster_name, cluster_version, managed, last_checked_at,
    total_hit_count, hits_by_total_risk
FROM crc_openshift_insights_aggregator_v2_clusters
WHERE total_hit_count > 4
```
