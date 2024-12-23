package service

import (
	"testing"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/model"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMaterialRepo struct {
	mock.Mock
}

func (m *MockMaterialRepo) AddMaterial(Material *model.Material) (uint, error) {
	args := m.Called(Material)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockMaterialRepo) FindMaterialByID(MaterialID uint) (*model.Material, error) {
	args := m.Called(MaterialID)
	return args.Get(0).(*model.Material), args.Error(1)
}

func (m *MockMaterialRepo) FindAllMaterial() (*[]model.Material, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Material), args.Error(1)
}

func (m *MockMaterialRepo) UpdateMaterial(Material *model.Material) error {
	args := m.Called(Material)
	return args.Error(0)
}

func (m *MockMaterialRepo) DeleteMaterial(MaterialID uint) error {
	args := m.Called(MaterialID)
	return args.Error(0)
}

func (m *MockMaterialRepo) UpdateMaterialStock(materialID uint, quantity uint) error {
	args := m.Called(materialID, quantity)
	return args.Error(0)
}

// Mocking the additional methods from the interface
func (m *MockMaterialRepo) CreateItem(Item *model.Item) (uint, error) {
	args := m.Called(Item)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockMaterialRepo) FindItemByID(ItemsID uint) (*model.Item, error) {
	args := m.Called(ItemsID)
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockMaterialRepo) FindAllItem() (*[]model.Item, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockMaterialRepo) UpdateItem(Item *model.Item) error {
	args := m.Called(Item)
	return args.Error(0)
}

func (m *MockMaterialRepo) DeletItem(ItemID uint) error {
	args := m.Called(ItemID)
	return args.Error(0)
}

func (m *MockMaterialRepo) FindAllItemByUsers(userID uint) (*[]model.Item, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockMaterialRepo) CreateOrders(Orders *model.Order) (uint, error) {
	args := m.Called(Orders)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockMaterialRepo) FindOrdersByID(OrdersID uint) (*model.Order, error) {
	args := m.Called(OrdersID)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockMaterialRepo) UpdateOrders(Orders *model.Order) error {
	args := m.Called(Orders)
	return args.Error(0)
}

func (m *MockMaterialRepo) UpdateOrderStaus(OrderID uint, status string) error {
	args := m.Called(OrderID, status)
	return args.Error(0)
}

func (m *MockMaterialRepo) DeleteOrders(OrdersID uint) error {
	args := m.Called(OrdersID)
	return args.Error(0)
}

func (m *MockMaterialRepo) FindAllOrders() (*[]model.Order, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Order), args.Error(1)
}

func (m *MockMaterialRepo) FindOrdersByUser(userID uint) (*[]model.Order, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]model.Order), args.Error(1)
}

func (m *MockMaterialRepo) FindOrder(userID, ItemID uint) (*model.Order, error) {
	args := m.Called(userID, ItemID)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockMaterialRepo) GetLatestPaymentByOrderID(orderID int) (model.Payment, error) {
	args := m.Called(orderID)
	return args.Get(0).(model.Payment), args.Error(1)
}

func (m *MockMaterialRepo) SavePayment(payment *model.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockMaterialRepo) UpdatePaymentStatus(paymentID, status string) error {
	args := m.Called(paymentID, status)
	return args.Error(0)
}

func (m *MockMaterialRepo) UpdatePaymentAndOrderStatus(paymentID string, orderID int, paymentStatus string, orderStatus string) error {
	args := m.Called(paymentID, orderID, paymentStatus, orderStatus)
	return args.Error(0)
}

func (m *MockMaterialRepo) SaveCuttingResult(itemID uint, components []model.Component) error {
	args := m.Called(itemID, components)
	return args.Error(0)
}

func (m *MockMaterialRepo) GetCuttingResultByItemID(itemID uint) ([]model.Component, error) {
	args := m.Called(itemID)
	return args.Get(0).([]model.Component), args.Error(1)
}

func TestAddMaterialService(t *testing.T) {
	mockRepo := new(MockMaterialRepo)
	materialService := &MaterialService{
		Repo: mockRepo,
	}

	// Mock the repository response
	mockRepo.On("AddMaterial", mock.Anything).Return(uint(1), nil)

	// Create test data
	material := &pb.Material{
		Material_Name: "Test Material",
		Description:   "Description of material",
		Stock:         10,
		Price:         99.99,
	}

	// Call the service method
	resp, err := materialService.AddMaterialService(material)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Material added successfully", resp.Message)
	assert.Equal(t, pb.MaterialResponse_OK, resp.Status)
	mockRepo.AssertExpectations(t)
}
func TestEditMaterialService(t *testing.T) {
	mockRepo := new(MockMaterialRepo)
	materialService := &MaterialService{
		Repo: mockRepo,
	}

	// Mock the repository response
	mockRepo.On("FindMaterialByID", uint(1)).Return(&model.Material{Name: "Old Material", Description: "Old Description", Stock: 10, Price: 50.0}, nil)
	mockRepo.On("UpdateMaterial", mock.Anything).Return(nil)

	// Create test data
	material := &pb.Material{
		Material_ID:   1,
		Material_Name: "Updated Material",
		Description:   "Updated Description",
		Stock:         20,
		Price:         75.0,
	}

	// Call the service method
	resp, err := materialService.EditMaterialService(material)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), resp.Material_ID)
	assert.Equal(t, "Updated Material", resp.Material_Name)
	mockRepo.AssertExpectations(t)
}
func TestFindMaterialByIDService(t *testing.T) {
	mockRepo := new(MockMaterialRepo)
	materialService := &MaterialService{
		Repo: mockRepo,
	}

	// Mock the repository response
	mockRepo.On("FindMaterialByID", uint(1)).Return(&model.Material{
		Name:        "Test Material",
		Description: "Test Description",
		Stock:       10,
		Price:       100.0,
	}, nil)

	// Create test data
	materialID := &pb.MaterialID{ID: 1}

	// Call the service method
	material, err := materialService.FindMaterialByIDService(materialID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Test Material", material.Material_Name)
	assert.Equal(t, "Test Description", material.Description)
	mockRepo.AssertExpectations(t)
}
