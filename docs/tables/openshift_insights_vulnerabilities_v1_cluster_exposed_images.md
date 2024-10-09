---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_cluster_exposed_images - Query OpenShift Insights Vulnerabilities Cluster Exposed Images using SQL"
description: "Allows users to query Cluster Exposed Images in OpenShift Insights Vulnerabilities, providing information about potentially vulnerable images in specific clusters."
---

# Table: openshift_insights_vulnerabilities_v1_cluster_exposed_images - Query OpenShift Insights Vulnerabilities Cluster Exposed Images using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift clusters. The Cluster Exposed Images component specifically focuses on identifying container images within clusters that may be exposed to known vulnerabilities. This information is crucial for maintaining the security and integrity of your OpenShift environment.

## Table Usage Guide

The `openshift_insights_vulnerabilities_v1_cluster_exposed_images` table provides insights into potentially vulnerable container images within specific OpenShift clusters. As a security analyst or cluster administrator, explore image-specific details through this table, including image names, registries, and versions. Utilize it to uncover information about exposed images, assess the vulnerability status of your container ecosystem, and prioritize image updates or replacements.

### List all exposed images in a specific cluster
```sql
SELECT cluster_id, name AS image_name, registry, version
FROM crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
ORDER BY image_name;
```

### List exposed images with a specific registry in a cluster
```sql
SELECT cluster_id, name AS image_name, version
FROM crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
AND registry = 'registry.connect.redhat.com'
ORDER BY version;
```

### Count of exposed images in a specific cluster
```sql
SELECT COUNT(*) AS exposed_image_count
FROM crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2';
```

### List exposed images with a specific version in a cluster
```sql
SELECT cluster_id, name AS image_name, registry
FROM crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images
WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
AND version = 'v0.9.0'
ORDER BY image_name;
```

### List exposed images along with cluster details

```sql
SELECT 
    e.cluster_id, 
    c.display_name AS cluster_name, 
    e.name AS image_name, 
    e.registry, 
    e.version
FROM 
    crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images e
JOIN 
    crc_openshift_insights_vulnerabilities_v1_clusters c 
ON 
    e.cluster_id = c.cluster_id
ORDER BY 
    e.name;
```
