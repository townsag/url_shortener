from locust import (
    HttpUser,
    task,
    between,
)
from requests import Response



class SimpleUser(HttpUser):
    wait_time = between(0, 1)
    host = "http://localhost:8000"

    @task
    def create_access_mapping(self):
        create_resp: Response = self.client.post(
            url=f"/api/mapping",
            json={
                "longUrl":"http://google.com",
            }
        )
        short_url_id = create_resp.json()["shortUrl"]
        access_response: Response = self.client.get(
            url=f"/api/{short_url_id}",
            name="/[shortUrl]",
            allow_redirects=False,
            # ^this was interesting to debug
            # without disallowing redirects we were sending 20,000 rpm to google.com 
            # and being rate limited by google
        )