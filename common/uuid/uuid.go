/**
 *
 * @author liangjf
 * @create on 2020/5/25
 * @version 1.0
 */
package uuid

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func NewUuid() string {
	rand.Seed(time.Now().Unix())
	u, _ := uuid.NewUUID()
	return u.String()
}
