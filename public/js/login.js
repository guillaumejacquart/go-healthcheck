var Login = new Vue({
    el: '#login',
    data: {
        username: '',
        password: '',
        error: ''
    },
    methods: {
        login: function() {
            var url = '../login';

            var data = {
                username: this.username,
                password: this.password,
            }

            var options = {
                headers: {
                    "Content-Type": "application/json"
                }
            }

            this.$http.post(url, data, options).then(function(res) {
                    if (res.status == 200) {
                        return res.json().then(function(json) {
                            $('#login').modal('hide');
                            auth.setToken(json.token);
                            app.startTimer();
                        });
                    }
                },
                function(e) {
                    alert("Error submitting form!");
                });
        }
    }
})
