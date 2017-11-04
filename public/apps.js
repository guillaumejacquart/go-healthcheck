var app = new Vue({
    el: '#app',
    data: {
        apps: [],
        newApp: {
            pollTime: 5
        }
    },
    created: function () {
        this.getApps()
        window.setInterval(this.getApps, 2000)
    },
    methods: {
        getApps: function () {
            var that = this;
            fetch("../apps").then(function (response) {
                var contentType = response.headers.get("content-type");
                if (contentType && contentType.indexOf("application/json") !== -1) {
                    return response.json().then(function (json) {
                        that.apps = json;
                    });
                }
            });
        },
        saveApp: function () {
            var url = this.newApp.isUpdate ? ('../apps/' + this.newApp.name) : '../apps';
            var method = this.newApp.isUpdate ? 'PUT' : 'POST'
            fetch(url, {
                method: method,
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(this.newApp)
            }).then(function (res) {
                if (res.status == 200) {
                    alert("App has been saved !");
                    $('#add-app').modal('hide');
                }
            }, function (e) {
                alert("Error submitting form!");
            });
        },
        deleteApp: function (app) {
            fetch("../apps/" + app.name, {
                method: "DELETE"
            }).then(function (res) {
                if (res.status == 200) {
                    alert("App has been deleted !");
                }
            }, function (e) {
                alert("Error submitting form!");
            });
        },
        updateForm: function (app) {
            this.newApp = app;
            this.newApp.isUpdate = true;
            $('#add-app').modal('show');
        },
        resetForm: function(){            
            this.newApp = {
                pollTime: 5
            };
            $('#add-app').modal('show');
        },
        formatUpdate: function (date) {
            if (date.indexOf('0001') === 0) {
                return 'Never'
            }
            return new Date(date).toUTCString();
        },
        getIcon: function (app) {
            var result = [];
            switch (app.status) {
                case 'up':
                    result.push('oi-circle-check')
                    break;
                case 'down':
                    result.push('oi-x')
                    break;
            }
            return result;
        }
    }
})
