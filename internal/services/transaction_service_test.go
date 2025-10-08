package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"justTest/internal/models"
)

// ============= МОКИ =============

// MockTransactionRepository - мок для репозитория транзакций
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) (*models.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByBankAccountID(bankAccountID int64, limit, offset int) ([]*models.Transaction, error) {
	args := m.Called(bankAccountID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(accountID int64, limit, offset int) ([]*models.Transaction, error) {
	args := m.Called(accountID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByTransactionID(transactionID int64) (*models.Transaction, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTotalAmountByBankAccountID(bankAccountID int64) (float64, error) {
	args := m.Called(bankAccountID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTransactionRepository) GetByCategoryID(categoryID int64) ([]*models.Transaction, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransfersByAccountID(accountID int64) ([]*models.Transaction, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(categoryID *int64, fromAccountID *int64, toAccountCategoryID *int64, toAccountID int64, amount float64, description string, transferRate *float64) error {
	args := m.Called(categoryID, fromAccountID, toAccountCategoryID, toAccountID, amount, description, transferRate)
	return args.Error(0)
}

// MockBankAccountRepository - мок для репозитория банковских счетов
type MockBankAccountRepository struct {
	mock.Mock
}

func (m *MockBankAccountRepository) GetByBankAccountID(bankAccountID int64) (*models.BankAccount, error) {
	args := m.Called(bankAccountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) Create(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	args := m.Called(bankAccount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetByAccountID(accountID int64) ([]*models.BankAccount, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetActiveBankAccounts(accountID int64) ([]*models.BankAccount, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetBankAccountByCurrency(accountID int64, currency string) ([]*models.BankAccount, error) {
	args := m.Called(accountID, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) ExsitsAccountIDAndName(accountID int64, name string) (bool, error) {
	args := m.Called(accountID, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockBankAccountRepository) DeActiveBankAccount(bankAccountID int64) error {
	args := m.Called(bankAccountID)
	return args.Error(0)
}

func (m *MockBankAccountRepository) ActivateBankAccount(bankAccountID int64) error {
	args := m.Called(bankAccountID)
	return args.Error(0)
}

func (m *MockBankAccountRepository) DeleteBankAccount(bankAccountID int64) error {
	args := m.Called(bankAccountID)
	return args.Error(0)
}

// MockCategoryRepository - мок для репозитория категорий
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategory(category *models.Category) (*models.Category, error) {
	args := m.Called(category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) DeleteCategory(categoryID int64) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByAccountID(accountID int64) ([]*models.Category, error) {
	args := m.Called(accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(categoryID int64) (*models.Category, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

// MockAccountRepository - мок для репозитория аккаунтов
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account *models.Account) (*models.Account, error) {
	args := m.Called(account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(userID string) (*models.Account, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByID(id int64) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

// ============= ТЕСТЫ =============

// TestTransactionService_CreateTransaction_Success - тест успешного создания транзакции
func TestTransactionService_CreateTransaction_Success(t *testing.T) {
	// Arrange - подготовка данных
	mockTransactionRepository := new(MockTransactionRepository)
	mockBankAccountRepository := new(MockBankAccountRepository)
	mockCategoryRepository := new(MockCategoryRepository)
	mockAccountRepository := new(MockAccountRepository)

	// Тестовые данные
	userID := "user123"
	bankAccountID := int64(1)
	amount := 100.0
	description := "Test transaction"
	transactionType := "income"

	// Мок-объекты для ответов репозиториев
	mockAccount := &models.Account{
		ID:     10,
		UserID: userID,
	}

	mockBankAccount := &models.BankAccount{
		ID:        bankAccountID,
		AccountID: 10,
	}

	expectedTransaction := &models.Transaction{
		ID:              1,
		BankAccountID:   bankAccountID,
		Amount:          amount,
		Description:     description,
		TransactionType: transactionType,
	}

	// Настройка ожиданий для моков
	mockBankAccountRepository.On("GetByBankAccountID", bankAccountID).Return(mockBankAccount, nil)
	mockAccountRepository.On("GetByUserID", userID).Return(mockAccount, nil)
	mockTransactionRepository.On("Create", mock.AnythingOfType("*models.Transaction")).Return(expectedTransaction, nil)

	// Создание сервиса с моками
	service := NewTransactionService(mockTransactionRepository, mockBankAccountRepository, mockCategoryRepository, mockAccountRepository)

	// Act - выполнение действия
	transaction, err := service.CreateTransaction(userID, bankAccountID, amount, description, nil, transactionType)

	// Assert - проверка результата
	assert.NoError(t, err, "Ошибка не должна была произойти")
	assert.NotNil(t, transaction, "Транзакция не должна быть nil")
	assert.Equal(t, expectedTransaction.ID, transaction.ID, "ID транзакции должен совпадать")
	assert.Equal(t, amount, transaction.Amount, "Сумма должна совпадать")
	assert.Equal(t, description, transaction.Description, "Описание должно совпадать")
	assert.Equal(t, transactionType, transaction.TransactionType, "Тип транзакции должен совпадать")

	// Проверяем, что все ожидания были выполнены
	mockTransactionRepository.AssertExpectations(t)
	mockBankAccountRepository.AssertExpectations(t)
	mockAccountRepository.AssertExpectations(t)
}

// TestTransactionService_CreateTransaction_InvalidUserID - тест с невалидным userID
func TestTransactionService_CreateTransaction_InvalidUserID(t *testing.T) {
	// Arrange
	mockTransactionRepository := new(MockTransactionRepository)
	mockBankAccountRepository := new(MockBankAccountRepository)
	mockCategoryRepository := new(MockCategoryRepository)
	mockAccountRepository := new(MockAccountRepository)

	service := NewTransactionService(mockTransactionRepository, mockBankAccountRepository, mockCategoryRepository, mockAccountRepository)

	// Act
	transaction, err := service.CreateTransaction("", 1, 100.0, "Test", nil, "income")

	// Assert
	assert.Error(t, err, "Должна быть ошибка")
	assert.Nil(t, transaction, "Транзакция должна быть nil")
	assert.Contains(t, err.Error(), "invalid user id", "Сообщение об ошибке должно содержать 'invalid user id'")
}

// TestTransactionService_CreateTransaction_InvalidAmount - тест с невалидной суммой
func TestTransactionService_CreateTransaction_InvalidAmount(t *testing.T) {
	// Arrange
	mockTransactionRepository := new(MockTransactionRepository)
	mockBankAccountRepository := new(MockBankAccountRepository)
	mockCategoryRepository := new(MockCategoryRepository)
	mockAccountRepository := new(MockAccountRepository)

	service := NewTransactionService(mockTransactionRepository, mockBankAccountRepository, mockCategoryRepository, mockAccountRepository)

	// Act
	transaction, err := service.CreateTransaction("user123", 1, 0, "Test", nil, "income")

	// Assert
	assert.Error(t, err, "Должна быть ошибка")
	assert.Nil(t, transaction, "Транзакция должна быть nil")
	assert.Contains(t, err.Error(), "invalid amount", "Сообщение об ошибке должно содержать 'invalid amount'")
}

// TestTransactionService_GetTransactionHistory_Success - тест успешного получения истории транзакций
func TestTransactionService_GetTransactionHistory_Success(t *testing.T) {
	// Arrange
	mockTransactionRepository := new(MockTransactionRepository)
	mockBankAccountRepository := new(MockBankAccountRepository)
	mockCategoryRepository := new(MockCategoryRepository)
	mockAccountRepository := new(MockAccountRepository)

	userID := "user123"
	bankAccountID := int64(1)

	mockAccount := &models.Account{
		ID:     10,
		UserID: userID,
	}

	mockBankAccount := &models.BankAccount{
		ID:        bankAccountID,
		AccountID: 10,
	}

	expectedTransactions := []*models.Transaction{
		{
			ID:              1,
			BankAccountID:   bankAccountID,
			Amount:          100.0,
			Description:     "Test transaction 1",
			TransactionType: "income",
		},
		{
			ID:              2,
			BankAccountID:   bankAccountID,
			Amount:          -50.0,
			Description:     "Test transaction 2",
			TransactionType: "expense",
		},
	}

	// Настройка ожиданий
	mockBankAccountRepository.On("GetByBankAccountID", bankAccountID).Return(mockBankAccount, nil)
	mockAccountRepository.On("GetByUserID", userID).Return(mockAccount, nil)
	mockTransactionRepository.On("GetByBankAccountID", bankAccountID, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedTransactions, nil)

	service := NewTransactionService(mockTransactionRepository, mockBankAccountRepository, mockCategoryRepository, mockAccountRepository)

	// Act
	transactions, err := service.GetTransactionHistory(userID, bankAccountID)

	// Assert
	assert.NoError(t, err, "Ошибка не должна была произойти")
	assert.NotNil(t, transactions, "Список транзакций не должен быть nil")
	assert.Len(t, transactions, 2, "Должно быть 2 транзакции")
	assert.Equal(t, expectedTransactions[0].ID, transactions[0].ID, "ID первой транзакции должен совпадать")
	assert.Equal(t, expectedTransactions[1].ID, transactions[1].ID, "ID второй транзакции должен совпадать")

	mockTransactionRepository.AssertExpectations(t)
	mockBankAccountRepository.AssertExpectations(t)
	mockAccountRepository.AssertExpectations(t)
}
