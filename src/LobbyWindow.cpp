#include "LobbyWindow.h"
#include "imgui.h"

LobbyWindow::LobbyWindow() {
    gameStarted(false);
}

LobbyWindow::~LobbyWindow() {
    // clean up
}

void LobbyWindow::Render() {
    // Récupérer la taille de l'écran
    ImVec2 screenSize = ImGui::GetIO().DisplaySize;

    // Créer une fenêtre de menu
    ImGui::SetNextWindowPos(ImVec2(10, 10));
    ImGui::SetNextWindowSize(ImVec2(300, 500));
    ImGui::Begin("Menu du jeu", nullptr, ImGuiWindowFlags_AlwaysAutoResize);

    ImGui::Text("Bienvenue dans le jeu !");
    ImGui::Separator();

    if (ImGui::CollapsingHeader("Options du jeu")) {
        static float volume = 0.5f;
        ImGui::SliderFloat("Volume", &volume, 0.0f, 1.0f);
    }

    ImGui::Separator();

    // Boutons de contrôle du jeu
    if (ImGui::Button("Lancer une partie", ImVec2(280, 30))) {
        gameStarted = true;
    }

    ImGui::Separator();

    if (ImGui::Button("Quitter le jeu", ImVec2(280, 30))) {
        exit(0);  // ou Platform::Application::exit(0) pour une fermeture plus propre
    }

    ImGui::End();
}

bool LobbyWindow::IsGameStarted() const {
    return gameStarted;
}
