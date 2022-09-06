from fastapi import FastAPI
import schemas
import pika
import json


app = FastAPI()


# declare the publisher
def Publish_alert(alert: schemas.Alert):
    params = pika.URLParameters("amqps://tmwnunoc:EspITp82kgC9D7sVtaeJ4Cn2lF0QDnU4@shark.rmq.cloudamqp.com/tmwnunoc")
    conn = pika.BlockingConnection(parameters=params)
    channel = conn.channel()
    channel.queue_declare(queue="Arima", durable=False)
    channel.basic_publish(exchange="", 
        routing_key="Arima",
        body=json.dumps(alert),
        properties=pika.BasicProperties(
        )
    )
    print(f"message: {alert.message}, published succesfully!")
    conn.close()
    return


# create the api and declare http route
@app.post("/create/alert")
async def Create_alert(alert: schemas.Alert):
    empty = None
    if alert is empty:
        return {"error": "cannot process empty alert payload"}
    Publish_alert(alert=alert)
    return {"message": "your alert will be published shortly"}
    