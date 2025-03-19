//
// Created by Tarook on 18/03/2025.
//

#include "Rigidbody.h"

#include <iostream>

Rigidbody::Rigidbody(const Float mass, btCollisionShape *bShape, const std::shared_ptr<World> &world, const std::shared_ptr<GameObject> &gameObject)
: BaseComponent(gameObject), _bWorld(world->getBulletWorld()) {
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
    std::cout << "Rigidbody created" << std::endl;
}

Rigidbody::~Rigidbody() {
    _bWorld.removeRigidBody(_bRigidBody.get());
}
