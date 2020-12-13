from aiohttp import web
import tagger
import requests
import aggregation
from urllib.parse import urlparse
import re
from multiprocessing import Queue, Process

BASE_BACKEND = "http://backend:8080"
IGNORE_DOMAINS = [r"localhost:\d+", "localhost", "mood.fflab.co"]

page_queue = Queue()

def is_ignored(req):
    for regexp in IGNORE_DOMAINS:
        if re.match(regexp, urlparse(req["url"]).netloc):
            return True
    return False
                

def process_page(q):
    while True:
        req = q.get()
        print("Processing", req["url"])
        processed = tagger.tag(req)
        if not is_ignored(req):
            requests.post(BASE_BACKEND + "/data/newpage", json={
                "agent_token": req["agent_token"],
                "category": processed["category"],
                "url": processed["url"],
                "keywords": processed["keywords"],
                "timestamp": processed["timestamp"]
                })

async def handle_page(request):
    req = await request.json()
    for regexp in IGNORE_DOMAINS:
        if re.match(regexp, urlparse(req["url"]).netloc):
            return web.json_response({"ok": True})
    
    page_queue.put(req)
    return web.json_response({"ok": True})

async def handle_interest_chart(request):
    req = await request.json()
    if req.get("data"):
        try:
            sessions = aggregation.to_sessions(req["data"])
            interest_map = aggregation.get_interest_map(sessions)
        except Exception as e:
            print(e)
            return web.json_response({"error": str(e)})
        return web.json_response({"result": interest_map})
    else:
        return web.json_response({"error": "no data provided"})

async def handle_keyword_chart(request):
    req = await request.json()
    if req.get("data"):
        try:
            sessions = aggregation.to_sessions(req["data"])
            keyword_map = aggregation.get_category_map(sessions)
        except Exception as e:
            return web.json_response({"error": str(e)})
        return web.json_response({"result": keyword_map})
    else:
        return web.json_response({"error": "no data provided"})

app = web.Application()
app.add_routes([
    web.post("/page", handle_page),
    web.post("/interest_chart", handle_interest_chart),
    web.post("/keyword_chart", handle_keyword_chart)
])

if __name__ == "__main__":
    p = Process(target=process_page, args=(page_queue,))
    p.start()
    web.run_app(app, port=8080)