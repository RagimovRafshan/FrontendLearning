# Модуль 2: Основы CSS

## 📖 Теория

### Что такое CSS?
CSS (Cascading Style Sheets) — это язык стилей, который определяет внешний вид и расположение элементов на странице.

### Способы подключения CSS

```html
<!-- 1. Внешний файл (рекомендуется) -->
<link rel="stylesheet" href="styles.css">

<!-- 2. Встроенный стиль в head -->
<style>
    body { margin: 0; }
</style>

<!-- 3. Инлайн-стили (не рекомендуется) -->
<div style="color: red;">Текст</div>
```

### Синтаксис CSS

```css
/* Селектор { свойство: значение; } */
selector {
    property: value;
}

/* Пример */
.button {
    color: white;
    background-color: blue;
    padding: 10px 20px;
    border-radius: 5px;
}
```

### Типы селекторов

```css
/* По тегу */
p { color: black; }

/* По классу */
.text { font-size: 16px; }

/* По ID */
#header { height: 60px; }

/* По атрибуту */
input[type="text"] { border: 1px solid gray; }

/* Псевдоклассы */
button:hover { background-color: darkblue; }
a:visited { color: purple; }
li:first-child { font-weight: bold; }

/* Псевдоэлементы */
p::first-letter { font-size: 24px; }
p::before { content: "→ "; }

/* Комбинирование */
.container .item { }           /* потомок */
.container > .item { }         /* прямой потомок */
.item + .item { }              /* следующий сосед */
.item ~ .item { }              /* все соседи после */
```

### Каскад и специфичность

Порядок приоритета (от меньшего к большему):
1. Тег (`p`) - 0,0,0,1
2. Класс (`.class`), псевдоклассы, атрибуты - 0,0,1,0
3. ID (`#id`) - 0,1,0,0
4. `!important` - перебивает всё (избегать!)

```css
p { color: red; }                    /* специфичность: 0,0,0,1 */
.text { color: blue; }               /* специфичность: 0,0,1,0 */
#main { color: green; }              /* специфичность: 0,1,0,0 */
p.text { color: orange; }            /* специфичность: 0,0,1,1 */
```

### Box Model (Блочная модель)

Каждый элемент состоит из:
- **Content** - содержимое
- **Padding** - внутренний отступ
- **Border** - граница
- **Margin** - внешний отступ

```css
.box {
    width: 200px;
    padding: 20px;
    border: 2px solid black;
    margin: 10px;
    
    /* Важно для предсказуемых размеров */
    box-sizing: border-box; /* ширина включает padding и border */
}
```

### Flexbox

```css
.container {
    display: flex;
    
    /* Выравнивание по главной оси */
    justify-content: center; /* или flex-start, flex-end, space-between, space-around */
    
    /* Выравнивание по поперечной оси */
    align-items: center; /* или flex-start, flex-end, stretch */
    
    /* Направление */
    flex-direction: row; /* или column */
    
    /* Перенос строк */
    flex-wrap: wrap;
}

.item {
    flex: 1; /* пропорции роста/сжатия */
    /* или */
    flex-grow: 1;   /* коэффициент роста */
    flex-shrink: 0; /* коэффициент сжатия */
    flex-basis: 200px; /* базовый размер */
}
```

### Grid Layout

```css
.container {
    display: grid;
    
    /* Колонки */
    grid-template-columns: 200px 1fr 2fr;
    /* или */
    grid-template-columns: repeat(3, 1fr);
    
    /* Строки */
    grid-template-rows: auto 100px;
    
    /* Отступы между ячейками */
    gap: 20px;
    
    /* Именованные области */
    grid-template-areas:
        "header header header"
        "sidebar main main"
        "footer footer footer";
}

.item {
    grid-area: sidebar; /* размещение в именованной области */
    /* или */
    grid-column: 1 / 3; /* от какой до какой колонки */
    grid-row: 1 / 2;
}
```

### Адаптивный дизайн

```css
/* Медиа-запросы */
@media (max-width: 768px) {
    .container {
        flex-direction: column;
    }
}

@media (min-width: 769px) and (max-width: 1024px) {
    .container {
        width: 90%;
    }
}

/* Mobile First подход */
.base-styles { } /* стили для мобильных */

@media (min-width: 768px) {
    /* стили для планшетов */
}

@media (min-width: 1024px) {
    /* стили для десктопов */
}
```

### CSS-переменные

```css
:root {
    --primary-color: #3498db;
    --secondary-color: #2ecc71;
    --spacing-unit: 8px;
    --font-size-base: 16px;
}

.button {
    background-color: var(--primary-color);
    padding: calc(var(--spacing-unit) * 2);
    font-size: var(--font-size-base);
}

/* Можно переопределять в контексте */
.dark-theme {
    --primary-color: #2c3e50;
}
```

### Методология BEM

**B**lock - **E**lement - **M**odifier

```css
/* Блок */
.card { }

/* Элемент блока */
.card__title { }
.card__image { }

/* Модификатор */
.card--featured { }
.card__title--large { }
```

```html
<div class="card card--featured">
    <h2 class="card__title card__title--large">Заголовок</h2>
    <img class="card__image" src="..." alt="">
</div>
```

---

## ✅ Практическое задание

### Задание 1: Стилизация блога

Возьмите HTML из предыдущего модуля (`src/html-css/blog.html`) и добавьте стили:

1. Создайте файл `src/html-css/blog.css`
2. Реализуйте:
   - Базовые стили (сброс margin/padding, шрифты)
   - Стилизацию шапки с фиксированной высотой
   - Сетку для статей (Grid или Flexbox)
   - Карточки статей с тенями и hover-эффектами
   - Боковую панель
   - Адаптивность для мобильных устройств

### Задание 2: Компоненты

Создайте библиотеку компонентов в отдельном файле:

1. Кнопки (разные размеры, цвета, состояния)
2. Формы (стилизованные input, select, checkbox)
3. Карточки товаров
4. Навигационное меню

Используйте методологию BEM и CSS-переменные.

---

## 📝 Проверочные вопросы

1. Что такое каскад в CSS?
2. Как рассчитать специфичность селектора?
3. В чем разница между `display: none` и `visibility: hidden`?
4. Что делает `box-sizing: border-box`?
5. Когда использовать Flexbox, а когда Grid?
6. Что такое Mobile First подход?

---

## 🎯 Критерии выполнения

- [ ] Использованы Flexbox или Grid для макета
- [ ] Реализована адаптивность (минимум 2 брейкпоинта)
- [ ] Применены CSS-переменные для цветов и отступов
- [ ] Соблюдена методология BEM
- [ ] Код отформатирован и структурирован

---

## 🔗 Полезные ресурсы

- [CSS Tricks - Flexbox Guide](https://css-tricks.com/snippets/css/a-guide-to-flexbox/)
- [CSS Tricks - Grid Guide](https://css-tricks.com/snippets/css/complete-guide-grid/)
- [Specificity Calculator](https://specificity.keegan.st/)
- [Can I Use](https://caniuse.com/) - проверка поддержки браузеров

---

**Следующий модуль:** Продвинутый CSS →
