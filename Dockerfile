FROM python:3.8-alpine

COPY . /app

RUN pip install cloudflare==2.10.2

WORKDIR /app

CMD ["python", "ddns.py"]
