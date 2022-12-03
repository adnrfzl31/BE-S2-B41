package repositories

import (
	"BE-S2-B41/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	AddOrder(orders models.Order) (models.Order, error)
	GetOrder(ID int) (models.Order, error)
	FindOrder() ([]models.Order, error)
	DelOrder(orders models.Order) (models.Order, error)
	UpdateOrder(orders models.Order) (models.Order, error)
	GetProductOrder(ID int) (models.Product, error)
	GetTopingOrder(ID []int) ([]models.Toping, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddOrder(orders models.Order) (models.Order, error) {
	err := r.db.Create(&orders).Error

	return orders, err
}

func (r *repository) GetOrder(ID int) (models.Order, error) {
	var orders models.Order
	err := r.db.Preload("Product").Preload("Toping").Preload("Buyyer").First(&orders, ID).Error

	return orders, err
}

func (r *repository) FindOrder() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Product").Preload("Toping").Preload("Buyyer").Find(&orders).Error

	return orders, err
}

func (r *repository) DelOrder(orders models.Order) (models.Order, error) {
	err := r.db.Delete(&orders).Error

	return orders, err
}

func (r *repository) UpdateOrder(orders models.Order) (models.Order, error) {
	err := r.db.Save(&orders).Error

	return orders, err
}

func (r *repository) GetProductOrder(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ID).Error

	return product, err
}

func (r *repository) GetTopingOrder(ID []int) ([]models.Toping, error) {
	var toping []models.Toping
	err := r.db.Find(&toping, ID).Error

	return toping, err
}
