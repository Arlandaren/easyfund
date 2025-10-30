package banking

// ClientFactory — фабрика создания клиентов Virtual Bank и мок-клиента.
// Здесь намеренно не используется интерфейс, чтобы избежать несовпадений методов:
// возвращаем конкретные типы, которые реально существуют.
type ClientFactory struct {
    vbankClient *VBankClient
    mockClient  *MockClient
}

func NewClientFactory(baseURL, clientID, clientSecret string) *ClientFactory {
    return &ClientFactory{
        vbankClient: NewVBankClient(baseURL, clientID, clientSecret),
        mockClient:  NewMockClient(),
    }
}

// UseMock позволяет настраивать использование мока (например, из конфига)
func (f *ClientFactory) UseMock() bool {
    return false // при необходимости поменяйте на флаг из вашего конфига
}

// GetVBankClient — получить реальный VBank клиент
func (f *ClientFactory) GetVBankClient() *VBankClient {
    return f.vbankClient
}

// GetMockClient — получить мок клиент
func (f *ClientFactory) GetMockClient() *MockClient {
    return f.mockClient
}

// GetClient — выбрать клиента в рантайме.
// Возвращает конкретный тип (любой) через empty interface, чтобы можно было
// подставить туда, где ожидается конкретная реализация.
func (f *ClientFactory) GetClient(useMock bool) interface{} {
    if useMock {
        return f.mockClient
    }
    return f.vbankClient
}
