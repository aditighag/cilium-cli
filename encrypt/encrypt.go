// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package encrypt

import (
	"context"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium-cli/k8s"
)

// Parameters contains options for CLI
type Parameters struct {
	CiliumNamespace  string
	AgentPodSelector string
	NodeName         string
	PerNodeDetails   bool
	Writer           io.Writer
	WaitDuration     time.Duration
	Output           string
}

// Status is used to get encrypt status from cilium agents
type Status struct {
	client *k8s.Client
	params Parameters
}

// NewStatus returns new encrypt.Status struct
func NewStatus(client *k8s.Client, p Parameters) *Status {
	return &Status{
		client: client,
		params: p,
	}
}

// fetchCiliumPods returns slice of cilium agent pods.
// If option NodeName is specified then only that nodes' cilium-agent
// pod is returned else all cilium-agents in the cluster are returned.
func (s *Status) fetchCiliumPods(ctx context.Context) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: s.params.AgentPodSelector}
	if s.params.NodeName != "" {
		opts.FieldSelector = fmt.Sprintf("spec.nodeName=%s", s.params.NodeName)
	}

	pods, err := s.client.ListPods(ctx, s.params.CiliumNamespace, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to list Cilium pods: %w", err)
	}
	return pods.Items, nil
}
