package service

import (
	"fmt"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/utils"
)

const plywoodSize = 3000.0 // plywood sheet area
const wasteFactor = 0.1    // waste factor (10%)

// AddItemService implements interfaces.MaterialServiceInte.
func (m *MaterialService) AddItemService(p *pb.Item) (*pb.ItemResponse, error) {

	//material,err := m.Repo.FindMaterialByID()
	newItem := model.Item{
		ItemName:    p.Item_Name,
		MaterialID:  uint(p.Material_ID),
		Length:      uint(p.Length),
		Width:       uint(p.Width),
		FixedSizeID: uint(p.Fixed_Size_ID),
		IsCustom:    p.Is_Custom,
		UserID:      uint(p.User_ID),
	}
	itemID, err := m.Repo.CreateItem(&newItem)
	if err != nil {
		return &pb.ItemResponse{
			Status:  pb.ItemResponse_ERROR,
			Message: "failed to add item",
			Payload: &pb.ItemResponse_Error{
				Error: err.Error(),
			},
		}, err
	}

	// // Fetch material to get price per unit
	material, err := m.Repo.FindMaterialByID(newItem.MaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to find material: %v", err)
	}

	// Calculate the estimated price in a goroutine and update the database
	go func() {
		estPrice, err := utils.CalculateEstPrice(&newItem, material.Price, plywoodSize, wasteFactor)
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Printf("failed to calculate price: %v\n", err)
			return
		}
		// Update the item with the estimated price
		newItem.EstPrice = float32(estPrice)
		err = m.Repo.UpdateItem(&newItem)
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Printf("failed to update item price: %v\n", err)
		}
	}()

	return &pb.ItemResponse{
		Status:  pb.ItemResponse_OK,
		Message: "Item added successfully",
		Payload: &pb.ItemResponse_Data{
			Data: fmt.Sprintf("ItemID:%d", itemID),
		},
	}, nil

}

// EditItemService implements interfaces.MaterialServiceInte.
func (m *MaterialService) EditItemService(p *pb.Item) (*pb.Item, error) {
	item, err := m.Repo.FindItemByID(uint(p.Item_ID))
	if err != nil {
		return nil, err
	}
	item.ItemName = p.Item_Name
	item.MaterialID = uint(p.Material_ID)
	item.Length = uint(p.Length)
	item.Width = uint(p.Width)
	item.FixedSizeID = uint(p.Fixed_Size_ID)
	item.IsCustom = p.Is_Custom
	item.UserID = uint(p.User_ID)

	err = m.Repo.UpdateItem(item)
	if err != nil {
		return nil, err
	}

	// Fetch material to get price per unit
	material, err := m.Repo.FindMaterialByID(item.MaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to find material: %v", err)
	}

	// Calculate the estimated price in a goroutine and update the database
	go func() {
		estPrice, err := utils.CalculateEstPrice(item, material.Price, plywoodSize, wasteFactor)
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Printf("failed to calculate price: %v\n", err)
			return
		}
		// Update the item with the estimated price
		item.EstPrice = float32(estPrice)
		err = m.Repo.UpdateItem(item)
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Printf("failed to update item price: %v\n", err)
		}
	}()

	return p, nil
}

// FindAllItem implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindAllItem(p *pb.ItemNoParams) (*pb.ItemList, error) {
	result, err := m.Repo.FindAllItem()
	if err != nil {
		return nil, err
	}
	var items []*pb.Item

	for _, item := range *result {
		pbItem := &pb.Item{
			Item_ID:         uint32(item.ID),
			Item_Name:       item.ItemName,
			Length:          uint32(item.Length),
			Width:           uint32(item.Width),
			Fixed_Size_ID:   uint32(item.FixedSizeID),
			Is_Custom:       item.IsCustom,
			Estimated_Price: item.EstPrice,
			User_ID:         uint32(item.UserID),
		}
		items = append(items, pbItem)
	}
	return &pb.ItemList{
		Items: items,
	}, nil
}

// FindItemByID implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindItemByID(p *pb.ItemID) (*pb.Item, error) {
	item, err := m.Repo.FindItemByID(uint(p.ID))
	if err != nil {
		return nil, err
	}
	return &pb.Item{
		Item_ID:         uint32(item.ID),
		Item_Name:       item.ItemName,
		Length:          uint32(item.Length),
		Width:           uint32(item.Width),
		Fixed_Size_ID:   uint32(item.FixedSizeID),
		Is_Custom:       item.IsCustom,
		Estimated_Price: item.EstPrice,
		User_ID:         uint32(item.UserID),
	}, nil

}

// RemoveItemService implements interfaces.MaterialServiceInte.
func (m *MaterialService) RemoveItemService(p *pb.ItemID) (*pb.ItemResponse, error) {
	//Find by ItemsID
	item, err := m.Repo.FindItemByID(uint(p.ID))
	if err != nil {
		return &pb.ItemResponse{
			Status:  pb.ItemResponse_ERROR,
			Message: "Error in finding item",
			Payload: &pb.ItemResponse_Error{Error: err.Error()},
		}, err
	}

	//remove the items
	err = m.Repo.DeletItem(item.ID)
	if err != nil {
		return &pb.ItemResponse{
			Status:  pb.ItemResponse_ERROR,
			Message: "Error in removing item",
			Payload: &pb.ItemResponse_Error{Error: err.Error()},
		}, err
	}
	//Return success response after the successful removal
	return &pb.ItemResponse{
		Status:  pb.ItemResponse_OK,
		Message: "Item removes successfully",
	}, nil
}

// FindAllItemByUser implements interfaces.MaterialServiceInte.
func (m *MaterialService) FindAllItemByUser(p *pb.ItemID) (*pb.ItemList, error) {
	userID := p.ID
	item, err := m.Repo.FindAllItemByUsers(uint(userID))
	if err != nil {
		return nil, err
	}

	itemList := &pb.ItemList{
		Items: make([]*pb.Item, 0, len(*item)),
	}
	for _, itm := range *item {
		itemList.Items = append(itemList.Items, &pb.Item{
			Item_ID:       uint32(itm.ID),
			Item_Name:     itm.ItemName,
			Material_ID:   uint32(itm.MaterialID),
			Length:        uint32(itm.Length),
			Width:         uint32(itm.Width),
			Fixed_Size_ID: uint32(itm.FixedSizeID),
			//Is_Custom:       itm.IsCustom,
			Estimated_Price: itm.EstPrice,
			User_ID:         uint32(itm.UserID),
		})
	}
	return itemList, nil
}
