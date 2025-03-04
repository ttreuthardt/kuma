package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/kumahq/kuma/app/kumactl/cmd/apply"
	"github.com/kumahq/kuma/app/kumactl/cmd/completion"
	"github.com/kumahq/kuma/app/kumactl/cmd/config"
	"github.com/kumahq/kuma/app/kumactl/cmd/delete"
	"github.com/kumahq/kuma/app/kumactl/cmd/generate"
	"github.com/kumahq/kuma/app/kumactl/cmd/get"
	"github.com/kumahq/kuma/app/kumactl/cmd/inspect"
	"github.com/kumahq/kuma/app/kumactl/cmd/install"
	kumactl_cmd "github.com/kumahq/kuma/app/kumactl/pkg/cmd"
	kumactl_config "github.com/kumahq/kuma/app/kumactl/pkg/config"
	kumactl_errors "github.com/kumahq/kuma/app/kumactl/pkg/errors"
	"github.com/kumahq/kuma/pkg/api-server/types"
	kuma_cmd "github.com/kumahq/kuma/pkg/cmd"
	"github.com/kumahq/kuma/pkg/cmd/version"
	"github.com/kumahq/kuma/pkg/core"
	kuma_log "github.com/kumahq/kuma/pkg/log"
	kuma_version "github.com/kumahq/kuma/pkg/version"
	_ "github.com/kumahq/kuma/pkg/xds/envoy" // import Envoy protobuf definitions so (un)marshalling Envoy protobuf works
)

var (
	kumactlLog       = core.Log.WithName("kumactl")
	kumaBuildVersion *types.IndexResponse
)

// newRootCmd represents the base command when called without any subcommands.
func NewRootCmd(root *kumactl_cmd.RootContext) *cobra.Command {
	args := struct {
		logLevel string
	}{}
	cmd := &cobra.Command{
		Use:   "kumactl",
		Short: "Management tool for Kuma",
		Long:  `Management tool for Kuma.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			client, err := root.CurrentApiClient()
			if err != nil {
				kumactlLog.Error(err, "Unable to get index client")
			} else {
				kumaBuildVersion, _ = client.GetVersion()
			}
			level, err := kuma_log.ParseLogLevel(args.logLevel)
			if err != nil {
				return err
			}
			core.SetLogger(core.NewLogger(level))

			// once command line flags have been parsed,
			// avoid printing usage instructions
			cmd.SilenceUsage = true
			if kumaBuildVersion != nil && kumaBuildVersion.Version != kuma_version.Build.Version {
				cmd.Println("Warning: Your kumactl version is " + kuma_version.Build.Version + " which is not the same as " + kumaBuildVersion.Version + " for CP. Update your kumactl.")
			}
			if root.IsFirstTimeUsage() {
				root.Runtime.Config = kumactl_config.DefaultConfiguration()
				if err := root.SaveConfig(); err != nil {
					return err
				}
			}
			return root.LoadConfig()
		},
	}
	// root flags
	cmd.PersistentFlags().StringVar(&root.Args.ConfigFile, "config-file", "", "path to the configuration file to use")
	cmd.PersistentFlags().StringVarP(&root.Args.Mesh, "mesh", "m", "default", "mesh to use")
	cmd.PersistentFlags().StringVar(&args.logLevel, "log-level", kuma_log.OffLevel.String(), kuma_cmd.UsageOptions("log level", kuma_log.OffLevel, kuma_log.InfoLevel, kuma_log.DebugLevel))
	// sub-commands
	cmd.AddCommand(apply.NewApplyCmd(root))
	cmd.AddCommand(completion.NewCompletionCommand(root))
	cmd.AddCommand(config.NewConfigCmd(root))
	cmd.AddCommand(delete.NewDeleteCmd(root))
	cmd.AddCommand(generate.NewGenerateCmd(root))
	cmd.AddCommand(get.NewGetCmd(root))
	cmd.AddCommand(inspect.NewInspectCmd(root))
	cmd.AddCommand(install.NewInstallCmd(root))
	cmd.AddCommand(version.NewVersionCmd())
	kumactl_cmd.WrapRunnables(cmd, kumactl_errors.FormatErrorWrapper)
	return cmd
}

func DefaultRootCmd() *cobra.Command {
	return NewRootCmd(kumactl_cmd.DefaultRootContext())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := DefaultRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
