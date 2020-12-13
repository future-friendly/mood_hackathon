import itertools
from collections import Counter
from urllib.parse import urlparse

def to_sessions(frames, domain=False):
    frames.sort(key=lambda x: x["timestamp"])
    if domain:
        sessions_raw = [list(g) for _, g in itertools.groupby(frames, lambda x: urlparse(x['url']).netloc)]
    else:
        sessions_raw = [list(g) for _, g in itertools.groupby(frames, lambda x: x['url'])]
    
    sessions = []
    for sample in sessions_raw:
        temp = {
            'category': sample[0]['category'],
            'from': sample[0]['timestamp'],
            'to': sample[-1]['timestamp'],
            'url': sample[0]['url'],
            'keyword_distribution': compress_keywords([i["keywords"].split("::") for i in sample])
        }
        sessions.append(temp)
    return sessions

def compress_keywords(keywords_list):
    total = []
    for i in keywords_list:
        total.extend(i)
    
    common = Counter(total).most_common()
    return { word: count / len(total) for word,count in common }
    
def get_interest_map(sessions):
    return dict(Counter([i["category"] for i in sessions]).most_common())

def get_category_map(sessions, full=False):
    category_map = {}
    for session in sessions:
        for word, percent in session["keyword_distribution"].items():
            if not category_map.get(word):
                category_map[word] = percent
            else:
                category_map[word] += percent

    return category_map

