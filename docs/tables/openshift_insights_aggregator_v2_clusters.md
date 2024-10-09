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

### Count clusters by version

```sql
SELECT
    cluster_version, COUNT(*) AS cluster_count
FROM crc_openshift_insights_aggregator_v2_clusters
GROUP BY cluster_version
ORDER BY cluster_count DESC;
```

### Calculate min, max, and mean total_hit_count by cluster version

```sql
SELECT
    cluster_version,
    MIN(total_hit_count) AS min_total_hit_count,
    MAX(total_hit_count) AS max_total_hit_count,
    AVG(total_hit_count) AS mean_total_hit_count
FROM crc_openshift_insights_aggregator_v2_clusters
GROUP BY cluster_version
ORDER BY cluster_version DESC;
```

### Average total hit count for managed vs non managed clusters

```sql
> SELECT
    managed,
    AVG(total_hit_count) AS average_total_hit_count
FROM crc_openshift_insights_aggregator_v2_clusters
GROUP BY managed;
```

### Get clusters checked in the last 30 days

```sql
SELECT
    cluster_id, cluster_name, last_checked_at
FROM crc_openshift_insights_aggregator_v2_clusters
WHERE last_checked_at >= NOW() - INTERVAL '30 days'
```
