package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,  errorReponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}
	arg := db.CreateUserParams {
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name(){
				case "unique_violation":
					ctx.JSON(http.StatusForbidden, errorReponse(err))
					return
			}		
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	// Avoid returning the hashed password
	response := CreateUserResponse {
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}