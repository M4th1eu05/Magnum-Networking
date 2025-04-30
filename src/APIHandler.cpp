//
// Created by mathi on 29/04/2025.
//

#include "APIHandler.h"

#include <iostream>
#include <thread>


bool APIHandler::POST(const std::string &url, const std::string& postData, const std::function<void(std::string)>& callback) {
    try {
        // Lancer une requête asynchrone dans un thread
        std::thread([url, postData, callback]() {
            const cpr::Response response = cpr::Post(
                cpr::Url{url},
                cpr::Body{postData},
                cpr::Header{{"Content-Type", "application/json"}}
            );

            // Appeler le callback avec la réponse
            if (callback) {
                callback(response.text);
            }
        }).detach();

        return true;
    } catch (const std::exception &e) {
        std::cerr << "Erreur lors de l'appel API : " << e.what() << std::endl;
        return false;
    }
}

bool APIHandler::GET(const std::string &url, const std::function<void(std::string)> &callback) {
    try {
        // Lancer une requête asynchrone dans un thread
        std::thread([url, callback]() {
            const cpr::Response response = cpr::Post(
                cpr::Url{url},
                cpr::Header{{"Content-Type", "application/json"}}
            );

            // Appeler le callback avec la réponse
            if (callback) {
                callback(response.text);
            }
        }).detach();

        return true;
    } catch (const std::exception &e) {
        std::cerr << "Erreur lors de l'appel API : " << e.what() << std::endl;
        return false;
    }
}
}
