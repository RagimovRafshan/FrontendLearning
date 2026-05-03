# ⏱ Асинхронный JavaScript

## Теория

### Callbacks

```javascript
setTimeout(() => {
    console.log('Через 1 сек');
}, 1000);
```

### Promises

```javascript
fetch('/api/data')
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error(error));
```

### Async/Await

```javascript
async function getData() {
    try {
        const response = await fetch('/api/data');
        const data = await response.json();
        return data;
    } catch (error) {
        console.error(error);
    }
}
```

## Практика
1. Fetch API запросы
2. Параллельные запросы (Promise.all)
3. Retry логика

**Следующий шаг:** `05-es6-plus`
