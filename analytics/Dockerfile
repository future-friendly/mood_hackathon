FROM python:3.8

WORKDIR /usr/src/app
COPY requirements.txt requirements.txt
RUN  pip3 install -r requirements.txt

COPY install_nltk.py install_nltk.py
RUN python3 install_nltk.py

COPY . .

ENTRYPOINT [ "python3", "-u","api.py" ]