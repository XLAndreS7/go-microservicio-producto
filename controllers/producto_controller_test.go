package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-microservicio-producto/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock del servicio
type MockProductoService struct {
	mock.Mock
}

func (m *MockProductoService) CrearProducto(producto models.Producto) error {
	args := m.Called(producto)
	return args.Error(0)
}

func (m *MockProductoService) ObtenerProductos() ([]models.Producto, error) {
	args := m.Called()
	return args.Get(0).([]models.Producto), args.Error(1)
}

func (m *MockProductoService) ObtenerProductoPorID(id string) (*models.Producto, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Producto), args.Error(1)
}

func (m *MockProductoService) ActualizarProducto(id string, producto models.Producto) error {
	args := m.Called(id, producto)
	return args.Error(0)
}

func (m *MockProductoService) EliminarProducto(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestObtenerProductos(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	productoID := primitive.NewObjectID()
	productosEsperados := []models.Producto{
		{ID: productoID, Nombre: "Producto A", Precio: 100},
	}

	mockService.On("ObtenerProductos").Return(productosEsperados, nil)

	req, _ := http.NewRequest("GET", "/productos", nil)
	rr := httptest.NewRecorder()

	controller.ObtenerProductos(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var productos []models.Producto
	err := json.NewDecoder(rr.Body).Decode(&productos)
	assert.NoError(t, err)
	assert.Equal(t, productosEsperados, productos)
}

func TestObtenerProductos_Error(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	mockService.On("ObtenerProductos").Return([]models.Producto{}, errors.New("fallo"))

	req, _ := http.NewRequest("GET", "/productos", nil)
	rr := httptest.NewRecorder()

	controller.ObtenerProductos(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCrearProducto(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	productoID := primitive.NewObjectID()
	nuevoProducto := models.Producto{ID: productoID, Nombre: "Producto B", Precio: 200}
	body, _ := json.Marshal(nuevoProducto)

	mockService.On("CrearProducto", nuevoProducto).Return(nil)

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	controller.CrearProducto(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp map[string]string
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Producto creado correctamente", resp["mensaje"])
}

func TestCrearProducto_Error(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	productoID := primitive.NewObjectID()
	nuevoProducto := models.Producto{ID: productoID, Nombre: "Producto C", Precio: 300}
	body, _ := json.Marshal(nuevoProducto)

	mockService.On("CrearProducto", nuevoProducto).Return(errors.New("error"))

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	controller.CrearProducto(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCrearProducto_BodyInvalido(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	body := []byte(`{invalid-json}`)

	req, _ := http.NewRequest("POST", "/productos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	controller.CrearProducto(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestObtenerProductoPorID(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	id := primitive.NewObjectID().Hex()
	producto := &models.Producto{ID: primitive.NewObjectID(), Nombre: "Producto X", Precio: 55}

	mockService.On("ObtenerProductoPorID", id).Return(producto, nil)

	req, _ := http.NewRequest("GET", "/productos/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})

	rr := httptest.NewRecorder()
	controller.ObtenerProductoPorID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp models.Producto
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, producto.Nombre, resp.Nombre)
}

func TestObtenerProductoPorID_Error(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	id := primitive.NewObjectID().Hex()
	mockService.On("ObtenerProductoPorID", id).Return(nil, errors.New("no encontrado"))

	req, _ := http.NewRequest("GET", "/productos/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})

	rr := httptest.NewRecorder()
	controller.ObtenerProductoPorID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestActualizarProducto(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	producto := models.Producto{Nombre: "Actualizado", Precio: 100}
	jsonBody, _ := json.Marshal(producto)

	req, err := http.NewRequest(http.MethodPut, "/productos/actualizar?id=123", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	rr := httptest.NewRecorder()

	mockService.On("ActualizarProducto", "123", producto).Return(nil)
	mockService.On("ObtenerProductoPorID", "123").Return(&producto, nil)

	controller.ActualizarProducto(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestActualizarProducto_SinID(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	producto := models.Producto{Nombre: "Sin ID", Precio: 150}
	jsonBody, _ := json.Marshal(producto)

	req, _ := http.NewRequest(http.MethodPut, "/productos/actualizar", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	mockService.On("ActualizarProducto", "", producto).Return(errors.New("el ID no puede estar vacío"))

	controller.ActualizarProducto(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertExpectations(t)
}

func TestActualizarProducto_ErrorServicio(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	producto := models.Producto{Nombre: "Error", Precio: 150}
	jsonBody, _ := json.Marshal(producto)

	req, _ := http.NewRequest(http.MethodPut, "/productos/actualizar?id=123", bytes.NewBuffer(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	rr := httptest.NewRecorder()

	mockService.On("ActualizarProducto", "123", producto).Return(errors.New("fallo"))

	controller.ActualizarProducto(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestEliminarProducto(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	req, err := http.NewRequest(http.MethodDelete, "/productos/eliminar?id=123", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	rr := httptest.NewRecorder()

	mockService.On("EliminarProducto", "123").Return(nil)

	controller.EliminarProducto(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestEliminarProducto_SinID(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	req, _ := http.NewRequest(http.MethodDelete, "/productos/eliminar", nil)
	rr := httptest.NewRecorder()

	controller.EliminarProducto(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEliminarProducto_ErrorServicio(t *testing.T) {
	mockService := new(MockProductoService)
	controller := NewProductoController(mockService)

	req, _ := http.NewRequest(http.MethodDelete, "/productos/eliminar?id=123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	rr := httptest.NewRecorder()

	mockService.On("EliminarProducto", "123").Return(errors.New("fallo"))

	controller.EliminarProducto(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

