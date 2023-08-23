package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name         string
		accountID    int64
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// Sql No Rows Error
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		// Internal Server Error
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// Invalid Request
		{
			name:      "Invalid ID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewSever(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// hold response in recorder and request is the req*
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkReponse(t, recorder)
		},
		)
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	createparams := db.CreateAccountParams{
		Owner:    account.Owner,
		Balance:  0,
		Currency: account.Currency,
	}

	requestparams := createAccountRequest{
		Owner:    account.Owner,
		Currency: account.Currency,
	}

	testcases := []struct {
		name         string
		Createparams db.CreateAccountParams
		reqparams    createAccountRequest
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			Createparams: createparams,
			reqparams:    requestparams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(createparams)).
					Times(1).
					Return(account, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:         "Invalid Parameters",
			Createparams: db.CreateAccountParams{},
			reqparams:    createAccountRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "Internal Server Error",
			Createparams: createparams,
			reqparams:    requestparams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			request, err := http.NewRequest(http.MethodPost, "/accounts", construstJson(t, tc.reqparams))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)
		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	account := randomAccount()

	testcases := []struct {
		name         string
		accountID    int64
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Invalid Parameters",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	account := randomAccount()

	updateAccountparams := db.UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance,
	}

	requestupdateparmas := updateAccountRequest{
		ID:      account.ID,
		Balance: account.Balance,
	}

	testcases := []struct {
		name          string
		updateAccount db.UpdateAccountParams
		requestupdate updateAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkReponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			updateAccount: updateAccountparams,
			requestupdate: requestupdateparmas,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(updateAccountparams)).
					Times(1).
					Return(account, nil)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:          "Invalid Parameters",
			updateAccount: updateAccountparams,
			requestupdate: updateAccountRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "Internal Server Error",
			updateAccount: updateAccountparams,
			requestupdate: requestupdateparmas,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(updateAccountparams)).
					Times(1).
					Return(db.Account{},sql.ErrConnDone)
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			request, err := http.NewRequest(http.MethodPut, "/accounts/update", construstJson(t, tc.requestupdate))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func construstJson(t *testing.T, data interface{}) *bytes.Buffer {
	body, err := json.Marshal(data)
	require.NoError(t, err)
	return bytes.NewBuffer(body)
}
