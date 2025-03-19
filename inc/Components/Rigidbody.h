#include <BaseComponent.h>
#include <GameObject.h>
#include <World.h>

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
    Rigidbody(const Float mass, btCollisionShape *bShape, const std::shared_ptr<World> &world, const std::shared_ptr<GameObject> &gameObject);

    ~Rigidbody();

    btRigidBody &rigidBody() { return *_bRigidBody; }

    /* needed after changing the pose from Magnum side */
    void syncPose() {
        _bRigidBody->setWorldTransform(btTransform(gameObject->transformationMatrix()));
    }

private:
    btDynamicsWorld &_bWorld;
    Containers::Pointer<btRigidBody> _bRigidBody;
};