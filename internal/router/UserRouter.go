package router

import (
	"github.com/K1ender/MemeWhisper/internal/dto"
	"github.com/K1ender/MemeWhisper/internal/response"
	"github.com/K1ender/MemeWhisper/internal/store"
	"github.com/K1ender/MemeWhisper/internal/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type userRoutes struct {
	logger    *zap.Logger
	userStore store.IUserStore
}

func NewUserRouter(
	router fiber.Router,
	logger *zap.Logger,
	userStore store.IUserStore,
) {
	routes := userRoutes{
		logger:    logger,
		userStore: userStore,
	}

	router.Post("/register", routes.registerRoute)
}

func (r *userRoutes) registerRoute(ctx *fiber.Ctx) error {
	body := dto.RegisterUserDTO{}
	ctx.BodyParser(&body)

	errs, ok := validator.Validate(body)

	if !ok {
		return ctx.Status(422).JSON(response.NewFailResponse(422, errs))
	}

	return ctx.JSON(body)
}
