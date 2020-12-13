console.log("App init")

function htmlAlert(text, class_) {
    alerts.innerHTML += `<div class="alert alert-${class_} alert-dismissible fade show" role="alert">
    ${text}
    <button type="button" class="close" data-dismiss="alert" aria-label="Close">
      <span aria-hidden="true">&times;</span>
    </button>
  </div>`
}

function showMain() {
    chrome.storage.sync.get(["mon_status"], function(result) {
        let monStatus = undefined
        
        if (result.mon_status) {
            monStatus = "on"
        } else {
            monStatus = "off"
        }

        let onOffButtonHtml = ""
    
        if (monStatus === "on") {
            onOffButtonHtml = '<button id="onOffButton" style="margin: 50px; width: 150px; height: 150px; border-radius: 75px; border-style: solid; border-width: 5px; border-color: black; padding: 20px 32px; text-align: center; font-size: 20px;" class="btn btn-danger">Off</button>'
        } else {
            onOffButtonHtml = '<button id="onOffButton" style="margin: 50px; width: 150px; height: 150px; border-radius: 75px; border-style: solid; border-width: 5px; border-color: black; padding: 20px 32px; text-align: center; font-size: 20px;" class="btn btn-success">On</button>'
        }
        
        app.innerHTML = `<div class="card" style="width: 18rem;">
        <div class="card-body">
          <h5 class="card-title">Мониторинговый агент <b id="monStatus">${monStatus}</b></h5>
          ${onOffButtonHtml}
          <button id="logOutButton" class="btn btn-danger">Выйти</button>
        </div>
        </div>`
        
        document.getElementById("logOutButton").addEventListener("click", handleLogout)
        document.getElementById("onOffButton").addEventListener("click", handleMonButton)
    })
}

function showConnect() {
    app.innerHTML = `<div class="card" style="width: 18rem;">
    <div class="card-body">
      <h5 class="card-title">Hello Mood!</h5>
      <p class="card-text">Вы можете подключить это приложение к Mood на<a href="https://mood.com/connect">главной странице!</a></p>
      <form id="authForm">
        <div class="form-group">
          <label for="authToken">Авторизационный токен с сайта</label>
          <input type="text" name="authToken" class="form-control" id="authToken" placeholder="Токен">
        </div>
        <button type="submit" class="btn btn-primary">Подключить</button>
      </form>
    </div>
    </div>`
    document.getElementById("authForm").addEventListener("submit", handleAuth)
}

function handleAuth(event) {
    let token = document.getElementById("authToken").value
    axios.post("https://mood.fflab.co/api/agent/confirm", {
        agent_token: token
    }).then(function (response) {
        console.log(response)
        chrome.storage.sync.set({"mood_token": token}, function () {
            chrome.storage.sync.set({"mon_status": false}, function () {
                showMain()
                htmlAlert("Login successful!", "success")
            })
        })
    }).catch(function(error) {
        document.getElementById("authToken").value = ""
        htmlAlert(error.response.data.error, "warning")
    })
    event.preventDefault()
}

function handleLogout(event){
    chrome.storage.sync.set({"mood_token": ""}, function () {
        showConnect()
    })
    event.preventDefault()
}

function handleMonButton(event) {
    chrome.storage.sync.get(["mon_status"], function(result) {
        chrome.storage.sync.set({"mon_status": !result.mon_status}, function() {
            showMain()
        })
    })
}

let app = document.getElementById("app")
let alerts = document.getElementById("alerts")
chrome.storage.sync.get(["mood_token"], function (result) {
    result.mood_token != "" ? showMain() : showConnect()
})