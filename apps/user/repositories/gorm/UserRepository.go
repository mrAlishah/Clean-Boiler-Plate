package repositories

import (
	"boilerplate/core/infrastructures"
	"boilerplate/core/models"
	"boilerplate/core/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRepository -> database structure
type UserRepository struct {
	db     infrastructures.GormDB
	logger infrastructures.Logger
}

// NewUserRepository -> creates a new User repository
func NewUserRepository(db infrastructures.GormDB, logger infrastructures.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// Save -> User
func (c UserRepository) Create(User *models.User) error {
	return c.db.DB.Create(User).Error
}

func (c UserRepository) FindByField(field string, value string) (user models.User, err error) {
	err = c.db.DB.Where(fmt.Sprintf("%s= ?", field), value).First(&user).Error
	return
}

func (c UserRepository) DeleteByID(id uint) error {
	user := models.User{}
	c.db.DB.Where("id=?", id).First(&user)
	return c.db.DB.Delete(&user).Error
}

func (c UserRepository) IsExist(field string, value string) (bool, error) {
	_, err := c.FindByField(field, value)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// GetAllUser -> Get All users
func (c UserRepository) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	err := queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return users, totalRows, err
}

//update a single column by user model
func (c UserRepository) UpdateColumn(user *models.User, column string, value interface{}) error {
	return c.db.DB.Model(user).Update(column, value).Error
}

func (ur UserRepository) GetAuthenticatedUser(c *gin.Context) (models.User, error) {
	userId := c.MustGet("userId").(string)
	return ur.FindByField("id", userId)
}