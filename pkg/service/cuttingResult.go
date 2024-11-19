package service

import (
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

// GenerateCuttingResult generates the cutting result for an item
func (m *MaterialService) GenerateCuttingResult(p *pb.ItemID) (*pb.CuttingResultResponse, error) {
	// Retrieve the item details from the database using the Repo
	item, err := m.Repo.FindItemByID(uint(p.ID))
	if err != nil {
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "item not found: " + err.Error(),
		}, err
	}

	var length, width uint

	// Check if the item uses a predefined size
	if item.FixedSizeID != 0 {
		predefinedSize, exists := PredefinedSizes[item.FixedSizeID]
		if !exists {
			return &pb.CuttingResultResponse{
				Status:  pb.CuttingResultResponse_ERROR,
				Message: "invalid fixedSizeID",
			}, nil
		}
		length = predefinedSize.Length
		width = predefinedSize.Width
	} else {
		length = item.Length
		width = item.Width
	}

	// Use GenerateComponents to get the components
	components, err := utils.GenerateComponents(item.MaterialID, length, width)
	if err != nil {
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "failed to generate components: " + err.Error(),
		}, err
	}

	// Save the CuttingResult using the Repo layer
	err = m.Repo.SaveCuttingResult(item.ID, components)
	if err != nil {
		return &pb.CuttingResultResponse{
			Status:  pb.CuttingResultResponse_ERROR,
			Message: "failed to save cutting result: " + err.Error(),
		}, err
	}

	// Convert components to the protobuf representation
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
		Message: "cutting result generated successfully",
		CuttingResult: &pb.CuttingResult{
			Item_ID:    uint32(item.ID),
			Components: componentProtos,
		},
	}, nil
}

// GetCuttingResService retrieves the cutting result for an item
func (m *MaterialService) GetCuttingResService(p *pb.ItemID) (*pb.CuttingResultResponse, error) {
	// Retrieve the item details from the database using the Repo
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
