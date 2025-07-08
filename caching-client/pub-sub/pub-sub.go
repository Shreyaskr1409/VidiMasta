package pubsub

import (
	"log"

	"github.com/valkey-io/valkey-go"
)

var Client valkey.Client

func Init(l *log.Logger) {
	c, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		l.Fatalln("Error encountered while initializing client at :6379", err)
	}

	Client = c
}
