package html

import (
	"fmt"
	"log"
	"strings"
)

// Parser представляет HTML парсер
type Parser struct {}

// Document представляет DOM-дерево HTML документа
type Document struct {
	Title    string
	Body     string
	Elements []Element
	Head     *Element
	BodyElem *Element
}

// Element представляет HTML элемент
type Element struct {
	TagName    string
	Attributes map[string]string
	Children   []Element
	Parent     *Element
	Text       string
	ID         string
	ClassNames []string
}

// NewParser создает новый HTML парсер
func NewParser() *Parser {
	return &Parser{}
}

// Parse разбирает HTML строку и возвращает DOM-дерево
func (p *Parser) Parse(htmlContent string) (*Document, error) {
	log.Println("Парсинг HTML документа...")
	
	// Создаем документ
	doc := &Document{
		Elements: make([]Element, 0),
	}
	
	// Простой парсинг заголовка
	titleStart := strings.Index(htmlContent, "<title>")
	titleEnd := strings.Index(htmlContent, "</title>")
	if titleStart >= 0 && titleEnd > titleStart {
		doc.Title = htmlContent[titleStart+7 : titleEnd]
	} else {
		doc.Title = "Без заголовка"
	}
	
	// Простой парсинг тела
	bodyStart := strings.Index(htmlContent, "<body")
	bodyEnd := strings.Index(htmlContent, "</body>")
	if bodyStart >= 0 && bodyEnd > bodyStart {
		// Находим закрывающую скобку открывающего тега body
		bodyTagEnd := strings.Index(htmlContent[bodyStart:], ">")
		if bodyTagEnd > 0 {
			doc.Body = htmlContent[bodyStart+bodyTagEnd+1 : bodyEnd]
		}
	} else {
		// Если нет тегов body, используем весь документ
		doc.Body = htmlContent
	}
	
	// Создаем корневые элементы
	html := Element{
		TagName:    "html",
		Attributes: make(map[string]string),
		Children:   make([]Element, 0),
	}
	
	head := Element{
		TagName:    "head",
		Attributes: make(map[string]string),
		Children:   make([]Element, 0),
		Parent:     &html,
	}
	
	body := Element{
		TagName:    "body",
		Attributes: make(map[string]string),
		Children:   make([]Element, 0),
		Parent:     &html,
		Text:       doc.Body,
	}
	
	html.Children = append(html.Children, head, body)
	doc.Elements = append(doc.Elements, html)
	doc.Head = &head
	doc.BodyElem = &body
	
	// В реальном браузере здесь должен быть полноценный парсинг DOM-дерева
	// Для демонстрации используем упрощенную версию
	
	return doc, nil
}

// FindElementsByTagName находит все элементы с указанным тегом
func (d *Document) FindElementsByTagName(tagName string) []Element {
	result := make([]Element, 0)
	
	// Рекурсивный поиск по всем элементам
	for _, element := range d.Elements {
		findElementsByTagNameRecursive(&element, tagName, &result)
	}
	
	return result
}

// FindElementsByID находит элемент с указанным ID
func (d *Document) FindElementsByID(id string) *Element {
	// Рекурсивный поиск по всем элементам
	for _, element := range d.Elements {
		if found := findElementByIDRecursive(&element, id); found != nil {
			return found
		}
	}
	
	return nil
}

// Вспомогательная функция для рекурсивного поиска по тегу
func findElementsByTagNameRecursive(element *Element, tagName string, result *[]Element) {
	if strings.EqualFold(element.TagName, tagName) {
		*result = append(*result, *element)
	}
	
	for i := range element.Children {
		findElementsByTagNameRecursive(&element.Children[i], tagName, result)
	}
}

// Вспомогательная функция для рекурсивного поиска по ID
func findElementByIDRecursive(element *Element, id string) *Element {
	if element.ID == id {
		return element
	}
	
	for i := range element.Children {
		if found := findElementByIDRecursive(&element.Children[i], id); found != nil {
			return found
		}
	}
	
	return nil
}

// GetInnerHTML возвращает внутреннее HTML содержимое элемента
func (e *Element) GetInnerHTML() string {
	var sb strings.Builder
	
	// Добавляем текст элемента
	if e.Text != "" {
		sb.WriteString(e.Text)
	}
	
	// Добавляем HTML всех дочерних элементов
	for _, child := range e.Children {
		sb.WriteString(fmt.Sprintf("<%s", child.TagName))
		
		// Добавляем атрибуты
		for key, value := range child.Attributes {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", key, value))
		}
		
		sb.WriteString(">")
		sb.WriteString(child.GetInnerHTML())
		sb.WriteString(fmt.Sprintf("</%s>", child.TagName))
	}
	
	return sb.String()
}