# ProjectHash Application

This is a sample project showcasing how to create a password hashing and validation service in Go.

### Prerequisites

- Go (1.16+)
- Protocol Buffers (protoc)

### Generating Go Code from Proto

1. Install Protocol Buffers compiler (protoc). Follow the instructions in the [Protocol Buffers documentation](https://developers.google.com/protocol-buffers/docs/downloads).

2. Create a directory named `pb` in your project's root directory.

3. Write your `password_service.proto` file in the `pb` directory.

4. Run the following command to generate Go code from the proto file:

   ```bash
   protoc -I pb/ pb/password_service.proto --go_out=plugins=grpc:pb

### Running the Application

1. Clone this repository to your local machine.

    ```bash
    git clone https://github.com/emankashifplex/ProjectHash.git


2. Install the required Go dependencies.
  
    ```bash
    go mod download

3. Build and run the application.

    ```bash
    go run main.go

    
The HTTP server will be running on http://localhost:8080, and the gRPC server will be running on http://localhost:9090.

### Usage

#### Hashing a Password

To hash a password, make a POST request to the /hash endpoint using curl. Replace yourpassword with the password you want to hash:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"password": "yourpassword"}' http://localhost:8080/hash
```

#### Validating a Password

To validate a password against a hashed password, make a POST request to the /validate endpoint using curl. Replace yourhashedpassword with the hashed password and yourpassword with the password you want to validate:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"hashed_password": "yourhashedpassword", "password": "yourpassword"}' http://localhost:8080/validate
 ```

### gRPC Service
#### Hashing a Password
To hash a password using gRPC, you can use tools like grpcurl or create a client in your preferred programming language.

```bash
grpcurl -plaintext -d '{"password": "yourpassword"}' localhost:9090 passwordpb.PasswordService.HashPassword
```

#### Validating a Password

```bash
grpcurl -plaintext -d '{"hashed_password": "yourhashedpassword", "password": "yourpassword"}' localhost:9090 passwordpb.PasswordService.ValidatePassword
```

