package redis-utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

func NewRD(address string, pass string, DB int) *redis.Client { // коннект к редису
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass, // no password set if need
		DB:       DB,    // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}

func DeleteAllKeys(rd *redis.Client, key ...string) error { // удаление всех связанных записей
	for i := range key {
		iter := rd.Scan(0, key[i]+"*", 0).Iterator()
		for iter.Next() {
			err := rd.Del(iter.Val()).Err()
			if err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}
	return nil
}

func DeleteOneKey(rd *redis.Client, key string) error { // удаление одного ключа
	err := rd.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}

func SetKey(rd *redis.Client, key string, value interface{}, duration int) error { // создание записи по ключу
	js, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rd.Set(key, js, time.Duration(duration)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}


func GetVal(rd *redis.Client, key string, value interface{}) (interface{}, error) { // получение данных по ключу
	js, err := rd.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(js, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
