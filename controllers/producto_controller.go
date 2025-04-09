package controllers

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "go-microservicio-producto/models"
    "go-microservicio-producto/services"
)

type ProductoController struct {
    service services.IProductoService
}

// Nuevo constructor que usa la interfaz y no la implementación concreta
func NewProductoController(service services.IProductoService) *ProductoController {
    return &ProductoController{
        service: service,
    }
}

// CrearProducto maneja la solicitud POST para crear un producto
func (c *ProductoController) CrearProducto(w http.ResponseWriter, r *http.Request) {
    var producto models.Producto
    err := json.NewDecoder(r.Body).Decode(&producto)
    if err != nil {
        http.Error(w, "Error al decodificar el producto", http.StatusBadRequest)
        return
    }

    err = c.service.CrearProducto(producto)
    if err != nil {
        http.Error(w, "Error al crear el producto", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"mensaje": "Producto creado correctamente"})
}

// ObtenerProductos maneja la solicitud GET para listar productos
func (c *ProductoController) ObtenerProductos(w http.ResponseWriter, r *http.Request) {
    productos, err := c.service.ObtenerProductos()
    if err != nil {
        http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(productos)
}

// ObtenerProductoPorID maneja la solicitud GET para un producto por ID
func (c *ProductoController) ObtenerProductoPorID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if id == "" {
        http.Error(w, "ID del producto requerido", http.StatusBadRequest)
        return
    }

    producto, err := c.service.ObtenerProductoPorID(id)
    if err != nil {
        http.Error(w, "Producto no encontrado", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(producto)
}

// ActualizarProducto maneja la solicitud PUT para actualizar un producto
func (c *ProductoController) ActualizarProducto(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var producto models.Producto
    if err := json.NewDecoder(r.Body).Decode(&producto); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }

    err := c.service.ActualizarProducto(id, producto)
    if err != nil {
        switch err.Error() {
        case "el ID no puede estar vacío", "no hay campos para actualizar":
            http.Error(w, err.Error(), http.StatusBadRequest)
        case "producto no encontrado":
            http.Error(w, err.Error(), http.StatusNotFound)
        default:
            http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        }
        return
    }

    productoActualizado, _ := c.service.ObtenerProductoPorID(id)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "mensaje":  "Producto actualizado correctamente",
        "producto": productoActualizado,
    })
}

// EliminarProducto maneja la solicitud DELETE para eliminar un producto
func (c *ProductoController) EliminarProducto(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if id == "" {
        http.Error(w, "ID del producto requerido", http.StatusBadRequest)
        return
    }

    err := c.service.EliminarProducto(id)
    if err != nil {
        http.Error(w, "Error al eliminar el producto", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"mensaje": "Producto eliminado correctamente"})
}
