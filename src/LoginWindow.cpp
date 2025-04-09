#include "LoginWindow.h"
#include "imgui.h"
#include <cstring>  // Pour utiliser strcmp()

LoginWindow::LoginWindow()
    : loggedIn(false), showError(false) {
    username[0] = '\0';  // Initialisation du nom d'utilisateur
    password[0] = '\0';  // Initialisation du mot de passe
}

LoginWindow::~LoginWindow() {}

void LoginWindow::Render() {
    ImGui::Begin("Login");

    // Champ pour le nom d'utilisateur
    if (ImGui::InputText("Username", username, sizeof(username))) {
        showError = false;  // Si l'utilisateur modifie le champ, on réinitialise l'erreur
    }

    // Champ pour le mot de passe sans masque
    if (ImGui::InputText("Password", password, sizeof(password))) {
        showError = false;  // Si l'utilisateur modifie le mot de passe, on réinitialise l'erreur
    }

    // Affichage de l'erreur si le login échoue
    if (showError) {
        ImGui::TextColored(ImVec4(1.0f, 0.0f, 0.0f, 1.0f), "Invalid username or password");
    }

    // Bouton de connexion
    if (ImGui::Button("Login")) {
        Authenticate();
    }

    ImGui::End();
}

void LoginWindow::Authenticate() {
    // Vérification des informations de connexion avec les nouvelles données
    if (strcmp(username, "utilisateur1") == 0 && strcmp(password, "mdp1") == 0) {
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
