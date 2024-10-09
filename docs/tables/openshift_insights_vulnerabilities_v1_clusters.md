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

### List clusters with critical vulnerabilities
This query retrieves clusters that have critical vulnerabilities, allowing you to focus on the most urgent security issues.

```sql
SELECT
    cluster_id,
    display_name,
    critical_cves,
    last_seen
FROM crc_openshift_insights_vulnerabilities_v1_clusters
WHERE critical_cves > 0
ORDER BY critical_cves DESC
```

### Count of vulnerabilities by severity
This query provides a summary of the total number of vulnerabilities categorized by severity level across all clusters.

```sql
SELECT
    SUM(low_cves) AS total_low,
    SUM(moderate_cves) AS total_moderate,
    SUM(important_cves) AS total_important,
    SUM(critical_cves) AS total_critical
FROM crc_openshift_insights_vulnerabilities_v1_clusters
```

### List the most exposed images along with their cluster details
This query retrieves the most exposed images from the clusters, providing insights into which images are most frequently identified as exposed.

```sql
SELECT 
    e.name AS image_name, 
    e.registry, 
    COUNT(e.cluster_id) AS exposure_count, 
    c.display_name AS cluster_name
FROM 
    crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images e
JOIN 
    crc_openshift_insights_vulnerabilities_v1_clusters c 
ON 
    e.cluster_id = c.cluster_id
GROUP BY 
    e.name, e.registry, c.display_name
ORDER BY 
    exposure_count DESC;
```