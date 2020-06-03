/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import "github.com/chasex/redis-go-cluster"

func (p *RedisPool) Set(key string, value interface{}) error {
	if _, err := p.cluster.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) SetEX(key string, value interface{}, expireSecond int) error {
	if _, err := p.cluster.Do("SET", key, value, "EX", expireSecond); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) SetPX(key string, value interface{}, expireMS int) error {
	if _, err := p.cluster.Do("SET", key, value, "PX", expireMS); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) SetNX(key string, value interface{}) error {
	if _, err := p.cluster.Do("SET", key, value, "NX"); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) Get(key string) (interface{}, error) {
	return p.cluster.Do("GET", key)
}

func (p *RedisPool) GetString(key string) (string, error) {
	return redis.String(p.cluster.Do("GET", key))
}

func (p *RedisPool) GetInt(key string) (int, error) {
	return redis.Int(p.cluster.Do("GET", key))
}

func (p *RedisPool) Incr(key string) (int64, error) {
	reply, err := redis.Int64(p.cluster.Do("INCR", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func (p *RedisPool) IncrBy(key string, increment int64) (int64, error) {
	reply, err := redis.Int64(p.cluster.Do("INCRBY", key, increment))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func (p *RedisPool) Decr(key string) (int64, error) {
	reply, err := redis.Int64(p.cluster.Do("DECR", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func (p *RedisPool) DecrBy(key string, decrement int64) (int64, error) {
	reply, err := redis.Int64(p.cluster.Do("DECRBY", key, decrement))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func (p *RedisPool) Strlen(key string) (int, error) {
	reply, err := redis.Int(p.cluster.Do("STRLEN", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}
