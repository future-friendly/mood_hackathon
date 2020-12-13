chrome.runtime.onMessage.addListener(
    function (request, sender, sendResponse) {
        chrome.storage.sync.get(["mood_token", "mon_status"], function (result) {
                if (result.mon_status) {
                    axios.post("https://mood.fflab.co/analytics/page", {
                        url: request.url,
                        pageSource: request.pageSource,
                        timestamp: request.timestamp,
                        agent_token: result.mood_token
                    }).then(function (response) {
                        console.log(response)
                    }).catch(function(error) {
                        console.log(error)
                    })
                }
        })
    }
)