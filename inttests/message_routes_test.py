# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import TestBase

# TODO: Update, Delete, Create, Get - requires session fix
class TestMessageRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestMessageRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/messages"

    def teardown_method(self, test_method):
        super(TestMessageRoutes, self).teardown_method(test_method)
