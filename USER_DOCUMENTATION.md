# 📊 Fin-Core - Руководство пользователя

## 🎯 Что такое Fin-Core?

**Fin-Core** — это система управления личными финансами, которая помогает вам:

- 💰 **Отслеживать доходы и расходы** по разным банковским счетам
- 🏦 **Управлять несколькими счетами** в разных валютах
- 📊 **Планировать бюджеты** и контролировать траты
- 🔔 **Получать уведомления** о превышении бюджета и низком балансе
- 📈 **Анализировать финансовые привычки** через категории

---

## 🚀 Начало работы

### Первый запуск

1. **Запустите сервер:**
   ```bash
   go run cmd/server/main.go
   ```

2. **Проверьте работу:**
   - Откройте браузер: `http://localhost:8080/api/v1/health`
   - Должен появиться ответ: `{"status": "OK", "service": "fin-core"}`

3. **Посмотрите документацию API:**
   - Swagger UI: `http://localhost:8080/swagger/index.html`

### Настройка окружения

Создайте файл `.env` в корне проекта:

```env
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fin_db
DB_SSL_MODE=disable

# Аутентификация
AUTH_SERVICE_URL=http://localhost:3000

# RabbitMQ (опционально)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Порт сервера
PORT=8080
```

---

## 👤 Создание аккаунта

### Шаг 1: Создайте финансовый аккаунт

```http
POST /api/v1/account
Content-Type: application/json

{
  "display_name": "Мой финансовый аккаунт"
}
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": "user123",
    "name": "Мой финансовый аккаунт",
    "display_name": "Мой финансовый аккаунт",
    "timezone": "UTC",
    "base_currency": "KZT",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Шаг 2: Получите информацию об аккаунте

```http
GET /api/v1/account
```

---

## 🏦 Управление банковскими счетами

### Создание банковского счета

```http
POST /api/v1/bankAccounts
Content-Type: application/json

{
  "name": "Каспи рубли",
  "currency": "KZT",
  "account_type": "debit",
  "bank_name": "Kaspi"
}
```

**Типы счетов:**
- `cash` - наличные деньги
- `debit` - дебетовая карта
- `credit` - кредитная карта  
- `savings` - сберегательный счет

**Поддерживаемые валюты:**
- `KZT` - тенге
- `USD` - доллар США
- `EUR` - евро
- `RUB` - российский рубль

### Просмотр всех счетов

```http
GET /api/v1/bankAccounts
```

### Управление счетом

**Получить конкретный счет:**
```http
GET /api/v1/bankAccounts/{bank_account_id}
```

**Деактивировать счет:**
```http
PUT /api/v1/bankAccounts/{bank_account_id}/deactivate
```

**Активировать счет:**
```http
PUT /api/v1/bankAccounts/{bank_account_id}/activate
```

**Удалить счет:**
```http
DELETE /api/v1/bankAccounts/{bank_account_id}
```

---

## 💳 Работа с транзакциями

### Создание транзакции

**Доход:**
```http
POST /api/v1/transactions
Content-Type: application/json

{
  "bank_account_id": 1,
  "category_id": 1,
  "amount": 150000.00,
  "description": "Зарплата",
  "transaction_type": "income"
}
```

**Расход:**
```http
POST /api/v1/transactions
Content-Type: application/json

{
  "bank_account_id": 1,
  "category_id": 2,
  "amount": 5000.00,
  "description": "Покупка продуктов",
  "transaction_type": "expense"
}
```

**Важно:** Для расходов указывайте положительную сумму, система автоматически сделает её отрицательной.

### Перевод между счетами

```http
POST /api/v1/transfer
Content-Type: application/json

{
  "from_bank_account_id": 1,
  "to_bank_account_id": 2,
  "amount": 10000.00,
  "description": "Перевод на сберегательный счет"
}
```

### Просмотр транзакций

**Все транзакции:**
```http
GET /api/v1/transactions
```

**Конкретная транзакция:**
```http
GET /api/v1/transactions/{id}
```

**Транзакции по категории:**
```http
GET /api/v1/transactions/by-category/{category_id}
```

**Баланс счета:**
```http
GET /api/v1/bank_accounts/{account_id}/balance
```

---

## 📂 Система категорий

### Создание категории

```http
POST /api/v1/categories
Content-Type: application/json

{
  "name": "Продукты",
  "type": "expense",
  "color": "#FF5722",
  "icon": "shopping_cart"
}
```

**Типы категорий:**
- `income` - для доходов
- `expense` - для расходов

### Управление категориями

**Все категории:**
```http
GET /api/v1/categories
```

**Конкретная категория:**
```http
GET /api/v1/categories/{category_id}
```

**Удалить категорию:**
```http
DELETE /api/v1/categories/{category_id}
```

---

## 💰 Система бюджетов

### Создание бюджета

```http
POST /api/v1/budgets
Content-Type: application/json

{
  "budget_limit_name": "Бюджет на продукты",
  "category_id": 1,
  "amount": 50000.00,
  "period": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

**Периоды бюджета:**
- `monthly` - месячный
- `weekly` - недельный
- `yearly` - годовой

### Просмотр бюджетов

**Все бюджеты:**
```http
GET /api/v1/budgets
```

**Статус бюджета по категории:**
```http
GET /api/v1/budgets/{category_id}/status
```

**Сводка по бюджетам:**
```http
GET /api/v1/budgets/summary
```

### Как работают бюджеты

1. **Создайте бюджет** для категории (например, "Продукты")
2. **Укажите лимит** (например, 50,000 тенге в месяц)
3. **Система автоматически отслеживает** траты по этой категории
4. **При достижении 80%** лимита приходит предупреждение
5. **При превышении** лимита приходит уведомление

---

## 🔔 Система уведомлений

### Типы уведомлений

1. **Превышение бюджета** - когда траты превысили установленный лимит
2. **Предупреждение о бюджете** - когда потрачено 80% от лимита
3. **Низкий баланс** - когда на счету мало денег

### Просмотр уведомлений

**Все уведомления:**
```http
GET /api/v1/notification?limit=20&offset=0
```

**Отметить как прочитанное:**
```http
PUT /api/v1/notification/{id}/read
```

**Отметить все как прочитанные:**
```http
PUT /api/v1/notification/read-all
```

### Настройки уведомлений

**Получить настройки:**
```http
GET /api/v1/notification/settings
```

**Сохранить настройки:**
```http
PUT /api/v1/notification/settings
Content-Type: application/json

{
  "budget_alerts_enabled": true,
  "balance_alerts_enabled": true,
  "budget_warning_percent": 80,
  "low_balance_threshold": 10000.00,
  "preferred_channel": "email"
}
```

**Параметры настроек:**
- `budget_alerts_enabled` - включить уведомления о бюджете
- `balance_alerts_enabled` - включить уведомления о балансе
- `budget_warning_percent` - процент для предупреждения (0-100)
- `low_balance_threshold` - минимальный баланс для уведомления
- `preferred_channel` - способ уведомления: `email`, `push`, `sms`

---

## 📈 Примеры использования

### Сценарий 1: Настройка личного бюджета

1. **Создайте аккаунт:**
   ```http
   POST /api/v1/account
   {"display_name": "Мой бюджет"}
   ```

2. **Добавьте банковский счет:**
   ```http
   POST /api/v1/bankAccounts
   {
     "name": "Основная карта",
     "currency": "KZT",
     "account_type": "debit",
     "bank_name": "Kaspi"
   }
   ```

3. **Создайте категории:**
   ```http
   POST /api/v1/categories
   {"name": "Продукты", "type": "expense", "color": "#FF5722", "icon": "shopping_cart"}
   
   POST /api/v1/categories
   {"name": "Транспорт", "type": "expense", "color": "#2196F3", "icon": "directions_car"}
   ```

4. **Установите бюджеты:**
   ```http
   POST /api/v1/budgets
   {
     "budget_limit_name": "Продукты на месяц",
     "category_id": 1,
     "amount": 50000.00,
     "period": "monthly",
     "start_date": "2024-01-01",
     "end_date": "2024-01-31"
   }
   ```

5. **Настройте уведомления:**
   ```http
   PUT /api/v1/notification/settings
   {
     "budget_alerts_enabled": true,
     "budget_warning_percent": 80,
     "low_balance_threshold": 5000.00,
     "preferred_channel": "email"
   }
   ```

### Сценарий 2: Отслеживание расходов

1. **Добавьте расход:**
   ```http
   POST /api/v1/transactions
   {
     "bank_account_id": 1,
     "category_id": 1,
     "amount": 2500.00,
     "description": "Покупка хлеба и молока",
     "transaction_type": "expense"
   }
   ```

2. **Проверьте статус бюджета:**
   ```http
   GET /api/v1/budgets/1/status
   ```

3. **Посмотрите уведомления:**
   ```http
   GET /api/v1/notification
   ```

### Сценарий 3: Перевод между счетами

1. **Создайте второй счет:**
   ```http
   POST /api/v1/bankAccounts
   {
     "name": "Сбережения",
     "currency": "KZT",
     "account_type": "savings",
     "bank_name": "Halyk"
   }
   ```

2. **Переведите деньги:**
   ```http
   POST /api/v1/transfer
   {
     "from_bank_account_id": 1,
     "to_bank_account_id": 2,
     "amount": 20000.00,
     "description": "Откладываю на отпуск"
   }
   ```

---

## ❓ Часто задаваемые вопросы

### Q: Как добавить начальный баланс на счет?
A: Создайте транзакцию типа "income" с описанием "Начальный баланс".

### Q: Можно ли изменить категорию транзакции?
A: В текущей версии нельзя. Создайте новую транзакцию с правильной категорией.

### Q: Как работает система уведомлений?
A: Система автоматически отслеживает траты и отправляет уведомления при превышении лимитов или достижении пороговых значений.

### Q: Можно ли иметь несколько аккаунтов?
A: Нет, каждый пользователь может иметь только один финансовый аккаунт, но с несколькими банковскими счетами.

### Q: Поддерживаются ли валютные переводы?
A: Да, система поддерживает поле `transfer_rate` для конвертации валют при переводах.

### Q: Как удалить все данные?
A: Удалите банковские счета (это удалит все транзакции), затем удалите аккаунт.

---

## 🛠️ Техническая поддержка

### Проверка состояния системы

```http
GET /api/v1/health
```

### Логи и отладка

- Проверьте логи сервера для диагностики проблем
- Убедитесь, что база данных PostgreSQL запущена
- Проверьте подключение к RabbitMQ (если используется)

### Контакты

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **API документация:** `http://localhost:8080/swagger/doc.json`

---

## 📝 Заключение

Fin-Core поможет вам:

✅ **Контролировать финансы** - отслеживайте доходы и расходы  
✅ **Планировать бюджет** - устанавливайте лимиты и следите за тратами  
✅ **Получать уведомления** - не пропускайте важные события  
✅ **Анализировать траты** - понимайте, куда уходят деньги  
✅ **Управлять счетами** - работайте с несколькими банковскими счетами  

Начните с создания аккаунта и добавления первого банковского счета. Затем создайте категории и установите бюджеты. Система будет автоматически отслеживать ваши финансы и уведомлять о важных событиях!

---

