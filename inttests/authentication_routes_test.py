# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, json
from base import TestBase

from dateutil import parser


class TestAuthenticationRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestAuthenticationRoutes, self).setup_method(test_method)

    def teardown_method(self, test_method):
        super(TestAuthenticationRoutes, self).teardown_method(test_method)

    def test_last_login_gets_updated(self):
        url = f'{self.base}/users/4'

        cookies = self.get_authn_cookies_user_normal()
        r = requests.get(url, headers=self.read_headers, cookies=cookies)
        assert 200 == r.status_code
        first_login = parser.parse(r.json()['LastLogin'])

        cookies = self.get_authn_cookies_user_normal()
        r = requests.get(url, headers=self.read_headers, cookies=cookies)
        assert 200 == r.status_code
        second_login = parser.parse(r.json()['LastLogin'])

        assert first_login < second_login
