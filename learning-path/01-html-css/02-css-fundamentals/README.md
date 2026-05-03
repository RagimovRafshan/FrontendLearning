# 🎨 Основы CSS

## 📚 Теория

### Что такое CSS?

**CSS (Cascading Style Sheets)** — это каскадные таблицы стилей, язык описания внешнего вида документа. CSS отвечает за цвета, шрифты, отступы, позиционирование и анимации.

### Способы подключения CSS

```html
<!-- 1. Внешний файл (рекомендуется) -->
<link rel="stylesheet" href="styles.css">

<!-- 2. Внутренние стили -->
<style>
    body { font-family: Arial; }
</style>

<!-- 3. Инлайн стили (не рекомендуется) -->
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

### Селекторы CSS

#### 1. Базовые селекторы

```css
/* По тегу */
p { color: black; }

/* По классу (.) */
.text-red { color: red; }

/* По ID (#) */
#header { background: gray; }

/* Универсальный селектор */
* { margin: 0; padding: 0; }

/* Группировка селекторов */
h1, h2, h3 { font-weight: bold; }
```

#### 2. Комбинаторы

```css
/* Дочерний элемент (прямой потомок) */
.parent > .child { color: blue; }

/* Потомок (любой уровень вложенности) */
.container p { font-size: 16px; }

/* Смежный_sibling (следующий элемент) */
h1 + p { margin-top: 10px; }

/* Общий sibling */
h1 ~ p { color: gray; }
```

#### 3. Псевдоклассы

```css
/* Состояния ссылок */
a:link { color: blue; }
a:visited { color: purple; }
a:hover { color: red; }
a:active { color: orange; }

/*:first-child - первый элемент */
li:first-child { font-weight: bold; }

/*:last-child - последний элемент */
li:last-child { border-bottom: none; }

/*:nth-child(n) - n-ный элемент */
li:nth-child(odd) { background: #f0f0f0; }
li:nth-child(even) { background: #fff; }
li:nth-child(3n) { color: red; } /* каждый 3-й */

/*:not() - исключение */
button:not([disabled]) { cursor: pointer; }

/*:focus - элемент в фокусе */
input:focus { outline: 2px solid blue; }
```

#### 4. Псевдоэлементы

```css
/*::before - перед содержимым */
.quote::before {
    content: """;
    color: gray;
}

/*::after - после содержимого */
.clearfix::after {
    content: "";
    display: table;
    clear: both;
}

/*::first-letter - первая буква */
p::first-letter {
    font-size: 2em;
    color: red;
}

/*::selection - выделенный текст */
::selection {
    background: yellow;
    color: black;
}
```

### Блочная модель (Box Model)

```
┌─────────────────────────────────┐
│           Margin                │
│  ┌───────────────────────────┐  │
│  │         Border            │  │
│  │  ┌─────────────────────┐  │  │
│  │  │      Padding        │  │  │
│  │  │  ┌───────────────┐  │  │  │
│  │  │  │   Content     │  │  │  │
│  │  │  │               │  │  │  │
│  │  │  └───────────────┘  │  │  │
│  │  └─────────────────────┘  │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

```css
.box {
    width: 200px;           /* ширина контента */
    height: 100px;          /* высота контента */
    padding: 20px;          /* внутренний отступ */
    border: 2px solid black;/* граница */
    margin: 10px;           /* внешний отступ */
    
    /* box-sizing меняет расчет ширины */
    box-sizing: border-box; /* ширина включает padding и border */
}
```

**Важно:** Всегда используйте `box-sizing: border-box` для предсказуемых размеров!

### Цвета в CSS

```css
/* Названия цветов */
color: red;
color: transparent;

/* HEX */
color: #ff0000;
color: #f00; /* сокращенная форма */

/* RGB */
color: rgb(255, 0, 0);
color: rgba(255, 0, 0, 0.5); /* с прозрачностью */

/* HSL */
color: hsl(0, 100%, 50%);
color: hsla(0, 100%, 50%, 0.5);

/* Современный формат */
color: rgb(255 0 0 / 50%);
```

### Шрифты и текст

```css
.text {
    /* Семейство шрифтов */
    font-family: 'Arial', sans-serif;
    
    /* Размер */
    font-size: 16px;
    font-size: 1rem; /* относительно root */
    font-size: 1.2em; /* относительно родителя */
    
    /* Начертание */
    font-weight: normal;
    font-weight: bold;
    font-weight: 700;
    
    /* Стиль */
    font-style: italic;
    
    /* Междустрочный интервал */
    line-height: 1.5;
    
    /* Выравнивание текста */
    text-align: left;
    text-align: center;
    text-align: justify;
    
    /* Декорация */
    text-decoration: underline;
    text-decoration: none;
    
    /* Тень текста */
    text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
    
    /* Трансформация */
    text-transform: uppercase;
    text-transform: lowercase;
    text-transform: capitalize;
    
    /* Межбуквенный интервал */
    letter-spacing: 1px;
    
    /* Интервал между словами */
    word-spacing: 5px;
}

/* Подключение шрифтов */
@font-face {
    font-family: 'MyFont';
    src: url('myfont.woff2') format('woff2');
    font-weight: normal;
    font-style: normal;
}
```

### Фон (Background)

```css
.element {
    /* Цвет фона */
    background-color: #f0f0f0;
    
    /* Изображение фона */
    background-image: url('image.jpg');
    
    /* Повторение */
    background-repeat: no-repeat;
    background-repeat: repeat-x;
    
    /* Позиция */
    background-position: center center;
    background-position: top right;
    
    /* Размер */
    background-size: cover;
    background-size: contain;
    background-size: 100% 100%;
    
    /* Крепление */
    background-attachment: fixed;
    
    /* Короткая запись */
    background: #f0f0f0 url('image.jpg') no-repeat center center/cover;
    
    /* Градиенты */
    background: linear-gradient(to right, red, blue);
    background: radial-gradient(circle, red, blue);
}
```

### Размеры и единицы измерения

```css
.element {
    /* Абсолютные единицы */
    width: 300px;
    height: 200px;
    
    /* Относительные единицы */
    width: 50%;           /* % от родителя */
    width: 50vw;          /* % от ширины viewport */
    height: 50vh;         /* % от высоты viewport */
    
    font-size: 1rem;      /* относительно root (16px по умолчанию) */
    font-size: 1em;       /* относительно родителя */
    
    /* Мин/макс размеры */
    min-width: 200px;
    max-width: 1200px;
    min-height: 100px;
    max-height: 80vh;
}
```

### Display свойства

```css
/* block - блочный элемент */
display: block; /* div, p, h1-h6, section */

/* inline - строчный элемент */
display: inline; /* span, a, strong, em */

/* inline-block - гибрид */
display: inline-block;

/* none - скрыть элемент */
display: none;

/* flex - flexbox */
display: flex;

/* grid - css grid */
display: grid;

/* visibility: hidden скрывает, но оставляет место */
visibility: hidden;
```

### Позиционирование

```css
/* static - по умолчанию */
position: static;

/* relative - относительно себя */
position: relative;
top: 10px;
left: 20px;

/* absolute - относительно ближайшего position:relative */
position: absolute;
top: 0;
right: 0;

/* fixed - относительно viewport */
position: fixed;
bottom: 20px;
right: 20px;

/* sticky - липкое позиционирование */
position: sticky;
top: 0;

/* z-index - порядок наложения */
z-index: 10;
```

### Границы (Borders)

```css
.border {
    /* Короткая запись */
    border: 2px solid black;
    
    /* Полная запись */
    border-width: 2px;
    border-style: solid;
    border-color: black;
    
    /* Отдельные стороны */
    border-top: 1px dashed red;
    border-bottom: 2px dotted blue;
    
    /* Скругление */
    border-radius: 10px;
    border-radius: 50%; /* круг */
    border-radius: 10px 20px 30px 40px; /* TL TR BR BL */
}
```

### Отступы (Spacing)

```css
.spacing {
    /* Margin - внешние отступы */
    margin: 10px;              /* все стороны */
    margin: 10px 20px;         /* vertical horizontal */
    margin: 10px 20px 30px;    /* T H B */
    margin: 10px 20px 30px 40px; /* T R B L */
    
    margin-top: 10px;
    margin-right: 20px;
    margin-bottom: 30px;
    margin-left: 40px;
    
    /* Auto margin для центрирования */
    margin: 0 auto;
    
    /* Padding - внутренние отступы */
    padding: 10px;
    padding: 10px 20px;
    padding: 10px 20px 30px 40px;
    
    padding-top: 10px;
    padding-right: 20px;
    padding-bottom: 30px;
    padding-left: 40px;
}
```

---

## 🎯 Практическое задание

### Задание 1: Стилизируйте блог из HTML модуля

Создайте файл `blog.css` и подключите его к `blog.html`:

1. **Базовые стили:**
   - Сброс margins/paddings (*)
   - Установите box-sizing: border-box
   - Базовый шрифт для body

2. **Header:**
   - Flexbox для расположения лого и навигации
   - Отступы и фон
   - Стили для ссылок навигации (hover эффект)

3. **Статьи:**
   - Карточки с тенью
   - Скругленные углы
   - Изображения на всю ширину карточки
   - Типографика заголовков и текста

4. **Aside:**
   - Фоновый цвет
   - Стили списка тем
   - Hover эффекты

5. **Footer:**
   - Темный фон
   - Светлый текст
   - Центрирование контента

### Задание 2: Создайте стильную форму

Стилизируйте `registration.html`:

1. Контейнер формы по центру экрана
2. Красивые input поля с:
   - Отступами
   - Границами
   - Focus состоянием
   - Placeholder стилями
3. Стилизованные кнопки
4. Адаптивная верстка
5. Валидационные состояния (valid/invalid)

---

## ❓ Проверочные вопросы

1. В чем разница между `margin` и `padding`?
2. Что делает `box-sizing: border-box`?
3. Как выбрать все элементы с классом "active"?
4. В чем разница между `display: none` и `visibility: hidden`?
5. Что такое псевдокласс `:hover`?
6. Как сделать элемент фиксированным при прокрутке?
7. Какие существуют единицы измерения в CSS?
8. Что означает `rem` и чем отличается от `em`?
9. Как центрировать элемент горизонтально?
10. Что такое специфичность селекторов?

---

## ✅ Критерии выполнения

- [ ] Код организован логически (группировка свойств)
- [ ] Используются классы вместо ID для стилизации
- [ ] Реализованы hover состояния
- [ ] Форма выглядит современно
- [ ] Соблюдена визуальная иерархия
- [ ] Нет инлайн стилей в HTML
- [ ] CSS файл подключен правильно

---

## 📖 Дополнительные ресурсы

- [MDN CSS Guide](https://developer.mozilla.org/ru/docs/Web/CSS)
- [CSS Tricks](https://css-tricks.com/)
- [Flexbox Froggy](https://flexboxfroggy.com/) - игра для изучения Flexbox
- [CSS Specificity Calculator](https://specificity.keegan.st/)
- [Can I Use](https://caniuse.com/) - проверка поддержки браузеров

---

**Следующий шаг:** Изучите `03-layouts-flex-grid` для создания сложных макетов!
