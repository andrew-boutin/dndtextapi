# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, json, string, pytest
from random import choice, randint
from mockserver import MockServerClient, request, response

class TestBase:

    def setup_method(self, test_method):
        self.base = "http://app:8080"
        self.client = MockServerClient("http://mockserver:1080")
        self.client.reset()

        self.read_headers = {"Accept": "application/json"}
        self.read_write_headers = {"Content-Type": "application/json", **self.read_headers}

    def teardown_method(self, test_method):
        self.client.verify()

    def get_authn_cookies_user_normal(self):
        return self.get_authn_cookies('regularuser@fake.com')

    def get_authn_cookies_user_admin(self):
        return self.get_authn_cookies('adminuser@fake.com')

    def get_authn_cookies_user_new(self):
        allchar = string.ascii_letters + string.digits
        random_name = "".join(choice(allchar) for x in range(8))
        return self.get_authn_cookies(f'{random_name}user@fake.com')

    def get_authn_cookies(self, email):
        self.client.stub(
            request(method="GET", path="/o/oauth2/auth"),
            response(code=200, body="fake google auth")
        )

        url = f"{self.base}/login"
        r = requests.get(url)
        assert r.status_code == 200

        # Retrieve the state that our app sent to the mock server when it redirected
        # to the google oauth2 auth endpoint.
        # http://www.mock-server.com/mock_server/mockserver_clients.html#rest-api shows
        # the mock-server endpoint that can retrieve requests it received.
        # Requested this feature be added to the mock-server python client here
        # https://github.com/internap/python-mockserver-client/issues/16.
        data = json.dumps({"path": "/o/oauth2/auth", "method": "GET"})
        r = requests.put("http://mockserver:1080/retrieve?type=REQUESTS", data=data)
        assert 200 == r.status_code
        state = r.json()[-1]["queryStringParameters"]["state"][0]
        assert "" != state

        # Mock out the app attempting to get an access token using the state and code from Google
        tokenJson = json.dumps({
            "access_token": "myfakeaccesstoken",
            "token_type": "Bearer",
            "expires_in": 0,
        })
        self.client.stub(
            request(method="POST", path="/o/oauth2/token"),
            response(code=200, body=tokenJson, headers={"Content-Type": "application/json"})
        )

        # Mock out the app attempting to get profile data using the access token
        data = json.dumps({"email": email})
        self.client.stub(
            request(method="GET", path="/oauth2/v2/userinfo"),
            response(code=200, body=data)
        )

        # Perform the callback from Google to finish the authentication with the app, return the cookie
        code = "supersecretcode"
        r = requests.get(f"{self.base}/callback?state={state}&code={code}")
        assert 204 == r.status_code
        cookie = r.cookies['dndtextapisession']
        return dict(dndtextapisession=cookie)

    # TODO: Teardown delete the channel
    @pytest.fixture()
    def create_channel_normal_user(self): # TODO: Input variable for what user to create the channel for
        cookies = self.get_authn_cookies_user_normal()

        data = json.dumps({
            "OwnerID": 4,
            "DMID": 4,
            "Name": "test channel " + str(randint(0, 1000000)),
            "IsPrivate": False
        })

        r = requests.post(f'{self.base}/channels', data=data, headers=self.read_write_headers, cookies=cookies)
        assert 201 == r.status_code
        return r.json()
