/**
 *
 * @author liangjf
 * @create on 2020/9/9
 * @version 1.0
 */
package db

const (
	MsgSeqIDIncrKey = "MsgSeqID"
)

func (p *RedisPool) GenerateMsgSeq(uuid string) uint64 {
	seq, _ := p.HIncrBy(MsgSeqIDIncrKey, uuid, 1)
	return uint64(seq)
}
