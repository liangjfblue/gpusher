/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import "github.com/chasex/redis-go-cluster"

func (p *RedisPool) Del(key string) error {
	if _, err := p.cluster.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) MDel(keys ...string) (int, error) {
	return redis.Int(p.cluster.Do("DEL", keys))
}

func (p *RedisPool) Exists(key string) (bool, error) {
	return redis.Bool(p.cluster.Do("EXISTS", key))
}

func (p *RedisPool) MExists(keys string) (int, error) {
	return redis.Int(p.cluster.Do("EXISTS", keys))
}

func (p *RedisPool) Expire(key string, seconds int) error {
	if _, err := p.cluster.Do("EXPIRE", key, seconds); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) RenameNX(key, nKey string) error {
	if _, err := p.cluster.Do("RENAMENX", key, nKey); err != nil {
		return err
	}
	return nil
}

func (p *RedisPool) TTL(key string) (int, error) {
	return redis.Int(p.cluster.Do("TTL", key))
}

func (p *RedisPool) Type(key string) (string, error) {
	return redis.String(p.cluster.Do("TYPE", key))
}
