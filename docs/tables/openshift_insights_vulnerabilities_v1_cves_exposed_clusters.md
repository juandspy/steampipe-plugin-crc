---
title: "Steampipe Table: openshift_insights_vulnerabilities_v1_cves_exposed_clusters - Query OpenShift Insights Vulnerabilities CVEs Exposed Clusters using SQL"
description: "Allows users to query CVEs Exposed Clusters in OpenShift Insights Vulnerabilities, providing information about clusters affected by specific Common Vulnerabilities and Exposures (CVEs)."
---

# Table: openshift_insights_vulnerabilities_v1_cves_exposed_clusters - Query OpenShift Insights Vulnerabilities CVEs Exposed Clusters using SQL

OpenShift Insights Vulnerabilities is a feature that provides information about potential security vulnerabilities in your OpenShift environment. The CVEs Exposed Clusters component specifically focuses on identifying which clusters are affected by particular Common Vulnerabilities and Exposures (CVEs). This information is crucial for understanding the scope of vulnerability impact and prioritizing remediation efforts.

## Examples

### List clusters exposed to a specific CVE

```sql
SELECT display_name, id, version, provider, status
FROM crc_openshift_insights_vulnerabilities_v1_cves_exposed_clusters
WHERE cve_name = 'CVE-2023-2602'
```

### Get details of clusters exposed to high severity CVEs
This query focuses on clusters affected by high severity CVEs, helping prioritize critical security updates.

```sql
SELECT c.display_name, c.id, c.version, c.provider, c.status, v.severity, v.cvss3_score
FROM crc_openshift_insights_vulnerabilities_v1_cves_exposed_clusters c
JOIN crc_openshift_insights_vulnerabilities_v1_cves v ON c.cve_name = v.synopsis
WHERE v.severity IN ('Important', 'Critical')
ORDER BY v.cvss3_score DESC
```
