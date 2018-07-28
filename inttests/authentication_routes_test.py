# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, json
from base import BaseTest

from dateutil import parser


class AuthenticationRoutesTest(BaseTest):

    def setUp(self):
        super(AuthenticationRoutesTest, self).setUp()

    def tearDown(self):
        super(AuthenticationRoutesTest, self).tearDown()

    def test_last_login_gets_updated(self):
        url = f'{self.base}/users/4'

        cookies = self.get_authn_cookies_user_normal()
        r = requests.get(url, headers=self.read_headers, cookies=cookies)
        self.assertEqual(200, r.status_code)
        first_login = parser.parse(r.json()['LastLogin'])

        cookies = self.get_authn_cookies_user_normal()
        r = requests.get(url, headers=self.read_headers, cookies=cookies)
        self.assertEqual(200, r.status_code)
        second_login = parser.parse(r.json()['LastLogin'])

        self.assertTrue(first_login < second_login)
