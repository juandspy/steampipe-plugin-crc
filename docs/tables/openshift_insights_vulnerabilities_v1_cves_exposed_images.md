---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_cves_exposed_images - Query OpenShift Insights Vulnerabilities CVEs Exposed Images using SQL"
description: "Allows users to query CVEs Exposed Images in OpenShift Insights Vulnerabilities, providing information about container images affected by specific Common Vulnerabilities and Exposures (CVEs)."
---

# Table: openshift_insights_vulnerabilities_v1_cves_exposed_images - Query OpenShift Insights Vulnerabilities CVEs Exposed Images using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift environment. The CVEs Exposed Images component specifically focuses on identifying which container images are affected by particular Common Vulnerabilities and Exposures (CVEs). This information is crucial for understanding the scope of vulnerability impact on your container ecosystem and prioritizing image updates or replacements.

## Examples

### List images exposed to a specific CVE

```sql
SELECT name, registry, version, clusters_exposed
FROM crc_openshift_insights_vulnerabilities_v1_cves_exposed_images
WHERE cve_name = 'CVE-2023-2602'
ORDER BY clusters_exposed DESC
```

### List the exposed images for the top 5 CVEs in terms of severity 

```sql
SELECT 
    e.cve_name, e.name, e.registry, e.version
FROM (
    SELECT 
        synopsis,
        severity,
        cvss3_score
    FROM
        crc_openshift_insights_vulnerabilities_v1_cluster_cves
    WHERE cluster_id = 'a5192f07-c608-40bb-8166-cf012af8c5b2'
    ORDER BY severity DESC
    LIMIT 5
) AS cve, crc_openshift_insights_vulnerabilities_v1_cves_exposed_images AS e
WHERE e.cve_name = cve.synopsis
```
