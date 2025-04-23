#ifndef COLLIDER_H
#define COLLIDER_H

#pragma once

#include <BaseComponent.h>
#include <memory>

#include "../../cmake-build-debug/_deps/bullet-src/src/BulletCollision/CollisionShapes/btCollisionShape.h"

class Collider : public BaseComponent {
public:
    explicit Collider(btCollisionShape* shape, const std::shared_ptr<GameObject> &gameObject)
        : BaseComponent(gameObject), _collisionShape(shape) {}

    btCollisionShape* collisionShape() const { return _collisionShape; }

    void serialize(std::ostream& os) const override {
        // Serialize Collider-specific data (e.g., collision shape type and parameters)
        // Placeholder: You need to implement serialization for the collision shape
    }

    void deserialize(std::istream& is) override {
        // Deserialize Collider-specific data
        // Placeholder: You need to implement deserialization for the collision shape
    }

private:
    btCollisionShape* _collisionShape;
};

#endif // COLLIDER_H