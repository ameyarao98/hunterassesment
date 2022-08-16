FROM python:3.10-slim

WORKDIR /app

COPY pyproject.toml poetry.lock ./
RUN pip3 install poetry
RUN poetry config virtualenvs.create false
RUN poetry install

COPY src ./
CMD ["python","-m","sanic","server.app","--host=0.0.0.0"]