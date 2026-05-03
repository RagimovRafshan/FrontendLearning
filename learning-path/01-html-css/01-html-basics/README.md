# 📝 Основы HTML5

## 📚 Теория

### Что такое HTML?

**HTML (HyperText Markup Language)** — это язык разметки гипертекста, который используется для создания структуры веб-страницы. HTML не является языком программирования, он описывает, какие элементы находятся на странице и как они организованы.

### Структура HTML-документа

```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Заголовок страницы</title>
    <meta name="description" content="Описание страницы для поисковиков">
</head>
<body>
    <!-- Содержимое страницы -->
</body>
</html>
```

### Основные теги HTML5

#### 1. Семантические теги

Семантические теги описывают **смысл** содержимого, а не только его внешний вид:

```html
<header>Шапка сайта</header>
<nav>Навигация</nav>
<main>Основное содержимое</main>
<article>Независимая статья</article>
<section>Секция документа</section>
<aside>Боковая панель</aside>
<footer>Подвал сайта</footer>
```

**Почему семантика важна:**
- ✅ Улучшает доступность (скринридеры понимают структуру)
- ✅ Лучше для SEO (поисковики понимают контент)
- ✅ Чище и понятнее код

#### 2. Текстовые элементы

```html
<h1>Заголовок первого уровня</h1>
<h2>Заголовок второго уровня</h2>
<h3>Заголовок третьего уровня</h3>

<p>Это параграф текста.</p>

<strong>Жирный текст (важный)</strong>
<em>Курсив (акцент)</em>
<mark>Выделенный текст</mark>
<small>Мелкий текст</small>
<del>Удаленный текст</del>
<ins>Добавленный текст</ins>

<blockquote>Цитата из другого источника</blockquote>
<code>Код</code>
<pre>Предформатированный текст</pre>
```

#### 3. Списки

```html
<!-- Маркированный список -->
<ul>
    <li>Элемент 1</li>
    <li>Элемент 2</li>
    <li>Элемент 3</li>
</ul>

<!-- Нумерованный список -->
<ol>
    <li>Первый шаг</li>
    <li>Второй шаг</li>
    <li>Третий шаг</li>
</ol>

<!-- Список определений -->
<dl>
    <dt>HTML</dt>
    <dd>Язык разметки гипертекста</dd>
    <dt>CSS</dt>
    <dd>Каскадные таблицы стилей</dd>
</dl>
```

#### 4. Ссылки и изображения

```html
<!-- Ссылка -->
<a href="https://example.com" target="_blank" rel="noopener noreferrer">
    Перейти на сайт
</a>

<!-- Изображение -->
<img src="image.jpg" alt="Описание изображения" width="300" height="200">

<!-- Ссылка с изображением -->
<a href="page.html">
    <img src="thumb.jpg" alt="Превью">
</a>
```

**Важно:** Всегда указывайте `alt` для изображений — это нужно для доступности и SEO!

#### 5. Формы

```html
<form action="/submit" method="POST">
    <!-- Текстовое поле -->
    <label for="name">Имя:</label>
    <input type="text" id="name" name="name" required placeholder="Введите имя">
    
    <!-- Email -->
    <label for="email">Email:</label>
    <input type="email" id="email" name="email" required>
    
    <!-- Пароль -->
    <label for="password">Пароль:</label>
    <input type="password" id="password" name="password" minlength="8" required>
    
    <!-- Число -->
    <label for="age">Возраст:</label>
    <input type="number" id="age" name="age" min="18" max="100">
    
    <!-- Дата -->
    <label for="birthdate">Дата рождения:</label>
    <input type="date" id="birthdate" name="birthdate">
    
    <!-- Выпадающий список -->
    <label for="country">Страна:</label>
    <select id="country" name="country">
        <option value="">Выберите страну</option>
        <option value="ru">Россия</option>
        <option value="us">США</option>
        <option value="de">Германия</option>
    </select>
    
    <!-- Чекбокс -->
    <label>
        <input type="checkbox" name="subscribe" value="yes">
        Подписаться на рассылку
    </label>
    
    <!-- Радио кнопки -->
    <fieldset>
        <legend>Пол:</legend>
        <label>
            <input type="radio" name="gender" value="male"> Мужской
        </label>
        <label>
            <input type="radio" name="gender" value="female"> Женский
        </label>
    </fieldset>
    
    <!-- Текстовая область -->
    <label for="message">Сообщение:</label>
    <textarea id="message" name="message" rows="5" cols="40"></textarea>
    
    <!-- Кнопка отправки -->
    <button type="submit">Отправить</button>
</form>
```

#### 6. Таблицы

```html
<table>
    <thead>
        <tr>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Возраст</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>Иван</td>
            <td>Иванов</td>
            <td>25</td>
        </tr>
        <tr>
            <td>Петр</td>
            <td>Петров</td>
            <td>30</td>
        </tr>
    </tbody>
    <tfoot>
        <tr>
            <td colspan="2">Всего записей:</td>
            <td>2</td>
        </tr>
    </tfoot>
</table>
```

#### 7. Медиа элементы

```html
<!-- Видео -->
<video controls width="640" poster="poster.jpg">
    <source src="video.mp4" type="video/mp4">
    <source src="video.webm" type="video/webm">
    Ваш браузер не поддерживает видео.
</video>

<!-- Аудио -->
<audio controls>
    <source src="audio.mp3" type="audio/mpeg">
    <source src="audio.ogg" type="audio/ogg">
    Ваш браузер не поддерживает аудио.
</audio>

<!-- Iframe (встраивание) -->
<iframe 
    src="https://www.youtube.com/embed/VIDEO_ID" 
    title="YouTube video"
    allowfullscreen>
</iframe>
```

### Атрибуты HTML

**Глобальные атрибуты** (работают со всеми тегами):
- `id` — уникальный идентификатор элемента
- `class` — классы для стилизации
- `style` — инлайн стили (не рекомендуется)
- `title` — всплывающая подсказка
- `data-*` — пользовательские данные
- `hidden` — скрыть элемент
- `tabindex` — порядок навигации Tab

### Валидация HTML

Проверяйте свой код на валидность: [W3C Validator](https://validator.w3.org/)

---

## 🎯 Практическое задание

### Задание 1: Создайте семантическую структуру блога

Создайте файл `blog.html` со следующей структурой:

1. **Header** с:
   - Логотипом (текст или img)
   - Навигацией (Главная, О нас, Контакты)

2. **Main** с:
   - Заголовком страницы (H1)
   - Двумя статьями (`<article>`), каждая содержит:
     - Заголовок статьи (H2)
     - Дату публикации
     - Изображение
     - Текст статьи (2-3 параграфа)
     - Кнопку "Читать далее"

3. **Aside** с:
   - Заголовком "Популярные темы"
   - Списком из 5 тем

4. **Footer** с:
   - Копирайтом
   - Ссылками на соцсети

### Задание 2: Создайте форму регистрации

Создайте файл `registration.html` с формой, содержащей:
- Имя (required)
- Фамилия (required)
- Email (required, type="email")
- Пароль (minlength=8, required)
- Подтверждение пароля
- Дата рождения
- Пол (radio buttons)
- Интересы (checkboxes, минимум 3)
- Страна (select)
- О себе (textarea)
- Чекбокс "Согласен с правилами" (required)
- Кнопка отправки

---

## ❓ Проверочные вопросы

Ответьте на вопросы без подглядывания:

1. В чем разница между `<div>` и `<section>`?

section это секция документа а Div ты мне не описывал. ты в целом мало дал информации для ответов на эти вопросы и для выполнения задания
2. Зачем нужен атрибут `alt` у изображений?

для доступности(хотя ты не объяснил что это значит) и Seo
3. Какой тег использовать для главного заголовка страницы?

'h1'
4. В чем разница между `<strong>` и `<b>`?

'strong' это жирный шрифт, а 'b' предполагаю что курсив или подчеркивание. ты не объяснил в документе, я не ебу
5. Как сделать ссылку, открывающуюся в новой вкладке?

target="_blank"
6. Какие существуют типы input в формах? (назовите минимум 5)

text, email, password, number, date, country, checkbox, radio
7. Зачем нужны семантические теги?

для ориентации в html и внешнего вида
8. В чем разница между `id` и `class`?

id обозначает переменную, когда class работу с css
9. Как создать маркированный список?

'ul'
10. Что означает DOCTYPE html?

наверное то что это html верстка, а так - я не ебу ты не объяснял
---

## ✅ Критерии выполнения

- [ ] Код валидный (проверено на W3C Validator)
- [ ] Использованы семантические теги
- [ ] Все формы имеют label для каждого input
- [ ] Изображения имеют alt атрибут
- [ ] Ссылки имеют осмысленный текст
- [ ] Код отформатирован и читаем
- [ ] Правильная иерархия заголовков (h1 → h2 → h3)

---

## 📖 Дополнительные ресурсы

- [MDN HTML Guide](https://developer.mozilla.org/ru/docs/Web/HTML)
- [HTML5 Semantic Elements](https://www.w3schools.com/html/html5_semantic_elements.asp)
- [Web Accessibility Initiative](https://www.w3.org/WAI/fundamentals/accessibility-intro/)
- [HTML Living Standard](https://html.spec.whatwg.org/)

---

**Следующий шаг:** После изучения HTML переходите к `02-css-fundamentals` для стилизации ваших страниц!
