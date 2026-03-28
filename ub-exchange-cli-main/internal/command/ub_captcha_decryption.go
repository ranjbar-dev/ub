package command

import (
	"context"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"flag"
	"fmt"
)

type ubCaptchaDecryptionCmd struct {
	ubCaptchaManager user.UbCaptchaManager
	logger           platform.Logger
	data             string
}

func (cmd *ubCaptchaDecryptionCmd) Run(ctx context.Context, flags []string) {

	fmt.Println("start decrypting data")

	err := cmd.setNeededData(flags)
	if err != nil {
		fmt.Println("can not find data flag")
		return
	}

	fmt.Println("encrypted data: ")
	fmt.Println(cmd.data)

	plainData, err := cmd.ubCaptchaManager.Decrypt(cmd.data)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}

	fmt.Println("plain data: ")
	fmt.Println(plainData)

}

func (cmd *ubCaptchaDecryptionCmd) setNeededData(flags []string) error {

	if len(flags) < 1 {
		return fmt.Errorf("no data passed")
	}

	data := flag.String("data", "", "")

	err := flag.CommandLine.Parse(flags)
	if err != nil {
		return err
	}

	if len(*data) == 0 {
		return fmt.Errorf("data could not be empty")
	}

	cmd.data = *data

	return nil
}

func NewUbCaptchaDecryptionCmd(ubCaptchaManager user.UbCaptchaManager, logger platform.Logger) ConsoleCommand {
	cmd := &ubCaptchaDecryptionCmd{ubCaptchaManager: ubCaptchaManager, logger: logger}
	return cmd
}
