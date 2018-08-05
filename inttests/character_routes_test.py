# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests, pytest
from base import TestBase

# TODO: Update, Delete, Create, Get - requires session fix
class TestCharacterRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestCharacterRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/channels/%s/characters/" # TODO: Multi level path may not work w/ sessions

    def teardown_method(self, test_method):
        super(TestCharacterRoutes, self).teardown_method(test_method)
    
    def test_user_has_to_be_channel_owner_to_create_character(self):
        # TODO: Implement
        pass

    @pytest.mark.skip(reason="skip until cookie/session issue w/ nested paths is resolved")
    def test_can_only_create_one_character_per_user_per_channel(self,
                                                                create_channel_normal_user):
        cookies = self.get_authn_cookies_user_normal() # TODO: Move to fixture that returns cookies

        # Fixture creates a channel for the test and allows us to access the json result
        channel_id = create_channel_normal_user['id']

        # Create a character and expect OK
        url = self.url % channel_id
        data = json.dumps({}) # TODO: Fill in data
        r = requests.get(self.url, headers=self.read_headers, cookies=cookies)
        assert 200 == r.status_code # TODO: Use http codes

        # This time expect an error related to already having a character
        r = requests.get(self.url, headers=self.read_headers, cookies=cookies)
        assert 500 == r.status_code # TODO: Pass a 400 back for this scenario
