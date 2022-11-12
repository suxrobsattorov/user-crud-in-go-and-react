package repository

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morkid/paginate"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"user-crud-app/database/migrations"
	"user-crud-app/database/models"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(user models.User) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (repo *Repository) CreateUser(ctx *fiber.Ctx) error {
	user := models.User{}
	err := ctx.BodyParser(&user)

	if err != nil {
		err = ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})
		return err
	}

	errors := ValidateStruct(user)

	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if err := repo.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "User has been added", "data": user})
}

func (repo *Repository) UpdateUser(ctx *fiber.Ctx) error {
	user := models.User{}
	err := ctx.BodyParser(&user)

	if err != nil {
		err = ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})
		return err
	}

	errors := ValidateStruct(user)

	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := repo.DB
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID cannot be empty"})
	}

	if db.Model(&user).Where("id = ?", id).Updates(&user).RowsAffected == 0 {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get User with given id"})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "User successfully updated"})
}

func (repo *Repository) DeleteUser(ctx *fiber.Ctx) error {
	userModel := migrations.Users{}
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID cannot be empty"})
	}

	err := repo.DB.Delete(userModel, id)

	if err.Error != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete boo"})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "User delete successfully"})
}

func (repo *Repository) GetUsers(ctx *fiber.Ctx) error {
	db := repo.DB
	model := db.Model(&migrations.Users{})

	pg := paginate.New(&paginate.Config{
		DefaultSize:        20,
		CustomParamEnabled: true,
	})

	page := pg.With(model).Request(ctx.Request()).Response(&[]migrations.Users{})

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": page})
}

func (repo *Repository) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userModel := &migrations.Users{}

	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
	}

	err := repo.DB.Where("id = ?", id).First(userModel).Error

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the user"})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "User id fetched successfully", "data": userModel})
}
