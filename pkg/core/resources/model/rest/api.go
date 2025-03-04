package rest

import (
	"fmt"

	"github.com/kumahq/kuma/pkg/core/resources/apis/system"

	"github.com/pkg/errors"

	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/model"
)

type Api interface {
	GetResourceApi(model.ResourceType) (ResourceApi, error)
}

type ResourceApi interface {
	List(mesh string) string
	Item(mesh string, name string) string
}

func NewResourceApi(resType model.ResourceType, path string) ResourceApi {
	switch resType {
	case mesh.MeshType:
		return &meshApi{}
	case system.ZoneType:
		return &zoneApi{}
	default:
		return &resourceApi{CollectionPath: path}
	}
}

type resourceApi struct {
	CollectionPath string
}

func (r *resourceApi) List(mesh string) string {
	return fmt.Sprintf("/meshes/%s/%s", mesh, r.CollectionPath)
}

func (r resourceApi) Item(mesh string, name string) string {
	return fmt.Sprintf("/meshes/%s/%s/%s", mesh, r.CollectionPath, name)
}

type meshApi struct {
}

func (r *meshApi) List(string) string {
	return "/meshes"
}

func (r *meshApi) Item(string, name string) string {
	return fmt.Sprintf("/meshes/%s", name)
}

type zoneApi struct {
}

func (r *zoneApi) List(string) string {
	return "/zones"
}

func (r *zoneApi) Item(string, name string) string {
	return fmt.Sprintf("/zones/%s", name)
}

var _ Api = &ApiDescriptor{}

type ApiDescriptor struct {
	Resources map[model.ResourceType]ResourceApi
}

func (m *ApiDescriptor) GetResourceApi(typ model.ResourceType) (ResourceApi, error) {
	mapping, ok := m.Resources[typ]
	if !ok {
		return nil, errors.Errorf("unknown resource type: %q", typ)
	}
	return mapping, nil
}
