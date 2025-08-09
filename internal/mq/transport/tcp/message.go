package tcp

import (
	"time"

	"github.com/hoppermq/hopper/pkg/domain"
)

func (t *TCP) sendMessage(d []byte, conn domain.Connection) {
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(d)
	if err != nil {
		t.logger.Warn("error sending message", "error", err)
		return
	}
	t.logger.Info("message has been delivered")
}
