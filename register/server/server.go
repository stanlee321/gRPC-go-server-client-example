package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	// "gopkg.in/mgo.v2/bson"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"projects/elearning/microservice_basic/register/registerpb"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
}

type registerUserItem struct {

	ID 		 	primitive.ObjectID 	`bson:"_id,omitempty"`
	Email 		string 				`bson:"email"`
	Name 		string 				`bson:"name"`
	Username 	string				`bson:"username"`
	Password 	string 				`bson:"password"`
	LevelOfEducation string			`bson:"level_of_education"`
	Gender		string				`bson:"gender"`
	YearOfBirth	int32				`bson:"year_of_birth"`
	MailingAddress string			`bson:"mailing_address"`
	Goals		string				`bson:"goals"`
	Country		string				`bson:"country"`
	HonorCode  	bool				`bson:"honor_code"`
}



func (*server) CreateUser(ctx context.Context, req *registerUserpb.CreateUserRequest) (*registerUserpb.CreateUserResponse, error) {
	fmt.Println("Create User request")

	user := req.GetUser()

	data := registerUserItem{

		Email 		: user.GetEmail(),
		Name 		: user.GetName(),
		Username 	: user.GetUsername(),
		Password 	: user.GetPassword(), 
		LevelOfEducation : user.GetLevelOfEducation(),
		Gender		: user.GetGender(),
		YearOfBirth	: user.GetYearOfBirth(),
		MailingAddress : user.GetMailingAddress(),
		Goals		: user.GetGoals(),
		Country		: user.GetCountry(),		
		HonorCode  	: user.GetHonorCode(),

	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &registerUserpb.CreateUserResponse{
		User: &registerUserpb.User{
			Id			: oid.Hex(),
			Email 		: user.GetEmail(),
			Name 		: user.GetName(),
			Username 	: user.GetUsername(),
			Password 	: user.GetPassword(), 
			LevelOfEducation : user.GetLevelOfEducation(),
			Gender		: user.GetGender(),
			YearOfBirth	: user.GetYearOfBirth(),
			MailingAddress : user.GetMailingAddress(),
			Goals		: user.GetGoals(),
			Country		: user.GetCountry(),		
			HonorCode  	: user.GetHonorCode(),
		},
	}, nil

}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	// connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Blog Service Started")
	collection = client.Database("elearning").Collection("users")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	registerUserpb.RegisterRegisterUserServiceServer( s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	// First we close the connection with MongoDB:
	fmt.Println("Closing MongoDB Connection")
    	// client.Disconnect(context.TODO())	
	if err := client.Disconnect(context.TODO()); err != nil {
        	log.Fatalf("Error on disconnection with MongoDB : %v", err)
    	}
    	// Second step : closing the listener
    	fmt.Println("Closing the listener")
    	if err := lis.Close(); err != nil {
        	log.Fatalf("Error on closing the listener : %v", err)
	}
	// Finally, we stop the server
	fmt.Println("Stopping the server")
    	s.Stop()
    	fmt.Println("End of Program")
}