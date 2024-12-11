function deleteUser(userId) {
    if (confirm('Вы уверены, что хотите удалить этого пользователя?')) {
        fetch('/delete_user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: userId })
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Пользователь успешно удален.');
                    location.reload();
                } else {
                    alert('Ошибка при удалении пользователя.');
                }
            })
            .catch(error => console.error('Ошибка:', error));
    }
}

function deleteFilm(filmId) {
    if (confirm('Вы уверены, что хотите удалить этого пользователя?')) {
        fetch('/delete_film', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: filmId })
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Фильм успешно удален.');
                    location.reload();
                } else {
                    alert('Ошибка при удалении фильма.');
                }
            })
            .catch(error => console.error('Ошибка:', error));
    }
}