document.getElementById('recipe-form').addEventListener('submit', function(e) {
    e.preventDefault();
    
    // Собираем данные формы
    const formData = {
        name: document.getElementById('recipe-name').value,
        category: document.getElementById('recipe-category').value,
        ingredients: document.getElementById('recipe-ingredients').value
            .split('\n')
            .filter(item => item.trim() !== ''),
        instructions: document.getElementById('recipe-instructions').value,
        //image: document.getElementById('recipe-image').value || ''
    };
    
    // Отправляем POST-запрос
    fetch('http://localhost:8080/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    });
});