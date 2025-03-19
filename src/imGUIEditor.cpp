#include <Magnum/Math/Color.h>
#include <Magnum/Math/Time.h>
#include <Magnum/GL/DefaultFramebuffer.h>
#include <Magnum/GL/Renderer.h>
#include <Magnum/ImGuiIntegration/Context.hpp>
#include <imgui.h>
#include <fstream>
#include <nlohmann/json.hpp>

using namespace Magnum;
using namespace Math::Literals;

namespace Game {
    struct Cube {
        Magnum::Vector3 position;
    };

    class imGUIEditor {
        public:
            void drawEvent();

            void drawUI();

            void saveScene(const std::string& filename);

            void loadScene(const std::string& filename);

            void createDefaultSceneFile(const std::string& filename);

        private:
            ImGuiIntegration::Context _imgui{NoCreate};
            std::vector<Cube> _cubes; // Liste des cubes de la scène

    };

    void imGUIEditor::drawUI() {
        GL::defaultFramebuffer.clear(GL::FramebufferClear::Color);
        _imgui.newFrame();

        ImGui::Begin("Scene Editor");

        if(ImGui::Button("Add a cube")) {
            _cubes.push_back({Magnum::Vector3(0.0f, 0.0f, 10.0f)});
        }

        for(int i = 0; i < _cubes.size(); ++i) {
            ImGui::PushID(i);
            ImGui::InputFloat3("Position", &_cubes[i].position[0]);
            if(ImGui::Button("Delete")) {
                _cubes.erase(_cubes.begin() + i);
            }
            ImGui::PopID();
        }

        if(ImGui::Button("Save scene")) {
            saveScene("scene.json");
        }

        ImGui::End();
    }

    void imGUIEditor::createDefaultSceneFile(const std::string &filename) {
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

    void imGUIEditor::loadScene(const std::string& filename) { // TODO: a appeler au lancement du jeu
        std::ifstream file(filename);
        if (!file) return;

        nlohmann::json scene;
        file >> scene;

        _cubes.clear();
        for (const auto& cubeData : scene["cubes"]) {
            _cubes.push_back({{cubeData["x"], cubeData["y"], cubeData["z"]}});
        }
    }

    void imGUIEditor::saveScene(const std::string& filename) {
        nlohmann::json scene;
        for (const auto& cube : _cubes) {
            scene["cubes"].push_back({{"x", cube.position.x()}, {"y", cube.position.y()}, {"z", cube.position.z()}});
        }
        std::ofstream file(filename);
        file << scene.dump(4);
    }

}




