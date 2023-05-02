package httpd

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/rateit/internal/config"
	"github.com/buglloc/rateit/internal/provider"
)

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(cfg *config.Config) (*Server, error) {
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).JSON(ErrorRsp{
				Code:    code,
				Message: err.Error(),
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})
	app.Use(
		recover.New(),
		logger.New(),
	)

	providers := make(map[string]provider.Provider, len(cfg.Providers))
	for _, p := range cfg.Providers {
		pp, err := p.NewProvider()
		if err != nil {
			return nil, fmt.Errorf("unable to create provider: %w", err)
		}

		route := p.Route
		if route == "" {
			log.Warn().
				Str("provider", pp.Name()).
				Msg("no route configured, use default")

			route = pp.Name()
		}

		providers[route] = pp
	}

	app.Route("/api/v1/rate", func(r fiber.Router) {
		for path, p := range providers {
			r.Get(path, newRateHandler(p))
		}
	})

	return &Server{
		app:  app,
		addr: cfg.Addr,
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.app.Listen(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func newRateHandler(p provider.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rate, err := p.CurrentRate(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(RateRsp{
			Rate: rate,
		})
	}
}
