package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var E *echo.Echo

func SetupWebServer() {
	E = echo.New()
	log.Print("Hello World!")
	log.Print("Web here, we are running version 0.5 (hopefully xD)")
	E.Use(middleware.RequestLogger())
	E.Use(middleware.Recover())
	E.Use(middleware.RequestID())
	// e.Use(session.Middleware(sessions.NewCookieStore([]byte("supersecret123"))))
	wd, err := os.Getwd()
	if err != nil {
		E.Logger.Fatalf("Could not get the current path... %v", err)
	}
	sessionPath := filepath.Join(wd, "sessions")
	if err := os.MkdirAll(sessionPath, 0700); err != nil {
		E.Logger.Fatalf("Failed to make the directories... %v", err)
	}

	store := sessions.NewFilesystemStore(sessionPath, []byte(os.Getenv("WORLDGAME_AUTH_SECRET")))
	store.Options.HttpOnly = true
	E.Use(session.Middleware(store))

	E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "localhost:5173", "https://fmr.us.to:9090", "http://fmr.us.to:9090", "fmr.us.to:9090", "fmr.us.to"},
		AllowCredentials: true,
	}))

	E.Use(middleware.RemoveTrailingSlash())

	// Serve ONLY public assets explicitly, so HTML files remain protected
	E.Static("/assets", "web/assets")
	// Optional: If Vite outputs other public files you need exposed (like favicon), map them:
	// E.File("/vite.svg", "web/vite.svg")

	// Serve specific production-related HTML pages cleanly
	E.GET("/login", GetLogin)
	E.GET("/register", GetRegister)

	// Group all API requests under /api
	api := E.Group("/api")

	// Public API endpoints
	api.POST("/login", Login)
	api.POST("/register", Register)
	api.GET("/countries", GetCountries)

	// Create a single protected group for BOTH pages and API routes
	protected := E.Group("")

	// very important auth middleware! DO NOT REMOVE!!!
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				// Return JSON error for API requests, but redirect for page requests
				if strings.HasPrefix(c.Request().URL.Path, "/api") {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized. Please log in!"})
				}
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			c.Set("username", sess.Values["username"])

			return next(c)
		}
	})

	// Protected Page endpoints
	protected.GET("/game", GetGame)
	// redirect root "/" to "/game" natively
	protected.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusSeeOther, "/game") })

	// Protected API endpoints
	// protectedApi := protected.Group("/api")
	// protectedApi.GET()

	isProd, err := strconv.ParseBool(os.Getenv("PROD"))
	if err != nil {
		log.Fatal("Could not parse PROD Environment Variable!")
	}

	if !isProd {
		E.Logger.Fatal(E.Start(":8080"))
	} else {
		E.Use(middleware.HTTPSRedirect())
		log.Fatal(E.StartTLS(":9090", "/cert/fullchain.pem", "/cert/privkey.pem"))
	}
}

func GetLogin(c echo.Context) error {
	return c.File("web/login.html")
}

func GetRegister(c echo.Context) error {
	return c.File("web/register.html")
}

func GetGame(c echo.Context) error {
	return c.File("web/game.html")
}
