package watcher

import (
	core "github.com/0xBow-io/asp-go-buildkit/core"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
)

var log = logging.Logger("observer/watcher")

func ConstructModule(cfg *core.Config) fx.Option {
	cfgErr := cfg.Validate()

	return fx.Module("watcher",
		fx.Supply(cfg),
		fx.Error(cfgErr),
		fx.Provide(),
	)
}
