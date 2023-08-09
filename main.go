package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	controllers "ProjectHash/controllers"
	"ProjectHash/services"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	// Create a new instance of the PasswordService
	svc := services.NewPasswordService()

	// Create HTTP and gRPC servers with the service instance
	httpServer := NewHTTPServer(svc)
	grpcServer := services.NewGRPCServer(svc)

	// Start the HTTP and gRPC servers in goroutines
	go startHTTPServer(httpServer)
	go startGRPCServer(grpcServer)

	// Wait for a termination signal
	waitForShutdown()
}

func startHTTPServer(server *http.Server) {
	// Define the HTTP server address
	httpAddr := ":8080"
	fmt.Println("HTTP server is running on http://localhost" + httpAddr)

	// Start the HTTP server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func startGRPCServer(server *grpc.Server) {
	// Define the gRPC server address
	grpcAddr := ":9090"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		fmt.Println("Error starting gRPC server:", err)
		return
	}
	fmt.Println("gRPC server is running on http://localhost" + grpcAddr)

	// Start the gRPC server
	err = server.Serve(lis)
	if err != nil {
		fmt.Println("Error serving gRPC:", err)
	}
}

func waitForShutdown() {
	// Set up a channel to receive termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal and print a shutdown message
	<-sigCh
	fmt.Println("Shutting down gracefully...")
}

func NewHTTPServer(svc services.PasswordService) *http.Server {
	// Create a new HTTP router using the gorilla/mux library
	router := mux.NewRouter()
	router.HandleFunc("/hash", controllers.HashPasswordHandler(svc)).Methods("POST")
	router.HandleFunc("/validate", controllers.ValidatePasswordHandler(svc)).Methods("POST")

	// Create an HTTP server instance with the router
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	return httpServer
}
