package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	controllers "ProjectHash/controllers"
	"ProjectHash/pb"
	"ProjectHash/services"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svc := services.NewPasswordService()

	httpServer := NewHTTPServer(svc) // Use the NewHTTPServer function directly
	grpcServer := NewGRPCServer(svc) // Use the NewGRPCServer function directly

	go startHTTPServer(httpServer)
	go startGRPCServer(grpcServer, svc)

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

func startGRPCServer(server *grpc.Server, svc services.PasswordService) {
	pb.RegisterPasswordServiceServer(server, services.NewGRPCService(svc))
	reflection.Register(server)

	grpcAddr := ":9090"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	fmt.Println("gRPC server is running on http://localhost" + grpcAddr)
	if err := server.Serve(lis); err != nil {
		fmt.Println("Failed to serve:", err)
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

func NewGRPCServer(svc services.PasswordService) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterPasswordServiceServer(grpcServer, services.NewGRPCService(svc))
	reflection.Register(grpcServer)
	return grpcServer
}
