package command

import (
	"context"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"fmt"
)

type ubCaptchaKeyGeneratorCmd struct {
	ubCaptchaManager user.UbCaptchaManager
	logger           platform.Logger
}

func (cmd *ubCaptchaKeyGeneratorCmd) Run(ctx context.Context, flags []string) {

	fmt.Println("generating new ub-captcha private and public keys")

	key, err := cmd.ubCaptchaManager.NewKey()
	if err != nil {
		fmt.Println("error in generating new ub-captcha keys")
		fmt.Println("error: " + err.Error())
		return
	}

	//write public & private key to pem files
	err = cmd.ubCaptchaManager.SaveKeyToPemFile(key)
	if err != nil {
		fmt.Println("error in writing pem file to the disk")
		fmt.Println("error: " + err.Error())
		return
	}

	fmt.Println("private.pem and public.pem files have been updated")

}


func NewUbCaptchaKeyGeneratorCmd(ubCaptchaManager user.UbCaptchaManager, logger platform.Logger) ConsoleCommand {
	cmd := &ubCaptchaKeyGeneratorCmd{ubCaptchaManager: ubCaptchaManager, logger: logger}
	return cmd
}
