package store

import (
	"fmt"
	"time"

	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
)

type CreateOptions struct {
	Name         string
	Mesh         string
	CreationTime time.Time
	Owner        core_model.Resource
	Synced       bool
}

type CreateOptionsFunc func(*CreateOptions)

func NewCreateOptions(fs ...CreateOptionsFunc) *CreateOptions {
	opts := &CreateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func CreateBy(key core_model.ResourceKey) CreateOptionsFunc {
	return CreateByKey(key.Name, key.Mesh)
}

func CreateByKey(name, mesh string) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func CreatedAt(creationTime time.Time) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.CreationTime = creationTime
	}
}

func CreateWithOwner(owner core_model.Resource) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Owner = owner
	}
}

func CreateSynced() CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Synced = true
	}
}

type UpdateOptions struct {
	ModificationTime time.Time
	Synced           bool
}

func ModifiedAt(modificationTime time.Time) UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.ModificationTime = modificationTime
	}
}

func UpdateSynced() UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.Synced = true
	}
}

type UpdateOptionsFunc func(*UpdateOptions)

func NewUpdateOptions(fs ...UpdateOptionsFunc) *UpdateOptions {
	opts := &UpdateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type DeleteOptions struct {
	Name string
	Mesh string
}

type DeleteOptionsFunc func(*DeleteOptions)

func NewDeleteOptions(fs ...DeleteOptionsFunc) *DeleteOptions {
	opts := &DeleteOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func DeleteBy(key core_model.ResourceKey) DeleteOptionsFunc {
	return DeleteByKey(key.Name, key.Mesh)
}

func DeleteByKey(name, mesh string) DeleteOptionsFunc {
	return func(opts *DeleteOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

type DeleteAllOptions struct {
	Mesh string
}

type DeleteAllOptionsFunc func(*DeleteAllOptions)

func DeleteAllByMesh(mesh string) DeleteAllOptionsFunc {
	return func(opts *DeleteAllOptions) {
		opts.Mesh = mesh
	}
}

func NewDeleteAllOptions(fs ...DeleteAllOptionsFunc) *DeleteAllOptions {
	opts := &DeleteAllOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type GetOptions struct {
	Name    string
	Mesh    string
	Version string
}

type GetOptionsFunc func(*GetOptions)

func NewGetOptions(fs ...GetOptionsFunc) *GetOptions {
	opts := &GetOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func GetBy(key core_model.ResourceKey) GetOptionsFunc {
	return GetByKey(key.Name, key.Mesh)
}

func GetByKey(name, mesh string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func GetByVersion(version string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Version = version
	}
}

func (g *GetOptions) HashCode() string {
	return fmt.Sprintf("%s:%s", g.Name, g.Mesh)
}

type ListOptions struct {
	Mesh       string
	PageSize   int
	PageOffset string
}

type ListOptionsFunc func(*ListOptions)

func NewListOptions(fs ...ListOptionsFunc) *ListOptions {
	opts := &ListOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func ListByMesh(mesh string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.Mesh = mesh
	}
}

func ListByPage(size int, offset string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.PageSize = size
		opts.PageOffset = offset
	}
}

func (l *ListOptions) HashCode() string {
	return l.Mesh
}
