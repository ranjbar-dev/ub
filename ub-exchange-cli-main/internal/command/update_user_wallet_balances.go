package command

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"

	"go.uber.org/zap"
)

type ubUpdateUserWalletBalancesCmd struct {
	userService        user.Service
	userBalanceService userbalance.Service
	activeCoins        []currency.Coin
	logger             platform.Logger
}

func (cmd *ubUpdateUserWalletBalancesCmd) Run(ctx context.Context, flags []string) {
	page := int64(0)
	pageSize := 100
	filters := map[string]interface{}{
		"status": user.StatusVerified,
	}
	coinIds := make([]int64, len(cmd.activeCoins))
	for i, coin := range cmd.activeCoins {
		coinIds[i] = coin.ID
	}
	for {
		users := cmd.userService.GetUsersByPagination(page, pageSize, filters)
		if len(users) == 0 {
			break
		}

		for _, u := range users {
			userIds := []int{u.ID}
			userBalances := cmd.userBalanceService.GetBalancesOfUsersForCoins(userIds, coinIds)
			for _, ub := range userBalances {
				if ub.Address.Valid && ub.Address.String != "" {
					blockchainCoinID, blockchainCoinCode, err := cmd.getBlockchainCoinIDAndCodeByCoinID(ub.CoinID)
					if err != nil {
						cmd.logger.Error2("can not get blockchain coin id ", err,
							zap.String("service", "ubUpdateUserWalletBalancesCmd"),
							zap.String("method", "Run"),
							zap.Int64("userBalanceId", ub.ID),
							zap.Int("userId", u.ID),
							zap.Int64("coinID", ub.CoinID),
						)
						continue
					}
					params := userbalance.UpsertUserWalletBalancesParams{
						UserID:             u.ID,
						CoinID:             ub.CoinID,
						CoinCode:           ub.BalanceCoin,
						BlockchainCoinID:   blockchainCoinID,
						BlockchainCoinCode: blockchainCoinCode,
						Address:            ub.Address.String,
					}
					err = cmd.userBalanceService.UpsertUserWalletBalance(params)
					if err != nil {
						cmd.logger.Error2("can not upsert user wallet balance", err,
							zap.String("service", "ubUpdateUserWalletBalancesCmd"),
							zap.String("method", "Run"),
							zap.Int64("userBalanceId", ub.ID),
							zap.Int("userId", u.ID),
						)
						continue
					}
				}

				if ub.OtherAddresses.Valid {
					otherAddresses, err := ub.GetOtherAddresses()
					if err != nil {
						cmd.logger.Error2("can not get other addresses", err,
							zap.String("service", "ubUpdateUserWalletBalancesCmd"),
							zap.String("method", "Run"),
							zap.Int64("userBalanceId", ub.ID),
							zap.Int("userId", u.ID),
						)
						continue
					}

					for _, otherAddress := range otherAddresses {
						blockchainCoinID, err := cmd.getBlockchainCoinIDByBlockchainCoinCode(otherAddress.Code)
						if err != nil {
							cmd.logger.Error2("can not get blockchain coin id", err,
								zap.String("service", "ubUpdateUserWalletBalancesCmd"),
								zap.String("method", "Run"),
								zap.Int64("userBalanceId", ub.ID),
								zap.Int("userId", u.ID),
								zap.String("coinCode", otherAddress.Code),
							)
							continue
						}
						params := userbalance.UpsertUserWalletBalancesParams{
							UserID:             u.ID,
							CoinID:             ub.CoinID,
							CoinCode:           ub.BalanceCoin,
							BlockchainCoinID:   blockchainCoinID,
							BlockchainCoinCode: otherAddress.Code,
							Address:            otherAddress.Address,
						}
						err = cmd.userBalanceService.UpsertUserWalletBalance(params)
						if err != nil {
							cmd.logger.Error2("can not upsert user wallet balance for other networks", err,
								zap.String("service", "ubUpdateUserWalletBalancesCmd"),
								zap.String("method", "Run"),
								zap.Int64("userBalanceId", ub.ID),
								zap.Int("userId", u.ID),
								zap.String("coinCode", otherAddress.Code),
							)
							continue
						}

					}

				}
			}
		}

		if len(users) < pageSize {
			break
		}
		page++
	}

}

func (cmd *ubUpdateUserWalletBalancesCmd) getBlockchainCoinIDAndCodeByCoinID(coinID int64) (int64, string, error) {
	for _, coin := range cmd.activeCoins {
		if coinID == coin.ID {
			if coin.BlockchainNetwork.Valid && coin.BlockchainNetwork.String != "" {
				for _, c := range cmd.activeCoins {
					if coin.BlockchainNetwork.String == c.Code {
						return c.ID, c.Code, nil
					}
				}
			}
			return coinID, coin.Code, nil
		}
	}
	return 0, "", fmt.Errorf("blockchain coin not found")
}

func (cmd *ubUpdateUserWalletBalancesCmd) getBlockchainCoinIDByBlockchainCoinCode(blockchainCoinCode string) (int64, error) {
	for _, coin := range cmd.activeCoins {
		if blockchainCoinCode == coin.Code {
			return coin.ID, nil
		}
	}
	return 0, fmt.Errorf("blockchain coin not found")
}

func NewUbUpdateUserWalletBalances(userService user.Service, userBalanceService userbalance.Service, currencyService currency.Service,
	logger platform.Logger) ConsoleCommand {
	activeCoins := currencyService.GetActiveCoins()
	return &ubUpdateUserWalletBalancesCmd{
		userService:        userService,
		userBalanceService: userBalanceService,
		activeCoins:        activeCoins,
		logger:             logger,
	}
}
