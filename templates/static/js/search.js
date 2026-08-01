document.addEventListener('DOMContentLoaded', function() {
    const recipesList = document.getElementById('recipes-list');
    const recipeNameInput = document.getElementById('recipe-name');
    const recipeCategoryInput = document.getElementById('recipe-category');
    const recipeInstructionsInput = document.getElementById('recipe-instructions');
    const recipeIngredientsList = document.getElementById('recipe-ingredients');
    const deleteBtn = document.getElementById('delete-btn');
    const editBtn = document.getElementById('edit-btn');
    const saveBtn = document.getElementById('save-btn');
    const cancelBtn = document.getElementById('cancel-btn');
    const downloadBtn = document.getElementById('download-btn')
    const host = window.location.hostname
    
    // Переменные для управления состоянием
    let currentRecipeId = null;
    let isEditMode = false;
    let originalRecipeData = null;

    // Загружаем список рецептов при загрузке страницы
    loadRecipes();
    
    // Обработчик выбора рецепта из списка
    recipesList.addEventListener('change', async function() {
        const recipeId = this.value;
        if (recipeId) {
            currentRecipeId = recipeId;
            await loadRecipeDetails(recipeId);
            // При загрузке нового рецепта выходим из режима редактирования
            if (isEditMode) {
                disableEditMode();
            }
        } else {
            clearRecipeForm();
            currentRecipeId = null;
            if (isEditMode) {
                disableEditMode();
            }
        }
    });

    // Обработчик кнопки удаления
    deleteBtn.addEventListener('click', function() {
        if (currentRecipeId) {
            deleteRecipe(currentRecipeId);
        } else {
            alert('Сначала выберите рецепт для удаления');
        }
    });

    // Обработчик кнопки редактирования
    editBtn.addEventListener('click', function() {
        if (currentRecipeId) {
            enableEditMode();
        } else {
            alert('Сначала выберите рецепт для редактирования');
        }
    });
    
    // Обработчик кнопки сохранения
    saveBtn.addEventListener('click', function() {
        saveRecipeChanges();
    });
    
    // Обработчик кнопки отмены
    cancelBtn.addEventListener('click', function() {
        disableEditMode();
        restoreOriginalData();
    });

    // Обработчик кнопки скачивания (для ссылки)
    downloadBtn.addEventListener('click', function(e) {
        // Предотвращаем переход по ссылке
        e.preventDefault();
        
        if (currentRecipeId) {
            downloadRecipe(currentRecipeId);
        } else {
            alert('Сначала выберите рецепт для скачивания');
        }
    });

        // Функция для скачивания рецепта
    function downloadRecipe(recipeId) {
        try {
            const downloadUrl = `http://${host}:8080/api/download?id=${recipeId}`;
            console.log(`📥 Скачивание рецепта с ID: ${recipeId}`);
            console.log(`🔗 URL для скачивания: ${downloadUrl}`);
            
            // Открываем ссылку для скачивания в новом окне или текущем
            window.open(downloadUrl, '_blank');
            
            // Альтернативный способ через создание временной ссылки
            // const link = document.createElement('a');
            // link.href = downloadUrl;
            // link.target = '_blank';
            // link.download = `recipe_${recipeId}.txt`; // Можно указать имя файла
            // document.body.appendChild(link);
            // link.click();
            // document.body.removeChild(link);
            
        } catch (error) {
            console.error('❌ Ошибка скачивания:', error);
            alert('Ошибка при скачивании рецепта: ' + error.message);
        }
    }

    // Функция для удаления рецепта
    async function deleteRecipe(recipeId) {
        if (!confirm('Вы уверены, что хотите удалить этот рецепт?')) {
            return;
        }

        try {
            console.log(`🗑️ Пытаемся удалить рецепт с ID: ${recipeId}`);
            
            const response = await fetch(`http://${host}:8080/api/deleteRecipe/${recipeId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            // Получаем текст ответа от сервера
            const responseText = await response.text();
            
            if (!response.ok) {
                // Если статус не 200-299, показываем ошибку от сервера
                throw new Error(responseText || `Ошибка сервера: ${response.status}`);
            }

            // Если сервер вернул OK и какой-то текст
            const successMessage = responseText || 'Рецепт успешно удален!';
            console.log('✅ Рецепт успешно удален');
            alert(successMessage);
            
            // Очищаем форму и обновляем список
            clearRecipeForm();
            currentRecipeId = null;
            await loadRecipes();
            
        } catch (error) {
            console.error('❌ Ошибка удаления:', error);
            alert('Ошибка при удалении: ' + error.message);
        }
    }

    // Функция включения режима редактирования
    function enableEditMode() {
        if (!currentRecipeId) return;
        
        isEditMode = true;
        originalRecipeData = getCurrentFormData();
        
        // Делаем поля редактируемыми
        document.getElementById('recipe-name').readOnly = false;
        document.getElementById('recipe-name').classList.remove('readonly-input');
        document.getElementById('recipe-name').classList.add('editable-input');
        
        document.getElementById('recipe-category').disabled = false;
        document.getElementById('recipe-category').classList.remove('readonly-input');
        document.getElementById('recipe-category').classList.add('editable-input');
        
        document.getElementById('recipe-instructions').readOnly = false;
        document.getElementById('recipe-instructions').classList.remove('readonly-input');
        document.getElementById('recipe-instructions').classList.add('editable-input');
        
        // Переключаем ингредиенты на текстовое поле
        const ingredientsList = document.getElementById('recipe-ingredients');
        const ingredientsEdit = document.getElementById('recipe-ingredients-edit');
        
        ingredientsList.style.display = 'none';
        ingredientsEdit.style.display = 'block';
        
        // Заполняем текстовое поле ингредиентами
        const ingredients = Array.from(ingredientsList.children).map(li => li.textContent);
        ingredientsEdit.value = ingredients.join('\n');
        
        // Показываем кнопки сохранения/отмены
        document.getElementById('save-buttons').style.display = 'flex';
        document.getElementById('edit-btn').style.display = 'none';
        
        // Добавляем индикатор режима редактирования
        document.querySelector('.recipe-form').classList.add('edit-mode');
        
        console.log('📝 Режим редактирования включен');
    }

    // Функция выключения режима редактирования
    function disableEditMode() {
        isEditMode = false;
        
        // Возвращаем поля в режим чтения
        document.getElementById('recipe-name').readOnly = true;
        document.getElementById('recipe-name').classList.remove('editable-input');
        document.getElementById('recipe-name').classList.add('readonly-input');
        
        document.getElementById('recipe-category').disabled = true;
        document.getElementById('recipe-category').classList.remove('editable-input');
        document.getElementById('recipe-category').classList.add('readonly-input');
        
        document.getElementById('recipe-instructions').readOnly = true;
        document.getElementById('recipe-instructions').classList.remove('editable-input');
        document.getElementById('recipe-instructions').classList.add('readonly-input');
        
        // Возвращаем список ингредиентов
        document.getElementById('recipe-ingredients').style.display = 'block';
        document.getElementById('recipe-ingredients-edit').style.display = 'none';
        
        // Показываем кнопку редактирования
        document.getElementById('save-buttons').style.display = 'none';
        document.getElementById('edit-btn').style.display = 'flex';
        
        // Убираем индикатор режима редактирования
        document.querySelector('.recipe-form').classList.remove('edit-mode');
        
        console.log('📝 Режим редактирования выключен');
    }

    // Функция получения текущих данных формы
    function getCurrentFormData() {
        return {
            name: document.getElementById('recipe-name').value,
            category: document.getElementById('recipe-category').value,
            instructions: document.getElementById('recipe-instructions').value,
            ingredients: Array.from(document.getElementById('recipe-ingredients').children).map(li => li.textContent)
        };
    }

    // Функция восстановления исходных данных
    function restoreOriginalData() {
        if (!originalRecipeData) return;
        
        document.getElementById('recipe-name').value = originalRecipeData.name;
        document.getElementById('recipe-category').value = originalRecipeData.category;
        document.getElementById('recipe-instructions').value = originalRecipeData.instructions;
        
        const ingredientsList = document.getElementById('recipe-ingredients');
        ingredientsList.innerHTML = '';
        originalRecipeData.ingredients.forEach(ingredient => {
            const li = document.createElement('li');
            li.textContent = ingredient;
            ingredientsList.appendChild(li);
        });
    }

    // Функция сохранения изменений
    async function saveRecipeChanges() {
    if (!currentRecipeId) {
        alert('Рецепт не выбран');
        return;
    }

    try {
        console.log('💾 Сохраняем изменения рецепта:', currentRecipeId);
        
        // Собираем данные из формы
        const updatedRecipe = {
            id: parseInt(currentRecipeId),
            name: document.getElementById('recipe-name').value.trim(),
            category: document.getElementById('recipe-category').value,
            instructions: document.getElementById('recipe-instructions').value.trim(),
            ingredients: document.getElementById('recipe-ingredients-edit').value
                .split('\n')
                .filter(item => item.trim() !== '')
        };

        // Валидация данных
        if (!updatedRecipe.name) {
            alert('Название рецепта не может быть пустым');
            return;
        }
        if (!updatedRecipe.category) {
            alert('Выберите категорию рецепта');
            return;
        }
        if (updatedRecipe.ingredients.length === 0) {
            alert('Добавьте хотя бы один ингредиент');
            return;
        }
        if (!updatedRecipe.instructions) {
            alert('Инструкция приготовления не может быть пустой');
            return;
        }

        // Показываем индикатор загрузки
        const saveBtn = document.getElementById('save-btn');
        const originalText = saveBtn.innerHTML;
        saveBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Сохранение...';
        saveBtn.disabled = true;

        // Отправляем PUT запрос на сервер
        const response = await fetch(`http://${host}:8080/api/updateRecipe/${currentRecipeId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedRecipe)
        });

        // Получаем ответ от сервера
        const responseText = await response.text();
        
        if (!response.ok) {
            throw new Error(responseText || `Ошибка сервера: ${response.status}`);
        }

        // Успешное сохранение
        const successMessage = responseText || 'Рецепт успешно обновлен!';
        console.log('✅ Рецепт обновлен:', currentRecipeId);
        alert(successMessage);
        
        // Обновляем список рецептов (на случай если изменилось название)
        await loadRecipes();
        
        // Перезагружаем детали рецепта чтобы получить актуальные данные
        await loadRecipeDetails(currentRecipeId);
        
        // Выходим из режима редактирования
        disableEditMode();
        
    } catch (error) {
        console.error('❌ Ошибка сохранения:', error);
        alert('Ошибка при сохранении: ' + error.message);
        
        // Восстанавливаем кнопку сохранения
        const saveBtn = document.getElementById('save-btn');
        saveBtn.innerHTML = '<i class="fas fa-save"></i> Сохранить';
        saveBtn.disabled = false;
    }
}
    // Функция загрузки списка рецептов
    async function loadRecipes() {
    try {
        console.log('🔄 Загружаем список рецептов...');
        
        const response = await fetch(`http://${host}:8080/api/recipes`);
        
        if (!response.ok) {
            throw new Error('Ошибка загрузки списка рецептов: ' + response.status);
        }
        
        const recipes = await response.json();
        updateRecipesList(recipes);
        
        console.log('✅ Список рецептов загружен');
        
    } catch (error) {
        console.error('❌ Ошибка:', error);
        recipesList.innerHTML = '<option value="">Ошибка загрузки рецептов</option>';
    }
}

    // Функция загрузки деталей рецепта
    async function loadRecipeDetails(recipeId) {
        try {
            console.log(`🔄 Загружаем рецепт с ID: ${recipeId}`);
            
            // Показываем что идет загрузка
            recipesList.disabled = true;
            
            const response = await fetch(`http://${host}:8080/api/recipes/${recipeId}`);
            
            if (!response.ok) {
                throw new Error('Ошибка загрузки рецепта: ' + response.status);
            }
            
            const recipe = await response.json();
            fillRecipeForm(recipe);
            
            console.log('✅ Рецепт загружен:', recipe.name);
            
        } catch (error) {
            console.error('❌ Ошибка загрузки рецепта:', error);
            alert('Не удалось загрузить рецепт: ' + error.message);
            clearRecipeForm();
        } finally {
            recipesList.disabled = false;
        }
    }

    // Функция обновления списка рецептов
    function updateRecipesList(recipes) {
        recipesList.innerHTML = '';
        
        const defaultOption = document.createElement('option');
        defaultOption.value = '';
        defaultOption.textContent = 'Выберите рецепт...';
        recipesList.appendChild(defaultOption);
        
        // Предполагаем, что сервер возвращает массив объектов {id, name}
        recipes.forEach(recipe => {
            const option = document.createElement('option');
            option.value = recipe.id;          // ID для запросов
            option.textContent = recipe.name;  // Название для показа
            recipesList.appendChild(option);
        });
        
        console.log(`📝 Загружено рецептов: ${recipes.length}`);
    }

    // Функция заполнения формы данными рецепта
    function fillRecipeForm(recipe) {
        recipeNameInput.value = recipe.name || '';
        recipeCategoryInput.value = recipe.category || '';
        recipeInstructionsInput.value = recipe.instructions || '';
        
        // Очищаем и заполняем список ингредиентов
        recipeIngredientsList.innerHTML = '';
        if (recipe.ingredients && Array.isArray(recipe.ingredients)) {
            recipe.ingredients.forEach(ingredient => {
                const li = document.createElement('li');
                li.textContent = ingredient;
                recipeIngredientsList.appendChild(li);
            });
        }
    }

    // Функция очистки формы
    function clearRecipeForm() {
        recipeNameInput.value = '';
        recipeCategoryInput.value = '';
        recipeInstructionsInput.value = '';
        recipeIngredientsList.innerHTML = '';
        recipesList.value = '';
    }
});