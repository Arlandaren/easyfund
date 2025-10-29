package banking

// ClientFactory - фабрика для создания клиентов банков
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

// GetClient - получить клиент в зависимости от режима
func (f *ClientFactory) GetClient(useMock bool) BankingClient {
	if useMock {
		return f.mockClient
	}
	return f.vbankClient
}

// GetVBankClient - получить реальный VBank клиент
func (f *ClientFactory) GetVBankClient() *VBankClient {
	return f.vbankClient
}

// GetMockClient - получить мок клиент
func (f *ClientFactory) GetMockClient() *MockClient {
	return f.mockClient
}
