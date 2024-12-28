package router

import (
	"github.com/K1ender/MemeWhisper/internal/dto"
	"github.com/K1ender/MemeWhisper/internal/models"
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

	logger.Debug("Registering user routes...")
	router.Post("/register", routes.registerRoute)
}

func (r *userRoutes) registerRoute(ctx *fiber.Ctx) error {
	body := dto.RegisterUserDTO{}

	err := ctx.BodyParser(&body)

	if err != nil {
		r.logger.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(500).JSON(response.NewErrorResponse(500, "Failed to parse request body", response.ErrorResponseOpts{}))
	}

	errs, ok := validator.Validate(body)

	if !ok {
		return ctx.Status(422).JSON(response.NewFailResponse(422, errs))
	}

	user := models.User{
		Username:       *body.Username,
		HashedPassword: *body.Password,
	}

	userID, err := r.userStore.CreateUser(user)

	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err))
		if err == store.ErrUserAlreadyExists {
			return ctx.Status(409).JSON(response.NewFailResponse(409, "User already exists"))
		}
		return ctx.Status(500).JSON(response.NewErrorResponse(500, "Failed to create user", response.ErrorResponseOpts{}))
	}

	r.logger.Debug("User created", zap.Any("userID", userID))

	return ctx.JSON(response.NewSuccessResponse(200, userID))
}
