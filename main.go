package main

import (
	"context"
	"database/sql"
	"github.com/getsentry/sentry-go"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/the-medo/golang-migrate-objects/migrator"
	"github.com/the-medo/talebound-backend/api/srv"
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
	"strings"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config:")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		TracesSampleRate: config.SentryTracesSampleRate,
	})
	if err != nil {
		log.Fatal().Msgf("sentry.Init: %s", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db:")
	}

	runDBMigration(config.MigrationURL, config.DBSource, config.MigrationObjectsURL, config.MigrationCreateObjectsFilename, config.MigrationDropObjectsFilename)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string, migrationObjectsURL string, createObjectsFilename string, dropObjectsFilename string) {

	dbConn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect! ")
	}

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	migration, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"talebound", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new migrate instance! ")
	}

	path := strings.TrimPrefix(migrationObjectsURL, "file://")
	log.Info().Msgf("Path: %s", path)

	mg, err := migrator.New(&migrator.Config{
		DB:                    dbConn,
		DbObjectPath:          path,
		MigrationFilesPath:    strings.TrimPrefix(migrationURL, "file://"),
		CreateObjectsFilename: createObjectsFilename,
		DropObjectsFilename:   dropObjectsFilename,
	}, migration)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new mg instance! ")
		return
	}

	err = mg.RunAll()
	if err != nil {
		return
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
	server, err := srv.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(util.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterLocationsServer(grpcServer, server)
	pb.RegisterMapsServer(grpcServer, server)
	pb.RegisterModulesServer(grpcServer, server)
	pb.RegisterEntitiesServer(grpcServer, server)
	pb.RegisterTagsServer(grpcServer, server)
	pb.RegisterUsersServer(grpcServer, server)
	pb.RegisterMenusServer(grpcServer, server)
	pb.RegisterPostsServer(grpcServer, server)
	pb.RegisterEvaluationsServer(grpcServer, server)
	pb.RegisterImagesServer(grpcServer, server)
	pb.RegisterAuthServer(grpcServer, server)
	pb.RegisterWorldsServer(grpcServer, server)
	pb.RegisterSystemsServer(grpcServer, server)
	pb.RegisterCharactersServer(grpcServer, server)
	pb.RegisterQuestsServer(grpcServer, server)
	pb.RegisterFetcherServer(grpcServer, server)

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
				UseProtoNames:   true,
				EmitUnpopulated: true,
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

	server, err := srv.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server:")
	}

	err = pb.RegisterLocationsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterMapsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterModulesHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterEntitiesHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterTagsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterUsersHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterMenusHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterPostsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterEvaluationsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterImagesHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterAuthHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterWorldsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterSystemsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterCharactersHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterQuestsHandlerServer(ctx, grpcMux, server)
	err = pb.RegisterFetcherHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	var muxWithCORS http.Handler

	//TODO - check Fetch-Ids header - should be in allowed or exposed?
	if config.Environment == "development" {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   []string{config.FullDomain},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Set-Cookie", "Fetch-Ids"},
			ExposedHeaders:   []string{"Fetch-Ids"},
			AllowCredentials: true,
		})
		muxWithCORS = corsMiddleware.Handler(mux)
	} else {
		//TODO PROD - check if everything works in production
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
	handler := util.HttpLogger(muxWithCORS)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start HTTP gateway server:")
	}
}
