package env

import (
	"os"

	"github.com/exim-id/go-exim-base-lib/errors"
	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func GetGrpcPort() string {
	defer errors.NormalError()
	return os.Getenv("GRPC_PORT")
}

func GetMongoDbName() string {
	defer errors.NormalError()
	return os.Getenv("MONGO_DBNAME")
}

func GetMongoDBUrl() string {
	defer errors.NormalError()
	return os.Getenv("MONGO_URL")
}

func GetServerPort() string {
	defer errors.NormalError()
	return os.Getenv("SERVER_PORT")
}
