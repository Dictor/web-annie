let Tasks = [];
let Toasts = [];
let App = {};

window.onload = function () { 
	App = new Vue({
	  el: '#app',
	  data: {
		  tasks: Tasks,
		  taskStatus: ["대기", "진행", "완료", "실패", "취소"],
		  taskStatusColor: ["gray", "deepskyblue", "green", "red", "yellow"],
		  toasts: Toasts
	  },
	  methods: {
		  alert: function(m) {
			  alert(m);
		  },
		  deleteTask: function(task) {
			  if (confirm(task.name + " (" + task.address + ") 작업을 삭제합니까?")) {
			  	webAnnie.deleteTask(task.id);
			  }
		  }
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
			webAnnie.addToast("alert-danger", "갱신 오류 : " + error.message)
		});
	}, 1000)

};

var webAnnie = {
	addTask: function() {
		axios.post("./tasks", {
			address: prompt('추가할 동영상의 주소를 입력하세요.')
		}).then(function (response) {
			webAnnie.addToast("alert-success", "추가 성공!");
		}).catch(function (error) {
			let msg = "";
			if(error.response) {
				switch (error.response.status) {
					case 400:
						msg = "유효하지 않은 주소";
						break;
					case 500:
						msg = "서버 내부 오류";
						break;
					default:
						msg = "예기치 못한 오류 = " + error.response.status + " " + error.response.data;
						break;
				}
			} else {
				msg = "요청 중 오류";
			}
			webAnnie.addToast("alert-danger", "추가 오류 : " + msg);
		});
	},
	deleteTask: function(id) {
		axios.delete("./tasks/" + String(id)).then(function (response) {
			webAnnie.addToast("alert-success", "삭제 성공!");
		}).catch(function (error) {
			let msg = "";
			if(error.response) {
				switch (error.response.status) {
					case 400:
						msg = "유효하지 않은 ID";
						break;
					case 500:
						msg = "서버 내부 오류";
						break;
					default:
						msg = "예기치 못한 오류 = " + error.response.status + " " + error.response.data;
						break;
				}
			} else {
				msg = "요청 중 오류";
			}
			webAnnie.addToast("alert-danger", "삭제 오류 : " + msg);
		});
	},
	addToast: function (color_class, message) {
		let i = Toasts.push({"class": color_class, "message": message, "visible": true});
		setTimeout(function() {
			Vue.set(App.toasts, i - 1, {"class": color_class, "message": message, "visible": false});
		}, 2000);
	}
}



