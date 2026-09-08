package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct {
  	ID int
  	Title string
	GPU string
	RAM string
	Storage string
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID:    1,
			Title: "2080 Ti",
			GPU: "2080 Ti",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
		},
		{
			ID:    2,
			Title: "3090",
			GPU: "3090",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
		},
		{
			ID:    3,
			Title: "4090",
			GPU: "4090",
			RAM: "32768 МБ ОЗУ",
			Storage: "160 ГБ SSD",
		},
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}