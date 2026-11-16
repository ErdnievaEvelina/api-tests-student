# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -v ./...
```

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```

# API-тесты для TestBoard (api-tests-student)

Проект содержит API-тесты для сервиса объявлений TestBoard. Тесты покрывают два блока функциональности:

1. **Создание объявления** - позитивные и негативные сценарии, проверка через `GET /advertisement`, поиск и фотографии.
2. **Работа с избранным**- полный цикл (добавление, проверка, удаление), бизнес-ограничения и негативные сценарии.

Чтобы тесты запускались нужно в  файл .env.override вставить переменные окружения. 

```bash
# Пользователь 1 (основной) 
# Владелец избранного. От его имени выполняются тесты Части 1 и добавление в избранное.
TEST_LOGIN=your_login1
TEST_PASSWORD=your_password_1

# Пользователь 2 (вспомогательный)
# Нужен ТОЛЬКО для тестов избранного (Часть 2).
# Владелец «чужого» объявления, которое Пользователь 1 добавляет в избранное.
TEST_LOGIN_2=your_login2
TEST_PASSWORD_2=your_password_2
```
# Зачем нужен второй пользователь?
В API TestBoard действуют строгие бизнес-правила:

Нельзя добавить в избранное своё собственное объявление-сервер вернёт 403 Forbidden.
Удалить объявление может только его владелец-иначе 403 Forbidden.
Поэтому позитивный сценарий избранного физически невозможно пройти с одним аккаунтом. 

# Как зарегистрировать второго пользователя
1. Откройте TestBoard в браузере.
2. Выйдите из текущего аккаунта.
3. Нажмите «Регистрация» и создайте нового пользователя (например, testuser2@example.com).
4. Скопируйте его login и пароль в .env.override в переменные TEST_LOGIN_2 и TEST_PASSWORD_2.

# Тестовые данные
Для тестов создания объявления нужен файл-изображение. Он уже лежит в репозитории:
```bash
testdata/test_image.jpg
```
# Запуск тестов 
Все тесты
```bash
go test -v ./...
```
Только Часть 1 Создание объявлений
```bash
go test -v ./tests/scenarios/advertisements/
```
Только Часть 2 Избранное
```bash
go test -v ./tests/scenarios/favorites/
```
