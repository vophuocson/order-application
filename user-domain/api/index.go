package api

import "user-domain/api/inbound"

type Handler interface {
	inbound.UserApi
}
