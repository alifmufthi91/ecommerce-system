package internal

import (
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user"
)

type Modules struct {
	User *user.UserModule
}

type InitOptions struct {
	_options.DefaultOptions
}

func InitModules(opts InitOptions) *Modules {

	userModule := user.NewUserModule(user.Options{
		DefaultOptions: opts.DefaultOptions,
	})

	return &Modules{
		User: userModule,
	}
}
