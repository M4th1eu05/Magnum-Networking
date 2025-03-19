#include <BaseComponent.h>
#include <GameObject.h>
#include "BulletCollision/CollisionDispatch/btCollisionObject.h"
#include "BulletCollision/CollisionShapes/btCollisionShape.h"
#include "BulletDynamics/Dynamics/btDynamicsWorld.h"
#include "BulletDynamics/Dynamics/btRigidBody.h"
#include "Corrade/Containers/Pointer.h"
#include "LinearMath/btVector3.h"
#include "Magnum/Magnum.h"
#include "Magnum/BulletIntegration/MotionState.h"

class btDynamicsWorld;
class btCollisionShape;
class btRigidBody;

class Rigidbody : public BaseComponent {
public:
    Rigidbody(const Float mass, btCollisionShape *bShape, btDynamicsWorld &bWorld, std::shared_ptr<GameObject> gameObject)
        : BaseComponent(std::move(gameObject)), _bWorld(bWorld) {
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
        bWorld.addRigidBody(_bRigidBody.get());
    }

    ~Rigidbody() {
        _bWorld.removeRigidBody(_bRigidBody.get());
    }

    btRigidBody &rigidBody() { return *_bRigidBody; }

    /* needed after changing the pose from Magnum side */
    void syncPose() {
        const auto& matrix = gameObject->transformationMatrix();
        btTransform transform;
        transform.setFromOpenGLMatrix(matrix.data());
        _bRigidBody->setWorldTransform(transform);
    }

private:
    btDynamicsWorld &_bWorld;
    Containers::Pointer<btRigidBody> _bRigidBody;
};