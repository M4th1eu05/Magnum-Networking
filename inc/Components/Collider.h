#ifndef COLLIDER_H
#define COLLIDER_H

#pragma once

#include <Components/BaseComponent.h>
#include <memory>

#include "BulletCollision/CollisionShapes/btCollisionShape.h"


class Collider : public BaseComponent {
public:
    explicit Collider(btCollisionShape* shape, const std::shared_ptr<GameObject> &gameObject)
        : BaseComponent(gameObject), _collisionShape(shape) {}

    btCollisionShape* collisionShape() const { return _collisionShape; } // this is a reference to the shape, there is no need to instantiate a shape every time

    void serialize(std::ostream& os) const override {
        // Serialize Collider-specific data
        // Placeholder: You need to implement serialization for the collision shape
        os.write(reinterpret_cast<const char*>(&_collisionShape), sizeof(_collisionShape));
    }

    void deserialize(std::istream& is) override {
        // Deserialize Collider-specific data
        // Placeholder: You need to implement deserialization for the collision shape
        is.read(reinterpret_cast<char*>(&_collisionShape), sizeof(_collisionShape));
    }

    void destroy() override {
        if (_collisionShape) {
            _collisionShape = nullptr;
        }
    }


private:
    btCollisionShape* _collisionShape;
};

#endif // COLLIDER_H