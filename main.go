package main

import (
	"context"
	"log"
	"net"
	"github.com/nasparria/GetData/protogen"
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type server struct {
	protogen.UnimplementedPortfolioServiceServer
	db *mongo.Database
}

func (s *server) GetOrdersByTicker(ctx context.Context, req *protogen.TickerRequest) (*protogen.OrdersResponse, error) {
	collection := s.db.Collection("orders")
	filter := bson.M{"ticker": req.GetTicker()}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*protogen.Order
	for cursor.Next(ctx) {
		var order protogen.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &protogen.OrdersResponse{Orders: orders}, nil
}

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	db := client.Database("portfolio")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protogen.RegisterPortfolioServiceServer(grpcServer, &server{db: db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
