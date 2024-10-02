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
