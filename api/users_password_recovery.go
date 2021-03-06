package api

import (
	"context"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dominik-zeglen/inkster/core"
	gql "github.com/graph-gophers/graphql-go"
)

type UserChangePasswordArgs struct {
	ID       gql.ID
	Password string
}

func (res *Resolver) ChangeUserPassword(
	ctx context.Context,
	args UserChangePasswordArgs,
) (bool, error) {
	if !checkPermission(ctx) {
		return false, errNoPermissions
	}
	localID, err := fromGlobalID("user", string(args.ID))
	if err != nil {
		return false, err
	}
	_, err = res.dataSource.UpdateUser(localID, core.UserInput{
		Password: &args.Password,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

type ResetUserPasswordArgs struct {
	Password string
	Token    string
}

func (res *Resolver) ResetUserPassword(
	ctx context.Context,
	args ResetUserPasswordArgs,
) (bool, error) {
	tokenObject, err := jwt.ParseWithClaims(
		args.Token,
		&ActionTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
				return nil, errors.New("Invalid signing method")
			}

			claims, ok := token.Claims.(*ActionTokenClaims)
			if !ok {
				return nil, errors.New("Invalid token claims")
			}

			user, err := res.dataSource.GetUser(claims.ID)
			if err != nil {
				return nil, err
			}

			key := fmt.Sprintf("%x", user.Password)

			return []byte(key), nil
		},
	)
	if err != nil {
		return false, err
	}

	if claims, ok := tokenObject.Claims.(*ActionTokenClaims); ok {
		_, err = res.dataSource.UpdateUser(claims.ID, core.UserInput{
			Password: &args.Password,
		})
		return true, nil
	}
	return false, nil
}

type SendUserPasswordResetTokenArgs struct {
	Email string
}

func (res *Resolver) SendUserPasswordResetToken(
	ctx context.Context,
	args SendUserPasswordResetTokenArgs,
) (bool, error) {
	user, err := res.dataSource.GetUserByEmail(args.Email)
	if err != nil {
		return false, nil
	}

	claims := ActionTokenClaims{
		ID:            user.ID,
		AllowedAction: RESET_PASSWORD,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := fmt.Sprintf("%x", user.Password)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return false, err
	}

	err = res.mailer.Send(args.Email, "Inkster reset password", tokenString)
	if err != nil {
		return false, err
	}

	return true, nil
}
