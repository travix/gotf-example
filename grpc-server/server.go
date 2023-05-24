package grpc_server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/travix/gotf-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var users = map[string]*pb.User{
	"someone": {
		Username: "someone",
		Email:    "someone@example.com",
	},
}

var groups = map[string]*pb.Group{
	"somegroup": {
		Name:  "somegroup",
		Email: proto.String("somegroup@example.com"),
		Users: []*pb.User{
			{
				Username: "someone",
				Email:    "someone@example.com",
			},
		},
	},
}

type Server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedGroupServiceServer
}

func (s Server) RegisterGRPC(server *grpc.Server) {
	pb.RegisterGroupServiceServer(server, s)
	pb.RegisterUserServiceServer(server, s)
}

func (s Server) GetUser(_ context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	for _, user := range users {
		if user.Username == request.Username {
			log.Info().Msgf("Found user %s", user.Username)
			return user, nil
		}
	}
	return nil, status.Error(codes.NotFound, fmt.Sprintf("user with username %s not found", request.Username))
}

func (s Server) CreateUser(_ context.Context, request *pb.User) (*pb.User, error) {
	_, ok := users[request.Username]
	users[request.Username] = request
	if ok {
		log.Info().Msgf("updated user %s", request.Username)
	} else {
		log.Info().Msgf("created user %s", request.Username)
	}
	return users[request.Username], nil
}

func (s Server) ListUsers(context.Context, *pb.Empty) (*pb.Users, error) {
	resp := &pb.Users{}
	for _, user := range users {
		resp.Users = append(resp.Users, user)
	}
	log.Info().Msgf("Listed %d users", len(resp.Users))
	return resp, nil
}

func (s Server) UpdateUser(_ context.Context, request *pb.User) (*pb.User, error) {
	users[request.Username] = request
	log.Info().Msgf("Updated user %s", request.Username)
	return request, nil
}

func (s Server) DeleteUser(_ context.Context, request *pb.User) (*pb.Empty, error) {
	delete(users, request.Username)
	log.Info().Msgf("Deleted user %s", request.Username)
	return &pb.Empty{}, nil
}

func (s Server) GetGroup(_ context.Context, request *pb.GetGroupRequest) (*pb.Group, error) {
	for _, group := range groups {
		if group.Name == request.Name {
			log.Info().Msgf("Found group %s", group.Name)
			return group, nil
		}
	}
	return nil, status.Error(codes.NotFound, fmt.Sprintf("group with name %s not found", request.Name))
}

func (s Server) CreateGroup(_ context.Context, request *pb.Group) (*pb.Group, error) {
	_, ok := groups[request.Name]
	groups[request.Name] = request
	if ok {
		log.Info().Msgf("updated group %s", request.Name)
	} else {
		log.Info().Msgf("created group %s", request.Name)
	}
	return groups[request.Name], nil
}

func (s Server) ListGroups(context.Context, *pb.Empty) (*pb.Groups, error) {
	resp := &pb.Groups{}
	for _, group := range groups {
		resp.Groups = append(resp.Groups, group)
	}
	log.Info().Msgf("Listed %d groups", len(resp.Groups))
	return resp, nil
}

func (s Server) UpdateGroup(_ context.Context, request *pb.Group) (*pb.Group, error) {
	groups[request.Name] = request
	log.Info().Msgf("Updated group %s", request.Name)
	return request, nil
}

func (s Server) DeleteGroup(_ context.Context, request *pb.Group) (*pb.Empty, error) {
	delete(groups, request.Name)
	log.Info().Msgf("Deleted group %s", request.Name)
	return &pb.Empty{}, nil
}
