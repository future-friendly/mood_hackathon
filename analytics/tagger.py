from bs4 import BeautifulSoup, Comment, Doctype, Declaration, CData, NavigableString
import re
from summa import keywords
from nltk.stem import WordNetLemmatizer
import pandas as pd
from nltk.corpus import stopwords
from nltk.tag import pos_tag

#URL Classification DB
df = pd.read_csv('final_classified.csv').drop(columns=['Unnamed: 0'])

def clean_html(html):
    soup = BeautifulSoup(html, 'lxml')

    [s.extract() for s in soup('script')]
    [s.extract() for s in soup('style')]
    [s.extract() for s in soup.find_all(string=lambda text: isinstance(text, Comment))]
    [s.extract() for s in soup.find_all(string=lambda text: isinstance(text, Doctype))]
    [s.extract() for s in soup.find_all(string=lambda text: isinstance(text, Declaration))]
    [s.extract() for s in soup.find_all(string=lambda text: isinstance(text, CData))]


    clean = ''

    for i in soup.children:
        if isinstance(i, NavigableString):
            clean += str(i)
        elif isinstance(i, Doctype):
            pass
        elif isinstance(i, Comment):
            pass
        else:
            clean += i.text
        
        clean += ' '

    return re.sub(r'\s+', ' ', clean)

def classify_page(url):
    domain = url.split('/')[2]
    # print(domain)
    resp = df[df['url'] == domain]
    
    if len(resp) != 0:
        website_class = list(resp['name'])[0]
    else:
        resp = df[df['url'] == '.'.join(domain.split('.')[-2:])]

        if len(resp) != 0:
            website_class = list(resp['name'])[0]
        else:
            website_class = 'Other'

    return website_class
    

def tag(data):
    lmtz = WordNetLemmatizer()
    url_class = classify_page(data['url'])

    clean_page = clean_html(data['pageSource'])
    keyword_count = int(len(clean_page.split(' ')) * 0.01)
    if keyword_count == 0:
        keyword_count = 1

    if keyword_count > 5:
        keyword_count = 5

    keywords_from_page = []
    for i in keywords.keywords(clean_page, words=keyword_count).split('\n'):
        if i != "":
            try:
                lemma = lmtz.lemmatize(i)
                if lemma not in keywords_from_page and lemma not in stopwords.words('english') and len(lemma) >= 3:
                    keywords_from_page.append(lemma)
            except Exception as e:
                print(e)

    return {"category": url_class, "keywords": keywords_from_page, 'timestamp': data['timestamp'], 'url': data['url']}