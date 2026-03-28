package order

import (
	"context"
	"errors"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *service) CreateOrder(u *user.User, params CreateOrderParams) (apiResponse response.APIResponse, statusCode int) {
	ctx := context.Background()

	if u.Status != user.StatusVerified {
		return response.Error("email status is not verified", http.StatusUnprocessableEntity, nil)
	}

	pair, err := s.currencyService.GetPairByID(params.PairID)
	if err != nil {
		return response.Error("pair currency id is not valid", http.StatusUnprocessableEntity, nil)
	}

	userConfig, err := s.userConfigService.GetUserConfig(u.ID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error in finding userConfig", err,
			zap.String("service", "orderService"),
			zap.String("method", "CreateOrder"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if err == nil && userConfig.ID > 0 {
		if userConfig.IsReadOnly {
			return response.Error("this account is in read only mode", http.StatusUnprocessableEntity, nil)
		}
	}

	hasPermission := s.userPermissionManager.IsPermissionGrantedToUserFor(*u, user.PermissionExchange)

	if !hasPermission {
		return response.Error("permission is not granted to place the order", http.StatusUnprocessableEntity, nil)
	}

	currentPrice, err := s.priceGenerator.GetPrice(ctx, pair.Name)
	if err != nil {
		s.logger.Error2("error in getting currentPrice", err,
			zap.String("service", "orderService"),
			zap.String("method", "CreateOrder"),
			zap.String("pairName", pair.Name),
		)
		return response.Error("can not get price right now", http.StatusUnprocessableEntity, nil)
	}

	// todo  because all we have now are instant should be removed later
	params.IsInstant = true

	data := CreateRequiredData{
		User:           u,
		Pair:           &pair,
		Amount:         params.Amount,
		OrderType:      strings.ToUpper(params.Type),
		ExchangeType:   strings.ToUpper(params.ExchangeType),
		Price:          params.Price,
		StopPointPrice: params.StopPointPrice,
		UserAgentInfo:  params.UserAgentInfo,
		CurrentPrice:   currentPrice,
		IsInstant:      params.IsInstant,
		IsFastExchange: params.IsFastExchange,
	}

	o, err := s.orderCreateManager.CreateOrder(data)
	if err != nil {
		msg := "error in creating order"
		s.logger.Error2("error in create order", err,
			zap.String("service", "orderService"),
			zap.String("method", "CreateOrder"),
			zap.String("amount", params.Amount),
			zap.String("price", params.Price),
			zap.Int64("pairId", params.PairID),
		)

		if errors.Is(err, platform.OrderCreateValidationError{}) {
			msg = err.Error()
		}

		return response.Error(msg, http.StatusUnprocessableEntity, nil)
	}

	//adding other fields to order
	o.Pair = pair
	o.User = *u

	//preparing response
	price := ""
	if o.Price.Valid {
		price = o.Price.String
	}
	res := CreateOrderResponse{
		ID:        o.ID,
		CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
		Price:     price,
	}

	go s.eventsHandler.HandleOrderCreation(*o, false)

	return response.Success(res, "")
}
