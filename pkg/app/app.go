package app

import (
	"log"
	"net/http"

	"github.com/dibyendu/Authentication-Authorization/pkg/domain"
	"github.com/dibyendu/Authentication-Authorization/pkg/handler"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"github.com/dibyendu/Authentication-Authorization/pkg/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartApp() {
	log.Println("Starting app")
	router := mux.NewRouter()
	// Create a private router for authenticated routes
	privateRouter := router.PathPrefix("/").Subrouter()

	// Create instances of your handlers and services
	userRepo := domain.NewUserRepositoryDb()
	userService := service.NewUserService(userRepo)
	userHandler := handler.UserHandler{Service: userService}

	// Define routes and corresponding handler methods
	router.HandleFunc("/user/sign-in", userHandler.SignIn).Methods(http.MethodPost)
	privateRouter.HandleFunc("/user/home", userHandler.GetBookList).Methods(http.MethodGet)
	privateRouter.HandleFunc("/user/add-book", userHandler.AddBook).Methods(http.MethodPost)
	privateRouter.HandleFunc("/user/delete-book", userHandler.DeleteBook).Methods(http.MethodDelete)

	// Setup middleware
	privateRouter.Use(middleware.Authentication)

	// CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	// Start the HTTP server
	address := ":8080" // Replace with your desired port
	log.Printf("Server is listening on %s...\n", address)

	// Wrap the router with CORS middleware and start the server
	err := http.ListenAndServe(address, handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Fatal("HTTP server error:", err)
	}
}
