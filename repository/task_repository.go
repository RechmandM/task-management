package repository

import (
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/models"
)

// ======================================================
// Task Repository
// Berisi seluruh proses akses database (CRUD)
// menggunakan GORM.
// ======================================================

// Mengambil daftar task berdasarkan parameter filter,
// sorting dan pagination.
// Method ini juga mengembalikan total data untuk
// kebutuhan perhitungan total halaman (pagination).
func GetAllTasks(keyword string, status string, assignee string, page int, limit int, sort string) ([]models.Task, int64, error) {

	var tasks []models.Task
	var totalRecords int64

	db := config.DB.Model(&models.Task{})

	// Terapkan filter sesuai query parameter yang dikirim client.
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

	// Menghitung total data sebelum pagination diterapkan.
	// Total data digunakan untuk menghitung total halaman.
	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Menentukan urutan data berdasarkan parameter sort.
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

	// Menghitung offset berdasarkan halaman yang dipilih.
	offset := (page - 1) * limit

	// Mengambil data sesuai filter, sorting dan pagination.
	err := db.
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error

	return tasks, totalRecords, err
}

// Menyimpan data task baru ke database.
func CreateTask(task *models.Task) error {
	// Menyimpan data task baru ke database.
	return config.DB.Create(task).Error

}

// Mengambil satu data task berdasarkan ID.
func GetTaskByID(id uint) (models.Task, error) {

	var task models.Task

	err := config.DB.First(&task, id).Error

	return task, err

}

// Memperbarui data task yang sudah ada.
func UpdateTask(task *models.Task) error {

	return config.DB.Save(task).Error

}

// Melakukan soft delete.
// GORM hanya mengisi kolom deleted_at sehingga
// data masih dapat dipulihkan jika diperlukan.
func DeleteTask(id uint) error {

	// lebih pendek jika menggunakan ini
	// return config.DB.Delete(&models.Task{}, id).Error

	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&task).Error

}
