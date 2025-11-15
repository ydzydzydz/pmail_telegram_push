package repository

import (
	"errors"

	"xorm.io/xorm"
)

// IRepository 仓库接口
type IRepository[T any] interface {
	Create(item *T) error
	Update(userID int, item *T) error
	FindOne(userID int) (*T, error)
	Exist(userID int) bool
}

// Repository 仓库实现
var _ IRepository[any] = (*Repository[any])(nil)

// Repository 仓库实现
type Repository[T any] struct {
	db *xorm.Engine
}

// NewRepository 创建仓库实例
func NewRepository[T any](db *xorm.Engine) *Repository[T] {
	return &Repository[T]{db: db}
}

// Create 创建记录
func (r *Repository[T]) Create(item *T) error {
	_, err := r.db.Insert(item)
	return err
}

// Update 更新记录
func (r *Repository[T]) Update(userID int, item *T) error {
	_, err := r.db.Where("user_id = ?", userID).AllCols().Update(item)
	return err
}

// Exist 检查记录是否存在
func (r *Repository[T]) Exist(userID int) bool {
	has, _ := r.db.Where("user_id = ?", userID).Exist(new(T))
	return has
}

// FindOne 获取记录
func (r *Repository[T]) FindOne(userID int) (*T, error) {
	item := new(T)
	has, err := r.db.Where("user_id = ?", userID).Get(item)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("item not found")
	}
	return item, nil
}
