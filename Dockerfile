FROM golang:latest

WORKDIR /test_ozon

COPY go.mod go.sum ./
COPY main.go .
COPY handlers.go .
COPY database.go .
COPY make.go .


# Сборка приложения
RUN go build -o app

# Установка переменной окружения для указания порта, на котором будет работать приложение
ENV PORT=8080

# Установка переменной окружения для указания типа хранилища
ENV STORAGE=postgres

# Установка переменной окружения для указания адреса in-memory или Postgres
ENV DB_ADDRESS=localhost

# Установка переменной окружения для указания пароля Postgres
ENV DB_PASSWORD=default

# Установка переменной окружения для указания имени базы данных Postgres
ENV DB_NAME=urls

# Открытие порта в контейнере
EXPOSE 8080

# Запуск приложения при старте контейнера
CMD ["./app"]