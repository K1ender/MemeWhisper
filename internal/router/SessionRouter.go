package router

import (
	"github.com/K1ender/MemeWhisper/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewSessionRouter(
	router fiber.Router,
	logger *zap.Logger,
	userStore store.IUserStore,
	sessionStore store.ISessionStore,
) {

}
