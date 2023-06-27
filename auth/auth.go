// SPDX-License-Identifier: ice License 1.0

package auth

import (
	"context"

	"github.com/pkg/errors"

	firebaseAuth "github.com/ice-blockchain/wintr/auth/internal/firebase"
	iceAuth "github.com/ice-blockchain/wintr/auth/internal/ice"
	"github.com/ice-blockchain/wintr/time"
)

func New(ctx context.Context, applicationYAMLKey string) Client {
	return &auth{
		fb:  firebaseAuth.New(ctx, applicationYAMLKey),
		ice: iceAuth.New(applicationYAMLKey),
	}
}

func (a *auth) VerifyToken(ctx context.Context, token string) (*Token, error) {
	var authToken *Token
	if err := iceAuth.DetectIceToken(token); err != nil {
		authToken, err = a.fb.VerifyToken(ctx, token)

		return authToken, errors.Wrapf(err, "can't verify fb token:%v", token)
	}
	authToken, err := a.ice.VerifyToken(token)

	return authToken, errors.Wrapf(err, "can't verify ice token:%v", token)
}

func (a *auth) UpdateCustomClaims(ctx context.Context, userID string, customClaims map[string]any) error {
	if usr, err := a.fb.GetUser(ctx, userID); err == nil && usr != nil {
		return errors.Wrapf(a.fb.UpdateCustomClaims(ctx, userID, customClaims), "failed to update custom claims for user:%v using firebase auth", userID)
	}

	return nil
}

func (a *auth) DeleteUser(ctx context.Context, userID string) error {
	if usr, err := a.fb.GetUser(ctx, userID); err == nil && usr != nil {
		return errors.Wrapf(a.fb.DeleteUser(ctx, userID), "failed to delete user:%v using firebase auth", userID)
	}

	return nil
}

func (a *auth) GenerateTokens( //nolint:revive // We need to have these parameters.
	now *time.Time, userID, deviceUniqueID, email string, hashCode, seq int64, claims map[string]any,
) (accessToken, refreshToken string, err error) {
	accessToken, refreshToken, err = a.ice.GenerateTokens(now, userID, deviceUniqueID, email, hashCode, seq, claims)
	err = errors.Wrapf(err, "can't generate tokens for userID:%v, email:%v", userID, email)

	return
}

func (a *auth) ParseToken(token string) (*IceToken, error) {
	res := new(IceToken)
	err := a.ice.VerifyTokenFields(token, res)

	return res, errors.Wrapf(err, "can't verify token fields for:%v", token)
}
