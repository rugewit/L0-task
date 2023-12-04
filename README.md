## Задание

Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. Модель данных в формате JSON прилагается к заданию.	

Что нужно сделать:

1. Развернуть локально PostgreSQL
   1. Создать свою БД
   2. Настроить своего пользователя
   3. Создать таблицы для хранения полученных данных
2. Разработать сервис
   1. Реализовать подключение и подписку на канал в nats-streaming
   2. Полученные данные записывать в БД
   3. Реализовать кэширование полученных данных в сервисе (сохранять in memory)
   4. В случае падения сервиса необходимо восстанавливать кэш из БД
   5. Запустить http-сервер и выдавать данные по id из кэша
3. Разработать простейший интерфейс отображения полученных данных по id заказа

## Как запустить

1. Запустить docker compose: 
   1. sudo docker compose up

2. Запустить subscriber: 
   1. cd subscriber
   2. go run cmd/main.go
3. Запустить pusblisher:
   1. cd publisher
   2. go run cmd/main.go
4. Запустить react frontend
   1. cd frontend/l0-frontend
   2. npm install && npm start

## Пример работы

<image src="https://github.com/rugewit/L0-task/blob/main/github_images/1.png" alt="">

<image src="https://github.com/rugewit/L0-task/blob/main/github_images/2.png" alt="">

<image src="https://github.com/rugewit/L0-task/blob/main/github_images/3.png" alt="">