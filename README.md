# Magnum Bootstrap

A self-contained bootstrap for Magnum GLFW with Bullet that compile on Windows, Linux and MacOS

## Equipe : 
* Tarek Bouchema (BOUT06050300)
* Mathieu Sparfel-Monnot (SPAM13050100)
* Natal Housset (HOUN25110300)

## Contenu : 
- [x] Connection Client/Serveur
- [x] Serialisation/Synchronisation des gameobjects
- [ ] Serialisation/Synchronisation des composants
- [x]  Serialisation/Synchronisation du world
- [x] Mise en place de l'api
- [ ] Utilisation de l'api
- [ ] Boucle de jeu

Les inputs clavier des utilisateurs sont serialisé et envoyé au serveur mais ne sont pas utilisé hors d'un print dans les logs coté serveur.


## Howto : 
- Lancer la target "DedicatedServer" pour démarrer le serveur en localhost
- Lancer la target "Client", les clients se connecteront automatiquement au localhost
