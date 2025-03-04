package framework

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kumahq/kuma/pkg/config/core"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	util_net "github.com/kumahq/kuma/pkg/util/net"
)

type PortFwd struct {
	lowFwdPort   uint32
	hiFwdPort    uint32
	localAPIPort uint32
}

type K8sControlPlane struct {
	t          testing.TestingT
	mode       core.CpMode
	name       string
	kubeconfig string
	kumactl    *KumactlOptions
	cluster    *K8sCluster
	portFwd    PortFwd
	verbose    bool
}

func NewK8sControlPlane(t testing.TestingT, mode core.CpMode, clusterName string,
	kubeconfig string, cluster *K8sCluster,
	loPort, hiPort uint32,
	verbose bool) *K8sControlPlane {
	name := clusterName + "-" + mode
	kumactl, _ := NewKumactlOptions(t, name, verbose)
	return &K8sControlPlane{
		t:          t,
		mode:       mode,
		name:       name,
		kubeconfig: kubeconfig,
		kumactl:    kumactl,
		cluster:    cluster,
		portFwd: PortFwd{
			localAPIPort: loPort,
		},
		verbose: verbose,
	}
}

func (c *K8sControlPlane) GetName() string {
	return c.name
}

func (c *K8sControlPlane) GetKubectlOptions(namespace ...string) *k8s.KubectlOptions {
	options := &k8s.KubectlOptions{
		ConfigPath: c.kubeconfig,
	}
	for _, ns := range namespace {
		options.Namespace = ns
		break
	}

	return options
}

func (c *K8sControlPlane) PortForwardKumaCP() error {
	var apiPort uint32
	var err error

	kumacpPods := c.GetKumaCPPods()
	if len(kumacpPods) != 1 {
		return errors.Errorf("Kuma CP pods: %d", len(kumacpPods))
	}

	kumacpPodName := kumacpPods[0].Name

	// API
	apiPort, err = util_net.PickTCPPort("", c.portFwd.lowFwdPort+1, c.portFwd.hiFwdPort)
	if err != nil {
		return errors.Errorf("No free port found in range:  %d - %d", c.portFwd.lowFwdPort, c.portFwd.hiFwdPort)
	}

	c.cluster.PortForwardPod(kumaNamespace, kumacpPodName, apiPort, kumaCPAPIPort)
	c.portFwd.localAPIPort = apiPort

	return nil
}

func (c *K8sControlPlane) GetKumaCPPods() []v1.Pod {
	return k8s.ListPods(c.t,
		c.GetKubectlOptions(kumaNamespace),
		metav1.ListOptions{
			LabelSelector: "app=" + kumaServiceName,
		},
	)
}

func (c *K8sControlPlane) VerifyKumaCtl() error {
	if c.portFwd.localAPIPort == 0 {
		return errors.Errorf("API port not forwarded")
	}

	output, err := c.kumactl.RunKumactlAndGetOutputV(c.verbose, "get", "dataplanes")
	fmt.Println(output)

	return err
}

func (c *K8sControlPlane) VerifyKumaREST() error {
	if c.portFwd.localAPIPort == 0 {
		return errors.Errorf("API port not forwarded")
	}

	return http_helper.HttpGetWithRetryWithCustomValidationE(
		c.t,
		"http://localhost:"+strconv.FormatUint(uint64(c.portFwd.localAPIPort), 10),
		&tls.Config{},
		DefaultRetries,
		DefaultTimeout,
		func(statusCode int, body string) bool {
			return statusCode == http.StatusOK
		},
	)
}

func (c *K8sControlPlane) VerifyKumaGUI() error {
	if c.mode == core.Remote {
		return nil
	}

	return http_helper.HttpGetWithRetryWithCustomValidationE(
		c.t,
		"http://localhost:"+strconv.FormatUint(uint64(c.portFwd.localAPIPort), 10)+"/gui",
		&tls.Config{},
		3,
		DefaultTimeout,
		func(statusCode int, body string) bool {
			return statusCode == http.StatusOK
		},
	)
}

func (c *K8sControlPlane) GetKumaCPLogs() (string, error) {
	logs := ""

	pods := c.GetKumaCPPods()
	if len(pods) < 1 {
		return "", errors.Errorf("no kuma-cp pods found for logs")
	}

	for _, p := range pods {
		log, err := c.cluster.GetPodLogs(p)
		if err != nil {
			return "", err
		}

		logs = logs + "\n >>> " + p.Name + "\n" + log
	}

	return logs, nil
}

func (c *K8sControlPlane) FinalizeAdd() error {
	if err := c.PortForwardKumaCP(); err != nil {
		return err
	}

	kumacpURL := "http://localhost:" + strconv.FormatUint(uint64(c.portFwd.localAPIPort), 10)

	return c.kumactl.KumactlConfigControlPlanesAdd(c.name, kumacpURL)
}

func (c *K8sControlPlane) InstallCP(args ...string) (string, error) {
	return c.kumactl.KumactlInstallCP(c.mode, args...)
}

func (c *K8sControlPlane) InjectDNS() error {
	// store the kumactl environment
	oldEnv := c.kumactl.Env
	c.kumactl.Env["KUBECONFIG"] = c.GetKubectlOptions().ConfigPath

	yaml, err := c.kumactl.RunKumactlAndGetOutput("install", "dns")
	if err != nil {
		return err
	}

	// restore kumactl environment
	c.kumactl.Env = oldEnv

	return k8s.KubectlApplyFromStringE(c.t,
		c.GetKubectlOptions(),
		yaml)
}

// A naive implementation to find the URL where Remote CP exposes its API
func (c *K8sControlPlane) GetKDSServerAddress() string {
	pod := c.GetKumaCPPods()[0]

	return "grpcs://" + pod.Status.HostIP + ":" + strconv.FormatUint(uint64(kdsPort), 10)
}

func (c *K8sControlPlane) GetIngressAddress() string {
	ctx := context.Background()
	cs, err := k8s.GetKubernetesClientFromOptionsE(c.t, c.GetKubectlOptions())
	if err != nil {
		return "invalid"
	}
	ingressSvc, err := cs.CoreV1().Services(kumaNamespace).Get(ctx, "kuma-ingress", metav1.GetOptions{})
	if err != nil {
		return "invalid"
	}
	port := ingressSvc.Spec.Ports[0].NodePort

	nodes, err := cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "invalid"
	}
	// assume that we have single node cluster
	for _, addr := range nodes.Items[0].Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address + ":" + strconv.Itoa(int(port))
		}
	}

	return "invalid"
}

func (c *K8sControlPlane) GetGlobaStatusAPI() string {
	return "http://localhost:" + strconv.FormatUint(uint64(c.portFwd.localAPIPort), 10) + "/status/zones"
}
