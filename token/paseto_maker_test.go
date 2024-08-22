package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.Randomstring(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker((util.Randomstring(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	fmt.Println(err)
	require.Error(t, err, ErrExpiredToken)
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker((util.Randomstring(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, err := maker.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// tamper with the token
	token = token[0 : len(token)-1]
	
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Error(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}