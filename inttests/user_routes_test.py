# Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

import requests
from base import BaseTest

class UserRoutesTest(BaseTest):

    def setUp(self):
        super(UserRoutesTest, self).setUp()
        self.url = f"{self.base}/users"

    def tearDown(self):
        super(UserRoutesTest, self).tearDown()