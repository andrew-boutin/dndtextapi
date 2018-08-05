# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, json
from base import TestBase

# TODO: Created, Delete
class TestChannelRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestChannelRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/channels"

    def teardown_method(self, test_method):
        super(TestChannelRoutes, self).teardown_method(test_method)

    def test_get_channels(self):
        cookies = self.get_authn_cookies_user_normal()

        r = requests.get(self.url, headers=self.read_headers, cookies=cookies)
        assert 200 == r.status_code
        # TODO: At least one, verify fields, etc.
        # [{"Name":"my public channel","Description":"my public channel description","Topic":"some topic","ID":1,"OwnerID":1,"IsPrivate":false,"CreatedOn":"2018-07-24T20:02:49.089425Z","LastUpdated":"2018-07-24T20:02:49.089425Z","DMID":1},{"Name":"my private channel","Description":"my private channel description","Topic":"","ID":2,"OwnerID":1,"IsPrivate":true,"CreatedOn":"2018-07-24T20:02:49.089425Z","LastUpdated":"2018-07-24T20:02:49.089425Z","DMID":1}]

    def test_update_channel(self,
                            create_channel_normal_user):
        cookies = self.get_authn_cookies_user_normal()

        channel_id = create_channel_normal_user['ID']
        url = self.url + f'/{channel_id}'

        data = create_channel_normal_user
        data['Topic'] = 'Updated topic'
        data = json.dumps(data)

        r = requests.put(url, data=data, headers=self.read_write_headers, cookies=cookies)
        assert 200 == r.status_code
        assert 'Updated topic' == r.json()['Topic']

    def test_get_channel(self,
                         create_channel_normal_user):
        cookies = self.get_authn_cookies_user_normal()

        channel_id = create_channel_normal_user['ID']
        url = self.url + f'/{channel_id}'

        r = requests.get(url, cookies=cookies)
        assert 200 == r.status_code
        assert channel_id == r.json()['ID']
