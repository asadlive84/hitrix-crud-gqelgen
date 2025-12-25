package main

import (
	"fmt"
	"hitrix-crud/entity"
	"hitrix-crud/graph"

	"git.ice.global/packages/hitrix"
	"git.ice.global/packages/hitrix/pkg/middleware"
	"git.ice.global/packages/hitrix/service/component/app"
	"git.ice.global/packages/hitrix/service/registry"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	s, deferFunc := hitrix.New(
		"app-name", "your-secret-key",
	).RegisterDIGlobalService(
		registry.ServiceProviderConfigDirectory("../config"),
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrmEngine(),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderGenerator(), 

	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(true),
	).RegisterRedisPools(&app.RedisPools{Persistent: "default"}).
		Build()

	defer deferFunc()

	// GraphQL schema
	executableSchema := graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	})

	fmt.Println("🚀 Server starting on :9999")
	fmt.Println("🎮 Playground: http://localhost:9999/playground")
	fmt.Println("🔗 GraphQL Endpoint: http://localhost:9999/query")

	s.RunServer(9999, executableSchema, ginSetup, gqlSetup)
}

func ginSetup(engine *gin.Engine) {
	middleware.Cors(engine)

	engine.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}

func gqlSetup(srv *handler.Server) {
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

}
