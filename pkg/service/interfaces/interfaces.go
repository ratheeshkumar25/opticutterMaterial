package interfaces

import pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"

type MaterialServiceInte interface {
	// Service to handle material management
	AddMaterialService(p *pb.Material) (*pb.MaterialResponse, error)
	EditMaterialService(p *pb.Material) (*pb.Material, error)
	RemoveMaterialService(p *pb.MaterialID) (*pb.MaterialResponse, error)
	FindMaterialByIDService(p *pb.MaterialID) (*pb.Material, error)
	FindAllMaterialService(p *pb.MaterialNoParams) (*pb.MaterialList, error)

	// Service to handle item management
	AddItemService(p *pb.Item) (*pb.ItemResponse, error)
	EditItemService(p *pb.Item) (*pb.Item, error)
	RemoveItemService(p *pb.ItemID) (*pb.ItemResponse, error)
	FindItemByID(p *pb.ItemID) (*pb.Item, error)
	FindAllItem(p *pb.ItemNoParams) (*pb.ItemList, error)

	// Service to handle orders
	PlaceOrderService(p *pb.Order) (*pb.OrderResponse, error)
	FindAllOrdersSvc(p *pb.ItemNoParams) (*pb.OrderList, error)
	FindOrdersByUserSvc(p *pb.ItemID) (*pb.OrderList, error)
	FindOrderSvc(p *pb.ItemID) (*pb.Order, error)
}
