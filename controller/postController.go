package controller

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ketanpolawar/blogbackend/database"
	"github.com/ketanpolawar/blogbackend/models"
	"github.com/ketanpolawar/blogbackend/util"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		fmt.Println("Unable to parse body")
	}
	if err := database.DB.Create(&blogpost).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})

	}
	return c.JSON(fiber.Map{
		"message": "Congratulation ,Your post is live",
	})
}
func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getblog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getblog)
	database.DB.Model(&models.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"data": getblog,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(int(total) / limit)),
		},
	})

}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blogpost models.Blog
	database.DB.Where("id=?", id).Preload("user").First(&blogpost)
	return c.JSON(fiber.Map{
		"data": blogpost,
	})

}

func Updatepost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}
	if err := c.BodyParser(&blog); err != nil {
		fmt.Println("Unable to parse body")
	}
	database.DB.Model(&blog).Updates(blog)
	return c.JSON(fiber.Map{
		"message": "Updated successfully",
	})

}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := util.Parsejwt(cookie)
	var blog []models.Blog
	database.DB.Model(&blog).Where("user_id=?", id).Preload("user").Find(&blog)

	return c.JSON(blog)

}

func DeletePost(c *fiber.Ctx) error {
	// Convert the id parameter to an integer
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	// Create a blog instance with the given id
	blog := models.Blog{
		Id: uint(id),
	}

	// Perform the delete operation
	deleteQuery := database.DB.Delete(&blog)

	// Check if the record was not found
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Oops, record not found",
		})
	}

	// Check for other possible errors
	if deleteQuery.Error != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Failed to delete the post",
		})
	}

	// Return success message
	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
