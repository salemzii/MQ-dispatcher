from email import message
from fastapi import Body
from pydantic import BaseModel


# declare the schema
class Alert(BaseModel):
    topic: str
    message: str
    #platform: list[str] = []