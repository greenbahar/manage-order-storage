package redis_db

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/greenbahar/manage-order-storage/config/env"
	"github.com/greenbahar/manage-order-storage/order"
	"log"
	"sync"
	"time"
)

const (
	maxIdleConnections   = 10
	idleTimeout          = 20 * time.Second
	maxActiveConnections = 20
	name_of_orderIDs_set = "orderIDs"
)

var (
	orderScore int
	mu         sync.Mutex
)

type RedisRepo struct {
	Pool *redis.Pool
}

type OrderIdScore struct {
	Score int
	Id    int
}

func NewRepo() *RedisRepo {
	redisConf := env.GetRedisConfig()

	pool := &redis.Pool{
		MaxIdle:     maxIdleConnections,
		MaxActive:   maxActiveConnections,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisConf.RedisContainer, redisConf.Port))
			if err != nil {
				panic(err)
			}
			return conn, err
		},
	}

	repo := &RedisRepo{pool}
	ping(repo)

	return repo
}

func ping(r *RedisRepo) {
	con := r.Pool.Get()
	defer con.Close()

	_, err := redis.String(con.Do("PING"))
	if err != nil {
		panic(err)
	}
}

func (r *RedisRepo) SetOrder(order *order.Order) error {
	// get conn and put back when exit from method
	conn := r.Pool.Get()
	defer conn.Close()

	key := order.OrderId
	val, jErr := json.Marshal(order)
	if jErr != nil {
		log.Printf("ERROR: fail to marshal val %s, error %s", val, jErr.Error())
		return jErr
	}
	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	mu.Lock()
	ordScore := orderScore
	orderScore++
	mu.Unlock()
	_, err = conn.Do("ZADD", name_of_orderIDs_set, ordScore, key)
	if err != nil {
		log.Printf("ERROR: fail to add key %d to the set %s, error %s", key, name_of_orderIDs_set, err.Error())
		return err
	}

	return nil
}

func (r *RedisRepo) GetOrder(key string) (*order.Order, error) {
	// get conn and put back when exit from method
	conn := r.Pool.Get()
	defer conn.Close()

	val, err := conn.Do("GET", key)
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return nil, err
	}
	ord := &order.Order{}
	if unmarshalErr := json.Unmarshal(val.([]byte), ord); unmarshalErr != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return nil, unmarshalErr
	}

	return ord, nil
}

func (r *RedisRepo) GetOrders(keys []string) ([]*order.Order, error) {
	// get conn and put back when exit from method
	conn := r.Pool.Get()
	defer conn.Close()

	orders := make([]*order.Order, 0)
	for _, key := range keys {
		ord, err := r.GetOrder(key)
		if err != nil {
			log.Printf("ERROR: fail get value of key %s, error %s", key, err.Error())
			return nil, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}

func (r *RedisRepo) GetOnePageOfKeys(startIndex, rangeIndex int) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	result, err := conn.Do("ZRANGE", name_of_orderIDs_set, startIndex, rangeIndex)
	if err != nil {
		log.Printf("ERROR: fail get ZRANGE from set %s, error %s", name_of_orderIDs_set, err.Error())
		return nil, err
	}

	fmt.Println(result.([]int))
	return result.([]string), nil
}
