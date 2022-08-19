FROM python:3.10-slim

WORKDIR /app

COPY pyproject.toml poetry.lock ./
RUN pip3 install poetry
RUN poetry config virtualenvs.create false
RUN poetry install

COPY server.py ./
CMD ["sanic","server.app","--host=0.0.0.0","--port=8000"]