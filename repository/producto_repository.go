package repository

import (
    "context"
    "errors"
    "go-microservicio-producto/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type ProductoRepository interface {
    CrearProducto(producto models.Producto) error
    ObtenerProductos() ([]models.Producto, error)
    ObtenerProductoPorID(id string) (*models.Producto, error)
    ActualizarProducto(id string, producto models.Producto) error
    EliminarProducto(id string) error
}

type productoRepositoryImpl struct {
    collection *mongo.Collection
}

func NewProductoRepository(db *mongo.Database) ProductoRepository {
    return &productoRepositoryImpl{
        collection: db.Collection("productos"),
    }
}

func (r *productoRepositoryImpl) CrearProducto(producto models.Producto) error {
    _, err := r.collection.InsertOne(context.TODO(), producto)
    return err
}

func (r *productoRepositoryImpl) ObtenerProductos() ([]models.Producto, error) {
    var productos []models.Producto
    cursor, err := r.collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var producto models.Producto
        if err := cursor.Decode(&producto); err != nil {
            return nil, err
        }
        productos = append(productos, producto)
    }

    return productos, nil
}

func (r *productoRepositoryImpl) ObtenerProductoPorID(id string) (*models.Producto, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var producto models.Producto
    err = r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&producto)
    if err != nil {
        return nil, err
    }

    return &producto, nil
}

func (r *productoRepositoryImpl) ActualizarProducto(id string, producto models.Producto) error {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    update := bson.M{}

    if producto.Nombre != "" {
        update["nombre"] = producto.Nombre
    }
    if producto.Descripcion != "" {
        update["descripcion"] = producto.Descripcion
    }
    if producto.Precio != 0 {
        update["precio"] = producto.Precio
    }

    if len(update) == 0 {
        return errors.New("no hay campos válidos para actualizar")
    }

    result, err := r.collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": objectId},
        bson.M{"$set": update},
    )
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return mongo.ErrNoDocuments
    }

    return nil
}

func (r *productoRepositoryImpl) EliminarProducto(id string) error {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
    return err
}

