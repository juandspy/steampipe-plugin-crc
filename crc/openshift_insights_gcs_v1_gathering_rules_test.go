package crc

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockGatheringResponse = `
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

func TestDecodeGatheringRules(t *testing.T) {
	var rules gatheringRules
	body := io.NopCloser(strings.NewReader(mockGatheringResponse))
	rules, err := decodeGatheringRules(body)
	assert.NoError(t, err)
	assert.Equal(t, rules.Version, "1.0.1")
	assert.Len(t, rules.Rules, 9)
}
