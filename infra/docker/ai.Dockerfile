FROM python:3.12-alpine

WORKDIR /workspace/services/ai
COPY services/ai ./
EXPOSE 8090
CMD ["python", "-m", "nexus_ai"]

