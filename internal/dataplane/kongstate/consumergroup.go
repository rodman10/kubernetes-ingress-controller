package kongstate

import (
	"github.com/kong/go-kong/kong"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// ConsumerGroup holds a Kong Consumer.
type ConsumerGroup struct {
	kong.ConsumerGroup

	K8sKongConsumerGroup kongv1beta1.KongConsumerGroup
}
