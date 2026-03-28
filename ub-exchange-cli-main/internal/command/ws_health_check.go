package command

import (
	"context"
	"exchange-go/internal/livedata"
	"exchange-go/internal/platform"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DiffSeconds = 20
	CommandName = "supervisorctl"
	//CommandName = "/home/majid/go/src/exchange-go/ws"
)

type WsHealthCheckCmd struct {
	liveDataService livedata.Service
	logger          platform.Logger
}

func (cmd *WsHealthCheckCmd) Run(ctx context.Context, flags []string) {
	cmd.checkKline(ctx)
	cmd.checkDepth(ctx)
	cmd.checkTrade(ctx)
}

func (cmd *WsHealthCheckCmd) checkKline(ctx context.Context) {
	//using this pair because its data updates more frequently
	pairName := "BTC-USDT"
	klineLastInsertTimeString, err := cmd.liveDataService.GetLastInsertTime(ctx, pairName, livedata.KlineLastInsertTime)
	if err != nil && err != redis.Nil {
		cmd.logger.Error2("error in WsHealthCheckCmd:checkKline", err)
		return
	}
	if err != redis.Nil {
		//the key exists
		klineLastInsertTime, err := strconv.ParseInt(klineLastInsertTimeString, 10, 64)
		if err != nil {
			cmd.logger.Error2("error in WsHealthCheckCmd:checkKline", err)
		}
		nowUnix := time.Now().Unix()
		if nowUnix-klineLastInsertTime > DiffSeconds {
			//it means the process is stopped and we should restart it
			//args := []string{"kline_1m", "kline_5m", "kline_1h", "kline_1d"}
			args := []string{"restart", "kline-stream"}
			err := cmd.restartCommand(ctx, args)
			if err != nil {
				cmd.logger.Error2("error in WsHealthCheckCmd:checkKline", err)
			}
		}
	}
}

func (cmd *WsHealthCheckCmd) checkDepth(ctx context.Context) {
	//using this pair because its data updates more frequently
	pairName := "BTC-USDT"
	depthInsertTimeString, err := cmd.liveDataService.GetLastInsertTime(ctx, pairName, livedata.DepthSnapshotLastInsertTime)
	if err != nil && err != redis.Nil {
		cmd.logger.Error2("error in WsHealthCheckCmd:checkDepth", err)
		return
	}
	if err != redis.Nil {
		//the key exists
		depthInsertTime, err := strconv.ParseInt(depthInsertTimeString, 10, 64)
		if err != nil {
			cmd.logger.Error2("error in WsHealthCheckCmd:checkDepth", err)
		}
		nowUnix := time.Now().Unix()
		if nowUnix-depthInsertTime > DiffSeconds {
			//it means the process is stopped and we should restart it
			args := []string{"restart", "depth-stream"}
			err := cmd.restartCommand(ctx, args)
			if err != nil {
				cmd.logger.Error2("error in WsHealthCheckCmd:checkDepth", err)
			}
		}
	}
}

func (cmd *WsHealthCheckCmd) checkTrade(ctx context.Context) {
	//using this pair because its data updates more frequently
	pairName := "BTC-USDT"
	tradeInsertTimeString, err := cmd.liveDataService.GetLastInsertTime(ctx, pairName, livedata.TradeBookLastInsertTime)
	if err != nil && err != redis.Nil {
		cmd.logger.Error2("error in WsHealthCheckCmd:checkTrade", err)
		return
	}
	if err != redis.Nil {
		//the key exists
		tradeInsertTime, err := strconv.ParseInt(tradeInsertTimeString, 10, 64)
		if err != nil {
			cmd.logger.Error2("error in WsHealthCheckCmd:checkTrade", err)
		}
		nowUnix := time.Now().Unix()
		if nowUnix-tradeInsertTime > DiffSeconds {
			//it means the process is stopped and we should restart it
			//ticker and trade run in same process
			args := []string{"restart", "ticker-trade-stream"}
			err := cmd.restartCommand(ctx, args)
			if err != nil {
				cmd.logger.Error2("error in WsHealthCheckCmd:checkTrade", err)
			}
		}
	}
}

func (cmd *WsHealthCheckCmd) restartCommand(ctx context.Context, args []string) error {
	command := exec.Command(CommandName, args...)
	return command.Start()
}

func NewWsHealthCheckCmd(liveDataService livedata.Service, logger platform.Logger) ConsoleCommand {
	return &WsHealthCheckCmd{
		liveDataService: liveDataService,
		logger:          logger,
	}

}
