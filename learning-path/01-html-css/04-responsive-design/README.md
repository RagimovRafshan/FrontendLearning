# 📱 Адаптивный дизайн (Responsive Design)

## 📚 Теория

### Что такое адаптивный дизайн?

Адаптивный дизайн — это подход к веб-разработке, при котором сайт автоматически подстраивается под размер экрана устройства.

### Mobile First подход

Сначала стили для мобильных устройств, затем для больших экранов:

```css
/* Базовые стили (mobile) */
.container {
    padding: 10px;
}

/* Планшеты */
@media (min-width: 768px) {
    .container {
        padding: 20px;
    }
}

/* Десктоп */
@media (min-width: 1024px) {
    .container {
        padding: 40px;
    }
}
```

### Медиа-запросы (Media Queries)

```css
/* По ширине экрана */
@media (max-width: 576px) { /* Mobile */ }
@media (min-width: 577px) and (max-width: 768px) { /* Tablet */ }
@media (min-width: 769px) and (max-width: 1024px) { /* Small Desktop */ }
@media (min-width: 1025px) { /* Large Desktop */ }

/* По ориентации */
@media (orientation: portrait) { }
@media (orientation: landscape) { }

/* По плотности пикселей */
@media (min-resolution: 2dppx) { /* Retina */ }

/* По предпочтениям пользователя */
@media (prefers-color-scheme: dark) {
    body {
        background: #1a1a1a;
        color: #fff;
    }
}

@media (prefers-reduced-motion: reduce) {
    * {
        animation: none !important;
        transition: none !important;
    }
}
```

### Breakpoints (Контрольные точки)

Стандартные breakpoints:

```css
/* Extra small devices (phones, 576px and down) */
@media (max-width: 575.98px) { }

/* Small devices (landscape phones, 576px and up) */
@media (min-width: 576px) and (max-width: 767.98px) { }

/* Medium devices (tablets, 768px and up) */
@media (min-width: 768px) and (max-width: 991.98px) { }

/* Large devices (desktops, 992px and up) */
@media (min-width: 992px) and (max-width: 1199.98px) { }

/* Extra large devices (large desktops, 1200px and up) */
@media (min-width: 1200px) { }
```

### Относительные единицы для адаптивности

```css
.responsive {
    /* Процентная ширина */
    width: 100%;
    max-width: 1200px;
    
    /* Viewport units */
    width: 100vw;
    height: 100vh;
    font-size: 2vw;
    
    /* Container queries (современный подход) */
    container-type: inline-size;
}

@container (min-width: 400px) {
    .card {
        display: flex;
    }
}
```

### Адаптивные изображения

```html
<!-- srcset для разных разрешений -->
<img 
    src="image-800.jpg"
    srcset="image-400.jpg 400w,
            image-800.jpg 800w,
            image-1200.jpg 1200w"
    sizes="(max-width: 600px) 400px,
           (max-width: 1000px) 800px,
           1200px"
    alt="Описание">

<!-- picture для разных форматов -->
<picture>
    <source media="(min-width: 1024px)" srcset="large.webp" type="image/webp">
    <source media="(min-width: 768px)" srcset="medium.webp" type="image/webp">
    <img src="fallback.jpg" alt="Описание">
</picture>
```

### Адаптивная типографика

```css
/* Fluid typography */
html {
    font-size: calc(14px + 0.5vw);
}

/* Clamp для ограничения диапазона */
h1 {
    font-size: clamp(1.5rem, 5vw, 3rem);
}

/* Media queries для типографии */
@media (min-width: 768px) {
    h1 { font-size: 2.5rem; }
}

@media (min-width: 1024px) {
    h1 { font-size: 3.5rem; }
}
```

---

## 🎯 Практическое задание

### Задание 1: Адаптируйте блог

Сделайте ваш блог полностью адаптивным:
- Mobile: 1 колонка
- Tablet: 2 колонки (контент + aside)
- Desktop: 3 колонки

### Задание 2: Адаптивная навигация

Создайте навигацию, которая:
- На десктопе: горизонтальное меню
- На мобильном: гамбургер-меню (без JS, используя checkbox hack)

### Задание 3: Адаптивная карточка товара

Карточка должна:
- Менять расположение элементов (горизонтально/вертикально)
- Адаптировать размеры шрифтов
- Показывать разное количество информации на разных экранах

---

## ❓ Проверочные вопросы

1. Что такое Mobile First и почему это важно?
2. В чем разница между `max-width` и `min-width` в медиа-запросах?
3. Какие стандартные breakpoints вы знаете?
4. Что такое viewport units?
5. Как сделать адаптивное изображение?
6. Что делает `clamp()`?
7. Как учесть темную тему пользователя?
8. Что такое container queries?
9. Почему не стоит использовать `!important` в адаптивных стилях?
10. Как тестировать адаптивность без реальных устройств?

---

## ✅ Критерии выполнения

- [ ] Сайт корректно отображается на всех размерах
- [ ] Использован Mobile First подход
- [ ] Изображения адаптивные
- [ ] Типография читается на всех экранах
- [ ] Нет горизонтального скролла
- [ ] Touch-friendly элементы на мобильных

---

## 📖 Дополнительные ресурсы

- [MDN Media Queries](https://developer.mozilla.org/ru/docs/Web/CSS/Media_Queries)
- [Responsive Design Basics](https://web.dev/responsive-web-design-basics/)
- [CSS-Tricks Responsive Design](https://css-tricks.com/responsive-web-design/)
- [Can I Use Container Queries](https://caniuse.com/css-container-queries)

---

**Следующий шаг:** Изучите `05-advanced-css` для продвинутых техник!
