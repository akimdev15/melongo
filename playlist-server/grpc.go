package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/akimdev15/melongo/playlist-server/internal/database"
	"github.com/akimdev15/melongo/playlist-server/proto"
	"google.golang.org/grpc"
)

const gRPCPORT = "50002"

type PlaylistServer struct {
	proto.UnimplementedPlaylistServiceServer
	DB *database.Queries
}

func (apiCfg *apiConfig) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPORT))
	if err != nil {
		log.Fatalf("Failed to listen for grpc %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterPlaylistServiceServer(grpcServer, &PlaylistServer{DB: apiCfg.DB})
	fmt.Printf("gRPC server start on PORT %s]\n", gRPCPORT)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC %v", err)
	}
}

func (playlistServer *PlaylistServer) CreatePlaylist(ctx context.Context, req *proto.CreatePlaylistRequest) (*proto.CreatePlaylistResponse, error) {
	fmt.Println("Received request to create playlist using gRPC")
	res := &proto.CreatePlaylistResponse{
		SpotifyPlaylistID: "playlistID",
		ExternalURL:       "TEST",
		Name:              "TEST",
	}
	return res, nil
}
