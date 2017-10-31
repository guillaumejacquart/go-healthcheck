var app = new Vue({
    el: '#app',
    data: {
        apps: [],
        newApp: {}
    },
    created: function() {
        this.getApps()
        window.setInterval(this.getApps, 2000)
    },
    methods: {
        getApps: function() {
            var that = this;
            fetch("../apps").then(function(response) {
                var contentType = response.headers.get("content-type");
                if (contentType && contentType.indexOf("application/json") !== -1) {
                    return response.json().then(function(json) {
                        that.apps = json;
                    });
                }
            });
        },
        saveApp: function() {
            var url = this.newApp.isUpdate ? ('../apps/' + this.newApp.name) : '../apps';
            var method = this.newApp.isUpdate ? 'PUT' : 'POST'
            fetch(url, {
                method: method,
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.newApp
            }).then(function(res) {
                if (res.status == 200) {
                    alert("App has been saved !");
                }
            }, function(e) {
                alert("Error submitting form!");
            });
        },
        deleteApp: function(name) {
            fetch("../apps/" + name, {
                method: "DELETE"
            }).then(function(res) {
                if (res.status == 200) {
                    alert("App has been deleted !");
                }
            }, function(e) {
                alert("Error submitting form!");
            });
        },
        updateForm: function(app) {
            this.newApp = app;
            this.newApp.isUpdate = true;
            $('#add-app').modal('show');
        },
    }
})
