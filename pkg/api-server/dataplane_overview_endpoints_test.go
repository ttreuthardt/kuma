package api_server_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/kumahq/kuma/api/mesh/v1alpha1"
	api_server "github.com/kumahq/kuma/pkg/api-server"
	config "github.com/kumahq/kuma/pkg/config/api-server"
	mesh_core "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/plugins/resources/memory"
	"github.com/kumahq/kuma/pkg/util/proto"
)

var _ = Describe("Dataplane Overview Endpoints", func() {
	var apiServer *api_server.ApiServer
	var resourceStore store.ResourceStore
	var stop chan struct{}
	t1, _ := time.Parse(time.RFC3339, "2018-07-17T16:05:36.995+00:00")
	BeforeEach(func() {
		resourceStore = memory.NewStore()
		apiServer = createTestApiServer(resourceStore, config.DefaultApiServerConfig(), true)
		client := resourceApiClient{
			address: apiServer.Address(),
			path:    "/meshes",
		}
		stop = make(chan struct{})
		go func() {
			defer GinkgoRecover()
			err := apiServer.Start(stop)
			Expect(err).ToNot(HaveOccurred())
		}()
		waitForServer(&client)
	}, 5)

	AfterEach(func() {
		close(stop)
	})

	BeforeEach(func() {
		err := resourceStore.Create(context.Background(), &mesh_core.MeshResource{}, store.CreateByKey("mesh1", "mesh1"), store.CreatedAt(t1))
		Expect(err).ToNot(HaveOccurred())
	})

	createDpWithInsights := func(name string, dp v1alpha1.Dataplane) {
		dpResource := mesh_core.DataplaneResource{
			Spec: dp,
		}
		err := resourceStore.Create(context.Background(), &dpResource, store.CreateByKey(name, "mesh1"), store.CreatedAt(t1))
		Expect(err).ToNot(HaveOccurred())

		sampleTime, _ := time.Parse(time.RFC3339, "2019-07-01T00:00:00+00:00")
		insightResource := mesh_core.DataplaneInsightResource{
			Spec: v1alpha1.DataplaneInsight{
				Subscriptions: []*v1alpha1.DiscoverySubscription{
					{
						Id:                     "stream-id-1",
						ControlPlaneInstanceId: "cp-1",
						ConnectTime:            proto.MustTimestampProto(sampleTime),
						Status:                 v1alpha1.NewSubscriptionStatus(),
					},
				},
			},
		}
		err = resourceStore.Create(context.Background(), &insightResource, store.CreateByKey(name, "mesh1"))
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		// given
		createDpWithInsights("dp-1", v1alpha1.Dataplane{
			Networking: &v1alpha1.Dataplane_Networking{
				Address: "127.0.0.1",
				Gateway: &v1alpha1.Dataplane_Networking_Gateway{
					Tags: map[string]string{
						"service": "gateway",
					},
				},
			},
		})

		createDpWithInsights("dp-2", v1alpha1.Dataplane{
			Networking: &v1alpha1.Dataplane_Networking{
				Address: "127.0.0.1",
				Inbound: []*v1alpha1.Dataplane_Networking_Inbound{
					{
						Port: 1234,
						Tags: map[string]string{
							"service": "backend",
							"version": "v1",
						},
					},
				},
			},
		})

		createDpWithInsights("dp-3", v1alpha1.Dataplane{
			Networking: &v1alpha1.Dataplane_Networking{
				Address: "127.0.0.1",
				Ingress: &v1alpha1.Dataplane_Networking_Ingress{},
				Inbound: []*v1alpha1.Dataplane_Networking_Inbound{
					{
						Port: 1234,
					},
				},
			},
		})
	})

	dp1Json := `
{
	"type": "DataplaneOverview",
	"name": "dp-1",
	"mesh": "mesh1",
	"creationTime": "2018-07-17T16:05:36.995Z",
	"modificationTime": "2018-07-17T16:05:36.995Z",
	"dataplane": {
		"networking": {
			"address": "127.0.0.1",
			"gateway": {
				"tags": {
					"service": "gateway"
				}
            }
		}
	},
	"dataplaneInsight": {
		"subscriptions": [
			{
				"id": "stream-id-1",
				"controlPlaneInstanceId": "cp-1",
				"connectTime": "2019-07-01T00:00:00Z",
				"status": {
					"total": {},
					"cds": {},
					"eds": {},
					"lds": {},
					"rds": {}
				}
			}
		]
	}
}`

	dp2Json := `
{
	"type": "DataplaneOverview",
	"name": "dp-2",
	"mesh": "mesh1",
	"creationTime": "2018-07-17T16:05:36.995Z",
	"modificationTime": "2018-07-17T16:05:36.995Z",
	"dataplane": {
		"networking": {
			"address": "127.0.0.1",
			"inbound": [
				{
					"port": 1234,
					"tags": {
						"service": "backend",
						"version": "v1"
					}
				}
			]
		}
	},
	"dataplaneInsight": {
		"subscriptions": [
			{
				"id": "stream-id-1",
				"controlPlaneInstanceId": "cp-1",
				"connectTime": "2019-07-01T00:00:00Z",
				"status": {
					"total": {},
					"cds": {},
					"eds": {},
					"lds": {},
					"rds": {}
				}
			}
		]
	}
}`

	dp3Json := `
{
	"type": "DataplaneOverview",
	"name": "dp-3",
	"mesh": "mesh1",
	"creationTime": "2018-07-17T16:05:36.995Z",
	"modificationTime": "2018-07-17T16:05:36.995Z",
	"dataplane": {
		"networking": {
			"address": "127.0.0.1",
			"ingress": {},
			"inbound": [
				{
					"port": 1234
				}
			]
		}
	},
	"dataplaneInsight": {
		"subscriptions": [
			{
				"id": "stream-id-1",
				"controlPlaneInstanceId": "cp-1",
				"connectTime": "2019-07-01T00:00:00Z",
				"status": {
					"total": {},
					"cds": {},
					"eds": {},
					"lds": {},
					"rds": {}
				}
			}
		]
	}
}`

	Describe("On GET", func() {
		It("should return an existing resource", func() {
			// when
			response, err := http.Get("http://" + apiServer.Address() + "/meshes/mesh1/dataplanes+insights/dp-1")
			Expect(err).ToNot(HaveOccurred())

			// then
			Expect(response.StatusCode).To(Equal(200))
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body).To(MatchJSON(dp1Json))
		})

		type testCase struct {
			url          string
			expectedJson string
		}

		DescribeTable("Listing resources filtering by tag",
			func(tc testCase) {
				// when
				response, err := http.Get("http://" + apiServer.Address() + tc.url)
				Expect(err).ToNot(HaveOccurred())

				// then
				Expect(response.StatusCode).To(Equal(200))
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(body)).To(MatchJSON(tc.expectedJson))
			},
			Entry("should list all when no tag is provided", testCase{
				url:          "/meshes/mesh1/dataplanes+insights",
				expectedJson: fmt.Sprintf(`{"total": 3, "items": [%s,%s,%s], "next": null}`, dp1Json, dp2Json, dp3Json),
			}),
			Entry("should list with only one matching tag", testCase{
				url:          "/meshes/mesh1/dataplanes+insights?tag=service:backend",
				expectedJson: fmt.Sprintf(`{"total": 1, "items": [%s], "next": null}`, dp2Json),
			}),
			Entry("should list all with all matching tags", testCase{
				url:          "/meshes/mesh1/dataplanes+insights?tag=service:backend&tag=version:v1",
				expectedJson: fmt.Sprintf(`{"total": 1, "items": [%s], "next": null}`, dp2Json),
			}),
			Entry("should not list when any tag is not matching", testCase{
				url:          "/meshes/mesh1/dataplanes+insights?tag=service:backend&tag=version:v2",
				expectedJson: `{"total": 0, "items": [], "next": null}`,
			}),
			Entry("should list only gateway dataplanes", testCase{
				url:          "/meshes/mesh1/dataplanes+insights?gateway=true",
				expectedJson: fmt.Sprintf(`{"total": 1, "items": [%s], "next": null}`, dp1Json),
			}),
			Entry("should list only ingress dataplanes", testCase{
				url:          "/meshes/mesh1/dataplanes+insights?ingress=true",
				expectedJson: fmt.Sprintf(`{"total": 1, "items": [%s], "next": null}`, dp3Json),
			}),
		)
	})
})
