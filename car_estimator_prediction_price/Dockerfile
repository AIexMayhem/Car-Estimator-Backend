FROM python:3.12-slim

ENV PYTHONUNBUFFERED=1

WORKDIR /app

COPY requirements.txt .

RUN python -m pip install --default-timeout=200 --no-cache-dir -r requirements.txt

COPY server.py prediction_pb2_grpc.py prediction_pb2.py ./
COPY predictor/ predictor/

EXPOSE 50051
CMD ["python", "server.py"]
