#include "GameWindow.h"
#include "imgui.h"

GameWindow::GameWindow() {
    // initialize
}

GameWindow::~GameWindow() {
    // clean up
}

void GameWindow::Render() {
    // Récupérer la taille de l'écran
    ImVec2 screenSize = ImGui::GetIO().DisplaySize;

    // Créer une fenêtre de menu
    ImGui::SetNextWindowPos(ImVec2(10, 10));
    ImGui::SetNextWindowSize(ImVec2(300, 500));
    ImGui::Begin("Stats", nullptr, ImGuiWindowFlags_AlwaysAutoResize);

    ImGui::Text("Bienvenue sur Epic Game Cube !");
    ImGui::Separator();

    if (ImGui::CollapsingHeader("Options du jeu")) {

        static float cube = 0.5f;
        ImGui::SliderFloat("Cube écartés", &cube, 0.0f, 1.0f);

        ImGui::Text("Autres options...");
    }

    ImGui::Separator();

    ImGui::End();
}
