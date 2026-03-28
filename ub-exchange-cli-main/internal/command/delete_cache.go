package command

import (
	"context"
	"exchange-go/internal/platform"
	"flag"

	"go.uber.org/zap"
)

type deleteCacheCmd struct {
	cache   platform.Cache
	logger  platform.Logger
	pattern string
}

func (cmd *deleteCacheCmd) Run(ctx context.Context, flags []string) {
	cmd.setNeededData(flags)
	cmd.cache.DeleteAll(ctx, cmd.pattern)
}

func (cmd *deleteCacheCmd) setNeededData(flags []string) {
	if len(flags) > 0 {
		pattern := flag.String("pattern", "", "")
		err := flag.CommandLine.Parse(flags)
		if err != nil {
			cmd.logger.Fatal("error in delete cache command", zap.Error(err))
		}
		cmd.pattern = *pattern
	}

}

func NewDeleteCacheCmd(cache platform.Cache, logger platform.Logger) ConsoleCommand {
	return &deleteCacheCmd{
		cache:  cache,
		logger: logger,
	}
}
