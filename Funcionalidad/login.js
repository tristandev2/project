const loginusuarios = document.querySelectorAll('#loginusuarios')
loginusuarios.addEventListener('submit', (e)=>{
    e.preventDefault()
    const email = document.querySelector('#email').value.trim();
    const password = document.querySelector('#password').value.trim();
        
        axios.post('/loguearse',{
            method:'POST',
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
                    
                });
            })
            .catch(function(error) {
                console.error('Error:', error);
                Swal.fire({
                    title: 'Error',
                    text: 'Error al iniciar sesión',
                    icon: 'error',
                    confirmButtonText: 'OK'
                    });
                    });
            });