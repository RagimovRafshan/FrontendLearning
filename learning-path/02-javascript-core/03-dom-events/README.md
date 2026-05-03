# 🌐 DOM и События

## Теория

### Выбор элементов

```javascript
document.getElementById('id');
document.querySelector('.class');
document.querySelectorAll('div');
document.getElementsByTagName('div');
```

### Манипуляции

```javascript
element.textContent = 'text';
element.innerHTML = '<span>HTML</span>';
element.setAttribute('data-id', '1');
element.classList.add('active');
element.style.color = 'red';
```

### События

```javascript
element.addEventListener('click', (e) => {
    e.preventDefault();
    console.log(e.target);
});

element.removeEventListener('click', handler);
```

## Практика
1. Создайте модальное окно
2. Реализуйте табы
3. Сделайте форму с валидацией

**Следующий шаг:** `04-asynchronous-js`
