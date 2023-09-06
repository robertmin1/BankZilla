package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := RandomUser(t)

	CreateUserParams := db.CreateUserParams{
		Username: 	 user.Username,
		HashedPassword: user.HashedPassword,
		FullName: 	 user.FullName,
		Email: 		 user.Email,
	}

	requestparams := CreateUserRequest{
		Username: user.Username,
		Password: password,
		FullName: user.FullName,
		Email:    user.Email,
	}
	
	testcases := []struct {
		name         string
		Createparams db.CreateUserParams
		reqparams    CreateUserRequest
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			Createparams: CreateUserParams,
			reqparams:    requestparams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(user, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewSever(store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/users", construstJson(t, tc.reqparams))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}

func RandomUser(t *testing.T) (user db.User, password string) {
	password = util.Randomstring(6)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return user, password
}