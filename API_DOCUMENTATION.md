# Fin-Core API Documentation

## 🚀 **Базовая информация**

- **Base URL**: `http://localhost:8080/api/v1`
- **Аутентификация**: Cookie `token` с JWT токеном
- **Content-Type**: `application/json`

## 📋 **Swagger UI**

Интерактивная документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## 🔐 **Аутентификация**

Все защищенные эндпоинты требуют токен в cookie:
```javascript
// Пример запроса с токеном
fetch('http://localhost:8080/api/v1/account', {
  method: 'GET',
  credentials: 'include', // Важно для отправки cookies
  headers: {
    'Content-Type': 'application/json'
  }
})
```

## 📊 **Основные эндпоинты**

### **Аккаунты**

#### Создать аккаунт
```http
POST /api/v1/account
Content-Type: application/json

{
  "display_name": "Мой финансовый аккаунт"
}
```

#### Получить аккаунт
```http
GET /api/v1/account
```

### **Банковские счета**

#### Получить все банковские счета
```http
GET /api/v1/bankAccounts
```

#### Создать банковский счет
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

#### Получить конкретный банковский счет
```http
GET /api/v1/bankAccounts/{bank_account_id}
```

#### Удалить банковский счет
```http
DELETE /api/v1/bankAccounts/{bank_account_id}
```

#### Деактивировать банковский счет
```http
PUT /api/v1/bankAccounts/{bank_account_id}/deactivate
```

#### Активировать банковский счет
```http
PUT /api/v1/bankAccounts/{bank_account_id}/activate
```

### **Транзакции**

#### Создать транзакцию
```http
POST /api/v1/transactions
Content-Type: application/json

{
  "bank_account_id": 1,
  "category_id": 1,
  "amount": 5000.00,
  "description": "Покупка продуктов",
  "transaction_type": "expense"
}
```

#### Получить все транзакции
```http
GET /api/v1/transactions
```

#### Получить конкретную транзакцию
```http
GET /api/v1/transactions/{id}
```

#### Перевод между счетами
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

#### История транзакций по аккаунту
```http
GET /api/v1/account/{account_id}/transactions
```

#### Баланс банковского счета
```http
GET /api/v1/bank_accounts/{account_id}/balance
```

### **Категории**

#### Создать категорию
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

#### Получить все категории
```http
GET /api/v1/categories
```

#### Получить конкретную категорию
```http
GET /api/v1/categories/{category_id}
```

#### Удалить категорию
```http
DELETE /api/v1/categories/{category_id}
```

## 📝 **Типы данных**

### **Типы транзакций**
- `income` - доход
- `expense` - расход  
- `transfer` - перевод

### **Типы банковских счетов**
- `cash` - наличные
- `debit` - дебетовая карта
- `credit` - кредитная карта
- `savings` - сберегательный счет

### **Валюты**
- `KZT` - тенге
- `USD` - доллар США
- `EUR` - евро
- `RUB` - российский рубль

### **Типы категорий**
- `income` - для доходов
- `expense` - для расходов

## ⚠️ **Коды ошибок**

- `400` - Некорректные данные запроса
- `401` - Не авторизован (нет токена или токен недействителен)
- `403` - Доступ запрещен
- `404` - Ресурс не найден
- `409` - Конфликт (например, попытка удалить счет с транзакциями)
- `500` - Внутренняя ошибка сервера

## 🔄 **Примеры ответов**

### Успешный ответ
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Каспи рубли",
    "currency": "KZT",
    "balance": 150000.00
  }
}
```

### Ошибка
```json
{
  "success": false,
  "error": "invalid bank account id",
  "details": "bank account not found"
}
```

## 🛠️ **Для разработки**

### Postman Collection
Можете импортировать Swagger JSON в Postman:
```
http://localhost:8080/swagger/doc.json
```

### Тестирование
1. Запустите сервер: `go run cmd/server/main.go`
2. Откройте Swagger UI: `http://localhost:8080/swagger/index.html`
3. Получите токен из auth-сервиса
4. Используйте токен в cookie для запросов
