FROM python:3.10

WORKDIR /app
COPY ./fastapi/requirements.txt .
RUN pip install --no-cache-dir --upgrade -r /app/requirements.txt

CMD ["uvicorn", "main:app", "--reload", "--host", "0.0.0.0", "--port", "8080"]
