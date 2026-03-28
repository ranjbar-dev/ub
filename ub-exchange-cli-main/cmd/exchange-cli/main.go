package main

import (
	"context"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"fmt"
	"os"

	sarulabsDI "github.com/sarulabs/di"
)

var container sarulabsDI.Container

var commands map[string]command.ConsoleCommand

func main() {
	ctx := context.Background()
	container = di.NewContainer()
	registerCommands()
	if len(os.Args) < 2 {
		fmt.Println("please insert command name")
		os.Exit(1)
	}
	commandName := os.Args[1]
	var flags []string
	if len(os.Args) > 2 {
		flags = os.Args[2:]
	}

	for name, cmd := range commands {
		if name == commandName {
			cmd.Run(ctx, flags)
			return
		}
	}

	fmt.Println("command not found")
	os.Exit(1)
}

func registerCommands() {
	commands = map[string]command.ConsoleCommand{
		//"ws-health-check": container.GetWsHealthCheckCommand(),
		"set-user-level":                 container.Get(di.SetUserLevelCommand).(command.ConsoleCommand),
		"initialize-balance":             container.Get(di.InitializeBalanceCommand).(command.ConsoleCommand),
		"generate-address":               container.Get(di.GenerateAddressCommand).(command.ConsoleCommand),
		"retrieve-open-orders":           container.Get(di.RetrieveOpenOrdersToRedisCommand).(command.ConsoleCommand),
		"submit-bot-orders":              container.Get(di.SubmitBotAggregatedOrderCommand).(command.ConsoleCommand),
		"retrieve-external-orders":       container.Get(di.RetrieveExternalOrdersToRedisCommand).(command.ConsoleCommand),
		"generate-kline-sync":            container.Get(di.GenerateKlineSyncCommand).(command.ConsoleCommand),
		"sync-kline":                     container.Get(di.KlineSyncCommand).(command.ConsoleCommand),
		"update-orders-from-external":    container.Get(di.UpdateOrdersInExternalExchangeCommand).(command.ConsoleCommand),
		"check-withdrawals":              container.Get(di.CheckWithdrawalsInExternalExchangeCommand).(command.ConsoleCommand),
		"ub-captcha-generate-keys":       container.Get(di.UbCaptchaKeyGeneratorCommand).(command.ConsoleCommand),
		"ub-captcha-encryption":          container.Get(di.UbCaptchaEncryptionCommand).(command.ConsoleCommand),
		"ub-captcha-decryption":          container.Get(di.UbCaptchaDecryptionCommand).(command.ConsoleCommand),
		"delete-cache":                   container.Get(di.DeleteCacheCommand).(command.ConsoleCommand),
		"ub-update-user-wallet-balances": container.Get(di.UpdateUserWalletBalancesCommand).(command.ConsoleCommand),
	}
}
