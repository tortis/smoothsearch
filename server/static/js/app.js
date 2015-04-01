var SmoothApp = angular.module('SmoothApp', ['ngAnimate']);

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
				} else if (d.type == "hit") {
                    $scope.clients[d.payload.hostname].smoothness = d.payload.smoothness;
                    $scope.clients[d.payload.hostname].smooth_num = d.payload.smooth_num;
                }

			});
		};
		conn.onopen = function(evt) {
			console.log("Connection established.");
		};
	}]);

SmoothApp.directive('animateOnChange', ['$animate', '$timeout', function($animate, $timeout) {
  return function(scope, elem, attr) {
    scope.$watchCollection(attr.animateOnChange, function() {
      $animate.addClass(elem, 'flashing-on').then(function() {
        $timeout(function(){
          $animate.removeClass(elem, 'flashing-on');
        }, 0);
      });
    });
  };
}]);
