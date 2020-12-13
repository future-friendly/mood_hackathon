MutationObserver = window.MutationObserver || window.WebKitMutationObserver;

var observer = new MutationObserver(function (mutations, observer) {
    chrome.runtime.sendMessage({
        url: window.location.href,
        pageSource: document.documentElement.innerHTML,
        timestamp: new Date().valueOf()
    });
    timeout = true;
});

observer.observe(document, {
	subtree: true,
	attributes: true
});
