package command

import (
	"context"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"flag"
	"fmt"

	"go.uber.org/zap"
)

type generateAddressCmd struct {
	userRepository        user.Repository
	userBalanceRepository userbalance.Repository
	userBalanceService    userbalance.Service
	logger                platform.Logger
	userID                string
}

func (cmd *generateAddressCmd) Run(ctx context.Context, flags []string) {
	cmd.setNeededData(flags)
	page := 0
	pageSize := 100
	filters := make(map[string]interface{})
	filters["page"] = page
	filters["pageSize"] = pageSize

	if cmd.userID != "" {
		filters["userId"] = cmd.userID
	}

	for {
		userBalances := cmd.userBalanceRepository.GetBalancesWithoutAddresses(filters)
		if len(userBalances) == 0 {
			break
		}
		for _, ub := range userBalances {
			fmt.Printf("generating address for user Id %d for coin %s", ub.UserID, ub.Coin.Code)
			_, err := cmd.userBalanceService.GenerateAddress(ub, ub.User, ub.Coin)
			if err != nil {
				cmd.logger.Error2("error generating address", err,
					zap.String("service", "generateAddressCmd"),
					zap.String("method", "Run"),
					zap.Int("userID", ub.UserID),
					zap.Int64("coinID", ub.CoinID),
				)
				continue
			}
		}
		page++
		filters["page"] = page
	}

}

func NewGenerateAddressCmd(userRepository user.Repository, userBalanceRepository userbalance.Repository,
	userBalanceService userbalance.Service, logger platform.Logger) ConsoleCommand {
	cmd := &generateAddressCmd{
		userRepository:        userRepository,
		userBalanceRepository: userBalanceRepository,
		userBalanceService:    userBalanceService,
		logger:                logger,
	}
	return cmd

}

func (cmd *generateAddressCmd) setNeededData(flags []string) {
	userID := flag.String("userid", "", "")
	err := flag.CommandLine.Parse(flags)
	if err != nil {
		cmd.logger.Fatal("error in generateAddressCmd", zap.Error(err))
	}
	cmd.userID = *userID
}
