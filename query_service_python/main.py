from fastapi import FastAPI
from fastapi_sqlalchemy import DBSessionMiddleware, db
from sqlalchemy import Column, Integer, String, Enum
from sqlalchemy.ext.declarative import as_declarative, declared_attr

import os
import logging

log = logging.getLogger("query_service")

@as_declarative()
class Base:
    @declared_attr
    def __tablename__(cls):
        return cls.__name__.lower()

class Person(Base):
    __tablename__ = "person"

    person_id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    first_name = Column(String(100), index=True)
    last_name = Column(String(100), index=True)
    sex = Column(Enum("male", "female", "other"), default="other")


app = FastAPI()

PW = os.popen('cat /run/secrets/db_root_password').read().replace("\n", "")
DB = os.getenv("MYSQL_DATABASE")
DATABASE_URL = f"mysql+mysqlconnector://root:{PW}@mysql:3306/{DB}"

app.add_middleware(DBSessionMiddleware, db_url=DATABASE_URL)



@app.get("/api/person")
def read_root():
    return {"welcome:": "to the person API"}

@app.get("/api/users")
def get_users(skip: int = 0, limit: int = 100):
    try:
        users = db.session.query(Person).offset(skip).limit(limit).all()
    except Exception as e:
        log.warning(e)

    for user in users:
        log.warning(user)

    return users