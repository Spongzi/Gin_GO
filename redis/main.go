package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func initClient() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100, // 连接池大小
	})
	_, err := rdb.Ping().Result()
	return err
}

// redisExample
func redisExample() {
	err := rdb.Set("score", 100, 0).Err()
	if err != nil {
		fmt.Println("redis set failed", err)
		return
	}
	val, err := rdb.Get("score").Result()
	if err != nil {
		fmt.Println("redis get failed", err)
		return
	}
	fmt.Printf("score value = %v\n", val)
	val2, err := rdb.Get("name").Result()
	if err == redis.Nil {
		fmt.Println("name is not exist")
	} else if err != nil {
		fmt.Println("redis get failed", err)
	} else {
		fmt.Println("name", val2)
	}
}

func main() {
	err := initClient()
	if err != nil {
		fmt.Println("init redis client failed", err)
		return
	}
	// 程序退出时释放相关资源
	defer func() {
		err = rdb.Close()
		if err != nil {
			fmt.Println("redis close failed!", err)
			return
		}
	}()
	redisExample()
	fmt.Println("connect redis success...")
}
