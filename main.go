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
	svc := services.NewPasswordService()

	httpServer := NewHTTPServer(svc)
	grpcServer := services.NewGRPCServer(svc)

	go startHTTPServer(httpServer)
	go startGRPCServer(grpcServer)

	waitForShutdown()
}

func startHTTPServer(server *http.Server) {
	httpAddr := ":8080"
	fmt.Println("HTTP server is running on http://localhost" + httpAddr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func startGRPCServer(server *grpc.Server) {
	grpcAddr := ":9090"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		fmt.Println("Error starting gRPC server:", err)
		return
	}
	fmt.Println("gRPC server is running on http://localhost" + grpcAddr)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println("Error serving gRPC:", err)
	}
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Shutting down gracefully...")
}

func NewHTTPServer(svc services.PasswordService) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/hash", controllers.HashPasswordHandler(svc)).Methods("POST")
	router.HandleFunc("/validate", controllers.ValidatePasswordHandler(svc)).Methods("POST")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	return httpServer
}
