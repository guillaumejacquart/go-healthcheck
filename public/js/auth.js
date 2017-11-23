var auth = {
    getToken: function(){
        return localStorage.getItem('go-healthcheck-token')
    },
    setToken: function(token){
        localStorage.setItem('go-healthcheck-token', token)
    },
    logout: function(){
        localStorage.removeItem('go-healthcheck-token')
        $('#login').modal('show');
    }
}

Vue.http.interceptors.push(function(request, next) {
    if (auth.getToken()) {
        request.headers.set('Authorization', 'Bearer ' + auth.getToken());
    }
    // continue to next interceptor
    next(function(response) {

        if (response.status == 401) {
            app.stopTimer();
            $('#login').modal('show');
        }

    });
});