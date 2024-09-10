package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	privacypool "github.com/0xBow-io/asp-go-buildkit/integrations/protocols/privacy-pool"
	. "github.com/0xBow-io/asp-go-buildkit/integrations/protocols/privacy-pool/cmd/srv"

	erpc "github.com/0xBow-io/asp-go-buildkit/internal/erpc"

	"github.com/spf13/cobra"
)

var observables = privacypool.Observables()

func findObservable(id string) watcher.Observable {
	for _, o := range observables {
		if o.ID() == id {
			return o
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use: "play [observable] [rpc] [from] [range]",
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

		Observe(context.Background(), observable, adapter, uint64(from), 10000,
			5*time.Second,
			privacypool.StateDeserializerFunc)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
