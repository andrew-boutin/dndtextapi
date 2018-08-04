# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import TestBase

class TestPublicRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestPublicRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/public"

    def teardown_method(self, test_method):
        super(TestPublicRoutes, self).teardown_method(test_method)

    def test_get_channels(self):
        """Test the public get channels route."""
        # Make sure a bogus accept header is rejected
        r = requests.get(f"{self.url}/channels", headers={"Accept": "Nothing"})
        assert r.status_code == 400

        # Add the correct header and make sure the call is ok
        r = requests.get(f"{self.url}/channels", headers=self.read_headers)
        assert r.status_code == 200
        # TODO: assert the contents

    def test_get_channel(self):
        """Test the public get channel route."""
        # Make sure using a bad accept header isn't allowed
        url = f"{self.url}/channels/"
        r = requests.get(f"{url}1", headers={"Accept": "Nothing"})
        assert r.status_code == 400

        # Add the header and make sure the call is ok
        r = requests.get(f"{url}1", headers=self.read_headers)
        assert r.status_code == 200
        # TODO: assert the contents

        # Verify we get a not found on a made up id
        r = requests.get(f"{url}999", headers=self.read_headers)
        assert r.status_code == 404

        # Make sure we can't get a private channel
        url = f"{url}2"
        r = requests.get(url, headers=self.read_headers)
        assert r.status_code == 401

    def test_get_messages_from_public_channel(self):
        """Test the public get messages route for getting messages from a channel."""
        messages_url = f"{self.url}/channels/%d/messages"

        # Use a private channel id verify denied
        url = messages_url % 2
        r = requests.get(url, headers=self.read_headers)
        assert r.status_code == 401

        # Make up a channel id verify not found
        url = messages_url % 999
        r = requests.get(url, headers=self.read_headers)
        assert r.status_code == 404

        # Make sure not including the accept header isn't allowed
        url = messages_url % 1
        r = requests.get(url)
        assert r.status_code == 400

        # Use a public channel for a valid request
        url = messages_url % 1
        r = requests.get(url, headers=self.read_headers)
        assert r.status_code == 200
        # TODO: assert the contents - no meta messages