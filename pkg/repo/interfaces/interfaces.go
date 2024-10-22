package interfaces

import "github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"

type MaterialRepoInter interface {
	//Material
	AddMaterial(Material *model.Material) (uint, error)
	FindMaterialByID(MaterialID uint) (*model.Material, error)
	FindAllMaterial() (*[]model.Material, error)
	UpdateMaterial(Material *model.Material) error
	DeleteMaterial(MaterialID uint) error

	//Items
	CreateItem(Item *model.Item) (uint, error)
	FindItemByID(ItemsID uint) (*model.Item, error)
	FindAllItem() (*[]model.Item, error)
	UpdateItem(Item *model.Item) error
	DeletItem(ItemID uint) error

	//Orders
	CreateOrders(Orders *model.Order) (uint, error)
	FindOrdersByID(OrdersID uint) (*model.Order, error)
	UpdateOrders(Orders *model.Order) error
	DeleteOrders(OrdersID uint) error
	FindAllOrders() (*[]model.Order, error)
	FindOrdersByUser(userID uint) (*[]model.Order, error)
	FindOrder(userID, ItemID uint) (*model.Order, error)
}