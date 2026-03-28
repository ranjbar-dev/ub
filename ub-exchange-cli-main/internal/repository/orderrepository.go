package repository

import (
	"exchange-go/internal/order"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	db *gorm.DB
}

func (or *orderRepository) GetOrdersByIds(orderIds []int64) []order.Order {
	var orders []order.Order
	or.db.Joins("User").Joins("Pair").Find(&orders, "orders.id IN ?", orderIds)
	return orders
}

func (or *orderRepository) GetOrdersDataByIdsWithJoinUsingTx(tx *gorm.DB, orderIds []int64) []order.MatchingNeededQueryFields {
	var result []order.MatchingNeededQueryFields
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("orders").Where("orders.id IN ?", orderIds).Select("" +
		"orders.id as OrderID," +
		"orders.price as Price," +
		"orders.type as OrderType," +
		"orders.exchange_type as OrderExchangeType," +
		"orders.path as Path," +
		"orders.status as Status," +
		"orders.current_market_price as MarketPrice," +
		"orders.created_at as CreatedAt," +
		"orders.creator_user_id as UserID," +
		"orders.demanded_money_amount as DemandedAmount," +
		"orders.payed_by_money_amount as PayedByAmount," +
		"orders_extra_info.user_agent_info as UserAgentInfo").
		Joins("join orders_extra_info on orders.extra_info_id = orders_extra_info.id").
		Scan(&result)
	return result
}

func (or *orderRepository) GetOrderByID(id int64, o *order.Order) error {
	return or.db.Joins("Pair").Where(order.Order{ID: id}).First(o).Error
}

func (or *orderRepository) GetOrderByIDUsingTx(tx *gorm.DB, id int64, o *order.Order) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(order.Order{ID: id}).First(o).Error
}

func (or *orderRepository) GetUserOpenOrders(userID int, pairID int64) []order.Order {
	var orders []order.Order
	filter := order.Order{
		UserID: userID,
		Status: order.StatusOpen,
	}
	if pairID != 0 {
		filter.PairID = pairID
	}

	or.db.Joins("Pair").Where(filter).Find(&orders)
	return orders
}

func (or *orderRepository) GetOrdersAncestors(ancestorsIds []int64) []order.Order {
	var orders []order.Order
	or.db.Joins("Pair").Find(&orders, "orders.id IN ?", ancestorsIds)
	return orders
}

func (or *orderRepository) GetLeafOrders(filters order.HistoryFilters) []order.HistoryNeededField {
	pageSize := 50
	var results []order.HistoryNeededField
	subQuery := or.db.Table("orders as o2").Select("DISTINCT(o2.parent_id)").Where("o2.parent_id IS NOT NULL AND o2.parent_id <= o.id")
	mainQuery := or.db.Table("orders as o").
		Joins("join pair_currencies as pairs on o.pair_currency_id = pairs.id").
		Select("o.id as OrderId, o.path as Path,o.created_at as CreatedAt").
		Where("o.creator_user_id =?", filters.UserID).
		Where("o.status <> ?", order.StatusOpen)

	if filters.IsFullHistory {
		if filters.LastID != 0 {
			subQuery.Where("o2.parent_id < ?", filters.LastID)
			mainQuery.Where("o.id < ?", filters.LastID)
		}

		if filters.PageSize != 0 {
			pageSize = filters.PageSize
		}
		mainQuery.Limit(pageSize)
	}

	if filters.Type != "" {
		mainQuery.Where("o.type =?", filters.Type)
	}

	if filters.Hide {
		mainQuery.Where("o.status <> ?", order.StatusCanceled)
	}
	//
	if filters.StartDate != "" {
		mainQuery.Where("o.created_at > ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		mainQuery.Where("o.updated_at <= ?", filters.EndDate)
	}

	if filters.IsFastExchange != nil {
		isFastExchange := *filters.IsFastExchange
		mainQuery.Where("o.is_fast_exchange = ?", isFastExchange)
	}

	//filter based on pair_id or pair_name. one of them will apply
	if filters.PairID != 0 {
		mainQuery.Where("o.pair_currency_id = ?", filters.PairID)
	} else if filters.PairName != "" {
		//prepare pair_name filter
		pairName := or.preparePairName(filters.PairName)
		if pairName != "" {
			mainQuery.Where("pairs.name LIKE ?", pairName)
		}
	}

	if filters.DependentCoinID != 0 {
		mainQuery.Where("pairs.dependent_currency_id = ?", filters.DependentCoinID)
	}

	if filters.BasisCoinID != 0 {
		mainQuery.Where("pairs.basis_currency_id = ?", filters.BasisCoinID)
	}

	mainQuery.Order("o.id desc")

	mainQuery.Where("(o.id NOT IN (?))", subQuery).Scan(&results)

	return results
}

//This method change the pairName based on this formats: `BTC-USDT`, `BTC-ALL`, `ALL-USDT`, `ALL-ALL`
// to this formats that can use in sql query: "BTC-USDT", "BTC-%", "%-USDT", ""
func (or *orderRepository) preparePairName(pairName string) string {
	coinNames := strings.Split(strings.ToUpper(pairName), "-")
	if coinNames[0] == "ALL" && coinNames[1] == "ALL" {
		return ""
	}

	if coinNames[0] == "ALL" {
		coinNames[0] = "%"
	}

	if coinNames[1] == "ALL" {
		coinNames[1] = "%"
	}

	return coinNames[0] + "-" + coinNames[1]
}

func (or *orderRepository) GetUserTradedOrders(filters order.TradeHistoryFilters) []order.Order {
	pageSize := 50
	var orders []order.Order
	q := or.db.Joins("Pair").
		Where("orders.status =?", order.StatusFilled).
		Where("orders.creator_user_id =?", filters.UserID)

	if filters.StartDate != "" {
		q.Where("orders.created_at > ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		q.Where("orders.created_at <= ?", filters.EndDate)
	}

	if filters.Type != "" {
		q.Where("orders.type =?", filters.Type)
	}

	//filter based on pair_id or pair_name. one of them will apply
	if filters.PairID != 0 {
		q.Where("orders.pair_currency_id = ?", filters.PairID)
	} else if filters.PairName != "" {
		//prepare pair_name filter
		pairName := or.preparePairName(filters.PairName)
		if pairName != "" {
			q.Where("Pair.name LIKE ?", pairName)
		}
	}

	if filters.DependentCoinID != 0 {
		q.Where("Pair.dependent_currency_id = ?", filters.DependentCoinID)
	}

	if filters.BasisCoinID != 0 {
		q.Where("Pair.basis_currency_id = ?", filters.BasisCoinID)
	}

	q.Order("orders.id desc")

	if filters.IsFullHistory {
		if filters.LastID != 0 {
			q.Where("orders.id < ?", filters.LastID)
		}

		if filters.PageSize != 0 {
			pageSize = filters.PageSize
		}

		q.Limit(pageSize)

	}

	q.Find(&orders)
	return orders
}

func (or *orderRepository) GetUserOrderDetailsByID(id int64, userID int) []order.Order {
	idString := strconv.FormatInt(id, 10)
	var orders []order.Order
	q := or.db.Joins("Pair").
		Where("orders.status = ?", order.StatusFilled).
		Where("orders.creator_user_id = ?", userID).
		Where("orders.path like ?", ""+idString+",%")
	q.Find(&orders)
	return orders
}

func (or *orderRepository) GetOpenOrders(date string) []order.Order {
	var orders []order.Order
	//or.db.Joins("User").Joins("Pair").Find(&orders, "orders.id IN ?", orderIds)
	q := or.db.Joins("Pair").Where("orders.status = ? ", order.StatusOpen)
	if date != "" {
		fromDate, err := time.Parse("2006-01-02 15:04:05", date)
		if err == nil {
			q.Where("orders.created_at >= ?", fromDate)
		}
	}

	q.Order("orders.id asc").Find(&orders)
	return orders
}

func NewOrderRepository(db *gorm.DB) order.Repository {
	return &orderRepository{db}
}
