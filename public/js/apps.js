var app = new Vue({
    el: '#app',
    data: {
        apps: [],
        newApp: {
            pollTime: 5,
            headers: []
        },
        search: '',
        checkType: 0,
        statusCode: 200
    },
    computed: {
        filteredApps: function() {
            var self = this;
            return this.apps.filter(function(app) {
                return app.name.toLowerCase().indexOf(self.search.toLowerCase()) >= 0 ||
                    app.url.toLowerCase().indexOf(self.search.toLowerCase()) >= 0;
            });
        }
    },
    created: function() {
        this.getApps()
        window.setInterval(this.getApps, 2000)
    },
    methods: {
        getApps: function() {
            var that = this;
            that.$http.get('../apps').then(response => {
                return response.json().then(function(json) {
                    that.apps = json;
                });
            });
        },
        getHistory: function(app) {
            var that = this;
            if (that.history && that.history.app_id === app.ID) {
                that.history = false;
                return;
            }

            that.$http.get("../apps/" + app.ID + "/history").then(function(response) {
                return response.json().then(function(json) {
                    that.history = {
                        app_id: app.ID,
                        items: json
                    };
                });
            });
        },
        changeStatus: function(status, app){
            app.checkStatus = status
            this.save(app);
        },
        saveApp: function(){
            this.save(this.newApp);
        },
        save: function(app) {
            var url = app.ID ? ('../apps/' + app.ID) : '../apps';
            
            var data = {
                name: app.name,
                url: app.url,
                pollTime: parseInt(app.pollTime),
                checkType: parseInt(app.checkType),
                statusCode: parseInt(app.statusCode),
                notify: app.notify,
                headers: app.headers,
                checkStatus: app.checkStatus
            }
            
            var options = {
                headers: {
                    "Content-Type": "application/json"
                }
            }
            var promise = app.ID ? this.$http.put(url, data, options) : this.$http.post(url, data, options)
            promise.then(function(res) {
                    if (res.status == 200) {
                        $('#add-app').modal('hide');
                    }
                },
                function(e) {
                    alert("Error submitting form!");
                });
        },
        deleteApp: function(app) {
            if (!confirm('Are you sure you want to delete this app ?')) {
                return;
            }

            this.$http.delete("../apps/" + app.ID)
                .then(function(res) {}, function(e) {
                    alert("Error submitting form!");
                });
        },
        updateForm: function(app) {
            this.newApp = app;
            this.newApp.isUpdate = true;
            $('#add-app').modal('show');
        },
        resetForm: function() {
            this.newApp = {
                pollTime: 5,
                headers: []
            };
            $('#add-app').modal('show');
        },
        formatUpdate: function(date) {
            if (date.indexOf('0001') === 0) {
                return 'Never'
            }
            return new Date(date).toUTCString();
        },
        getIcon: function(app) {
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
