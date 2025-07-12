package renderer

import (
	"fmt"
	"log"
	"strings"

	"github.com/baneronetwo/gluglu/internal/browser/html"
)

// Renderer представляет движок рендеринга
type Renderer struct {}

// Document представляет отрендеренный документ
type Document struct {
	Title    string
	Elements []RenderedElement
	Width    int
	Height   int
}

// RenderedElement представляет отрендеренный элемент
type RenderedElement struct {
	TagName    string
	Text       string
	X          int
	Y          int
	Width      int
	Height     int
	Color      string
	Background string
	Children   []RenderedElement
}

// NewRenderer создает новый движок рендеринга
func NewRenderer() *Renderer {
	log.Println("Инициализация движка рендеринга...")
	return &Renderer{}
}

// Render выполняет рендеринг HTML документа
func (r *Renderer) Render(doc *html.Document) *Document {
	log.Println("Рендеринг HTML документа...")
	
	// Создаем отрендеренный документ
	renderedDoc := &Document{
		Title:    doc.Title,
		Elements: make([]RenderedElement, 0),
		Width:    800, // Стандартная ширина
		Height:   600, // Стандартная высота
	}
	
	// Рендерим все элементы, начиная с корневого
	for _, element := range doc.Elements {
		renderedElement := r.renderElement(&element, 0, 0)
		renderedDoc.Elements = append(renderedDoc.Elements, renderedElement)
	}
	
	return renderedDoc
}

// renderElement рендерит отдельный HTML элемент
func (r *Renderer) renderElement(element *html.Element, x, y int) RenderedElement {
	// Создаем отрендеренный элемент
	renderedElement := RenderedElement{
		TagName:    element.TagName,
		Text:       element.Text,
		X:          x,
		Y:          y,
		Width:      100, // Стандартная ширина
		Height:     20,  // Стандартная высота
		Color:      "#000000", // Черный цвет текста по умолчанию
		Background: "#FFFFFF", // Белый фон по умолчанию
		Children:   make([]RenderedElement, 0),
	}
	
	// Применяем стили из атрибутов
	if style, ok := element.Attributes["style"]; ok {
		r.applyStyles(&renderedElement, style)
	}
	
	// Рендерим дочерние элементы
	currentY := y
	for _, child := range element.Children {
		childElement := r.renderElement(&child, x+10, currentY+renderedElement.Height)
		renderedElement.Children = append(renderedElement.Children, childElement)
		currentY += childElement.Height
	}
	
	// Обновляем высоту элемента с учетом дочерних элементов
	if len(renderedElement.Children) > 0 {
		lastChild := renderedElement.Children[len(renderedElement.Children)-1]
		renderedElement.Height = (lastChild.Y - y) + lastChild.Height
	}
	
	return renderedElement
}

// applyStyles применяет CSS стили к элементу
func (r *Renderer) applyStyles(element *RenderedElement, styleStr string) {
	// Разбиваем строку стилей на отдельные правила
	styles := strings.Split(styleStr, ";")
	
	for _, style := range styles {
		// Разбиваем правило на имя и значение
		parts := strings.Split(style, ":")
		if len(parts) != 2 {
			continue
		}
		
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Применяем стиль в зависимости от его имени
		switch name {
		case "color":
			element.Color = value
		case "background-color":
			element.Background = value
		case "width":
			// В реальном браузере здесь должен быть парсинг размеров (px, em, % и т.д.)
			if strings.HasSuffix(value, "px") {
				widthStr := strings.TrimSuffix(value, "px")
				var width int
				if _, err := fmt.Sscanf(widthStr, "%d", &width); err == nil {
					element.Width = width
				}
			}
		case "height":
			// В реальном браузере здесь должен быть парсинг размеров (px, em, % и т.д.)
			if strings.HasSuffix(value, "px") {
				heightStr := strings.TrimSuffix(value, "px")
				var height int
				if _, err := fmt.Sscanf(heightStr, "%d", &height); err == nil {
					element.Height = height
				}
			}
		}
	}
}

// GetTextRepresentation возвращает текстовое представление отрендеренного документа
func (d *Document) GetTextRepresentation() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Документ: %s\n", d.Title))
	sb.WriteString(fmt.Sprintf("Размер: %dx%d\n\n", d.Width, d.Height))
	
	// Рекурсивно добавляем текстовое представление всех элементов
	for _, element := range d.Elements {
		addElementTextRepresentation(&sb, element, 0)
	}
	
	return sb.String()
}

// addElementTextRepresentation добавляет текстовое представление элемента
func addElementTextRepresentation(sb *strings.Builder, element RenderedElement, indent int) {
	// Добавляем отступ
	indentStr := strings.Repeat("  ", indent)
	
	// Добавляем информацию об элементе
	sb.WriteString(fmt.Sprintf("%s<%s> (%d,%d) %dx%d\n", indentStr, element.TagName, element.X, element.Y, element.Width, element.Height))
	
	// Добавляем текст элемента, если он есть
	if element.Text != "" {
		sb.WriteString(fmt.Sprintf("%s  Текст: %s\n", indentStr, element.Text))
	}
	
	// Рекурсивно добавляем дочерние элементы
	for _, child := range element.Children {
		addElementTextRepresentation(sb, child, indent+1)
	}
	
	// Закрываем тег
	sb.WriteString(fmt.Sprintf("%s</%s>\n", indentStr, element.TagName))
}