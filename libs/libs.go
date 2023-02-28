package libs

import (
	"context"
	"time"

	"example.com/m/v2/keys"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB      = &mongo.Client{}
	RENNYMongoDB = &mongo.Client{}
)

type (
	MongoConfig struct {
		Host       string
		DBName     string
		Collection string
	}
)

func StartMongoServiceRENNY() {

	_ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	RENNYMongoDB, _ = mongo.Connect(_ctx, options.Client().ApplyURI("mongodb+srv://Renny:Rennyyy2M@cluster0.skgku.mongodb.net/"+keys.Database))
}

func StartMongoService(host, dbname, collection string, isDev bool) {
	//fmt.Println("Running mongo service")
	mc := new(MongoConfig)

	mc.Host = host
	mc.DBName = dbname
	mc.Collection = collection

	_ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()
	MongoDB, _ = mongo.Connect(_ctx, options.Client().ApplyURI(mc.Host+"/"+mc.DBName))

}
