package store

import (
	"github.com/progrium/cmd/com/core"

	"github.com/gliderlabs/gosper/pkg/com"
)

func init() {
	com.Register("store", struct{}{},
		com.Option("backend", "store.filesystem", "Store backend"))
}

func Selected() Backend {
	backend := com.Select(com.GetString("backend"), new(Backend))
	if backend == nil {
		panic("Unable to find selected backend: " + com.GetString("backend"))
	}
	return backend.(Backend)
}

type Backend interface {
	CmdBackend
	TokenBackend
}

type CmdBackend interface {
	List(user string) []*core.Command
	Get(user, name string) *core.Command
	Put(user, name string, cmd *core.Command) error
	Delete(user, name string) error
}

type TokenBackend interface {
	ListTokens(user string) ([]*core.Token, error)
	GetToken(key string) (*core.Token, error)
	PutToken(token *core.Token) error
	DeleteToken(key string) error
}