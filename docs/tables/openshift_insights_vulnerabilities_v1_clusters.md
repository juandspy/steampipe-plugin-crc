---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_clusters - Query OpenShift Insights Vulnerabilities using SQL"
description: "Allows users to query OpenShift Insights Vulnerabilities, providing information about cluster vulnerabilities and their severity levels."
---

# Table: openshift_insights_vulnerabilities_v1_clusters - Query OpenShift Insights Vulnerabilities using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift clusters. It helps identify and assess the severity of vulnerabilities, allowing administrators to take appropriate action to secure their clusters.

## Table Usage Guide

The `openshift_insights_vulnerabilities_v1_clusters` table provides insights into cluster vulnerabilities within OpenShift. As a security administrator, explore cluster-specific details through this table, including vulnerability counts of different severity levels, cluster status, and last seen timestamps. Utilize it to uncover information about clusters with high numbers of critical vulnerabilities, the overall security posture of your clusters, and the recency of vulnerability scans.

## Examples

### List your clusters
Explore the vulnerability status of your clusters to assess their security posture. This query provides a comprehensive overview of each cluster's vulnerability counts across different severity levels, helping you prioritize security efforts.

```sql
SELECT
    cluster_id,
    display_name,
    version,
    provider,
    last_seen,
    status,
    low_cves,
    moderate_cves,
    important_cves,
    critical_cves
FROM crc_openshift_insights_vulnerabilities_v1_clusters
```
