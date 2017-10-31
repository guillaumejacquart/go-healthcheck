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
        addApp: function() {
            fetch("../apps", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.newApp
            }).then(function(res) {
                if (res.status == 200) {
                    alert("Perfect! Your settings are saved.");
                }
            }, function(e) {
                alert("Error submitting form!");
            });
        }
    }
})