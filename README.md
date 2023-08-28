## API Документация

### Аутентификация
Прежде чем сделать API реквест, нужно аутентифицироваться используя валидный access токен. Для этого вставьте ваш токен полученый на эндпойнте /login в хедер вашего запроса
Authorization: Bearer ВАШ_ТОКЕН

### Эндпойнты 
##### POST /user/signup
Создание нового пользователя

##### POST /user/login
Вход

##### POST /notes/new
Создание новой заметки для аутентифицированного пользователя

##### GET /notes
Получение всех заметок для аутентифицированного пользователя


---

### Сборка проекта

Сборка проекта осуществляется через docker </br>
`docker-compose build` </br>
`docker-compose up -d`

---

### Postman коллекция

[![Run in Postman](https://run.pstmn.io/button.svg)](https://god.gw.postman.com/run-collection/29153180-3a2c71e3-9e62-453b-8714-38f132ab43bc?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D29153180-3a2c71e3-9e62-453b-8714-38f132ab43bc%26entityType%3Dcollection%26workspaceId%3Dcd958e8a-baa6-4844-a165-b5ac4e435231)