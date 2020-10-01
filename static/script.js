let Tasks = [];
let Toasts = [];
let App = {};

window.onload = function () { 
	App = new Vue({
	  i18n: i18n,
	  el: '#app',
	  data: {
		  version: "",
		  tasks: Tasks,
		  taskStatusColor: ["gray", "deepskyblue", "green", "red", "yellow"],
		  toasts: Toasts
	  },
	  methods: {
		  alert: function(m) {
			  alert(m);
		  },
		  deleteTask: function(task) {
			if (confirm(i18n.t("message.confirmDeleteTask", {"name": task.name, "address": task.address}))) {
			  	webAnnie.deleteTask(task.id);
			  }
		  },
		  getTaskStatusMessage: function(status) {
			const task_status = ["message.statusWait", "message.statusProgress", "message.statusComplete", "message.statusFail", "message.statusCancel"];
			return i18n.t(task_status[status]);
		  },
	  }
	});

	setInterval(function() {
		axios.get("./tasks").then(function (response) {
			let t = response.data;
			Tasks.splice(0, Tasks.length)
			let i = 0;
			for (let id in t) {
				t[id].id = id;
				Tasks.push(t[id]);
				i++;
			}
		}).catch(function (error) {
			webAnnie.addToast("alert-danger", i18n.t("message.errorRefreshTask", {"msg": error.message}));
		});
	}, 1000);
	axios.get("./version").then(function (response) {
		App.version = response.data.tag + " (" + response.data.date +  ")";
	});
};

var webAnnie = {
	addTask: function() {
		axios.post("./tasks", {
			address: prompt(i18n.t("message.promptTaskAddress"))
		}).then(function (response) {
			webAnnie.addToast("alert-success", i18n.t("message.infoAddSuccess"));
		}).catch(function (error) {
			webAnnie.errorToToast(error);
		});
	},
	deleteTask: function(id) {
		axios.delete("./tasks/" + String(id)).then(function (response) {
			webAnnie.addToast("alert-success", i18n.t("message.infoDeleteSuccess"));
		}).catch(function (error) {
			webAnnie.errorToToast(error);
		});
	},
	deleteCompleteTask: function() {
		if (!confirm(i18n.t("message.confirmDeleteCompletedTask"))) {
			return;
		}
		axios.delete("./tasks/complete").then(function (response) {
			let res = response.data;
			webAnnie.addToast("alert-success", i18n.t("message.infoDeleteCompletedTask", {"count": res.count}));
		}).catch(function (error) {
			webAnnie.errorToToast(error);
		});
	},
	errorToToast: function(error) {
		let msg = "";
		if(error.response) {
			switch (error.response.status) {
				case 400:
					msg = i18n.t("message.error400");
					break;
				case 500:
					msg = i18n.t("message.error500");
					break;
				default:
					msg = i18n.t("message.errorUnknown", {"status": error.response.status,"data": error.response.data});
					break;
			}
		} else {
			msg = "요청 중 오류";
		}
		webAnnie.addToast("alert-danger", i18n.t("message.errorGeneral", {"msg": msg}));
	},
	addToast: function (color_class, message) {
		let i = Toasts.push({"class": color_class, "message": message, "visible": true});
		setTimeout(function() {
			Vue.set(App.toasts, i - 1, {"class": color_class, "message": message, "visible": false});
		}, 2000);
	}
}



