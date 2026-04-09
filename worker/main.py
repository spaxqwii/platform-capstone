from fastapi import FastAPI
from celery import Celery
import os

app = FastAPI()

# Celery for background jobs (Redis as broker)
celery_app = Celery(
    "worker",
    broker=os.getenv("REDIS_URL", "redis://localhost:6379/0"),
    backend=os.getenv("REDIS_URL", "redis://localhost:6379/0")
)

@app.get("/health")
def health():
    return {"status": "ok", "service": "worker"}

@app.get("/ready")
def ready():
    return {"ready": True}

@celery_app.task
def process_job(data: str):
    """Background job placeholder"""
    return f"Processed: {data}"

@app.post("/job")
def create_job(data: str):
    """Trigger a background job"""
    task = process_job.delay(data)
    return {"task_id": task.id, "status": "queued"}

@app.get("/job/{task_id}")
def get_job(task_id: str):
    """Check job status"""
    result = celery_app.AsyncResult(task_id)
    return {"task_id": task_id, "status": result.status, "result": result.result}

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8080"))
    uvicorn.run(app, host="0.0.0.0", port=port)
