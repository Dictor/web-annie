let Tasks = [];
let App = {};

window.onload = function () { 
	App = new Vue({
	  el: '#app',
	  data: {
		  tasks: Tasks,
		  taskStatus: ["대기", "진행", "완료", "실패"],
		  taskStatusColor: ["gray", "deepskyblue", "green", "red"]
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
		});
	}, 1000)

};

function addTask() {
	axios.post("./task", {
		address: prompt('추가할 동영상의 주소를 입력하세요.')
  	})
}