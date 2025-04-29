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
    : Magnum::Platform::Application{args} {
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

    void ServerGameApp::drawEvent() {
    }

    void ServerGameApp::drawUI() const {
    }

    void ServerGameApp::keyPressEvent(KeyEvent &event) {
    }

    void ServerGameApp::keyReleaseEvent(KeyEvent &event) {
    }

    void ServerGameApp::pointerPressEvent(PointerEvent &event) {
    }

    void ServerGameApp::pointerReleaseEvent(PointerEvent &event) {
    }

    void ServerGameApp::scrollEvent(ScrollEvent &event) {
    }

    void ServerGameApp::pointerMoveEvent(PointerMoveEvent &event) {
    }

    void ServerGameApp::textInputEvent(TextInputEvent &event) {
    }

    void ServerGameApp::viewportEvent(ViewportEvent &event) {
    }
}

MAGNUM_APPLICATION_MAIN(Game::ServerGameApp)