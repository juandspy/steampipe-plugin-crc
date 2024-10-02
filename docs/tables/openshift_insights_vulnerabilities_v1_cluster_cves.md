---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_cluster_cves - Query OpenShift Insights Vulnerabilities for a given Cluster CVEs using SQL"
description: "Allows users to query Cluster CVEs in OpenShift Insights Vulnerabilities, providing detailed information about Common Vulnerabilities and Exposures (CVEs) affecting specific clusters."
---

# Table: openshift_insights_vulnerabilities_v1_cluster_cves - Query OpenShift Insights Vulnerabilities Cluster CVEs using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift clusters. The Cluster CVEs component specifically focuses on identifying and detailing Common Vulnerabilities and Exposures (CVEs) that affect individual clusters. This information is crucial for maintaining the security and integrity of your OpenShift environment.

## Table Usage Guide

The openshift_insights_vulnerabilities_v1_cluster_cves table provides insights into CVEs affecting specific OpenShift clusters. As a security analyst or cluster administrator, explore CVE-specific details through this table, including severity levels, CVSS scores, and affected clusters. Utilize it to uncover information about high-severity vulnerabilities, assess the overall security posture of specific clusters, and prioritize patching efforts.

## Examples

### List CVEs for a cluster

```sql
SELECT synopsis, severity, cvss3_score
FROM crc_openshift_insights_vulnerabilities_v1_cluster_cves
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
```

### List CVEs for a cluster with a severity filter

```sql
SELECT synopsis, severity, cvss3_score
FROM crc_openshift_insights_vulnerabilities_v1_cluster_cves
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
AND severity = 'Low'
```

