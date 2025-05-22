package router

import (
	"time"

	"github.com/cobanhub/lib/response"
	json "github.com/goccy/go-json"
	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type (
	Ctx fiber.Ctx

	Opts struct {
		AppName      string
		Prefix       string
		ReadTimeout  int                  //Read Timeout in seconds
		WriteTimeout int                  //Write Timeout in seconds
		RateLimit    map[string]PathLimit //Rate limit per path
		BodyMaxSize  int                  //Body Max Size in mb
		WithMetrics  bool                 //Enable metrics
	}

	Router struct {
		httprouter *fiber.App
		Opts       Opts
	}

	PathLimit struct {
		Rate map[string]Rate //configuration for rate limit per ip
	}

	Rate struct {
		Limit      int
		Expiration time.Duration
	}

	Handler func(ctx *Ctx) *response.JSONResponse
)

/*
*

	@opts: Opts - Options for configuring the router.
		- AppName: Name of the application.
		- Prefix: Prefix for all routes.
		- ReadTimeout: Read timeout in seconds.
		- WriteTimeout: Write timeout in seconds.
		- RateLimit: Rate limit configuration per path.
		- BodyMaxSize: Maximum body size in MB.
		- WithMetrics: Enable metrics endpoint.
	@returns: *Router - Pointer to the initialized Router struct.
	@description:
		New creates a new fiber router with the given options.
		It initializes the router with the specified configurations, including
		read/write timeouts, body size limits, and optional metrics.
		The router is set to strict routing and case sensitivity.
		It also sets up the JSON encoder and decoder for request/response handling.
		The function returns a pointer to the initialized Router struct.
		The router is ready to be used for adding routes and handling requests.
		The function also sets up a metrics endpoint if the WithMetrics option is enabled.
		The metrics endpoint is accessible at "/metrics" and provides monitoring information.
		The function also sets up a rate limiter for specific paths if the RateLimit option is provided.
		The rate limiter is configured to limit the number of requests per IP address.
		The rate limit configuration is specified in the Opts struct.
		The function returns a pointer to the initialized Router struct.
		The router is ready to be used for adding routes and handling requests.
*/
func New(opts Opts) *Router {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		BodyLimit:     opts.BodyMaxSize * 1024 * 1024,
		ReadTimeout:   time.Duration(opts.ReadTimeout) * time.Second,
		WriteTimeout:  time.Duration(opts.WriteTimeout) * time.Second,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		AppName:       opts.AppName,
	})

	if opts.WithMetrics {
		app.Get("/metrics", monitor.New())
	}

	return &Router{
		httprouter: app,
		Opts:       opts,
	}
}

func (r *Router) GET(path string, handler Handler) {
	r.AddRoute("GET", path, handler)
}

func (r *Router) AddRoute(method, path string, handler Handler) {
	middlewares := []fiber.Handler{}

	if r.Opts.RateLimit != nil {
		middlewares = append(middlewares, r.setLimiter(path))
	}

	switch method {
	case "GET":
		r.httprouter.Get(path, r.wrapHandler(handler), middlewares...)
	case "POST":
		r.httprouter.Post(path, r.wrapHandler(handler), middlewares...)
	case "PUT":
		r.httprouter.Put(path, r.wrapHandler(handler), middlewares...)
	case "DELETE":
		r.httprouter.Delete(path, r.wrapHandler(handler), middlewares...)
	case "PATCH":
		r.httprouter.Patch(path, r.wrapHandler(handler), middlewares...)
	case "OPTIONS":
		r.httprouter.Options(path, r.wrapHandler(handler), middlewares...)
	case "HEAD":
		r.httprouter.Head(path, r.wrapHandler(handler), middlewares...)
	}
}

func (router *Router) Group(path string, fn func(r *Router)) {
	group := router.httprouter.Group(path)
	groupRouter := &Router{
		httprouter: group.(*fiber.App), // or wrap in a custom struct
		Opts:       router.Opts,        // Consider cloning and updating prefix
	}
	fn(groupRouter)
}

func (router *Router) wrapHandler(h Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := Ctx(c)
		resp := h(&ctx)
		if resp != nil {
			return c.JSON(resp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&response.JSONResponse{
			Code:    500,
			Message: "Unexpected nil response from handler",
		})
	}
}
func (r *Router) setLimiter(path string) fiber.Handler {
	cfg, ok := r.Opts.RateLimit[path]
	if !ok {
		return func(c fiber.Ctx) error { return c.Next() } // no limiter for this path
	}

	return limiter.New(limiter.Config{
		MaxFunc: func(c fiber.Ctx) int {
			ip := c.IP()
			if rate, ok := cfg.Rate[ip]; ok {
				return rate.Limit
			}
			return 100 // default
		},
		Expiration: 30 * time.Second, // static expiration; you can enhance it
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&response.JSONResponse{
				Code:    fiber.StatusTooManyRequests,
				Message: "Rate limit exceeded",
			})
		},
	})
}
