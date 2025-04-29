//
// Created by Tarook on 18/03/2025.
//

#include "Components/Rigidbody.h"

#include <iostream>

#include "Magnum/BulletIntegration/MotionState.h"

Rigidbody::Rigidbody(const Float mass, btCollisionShape *bShape, const std::shared_ptr<GameObject> &gameObject)
: BaseComponent(gameObject), _bWorld(gameObject->getWorld()->getBulletWorld()) {
    /* Calculate inertia so the object reacts as it should with rotation and everything */
    btVector3 bInertia(0.0f, 0.0f, 0.0f);
    if (!Math::TypeTraits<Float>::equals(mass, 0.0f))
        bShape->calculateLocalInertia(mass, bInertia);

    /* Bullet rigid body setup */
    auto *motionState = new BulletIntegration::MotionState{*this->gameObject};
    _bRigidBody.emplace(btRigidBody::btRigidBodyConstructionInfo{
        mass, &motionState->btMotionState(), bShape, bInertia
    });
    _bRigidBody->forceActivationState(DISABLE_DEACTIVATION);
    _bWorld.addRigidBody(_bRigidBody.get());
}

Rigidbody::Rigidbody(const Float mass, const std::shared_ptr<GameObject>& gameObject)
    : BaseComponent(gameObject), _bWorld(gameObject->getWorld()->getBulletWorld()) {
    // Try to get the Collider component from the GameObject
    _collider = gameObject->getComponent<Collider>();
    if (!_collider) {
        throw std::runtime_error("Rigidbody requires a Collider component.");
    }

    btCollisionShape* bShape = _collider->collisionShape();
    if (!bShape) {
        throw std::runtime_error("Collider does not have a valid collision shape.");
    }

    // Create the rigid body
    btVector3 localInertia(0, 0, 0);
    if (mass != 0.0f) {
        bShape->calculateLocalInertia(mass, localInertia);
    }

    /* Bullet rigid body setup */
    auto *motionState = new BulletIntegration::MotionState{*this->gameObject};
    _bRigidBody.emplace(btRigidBody::btRigidBodyConstructionInfo{
        mass, &motionState->btMotionState(), bShape, localInertia
    });
    _bRigidBody->forceActivationState(DISABLE_DEACTIVATION);
    _bWorld.addRigidBody(_bRigidBody.get());
}

Rigidbody::~Rigidbody() {
    _bWorld.removeRigidBody(_bRigidBody.get());
}
