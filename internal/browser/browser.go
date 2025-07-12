package browser

import (
	"log"
	"net/http"
	"sync"

	"github.com/baneronetwo/gluglu/internal/browser/html"
	"github.com/baneronetwo/gluglu/internal/browser/js"
	"github.com/baneronetwo/gluglu/internal/browser/network"
	"github.com/baneronetwo/gluglu/internal/browser/renderer"
)

// Browser представляет основной движок браузера
type Browser struct {
	networkManager *network.Manager
	htmlParser     *html.Parser
	jsEngine       *js.Engine
	renderer       *renderer.Renderer
	currentURL     string
	currentPage    *Page
	history        []string
	historyPos     int
	mutex          sync.Mutex
}

// Page представляет загруженную страницу
type Page struct {
	URL              string
	Title            string
	Content          string
	DOM              *html.Document
	RenderedDocument *renderer.Document
}

// NewBrowser создает новый экземпляр браузера
func NewBrowser() *Browser {
	log.Println("Инициализация компонентов браузера...")
	
	// Создание компонентов
	networkMgr := network.NewManager()
	htmlParser := html.NewParser()
	jsEngine := js.NewEngine()
	renderer := renderer.NewRenderer()
	
	return &Browser{
		networkManager: networkMgr,
		htmlParser:     htmlParser,
		jsEngine:       jsEngine,
		renderer:       renderer,
		history:        make([]string, 0),
		historyPos:     -1,
	}
}

// LoadURL загружает указанный URL
func (b *Browser) LoadURL(url string) error {
	log.Printf("Загрузка URL: %s", url)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	// Получение содержимого страницы через сетевой модуль
	content, err := b.networkManager.Fetch(url)
	if err != nil {
		log.Printf("Ошибка загрузки URL %s: %v", url, err)
		return err
	}
	
	// Парсинг HTML
	doc, err := b.htmlParser.Parse(content)
	if err != nil {
		log.Printf("Ошибка парсинга HTML: %v", err)
		return err
	}
	
	// Выполнение JavaScript
	b.jsEngine.Execute(doc)
	
	// Рендеринг страницы
	renderedDoc := b.renderer.Render(doc)
	
	// Создание объекта страницы
	page := &Page{
		URL:              url,
		Title:            doc.Title,
		Content:          content,
		DOM:              doc,
		RenderedDocument: renderedDoc,
	}
	
	// Обновление текущей страницы и истории
	b.currentURL = url
	b.currentPage = page
	
	// Добавление URL в историю
	if b.historyPos < len(b.history)-1 {
		// Если мы перешли назад и затем на новую страницу, обрезаем историю
		b.history = b.history[:b.historyPos+1]
	}
	b.history = append(b.history, url)
	b.historyPos = len(b.history) - 1
	
	return nil
}

// GetCurrentPage возвращает текущую страницу
func (b *Browser) GetCurrentPage() *Page {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.currentPage
}

// GoBack переходит на предыдущую страницу в истории
func (b *Browser) GoBack() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	if b.historyPos <= 0 {
		return http.ErrNoLocation
	}
	
	b.historyPos--
	url := b.history[b.historyPos]
	
	// Разблокируем мьютекс перед вызовом LoadURL, который также блокирует мьютекс
	b.mutex.Unlock()
	err := b.LoadURL(url)
	b.mutex.Lock() // Восстанавливаем блокировку
	
	return err
}

// GoForward переходит на следующую страницу в истории
func (b *Browser) GoForward() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	if b.historyPos >= len(b.history)-1 {
		return http.ErrNoLocation
	}
	
	b.historyPos++
	url := b.history[b.historyPos]
	
	// Разблокируем мьютекс перед вызовом LoadURL, который также блокирует мьютекс
	b.mutex.Unlock()
	err := b.LoadURL(url)
	b.mutex.Lock() // Восстанавливаем блокировку
	
	return err
}

// CanGoBack проверяет, можно ли перейти назад в истории
func (b *Browser) CanGoBack() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.historyPos > 0
}

// CanGoForward проверяет, можно ли перейти вперед в истории
func (b *Browser) CanGoForward() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.historyPos < len(b.history)-1
}