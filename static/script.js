let Tasks = [];
let App = {};

window.onload = function () { 
	App = new Vue({
	  el: '#app',
	  data: {tasks: Tasks}
	});
	setInterval(function() {
		axios.get("./tasks").then(function (response) {
			let t = response.data;
			Tasks = [];
			let i = 0;
			for (let id in t) {
				//Tasks.push(t[id]);
				Vue.set(App.tasks, i, t[id]);
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