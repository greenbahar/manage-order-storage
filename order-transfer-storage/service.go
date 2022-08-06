package order_transfer_storage

import (
	"fmt"
	"github.com/greenbahar/manage-order-storage/order"
)

type Service interface {
	Run() error
	PersistOrdersFromRedisToMysql() error
}

type RedisRepository interface {
	SetOrder(order *order.Order) error
	GetOrder(key string) (*order.Order, error)
	GetOrders(key []string) ([]*order.Order, error)
	GetOnePageOfKeys(startIndex, rangeIndex int) ([]string, error)
}

type MysqlRepository interface {
	InsertOrder(orders []*order.Order) error
}

type service struct {
	redisRepo RedisRepository
	mysqlRepo MysqlRepository
}

func NewService(redisRepo RedisRepository, mysqlRepo MysqlRepository) Service {
	return &service{
		redisRepo: redisRepo,
		mysqlRepo: mysqlRepo,
	}
}

func (s *service) Run() error {
	return s.PersistOrdersFromRedisToMysql()
}

func (s *service) PersistOrdersFromRedisToMysql() error {
	orderIdCursor := 0
	pageSize := 2
	for {
		// get a page of keys (orderIds) from redis
		keys, err := s.redisRepo.GetOnePageOfKeys(orderIdCursor, pageSize)
		if err != nil {
			return fmt.Errorf("cannot get the list of orderIDs: %v", err)
		}
		orderIdCursor += pageSize

		// read the value of the keys from redis
		orders, getErr := s.redisRepo.GetOrders(keys)
		if getErr != nil {
			return fmt.Errorf("cannot get the orders of keys: %v", getErr)
		}

		// insert orders into mysql database
		if sqlInsertErr := s.mysqlRepo.InsertOrder(orders); sqlInsertErr != nil {
			return sqlInsertErr
		}
	}
}
