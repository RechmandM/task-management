package repository

import (
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/models"
)

// func GetAllTasks() ([]models.Task, error) {
func GetAllTasks(keyword string, status string, assignee string, page int, limit int, sort string) ([]models.Task, int64, error) {

	var tasks []models.Task
	var totalRecords int64

	db := config.DB.Model(&models.Task{})

	// Filter Keyword
	if keyword != "" {
		db = db.Where(
			"title LIKE ? OR description LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	// Filter Status
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// Filter Assignee
	if assignee != "" {
		db = db.Where("assignee = ?", assignee)
	}

	// 4. Hitung TOTAL DATA (sebelum sorting, offset, & limit)
	// GORM otomatis mengabaikan data dengan DeletedAt != NULL
	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	switch sort {
	case "created_at_asc":
		db = db.Order("created_at ASC")
	case "created_at_desc":
		db = db.Order("created_at DESC")
	case "title_asc":
		db = db.Order("title ASC")
	case "title_desc":
		db = db.Order("title DESC")
	default:
		db = db.Order("id DESC")
	}

	// Pagination
	offset := (page - 1) * limit

	// get count

	err := db.
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error

	return tasks, totalRecords, err
}

// create
func CreateTask(task *models.Task) error {

	return config.DB.Create(task).Error

}

// for update
func GetTaskByID(id uint) (models.Task, error) {

	var task models.Task

	err := config.DB.First(&task, id).Error

	return task, err

}

func UpdateTask(task *models.Task) error {

	return config.DB.Save(task).Error

}

func DeleteTask(id uint) error {

	// lebih pendek
	// return config.DB.Delete(&models.Task{}, id).Error

	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&task).Error

}
