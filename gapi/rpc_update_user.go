package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/jnka9755/go-07SIMPLE-BANK/db/sqlc"
	"github.com/jnka9755/go-07SIMPLE-BANK/pb"
	"github.com/jnka9755/go-07SIMPLE-BANK/util"
	"github.com/jnka9755/go-07SIMPLE-BANK/validations"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)

	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {

		hashedPassword, err := util.HashPassword(req.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "cannot update user: %s", err)
	}

	response := &pb.UpdateUserResponse{

		User: convertToUser(user),
	}
	return response, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validations.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolationError("username", err))
	}

	if req.Password != nil {
		if err := validations.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolationError("password", err))
		}
	}

	if req.FullName != nil {
		if err := validations.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolationError("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validations.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolationError("email", err))
		}
	}

	return violations
}
