package main

import (
	"errors"
	"os/signal"
	"syscall"

	observerBuilder "github.com/0xBow-io/asp-go-buildkit/builder/observer"
	"github.com/spf13/cobra"
)

// Start constructs a CLI command to start Observer Service with the given flags.
// Starts Node daemon. First stopping signal gracefully stops the Node and second terminates it.
// Options passed on start override configuration options only on start and are not persisted in config
func Start(options ...func(*cobra.Command)) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start",
		Short: `Starts the observer service. First stopping signal gracefully stops the service and second terminates it.
		Options passed on start override configuration options only on start and are not persisted in config.`,
		Aliases:      []string{"run", "observer"},
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			ctx := cmd.Context()

			// TODO: Apply configuration for the Observer Service
			cfg := observerBuilder.ObserverConfig(ctx)

			// observerBuilder is a package that constructs the observer service
			// NodeType(ctx), Protocol(ctx), store, &cfg, NodeOptions(ctx)...
			// TODO: Create the observerBuilder pkg
			obs, err := observerBuilder.NewWithConfig(&cfg)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, obs.Close())
			}()

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			err = obs.Start(ctx)
			if err != nil {
				return err
			}

			<-ctx.Done()
			cancel() // ensure we stop reading more signals for start context

			ctx, cancel = signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			return obs.Stop(ctx)
		},
	}
	// Apply each passed option to the command
	for _, option := range options {
		option(cmd)
	}
	return cmd
}
