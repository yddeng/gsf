package agnet

import "github.com/yddeng/dnet"

type (
	Codec interface {
		dnet.Codec
	}

	ServerEvent interface {
		OnConnection()
		OnClose()
		OnMessage()
	}

	Agent struct {
		acceptor dnet.Acceptor
	}
)

func New(acceptor dnet.Acceptor)
