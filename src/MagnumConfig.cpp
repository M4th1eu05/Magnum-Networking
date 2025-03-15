#include <Magnum/Platform/Sdl2Application.h>
#include <Magnum/Math/Vector3.h>
#include <BulletCollision/CollisionShapes/btBoxShape.h>
#include <BulletDynamics/Dynamics/btDiscreteDynamicsWorld.h>
#include <ImGui/ImGui.h>
#include <Magnum/ImGuiIntegration/Context.h>

using namespace Magnum;
using namespace Magnum::Math::Literals;

class Game : public Platform::Application {
public:
    explicit Game(const Arguments& arguments);

    void drawEvent() override;
    void updateEvent() override;

private:
    void setupPhysics();
    void setupScene();
    void updateCamera();

    btDiscreteDynamicsWorld* world;
    btRigidBody* floor;
    btRigidBody* sphere;

    ImGuiIntegration::Context imguiContext;
    Vector3 cameraPosition;
    Matrix4 projectionMatrix;
};

Game::Game(const Arguments& arguments) : Platform::Application{ arguments, NoCreate } {
    // Configuration initiale de Magnum et Bullet
}

void Game::setupPhysics() {
    // Initialisation de Bullet : terrain, sphères, etc.
}

void Game::setupScene() {
    // Configuration du sol : cubes disposés en grille de 10x10 (100 cubes)
}

void Game::updateEvent() {
    updateCamera();
    world->stepSimulation(1.f / 60.f);

    // Ajoutez l'animation ici, comme le mouvement de la caméra et la simulation physique
}

void Game::updateCamera() {
    // Logique pour déplacer la caméra avec WASD et QE
}

void Game::drawEvent() {
    // Dessiner la scène avec Magnum et ImGUI
}

void Game::setupScene() {
    // Création de la forme du sol, chaque cube dans la grille
    for (int x = 0; x < 10; ++x) {
        for (int z = 0; z < 10; ++z) {
            auto cubeTransform = Matrix4::translation(Vector3{ float(x), 0.f, float(z) });
            auto cubeShape = new btBoxShape(btVector3(0.5f, 0.5f, 0.5f));

            // Ajout du cube à Bullet
            btRigidBody::btRigidBodyConstructionInfo info(0, nullptr, cubeShape);
            btRigidBody* body = new btRigidBody(info);
            world->addRigidBody(body);
        }
    }
}

void Game::updateCamera() {
    if (keyPressed(KeyEvent::Key::W)) {
        cameraPosition.z() -= 0.1f;
    }
    if (keyPressed(KeyEvent::Key::S)) {
        cameraPosition.z() += 0.1f;
    }
    if (keyPressed(KeyEvent::Key::A)) {
        cameraPosition.x() -= 0.1f;
    }
    if (keyPressed(KeyEvent::Key::D)) {
        cameraPosition.x() += 0.1f;
    }
    if (keyPressed(KeyEvent::Key::Q)) {
        cameraPosition.y() += 0.1f;
    }
    if (keyPressed(KeyEvent::Key::E)) {
        cameraPosition.y() -= 0.1f;
    }

    // Appliquer la nouvelle position de la caméra
    projectionMatrix = Matrix4::perspectiveProjection(45.0_degf, windowSize().aspectRatio(), 0.1f, 100.f);
}

void Game::mousePressEvent(MouseEvent& event) {
    // Récupérer la direction de la caméra
    Vector3 rayDirection = cameraDirection;  // Calculer la direction de la caméra

    // Créer une sphère et la lancer
    btSphereShape* sphereShape = new btSphereShape(0.5f); // Rayon de la sphère
    btRigidBody::btRigidBodyConstructionInfo info(1.f, nullptr, sphereShape);
    sphere = new btRigidBody(info);
    sphere->setLinearVelocity(btVector3(rayDirection.x(), rayDirection.y(), rayDirection.z()));

    world->addRigidBody(sphere);
}

void Game::drawEvent() {
    imguiContext.newFrame();

    ImGui::Begin("Stats");
    ImGui::Text("Simulation speed: %f", 1.f / deltaTime);
    ImGui::End();

    imguiContext.render();
}
