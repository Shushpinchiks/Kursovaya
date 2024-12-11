document.addEventListener('DOMContentLoaded', () => {
    const stars = document.querySelectorAll('.star');
    stars.forEach(star => {
        star.addEventListener('click', function () {
            const rating = this.getAttribute('data-value');
            const filmID = this.getAttribute('data-film-id');
            submitRating(rating, filmID);
            highlightStars(rating);
        });
    });
});

function highlightStars(rating) {
    const stars = document.querySelectorAll('.star');
    stars.forEach(star => {
        if (star.getAttribute('data-value') <= rating) {
            star.classList.add('selected');
        } else {
            star.classList.remove('selected');
        }
    });
}

function submitRating(rating, filmID) {
    fetch(`/rate/${filmID}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ rating: rating })
    })
        .then(response => response.json())
        .catch(error => console.error('Ошибка:', error));
}
