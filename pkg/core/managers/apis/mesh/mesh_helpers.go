package mesh

import (
	"context"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core"
	core_ca "github.com/kumahq/kuma/pkg/core/ca"
	mesh_core "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_manager "github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"

	"github.com/pkg/errors"
)

// FetchDefaultMeshIfExists will try to get the default mesh, if present, and return a boolean
// for whether or not it exists
func FetchDefaultMeshIfExists(resManager core_manager.ResourceManager, mesh *mesh_core.MeshResource) (bool, error) {
	key := core_model.ResourceKey{Mesh: core_model.DefaultMesh, Name: core_model.DefaultMesh}

	if err := resManager.Get(context.Background(), mesh, core_store.GetBy(key)); err != nil {
		if core_store.IsResourceNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func CreateDefaultMesh(resManager core_manager.ResourceManager, template mesh_proto.Mesh) error {
	defaultMesh := mesh_core.MeshResource{}

	key := core_model.ResourceKey{Mesh: core_model.DefaultMesh, Name: core_model.DefaultMesh}
	exists, err := FetchDefaultMeshIfExists(resManager, &defaultMesh)
	if err != nil {
		return err
	}

	if !exists {
		defaultMesh.Spec = template
		core.Log.WithName("bootstrap").Info("Creating default mesh from the settings", "mesh", defaultMesh.Spec)

		if err := resManager.Create(context.Background(), &defaultMesh, core_store.CreateBy(key)); err != nil {
			return errors.Wrapf(err, "Failed to create `default` Mesh resource in a given resource store")
		}
	}

	return nil
}

func EnsureEnabledCA(ctx context.Context, caManagers core_ca.Managers, mesh *mesh_core.MeshResource, meshName string) error {
	if mesh.GetEnabledCertificateAuthorityBackend() != nil {
		backend := mesh.GetEnabledCertificateAuthorityBackend()
		caManager, exist := caManagers[backend.Type]
		if !exist { // this should be caught by validator earlier
			return errors.Errorf("CA manager for type %s does not exist", backend.Type)
		}
		if err := caManager.Ensure(ctx, meshName, *backend); err != nil {
			return errors.Wrapf(err, "could not create CA of backend name %s", backend.Name)
		}
	}
	return nil
}
