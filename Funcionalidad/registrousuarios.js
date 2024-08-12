

document.addEventListener('DOMContentLoaded', function() {
    const registroForm = document.getElementById('registrousuario');

    registroForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Evita el envío normal del formulario

        const name = document.getElementById('name').value;
        const apellido = document.getElementById('apellido').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        axios.post('/insertar', {
            name: name,
            apellido: apellido,
            email: email,
            password: password
        })
        .then(function(response) {
            
            console.log('Bienvenido:', response.data);
            Swal.fire({
                title: 'Éxito!',
                text: response.data.message, // Muestra el mensaje de la respuesta
                icon: 'success',
                confirmButtonText: 'OK'
            }).then(() => {
                 // O la URL que desees
            });
        })
        .catch(function(error) {
            // Maneja los errores
            console.error('Error al registrar:', error);
            Swal.fire({
                title: 'Error!',
                text: 'Hubo un problema al registrar el usuario.',
                icon: 'error',
                confirmButtonText: 'OK'
            });
        });
    });
});
