package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/godfreyowidi/simple-ecomm-demo/db"
	"github.com/godfreyowidi/simple-ecomm-demo/gql-gateway/graph"
	"github.com/godfreyowidi/simple-ecomm-demo/gql-gateway/resolvers"
	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Init DB
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	productRepo := repo.NewProductRepo(database.Pool)
	customerRepo := repo.NewCustomerRepo(database.Pool)
	orderRepo := repo.NewOrderRepo(database.Pool)
	orderItemRepo := repo.NewOrderItemRepo(database.Pool)
	createCategoryRepo := repo.NewCategoryRepo(database.Pool)

	// Construct the resolver with all dependencies
	resolver := &resolvers.Resolver{
		ProductRepo:   productRepo,
		CustomerRepo:  customerRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		CategoryRepo:  createCategoryRepo,
	}

	// GraphQL server setup
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
