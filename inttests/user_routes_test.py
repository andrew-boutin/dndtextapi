# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import TestBase

class TestUserRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestUserRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/users"

    def teardown_method(self, test_method):
        super(TestUserRoutes, self).teardown_method(test_method)