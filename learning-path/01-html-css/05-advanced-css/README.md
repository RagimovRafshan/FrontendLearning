# 🚀 Продвинутый CSS

## 📚 Теория

### CSS Переменные (Custom Properties)

```css
:root {
    --primary-color: #3498db;
    --secondary-color: #2ecc71;
    --text-color: #333;
    --background-color: #fff;
    --spacing-unit: 8px;
    --border-radius: 4px;
    --transition-speed: 0.3s;
}

.button {
    background-color: var(--primary-color);
    color: white;
    padding: calc(var(--spacing-unit) * 2) calc(var(--spacing-unit) * 4);
    border-radius: var(--border-radius);
    transition: all var(--transition-speed) ease;
}

/* Использование с fallback */
.element {
    color: var(--text-color, #333);
}

/* Динамическое изменение в JS */
/* document.documentElement.style.setProperty('--primary-color', '#e74c3c'); */
```

### Методология BEM (Block Element Modifier)

```css
/* Block */
.card { }

/* Element */
.card__title { }
.card__image { }
.card__description { }

/* Modifier */
.card--featured { }
.card__title--large { }
.card--disabled { }

/* Пример использования */
/* HTML: <article class="card card--featured">
            <h2 class="card__title card__title--large">Заголовок</h2>
            <img class="card__image" src="...">
            <p class="card__description">Описание</p>
         </article> */
```

### CSS Анимации

```css
/* Keyframes */
@keyframes slideIn {
    from {
        transform: translateX(-100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

@keyframes pulse {
    0%, 100% { transform: scale(1); }
    50% { transform: scale(1.05); }
}

/* Применение анимации */
.animated {
    animation-name: slideIn;
    animation-duration: 0.5s;
    animation-timing-function: ease-out;
    animation-delay: 0.2s;
    animation-iteration-count: 1;
    animation-direction: normal;
    animation-fill-mode: forwards;
    
    /* Короткая запись */
    animation: slideIn 0.5s ease-out 0.2s 1 normal forwards;
}

/* Бесконечная анимация */
.pulse {
    animation: pulse 2s infinite;
}

/* Pause/Play */
.paused {
    animation-play-state: paused;
}
```

### CSS Transitions

```css
.button {
    background-color: blue;
    transform: scale(1);
    
    /* Плавный переход */
    transition: background-color 0.3s ease, 
                transform 0.2s ease-in-out;
}

.button:hover {
    background-color: darkblue;
    transform: scale(1.05);
}

/* Transition свойства */
.transition {
    transition-property: all; /* или specific: transform, opacity */
    transition-duration: 0.3s;
    transition-timing-function: ease; /* linear, ease-in, ease-out, cubic-bezier() */
    transition-delay: 0.1s;
}
```

### CSS Transform

```css
.transform-example {
    /* Перемещение */
    transform: translate(50px, 100px);
    transform: translateX(50px);
    transform: translateY(50px);
    transform: translateZ(50px); /* 3D */
    
    /* Вращение */
    transform: rotate(45deg);
    transform: rotateX(45deg); /* 3D */
    transform: rotateY(45deg);
    transform: rotateZ(45deg);
    
    /* Масштабирование */
    transform: scale(1.5);
    transform: scaleX(2);
    transform: scaleY(0.5);
    
    /* Наклон */
    transform: skew(30deg, 20deg);
    transform: skewX(30deg);
    transform: skewY(20deg);
    
    /* Комбинирование */
    transform: translate(50px, 0) rotate(45deg) scale(1.2);
    
    /* Точка трансформации */
    transform-origin: center center;
    transform-origin: top left;
    transform-origin: 50% 50%;
}
```

### CSS Filters

```css
.filtered {
    /* Размытие */
    filter: blur(5px);
    
    /* Яркость */
    filter: brightness(150%);
    
    /* Контраст */
    filter: contrast(200%);
    
    /* Оттенки серого */
    filter: grayscale(100%);
    
    /* Инверсия */
    filter: invert(100%);
    
    /* Сепия */
    filter: sepia(100%);
    
    /* Прозрачность */
    filter: opacity(50%);
    
    /* Насыщенность */
    filter: saturate(200%);
    
    /* Оттенок */
    filter: hue-rotate(90deg);
    
    /* Комбинирование */
    filter: blur(2px) grayscale(50%) brightness(120%);
    
    /* Drop shadow */
    filter: drop-shadow(5px 5px 5px rgba(0,0,0,0.5));
}
```

### CSS Grid продвинутые техники

```css
/* Auto-fit vs Auto-fill */
.grid-auto-fit {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
}

.grid-auto-fill {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

/* Subgrid (современный) */
.parent {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
}

.child {
    display: grid;
    grid-template-columns: subgrid;
}

/* Masonry layout */
.masonry {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    grid-auto-rows: 100px;
}

.item-tall {
    grid-row: span 2;
}

.item-wide {
    grid-column: span 2;
}
```

### CSS Functions

```css
.advanced-css {
    /* calc - вычисления */
    width: calc(100% - 40px);
    font-size: calc(1rem + 0.5vw);
    
    /* min/max/clamp */
    width: min(50%, 600px);
    width: max(300px, 50%);
    font-size: clamp(1rem, 2vw, 2rem);
    
    /* color functions */
    background: rgb(255 0 0 / 50%);
    background: hsl(210 100% 50% / 0.8);
    color: color-mix(in srgb, red, blue 50%);
    
    /* attr - получение атрибута */
    content: attr(data-label);
    
    /* var - переменные */
    color: var(--primary-color);
    
    /* env - безопасные зоны */
    padding: env(safe-area-inset-top) 
             env(safe-area-inset-right)
             env(safe-area-inset-bottom)
             env(safe-area-inset-left);
}
```

---

## 🎯 Практическое задание

### Задание 1: Тема сайта на CSS переменных

Создайте светлую и темную тему используя CSS переменные:
- Определите переменные в :root
- Переопределите для [data-theme="dark"]
- Добавьте переключатель темы

### Задание 2: Анимированные карточки

Создайте карточки товаров с:
- Hover эффектами (transform, shadow)
- Анимацией появления
- Micro-interactions при клике

### Задание 3: BEM структура

Рефакторите ваш блог используя методологию BEM:
- Переименуйте классы по BEM
- Создайте modifiers для разных состояний
- Документируйте структуру

---

## ❓ Проверочные вопросы

1. Что такое CSS переменные и чем они отличаются от препроцессорных?
2. Объясните структуру BEM
3. В чем разница между animation и transition?
4. Что делает transform-origin?
5. Какие timing functions вы знаете?
6. Как работает filter: drop-shadow vs box-shadow?
7. Что такое subgrid?
8. Для чего нужен clamp()?
9. Как создать бесконечную анимацию?
10. Что такое safe-area-inset?

---

## ✅ Критерии выполнения

- [ ] Правильное использование CSS переменных
- [ ] Код структурирован по BEM
- [ ] Анимации плавные и производительные
- [ ] Темы работают корректно
- [ ] Нет проблем с доступностью

---

## 📖 Дополнительные ресурсы

- [CSS Variables MDN](https://developer.mozilla.org/ru/docs/Web/CSS/Using_CSS_custom_properties)
- [BEM Methodology](http://getbem.com/)
- [CSS Animations MDN](https://developer.mozilla.org/ru/docs/Web/CSS/CSS_Animations)
- [Can I Use CSS Features](https://caniuse.com/)

---

**Следующий шаг:** Выполните `06-project` - финальный проект этапа HTML/CSS!
