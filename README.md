Yanis Chartier
# Application web en Go · B1 Lyon Ynov Campus

Fil rouge du cours : **GameDex**, un catalogue de jeux vidéo et de leurs tournois. Les jeux sont réels, les tournois sont inventés.

Ce dépôt contient la **mini API** que GameDex interroge, et un dossier `gamedex/` vide : c'est vous qui écrivez GameDex pendant les mini TP.

## Démarrer : fork, puis clone

1. Connectez-vous à GitHub et cliquez sur **Fork** en haut de cette page. Gardez votre fork **public**.
2. Clonez **votre fork**, pas ce dépôt :

```
git clone https://github.com/VOTRE-PSEUDO/go-web-2026.git
cd go-web-2026
```

3. Ouvrez le dossier dans VS Code.

## Rendre son travail

Votre fork **est** votre rendu : c'est lui qui est évalué.

À la fin de chaque mini TP :

```
git add .
git commit -m "Mini TP 2 : premier serveur"
git push
```

Poussez régulièrement, pas seulement le dernier jour : l'historique des commits fait partie du travail.

## Lancer la mini API

Depuis la racine du dépôt, dans un terminal qu'on laisse ouvert :

```
go run ./api
```

Elle répond sur http://localhost:8081/api/ et sa **documentation Swagger** est sur http://localhost:8081/api/docs (la description brute, au format OpenAPI, est sur `/api/openapi.yaml`).

Les endpoints :

| Endpoint | Ce qu'il renvoie |
| --- | --- |
| `GET /api/games` | La liste des jeux |
| `GET /api/games?q=cel&genre=course&year=2018` | Les jeux filtrés : `q` cherche dans le nom (sans tenir compte des majuscules), `genre` et `year` doivent correspondre exactement. Chaque filtre est facultatif |
| `GET /api/games/{id}` | Un jeu (404 s'il n'existe pas) |
| `GET /api/games/{id}/tournaments` | Les tournois d'un jeu |
| `GET /api/tournaments` | Tous les tournois |
| `GET /api/stats` | Des statistiques, protégées par l'en-tête `X-API-Key` (401 ou 403 sans la bonne clé) |
| `GET /api/broken` | Toujours une erreur 500, pour tester |
| `GET /api/slow` | La liste des jeux, après 5 secondes, pour tester les timeouts |

Votre GameDex tourne dans un **second terminal**, depuis le dossier `gamedex/`, sur http://localhost:8080.

## Tester l'API

Ouvrez http://localhost:8081/api/docs. Pour chaque endpoint : déroulez-le, cliquez sur **Try it out**, remplissez les paramètres, puis **Execute**. Swagger affiche l'URL appelée, le code de réponse, les en-têtes et le JSON reçu.

Pour `/api/stats`, la clé se donne avec le bouton **Authorize** en haut de la page.

Dans un terminal, ça marche aussi : `curl.exe -i http://localhost:8081/api/games/3` (sous PowerShell, bien taper `curl.exe`).

## Organisation du dépôt

```
go-web-2026/
├── api/          ← la mini API (ne pas modifier)
├── gamedex/      ← votre travail
├── README.md
└── LEXIQUE.md    ← le vocabulaire du cours
```

## Ressources

- Le protocole HTTP, en français : https://developer.mozilla.org/fr/docs/Web/HTTP/Overview
- Les codes de réponse HTTP : https://developer.mozilla.org/fr/docs/Web/HTTP/Status
- Go by Example, serveur et client HTTP, JSON, templates : https://gobyexample.com
- Le routage de Go 1.22 : https://go.dev/blog/routing-enhancements
- fetch en JavaScript : https://developer.mozilla.org/fr/docs/Web/API/Fetch_API/Using_Fetch
- Leaflet, prise en main : https://leafletjs.com/examples/quick-start/

Le vocabulaire est dans [LEXIQUE.md](LEXIQUE.md).
