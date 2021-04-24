package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var db *gorm.DB

func init() {
	tDB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	tDB.AutoMigrate(&Product{})
	db = tDB
}

// GetProduct 根据 pid(product_id) 获取 Product，不存在记录时返回 nil
func GetProduct(pid int) (*Product, error) {
	var product Product
	err := db.First(&product, pid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据不存在时 product 返回 nil
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("db.First meed error: %w", err)
	}
	return &product, nil
}

// CreateProduct 创建 Product
func CreateProduct(product *Product) error {
	return db.Create(product).Error
}
