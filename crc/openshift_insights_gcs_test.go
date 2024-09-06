package crc

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockGatheringResponseV1 = `
{
  "version": "1.0.1",
  "rules": [
    {
      "conditions": [
        {
          "alert": {
            "name": "APIRemovedInNextEUSReleaseInUse"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "api_request_counts_of_resource_from_alert": {
          "alert_name": "APIRemovedInNextEUSReleaseInUse"
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "AlertmanagerFailedReload"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "AlertmanagerFailedReload",
          "container": "alertmanager",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "AlertmanagerFailedToSendAlerts"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "AlertmanagerFailedToSendAlerts",
          "container": "alertmanager",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "KubePodCrashLooping"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "KubePodCrashLooping",
          "previous": true,
          "tail_lines": 20
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "KubePodNotReady"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "KubePodNotReady",
          "tail_lines": 100
        },
        "pod_definition": {
          "alert_name": "KubePodNotReady"
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "PrometheusOperatorSyncFailed"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "PrometheusOperatorSyncFailed",
          "container": "prometheus-operator",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "PrometheusTargetSyncFailure"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "PrometheusTargetSyncFailure",
          "container": "prometheus",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "SamplesImagestreamImportFailing"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "image_streams_of_namespace": {
          "namespace": "openshift-cluster-samples-operator"
        },
        "logs_of_namespace": {
          "namespace": "openshift-cluster-samples-operator",
          "tail_lines": 100
        }
      }
    },
    {
      "conditions": [
        {
          "alert": {
            "name": "ThanosRuleQueueIsDroppingAlerts"
          },
          "type": "alert_is_firing"
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "ThanosRuleQueueIsDroppingAlerts",
          "container": "thanos-ruler",
          "tail_lines": 50
        }
      }
    }
  ]
}
`
const mockGatheringResponseV2 = `
{
  "conditional_gathering_rules": [
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "AlertmanagerFailedReload"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "AlertmanagerFailedReload",
          "container": "alertmanager",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "AlertmanagerFailedToSendAlerts"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "AlertmanagerFailedToSendAlerts",
          "tail_lines": 50,
          "container": "alertmanager"
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "APIRemovedInNextEUSReleaseInUse"
          }
        }
      ],
      "gathering_functions": {
        "api_request_counts_of_resource_from_alert": {
          "alert_name": "APIRemovedInNextEUSReleaseInUse"
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "KubePodCrashLooping"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "KubePodCrashLooping",
          "tail_lines": 20,
          "previous": true
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "KubePodNotReady"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "KubePodNotReady",
          "tail_lines": 100
        },
        "pod_definition": {
          "alert_name": "KubePodNotReady"
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "PrometheusOperatorSyncFailed"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "PrometheusOperatorSyncFailed",
          "tail_lines": 50,
          "container": "prometheus-operator"
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "PrometheusTargetSyncFailure"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "PrometheusTargetSyncFailure",
          "container": "prometheus",
          "tail_lines": 50
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "SamplesImagestreamImportFailing"
          }
        }
      ],
      "gathering_functions": {
        "logs_of_namespace": {
          "namespace": "openshift-cluster-samples-operator",
          "tail_lines": 100
        },
        "image_streams_of_namespace": {
          "namespace": "openshift-cluster-samples-operator"
        }
      }
    },
    {
      "conditions": [
        {
          "type": "alert_is_firing",
          "alert": {
            "name": "ThanosRuleQueueIsDroppingAlerts"
          }
        }
      ],
      "gathering_functions": {
        "containers_logs": {
          "alert_name": "ThanosRuleQueueIsDroppingAlerts",
          "container": "thanos-ruler",
          "tail_lines": 50
        }
      }
    }
  ],
  "container_logs": [
    {
      "namespace": "test-namespace",
      "pod_name_regex": "test.*",
      "messages": [
        "test"
      ],
      "previous": true
    }
  ],
  "version": "1.1.0"
}
`

func TestDecodeGatheringRulesV1(t *testing.T) {
	var rules gatheringRulesV1
	body := io.NopCloser(strings.NewReader(mockGatheringResponseV1))
	rules, err := decodeGatheringRulesV1(body)
	assert.NoError(t, err)
	assert.Equal(t, rules.Version, "1.0.1")
	assert.Len(t, rules.Rules, 9)
}

func TestDecodeGatheringRulesV2(t *testing.T) {
	var rules gatheringRulesV2
	body := io.NopCloser(strings.NewReader(mockGatheringResponseV2))
	rules, err := decodeGatheringRulesV2(body)
	assert.NoError(t, err)
	assert.Equal(t, rules.Version, "1.1.0")
	assert.Len(t, rules.ConditionalGatheringRules, 9)
	assert.Len(t, rules.ContainerLogs, 1)
}
