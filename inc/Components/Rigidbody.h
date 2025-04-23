#pragma once

#include <BaseComponent.h>
#include <GameObject.h>
#include <World.h>

#include "Collider.h"
#include "BulletCollision/CollisionDispatch/btCollisionObject.h"
#include "BulletCollision/CollisionShapes/btCollisionShape.h"
#include "BulletDynamics/Dynamics/btDynamicsWorld.h"
#include "BulletDynamics/Dynamics/btRigidBody.h"
#include "Corrade/Containers/Pointer.h"

class btDynamicsWorld;
class btCollisionShape;
class btRigidBody;

class Rigidbody : public BaseComponent {
public:
    Rigidbody(const Float mass, btCollisionShape *bShape, const std::shared_ptr<GameObject> &gameObject);
    Rigidbody(const Float mass, const std::shared_ptr<GameObject>& gameObject);

    ~Rigidbody();

    btRigidBody &rigidBody() { return *_bRigidBody; }

    /* needed after changing the pose from Magnum side */
    void syncPose() {
        _bRigidBody->setWorldTransform(btTransform(gameObject->transformationMatrix()));
    }

private:
    btDynamicsWorld &_bWorld;
    Containers::Pointer<btRigidBody> _bRigidBody;
    std::shared_ptr<Collider> _collider;
};