package e2e_test

import (
	"fmt"
	"reflect"

	"github.com/gruntwork-io/terratest/modules/retry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/kumahq/kuma/pkg/config/core"
	. "github.com/kumahq/kuma/test/framework"
	"github.com/kumahq/kuma/test/framework/deployments/tracing"
)

var _ = Describe("Tracing Universal", func() {

	meshWithTracing := func(zipkinURL string) string {
		return fmt.Sprintf(`
type: Mesh
name: default
tracing:
  defaultBackend: zipkin
  backends:
  - name: zipkin
    type: zipkin
    conf:
      url: %s
`, zipkinURL)
	}

	traceAll := `
type: TrafficTrace
name: traffic-trace-all
mesh: default
selectors:
- match:
   kuma.io/service: "*"
`

	var cluster Cluster

	BeforeEach(func() {
		cluster = NewUniversalCluster(NewTestingT(), Kuma1, Silent)

		err := NewClusterSetup().
			Install(Kuma(core.Standalone)).
			Install(EchoServerUniversal()).
			Install(DemoClientUniversal()).
			Install(tracing.Install()).
			Setup(cluster)
		Expect(err).ToNot(HaveOccurred())
		err = cluster.VerifyKuma()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(cluster.DeleteKuma()).To(Succeed())
		Expect(cluster.DismissCluster()).To(Succeed())
	})

	It("should emit traces to jaeger", func() {
		// given TrafficTrace and mesh with tracing backend
		err := YamlUniversal(meshWithTracing(tracing.From(cluster).ZipkinCollectorURL()))(cluster)
		Expect(err).ToNot(HaveOccurred())
		err = YamlUniversal(traceAll)(cluster)
		Expect(err).ToNot(HaveOccurred())

		retry.DoWithRetry(cluster.GetTesting(), "check traced services", DefaultRetries, DefaultTimeout, func() (string, error) {
			// when client sends requests to server
			_, _, err := cluster.Exec("", "", "demo-client", "curl", "-v", "-m", "3", "localhost:4001")
			if err != nil {
				return "", err
			}

			// then traces are published
			services, err := tracing.From(cluster).TracedServices()
			if err != nil {
				return "", err
			}

			expectedServices := []string{"demo-client", "echo-server_kuma-test_svc_8080", "jaeger-query"}
			if !reflect.DeepEqual(services, expectedServices) {
				return "", errors.Errorf("services not traced. Expected %q, got %q", expectedServices, services)
			}
			return "", nil
		})
	})
})
