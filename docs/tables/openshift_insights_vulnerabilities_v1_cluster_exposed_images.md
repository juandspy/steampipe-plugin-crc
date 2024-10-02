---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_cluster_exposed_images - Query OpenShift Insights Vulnerabilities Cluster Exposed Images using SQL"
description: "Allows users to query Cluster Exposed Images in OpenShift Insights Vulnerabilities, providing information about potentially vulnerable images in specific clusters."
---

# Table: openshift_insights_vulnerabilities_v1_cluster_exposed_images - Query OpenShift Insights Vulnerabilities Cluster Exposed Images using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift clusters. The Cluster Exposed Images component specifically focuses on identifying container images within clusters that may be exposed to known vulnerabilities. This information is crucial for maintaining the security and integrity of your OpenShift environment.

## Table Usage Guide

The `openshift_insights_vulnerabilities_v1_cluster_exposed_images` table provides insights into potentially vulnerable container images within specific OpenShift clusters. As a security analyst or cluster administrator, explore image-specific details through this table, including image names, registries, and versions. Utilize it to uncover information about exposed images, assess the vulnerability status of your container ecosystem, and prioritize image updates or replacements.

### List all CVEs affecting the current workload

```sql
SELECT synopsis, severity, cvss3_score, clusters_exposed, images_exposed
FROM crc_openshift_insights_vulnerabilities_v1_cves
ORDER BY cvss3_score DESC
LIMIT 10
```
