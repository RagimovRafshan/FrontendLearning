# 📐 Layouts: Flexbox и CSS Grid

## 📚 Теория

### Flexbox (Flexible Box Layout)

Flexbox — это одномерная система раскладки, идеальная для выравнивания элементов в строку или колонку.

#### Контейнер flex

```css
.container {
    display: flex; /* или inline-flex */
    
    /* Направление оси */
    flex-direction: row;        /* по умолчанию: слева направо */
    flex-direction: row-reverse;/* справа налево */
    flex-direction: column;     /* сверху вниз */
    flex-direction: column-reverse;
    
    /* Перенос элементов */
    flex-wrap: nowrap;      /* по умолчанию: без переноса */
    flex-wrap: wrap;        /* перенос */
    flex-wrap: wrap-reverse;
    
    /* Короткая запись */
    flex-flow: row wrap;
    
    /* Выравнивание по главной оси */
    justify-content: flex-start;   /* начало */
    justify-content: flex-end;     /* конец */
    justify-content: center;       /* центр */
    justify-content: space-between;/* между элементами */
    justify-content: space-around; /* вокруг элементов */
    justify-content: space-evenly; /* равномерно */
    
    /* Выравнивание по поперечной оси */
    align-items: stretch;   /* по умолчанию: растянуть */
    align-items: flex-start;/* начало */
    align-items: flex-end;  /* конец */
    align-items: center;    /* центр */
    align-items: baseline;  /* по базовой линии текста */
    
    /* Выравнивание нескольких строк */
    align-content: flex-start;
    align-content: center;
    align-content: space-between;
}
```

#### Элементы flex

```css
.item {
    /* Порядок отображения */
    order: 0; /* по умолчанию */
    
    /* Коэффициент роста */
    flex-grow: 0; /* не растет */
    flex-grow: 1; /* растет пропорционально */
    
    /* Коэффициент сжатия */
    flex-shrink: 1; /* сжимается */
    flex-shrink: 0; /* не сжимается */
    
    /* Базовый размер */
    flex-basis: auto;
    flex-basis: 200px;
    
    /* Короткая запись */
    flex: 0 1 auto; /* grow shrink basis */
    flex: 1;        /* grow: 1, shrink: 1, basis: 0% */
    flex: 0 0 200px;/* фиксированный размер */
    
    /* Самостоятельное выравнивание */
    align-self: auto;        /* наследует от родителя */
    align-self: flex-start;
    align-self: center;
    align-self: stretch;
}
```

### CSS Grid

Grid — это двумерная система раскладки для строк и колонок одновременно.

#### Контейнер grid

```css
.container {
    display: grid; /* или inline-grid */
    
    /* Колонки */
    grid-template-columns: 200px 200px 200px;
    grid-template-columns: repeat(3, 1fr); /* 3 равные колонки */
    grid-template-columns: 1fr 2fr 1fr;    /* пропорции 1:2:1 */
    grid-template-columns: minmax(100px, 1fr);
    grid-template-columns: [start] 200px [middle] 1fr [end];
    
    /* Строки */
    grid-template-rows: 100px auto 100px;
    grid-template-rows: repeat(4, 1fr);
    
    /* Gap - отступы между ячейками */
    gap: 20px;
    row-gap: 10px;
    column-gap: 20px;
    
    /* Именованные области */
    grid-template-areas:
        "header header header"
        "sidebar content content"
        "footer footer footer";
    
    /* Выравнивание по основной оси (колонки) */
    justify-items: stretch;
    justify-items: start;
    justify-items: center;
    justify-items: end;
    
    /* Выравнивание по побочной оси (строки) */
    align-items: stretch;
    align-items: start;
    align-items: center;
    align-items: end;
    
    /* Выравнивание всего контента grid */
    justify-content: start;
    justify-content: center;
    justify-content: end;
    justify-content: space-between;
    
    align-content: start;
    align-content: center;
    align-content: end;
    align-content: space-between;
}
```

#### Элементы grid

```css
.item {
    /* Размещение в конкретных ячейках */
    grid-column-start: 1;
    grid-column-end: 3;
    grid-row-start: 1;
    grid-row-end: 2;
    
    /* Короткая запись */
    grid-column: 1 / 3; /* от линии 1 до линии 3 */
    grid-row: 1 / 2;
    grid-column: 1 / span 2; /* занять 2 колонки */
    grid-row: 1 / span 1;
    
    /* По имени области */
    grid-area: header;
    
    /* Полная запись grid-area */
    grid-area: row-start / column-start / row-end / column-end;
    
    /* Самостоятельное выравнивание */
    justify-self: stretch;
    justify-self: start;
    justify-self: center;
    justify-self: end;
    
    align-self: stretch;
    align-self: start;
    align-self: center;
    align-self: end;
}
```

---

## 🎯 Практическое задание

### Задание 1: Создайте адаптивную навигацию на Flexbox

Создайте компонент навигации:
- Логотип слева
- Навигационные ссылки по центру
- Кнопка входа справа
- На мобильных устройствах: гамбургер-меню

### Задание 2: Создайте галерею на Grid

Создайте фотогалерею:
- Сетка 3x3 на десктопе
- 2 колонки на планшете
- 1 колонка на мобильном
- Используйте `grid-template-areas` для сложной раскладки

### Задание 3: Карточки товаров

Создайте сетку карточек товаров:
- Flexbox или Grid на выбор
- Адаптивность: 4 колонки (desktop), 2 (tablet), 1 (mobile)
- Карточка должна содержать: изображение, название, цену, кнопку

---

## ❓ Проверочные вопросы

1. В чем разница между Flexbox и Grid?
2. Что делает `justify-content` vs `align-items` во Flexbox?
3. Как создать 3 равные колонки в Grid?
4. Что означает `fr` в Grid?
5. Как сделать элемент spanning через несколько колонок?
6. Что такое `flex-grow` и как он работает?
7. Как центрировать элемент по вертикали и горизонтали во Flexbox?
8. Для чего используется `gap`?
9. Что такое `minmax()` в Grid?
10. Когда использовать Flexbox, а когда Grid?

---

## ✅ Критерии выполнения

- [ ] Понимание разницы между осями во Flexbox
- [ ] Умение создавать сетки в Grid
- [ ] Адаптивность реализована корректно
- [ ] Код чистый и понятный
- [ ] Использованы современные свойства

---

## 📖 Дополнительные ресурсы

- [CSS Tricks Flexbox Guide](https://css-tricks.com/snippets/css/a-guide-to-flexbox/)
- [CSS Tricks Grid Guide](https://css-tricks.com/snippets/css/complete-guide-grid/)
- [Flexbox Froggy](https://flexboxfroggy.com/)
- [Grid Garden](https://cssgridgarden.com/)
- [MDN Flexbox](https://developer.mozilla.org/ru/docs/Web/CSS/CSS_Flexible_Box_Layout)
- [MDN Grid](https://developer.mozilla.org/ru/docs/Web/CSS/CSS_Grid_Layout)

---

**Следующий шаг:** Изучите `04-responsive-design` для создания адаптивных сайтов!
