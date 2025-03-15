#include <Magnum/GL/DefaultFramebuffer.h>
#include <Magnum/Platform/GlfwApplication.h>
#include <Magnum/GL/Mesh.h>
#include <Magnum/GL/Renderer.h>
#include <Magnum/Math/Angle.h>
#include <Magnum/Math/Color.h>
#include <Magnum/Math/Matrix4.h>
#include <Magnum/MeshTools/Compile.h>
#include <Magnum/Primitives/Cube.h>
#include <Magnum/Shaders/PhongGL.h>
#include <Magnum/Trade/MeshData.h>

using namespace Magnum;
using namespace Math::Literals;

namespace Game {

    class BasicCubeApp : public Platform::Application {
    public:
        virtual ~BasicCubeApp() = default;

        explicit BasicCubeApp(const Arguments &arguments);

    private:
        void drawEvent() override;

        void pointerReleaseEvent(PointerEvent &event) override;

        void pointerMoveEvent(PointerMoveEvent &event) override;

        GL::Mesh _mesh;
        Shaders::PhongGL _shader;

        Matrix4 _transformation, _projection;
        Color3 _color;
    };

    BasicCubeApp::BasicCubeApp(const Arguments &arguments): Platform::Application{arguments} {
        GL::Renderer::enable(GL::Renderer::Feature::DepthTest);
        GL::Renderer::enable(GL::Renderer::Feature::FaceCulling);
        _mesh = MeshTools::compile(Primitives::cubeSolid());
        _transformation =
                Matrix4::rotationX(30.0_degf) * Matrix4::rotationY(40.0_degf);
        _projection =
                Matrix4::perspectiveProjection(
                    35.0_degf, Vector2{windowSize()}.aspectRatio(), 0.01f, 100.0f) *
                Matrix4::translation(Vector3::zAxis(-10.0f));
        _color = Color3::fromHsv({35.0_degf, 1.0f, 1.0f});
    }

    void BasicCubeApp::drawEvent() {
        GL::defaultFramebuffer.clear(
            GL::FramebufferClear::Color | GL::FramebufferClear::Depth);

        _shader.setLightPositions({{1.4f, 1.0f, 0.75f, 0.0f}})
                .setDiffuseColor(_color)
                .setAmbientColor(Color3::fromHsv({_color.hue(), 1.0f, 0.3f}))
                .setTransformationMatrix(_transformation)
                .setNormalMatrix(_transformation.normalMatrix())
                .setProjectionMatrix(_projection)
                .draw(_mesh);

        swapBuffers();
    }

    void BasicCubeApp::pointerReleaseEvent(PointerEvent &event) {
        if (!event.isPrimary() ||
            !(event.pointer() & (Pointer::MouseLeft)))
            return;

        _color = Color3::fromHsv({_color.hue() + 50.0_degf, 1.0f, 1.0f});

        event.setAccepted();
        redraw();
    }

    void BasicCubeApp::pointerMoveEvent(PointerMoveEvent &event) {
        if (!event.isPrimary() ||
            !(event.pointers() & (Pointer::MouseLeft)))
            return;

        Vector2 delta = 3.0f * Vector2{event.relativePosition()} / Vector2{windowSize()};

        _transformation =
                Matrix4::rotationX(Rad{delta.y()}) *
                _transformation *
                Matrix4::rotationY(Rad{delta.x()});

        event.setAccepted();
        redraw();
    }
}

MAGNUM_APPLICATION_MAIN(Game::BasicCubeApp)
