package wallet

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	ValidationURI  = "/api/v1/address/is-valid"
	NewAddressURI  = "/api/v1/address/new"
	NewTxURI       = "/api/v1/tx/new"
	AddressBalanceURI = "/api/v1/address/balance"
	HostEnvKey     = "wallet.host"
	UsernameEnvKey = "wallet.username"
	PasswordEnvKey = "wallet.password"
)

// Service is a client for the external wallet microservice, providing address
// validation, address generation, transaction sending, and balance queries.
type Service interface {
	// IsAddressValid validates whether a cryptocurrency address is valid for the given coin and network.
	IsAddressValid(coin string, address string, network string) (bool, error)
	// GetAddressForUser requests a new deposit address for the given coin and user code.
	GetAddressForUser(coin string, userCode string) (string, error)
	// SendTransaction initiates a cryptocurrency withdrawal transaction via the wallet service.
	SendTransaction(coin, amount, toAddress, network, networkFee string) (string, error)
	// GetAddressBalance returns the balance of a specific address for the given coin and network.
	GetAddressBalance(code string, network string, address string, isHumanReadable bool) (string, error)
}

type service struct {
	authorizationService AuthorizationService
	httpClient           platform.HTTPClient
	logger               platform.Logger
	walletHost           string
	username             string
	password             string
	env                  string
}

type isAddressValidRequestBody struct {
	Coin    string `json:"code"`
	Address string `json:"address"`
	Network string `json:"network"`
}

type IsAddressValidBodyData struct {
	IsValid bool `json:"isValid"`
}

type IsAddressValidResponseBody struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    IsAddressValidBodyData `json:"data"`
}

type newAddressRequestBody struct {
	Coin     string `json:"code"`
	UserCode string `json:"user_id"`
}

type newAddressBodyData struct {
	Address string `json:"address"`
}

type newAddressResponseBody struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    newAddressBodyData `json:"data"`
}

type newTxRequestBody struct {
	Coin      string `json:"code"`
	Amount    string `json:"amount"`
	ToAddress string `json:"to"`
	Network   string `json:"network"`
	Fee       string `json:"fee"`
}

type newTxBodyData struct {
	Tx string `json:"txId"`
}

type newTxResponseBody struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    newTxBodyData `json:"data"`
}

type GetBalanceRequestBody struct {
	Code          string `json:"code"`
	Network       string `json:"network,omitempty"`
	Address       string `json:"address"`
	HumanReadable bool   `json:"human_readable"`
}

type GetBalanceResponseBody struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    GetBalanceDataBody `json:"data"`
}

type GetBalanceDataBody struct {
	Balance string `json:"balance"`
}

func (s *service) IsAddressValid(coin string, address string, network string) (bool, error) {
	ctx := context.Background()
	if s.env != platform.EnvProd {
		return true, nil
	}
	coin = strings.ToUpper(coin)
	url := s.walletHost + ValidationURI
	body := isAddressValidRequestBody{
		Coin:    coin,
		Address: address,
		Network: network,
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	token, err := s.authorizationService.GetToken(ctx)

	if err != nil {
		s.logger.Error2("can not get authorization token", err,
			zap.String("service", "walletService"),
			zap.String("method", "IsAddressValid"),
		)
		return false, fmt.Errorf("can not process your request")

	}
	header["Authorization"] = getAuthorizationHeader(token)
	resp, _, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)

	if err != nil {
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "IsAddressValid"),
		)
		return false, err
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("status code is %d", statusCode)
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "IsAddressValid"),
		)
		return false, fmt.Errorf("can not process your request")
	}

	resBody := IsAddressValidResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		s.logger.Error2("can not marshal wallet response", err,
			zap.String("service", "walletService"),
			zap.String("method", "IsAddressValid"),
		)
		return false, fmt.Errorf("can not process your request")
	}

	isValid := resBody.Status && resBody.Data.IsValid
	return isValid, nil
}

func getAuthorizationHeader(token string) string {
	return "Bearer " + token
}

func (s *service) GetAddressForUser(coin string, userCode string) (string, error) {
	if s.env != platform.EnvProd {
		return coin + "Address", nil
	}

	ctx := context.Background()
	url := s.walletHost + NewAddressURI
	body := newAddressRequestBody{
		Coin:     coin,
		UserCode: userCode,
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	token, err := s.authorizationService.GetToken(ctx)

	if err != nil {
		s.logger.Error2("can not get authorization token ", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressForUser"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	header["Authorization"] = getAuthorizationHeader(token)

	resp, _, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)

	if err != nil {
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressForUser"),
		)
		return "", err
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("status code is %d", statusCode)
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressForUser"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	resBody := newAddressResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		s.logger.Error2("can not marshal wallet response", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressForUser"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	if resBody.Status {
		return resBody.Data.Address, nil
	}
	err = fmt.Errorf("status field of response is false")
	s.logger.Error2("can not marshal wallet response", err,
		zap.String("service", "walletService"),
		zap.String("method", "GetAddressForUser"),
	)
	return "", fmt.Errorf("status field of response is false")
}

func (s *service) SendTransaction(coin string, amount string, toAddress string, network string, networkFee string) (string, error) {
	if s.env != platform.EnvProd {
		return coin + "TxId", nil
	}

	ctx := context.Background()
	url := s.walletHost + NewTxURI
	body := newTxRequestBody{
		Coin:      strings.ToUpper(coin),
		Amount:    amount,
		ToAddress: toAddress,
		Network:   network,
		Fee:       networkFee,
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	token, err := s.authorizationService.GetToken(ctx)

	if err != nil {
		s.logger.Error2("can not get authorization token ", err,
			zap.String("service", "walletService"),
			zap.String("method", "SendTransaction"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	header["Authorization"] = getAuthorizationHeader(token)

	resp, _, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)

	if err != nil {
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "SendTransaction"),
		)
		return "", err
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("status code is %d", statusCode)
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "SendTransaction"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	resBody := newTxResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		s.logger.Error2("can not marshal wallet response", err,
			zap.String("service", "walletService"),
			zap.String("method", "SendTransaction"),
		)
		return "", fmt.Errorf("can not process your request")
	}

	if resBody.Status {
		return resBody.Data.Tx, nil
	}
	err = fmt.Errorf("status field of response is false")
	s.logger.Error2("can not marshal wallet response", err,
		zap.String("service", "walletService"),
		zap.String("method", "GetAddressForUser"),
	)
	return "", fmt.Errorf("status field of response is false")
}

func (s *service) GetAddressBalance(code string, network string, address string, isHumanReadable bool) (string, error) {
	ctx := context.Background()
	if s.env != platform.EnvProd {
		return "0.1", nil
	}

	url := s.walletHost + AddressBalanceURI
	body := GetBalanceRequestBody{
		Code:          code,
		Network:       network,
		Address:       address,
		HumanReadable: isHumanReadable,
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	token, err := s.authorizationService.GetToken(ctx)

	if err != nil {
		s.logger.Error2("can not get authorization token", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressBalance"),
		)
		return "0", fmt.Errorf("can not process your request")

	}
	header["Authorization"] = getAuthorizationHeader(token)
	resp, _, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)

	if err != nil {
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressBalance"),
		)
		return "0", err
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("status code is %d", statusCode)
		s.logger.Error2("can not request to wallet", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressBalance"),
		)
		return "0", fmt.Errorf("can not process your request")
	}

	resBody := GetBalanceResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		s.logger.Error2("can not unmarshal wallet response", err,
			zap.String("service", "walletService"),
			zap.String("method", "GetAddressBalance"),
		)
		return "0", fmt.Errorf("can not process your request")
	}

	return resBody.Data.Balance, nil
}

func NewWalletService(authorizationService AuthorizationService, httpClient platform.HTTPClient, configs platform.Configs, logger platform.Logger) Service {
	walletHost := configs.GetString(HostEnvKey)
	username := configs.GetString(UsernameEnvKey)
	password := configs.GetString(PasswordEnvKey)
	env := configs.GetEnv()

	return &service{authorizationService, httpClient, logger, walletHost, username, password, env}
}
