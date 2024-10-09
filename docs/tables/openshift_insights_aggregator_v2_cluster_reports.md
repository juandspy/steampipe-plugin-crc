---
title: "Steampipe Table: openshift_insights_aggregator_v2_cluster_reports - List cluster reports using OpenShift Insights"
description: "Allows users to query OpenShift Insights to retrieve information about cluster reports."
---

# Table: openshift_insights_aggregator_v2_cluster_reports - Query Insights Aggregator reports using SQL

The Insights Aggregator is a service that identifies OpenShift clusters
in your organization and propose a set of recommendations to improve the
configuration of the clusters.

## Examples

### List cluster rules with total risk

```sql
SELECT cluster_id, rule_id, total_risk
FROM crc_openshift_insights_aggregator_v2_cluster_reports
WHERE cluster_id = '5a78700a-e3d3-4300-a796-75bf73fc1653'
```

### Count reports for a cluster

```sql
SELECT COUNT(*) AS report_count
FROM crc_openshift_insights_aggregator_v2_cluster_reports
WHERE cluster_id = '5a78700a-e3d3-4300-a796-75bf73fc1653'
GROUP BY cluster_id
ORDER BY report_count DESC
```

### Get rules with a high total risk

```sql
SELECT cluster_id, rule_id, total_risk
FROM crc_openshift_insights_aggregator_v2_cluster_reports
WHERE cluster_id = '5a78700a-e3d3-4300-a796-75bf73fc1653'
AND total_risk >= 3
ORDER BY total_risk DESC
```

### Get the highest total risk report for a specific cluster

```sql
SELECT rule_id, total_risk
FROM crc_openshift_insights_aggregator_v2_cluster_reports
WHERE cluster_id = '5a78700a-e3d3-4300-a796-75bf73fc1653'
ORDER BY total_risk DESC
LIMIT 1
```

### Count of rules across some clusters

These queries makes use of the `openshift_insights_aggregator_v2_clusters` 
table.
I've used `LIMIT` as my organization contains thousands of clusters and it
can get slow.

#### Total rules per cluster

```sql
SELECT c.cluster_id, c.cluster_name, COUNT(r.rule_id) AS rule_count
FROM (SELECT cluster_id, cluster_name, total_hit_count
      FROM crc_openshift_insights_aggregator_v2_clusters
      WHERE total_hit_count > 0
      LIMIT 10) AS c
LEFT JOIN crc_openshift_insights_aggregator_v2_cluster_reports AS r
ON c.cluster_id = r.cluster_id
GROUP BY c.cluster_id, c.cluster_name
ORDER BY rule_count DESC;
```

#### Occurrences of each rule across the fleet

```sql
SELECT r.rule_id, COUNT(*) AS occurrence_count
FROM (SELECT cluster_id, cluster_name, total_hit_count
      FROM crc_openshift_insights_aggregator_v2_clusters
      WHERE total_hit_count > 0
      LIMIT 100) AS c
LEFT JOIN crc_openshift_insights_aggregator_v2_cluster_reports AS r
ON c.cluster_id = r.cluster_id
GROUP BY r.rule_id
ORDER BY occurrence_count DESC;
```
