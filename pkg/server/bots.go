package server

import "github.com/gofiber/fiber/v2"

func (s *Server) bots(ctx *fiber.Ctx) error {
	bots, err := s.botRepository.Bots(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	return ctx.JSON(bots)
}
