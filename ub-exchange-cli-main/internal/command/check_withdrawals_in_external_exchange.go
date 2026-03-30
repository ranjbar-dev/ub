package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/payment"
	"exchange-go/internal/platform"

	"go.uber.org/zap"
)

type checkWithdrawalsCmd struct {
	paymentService          payment.Service
	externalExchangeService externalexchange.Service
	internalTransferService payment.InternalTransferService
	configs                 platform.Configs
	logger                  platform.Logger

	status string //just for test purposes we set this as flag
}

func (cmd *checkWithdrawalsCmd) Run(ctx context.Context, flags []string) {
	cmd.logger.Info("start of check withdrawals command")

	withdrawals := cmd.paymentService.GetInProgressWithdrawalsInExternalExchange()
	if len(withdrawals) > 0 {
		firstWithdrawal := withdrawals[0]
		fetchWithdrawalsParams := externalexchange.FetchWithdrawalsParams{
			From: firstWithdrawal.UpdatedAt,
		}

		if cmd.configs.GetEnv() == platform.EnvTest {
			fetchWithdrawalsParams.Status = flags[0]
		}

		externalWithdrawals, err := cmd.externalExchangeService.FetchWithdrawals(fetchWithdrawalsParams)
		if err != nil {
			cmd.logger.Error2("error fetching withdrawals ", err,
				zap.String("service", "checkWithdrawalCmd"),
				zap.String("method", "Run"),
			)
		}

		for _, w := range externalWithdrawals {
			externalWithdrawID := w.ExternalWithdrawID
			for _, withdraw := range withdrawals {
				if withdraw.ExternalExchangeWithdrawID == externalWithdrawID {
					params := payment.UpdatePaymentInExternalExchangeParams{
						PaymentID:          withdraw.PaymentID,
						PaymentExtraInfoID: withdraw.PaymentExtraInfoID,
						TxID:               w.TxID,
						Status:             w.Status,
						Data:               w.Data,
					}

					cmd.paymentService.UpdatePaymentInExternalExchange(params)
					break
				}
			}
		}
	}

	internalTransfers := cmd.internalTransferService.GetFromExternalInProgressTransfers()
	if len(internalTransfers) < 1 {
		cmd.logger.Info("end of check withdrawals command")
		return
	}

	firstInternalTransfer := internalTransfers[0]
	fetchWithdrawalsParams := externalexchange.FetchWithdrawalsParams{
		From: firstInternalTransfer.CreatedAt,
	}

	if cmd.configs.GetEnv() == platform.EnvTest {
		fetchWithdrawalsParams.Status = flags[0]
	}

	externalWithdrawals, err := cmd.externalExchangeService.FetchWithdrawals(fetchWithdrawalsParams)
	if err != nil {
		cmd.logger.Error2("error fetching withdrawals ", err,
			zap.String("service", "checkWithdrawalCmd"),
			zap.String("method", "Run"),
		)
	}

	for _, w := range externalWithdrawals {
		externalWithdrawID := w.ExternalWithdrawID
		for _, it := range internalTransfers {
			id, err := cmd.getIDFromMetadataForInternalTransfer(it.Metadata.String)
			if err != nil {
				cmd.logger.Error2("error geting id from metadata for internal transfer", err,
					zap.String("service", "checkWithdrawalCmd"),
					zap.String("method", "Run"),
				)
				continue
			}

			if id == externalWithdrawID {
				if w.Status == payment.StatusInProgress {
					break
				}

				it.Status = w.Status
				it.TxID = sql.NullString{String: w.TxID, Valid: true}
				err := cmd.internalTransferService.Update(&it)
				if err != nil {
					cmd.logger.Error2("error updatings internalTransfer", err,
						zap.String("service", "checkWithdrawalCmd"),
						zap.String("method", "Run"),
						zap.Int64("internalTransferID", it.ID),
					)
				}
				break
			}
		}
	}

	cmd.logger.Info("end of check withdrawals command")
}

func (cmd *checkWithdrawalsCmd) getIDFromMetadataForInternalTransfer(metadata string) (id string, err error) {
	if metadata != "" {
		md := make(map[string]interface{})
		err = json.Unmarshal([]byte(metadata), &md)
		if err != nil {
			return id, err
		}
		infoField, ok := md["info"]
		if ok {
			id, ok := infoField.(map[string]interface{})["id"]
			if ok {
				return id.(string), nil
			}
		}
	}
	return id, err
}

func NewCheckWithdrawalsCmd(paymentService payment.Service, externalExchangeService externalexchange.Service,
	internalTransferService payment.InternalTransferService, configs platform.Configs, logger platform.Logger) ConsoleCommand {
	return &checkWithdrawalsCmd{
		paymentService:          paymentService,
		externalExchangeService: externalExchangeService,
		internalTransferService: internalTransferService,
		configs:                 configs,
		logger:                  logger,
	}
}
