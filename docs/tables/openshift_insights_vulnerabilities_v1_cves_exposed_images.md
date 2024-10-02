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
