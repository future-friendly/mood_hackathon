from aiohttp import web
from aiohttp_session import get_session, setup
from aiohttp_session.cookie_storage import EncryptedCookieStorage
import requests

BACKEND_BASE = "http://backend:8080"

agent_types = {
    0: "Chrome Extension",
    1: "Android App",
    2: "iOS App"
}

async def index(request):
    session = await get_session(request)
    if session.get("auth_token"):
        return web.FileResponse("templates/index.html")
    else:
        raise web.HTTPFound("/login")

async def signup(request):
    return web.FileResponse("templates/signup.html")

async def login(request):
    return web.FileResponse("templates/login.html")

async def logout(request):
    session = await get_session(request)
    if session.get("auth_token"):
        token = session["auth_token"]
        session["auth_token"] = None
        return web.json_response(requests.post(BACKEND_BASE + "/auth/logout", json={"token": token}).json())

async def new_agent(request):
    session = await get_session(request)
    if session.get("auth_token"):
        data = await request.json()
        data["token"] = session["auth_token"]
        return web.json_response(requests.post(BACKEND_BASE + "/agent/add", json=data).json())

async def agents(request):
    session = await get_session(request)
    if session.get("auth_token"):
        token = session["auth_token"]
        resp = requests.post(BACKEND_BASE + "/agent/get", json={"token": token}).json()
        if resp.get("error"):
            return web.json_response(resp)
        for i in resp["agents"]:
            i["agent_type"] = agent_types[i["agent_type"]]
        return web.json_response(resp)

async def delete_agent(request):
    session = await get_session(request)
    if session.get("auth_token"):
        token = session["auth_token"]
        data = await request.json()
        data["token"] = token
        return web.json_response(requests.post(BACKEND_BASE + "/agent/delete", json=data).json())

async def get_interest_map(request):
    session = await get_session(request)
    if session.get("auth_token"):
        token = session["auth_token"]
        resp = requests.post(BACKEND_BASE + "/chart/get", json={"token": token, "chart_type": 0}).json()
        if resp.get("result"):
            resp = resp["result"]
            return web.json_response({"data": list(resp.values()), "labels": list(resp.keys())})
        else:
            return web.json_response(resp)

async def get_keyword_map(request):
    session = await get_session(request)
    if session.get("auth_token"):
        token = session["auth_token"]
        data = await request.json()
        resp = requests.post(BACKEND_BASE + "/chart/get", json={"token": token, "chart_type": 1, "category": data["category"]}).json()
        if resp.get("result"):
            resp = resp["result"]
            resp = dict(sorted(list(resp.items()), key=lambda x: x[1], reverse=True))
            top_keywords = dict(sorted(list(resp.items()), key=lambda x: x[1], reverse=True)[0:20])
            
            total = sum(resp.values())
            top_total = sum(top_keywords.values())
            other = ((total - top_total) / total) * 100

            for i in top_keywords:
                top_keywords[i] = (top_keywords[i] / total) * 100
            top_keywords["Other less significant"] = other
            return web.json_response({"data": [round(i,2) for i in top_keywords.values()], "labels": list(top_keywords.keys()), "category": data["category"],
                "full_data": [round(i,2) for i in resp.values()], "full_labels": list(resp.keys())
            })
        else:
            return web.json_response(resp)

async def register_token(request):
    data = await request.json()
    session = await get_session(request)
    session["auth_token"] = data["token"]
    return web.json_response({'ok': True})

def init():
    app = web.Application()
    setup(app,
        EncryptedCookieStorage(b'Thirty  two  length  bytes  key.'))
    app.add_routes([
        web.get("/", index),

        web.get("/login", login),
        web.post("/newtoken", register_token),
        web.get("/signup", signup),
        web.post("/logout", logout),

        web.get("/agents", agents),
        web.post("/newagent", new_agent),
        web.post("/deleteagent", delete_agent),

        web.post("/interestmap", get_interest_map),
        web.post("/keywordmap", get_keyword_map),

        web.static("/static", "./static")
    ])
    return app

if __name__ == "__main__":
    web.run_app(init(), port=8080)