package handler

import (
	"context"
	"fmt"

	"log"

	"github.com/susilo001/simple-wallet-system/user/entity"
	pb "github.com/susilo001/simple-wallet-system/user/proto/user/v1"
	"github.com/susilo001/simple-wallet-system/user/service"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandler is used to implement UnimplementedUserServiceServer
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var usersProto []*pb.User
	for _, user := range users {
		usersProto = append(usersProto, &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.GetUsersResponse{
		User: usersProto,
	}, nil
}
func (u *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := u.userService.GetUserByID(ctx, int(req.GetId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pb.GetUserResponse{
		User: &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return res, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.MutationResponse, error) {
	createdUser, err := u.userService.CreateUser(ctx, &entity.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created user with ID %d", createdUser.ID),
	}, nil
}
func (u *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MutationResponse, error) {
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success update user with ID %d", req.User.Id),
	}, nil
}
func (u *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MutationResponse, error) {
	if err := u.userService.DeleteUser(ctx, int(req.GetId())); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success deleted user with ID %d", req.GetId()),
	}, nil
}
