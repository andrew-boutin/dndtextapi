# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import BaseTest

class AdminRoutesTest(BaseTest):

    def setUp(self):
        super(AdminRoutesTest, self).setUp()
        self.url = f"{self.base}/admin"

    def tearDown(self):
        super(AdminRoutesTest, self).tearDown()