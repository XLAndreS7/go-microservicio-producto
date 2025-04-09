package services

import (
    "errors"
    "strings"

    "go-microservicio-producto/models"
    "go-microservicio-producto/repository"

    "go.mongodb.org/mongo-driver/mongo"
)

// Interfaz esperada por los controladores
type IProductoService interface {
    CrearProducto(producto models.Producto) error
    ObtenerProductos() ([]models.Producto, error)
    ActualizarProducto(id string, producto models.Producto) error
    EliminarProducto(id string) error
    ObtenerProductoPorID(id string) (*models.Producto, error)
}

type ProductoService struct {
    repo repository.ProductoRepository
}

func NewProductoService(repo repository.ProductoRepository) *ProductoService {
    return &ProductoService{repo: repo}
}

func (s *ProductoService) CrearProducto(producto models.Producto) error {
    return s.repo.CrearProducto(producto)
}

func (s *ProductoService) ObtenerProductos() ([]models.Producto, error) {
    return s.repo.ObtenerProductos()
}

func (s *ProductoService) ObtenerProductoPorID(id string) (*models.Producto, error) {
    return s.repo.ObtenerProductoPorID(id)
}

func (s *ProductoService) EliminarProducto(id string) error {
    return s.repo.EliminarProducto(id)
}

func (s *ProductoService) ActualizarProducto(id string, producto models.Producto) error {
    if strings.TrimSpace(id) == "" {
        return errors.New("el ID no puede estar vacío")
    }

    // Validación: al menos un campo debe venir con datos
    if producto.Nombre == "" && producto.Descripcion == "" && producto.Precio == 0 {
        return errors.New("no hay campos para actualizar")
    }

    err := s.repo.ActualizarProducto(id, producto)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return errors.New("producto no encontrado")
        }
        return err
    }
    return nil
}

