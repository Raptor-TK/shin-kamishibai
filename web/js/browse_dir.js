/* jshint asi: true */
/* jshint esversion: 3 */

/*!
License: refer to LICENSE file
 */

// load dir sources
function sourcesReload(reloadFromCache) {
	// internal func to rebuild dom
	var rebuild = function(srcs) {
		var ul = document.getElementById("div-sources");

		// remove all child
		while (ul.hasChildNodes()) {
			ul.removeChild(ul.lastChild);
		}

		for (var i = 0; i < srcs.length; i++) {
			var a = document.createElement("a");
			a.href = "#";
			a.setAttribute("srcNum", i);
			a.onclick = aSourceSelect;
			a.innerText = srcs[i];

			ul.appendChild(a);
		}
	};

	// just reload locally, no retrive from server
	if (reloadFromCache === true) {
		rebuild(dirSources);
		return;
	}

	ajaxGet("/api/list_sources", {}, function(dat) {
		// update global sources
		dirSources = JSON.parse(dat);

		rebuild(dirSources);
	});
}

function dirOrderBy(str) {
	window.localStorage.orderBy = str;

	dirListReload(dirPath, document.getElementById("searchbox").value, dirPage);
}

// used for construct dir listing from data
// returns file/dir elements
function dirParseList(files) {
	// not files
	if (files instanceof Array === false) {
		console.error("input not array");
		return;
	}

	// return value
	var fragment = document.createDocumentFragment();

	var item, a, img, text, span;

	// first block contains info
	var dirInfo = files[0];

	for (var i = 0; i < files.length; i++) {
		// hack, skip first one
		if (i === 0) continue;

		var file = files[i];

		// skip dot file
		if (file.name && file.name[0] === ".") continue;

		var full_path = dirInfo.path + "/" + file.name;

		if (file.is_dir) {
			// dir

			var icon = "folder-mini.png";

			if (file.name === "Trash") {
				// trash
				// icon = "folder-trash.png";
			}

			item = document.createElement("div");
			item.className = "directory";

			a = document.createElement("a");
			a.setAttribute("dir", full_path);
			a.href = "#dir=" + full_path + "&page=1";
			a.onclick = dirClicked;

			img = document.createElement("img");
			img.src = "/images/" + icon;
			img.alt = "folder";

			text = document.createElement("div");
			text.className = "text";
			text.innerText = file.name;

			item.appendChild(a);
			a.appendChild(img);
			a.appendChild(text);

			fragment.appendChild(item);
		} else if (file.name) {
			// file

			var href = "";
			var readstate = "read";

			if (file.pages && file.page) {
				// file read
				// read5 10 20 30 40 ... 100

				var bn = ((1.0 * file.page) / file.pages) * 100;
				var pc = bn - (bn % 10);

				//  read percentage css class
				if (pc > 0) {
					readstate += pc;
				} else {
					readstate += "5";
				}

				href = "/read.html?book=" + file.id + "&page=" + file.page;
			} else {
				// unread
				readstate += "0";
				href = "/read.html?book=" + file.id;
			}

			item = document.createElement("div");
			item.className = "file";

			a = document.createElement("a");
			a.setAttribute("bookcode", file.id);
			a.href = href;
			a.onclick = bookClicked;

			img = document.createElement("img");
			img.src = "/api/thumbnail/" + file.id;
			img.alt = "book";

			text = document.createElement("div");
			text.className = "text " + readstate;
			text.innerText = file.name;

			span = document.createElement("span");
			span.className = "book-pages";
			span.innerHTML = file.pages;

			a.appendChild(img);
			a.appendChild(text);
			a.appendChild(span);

			item.appendChild(a);

			fragment.appendChild(item);
		} else if (file.more) {
			item = document.createElement("div");
			item.className = "directory";

			text = document.createElement("div");
			text.className = "text";
			text.innerText = "More...";

			item.appendChild(text);
			fragment.appendChild(item);
		}
	}

	// indicate eof or more of dir list
	if (files.length === 1) {
		item = document.createElement("div");
		item.className = "directory";

		text = document.createElement("div");
		text.className = "text";
		text.innerText = "EOF";

		item.appendChild(text);
		fragment.appendChild(item);
	}

	return fragment;
}

function dirListPrev() {
	// stop if dir not defined
	if (dirPath == undefined) {
		return;
	}

	// get keyword from searchbox
	var keyword = document.getElementById("searchbox").value;

	// get page
	if (typeof dirPage !== "number") {
		dirPage = 1;
	}
	if (isNaN(dirPage) || dirPage < 0) {
		dirPage = 1;
	}
	if (dirPage <= 1) {
		return;
	}

	dirPage = dirPage - 1;

	dirListReload(dirPath, keyword, dirPage);
}

function dirListNext() {
	// stop if dir not defined
	if (dirPath == undefined) {
		return;
	}

	// get keyword from searchbox
	var keyword = document.getElementById("searchbox").value;

	// get page
	if (typeof dirPage !== "number") {
		dirPage = 1;
	}
	if (isNaN(dirPage) || dirPage < 0) {
		dirPage = 0;
	}
	if (dirList.length <= 1) {
		return;
	}

	dirPage = dirPage + 1;

	dirListReload(dirPath, keyword, dirPage);
}

// selection from bookmark
function aSourceSelect(evt) {
	var srcNum = Number(this.getAttribute("srcNum"));
	dirListReload(dirSources[srcNum], "", 1);
}

// reload listing
function dirListReload(dir_path, keyword, page, loadFromCache) {
	// set default to name for order_by
	var order_by = "name";
	var co = window.sessionStorage.orderBy;
	if (co) {
		switch (co) {
			case "name":
			case "size":
			case "date":
				order_by = co;
				break;
		}
	}

	// load from cache
	if (loadFromCache === true) {
		updatePathLabel(dir_path);

		var els = dirParseList(dirList);
		if (!els) return;

		var el = document.getElementById("dir-lists");
		while (el.hasChildNodes()) {
			el.removeChild(el.lastChild);
		}

		el.appendChild(els);
		return;
	}

	// remember on cookie
	window.sessionStorage.lastPath = dir_path;
	window.sessionStorage.lastPage = page;
	window.sessionStorage.orderBy = order_by;

	// update values
	dirPath = dir_path;
	dirPage = page;
	document.getElementById("span-page").textContent = dirPage;
	// replace url without adding history
	hashParamSet("dir", dirPath);
	hashParamSet("page", dirPage);

	var el = document.getElementById("dir-lists");
	// delete all child
	while (el.hasChildNodes()) {
		el.removeChild(el.lastChild);
	}

	ajaxGet(
		"/api/lists_dir",
		{
			dir: dir_path,
			page: page,
			keyword: keyword
		},
		function(data) {
			dirList = JSON.parse(data);

			updatePathLabel(dir_path);

			var els = dirParseList(dirList);

			if (!els) {
				console.error("dirParseList() failed");
				return;
			}

			// add
			el.appendChild(els);

			// // get to the last selected item
			// var el_lsi = $('span:contains("' + window.sessionStorage.lastSelectedItem + '")').parent();
			// if (el_lsi.length == 1) {
			// 	$(el_lsi).addClass("last-selected-item");
			// 	$(document).scrollTo(el_lsi, { offset: -$(".navbar-inner").height() });
			// }

			// make sure files are deleteable if in delete mode
			if (isDeleteMode) {
				deleteEnable();
			}
		},
		function(dat) {
			// fail callback
			el.innerHTML = dat;
		}
	);
}

function dirUp() {
	var dirs = dirPath.split("/");
	dirs.pop();

	var dir = dirs.join("/");

	// make sure dir is in allowed bookmark or stop
	for (var i = 0; i < dirSources.length; i++) {
		var dirSS = dirSources[i];

		if (dir.includes(dirSS)) {
			dirListReload(dir, "", 1);
			return;
		}
	}
}

function dirClicked(evt) {
	dirListReload(this.getAttribute("dir"), "", 1);
}

function bookClicked(evt) {
	// console.log(777, evt, this.getAttribute("bookcode"), this.getAttribute("dir"));
	// window.sessionStorage.lastSelectedItem = this.innerText;
	// return false;
}