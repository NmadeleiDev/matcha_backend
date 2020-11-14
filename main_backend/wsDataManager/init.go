package wsDataManager

import (
	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy"
)

type wsData struct {
	conn *kafka.Conn
}


var manager wsData

