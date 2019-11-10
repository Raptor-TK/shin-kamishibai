/* jshint asi: true */
/* jshint esversion: 3 */

/*!
License: refer to LICENSE file
 */

// function to format the hash to object
// example: #book=abc.zip&page=7 -> $['book']='abc.zip', $['page']=7
function hashParamGet(key) {
	var hashParams = {};
	var e,
		a = /\+/g, // Regex for replacing addition symbol with a space
		r = /([^&;=]+)=?([^&;]*)/g,
		// d = function (s) { return decodeURIComponent(s.replace(a, " ")); },
		d = function(s) {
			return decodeURIComponent(s);
		},
		q = window.location.hash.substring(1);

	while ((e = r.exec(q))) hashParams[d(e[1])] = d(e[2]);

	if (key) {
		return hashParams[key];
	}

	return hashParams;
}

// replace hash without adding history
function hashParamSet(key, val) {
	var clone = JSON.parse(JSON.stringify(hashParamGet()));

	clone[key] = val;

	var rp = window.location.protocol + "//" + window.location.host + window.location.pathname + window.location.search + "#";

	for (var k in clone) {
		rp += k + "=" + encodeURIComponent(clone[k]) + "&";
	}
	rp = rp.substring(0, rp.length - 1);

	// replace without adding history
	window.location.replace(rp);
}

// return fully formatted hash
function fullhash(page) {
	return "book=" + hashParamGet("book") + "&page=" + page;
}

function uport() {
	// get the url port
	var port = window.location.port || "80";
	return port;
}

// set cookie with prepending port number for key
function cookiep(name, setVal) {
	name = uport() + "." + name;
	return cookie(name, setVal);
}

function cookie(name, setVal) {
	var parseCookie = function(str) {
		var items = str.split("; ");
		var obj = {};

		for (var i = 0; i < items.length; i++) {
			var item = items[i];

			var pos = item.indexOf("=");

			var key = item.slice(0, pos);
			var val = decodeURIComponent(item.slice(pos + 1));

			obj[key] = val;
		}

		return obj;
	};

	var setCookie = function(skey, sval) {
		// max-age is 1 year
		document.cookie = skey + "=" + encodeURIComponent(sval) + "; max-age=31536000";
	};

	var cObj = parseCookie(document.cookie);

	if (arguments.length === 0) {
		return cObj;
	}
	if (arguments.length === 1) {
		return cObj[name];
	}

	return setCookie(name, setVal);
}

/*
  full screen functions
*/
function goFullScreen(i) {
	var elem;

	// if out what i is
	if (typeof i == "object" || i instanceof Object) {
		// i is a DOM element
		elem = i;
	} else if (typeof i == "string" || i instanceof String) {
		// i is an ID of DOM element
		elem = document.getElementById(i);
	} else {
		alert("goFullScreen(): unknown i");
	}

	// go full screen
	if (elem.mozRequestFullScreen) {
		elem.mozRequestFullScreen();
	} else if (elem.webkitRequestFullScreen) {
		elem.webkitRequestFullScreen();
	} else {
		alert("cannot go full screen");
	}
}

function fullScreenSupported() {
	if (document.mozCancelFullScreen) return true;
	if (document.webkitCancelFullScreen) return true;
	return false;
}

function exitFullScreen() {
	if (document.mozCancelFullScreen) {
		document.mozCancelFullScreen();
	} else if (document.webkitCancelFullScreen) {
		document.webkitCancelFullScreen();
	} else {
		alert("cannot exit full screen");
	}
}

function toggleFullScreen(id) {
	if (isFullScreen()) {
		exitFullScreen();
	} else {
		goFullScreen(id);
	}
}

function isFullScreen() {
	if (typeof document.mozFullScreen != "undefined") {
		return document.mozFullScreen;
	} else if (typeof document.webkitIsFullScreen != "undefined") {
		return document.webkitIsFullScreen;
	} else if (screen.width == window.innerWidth && screen.height == window.innerHeight) {
		return true;
	} else {
		return false;
	}
}

function isImageCached(src) {
	var image = new Image();
	image.src = src;

	return image.complete;
}

// for performance on addEventListener
// ref https://developers.google.com/web/tools/lighthouse/audits/passive-event-listeners
// ref https://github.com/WICG/EventListenerOptions/pull/30
function isPassive() {
	var supportsPassiveOption = false;
	try {
		addEventListener(
			"test",
			null,
			Object.defineProperty({}, "passive", {
				get: function() {
					supportsPassiveOption = true;
				}
			})
		);
	} catch (e) {}
	return supportsPassiveOption;
}

// attempt to use passive listener
function tryPassiveListner() {
	if (isPassive()) {
		return {
			capture: false,
			passive: false
		};
	}
	return false;
}

function objectToFormData(obj, form, namespace) {
	var fd = form || new FormData();
	var formKey;

	for (var property in obj) {
		if (obj.hasOwnProperty(property)) {
			if (namespace) {
				formKey = namespace + "[" + property + "]";
			} else {
				formKey = property;
			}

			// if the property is an object, but not a File,
			// use recursivity.
			if (typeof obj[property] === "object" && !(obj[property] instanceof File)) {
				objectToFormData(obj[property], fd, property);
			} else {
				// if it's a string or a File object
				fd.append(formKey, obj[property]);
			}
		}
	}

	return fd;
}

// ajax helper
function ajaxGet(url, queries, callback, callbackFail) {
	return ajax(url, {
		get: queries,
		callback: callback,
		callbackFail: callbackFail
	});
}
function ajaxPost(url, data, callback, callbackFail) {
	return ajax(url, {
		post: JSON.stringify(data),
		callback: callback,
		callbackFail: callbackFail
	});
}

// good ol' fash ajax using xmlhttprequest
function ajax(url, parms) {
	parms = parms || {};
	var req = new XMLHttpRequest(),
		post = parms.post || null,
		get = parms.get || null,
		contentType = parms.contentType || null,
		callback = parms.callback || null,
		callbackFail = parms.callbackFail || null,
		timeout = parms.timeout || null;

	req.onreadystatechange = function() {
		if (req.readyState != 4) return;

		// Error
		if (req.status != 200 && req.status != 304) {
			if (callbackFail) callbackFail(req.responseText);
			return;
		}

		if (callback) callback(req.responseText);
	};

	if (post) {
		req.open("POST", url, true);
		// req.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
		if (contentType) {
			req.setRequestHeader("Content-type", contentType);
		} else {
			req.setRequestHeader("Content-type", "application/json");
		}
	} else {
		if (get) {
			var result = [];
			for (var key in get) {
				result.push(key + "=" + encodeURI(get[key]));
			}
			url += "?" + result.join("&");
		}
		req.open("GET", url, true);
	}

	req.setRequestHeader("X-Requested-With", "XMLHttpRequest");

	req.send(post);

	if (timeout) {
		setTimeout(function() {
			req.onreadystatechange = function() {};
			req.abort();
			if (callbackFail) callbackFail(false);
		}, timeout);
	}

	return req;
}

// remember screen size
function setScreenSize() {
	cookie("width", window.innerWidth);
	cookie("height", window.innerHeight);
}

function checkLogin(callbackSuccess, callbackFail) {
	ajaxGet("/api/check", null, callbackSuccess, callbackFail);
}

function doLogin() {
	var password = document.getElementById("password").value;

	var callback = function(txt) {
		if (queryParams("referer") !== undefined) {
			// go to page
			var referer = queryParams("referer");
			window.location.href = window.location.origin + referer + window.location.hash;
			// var z = window.location.origin + referer + window.location.hash;
			// console.log(z);
		} else if (window.location.pathname === "/login.html") {
			window.location.pathname = "/browse.html";
		} else {
			// reload page without history
			location.reload();

			// window.location.replace(window.location.pathname + window.location.search + window.location.hash);
		}
	};

	var callbackFail = function(txt) {
		alert("invalid password or error");
	};

	var params = {
		post: JSON.stringify({
			password: password
		}),
		callback: callback,
		callbackFail: callbackFail
	};

	ajax("/login", params);
}

// alternative to URLSearchParams
function queryParams(key) {
	var hash = window.location.search
		.substr(1)
		.split("&")
		.reduce(function(q, query) {
			var chunks = query.split("=");
			var key = chunks[0];
			var value = chunks[1];
			return (q[key] = value), q;
		}, {});

	if (key !== undefined) {
		return hash[key];
	} else {
		return hash;
	}
}

// // set value, incomplete
// function queryParamSet(key, val) {
// 	var clone = JSON.parse(JSON.stringify(queryParams()));

// 	clone[key] = val;
// }
