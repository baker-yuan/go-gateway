package registry

import "github.com/google/wire"

var RegistrarSet = wire.NewSet(NewRegistrar)
