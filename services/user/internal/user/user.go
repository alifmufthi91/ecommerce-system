package user

import (
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/handler"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/repository"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/service"
)

type UserModule struct {
	userService service.UserService
}

type Options struct {
	_options.DefaultOptions
}

func NewUserModule(opts Options) *UserModule {

	userRepo := repository.NewUserRepository(opts.Db)

	userService := service.NewUserService(opts.Config, userRepo)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, userService))

	return &UserModule{
		userService: userService,
	}
}
