package get

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"

	"github.com/spf13/cobra"

	"github.com/kumahq/kuma/app/kumactl/pkg/output"
	"github.com/kumahq/kuma/app/kumactl/pkg/output/printers"
	rest_types "github.com/kumahq/kuma/pkg/core/resources/model/rest"
	"github.com/kumahq/kuma/pkg/core/resources/store"
)

func newGetDataplaneCmd(pctx *getContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataplane NAME",
		Short: "Show a single Dataplane resource",
		Long:  `Show a single Dataplane resource.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := pctx.CurrentResourceStore()
			if err != nil {
				return err
			}
			name := args[0]
			currentMesh := pctx.CurrentMesh()
			dataplane := &mesh.DataplaneResource{}
			if err := rs.Get(context.Background(), dataplane, store.GetByKey(name, currentMesh)); err != nil {
				if store.IsResourceNotFound(err) {
					return errors.Errorf("No resources found in %s mesh", currentMesh)
				}
				return errors.Wrapf(err, "failed to get mesh %s", currentMesh)
			}
			dataplanes := mesh.DataplaneResourceList{
				Items: []*mesh.DataplaneResource{dataplane},
			}
			switch format := output.Format(pctx.args.outputFormat); format {
			case output.TableFormat:
				return printDataplanes(pctx.Now(), &dataplanes, cmd.OutOrStdout())
			default:
				printer, err := printers.NewGenericPrinter(format)
				if err != nil {
					return err
				}
				return printer.Print(rest_types.From.Resource(dataplane), cmd.OutOrStdout())
			}
		},
	}
	return cmd
}
