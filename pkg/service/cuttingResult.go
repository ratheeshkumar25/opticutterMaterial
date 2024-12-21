package service

import (
	"log"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/utils"
)

// PredefinedSizes holds predefined sizes (fixed size ID => size)
var PredefinedSizes = map[uint]model.PredefinedSize{
	1: {Length: 50, Width: 30, Name: "Small"},
	2: {Length: 100, Width: 60, Name: "Medium"},
	3: {Length: 150, Width: 90, Name: "Large"},
}

// GenerateCuttingResult generates the cutting result for a given item ID
func (m *MaterialService) GenerateCuttingResult(p *pb.ItemID) (*pb.CuttingResultResponse, error) {
	log.Printf("Triggering cutting result generation for item ID: %d", p.ID)

	// Retrieve the item details from the database using the Repo
	item, err := m.Repo.FindItemByID(uint(p.ID))
	if err != nil {
		// Log the error with item ID for debugging
		log.Printf("Error finding item with ID %d: %v", p.ID, err)
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "Item not found: " + err.Error(),
		}, err
	}

	// Log the found item to ensure we have the correct data
	log.Printf("Found item: %+v", item)

	var length, width uint

	// Check if the item uses a predefined size
	if item.FixedSizeID != 0 {
		// Ensure that we have a valid fixed size for the item
		predefinedSize, exists := utils.PredefinedSizes[item.FixedSizeID]
		if !exists {
			log.Printf("Invalid fixedSizeID %d for item ID %d", item.FixedSizeID, p.ID)
			return &pb.CuttingResultResponse{
				Status:  pb.CuttingResultResponse_ERROR,
				Message: "invalid fixedSizeID",
			}, nil
		}
		length = predefinedSize.Length
		width = predefinedSize.Width
	} else {
		// Use the item's specified length and width
		length = item.Length
		width = item.Width
	}

	// Log the size being used for cutting result generation
	log.Printf("Generating components for item ID %d with Length: %d, Width: %d", item.ID, length, width)

	// Generate the components using the utility function
	components, err := utils.GenerateComponents(item.MaterialID, item.ItemName, length, width)
	if err != nil {
		log.Printf("Failed to generate components for item ID %d: %v", item.ID, err)
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "Failed to generate components: " + err.Error(),
		}, err
	}

	// Log the components generated
	log.Printf("Generated components: %+v", components)

	// Save the cutting result to the database
	err = m.Repo.SaveCuttingResult(item.ID, components)
	if err != nil {
		log.Printf("Error saving cutting result for item ID %d: %v", item.ID, err)
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "Failed to save cutting result: " + err.Error(),
		}, err
	}

	// Convert the components to the protobuf representation
	var componentProtos []*pb.Component
	for _, component := range components {
		componentProtos = append(componentProtos, &pb.Component{
			Material_ID:   uint32(component.MaterialID),
			DoorPanel:     component.DoorPanel,
			BackSidePanel: component.BackSidePanel,
			SidePanel:     component.SidePanel,
			TopPanel:      component.TopPanel,
			BottomPanel:   component.BottomPanel,
			ShelvesPanel:  component.ShelvesPanel,
			Panel_Count:   component.PanelCount,
		})
	}

	var componentPayloads []utils.ComponentPayload
	for _, component := range components {
		componentPayloads = append(componentPayloads, utils.ComponentPayload{
			MaterialID:    component.MaterialID,
			DoorPanel:     component.DoorPanel,
			BackSidePanel: component.BackSidePanel,
			SidePanel:     component.SidePanel,
			TopPanel:      component.TopPanel,
			BottomPanel:   component.BottomPanel,
			ShelvesPanel:  component.ShelvesPanel,
			PanelCount:    component.PanelCount,
		})
	}

	// Call the function to handle the cutting result notification
	err = utils.HandleCuttingResultNotification(item.ID, componentPayloads, item.ID)
	if err != nil {
		log.Printf("Error notifying cutting result for item ID %d: %v", item.ID, err)
	}

	// Return the successful response with the generated cutting result
	return &pb.CuttingResultResponse{
		Status:  pb.CuttingResultResponse_OK,
		Message: "Cutting result generated successfully",
		CuttingResult: &pb.CuttingResult{
			Item_ID:    uint32(item.ID),
			Components: componentProtos,
		},
	}, nil
}

// GetCuttingResService retrieves the cutting result for an item
func (m *MaterialService) GetCuttingResService(p *pb.ItemID) (*pb.CuttingResultResponse, error) {
	// Retrieve the item details from the database using the Repo
	log.Println("itemid", p.ID)
	item, err := m.Repo.FindItemByID(uint(p.ID))
	if err != nil {
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "item not found: " + err.Error(),
		}, err
	}

	// Retrieve the cutting result components from the database
	components, err := m.Repo.GetCuttingResultByItemID(item.ID)
	if err != nil {
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "failed to retrieve cutting result: " + err.Error(),
		}, err
	}

	// Convert the components into the Protobuf representation
	var componentProtos []*pb.Component
	for _, component := range components {
		componentProtos = append(componentProtos, &pb.Component{
			Material_ID:   uint32(component.MaterialID),
			DoorPanel:     component.DoorPanel,
			BackSidePanel: component.BackSidePanel,
			SidePanel:     component.SidePanel,
			TopPanel:      component.TopPanel,
			BottomPanel:   component.BottomPanel,
			ShelvesPanel:  component.ShelvesPanel,
			Panel_Count:   component.PanelCount,
		})
	}

	// Return the successful response
	return &pb.CuttingResultResponse{
		Status:  pb.CuttingResultResponse_OK,
		Message: "cutting result retrieved successfully",
		CuttingResult: &pb.CuttingResult{
			Item_ID:    uint32(item.ID),
			Components: componentProtos,
		},
	}, nil
}

// func (m *MaterialService) GenerateCuttingResult(p *pb.ItemID) (*pb.CuttingResultResponse, error) {
// 	log.Println("itemmms id", p.ID)
// 	// Retrieve the item details from the database using the Repo
// 	item, err := m.Repo.FindItemByID(uint(p.ID))
// 	if err != nil {
// 		return &pb.CuttingResultResponse{
// 			Status:  pb.CuttingResultResponse_ERROR,
// 			Message: "item not found: " + err.Error(),
// 		}, err
// 	}

// 	log.Println("find the itemID", item)

// 	var length, width uint

// 	// Check if the item uses a predefined size
// 	if item.FixedSizeID != 0 {
// 		predefinedSize, exists := PredefinedSizes[item.FixedSizeID]
// 		if !exists {
// 			return &pb.CuttingResultResponse{
// 				Status:  pb.CuttingResultResponse_ERROR,
// 				Message: "invalid fixedSizeID",
// 			}, nil
// 		}
// 		length = predefinedSize.Length
// 		width = predefinedSize.Width
// 	} else {
// 		length = item.Length
// 		width = item.Width
// 	}

// 	// Use GenerateComponents to get the components
// 	components, err := utils.GenerateComponents(item.MaterialID, length, width)
// 	if err != nil {
// 		return &pb.CuttingResultResponse{
// 			Status:  pb.CuttingResultResponse_ERROR,
// 			Message: "failed to generate components: " + err.Error(),
// 		}, err
// 	}

// 	err = m.Repo.SaveCuttingResult(item.ID, components)
// 	if err != nil {
// 		log.Printf("Error saving cutting result for item %d: %v", item.ID, err)
// 		return &pb.CuttingResultResponse{
// 			Status:  pb.CuttingResultResponse_ERROR,
// 			Message: "failed to save cutting result: " + err.Error(),
// 		}, err
// 	}

// 	// log.Println("itemsID", item.ID)
// 	log.Printf("Item retrieved: %+v", item)
// 	log.Printf("Generated components: %+v", components)
// 	log.Printf("Saving components for ItemID %d", item.ID)
// 	log.Printf("Retrieved components from DB: %+v", components)

// 	// Convert components to the protobuf representation
// 	var componentProtos []*pb.Component
// 	for _, component := range components {
// 		componentProtos = append(componentProtos, &pb.Component{
// 			Material_ID:   uint32(component.MaterialID),
// 			DoorPanel:     component.DoorPanel,
// 			BackSidePanel: component.BackSidePanel,
// 			SidePanel:     component.SidePanel,
// 			TopPanel:      component.TopPanel,
// 			BottomPanel:   component.BottomPanel,
// 			ShelvesPanel:  component.ShelvesPanel,
// 			Panel_Count:   component.PanelCount,
// 		})
// 	}

// 	// Return the successful response
// 	return &pb.CuttingResultResponse{
// 		Status:  pb.CuttingResultResponse_OK,
// 		Message: "cutting result generated successfully",
// 		CuttingResult: &pb.CuttingResult{
// 			Item_ID:    uint32(item.ID),
// 			Components: componentProtos,
// 		},
// 	}, nil
// }
