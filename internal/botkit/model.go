package botkit

import (
	"context"
	"github.com/vvv9912/sddb"
)

type OrderStorager interface {
	AddOrder(ctx context.Context, order sddb.Orders) error
	GetOrdersByTgID(ctx context.Context, tgId int64) ([]sddb.Orders, error)
}
type CorzinaStorager interface {
	AddShopCart(ctx context.Context, ShopCart sddb.ShopCart) error
	GetShopCartByTgID(ctx context.Context, tgId int64) ([]sddb.ShopCart, error)
	UpdateShopCartByTgId(ctx context.Context, tgId int64, article int, quantity int) error
	GetShopCartByTgIdAndArticle(ctx context.Context, tgId int64, article int) (sddb.ShopCart, error)
	GetShopCartDetailByTgId(ctx context.Context, tgId int64) ([]sddb.DbCorzineCatalog, error)
	DeleteShopCartByTgId(ctx context.Context, tgId int64) error
	DeleteShopCartByTgIdAndArticle(ctx context.Context, tgId int64, article int) error
}
type ProductsStorager interface {
	AddProduct(ctx context.Context, product sddb.Products) error

	ChangeProductByArticle(ctx context.Context, product sddb.Products) error

	GetCatalogNames(ctx context.Context) ([]string, error)
	GetCatalogNamesIsAvailable(ctx context.Context) ([]string, error)

	GetAllProducts(ctx context.Context) ([]sddb.Products, error)

	GetProductsByCatalogIsAvailable(ctx context.Context, ctlg string) ([]sddb.Products, error)

	GetProductsByCatalog(ctx context.Context, ctlg string) ([]sddb.Products, error)

	GetProductByArticle(ctx context.Context, article int) (sddb.Products, error)
}
type UsersStorager interface {
	GetStatusUserByTgID(ctx context.Context, tgID int64) (int, int, error)
	AddUser(ctx context.Context, users sddb.Users) error
	UpdateStateByTgID(ctx context.Context, tgId int64, state int) error
	//GetCorzinaByTgID(ctx context.Context, tgID int64) ([]int64, error)
	//UpdateShopCartByTgId(ctx context.Context, tgId int64, corzina []int64) error
}
