# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, json
from base import TestBase

# TODO: Delete, Create
class TestUserRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestUserRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/users/{4}"

    def teardown_method(self, test_method):
        super(TestUserRoutes, self).teardown_method(test_method)

    def test_get_user(self):
        cookies = self.get_authn_cookies_user_normal()

        r = requests.get(self.url, cookies=cookies)
        assert 200 == r.status_code
        assert 'regularuser@fake.com' == r.json()['Email']

    def test_update_user(self):
        cookies = self.get_authn_cookies_user_normal()

        data = json.dumps({
            'Username': 'regularuser',
            'Email': 'regularuser@fake.com',
            'IsAdmin': False,
            'IsBanned': False,
            'Bio': 'updated bio'
        })

        r = requests.put(self.url, data=data, cookies=cookies, headers=self.read_write_headers)
        assert 200 == r.status_code
        assert 'updated bio' == r.json()['Bio']
