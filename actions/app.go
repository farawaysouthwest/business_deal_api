package actions

import (
	"os"
	"sync"

	"business_deal_api/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_business_deal_api_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)


		// // Wraps each request in a transaction.
		//c.Value("tx").(*pop.Connection)
		// // Remove to disable this.
		//app.Use(popmw.Transaction(models.DB))
		// // Setup and use translations:
		// app.Use(translations())




		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "https://studio.apollographql.com"},
			AllowCredentials: true,
		})

		graphqlHandler := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
		graphqlHandler.AddTransport(transport.POST{})
		graphqlHandler.AddTransport(transport.Options{})
		if os.Getenv("ENVIRONMENT") == "development" {
    	graphqlHandler.Use(extension.Introspection{})
		}

		app.ANY("/graphql", buffalo.WrapHandlerFunc(c.Handler(graphqlHandler).ServeHTTP))

		//app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory

	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// // for more information: https://gobuffalo.io/en/docs/localization
// func translations() buffalo.MiddlewareFunc {
// 	var err error
// 	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
// 		app.Stop(err)
// 	}
// 	return T.Middleware()
// }

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
