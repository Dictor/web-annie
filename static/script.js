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
	/*
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
	*/
	
	Vue.set(App.tasks, 0, {"address":"https://www.youtube.com/watch?v=s-wdegCT3qY","status":1,"info":"\n Site:      YouTube youtube.com\n Title:     구워먹는 삼겹살에 질린 여러분을 위한 영상\n Type:      video\n Stream:   \n     [303]  -------------------\n     Quality:         1080p60 video/webm; codecs=\"vp9\"\n     Size:            208.73 MiB (218868305 Bytes)\n     # download with: annie -f 303 ...\n\n\r","progress":{"total":"208.73 MiB","current":"19.98 MiB","speed":"10 MiB/s","percentage":"10%","time_left":"1m"}});
};

function addTask() {
	axios.post("./task", {
		address: prompt('추가할 동영상의 주소를 입력하세요.')
  	})
}