package js

import (
	"log"

	"github.com/baneronetwo/gluglu/internal/browser/html"
	"github.com/robertkrimen/otto"
)

// Engine представляет JavaScript движок
type Engine struct {
	vm *otto.Otto
}

// NewEngine создает новый JavaScript движок
func NewEngine() *Engine {
	log.Println("Инициализация JavaScript движка...")
	
	// Создаем новый экземпляр Otto VM
	vm := otto.New()
	
	// Возвращаем новый движок
	return &Engine{
		vm: vm,
	}
}

// Execute выполняет JavaScript код в контексте документа
func (e *Engine) Execute(doc *html.Document) {
	log.Println("Выполнение JavaScript...")
	
	// Находим все скрипты в документе
	scripts := doc.FindElementsByTagName("script")
	
	// Создаем объект document для доступа из JavaScript
	e.setupDocumentObject(doc)
	
	// Выполняем каждый скрипт
	for _, script := range scripts {
		// Проверяем, является ли скрипт внешним (имеет атрибут src)
		if src, ok := script.Attributes["src"]; ok && src != "" {
			// В реальном браузере здесь должна быть загрузка внешнего скрипта
			log.Printf("Внешний скрипт не загружен: %s", src)
			continue
		}
		
		// Выполняем встроенный скрипт
		if script.Text != "" {
			_, err := e.vm.Run(script.Text)
			if err != nil {
				log.Printf("Ошибка выполнения JavaScript: %v", err)
			}
		}
	}
}

// setupDocumentObject настраивает объект document для доступа из JavaScript
func (e *Engine) setupDocumentObject(doc *html.Document) {
	// Создаем объект document
	documentObj, _ := e.vm.Object("document = {}")
	
	// Добавляем свойства
	documentObj.Set("title", doc.Title)
	
	// Добавляем методы
	documentObj.Set("getElementById", func(call otto.FunctionCall) otto.Value {
		// Получаем ID из аргументов
		id, _ := call.Argument(0).ToString()
		
		// Находим элемент по ID
		element := doc.FindElementsByID(id)
		
		if element == nil {
			// Если элемент не найден, возвращаем null
			nullValue, _ := otto.NullValue().ToValue()
			return nullValue
		}
		
		// Создаем объект элемента
		elementObj, _ := e.vm.Object("({})")
		
		// Добавляем свойства элемента
		elementObj.Set("tagName", element.TagName)
		elementObj.Set("id", element.ID)
		elementObj.Set("innerHTML", element.GetInnerHTML())
		
		// Добавляем метод setAttribute
		elementObj.Set("setAttribute", func(call otto.FunctionCall) otto.Value {
			name, _ := call.Argument(0).ToString()
			value, _ := call.Argument(1).ToString()
			
			element.Attributes[name] = value
			
			// Специальная обработка для id
			if name == "id" {
				element.ID = value
			}
			
			// Возвращаем undefined
			undefinedValue, _ := otto.UndefinedValue().ToValue()
			return undefinedValue
		})
		
		// Добавляем свойство textContent с геттером и сеттером
		elementObj.Set("textContent", element.Text)
		
		// Возвращаем объект элемента
		return elementObj.Value()
	})
	
	// Добавляем метод getElementsByTagName
	documentObj.Set("getElementsByTagName", func(call otto.FunctionCall) otto.Value {
		// Получаем имя тега из аргументов
		tagName, _ := call.Argument(0).ToString()
		
		// Находим элементы по имени тега
		elements := doc.FindElementsByTagName(tagName)
		
		// Создаем массив для результатов
		resultArray, _ := e.vm.Object("([])")
		
		// Добавляем каждый элемент в массив
		for i, element := range elements {
			// Создаем объект элемента
			elementObj, _ := e.vm.Object("({})")
			
			// Добавляем свойства элемента
			elementObj.Set("tagName", element.TagName)
			elementObj.Set("id", element.ID)
			elementObj.Set("innerHTML", element.GetInnerHTML())
			
			// Добавляем элемент в массив
			resultArray.Set(i, elementObj)
		}
		
		// Возвращаем массив элементов
		return resultArray.Value()
	})
	
	// Добавляем метод createElement
	documentObj.Set("createElement", func(call otto.FunctionCall) otto.Value {
		// Получаем имя тега из аргументов
		tagName, _ := call.Argument(0).ToString()
		
		// Создаем новый элемент
		element := html.Element{
			TagName:    tagName,
			Attributes: make(map[string]string),
			Children:   make([]html.Element, 0),
		}
		
		// Создаем объект элемента
		elementObj, _ := e.vm.Object("({})")
		
		// Добавляем свойства элемента
		elementObj.Set("tagName", element.TagName)
		
		// Возвращаем объект элемента
		return elementObj.Value()
	})
}

// EvaluateScript выполняет JavaScript код и возвращает результат
func (e *Engine) EvaluateScript(script string) (string, error) {
	value, err := e.vm.Run(script)
	if err != nil {
		return "", err
	}
	
	result, _ := value.ToString()
	return result, nil
}