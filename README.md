# EffectiveMobileGoLang
Тестовое задание:

Реализовать каталог автомобилей. Необходимо реализовать следующее
1. Выставить rest методы
	1. Получение данных с фильтрацией по всем полям и пагинацией 
	2. Удаления по идентификатору
	3. Изменение одного или нескольких полей по идентификатору
	4. Добавления новых автомобилей в формате
```json
{
    "regNums": ["X123XX150"] // массив гос. номеров
}
```
2. При добавлении сделать запрос во внешнее АПИ, описанного сваггером (это описание некоторого внешнего АПИ, которого нет, но к которому надо обращаться. Реализованное, согласно описанию, АПИ будет использоваться при проверке)

```yaml
openapi: 3.0.3
info:
  title: Car info
  version: 0.0.1
paths:
  /info:
    get:
      parameters:
        - name: regNum
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Car'
        '400':
          description: Bad request
        '500':
          description: Internal server error
components:
  schemas:
    Car:
      required:
        - regNum
        - mark
        - model
        - owner
      type: object
      properties:
        regNum:
          type: string
          example: X123XX150
        mark:
          type: string
          example: Lada
        model:
          type: string
          example: Vesta
        year:
          type: integer
          example: 2002
        owner:
          $ref: '#/components/schemas/People'
    People:
      required:
        - name
        - surname
      type: object
      properties:
        name:
          type: string
        surname:
          type: string
        patronymic:
          type: string
```
3. Обогащенную информацию положить в БД postgres (структура БД должна быть создана путем миграций при старте сервиса)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env-файл
6. Сгенерировать сваггер на реализованное АПИ

# Выполненная работа 

Описание сервиса:

Во время запуска сервера происходят следующие основные действия:
1. Инициализация логгера
2. Загрузка конфига
3. Конект к БД
4. Инициализация репозиториев (был использован паттерн проектирования "Репозиторий")
5. Откат миграций (удаление таблиц из БД, если существуют)
6. Использование миграций (создание таблиц как в задании)
7. Заполнение БД данными для тестов
   
Уточнения к REST-методам:
1. "Получение данных с фильтрацией по всем полям и пагинацией" - не совсем понятно, что есть пагинация в данном случае, по этому просто был реализован метод, возвращающий всю информацию о автомобиле по его индификатору.
2. "Добавления новых автомобилей в формате" - формат был выбран следующий:
```json
{
  "regNums": ["X000XX150", "Y456YY789"],
  "mark": ["Lada", "Lada"],
  "model": ["Vesta", "Sport"],
  "year": [2010, 2020],
  "owner": [{
    "name": "Сергей",
    "surname": "Сидоров",
    "patronymic": "Витальевич"
  }, {
    "name": "Антон",
    "surname": "Губин",
    "patronymic": ""
  }]
}
```

Уточнения к БД:
1. Так как в сваггере не было указано поле с ID в таблице Car, считаем поле "regNum" уникальным индификатором.
2. Так как в сваггере не было указано поле с ID в таблице People, считаем поля "name" и "surname" уникальными индификатороми. Так же считаем, что полных тёзок в БД нет.

Сгенериованный сваггер (дублируется в файлах проекта):
```yaml
openapi: 3.0.3
info:
  title: Каталог автомобилей OpenAPI спецификация
  version: 0.0.1
servers:
  - url: http://localhost:8080/
    description: Dev server
paths:
  /info:
    get:
      summary: Метод получения автомобиля по его идентификатору
      tags:
        - Car
      parameters:
        - name: regNum
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Успешный ответ с информацией об автомобиле
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Car'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                type: string
                example: Parameter 'regNum' is required
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                type: string
                example: 'Failed to fetch car information: [error]'
  /add:
    post:
      summary: Метод добавления новых автомобилей
      tags:
        - Car
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                regNums:
                  type: array
                  items:
                    type: string
                  example: ["X000XX150", "Y456YY789"]
                mark:
                  type: array
                  items:
                    type: string
                  example: ["Lada", "Lada"]
                model:
                  type: array
                  items:
                    type: string
                  example: ["Vesta", "Sport"]
                year:
                  type: array
                  items:
                    type: integer
                  example: [2010, 2020]
                owner:
                  type: array
                  items:
                    type: object
                    properties:
                      name:
                        type: string
                        example: "Сергей"
                      surname:
                        type: string
                        example: "Сидоров"
                      patronymic:
                        type: string
                        example: "Витальевич"
                    example: 
                      - name: "Сергей"
                        surname: "Сидоров"
                        patronymic: "Витальевич"
                      - name: "Антон"
                        surname: "Губин"
                        patronymic: ""
      responses:
        '200':
          description: Успешный ответ с информацией об автомобиле
          content:
            application/json:
              schema:
                type: string
                example: 'Cars added successfully'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                  type: string
                  example: 'Failed to add car: data is incomplete'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                  type: string
                  example: 'Failed to add car: [error]'
  /update:
    put:
      summary: Метод изменения одного или нескольких полей по идентификатору
      tags:
        - Car
      requestBody:
        description: Некоторые поля могут быть пустыми
        required: false
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Car'
      responses:
        '200':
          description: Успешный ответ с информацией об автомобиле
          content:
            application/json:
              schema:
                type: string
                example: 'Car updated successfully'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                type: string
                example: Parameter 'regNum' is required
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                type: string
                example: 'Failed to update car: [error]'
  /del:
    delete:
      summary: Метод удаления автомобиля по его индификатору 
      tags:
        - Car
      parameters:
        - name: regNum
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Успешный ответ с информацией об автомобиле
          content:
            application/json:
              schema:
                type: string
                example: 'Car deleted successfully'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                type: string
                example: Parameter 'regNum' is required
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                type: string
                example: 'Failed to delete car: [error]'
components:
  schemas:
    Car:
      required:
        - regNum
        - mark
        - model
        - owner
      type: object
      properties:
        regNum:
          type: string
          example: X123XX150
        mark:
          type: string
          example: Lada
        model:
          type: string
          example: Vesta
        year:
          type: integer
          example: 2002
        owner:
          $ref: '#/components/schemas/People'
    People:
      required:
        - name
        - surname
      type: object
      properties:
        name:
          type: string
          example: Сергей
        surname:
          type: string
          example: Сидоров
        patronymic:
          type: string
          example: Витальевич
```
