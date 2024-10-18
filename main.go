package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jnka9755/go-07SIMPLE-BANK/api"
	db "github.com/jnka9755/go-07SIMPLE-BANK/db/sqlc"
	_ "github.com/jnka9755/go-07SIMPLE-BANK/doc/statik"
	"github.com/jnka9755/go-07SIMPLE-BANK/gapi"
	"github.com/jnka9755/go-07SIMPLE-BANK/pb"
	"github.com/jnka9755/go-07SIMPLE-BANK/util"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load configurations:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal("cannot create migration:", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot migrate:", err)
	}

	log.Println("database migration done")
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Println("gRPC server listening on", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()

	if err != nil {
		log.Fatal("cannot create statikFS:", err)
	}

	swagerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swagerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Println("HTTP gateway server listening on", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
