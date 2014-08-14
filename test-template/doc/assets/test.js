;(function() {

function testFunc(){
	alert("testFunc");
}

//	全局参数
this.test = testFunc;

}).call(function() {
	// search javascript call
	// 让(window) 继承闭包函数，这也外部就可以通过全局函数this.test()进行调用。testFunc() 外部是无法访问的。
	return this || (typeof window !== 'undefined' ? window : global);
}());

/*
	<script src="assets/test.js"></script>
	<script>
	this.test();
	</script>
*/