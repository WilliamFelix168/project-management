package controllers

import (
	"math"
	"strconv"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// membuat struct UserController yang berisi dependency ke UserService
// membuat request handler untuk user
type UserController struct {
	service services.UserService
	//dependency yang dibutuhkan controller utk logika bisnis
}

// fungsi untuk mengembalikan objek UserController dengan dependency UserService
// tujuannya untuk menginisialisasi controller dengan service yang diberikan
func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	//handler untuk register user
	user := new(models.User)

	// untuk memparsing body request ke struct user
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing data user", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)

	return utils.Success(ctx, "Register Success", userResp)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	//handler untuk login user
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// untuk memparsing body request ke struct body
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)

	return utils.Success(ctx, "Login Successful", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResp,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User Not Found", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)

	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "User Found", userResp)
}

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	// /users/page?page=1&limit=10&sort=name&filter=triady
	// 100/ 10 = 10 pages

	//mengambil query parameter untuk filtering, sorting, pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to get users", err.Error())
	}

	var usersResp []models.UserResponse
	_ = copier.Copy(&usersResp, &users)

	meta := utils.PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: int(total),
		//TotalPages = total / limit 100 / 10 = 10 ||  100 / 3 = 33.33 -> 34
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Filter:     filter,
		Sorting:    sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "No Users Found", usersResp, meta)
	}

	return utils.SuccessPagination(ctx, "Users Found", usersResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	//uuid.Parse validates the format of UUID
	publicID, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	user.PublicID = publicID

	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to update user", err.Error())
	}

	userUpdated, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to get updated user", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Error Parsing Data", err.Error())
	}

	return utils.Success(ctx, "User Updated Successfully", userResp)

}
