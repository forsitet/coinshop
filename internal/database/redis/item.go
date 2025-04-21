package redis

import (
	"bufio"
	"coin/domain"
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	redis := &RedisClient{Client: rdb}
	if err = redis.insertItems(); err != nil {
		log.Fatal("failed to fill redis")
	}
	return redis
}

func (r *RedisClient) insertItems() error {
	const op = "database.redis.insertItems"
	log.Println(op, "attempting to read items.csv")
	file, err := os.Open("items.csv")
	if err != nil {
		log.Fatal(op, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		log.Println(op, "parsed line:", fields)

		err := r.Client.HSet(ctx, "item:"+fields[0], map[string]interface{}{
			"price": fields[1],
		}).Err()
		if err != nil {
			log.Println(op, err)
			continue
		}
		log.Println(op, "successfully set item:", fields[0])
	}
	return nil
}

func (r *RedisClient) GetItemPrice(itemName string) (int, error) {
	priceStr, err := r.Client.HGet(ctx, "item:"+itemName, "price").Result()
	if err != nil {
		return 0, err
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (r *RedisClient) GetItems() ([]domain.Item, error) {
	var cursor uint64
	var items []domain.Item

	for {
		keys, newCursor, err := r.Client.Scan(ctx, cursor, "item:*", 10).Result()
		if err != nil {
			return nil, err
		}
		cursor = newCursor

		for _, key := range keys {
			data, err := r.Client.HGetAll(ctx, key).Result()
			if err != nil {
				log.Println("HGetAll error:", err)
				continue
			}

			price, err := strconv.Atoi(data["price"])
			if err != nil {
				log.Println("Invalid price for key", key)
				continue
			}

			parts := strings.SplitN(key, "item:", 2)
			if len(parts) != 2 {
				log.Println("Invalid key format:", key)
				continue
			}

			item := domain.Item{
				Name:  parts[1],
				Price: price,
			}
			items = append(items, item)
		}

		if cursor == 0 {
			break
		}
	}

	return items, nil
}
