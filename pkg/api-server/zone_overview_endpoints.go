package api_server

import (
	"context"

	"github.com/emicklei/go-restful"
	"github.com/golang/protobuf/proto"

	system_proto "github.com/kumahq/kuma/api/system/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/system"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/model/rest"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	rest_errors "github.com/kumahq/kuma/pkg/core/rest/errors"
)

type zoneOverviewEndpoints struct {
	publicURL  string
	resManager manager.ResourceManager
}

func (r *zoneOverviewEndpoints) addFindEndpoint(ws *restful.WebService) {
	ws.Route(ws.GET("/zones+insights/{name}").To(r.inspectZone).
		Doc("Inspect a zone").
		Param(ws.PathParameter("name", "Name of a zone").DataType("string")).
		Returns(200, "OK", nil).
		Returns(404, "Not found", nil))
}

func (r *zoneOverviewEndpoints) addListEndpoint(ws *restful.WebService) {
	ws.Route(ws.GET("/zones+insights").To(r.inspectZones).
		Doc("Inspect all zones").
		Returns(200, "OK", nil))
}

func (r *zoneOverviewEndpoints) inspectZone(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	overview, err := r.fetchOverview(request.Request.Context(), name)
	if err != nil {
		rest_errors.HandleError(response, err, "Could not retrieve a zone overview")
		return
	}

	res := rest.From.Resource(overview)
	if err := response.WriteAsJson(res); err != nil {
		rest_errors.HandleError(response, err, "Could not retrieve a zone overview")
	}
}

func (r *zoneOverviewEndpoints) fetchOverview(ctx context.Context, name string) (*system.ZoneOverviewResource, error) {
	zone := system.ZoneResource{}
	if err := r.resManager.Get(ctx, &zone, store.GetByKey(name, core_model.DefaultMesh)); err != nil {
		return nil, err
	}

	insight := system.ZoneInsightResource{}
	err := r.resManager.Get(ctx, &insight, store.GetByKey(name, core_model.DefaultMesh))
	if err != nil && !store.IsResourceNotFound(err) { // It's fine to have zone without insight
		return nil, err
	}

	return &system.ZoneOverviewResource{
		Meta: zone.Meta,
		Spec: system_proto.ZoneOverview{
			Zone:        proto.Clone(&zone.Spec).(*system_proto.Zone),
			ZoneInsight: proto.Clone(&insight.Spec).(*system_proto.ZoneInsight),
		},
	}, nil
}

func (r *zoneOverviewEndpoints) inspectZones(request *restful.Request, response *restful.Response) {
	page, err := pagination(request)
	if err != nil {
		rest_errors.HandleError(response, err, "Could not retrieve dataplane overviews")
		return
	}

	overviews, err := r.fetchOverviews(request.Request.Context(), page, core_model.DefaultMesh)
	if err != nil {
		rest_errors.HandleError(response, err, "Could not retrieve dataplane overviews")
		return
	}

	// pagination is not supported yet so we need to override pagination total items after retaining dataplanes
	overviews.GetPagination().SetTotal(uint32(len(overviews.Items)))
	restList := rest.From.ResourceList(&overviews)
	next, err := nextLink(request, r.publicURL, &overviews)
	if err != nil {
		rest_errors.HandleError(response, err, "Could not list dataplane overviews")
		return
	}
	restList.Next = next
	if err := response.WriteAsJson(restList); err != nil {
		rest_errors.HandleError(response, err, "Could not list dataplane overviews")
	}
}

func (r *zoneOverviewEndpoints) fetchOverviews(ctx context.Context, p page, meshName string) (system.ZoneOverviewResourceList, error) {
	zones := system.ZoneResourceList{}
	if err := r.resManager.List(ctx, &zones, store.ListByMesh(meshName), store.ListByPage(p.size, p.offset)); err != nil {
		return system.ZoneOverviewResourceList{}, err
	}

	// we cannot paginate insights since there is no guarantee that the elements will be the same as dataplanes
	insights := system.ZoneInsightResourceList{}
	if err := r.resManager.List(ctx, &insights, store.ListByMesh(meshName)); err != nil {
		return system.ZoneOverviewResourceList{}, err
	}

	return system.NewZoneOverviews(zones, insights), nil
}
