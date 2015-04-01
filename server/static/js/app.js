var SmoothApp = angular.module('SmoothApp', []);

SmoothApp.controller('SmoothCtrl', ['$scope', 
	function($scope) {
		$scope.messages = [];
		$scope.clients = {};
		conn = new WebSocket("ws://"+location.host+"/ws");
		conn.onclose = function(evt) {
			console.log("Connection closed");
		};
		conn.onmessage = function(evt) {
			$scope.$apply(function() {
				console.log(evt.data);
				var d = JSON.parse(evt.data);
				if (d.type == "itr") {
					$scope.clients[d.payload.hostname] = d.payload;	
				} else if(d.type == "ini") {
					$scope.clients = d.payload;
				}
			});
		};
		conn.onopen = function(evt) {
			console.log("Connection established.");
		};
	}]);
