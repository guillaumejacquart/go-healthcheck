var app = new Vue({
    el: '#app',
    data: {
        apps: [],
        newApp: {
            pollTime: 5
        },
        search: '',
        checkType: 0,
        statusCode: 200
    },
    computed:
    {
        filteredApps:function()
        {
            var self=this;
            return this.apps.filter(function(app){
                return app.name.toLowerCase().indexOf(self.search.toLowerCase())>=0
                || app.url.toLowerCase().indexOf(self.search.toLowerCase())>=0;
            });
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
        getHistory: function (app) {
            var that = this;
            if(that.history && that.history.app_id === app.ID){
                that.history = false;
                return;
            }

            fetch("../apps/" + app.ID + "/history").then(function (response) {
                var contentType = response.headers.get("content-type");
                if (contentType && contentType.indexOf("application/json") !== -1) {
                    return response.json().then(function (json) {
                        that.history = {
                            app_id: app.ID,
                            items: json
                        };
                    });
                }
            });
        },
        saveApp: function () {
            var url = this.newApp.isUpdate ? ('../apps/' + this.newApp.ID) : '../apps';
            var method = this.newApp.isUpdate ? 'PUT' : 'POST'
            fetch(url, {
                method: method,
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    name: this.newApp.name,
                    url: this.newApp.url,
                    pollTime: parseInt(this.newApp.pollTime),
                    checkType: parseInt(this.newApp.checkType),
                    statusCode: parseInt(this.newApp.statusCode),
                    notify: this.newApp.notify
                })
            }).then(function (res) {
                if (res.status == 200) {
                    $('#add-app').modal('hide');
                }
            }, function (e) {
                alert("Error submitting form!");
            });
        },
        deleteApp: function (app) {
            if (!confirm('Are you sure you want to delete this app ?')) {
                return;
            }

            fetch("../apps/" + app.ID, {
                method: "DELETE"
            }).then(function (res) {
            }, function (e) {
                alert("Error submitting form!");
            });
        },
        updateForm: function (app) {
            this.newApp = app;
            this.newApp.isUpdate = true;
            $('#add-app').modal('show');
        },
        resetForm: function () {
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
