//
// Created by Tarook on 18/03/2025.
//

#ifndef WORLD_H
#define WORLD_H
#include "GameObject.h"
#include "BulletCollision/BroadphaseCollision/btDbvtBroadphase.h"
#include "BulletCollision/CollisionDispatch/btDefaultCollisionConfiguration.h"
#include "BulletDynamics/ConstraintSolver/btSequentialImpulseConstraintSolver.h"
#include "BulletDynamics/Dynamics/btDiscreteDynamicsWorld.h"
#include "Magnum/BulletIntegration/DebugDraw.h"
#include "Magnum/SceneGraph/MatrixTransformation3D.h"
#include "Magnum/SceneGraph/Scene.h"

using namespace Magnum;
typedef SceneGraph::Scene<SceneGraph::MatrixTransformation3D> Scene3D;

class World {
public:
    World(Magnum::Timeline& timeline);

    Scene3D& getScene() { return _scene; }

    std::shared_ptr<GameObject> createGameObject(std::shared_ptr<GameObject> parent = nullptr);
    void addObject(const std::shared_ptr<GameObject> &object);

    void removeObject(std::shared_ptr<GameObject> object);

    void update();

    btDiscreteDynamicsWorld &getBulletWorld() { return _bWorld; }

public:
    BulletIntegration::DebugDraw _debugDraw{NoCreate};

private:
    Timeline& _timeline;

    Scene3D _scene;
    std::vector<std::shared_ptr<GameObject>> _objects;

    /* PHYSICS */
    btDbvtBroadphase _bBroadphase;
    btDefaultCollisionConfiguration _bCollisionConfig;
    btCollisionDispatcher _bDispatcher{&_bCollisionConfig};
    btSequentialImpulseConstraintSolver _bSolver;

    /* The world has to live longer than the scene because RigidBody
           instances have to remove themselves from it on destruction */
    btDiscreteDynamicsWorld _bWorld{&_bDispatcher, &_bBroadphase, &_bSolver, &_bCollisionConfig};

    float physicsTimeStep = 1.0f / 60.0f;
};


#endif //WORLD_H
