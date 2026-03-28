package command

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"flag"
	"strings"

	"go.uber.org/zap"
)

type initializeBalanceCmd struct {
	userRepository        user.Repository
	currencyService       currency.Service
	userBalanceRepository userbalance.Repository
	userBalanceService    userbalance.Service
	logger                platform.Logger
	coinCode              string
}

func (cmd *initializeBalanceCmd) Run(ctx context.Context, flags []string) {
	cmd.setNeededData(flags)
	page := 0
	pageSize := 100
	filters := map[string]interface{}{}

	coin, err := cmd.currencyService.GetCoinByCode(cmd.coinCode)
	if err != nil {
		cmd.logger.Error2("error getting coin", err,
			zap.String("service", "initializeBalanceCmd"),
			zap.String("method", "Run"),
			zap.String("coinCode", cmd.coinCode),
		)
	}

	for {
		users := cmd.userRepository.GetUsersByPagination(int64(page), pageSize, filters)
		if len(users) == 0 {
			break
		}

		var userIds []int
		for _, u := range users {
			userIds = append(userIds, u.ID)
		}

		userBalances := cmd.userBalanceRepository.GetBalancesOfUsersForCoins(userIds, []int64{coin.ID})

		//trying to find the user which has no balance for specific coin
		for _, u := range users {
			exists := false
			for _, ub := range userBalances {
				if u.ID == ub.UserID {
					exists = true
					break
				}
			}

			if !exists {
				ub, err := cmd.userBalanceService.GenerateSingleUserBalanceForCoin(u, coin)
				if err != nil {
					cmd.logger.Error2("error generating user balance", err,
						zap.String("service", "initializeBalanceCmd"),
						zap.String("method", "Run"),
						zap.String("coinCode", coin.Code),
						zap.Int("userID", u.ID),
					)
					continue
				}
				_, err = cmd.userBalanceService.GenerateAddress(ub, u, coin)
				if err != nil {
					cmd.logger.Error2("error generating address", err,
						zap.String("service", "initializeBalanceCmd"),
						zap.String("method", "Run"),
						zap.String("coinCode", coin.Code),
						zap.Int("userID", u.ID),
					)
				}
			}
		}
		page++
	}

}

func (cmd *initializeBalanceCmd) setNeededData(flags []string) {
	coin := flag.String("coin", "", "")
	err := flag.CommandLine.Parse(flags)
	if err != nil {
		cmd.logger.Fatal("error in initializeBalanceCmd", zap.Error(err))
	}
	cmd.coinCode = strings.ToUpper(*coin)
}

func NewInitializedBalanceCmd(userRepository user.Repository, currencyService currency.Service,
	userBalanceRepository userbalance.Repository, userBalanceService userbalance.Service,
	logger platform.Logger) ConsoleCommand {
	cmd := &initializeBalanceCmd{
		userRepository:        userRepository,
		currencyService:       currencyService,
		userBalanceRepository: userBalanceRepository,
		userBalanceService:    userBalanceService,
		logger:                logger,
	}
	return cmd

}
