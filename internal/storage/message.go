package storage

import (
	"messenger/proto"
)

type Store interface {
	Save(*messenger.ChatMessage)
	List() []*messenger.ChatMessage
	History() *messenger.ChatHistory
}
