package controllers

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	kube_core "k8s.io/api/core/v1"
	kube_client "sigs.k8s.io/controller-runtime/pkg/client"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	mesh_k8s "github.com/kumahq/kuma/pkg/plugins/resources/k8s/native/api/v1alpha1"
)

func (p *PodConverter) OutboundInterfacesFor(pod *kube_core.Pod, others []*mesh_k8s.Dataplane) ([]*mesh_proto.Dataplane_Networking_Outbound, error) {
	var outbounds []*mesh_proto.Dataplane_Networking_Outbound
	endpoints := endpointsByService(others)
	for _, serviceTag := range endpoints.Services() {
		service, port, err := p.k8sService(serviceTag)
		if err != nil {
			converterLog.Error(err, "could not get K8S Service for service tag")
			continue // one invalid Dataplane definition should not break the entire mesh
		}
		if isHeadlessService(service) {
			// Generate outbound listeners for every endpoint of services.
			for _, endpoint := range endpoints[serviceTag] {
				if endpoint.Address == pod.Status.PodIP {
					continue // ignore generating outbound for itself, otherwise we've got a conflict with inbound
				}
				outbounds = append(outbounds, &mesh_proto.Dataplane_Networking_Outbound{
					Address: endpoint.Address,
					Port:    endpoint.Port,
					Tags: map[string]string{
						mesh_proto.ServiceTag:  serviceTag,
						mesh_proto.InstanceTag: endpoint.Instance,
					},
				})
			}
		} else {
			// generate outbound based on ClusterIP. Transparent Proxy will work only if DNS name that resolves to ClusterIP is used
			outbounds = append(outbounds, &mesh_proto.Dataplane_Networking_Outbound{
				Address: service.Spec.ClusterIP,
				Port:    port,
				Tags: map[string]string{
					mesh_proto.ServiceTag: serviceTag,
				},
			})
		}
	}
	return outbounds, nil
}

func isHeadlessService(svc *kube_core.Service) bool {
	return svc.Spec.ClusterIP == "None"
}

func (p *PodConverter) k8sService(serviceTag string) (*kube_core.Service, uint32, error) {
	name, ns, port, err := ParseService(serviceTag)
	if err != nil {
		return nil, 0, errors.Errorf("failed to parse `service` host %q as FQDN", serviceTag)
	}

	svc := &kube_core.Service{}
	svcKey := kube_client.ObjectKey{Namespace: ns, Name: name}
	if err := p.ServiceGetter.Get(context.Background(), svcKey, svc); err != nil {
		return nil, 0, errors.Wrapf(err, "failed to get Service %q", svcKey)
	}
	return svc, port, nil
}

func ParseService(host string) (name string, namespace string, port uint32, err error) {
	// split host into <name>_<namespace>_svc_<port>
	segments := strings.Split(host, "_")
	if len(segments) != 4 {
		return "", "", 0, errors.Errorf("service tag in unexpected format")
	}
	p, err := strconv.Atoi(segments[3])
	if err != nil {
		return "", "", 0, err
	}
	port = uint32(p)
	name, namespace = segments[0], segments[1]
	return
}
