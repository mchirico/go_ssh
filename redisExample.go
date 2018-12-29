package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

var client   *redis.Client

func Client() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDISPWD"),
		DB:       0,                     // use default DB
	})

	return client
}


func ExampleNewClient() {


	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	err = client.Set("key", "Yes worked", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func ExamplePipe() {

	redisdb := client
	pipe := redisdb.Pipeline()

	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)

	// Execute
	//
	//     INCR pipeline_counter
	//     EXPIRE pipeline_counts 3600
	//
	// using one redisdb-server roundtrip.
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)



	redisdb.WrapProcessPipeline(func(old func([]redis.Cmder) error) func([]redis.Cmder) error {
		return func(cmds []redis.Cmder) error {
			fmt.Printf("pipeline starting processing: %v\n", cmds)
			err := old(cmds)
			fmt.Printf("pipeline finished processing: %v\n", cmds)
			return err
		}
	})

	redisdb.Pipelined(func(pipe redis.Pipeliner) error {

		pipe.Set("one","1", 5*time.Minute)
		pipe.Incr("X")
		pipe.Get("one")
		pipe.Get("key")

		return nil
	})
}


func main() {
	client = Client()
	ExampleNewClient()
	ExamplePipe()
}
