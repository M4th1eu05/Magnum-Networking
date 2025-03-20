#include <btBulletDynamicsCommon.h>
#include <Corrade/Containers/GrowableArray.h>
#include <Corrade/Containers/Optional.h>
#include <Corrade/Containers/Pointer.h>
#include <Magnum/Timeline.h>
#include <Magnum/BulletIntegration/Integration.h>
#include <Magnum/BulletIntegration/MotionState.h>
#include <Magnum/BulletIntegration/DebugDraw.h>
#include <Magnum/GL/DefaultFramebuffer.h>
#include <Magnum/GL/Mesh.h>
#include <Magnum/GL/Renderer.h>
#include <Magnum/Math/Constants.h>
#include <Magnum/Math/Color.h>
#include <Magnum/Math/Time.h>
#include <Magnum/MeshTools/Compile.h>
#include <Magnum/MeshTools/Transform.h>
#include <Magnum/Platform/GlfwApplication.h>
#include <Magnum/Primitives/Cube.h>
#include <Magnum/Primitives/UVSphere.h>
#include <Magnum/SceneGraph/Camera.h>
#include <Magnum/SceneGraph/Drawable.h>
#include <Magnum/SceneGraph/MatrixTransformation3D.h>
#include <Magnum/SceneGraph/Scene.h>
#include <Magnum/Shaders/PhongGL.h>
#include <Magnum/Trade/MeshData.h>
#include <imgui.h>
#include <fstream>
#include <nlohmann/json.hpp>
#include <Magnum/ImGuiIntegration/Context.hpp>
#include <iostream>

using namespace Magnum;
using namespace Math::Literals;

typedef SceneGraph::Object<SceneGraph::MatrixTransformation3D> Object3D;
typedef SceneGraph::Scene<SceneGraph::MatrixTransformation3D> Scene3D;

namespace Game {
    struct InstanceData {
        Matrix4 transformationMatrix;
        Matrix3x3 normalMatrix;
        Color3 color;
    };

    struct Cube {
        Magnum::Vector3 position;
    };

    class GameApp : public Platform::Application {
    public:
        virtual ~GameApp() = default;

        explicit GameApp(const Arguments &arguments);

    private:
        void drawEvent() override;

        void saveScene(const std::string& filename);

        void loadScene(const std::string& filename);

        void createDefaultSceneFile(const std::string& filename);

        void drawUI();

        void keyPressEvent(KeyEvent &event) override;

        void pointerPressEvent(PointerEvent &event) override;

        void pointerMoveEvent(PointerMoveEvent &event) override;

        GL::Mesh _box{NoCreate}, _sphere{NoCreate};
        GL::Buffer _boxInstanceBuffer{NoCreate}, _sphereInstanceBuffer{NoCreate};
        Shaders::PhongGL _shader{NoCreate};
        BulletIntegration::DebugDraw _debugDraw{NoCreate};
        Containers::Array<InstanceData> _boxInstanceData, _sphereInstanceData;

        btDbvtBroadphase _bBroadphase;
        btDefaultCollisionConfiguration _bCollisionConfig;
        btCollisionDispatcher _bDispatcher{&_bCollisionConfig};
        btSequentialImpulseConstraintSolver _bSolver;

        /* The world has to live longer than the scene because RigidBody
           instances have to remove themselves from it on destruction */
        btDiscreteDynamicsWorld _bWorld{&_bDispatcher, &_bBroadphase, &_bSolver, &_bCollisionConfig};

        Scene3D _scene;
        SceneGraph::Camera3D *_camera;
        SceneGraph::DrawableGroup3D _drawables;
        Timeline _timeline;

        Object3D *_cameraRig, *_cameraObject;

        btBoxShape _bBoxShape{{0.5f, 0.5f, 0.5f}};
        btSphereShape _bSphereShape{0.25f};
        btBoxShape _bGroundShape{{100.0f, 0.5f, 100.0f}};

        bool _drawCubes{true}, _drawDebug{true}, _shootBox{true};

        float _cameraRotationSpeed{0.01f};
        float _cameraMoveSpeed{0.1f};

        std::vector<Cube> _cubes; // Liste des cubes de la scène
        ImGuiIntegration::Context _imgui{NoCreate};
    };

    class ColoredDrawable : public SceneGraph::Drawable3D {
    public:
        explicit ColoredDrawable(Object3D &object, Containers::Array<InstanceData> &instanceData,
                                 const Color3 &color, const Matrix4 &primitiveTransformation,
                                 SceneGraph::DrawableGroup3D &drawables): SceneGraph::Drawable3D
                                                                          {object, &drawables},
                                                                          _instanceData(instanceData),
                                                                          _color{color},
                                                                          _primitiveTransformation{
                                                                              primitiveTransformation
                                                                          } {
        }

    private:
        void draw(const Matrix4 &transformation, SceneGraph::Camera3D &) override {
            const Matrix4 t = transformation * _primitiveTransformation;
            arrayAppend(_instanceData, InPlaceInit, t, t.normalMatrix(), _color);
        }

        Containers::Array<InstanceData> &_instanceData;
        Color3 _color;
        Matrix4 _primitiveTransformation;
    };

    class RigidBody : public Object3D {
    public:
        RigidBody(Object3D *parent, Float mass, btCollisionShape *bShape,
                  btDynamicsWorld &bWorld): Object3D{parent}, _bWorld(bWorld) {
            /* Calculate inertia so the object reacts as it should with
               rotation and everything */
            btVector3 bInertia(0.0f, 0.0f, 0.0f);
            if (!Math::TypeTraits<Float>::equals(mass, 0.0f))
                bShape->calculateLocalInertia(mass, bInertia);

            /* Bullet rigid body setup */
            auto *motionState = new BulletIntegration::MotionState{*this};
            _bRigidBody.emplace(btRigidBody::btRigidBodyConstructionInfo{
                mass, &motionState->btMotionState(), bShape, bInertia
            });
            _bRigidBody->forceActivationState(DISABLE_DEACTIVATION);
            bWorld.addRigidBody(_bRigidBody.get());
        }

        ~RigidBody() {
            _bWorld.removeRigidBody(_bRigidBody.get());
        }

        btRigidBody &rigidBody() { return *_bRigidBody; }

        /* needed after changing the pose from Magnum side */
        void syncPose() {
            _bRigidBody->setWorldTransform(btTransform(transformationMatrix()));
        }

    private:
        btDynamicsWorld &_bWorld;
        Containers::Pointer<btRigidBody> _bRigidBody;
    };

    GameApp::GameApp(const Arguments &arguments): Platform::Application(arguments, NoCreate) {
        /* Try 8x MSAA, fall back to zero samples if not possible. Enable only 2x
           MSAA if we have enough DPI. */
        {
            const Vector2 dpiScaling = this->dpiScaling({});
            Configuration conf;
            conf.setTitle("Game App")
                    .setSize(conf.size(), dpiScaling);
            GLConfiguration glConf;
            glConf.setSampleCount(dpiScaling.max() < 2.0f ? 8 : 2);
            if (!tryCreate(conf, glConf))
                create(conf, glConf.setSampleCount(0));
        }

        /* Camera setup */
        (*(_cameraRig = new Object3D{&_scene}))
                .translate(Vector3::yAxis(3.0f))
                .rotateY(40.0_degf);
        (*(_cameraObject = new Object3D{_cameraRig}))
                .translate(Vector3::zAxis(20.0f))
                .rotateX(-25.0_degf);
        (_camera = new SceneGraph::Camera3D(*_cameraObject))
                ->setAspectRatioPolicy(SceneGraph::AspectRatioPolicy::Extend)
                .setProjectionMatrix(Matrix4::perspectiveProjection(35.0_degf, 1.0f, 0.1f, 1000.0f))
                .setViewport(GL::defaultFramebuffer.viewport().size());


        /* Create an instanced shader */
        _shader = Shaders::PhongGL{
            Shaders::PhongGL::Configuration{}
            .setFlags(Shaders::PhongGL::Flag::VertexColor |
                      Shaders::PhongGL::Flag::InstancedTransformation)
        };
        _shader.setAmbientColor(0x111111_rgbf)
                .setSpecularColor(0x330000_rgbf)
                .setLightPositions({{10.0f, 15.0f, 5.0f, 0.0f}});

        /* Box and sphere mesh, with an (initially empty) instance buffer */
        _box = MeshTools::compile(Primitives::cubeSolid());
        _sphere = MeshTools::compile(Primitives::uvSphereSolid(16, 32));
        _boxInstanceBuffer = GL::Buffer{};
        _sphereInstanceBuffer = GL::Buffer{};
        _box.addVertexBufferInstanced(_boxInstanceBuffer, 1, 0,
                                      Shaders::PhongGL::TransformationMatrix{},
                                      Shaders::PhongGL::NormalMatrix{},
                                      Shaders::PhongGL::Color3{});
        _sphere.addVertexBufferInstanced(_sphereInstanceBuffer, 1, 0,
                                         Shaders::PhongGL::TransformationMatrix{},
                                         Shaders::PhongGL::NormalMatrix{},
                                         Shaders::PhongGL::Color3{});

        /* Set up the renderer so we can draw the debug lines on top */
        GL::Renderer::enable(GL::Renderer::Feature::DepthTest);
        GL::Renderer::enable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::enable(GL::Renderer::Feature::PolygonOffsetFill);
        GL::Renderer::setPolygonOffset(2.0f, 0.5f);

        /* Bullet setup */
        _debugDraw = BulletIntegration::DebugDraw{};
        _debugDraw.setMode(BulletIntegration::DebugDraw::Mode::DrawWireframe);
        _bWorld.setGravity({0.0f, -10.0f, 0.0f});
        _bWorld.setDebugDrawer(&_debugDraw);

        /* Create the ground */
        auto *ground = new RigidBody{&_scene, 0.0f, &_bGroundShape, _bWorld};
        new ColoredDrawable{
            *ground, _boxInstanceData, 0xffffff_rgbf,
            Matrix4::scaling({100.0f, 0.5f, 100.0f}), _drawables
        };

        /* Create boxes with random colors */
        Deg hue = 42.0_degf;
        for (Int i = 0; i != 10; ++i) {
            for (Int j = 0; j != 10; ++j) {
                for (Int k = 0; k != 5; ++k) {
                    auto *o = new RigidBody{&_scene, 1.0f, &_bBoxShape, _bWorld};
                    o->translate({i + 1.0f , j + 5.0f, k + 1.0f});
                    o->syncPose();
                    new ColoredDrawable{
                        *o, _boxInstanceData,
                        Color3::fromHsv({hue += 137.5_degf, 0.75f, 0.9f}),
                        Matrix4::scaling(Vector3{0.5f}), _drawables
                    };
                }
            }
        }

        /* Loop at 60 Hz max */
        setSwapInterval(1);
        setMinimalLoopPeriod(16.0_msec);
        _timeline.start();
        _imgui = ImGuiIntegration::Context(Vector2{windowSize()}/dpiScaling(),
        windowSize(), framebufferSize());

        GL::Renderer::setBlendEquation(GL::Renderer::BlendEquation::Add,
        GL::Renderer::BlendEquation::Add);
        GL::Renderer::setBlendFunction(GL::Renderer::BlendFunction::SourceAlpha,
        GL::Renderer::BlendFunction::OneMinusSourceAlpha);
    }

    void GameApp::drawEvent() {
        GL::defaultFramebuffer.clear(GL::FramebufferClear::Color | GL::FramebufferClear::Depth);

        /* Housekeeping: remove any objects which are far away from the origin */
        for (Object3D *obj = _scene.children().first(); obj;) {
            Object3D *next = obj->nextSibling();
            if (obj->transformation().translation().dot() > 100 * 100)
                delete obj;

            obj = next;
        }

        /* Step bullet simulation */
        _bWorld.stepSimulation(_timeline.previousFrameDuration(), 5);

        if (_drawCubes) {
            /* Populate instance data with transformations and colors */
            arrayResize(_boxInstanceData, 0);
            arrayResize(_sphereInstanceData, 0);

            /* Draw the objects */
            _camera->draw(_drawables);

            _shader.setProjectionMatrix(_camera->projectionMatrix());

            /* Upload instance data to the GPU (orphaning the previous buffer
               contents) and draw all cubes in one call, and all spheres (if any)
               in another call */
            _boxInstanceBuffer.setData(_boxInstanceData, GL::BufferUsage::DynamicDraw);
            _box.setInstanceCount(_boxInstanceData.size());
            _shader.draw(_box);

            _sphereInstanceBuffer.setData(_sphereInstanceData, GL::BufferUsage::DynamicDraw);
            _sphere.setInstanceCount(_sphereInstanceData.size());
            _shader.draw(_sphere);
        }

        /* Debug draw. If drawing on top of cubes, avoid flickering by setting
           depth function to <= instead of just <. */
        if (_drawDebug) {
            if (_drawCubes)
                GL::Renderer::setDepthFunction(GL::Renderer::DepthFunction::LessOrEqual);

            _debugDraw.setTransformationProjectionMatrix(
                _camera->projectionMatrix() * _camera->cameraMatrix());
            _bWorld.debugDrawWorld();

            if (_drawCubes)
                GL::Renderer::setDepthFunction(GL::Renderer::DepthFunction::Less);
        }

        drawUI();
        swapBuffers();
        _timeline.nextFrame();
        redraw();


    }

    void GameApp::drawUI() {
        _imgui.newFrame();

        ImGui::Begin("Scene Editor");
        ImGui::SetWindowSize(ImVec2(500,100), ImGuiCond_FirstUseEver);
        if(ImGui::Button("Add a cube")) {
            std::cout << "[DEBUG] add button pressed !" << std::endl;
            _cubes.push_back({Magnum::Vector3(0.0f, 0.0f, 10.0f)});
        }

        for(int i = 0; i < _cubes.size(); ++i) {
            ImGui::PushID(i);
            ImGui::InputFloat3("Position", &_cubes[i].position[0]);
            if(ImGui::Button("Delete")) {
                std::cout << "[DEBUG] delete button pressed !" << std::endl;
                _cubes.erase(_cubes.begin() + i);
            }
            ImGui::PopID();
        }

        if(ImGui::Button("Save scene")) {
            std::cout << "[DEBUG] save button pressed !" << std::endl;
            saveScene("scene.json");
        }

        ImGui::End();

        _imgui.updateApplicationCursor(*this);

        GL::Renderer::enable(GL::Renderer::Feature::Blending);
        GL::Renderer::enable(GL::Renderer::Feature::ScissorTest);
        GL::Renderer::disable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::disable(GL::Renderer::Feature::DepthTest);

        _imgui.drawFrame();

        GL::Renderer::enable(GL::Renderer::Feature::DepthTest);
        GL::Renderer::enable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::disable(GL::Renderer::Feature::ScissorTest);
        GL::Renderer::disable(GL::Renderer::Feature::Blending);
    }

    void GameApp::createDefaultSceneFile(const std::string &filename) {
        nlohmann::json scene;
        scene["cubes"] = {
            {{"x", 0.0}, {"y", 0.0}, {"z", 0.0}},
            {{"x", 1.0}, {"y", 0.0}, {"z", 0.0}}
        };

        std::ofstream file(filename);
        if (file) {
            file << scene.dump(4); // Écrit en format JSON avec indentation
            file.close();
        }
    }

    void GameApp::loadScene(const std::string& filename) { // TODO: a appeler au lancement du jeu
        std::ifstream file(filename);
        if (!file) return;

        nlohmann::json scene;
        file >> scene;

        _cubes.clear();
        for (const auto& cubeData : scene["cubes"]) {
            _cubes.push_back({{cubeData["x"], cubeData["y"], cubeData["z"]}});
        }
    }

    void GameApp::saveScene(const std::string& filename) {
        nlohmann::json scene;
        for (const auto& cube : _cubes) {
            scene["cubes"].push_back({{"x", cube.position.x()}, {"y", cube.position.y()}, {"z", cube.position.z()}});
        }
        std::ofstream file(filename);
        file << scene.dump(4);
    }


    void GameApp::keyPressEvent(KeyEvent &event) {
        /* Movement */
        if (event.key() == Key::W) {
            _cameraObject->translate(Vector3::zAxis(-_cameraMoveSpeed));
        } else if (event.key() == Key::S) {
            _cameraObject->translate(Vector3::zAxis(_cameraMoveSpeed));
        } else if (event.key() == Key::A) {
            _cameraObject->translate(Vector3::xAxis(-_cameraMoveSpeed));
        } else if (event.key() == Key::D) {
            _cameraObject->translate(Vector3::xAxis(_cameraMoveSpeed));
        } else if (event.key() == Key::Q) {
            _cameraObject->translate(Vector3::yAxis(-_cameraMoveSpeed));
        } else if (event.key() == Key::E) {
            _cameraObject->translate(Vector3::yAxis(_cameraMoveSpeed));

            /* Toggling draw modes */
        } else if (event.key() == Key::D) {
            if (_drawCubes && _drawDebug) {
                _drawDebug = false;
            } else if (_drawCubes && !_drawDebug) {
                _drawCubes = false;
                _drawDebug = true;
            } else if (!_drawCubes && _drawDebug) {
                _drawCubes = true;
                _drawDebug = true;
            }

            /* What to shoot */
        } else if (event.key() == Key::C) {
            _shootBox ^= true;
        } else return;

        event.setAccepted();
    }

    void GameApp::pointerPressEvent(PointerEvent &event) {

        if (ImGui::GetIO().WantCaptureMouse) return;
        /* Shoot an object on click */
        if (!event.isPrimary() || !(event.pointer() & (Pointer::MouseLeft)))
            return;

        /* First scale the position from being relative to window size to being
           relative to framebuffer size as those two can be different on HiDPI
           systems */
        const Vector2 position = event.position() * Vector2{framebufferSize()} / Vector2{windowSize()};
        const Vector2 clickPoint = Vector2::yScale(-1.0f) * (position / Vector2{framebufferSize()} - Vector2{0.5f})
                                   * _camera->projectionSize();
        const Vector3 direction = (_cameraObject->absoluteTransformation().rotationScaling() * Vector3{
                                       clickPoint, -1.0f
                                   }).normalized();

        auto *object = new RigidBody{
            &_scene,
            _shootBox ? 1.0f : 5.0f,
            _shootBox ? static_cast<btCollisionShape *>(&_bBoxShape) : &_bSphereShape,
            _bWorld
        };
        object->translate(_cameraObject->absoluteTransformation().translation());
        /* Has to be done explicitly after the translate() above, as Magnum ->
           Bullet updates are implicitly done only for kinematic bodies */
        object->syncPose();

        /* Create either a box or a sphere */
        new ColoredDrawable{
            *object,
            _shootBox ? _boxInstanceData : _sphereInstanceData,
            _shootBox ? 0x880000_rgbf : 0x220000_rgbf,
            Matrix4::scaling(Vector3{_shootBox ? 0.5f : 0.25f}), _drawables
        };

        /* Give it an initial velocity */
        object->rigidBody().setLinearVelocity(btVector3{direction * 25.f});

        event.setAccepted();
    }

    void GameApp::pointerMoveEvent(PointerMoveEvent &event) {
        /* Rotate the camera on mouse drag */
        if (!event.isPrimary() ||
            !(event.pointers() & (Pointer::MouseRight)))
            return;

        Vector2 delta = _cameraRotationSpeed * Vector2{event.relativePosition()};

        _cameraObject->rotateX(Rad{delta.y()});
        _cameraRig->rotateY(Rad{delta.x()});

        event.setAccepted();
        redraw();
    }
}

MAGNUM_APPLICATION_MAIN(Game::GameApp)
