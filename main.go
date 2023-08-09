package main

import (
	"context"
	"log"
	"net"
	"time"
	"github.com/nasparria/GetData/protogen"
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type server struct {
	protogen.UnimplementedPortfolioServiceServer
	db *mongo.Database
}

type MongoOrder struct {
	Account      string             `bson:"account"`
	Action       string             `bson:"action"`
	AveragePrice string             `bson:"average_price"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
	Fee          string             `bson:"fee"`
	IsPrime      bool               `bson:"is_prime"`
	LimitPrice   string             `bson:"limit_price"`
	MarketTime   string             `bson:"market_time"`
	Notional     string             `bson:"notional"`
	OrderID      string             `bson:"order_id"`
	OrderType    string             `bson:"order_type"`
	Quantity     string             `bson:"quantity"`
	QuantityShares string           `bson:"quantity_shares"`
	QuantityType string             `bson:"quantity_type"`
	Status      string              `bson:"status"`
	Ticker      string              `bson:"ticker"`
	TradingType string              `bson:"trading_type"`
	UpdatedAt   primitive.DateTime  `bson:"updated_at"`
	UserID      string              `bson:"user_id"`
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
		var mongoOrder MongoOrder
		if err := cursor.Decode(&mongoOrder); err != nil {
			return nil, err
		}

		order := &protogen.Order{
			Account:       mongoOrder.Account,
			Action:        mongoOrder.Action,
			AveragePrice:  mongoOrder.AveragePrice,
			Fee:           mongoOrder.Fee,
			IsPrime:       mongoOrder.IsPrime,
			LimitPrice:    mongoOrder.LimitPrice,
			MarketTime:    mongoOrder.MarketTime,
			Notional:      mongoOrder.Notional,
			OrderId:       mongoOrder.OrderID,
			OrderType:     mongoOrder.OrderType,
			Quantity:      mongoOrder.Quantity,
			QuantityShares:mongoOrder.QuantityShares,
			QuantityType:  mongoOrder.QuantityType,
			Status:        mongoOrder.Status,
			Ticker:        mongoOrder.Ticker,
			TradingType:   mongoOrder.TradingType,
			UserId:        mongoOrder.UserID,
		}

		order.CreatedAt = mongoOrder.CreatedAt.Time().Format(time.RFC3339)
		order.UpdatedAt = mongoOrder.UpdatedAt.Time().Format(time.RFC3339)
		orders = append(orders, order)
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
