# Livrables_projet_groupie_tracker

Ce dépôt contient une solution de gestion de ressources musicales, offrant diverses fonctionnalités telles que la recherche, la consultation en détails, la gestion de favoris et bien plus encore. Ce guide vous aidera à démarrer avec la solution et vous fournira une liste détaillée des routes disponibles.

## Lancement de la Solution

Pour lancer la solution localement, suivez ces étapes :

1. **Téléchargement :** Téléchargez le dossier du dépôt sur votre machine.

2. **Ouverture dans un IDE :** Ouvrez le dossier dans votre environnement de développement intégré (IDE) préféré.

3. **Exécution :** Lancez la commande suivante dans votre terminal :

    ```bash
    go run .
    ```

4. **Accès au Serveur :** Suivez les instructions affichées dans le terminal pour accéder au serveur. Pour arrêter le serveur, utilisez `CTRL + C`. Pour accéder au site, cliquez sur le lien affiché dans le terminal.

## Liste des Routes et Leurs Fonctionnalités

- **"/index"** :
  - Page d'accueil affichant des playlists générées aléatoirement.

- **"/detail"** :
  - Page de détails fournissant des informations sur les playlists, albums, artistes et morceaux.

- **"/category"** :
  - Page de catégorie affichant des ressources spécifiques à un endpoint avec pagination.

- **"/search"** :
  - Page de recherche permettant de rechercher dans les playlists, albums, artistes et morceaux avec des filtres sur le nombre de followers, le genre musical,et trie par ordre alphabétique.

- **"/propos"** :
  - Page à propos expliquant le projet.

- **"/favoris"** :
  - Page du compte permettant de visualiser les ressources mises en favoris et de changer l'image de profil.

- **"/treatment/favoris"** :
  - Route de traitement pour les favoris.

- **"/suppr"** :
  - Route pour supprimer un favori.

- **"/img"** :
  - Route pour initialiser la nouvelle image de profil de l'utilisateur.

- **"/play"** :
  - Route pour écouter un morceau.

- **"/url"** :
  - Route pour mettre à jour le navigateur (système de retour en arrière ou en avant sur une page).

- **"/login"** :
  - Page de connexion.

- **"/login/treatment"** :
  - Route de traitement des connexions.

- **"/inscription"** :
  - Page d'inscription.

- **"/inscription/treatment"** :
  - Route de traitement des inscriptions.

- **"/logout"** :
  - Page de déconnexion.

- **"/"** :
  - Page d'erreur 404.

Explorez ces routes pour découvrir les fonctionnalités offertes par la solution !

**Remarque :** Si vous avez des questions ou des retours, n'hésitez pas à me contacter.

**Bonne Exploration !** 🎵🚀

