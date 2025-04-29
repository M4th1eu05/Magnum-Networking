//
// Created by Tarook on 29/04/2025.
//
#include "../inc/DedicatedServer.h"
#include <iostream>
#include <btBulletDynamicsCommon.h>
#include <Corrade/Containers/GrowableArray.h>
#include <Corrade/Containers/Optional.h>
#include <Corrade/Containers/Pointer.h>
#include <Magnum/Timeline.h>
#include <Magnum/BulletIntegration/Integration.h>
#include <Magnum/BulletIntegration/DebugDraw.h>
#include <Magnum/GL/DefaultFramebuffer.h>
#include <Magnum/GL/Mesh.h>
#include <Magnum/GL/Renderer.h>
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
#include <World.h>
#include <nlohmann/json.hpp>

#include "Components/Collider.h"
#include "Components/Rigidbody.h"
#include "Magnum/ImGuiIntegration/Context.hpp"

using namespace Magnum;
using namespace Math::Literals;

typedef SceneGraph::Object<SceneGraph::MatrixTransformation3D> Object3D;
typedef SceneGraph::Scene<SceneGraph::MatrixTransformation3D> Scene3D;

/*
int main() {
    try {
        const uint16_t port = 1234; // Example port
        const size_t maxClients = 32; // Example max clients

        DedicatedServer server(port, maxClients);
        server.start();

        std::cout << "Press Enter to stop the server..." << std::endl;
        std::cin.get();

        server.stop();
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}
*/
namespace Game {

    struct InstanceData {
        Matrix4 transformationMatrix;
        Matrix3x3 normalMatrix;
        Color3 color;
    };


    class ServerGameApp : public Platform::Application
    {
    public:
        virtual ~ServerGameApp() = default;

        explicit ServerGameApp(const Arguments &arguments);

    private:
        void drawEvent() override;

        void drawUI() const;

        void keyPressEvent(KeyEvent &event) override;
        void keyReleaseEvent(KeyEvent& event) override;

        void pointerPressEvent(PointerEvent &event) override;
        void pointerReleaseEvent(PointerEvent& event) override;
        void scrollEvent(ScrollEvent& event) override;
        void pointerMoveEvent(PointerMoveEvent &event) override;
        void textInputEvent(TextInputEvent& event) override;

        void viewportEvent(ViewportEvent& event) override;

        ImGuiIntegration::Context _imgui{NoCreate};

        GL::Mesh _box{NoCreate}, _sphere{NoCreate};
        GL::Buffer _boxInstanceBuffer{NoCreate}, _sphereInstanceBuffer{NoCreate};
        Shaders::PhongGL _shader{NoCreate};
        Containers::Array<InstanceData> _boxInstanceData, _sphereInstanceData;

        btDbvtBroadphase _bBroadphase;
        btDefaultCollisionConfiguration _bCollisionConfig;
        btCollisionDispatcher _bDispatcher{&_bCollisionConfig};
        btSequentialImpulseConstraintSolver _bSolver;

        SceneGraph::Camera3D *_camera;
        SceneGraph::DrawableGroup3D _drawables;
        Timeline _timeline;

        std::shared_ptr<GameObject>_cameraRig, _cameraObject;

        btBoxShape _bBoxShape{{0.5f, 0.5f, 0.5f}};
        btSphereShape _bSphereShape{0.25f};
        btBoxShape _bGroundShape{{100.0f, 0.5f, 100.0f}};

        bool _drawCubes{true}, _drawDebug{true}, _shootBox{true};

        float _cameraRotationSpeed{0.01f};
        float _cameraMoveSpeed{0.1f};

        std::shared_ptr<World> _world;
        std::shared_ptr<DedicatedServer> _server;
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


    ServerGameApp::ServerGameApp(const Magnum::Platform::GlfwApplication::Arguments& args)
    : Magnum::Platform::Application{args, NoCreate} {
        {
            const Vector2 dpiScaling = this->dpiScaling({});
            Configuration conf;
            conf.setTitle("Server Game App")
                    .setSize(conf.size(), dpiScaling);
            GLConfiguration glConf;
            glConf.setSampleCount(dpiScaling.max() < 2.0f ? 8 : 2);
            if (!tryCreate(conf, glConf))
                create(conf, glConf.setSampleCount(0));

            _imgui = ImGuiIntegration::Context(
                Vector2{windowSize()}/dpiScaling,
                    windowSize(),
                    framebufferSize());

            /* Set up proper blending to be used by ImGui. There's a great chance
               you'll need this exact behavior for the rest of your scene. If not, set
               this only for the drawFrame() call. */
            GL::Renderer::setBlendEquation(GL::Renderer::BlendEquation::Add,
                GL::Renderer::BlendEquation::Add);
            GL::Renderer::setBlendFunction(GL::Renderer::BlendFunction::SourceAlpha,
                GL::Renderer::BlendFunction::OneMinusSourceAlpha);
        }

        _server = std::make_shared<DedicatedServer>(5555, 4);

        _world = std::make_shared<World>(_timeline);

        /* Camera setup */
        _cameraRig = _world->createGameObject();
        _cameraRig->translate(Vector3::yAxis(3.0f)).rotateY(40.0_degf);

        _cameraObject = _world->createGameObject(_cameraRig);
        _cameraObject->translate(Vector3::zAxis(20.0f)).rotateX(-25.0_degf);

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

        /* Setup the renderer so we can draw the debug lines on top */
        GL::Renderer::enable(GL::Renderer::Feature::DepthTest);
        GL::Renderer::enable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::enable(GL::Renderer::Feature::PolygonOffsetFill);
        GL::Renderer::setPolygonOffset(2.0f, 0.5f);

        /* Bullet setup */
        _world->_debugDraw = BulletIntegration::DebugDraw{};
        _world->_debugDraw.setMode(BulletIntegration::DebugDraw::Mode::DrawWireframe);

        /* Create the ground */
        //auto *ground = new RigidBody{&_scene, 0.0f, &_bGroundShape, _bWorld};
        const std::shared_ptr<GameObject> ground = _world->createGameObject();
        ground->addComponent<Collider>(reinterpret_cast<btCollisionShape*>(&_bGroundShape));
        ground->addComponent<Rigidbody>(0.0f);

        new ColoredDrawable{
            *ground, _boxInstanceData, 0xffffff_rgbf,
            Matrix4::scaling({100.0f, 0.5f, 100.0f}), _drawables
        };

        /* Create boxes with random colors */
        Deg hue = 42.0_degf;
        for (Int i = 0; i != 10; ++i) {
            for (Int j = 0; j != 10; ++j) {
                for (Int k = 0; k != 5; ++k) {
                    const std::shared_ptr<GameObject> newBox = _world->createGameObject();
                    auto collider = newBox->addComponent<Collider>(reinterpret_cast<btCollisionShape*>(&_bBoxShape));
                    const auto rb = newBox->addComponent<Rigidbody>(1.0f);
                    newBox->translate({i + 1.0f , j + 5.0f, k + 1.0f});
                    rb->syncPose();
                    new ColoredDrawable{
                        *newBox, _boxInstanceData,
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
    }

    void ServerGameApp::drawEvent() {
        GL::defaultFramebuffer.clear(GL::FramebufferClear::Color | GL::FramebufferClear::Depth);

        /* Housekeeping: remove any objects which are far away from the origin */
        /*
        for (Object3D *obj = _scene.children().first(); obj;) {
            Object3D *next = obj->nextSibling();
            if (obj->transformation().translation().dot() > 100 * 100)
                delete obj;

            obj = next;
        }
        */

        /* Step bullet simulation */
        _world->update();

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

            _world->_debugDraw.setTransformationProjectionMatrix(
                _camera->projectionMatrix() * _camera->cameraMatrix());
            _world->getBulletWorld().debugDrawWorld();

            if (_drawCubes)
                GL::Renderer::setDepthFunction(GL::Renderer::DepthFunction::Less);
        }

        _imgui.newFrame();
        drawUI();
        _imgui.updateApplicationCursor(*this);

        /* Set appropriate states. If you only draw ImGui, it is sufficient to
       just enable blending and scissor test in the constructor. */
        GL::Renderer::enable(GL::Renderer::Feature::Blending);
        GL::Renderer::enable(GL::Renderer::Feature::ScissorTest);
        GL::Renderer::disable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::disable(GL::Renderer::Feature::DepthTest);

        _imgui.drawFrame();

        /* Reset state. Only needed if you want to draw something else with
       different state after. */
        GL::Renderer::enable(GL::Renderer::Feature::DepthTest);
        GL::Renderer::enable(GL::Renderer::Feature::FaceCulling);
        GL::Renderer::disable(GL::Renderer::Feature::ScissorTest);
        GL::Renderer::disable(GL::Renderer::Feature::Blending);

        swapBuffers();
        _timeline.nextFrame();
        redraw();
    }

    void ServerGameApp::drawUI() const {
    }

    void ServerGameApp::keyPressEvent(KeyEvent &event) {
        if(_imgui.handleKeyPressEvent(event)) return;
    }

    void ServerGameApp::keyReleaseEvent(KeyEvent &event) {
        if(_imgui.handleKeyReleaseEvent(event)) return;
    }

    void ServerGameApp::pointerPressEvent(PointerEvent &event) {
        if(_imgui.handlePointerPressEvent(event)) return;
    }

    void ServerGameApp::pointerReleaseEvent(PointerEvent &event) {
        if (_imgui.handlePointerReleaseEvent(event)) return;
    }

    void ServerGameApp::scrollEvent(ScrollEvent &event) {
        if(_imgui.handleScrollEvent(event)) {
            /* Prevent scrolling the page */
            event.setAccepted();
            return;
        }
    }

    void ServerGameApp::pointerMoveEvent(PointerMoveEvent &event) {
        if (_imgui.handlePointerMoveEvent(event)) return;
    }

    void ServerGameApp::textInputEvent(TextInputEvent &event) {
        if(_imgui.handleTextInputEvent(event)) return;
    }

    void ServerGameApp::viewportEvent(ViewportEvent &event) {
        GL::defaultFramebuffer.setViewport({{}, event.framebufferSize()});

        _imgui.relayout(Vector2{event.windowSize()}/event.dpiScaling(),
            event.windowSize(), event.framebufferSize());
    }
}

MAGNUM_APPLICATION_MAIN(Game::ServerGameApp)