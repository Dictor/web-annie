let Tasks = [];
let Toasts = [];
let App = {};

window.onload = function () { 
	App = new Vue({
	  el: '#app',
	  data: {
		  tasks: Tasks,
		  taskStatus: ["대기", "진행", "완료", "실패"],
		  taskStatusColor: ["gray", "deepskyblue", "green", "red"],
		  toasts: Toasts
	  },
	  methods: {
		  alert: function(m) {
			  alert(m);
		  }
	  }
	});
	setInterval(function() {
		axios.get("./tasks").then(function (response) {
			let t = response.data;
			Tasks.splice(0, Tasks.length)
			let i = 0;
			for (let id in t) {
				Tasks.push(t[id]);
				i++;
			}
		}).catch(function (error) {
			addToast("alert-danger", "갱신 오류 :" + error)
		});
	}, 1000)

};

function addTask() {
	axios.post("./task", {
		address: prompt('추가할 동영상의 주소를 입력하세요.')
	}).then(function (response) {
		addToast("alert-success", "추가 성공!");
	}).catch(function (error) {
		switch (error.response.status) {
			case 400:
				addToast("alert-danger", "추가 오류 : 유효하지 않은 주소");
				break;
			case 500:
				addToast("alert-danger", "추가 오류 : 서버 내부 오류");
				break;
		}
	});
}

function addToast(color_class, message) {
	let i = Toasts.push({"class": color_class, "message": message, "visible": true});
	setTimeout(function() {
		Vue.set(App.toasts, i - 1, {"class": color_class, "message": message, "visible": false});
	}, 2000);
}