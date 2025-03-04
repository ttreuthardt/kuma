package api_server_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	system_proto "github.com/kumahq/kuma/api/system/v1alpha1"
	api_server "github.com/kumahq/kuma/pkg/api-server"
	config "github.com/kumahq/kuma/pkg/config/api-server"
	mesh_core "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/apis/system"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/plugins/resources/memory"
	"github.com/kumahq/kuma/pkg/util/proto"
)

var _ = Describe("Zone Overview Endpoints", func() {
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
		err := resourceStore.Create(context.Background(), &mesh_core.MeshResource{}, store.CreateByKey(core_model.DefaultMesh, core_model.DefaultMesh), store.CreatedAt(t1))
		Expect(err).ToNot(HaveOccurred())
	})

	createZoneWithInsights := func(name string, zone system_proto.Zone) {
		zoneResource := system.ZoneResource{
			Spec: zone,
		}
		err := resourceStore.Create(context.Background(), &zoneResource, store.CreateByKey(name, core_model.DefaultMesh), store.CreatedAt(t1))
		Expect(err).ToNot(HaveOccurred())

		sampleTime, _ := time.Parse(time.RFC3339, "2019-07-01T00:00:00+00:00")
		insightResource := system.ZoneInsightResource{
			Spec: system_proto.ZoneInsight{
				Subscriptions: []*system_proto.KDSSubscription{
					{
						Id:               "stream-id-1",
						GlobalInstanceId: "cp-1",
						ConnectTime:      proto.MustTimestampProto(sampleTime),
						Status:           system_proto.NewSubscriptionStatus(),
					},
				},
			},
		}
		err = resourceStore.Create(context.Background(), &insightResource, store.CreateByKey(name, core_model.DefaultMesh))
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		createZoneWithInsights("zone-1", system_proto.Zone{
			Ingress: &system_proto.Zone_Ingress{Address: "10.20.1.1:10001"},
		})

		createZoneWithInsights("zone-2", system_proto.Zone{
			Ingress: &system_proto.Zone_Ingress{Address: "10.20.1.2:10002"},
		})

		createZoneWithInsights("zone-3", system_proto.Zone{
			Ingress: &system_proto.Zone_Ingress{Address: "10.20.1.3:10003"},
		})
	})

	zone1Json := `
{
 "type": "ZoneOverview",
 "mesh": "default",
 "name": "zone-1",
 "creationTime": "2018-07-17T16:05:36.995Z",
 "modificationTime": "2018-07-17T16:05:36.995Z",
 "zone": {
  "ingress": {
   "address": "10.20.1.1:10001"
  }
 },
 "zoneInsight": {
  "subscriptions": [
   {
    "id": "stream-id-1",
    "globalInstanceId": "cp-1",
    "connectTime": "2019-07-01T00:00:00Z",
    "status": {
     "total": {}
    }
   }
  ]
 }
}`

	zone2Json := `
{
 "type": "ZoneOverview",
 "mesh": "default",
 "name": "zone-2",
 "creationTime": "2018-07-17T16:05:36.995Z",
 "modificationTime": "2018-07-17T16:05:36.995Z",
 "zone": {
  "ingress": {
   "address": "10.20.1.2:10002"
  }
 },
 "zoneInsight": {
  "subscriptions": [
   {
    "id": "stream-id-1",
    "globalInstanceId": "cp-1",
    "connectTime": "2019-07-01T00:00:00Z",
    "status": {
     "total": {}
    }
   }
  ]
 }
}`

	zone3Json := `
{
 "type": "ZoneOverview",
 "mesh": "default",
 "name": "zone-3",
 "creationTime": "2018-07-17T16:05:36.995Z",
 "modificationTime": "2018-07-17T16:05:36.995Z",
 "zone": {
  "ingress": {
   "address": "10.20.1.3:10003"
  }
 },
 "zoneInsight": {
  "subscriptions": [
   {
    "id": "stream-id-1",
    "globalInstanceId": "cp-1",
    "connectTime": "2019-07-01T00:00:00Z",
    "status": {
     "total": {}
    }
   }
  ]
 }
}`

	Describe("On GET", func() {
		It("should return an existing resource", func() {
			// when
			response, err := http.Get("http://" + apiServer.Address() + "/zones+insights/zone-1")
			Expect(err).ToNot(HaveOccurred())

			// then
			Expect(response.StatusCode).To(Equal(200))
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body).To(MatchJSON(zone1Json))
		})

		It("should list resources", func() {
			// when
			response, err := http.Get("http://" + apiServer.Address() + "/zones+insights")
			Expect(err).ToNot(HaveOccurred())

			// then
			Expect(response.StatusCode).To(Equal(200))
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(fmt.Sprintf(`{"total": 3, "items": [%s,%s,%s], "next": null}`, zone1Json, zone2Json, zone3Json)))
		})
	})
})
