package externalexchange

// OrderFromExternalService provides CRUD operations for orders and trades
// that are synchronized from the external exchange into the local database.
type OrderFromExternalService interface {
	// GetLastOrderFromExternalByPairID retrieves the most recent synced order for a pair.
	GetLastOrderFromExternalByPairID(pairID int64, orderFromExternal *OrderFromExternal) error
	// GetLastTradeFromExternalByPairID retrieves the most recent synced trade for a pair.
	GetLastTradeFromExternalByPairID(pairID int64, tradeFromExternal *TradeFromExternal) error
	// CreateOrder inserts a new order record synced from the external exchange.
	CreateOrder(orderFromExternal *OrderFromExternal) error
	// CreateTrade inserts a new trade record synced from the external exchange.
	CreateTrade(tradeFromExternal *TradeFromExternal) error
	// GetOrderByExternalOrderID looks up a synced order by its external exchange order ID.
	GetOrderByExternalOrderID(externalOrderID int64, orderFromExternal *OrderFromExternal) error
}

type orderFromExternalService struct {
	orderFromExternalRepository OrderFromExternalRepository
	tradeFromExternalRepository TradeFromExternalRepository
}

func (s *orderFromExternalService) GetLastOrderFromExternalByPairID(id int64, orderFromExternal *OrderFromExternal) error {
	return s.orderFromExternalRepository.GetLastOrderFromExternalByPairID(id, orderFromExternal)
}

func (s *orderFromExternalService) GetLastTradeFromExternalByPairID(id int64, tradeFromExternal *TradeFromExternal) error {
	return s.tradeFromExternalRepository.GetLastTradeFromExternalByPairID(id, tradeFromExternal)
}

func (s *orderFromExternalService) CreateOrder(orderFromExternal *OrderFromExternal) error {
	return s.orderFromExternalRepository.Create(orderFromExternal)
}

func (s *orderFromExternalService) CreateTrade(tradeFromExternal *TradeFromExternal) error {
	return s.tradeFromExternalRepository.Create(tradeFromExternal)
}

func (s *orderFromExternalService) GetOrderByExternalOrderID(externalOrderID int64, orderFromExternal *OrderFromExternal) error {
	return s.orderFromExternalRepository.GetOrderByExternalOrderID(externalOrderID, orderFromExternal)
}

func NewOrderFromExternalService(orderFromExternalRepository OrderFromExternalRepository, tradeFromExternalRepository TradeFromExternalRepository) OrderFromExternalService {
	return &orderFromExternalService{
		orderFromExternalRepository: orderFromExternalRepository,
		tradeFromExternalRepository: tradeFromExternalRepository,
	}
}
