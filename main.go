package main

import (
	"BlogApplication/repository"
	"BlogApplication/server"
	"BlogApplication/service"
	"encoding/json"
	"log"
	"net"

	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://blog-database:27017"))
	if err != nil {
		return nil
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	println("Connected to MongoDB!")

	return client
}

func initTracer() (func(context.Context) error, error) {

	jaegerExporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint("jaeger:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "blog-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func startServer(blogService *service.BlogService, commentService *service.CommentService, reportService *service.ReportService, natsConn *nats.Conn) {

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	blogMicroservice := &server.BlogMicroservice{
		BlogService:    blogService,
		CommentService: commentService,
		ReportService:  reportService,
		NatsConn:       natsConn,
	}

	server.RegisterBlogMicroserviceServer(grpcServer, blogMicroservice)

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port :8088")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
func handleRollback(nc *nats.Conn, commentService *service.CommentService) {
	nc.Subscribe("comment.creation.rollback", func(m *nats.Msg) {
		var event struct {
			CommentID int `json:"comment_id"`
		}
		json.Unmarshal(m.Data, &event)

		// Delete the comment
		err := commentService.Delete(context.Background(), int64(event.CommentID))
		if err != nil {
			log.Printf("Failed to rollback comment creation: %v", err)
		} else {
			log.Printf("Successfully rolled back comment creation with ID: %d", event.CommentID)
		}
	})
}
func main() {
	client := initDB()
	if client == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}
	shutdown, err := initTracer()
	if err != nil {
		log.Fatalf("FAILED TO INITIALIZE TRACER: %v", err)
	}
	defer shutdown(context.Background())

	conn := Conn()
	defer conn.Close()

	blogRepository := repository.NewBlogRepository(client)
	blogService := &service.BlogService{BlogRepository: blogRepository}

	commentRepository := repository.NewCommentRepository(client)
	commentService := &service.CommentService{CommentRepo: commentRepository}

	reportRepository := repository.NewReportRepository(client)
	reportService := &service.ReportService{ReportRepository: reportRepository}

	handleRollback(conn, commentService)
	startServer(blogService, commentService, reportService, conn)

	select {}
}
