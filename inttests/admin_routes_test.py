# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import TestBase

class TestAdminRoutes(TestBase):

    def setup_method(self, test_method):
        super(TestAdminRoutes, self).setup_method(test_method)
        self.url = f"{self.base}/admin"

    def teardown_method(self, test_method):
        super(TestAdminRoutes, self).teardown_method(test_method)