document.getElementById('recipe-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    // Собираем данные формы
    const formData = {
        name: document.getElementById('recipe-name').value,
        category: document.getElementById('recipe-category').value,
        ingredients: document.getElementById('recipe-ingredients').value
            .split('\n')
            .filter(item => item.trim() !== ''),
        instructions: document.getElementById('recipe-instructions').value
    };
    
    // Кнопка отправки
    const submitBtn = document.querySelector('#recipe-form button[type="submit"]');
    const originalBtnText = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Отправка...';
    
    try {
        // Отправляем POST-запрос
        const response = await fetch('http://192.168.1.9:8080/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });
        
        // Получаем текст ответа (для ошибок и успеха)
        const responseText = await response.text();
        
        if (!response.ok) {
            // Показываем ошибку
            showAlert(responseText, 'error');
            return;
        }
        
        // Показываем успех
        showAlert('Рецепт успешно добавлен!', 'success');
        this.reset(); // Очищаем форму
        
    } catch (error) {
        // Обработка сетевых ошибок
        showAlert('Ошибка сети: ' + error.message, 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalBtnText;
    }
  });
  
  // Функция показа уведомлений
  function showAlert(message, type = 'success') {
      // Создаем контейнер, если его еще нет
      let container = document.querySelector('.notification-container');
      if (!container) {
          container = document.createElement('div');
          container.className = 'notification-container';
          document.body.appendChild(container);
      }
  
      const notification = document.createElement('div');
      notification.className = `notification notification-${type}`;
      notification.innerHTML = `
          <span class="notification-message">${message}</span>
          <button class="notification-close">&times;</button>
      `;
      
      // Добавляем уведомление в контейнер
      container.appendChild(notification);
      
      // Автоудаление через 3.5 секунды (с учетом анимации)
      setTimeout(() => {
          notification.remove();
      }, 3500);
      
      // Закрытие по клику
      notification.querySelector('.notification-close').addEventListener('click', () => {
          notification.remove();
      });
  }