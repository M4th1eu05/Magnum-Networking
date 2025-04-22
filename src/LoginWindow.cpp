#include "LoginWindow.h"
#include "imgui.h"
#include <cstring>  // Pour utiliser strcmp()

LoginWindow::LoginWindow()
    : loggedIn(false), showError(false) {
    username[0] = '\0';
    password[0] = '\0';
}

LoginWindow::~LoginWindow() {}

void LoginWindow::Render() {
    // Récupérer la taille de l'écran
    ImVec2 screenSize = ImGui::GetIO().DisplaySize;

    // Configurations pour une fenêtre plein écran
    ImGui::SetNextWindowPos(ImVec2(0, 0));
    ImGui::SetNextWindowSize(screenSize);

    // Flags pour une fenêtre plein écran sans éléments de décoration
    ImGuiWindowFlags window_flags = ImGuiWindowFlags_NoTitleBar |
                                    ImGuiWindowFlags_NoResize |
                                    ImGuiWindowFlags_NoMove |
                                    ImGuiWindowFlags_NoScrollbar |
                                    ImGuiWindowFlags_NoCollapse |
                                    ImGuiWindowFlags_NoSavedSettings |
                                    ImGuiWindowFlags_NoBackground;

    // Commencer la fenêtre avec les flags spécifiés
    ImGui::Begin("FullScreenBackground", nullptr, window_flags);

    // Créer un fond semi-transparent
    ImVec2 winPos = ImGui::GetWindowPos();
    ImDrawList* drawList = ImGui::GetWindowDrawList();
    drawList->AddRectFilled(winPos,
                          ImVec2(winPos.x + screenSize.x, winPos.y + screenSize.y),
                          ImGui::ColorConvertFloat4ToU32(ImVec4(0.1f, 0.1f, 0.1f, 0.8f)));

    // Positionnement de la boîte de login
    float loginBoxWidth = 300.0f;
    float loginBoxHeight = 200.0f; // Augmenté pour l'espace vertical supplémentaire
    ImVec2 loginBoxPos = ImVec2((screenSize.x - loginBoxWidth) * 0.5f,
                               (screenSize.y - loginBoxHeight) * 0.5f);

    // Dessiner la boîte de login
    drawList->AddRectFilled(loginBoxPos,
                          ImVec2(loginBoxPos.x + loginBoxWidth, loginBoxPos.y + loginBoxHeight),
                          ImGui::ColorConvertFloat4ToU32(ImVec4(0.2f, 0.2f, 0.2f, 0.9f)),
                          10.0f); // Coins arrondis

    // Positionner les éléments de login
    ImGui::SetCursorPos(ImVec2(loginBoxPos.x + 20, loginBoxPos.y + 20));

    ImGui::BeginGroup();
    ImGui::Text("Login");
    ImGui::Spacing();
    ImGui::Spacing();
    ImGui::Spacing();

    // Style pour les champs de saisie
    ImGui::PushItemWidth(loginBoxWidth - 40);

    // Étiquette au-dessus du champ de saisie pour le nom d'utilisateur
    ImGui::Text("Username");
    ImGui::InputText("##username", username, sizeof(username));

    ImGui::Spacing();

    // Étiquette au-dessus du champ de saisie pour le mot de passe
    ImGui::Text("Password");
    ImGui::InputText("##password", password, sizeof(password));

    ImGui::PopItemWidth();

    if (showError) {
        ImGui::Spacing();
        ImGui::TextColored(ImVec4(1.0f, 0.0f, 0.0f, 1.0f), "Invalid username or password");
    }

    ImGui::Spacing();
    ImGui::SetCursorPos(ImVec2(loginBoxPos.x + 20, loginBoxPos.y + loginBoxHeight - 40));
    if (ImGui::Button("Login", ImVec2(loginBoxWidth - 40, 30))) {
        Authenticate();
    }

    ImGui::EndGroup();

    ImGui::End();
}

void LoginWindow::Authenticate() {
    // check inputs
    if (strcmp(username, "admin") == 0 && strcmp(password, "password123") == 0) {
        loggedIn = true;
        showError = false;
    } else {
        loggedIn = false;
        showError = true;
    }
}

bool LoginWindow::IsLoggedIn() const {
    return loggedIn;
}
