/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import "github.com/chasex/redis-go-cluster"

func (p *RedisPool) HGet(key, field string) (interface{}, error) {
	return p.cluster.Do("HGET", key, field)
}

func (p *RedisPool) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(p.cluster.Do("HGETALL", key))
}

func (p *RedisPool) HKeys(key string) ([]string, error) {
	return redis.Strings(p.cluster.Do("HKEYS", key))
}

func (p *RedisPool) HLen(key string) (int, error) {
	return redis.Int(p.cluster.Do("HLEN", key))
}

func (p *RedisPool) HMGet(key string, fields ...string) ([]string, error) {
	reply, err := redis.Strings(p.cluster.Do("HMGET", key))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (p *RedisPool) HSet(key, field string, value interface{}) error {
	if _, err := p.cluster.Do("HSET", key, field, value); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) HMSet(key string, pairs ...string) error {
	if _, err := p.cluster.Do("HSET", key, pairs); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) HSetNX(key, field string, value interface{}) (bool, error) {
	return redis.Bool(p.cluster.Do("HSETNX", key, field, value))
}

func (p *RedisPool) HIncrBy(key, field string, increment int) (int64, error) {
	return redis.Int64(p.cluster.Do("HINCRBY", key, field, increment))
}

func (p *RedisPool) HDel(key string, fields ...string) (int, error) {
	return redis.Int(p.cluster.Do("HDEL", key, fields))
}

func (p *RedisPool) HExist(key, field string) (bool, error) {
	return redis.Bool(p.cluster.Do("HEXISTS", key, field))
}

func (p *RedisPool) HStrlen(key, field string) (int, error) {
	return redis.Int(p.cluster.Do("HSTRLEN", key, field))
}
