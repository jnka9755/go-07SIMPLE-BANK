package gapi

import (
	"context"

	db "github.com/jnka9755/go-07SIMPLE-BANK/db/sqlc"
	"github.com/jnka9755/go-07SIMPLE-BANK/pb"
	"github.com/jnka9755/go-07SIMPLE-BANK/util"
	"github.com/jnka9755/go-07SIMPLE-BANK/validations"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)

	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %s", err)
	}

	response := &pb.CreateUserResponse{

		User: convertToUser(user),
	}
	return response, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validations.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolationError("username", err))
	}

	if err := validations.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolationError("password", err))
	}

	if err := validations.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolationError("full_name", err))
	}

	if err := validations.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolationError("email", err))
	}

	return violations
}
