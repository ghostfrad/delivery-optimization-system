# Delivery Optimization System

Система для управления заказами и оптимизации маршрутов доставки.

## Архитектура

- **Order Service** — управление заказами (создание, статусы, история)
- **Pool Service** — управление пулом заказов и назначение курьерам
- **Route Service** — построение оптимальных маршрутов (TSP + GraphHopper)
- **Notification Service** — push-уведомления
- **API Gateway** — единая точка входа
- 
## Быстрый старт

```bash
# Клонировать репозиторий
git clone https://github.com/ghostfrad/delivery-optimization-system.git
cd delivery-optimization-system/services/order-service

# Установить зависимости
make deps

# Запустить сервис (с моками)
В файле config.yaml выставить параметр use_mocks: true

make run

# Запустить сервис с реальным postgres и kafka
В файле конфига выставить параметр use_mocks: false 

Поднять контейнеры используя файл docker-compose

make run

# Запустить тесты
make test
