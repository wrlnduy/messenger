package storage

import (
	"messenger/proto"
)

type Store interface {
	Save(*message.ChatMessage)
	List() []*message.ChatMessage
	History() *message.ChatHistory
}
