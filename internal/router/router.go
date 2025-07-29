package router

import (
	"net/http"
	"strings"

	"github.com/nshmdayo/nft-platform-sample/internal/config"
	"github.com/nshmdayo/nft-platform-sample/internal/handlers"
	"github.com/nshmdayo/nft-platform-sample/internal/middleware"
)

type Router struct {
	cfg          *config.Config
	routeHandler *handlers.RouteHandler
}

func NewRouter(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	paperHandler *handlers.PaperHandler,
	reviewHandler *handlers.ReviewHandler,
) *Router {
	return &Router{
		cfg:          cfg,
		routeHandler: handlers.NewRouteHandler(authHandler, paperHandler, reviewHandler),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", r.routeHandler.HealthHandler.Health)

	// Auth routes (public)
	mux.HandleFunc("/api/v1/auth/register", r.routeHandler.AuthHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", r.routeHandler.AuthHandler.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(r.cfg)

	// Auth protected routes
	mux.Handle("/api/v1/auth/profile", authMiddleware(http.HandlerFunc(r.routeHandler.AuthHandler.GetProfile)))

	// Paper routes
	mux.Handle("/api/v1/papers/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "/reviews") {
			r.routeHandler.HandlePaperReviews(w, req)
		} else {
			r.routeHandler.HandlePapers(w, req)
		}
	})))
	mux.Handle("/api/v1/papers/my", authMiddleware(http.HandlerFunc(r.routeHandler.PaperHandler.GetMyPapers)))

	// Review routes
	mux.Handle("/api/v1/reviews/", authMiddleware(http.HandlerFunc(r.routeHandler.HandleReviews)))
	mux.Handle("/api/v1/reviews/my", authMiddleware(http.HandlerFunc(r.routeHandler.ReviewHandler.GetMyReviews)))
	mux.Handle("/api/v1/reviews/pending", authMiddleware(http.HandlerFunc(r.routeHandler.ReviewHandler.GetPendingReviews)))

	// Apply global middleware
	handler := middleware.LoggerMiddleware()(mux)
	handler = middleware.CORSMiddleware()(handler)

	return handler
}
