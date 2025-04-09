
#ifndef LOGINWINDOW_H
#define LOGINWINDOW_H

#include "imgui.h"
#include <string>

class LoginWindow {
public:
    LoginWindow();
    ~LoginWindow();

    void Render();  // method to render the login window
    bool IsLoggedIn() const;  // check if user is logged in

private:
    char username[128];  // buffer for the username
    char password[128];  // buffer for the password
    bool loggedIn;
    bool showError;

    void Authenticate();  // void for authentication logic
};

#endif // LOGINWINDOW_H
