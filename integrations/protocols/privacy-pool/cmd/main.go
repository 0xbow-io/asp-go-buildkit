package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	privacypool "github.com/0xBow-io/asp-go-buildkit/integrations/protocols/privacy-pool"

	erpc "github.com/0xBow-io/asp-go-buildkit/internal/erpc"

	"github.com/spf13/cobra"
)

var observabels = privacypool.Observables()

func findObservable(id string) watcher.Observable {
	for _, o := range observabels {
		if o.ID() == id {
			return o
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use: "play [observable] [rpc] [from]",
	Short: `
 ██████╗ ██████╗ ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗
██╔═══██╗██╔══██╗██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗
██║   ██║██████╔╝███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝
██║   ██║██╔══██╗╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗
╚██████╔╝██████╔╝███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║
 ╚═════╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝
	`,
	Long: `stream state transitions of an observable instance`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err        error
			observable watcher.Observable
			adapter    erpc.Backend
			from       int64
		)
		if len(args) < 4 {
			cmd.Usage()
			os.Exit(1)
		}

		observable = findObservable(args[1])
		if observable == nil {
			fmt.Printf("observable %+v not found \n", args[0])
			cmd.Usage()
			os.Exit(1)
		}

		adapter, err = erpc.NewERPC(args[2])
		if err != nil {
			fmt.Printf("failed to create adapter %+v \n", args[1])
			cmd.Usage()
			os.Exit(1)
		}

		from, err = strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			fmt.Printf("failed to parse from %+v \n", args[3])
			cmd.Usage()
			os.Exit(1)
		}
		end := uint64(6314742)

		out, err := watcher.NewService(
			func(ctx context.Context) (uint64, uint64) {
				return uint64(from), end
			},
			adapter,
		).Watch(observable)

		for _, s := range out {
			event := s.Event()
			if event == nil {
				fmt.Println("failed to extract event !")
				os.Exit(1)
			}
			fmt.Printf("New State --> hash: %+v, event: %+v \n", s.Hash(), event.Format())
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
