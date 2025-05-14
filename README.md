# Backend для приложения T-Prep
Frontend приложения находится в [другом репозитории](https://github.com/s0roh/T-Prep)

## Архитектура слоев проекта
- **Router**: роутинг
- **Controller**: бизнес-логика
- **Usecase**: основные сценарии пользователя
- **Repository и Storage**: взаимодействие с хранилищами (MongoDB и MinIO)
- **Domain**: все основные структуры и интерфейсы
  
![image](https://github.com/user-attachments/assets/be33c335-7afa-42f6-9772-187d0a3b9d53)

## Технологии
- **Golang**: основной язык программирования
- **MongoDB**: база данных
- **MinIO**: s3 для медиа
- **Docker**: контейнеризация
- **Prometheus**: скрапинг метрик
- **Grafana**: визуализация полученных метрик

![image](https://github.com/user-attachments/assets/813dca87-bbe6-4e0a-ab4d-363696f60795)

## Основные библиотеки
- **chi**: роутинг
- **mongo go driver**: официальный драйвер MongoDB
- **jwt**: авторизация
- **viper**: загрузка конфига из `.env` файла
- **easyjson**: быстрая работа с json
- **bcrypt**: хеширование паролей
- **testing**: юнит-тестирование
- **mockery**: генерация моков
- **minio**: официальный драйвер MinIO
- **prometheus**: скрапинг метрик
- Увидеть больше можно в `go.mod`
