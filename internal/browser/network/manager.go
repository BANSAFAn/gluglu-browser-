package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Manager управляет сетевыми запросами
type Manager struct {
	client *http.Client
	cache  map[string]CacheEntry
}

// CacheEntry представляет кэшированный ответ
type CacheEntry struct {
	Content   string
	Timestamp time.Time
	Headers   http.Header
}

// NewManager создает новый менеджер сетевых запросов
func NewManager() *Manager {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Manager{
		client: client,
		cache:  make(map[string]CacheEntry),
	}
}

// Fetch загружает содержимое по указанному URL
func (m *Manager) Fetch(url string) (string, error) {
	log.Printf("Сетевой запрос: %s", url)
	
	// Проверка кэша
	if entry, ok := m.cache[url]; ok {
		// Проверяем, не устарел ли кэш (простая реализация - 5 минут)
		if time.Since(entry.Timestamp) < 5*time.Minute {
			log.Printf("Использование кэшированного ответа для %s", url)
			return entry.Content, nil
		}
	}
	
	// Создание запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	
	// Установка заголовков
	req.Header.Set("User-Agent", "GluGlu Browser/0.1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	
	// Выполнение запроса
	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()
	
	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неверный статус ответа: %d %s", resp.StatusCode, resp.Status)
	}
	
	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	
	content := string(body)
	
	// Кэширование ответа
	m.cache[url] = CacheEntry{
		Content:   content,
		Timestamp: time.Now(),
		Headers:   resp.Header,
	}
	
	return content, nil
}

// ClearCache очищает кэш
func (m *Manager) ClearCache() {
	m.cache = make(map[string]CacheEntry)
}