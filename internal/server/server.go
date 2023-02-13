package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/artemmarkaryan/exlex-backend/graph"
	"github.com/artemmarkaryan/exlex-backend/internal/service"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func Serve(ctx context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	serviceContainer, err := service.MakeContainer(ctx)
	if err != nil {
		return fmt.Errorf("initing service container: %w", err)
	}

	router := chi.NewRouter()
	router.Use(
		ContextPropagateMiddleware(ctx),
		// todo: auth
		// todo: log errors
	)

	graphqlSchema := graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{
			ServiceContainer: serviceContainer,
		}},
	)

	playgroundPath := "playground"
	router.Handle("/"+playgroundPath, playground.Handler("playground", "/query"))
	router.Handle("/query", handler.NewDefaultServer(graphqlSchema))

	log.Printf("connect to http://localhost:%s/%s for GraphQL playground", port, playgroundPath)
	return http.ListenAndServe(":"+port, router)
}

func ContextPropagateMiddleware(ctx context.Context) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx = database.Propagate(r.Context(), database.C(ctx))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
