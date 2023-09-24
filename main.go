package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/the-medo/talebound-backend/api"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	_ "github.com/the-medo/talebound-backend/doc/statik"
	"github.com/the-medo/talebound-backend/mail"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config:")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db:")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new migrate instance: ")
		return
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up:")
	}

	log.Info().Msg("Migration successful")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewAwsSesSender(config.EmailSenderName, config.EmailSenderAddress, config.SmtpUsername, config.SmtpPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer, config)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(api.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterChatServer(grpcServer, server)
	pb.RegisterEvaluationsServer(grpcServer, server)
	pb.RegisterImagesServer(grpcServer, server)
	pb.RegisterMapsServer(grpcServer, server)
	pb.RegisterMenusServer(grpcServer, server)
	pb.RegisterPostTypesServer(grpcServer, server)
	pb.RegisterPostsServer(grpcServer, server)
	pb.RegisterTagsServer(grpcServer, server)
	pb.RegisterUsersServer(grpcServer, server)
	pb.RegisterVerifyServer(grpcServer, server)
	pb.RegisterWorldsServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener:")
	}

	log.Info().Msgf("Starting gRPC server at %s", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start gRPC server:")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
			return s[2:], true
		}),
		runtime.WithForwardResponseOption(util.CreateFilterTokensToCookies(config)),
		runtime.WithMetadata(util.CookieAnnotator),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server:")
	}

	err = pb.RegisterChatHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterEvaluationsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterImagesHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterMapsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterMenusHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterPostTypesHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterPostsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterTagsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterUsersHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterVerifyHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterWorldsHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	var muxWithCORS http.Handler

	if config.Environment == "development" {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   []string{config.FullDomain},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
			AllowCredentials: true,
		})
		muxWithCORS = corsMiddleware.Handler(mux)
	} else {
		muxWithCORS = cors.Default().Handler(mux)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create statik file system:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener:")
	}

	log.Info().Msgf("Starting HTTP gateway server at %s", listener.Addr().String())
	handler := api.HttpLogger(muxWithCORS)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start HTTP gateway server:")
	}
}
