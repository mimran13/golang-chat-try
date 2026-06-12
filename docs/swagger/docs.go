// Package swagger contains auto-generated Swagger documentation.
//
// This is a placeholder file. Run the following command to generate
// the actual Swagger documentation from your code annotations:
//
//   make swagger
//
// Or manually:
//
//   swag init -g cmd/api/main.go -o docs/swagger
//
// After generation, this file will be replaced with the actual
// Swagger spec containing your API documentation.
package swagger

// SwaggerInfo holds the API metadata for Swagger documentation.
// This will be overwritten by swag init.
type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo contains the Swagger metadata.
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Title:       "Task Manager API",
	Description: "A production-grade REST API for task management.",
}
